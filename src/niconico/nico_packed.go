package niconico

import (
	"bytes"
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

type PackedSegment struct {
	uri                  string
	headers              map[string]string
	isDisconnect         bool
	unexpectedDisconnect bool
	//stream               *BinaryStream
	buffer               *bytes.Buffer
	cancellationCtx      context.Context
	cancelFunc           context.CancelFunc
	mu                   sync.Mutex
	segment              chan<- *pb.PackedSegment
}

func NewPackedSegment(uri string, segment chan<- *pb.PackedSegment) *PackedSegment {
	ctx, cancel := context.WithCancel(context.Background())
	headers := map[string]string{}

	return &PackedSegment{
		uri:                  uri,
		headers:              headers,
		isDisconnect:         false,
		unexpectedDisconnect: false,
		//stream:               NewBinaryStream(),
		buffer:               new(bytes.Buffer),
		cancellationCtx:      ctx,
		cancelFunc:           cancel,
		segment:              segment,
	}
}

func (ps *PackedSegment) Connect() error {
	client := &http.Client{
		Timeout: 45 * time.Second,
	}

	messageUri := ps.uri
	req, err := http.NewRequestWithContext(ps.cancellationCtx, http.MethodGet, messageUri, nil)
	if err != nil {
		return err
	}

	for key, value := range ps.headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		ps.mu.Lock()
		defer ps.mu.Unlock()
		ps.unexpectedDisconnect = true
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buffer := make([]byte, 8192)
ExitPacked:
	for {
		select {
		case <-ps.cancellationCtx.Done():
			fmt.Println("Read operation was canceled due to a timeout or external cancellation.")
			return nil
		default:
			//resp.Body.SetReadDeadline(time.Now().Add(45 * time.Second))
			n, err := resp.Body.Read(buffer)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Read EOF.")
					//return nil
					break ExitPacked
				}
				ps.mu.Lock()
				defer ps.mu.Unlock()
				ps.unexpectedDisconnect = true
				return err
			}

			if n > 0 {
				dataChunk := make([]byte, n)
				copy(dataChunk, buffer[:n])
				ps.buffer.Write(dataChunk)
				err = ps.packedData(dataChunk)
				if err != nil {
					return err
				}
			}
		}
	}

	if err != nil {
		return err
	}
	ps.isDisconnect = true

	//if ps.IsUnexpectedDisconnect() {
		//if err := ps.onNetworkError(); err != nil {
		//	return err
		//}
		//return nil
	//}
	buf := ps.getBuffers()
	segment := &pb.PackedSegment{}
	if err := proto.Unmarshal(buf, segment); err != nil {
		return err
	}
	ps.segment <- segment

	return nil
}

func (ps *PackedSegment) getBuffers() []byte {
	return ps.buffer.Bytes()
}

func (ps *PackedSegment) stopReceiving() {
	ps.cancelFunc()
}

func (ps *PackedSegment) IsUnexpectedDisconnect() bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.unexpectedDisconnect
}

func (ps *PackedSegment) IsDisconnect() bool {
	return ps.isDisconnect
}

func (ps *PackedSegment) Disconnect() bool {
	ps.stopReceiving()
	ps.isDisconnect = true
	log.Println("disconnect packedSegment server.")
	return true
}

func (ps *PackedSegment) packedData(data []byte) error {
	log.Printf("packedSegment received %d bytes.\n", len(data))
	return nil
}

func (ps *PackedSegment) GetNextUri() string {
	return ps.uri
}

func (ps *PackedSegment) SetNextUri(nexturi string) error {
	if len(nexturi) > 0 {
		ps.uri = nexturi
	}
	return nil
}

func packedNetworkError() error {
	log.Println("PackedSegment Network error")
	return nil
}
