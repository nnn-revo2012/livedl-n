package niconico

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type StreamReceiver[T any] struct {
	cancellationCtx    context.Context
	cancelFunc         context.CancelFunc
	buffer             *bytes.Buffer
	processData        func([]byte, *T, *NicoHls) error
	unexpectedDisconnect bool
	mu                 sync.Mutex
}

func NewStreamReceiver[T any] (processData func([]byte, *T, *NicoHls) error) *StreamReceiver[T] {
	ctx, cancel := context.WithCancel(context.Background())

	return &StreamReceiver[T]{
		cancellationCtx:    ctx,
		cancelFunc:         cancel,
		buffer:             new(bytes.Buffer),
		processData:        processData,
		unexpectedDisconnect: false,
	}
}

func (sr *StreamReceiver[T]) Receive (url string, headers map[string]string, msc *T, hls *NicoHls) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(sr.cancellationCtx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		sr.mu.Lock()
		defer sr.mu.Unlock()
		sr.unexpectedDisconnect = true
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buffer := make([]byte, 8192)
	for {
		select {
		case <-sr.cancellationCtx.Done():
			fmt.Println("Read operation was canceled due to a timeout or external cancellation.")
			return nil
		default:
			//resp.Body.SetReadDeadline(time.Now().Add(30 * time.Second))
			n, err := resp.Body.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return nil
				}
				sr.mu.Lock()
				defer sr.mu.Unlock()
				sr.unexpectedDisconnect = true
				return err
			}

			if n > 0 {
				dataChunk := make([]byte, n)
				copy(dataChunk, buffer[:n])
				sr.buffer.Write(dataChunk)

				//err = sr.processData(dataChunk, stream)
				err = sr.processData(dataChunk, msc, hls)
				if err != nil {
					return err
				}
			}
		}
	}
}

func (sr *StreamReceiver[T]) GetBuffers() []byte {
	return sr.buffer.Bytes()
}

func (sr *StreamReceiver[T]) StopReceiving() {
	sr.cancelFunc()
}

func (sr *StreamReceiver[T]) UnexpectedDisconnect() bool {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	return sr.unexpectedDisconnect
}
