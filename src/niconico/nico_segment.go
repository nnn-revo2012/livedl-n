package niconico

import (
	//"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	pb "github.com/nnn-revo2012/livedl/proto"

    "google.golang.org/protobuf/proto"
)

type SegmentServer struct {
	uri                  string
	headers              map[string]string
	isDisconnect         bool
	unexpectedDisconnect bool
	stream               *BinaryStream
	//buffer               *bytes.Buffer
	cancellationCtx      context.Context
	cancelFunc           context.CancelFunc
	//onNetworkError       func() error
	servername           string
	mu                   sync.Mutex
	message              chan<- *pb.ChunkedMessage
}

func NewSegmentServer(uri, servername string, message chan<- *pb.ChunkedMessage) *SegmentServer {
	ctx, cancel := context.WithCancel(context.Background())
	headers := map[string]string{}

	return &SegmentServer{
		uri:                  uri,
		headers:              headers,
		isDisconnect:         false,
		unexpectedDisconnect: false,
		stream:               NewBinaryStream(),
		//buffer:             new(bytes.Buffer),
		cancellationCtx:      ctx,
		cancelFunc:           cancel,
		servername:           servername,
		message:              message,
	}
}

func (ssc *SegmentServer) Connect() error {
	client := &http.Client{
		Timeout: 45 * time.Second,
	}

	req, err := http.NewRequestWithContext(ssc.cancellationCtx, http.MethodGet, ssc.uri, nil)
	if err != nil {
		return err
	}

	for key, value := range ssc.headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		ssc.mu.Lock()
		defer ssc.mu.Unlock()
		ssc.unexpectedDisconnect = true
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buffer := make([]byte, 8192)
	for {
		select {
		case <-ssc.cancellationCtx.Done():
			fmt.Println("Read operation was canceled due to a timeout or external cancellation.")
			return nil
		default:
			//resp.Body.SetReadDeadline(time.Now().Add(45 * time.Second))
			n, err := resp.Body.Read(buffer)
			if err != nil {
				if err == io.EOF {
					log.Println(ssc.servername+"Read EOF.")
					ssc.isDisconnect = true
					return nil
				}
				ssc.mu.Lock()
				defer ssc.mu.Unlock()
				ssc.unexpectedDisconnect = true
				return err
			}

			if n > 0 {
				dataChunk := make([]byte, n)
				copy(dataChunk, buffer[:n])
				//ssc.buffer.Write(dataChunk)
				err = ssc.segmentData(dataChunk)
				if err != nil {
					return err
				}
			}
		}
	}

	if err != nil {
		return err
	}

	if ssc.IsUnexpectedDisconnect() {
		return segmentNetworkError()
	}

	return nil
}

//func (ssc *SegmentServer) getBuffers() []byte {
//	return ssc.buffer.Bytes()
//}

func (ssc *SegmentServer) stopReceiving() {
	ssc.cancelFunc()
}

func (ssc *SegmentServer) IsUnexpectedDisconnect() bool {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()
	return ssc.unexpectedDisconnect
}

func (ssc *SegmentServer) IsDisconnect() bool {
	return ssc.isDisconnect
}

func (ssc *SegmentServer) Disconnect() bool {
	ssc.stopReceiving()
	ssc.isDisconnect = true
	log.Println("disconnect "+ssc.servername+" server.")
	return true
}

func (ssc *SegmentServer) segmentData(data []byte) error {
	//log.Printf(ssc.servername+" received %d bytes.\n", len(data))

	ssc.stream.AddBuffer(data)

	items := ssc.stream.Read()

	for item := range items {
		message := &pb.ChunkedMessage{}
		if err := proto.Unmarshal(item, message); err != nil {
			fmt.Println(err)
			continue
		}
		//fmt.Println(message)
		if len(message.String()) > 0 {
			ssc.message <- message
		}
	}

	ssc.stream.ClearBuffer()

	return nil
}

func segmentNetworkError() error {
	log.Println("segment(or previous)Server Network error")
	return nil
}
