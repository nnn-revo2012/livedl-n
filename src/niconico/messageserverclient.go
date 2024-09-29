package niconico

import (
	//"fmt"
	"log"
	"sync"
	//"regexp"

	pb "github.com/nnn-revo2012/livedl/proto"

    "google.golang.org/protobuf/proto"
)

type MessageServerClient struct {
	nextStreamAt        string
	processData         func(*pb.ChunkedEntry, *MessageServerClient) error
	streamReceiver      *StreamReceiver[MessageServerClient]
	uri                 string
	headers             map[string]string
	isDisconnect        bool
	stream              *BinaryStream
	//hls                 *NicoHls
	onNetworkError      func() error
	beforeNextStreamAt  string
	mu                  sync.Mutex
	entry               chan<- *pb.ChunkedEntry
}

func NewMessageServerClient(uri string, processData func(*pb.ChunkedEntry, *MessageServerClient) error, onNetworkError func() error, entry chan<- *pb.ChunkedEntry) *MessageServerClient {
	headers := map[string]string{
		"header": "u=1, i",
	}

	return &MessageServerClient{
		nextStreamAt:       "now",
		processData:        processData,
		uri:                uri,
		headers:            headers,
		stream:             NewBinaryStream(),
		streamReceiver:     NewStreamReceiver(ProcessRawData),
		//hls:                hls,
		onNetworkError:     onNetworkError,
		beforeNextStreamAt: "",
		entry:              entry,
	}
}

func (msc *MessageServerClient) DoConnect() error {
	//fmt.Println("msc.DoConnect")
	//for !msc.isDisconnect && !msc.IsUnexpectedDisconnect() {
		err := msc.streamReceiver.Receive(msc.uri + "?at=" + msc.nextStreamAt, msc.headers, msc)
		if err != nil {
			return err
		}

		msc.mu.Lock()
		if msc.beforeNextStreamAt == msc.nextStreamAt {
			msc.mu.Unlock()
			return msc.onNetworkError()
		}
		msc.beforeNextStreamAt = msc.nextStreamAt
		msc.mu.Unlock()
	//}

	if msc.IsUnexpectedDisconnect() {
		return msc.onNetworkError()
	}

	return nil
}

func (msc *MessageServerClient) IsUnexpectedDisconnect() bool {
	return msc.streamReceiver.UnexpectedDisconnect()
}

func (msc *MessageServerClient) IsDisconnect() bool {
	return msc.isDisconnect
}

func (msc *MessageServerClient) Disconnect() bool {
	msc.streamReceiver.StopReceiving()
	msc.isDisconnect = true
	log.Println("disconnect message server.")
	return true
}

func (msc *MessageServerClient) GetNextStreamAt() string {
	return msc.nextStreamAt
}

func (msc *MessageServerClient) SetNextStreamAt(nextat string) error {
	if len(nextat) > 0 {
		msc.nextStreamAt = nextat
	}
	return nil
}

func ProcessRawData(data []byte, msc *MessageServerClient) error {
	//fmt.Println("ProcessRawData")
	log.Printf("message received %d bytes.\n", len(data))

	msc.stream.AddBuffer(data)

	items := msc.stream.Read()

	for item := range items {
		entry := &pb.ChunkedEntry{}
		if err := proto.Unmarshal(item, entry); err != nil {
			return err
		}
		//fmt.Println(entry)
		if len(entry.String()) <= 0 {
			return nil
		}
		msc.entry <- entry
		if err := msc.processData(entry, msc); err != nil {
			return err
		}
	}

	msc.stream.CheckClearBuffer()

	return nil
}

func ProcessMessageData(entry *pb.ChunkedEntry, msc *MessageServerClient) error {
	return nil
}

func NetworkError() error {
	log.Println("messageServer Network error")

	return nil
}
