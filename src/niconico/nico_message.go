package niconico

import (
	//"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/nnn-revo2012/livedl/httpbase"
	pb "github.com/nnn-revo2012/livedl/proto"

    "google.golang.org/protobuf/proto"
)

type MessageServer struct {
	uri                  string
	headers              map[string]string
	nextStreamAt         string
	beforeNextStreamAt   string
	isDisconnect         bool
	unexpectedDisconnect bool
	stream               *BinaryStream
	//buffer               *bytes.Buffer
	cancellationCtx      context.Context
	cancelFunc           context.CancelFunc
	//onNetworkError       func() error
	mu                   sync.Mutex
	entry                chan<- *pb.ChunkedEntry
}

var ClientMessage = &http.Client{}

func setMessageProxy(rawurl string) (err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return
	}
	ClientMessage.Transport.(*http.Transport).Proxy = http.ProxyURL(u)
	return
}

func NewMessageServer(uri, proxy string, entry chan<- *pb.ChunkedEntry) *MessageServer {
	ctx, cancel := context.WithCancel(context.Background())
	headers := map[string]string{
		"header": "u=1, i",
	}

	//dt := http.DefaultTransport.(*http.Transport).Clone()
	dt := httpbase.ClonedTransport(http.DefaultTransport)
	dt.MaxIdleConnsPerHost = 10

	if proxy != "" {
		setMessageProxy(proxy)
	}

	ClientMessage = &http.Client{
		Transport: dt,
		Timeout:   45 * time.Second,
	}

	return &MessageServer{
		uri:                  uri,
		headers:              headers,
		isDisconnect:         false,
		unexpectedDisconnect: false,
		stream:               NewBinaryStream(),
		//buffer:             new(bytes.Buffer),
		cancellationCtx:      ctx,
		cancelFunc:           cancel,
		nextStreamAt:         "now",
		beforeNextStreamAt:   "",
		entry:                entry,
	}
}

func (msc *MessageServer) Connect() error {
	messageUri := msc.uri + "?at=" + msc.nextStreamAt
	req, err := http.NewRequestWithContext(msc.cancellationCtx, http.MethodGet, messageUri, nil)
	if err != nil {
		return err
	}

	for key, value := range msc.headers {
		req.Header.Set(key, value)
	}

	resp, err := ClientMessage.Do(req)
	if err != nil {
		msc.mu.Lock()
		defer msc.mu.Unlock()
		msc.unexpectedDisconnect = true
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buffer := make([]byte, 8192)
	for {
		select {
		case <-msc.cancellationCtx.Done():
			fmt.Println("Read operation was canceled due to a timeout or external cancellation.")
			return nil
		default:
			//resp.Body.SetReadDeadline(time.Now().Add(45 * time.Second))
			n, err := resp.Body.Read(buffer)
			if err != nil {
				if err == io.EOF {
					//fmt.Println("Read EOF.")
					return nil
				}
				msc.mu.Lock()
				defer msc.mu.Unlock()
				msc.unexpectedDisconnect = true
				return err
			}

			if n > 0 {
				dataChunk := make([]byte, n)
				copy(dataChunk, buffer[:n])
				//msc.buffer.Write(dataChunk)
				err = msc.messageData(dataChunk)
				if err != nil {
					return err
				}
			}
		}
	}

	if err != nil {
		return err
	}

	if msc.IsUnexpectedDisconnect() {
		return messageNetworkError()
	}

	return nil
}

//func (msc *MessageServer) getBuffers() []byte {
//	return msc.buffer.Bytes()
//}

func (msc *MessageServer) stopReceiving() {
	msc.cancelFunc()
}

func (msc *MessageServer) IsUnexpectedDisconnect() bool {
	msc.mu.Lock()
	defer msc.mu.Unlock()
	return msc.unexpectedDisconnect
}

func (msc *MessageServer) IsDisconnect() bool {
	return msc.isDisconnect
}

func (msc *MessageServer) Disconnect() bool {
	msc.stopReceiving()
	msc.isDisconnect = true
	log.Println("disconnect message server.")
	return true
}

func (msc *MessageServer) GetNextStreamAt() string {
	return msc.nextStreamAt
}

func (msc *MessageServer) SetNextStreamAt(nextat string) error {
	if len(nextat) > 0 {
		msc.nextStreamAt = nextat
	}
	return nil
}

func (msc *MessageServer) messageData(data []byte) error {
	//log.Printf("message received %d bytes.\n", len(data))

	msc.stream.AddBuffer(msc.stream, data)

	items := msc.stream.Read(msc.stream)

	for item := range items {
		entry := &pb.ChunkedEntry{}
		if err := proto.Unmarshal(item, entry); err != nil {
			fmt.Println(err)
			continue
		}
		//fmt.Println(entry)
		if len(entry.String()) > 0  && !msc.isDisconnect {
			msc.entry <- entry
		}
	}

	msc.stream.ClearBuffer(msc.stream)

	return nil
}

func messageNetworkError() error {
	log.Println("messageServer Network error")
	return nil
}
