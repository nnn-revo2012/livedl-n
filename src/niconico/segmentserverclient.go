package niconico

import (
	"fmt"
	"log"
	"sync"
	"regexp"

	pb "github.com/nnn-revo2012/livedl/proto"

    "google.golang.org/protobuf/encoding/protojson"
    "google.golang.org/protobuf/proto"
)

type SegmentServerClient struct {
	uri                    string
	headers                map[string]string
	streamReceiver         *StreamReceiver[SegmentServerClient]
	isDisconnect           bool
	stream                 *BinaryStream
	processData            func(*pb.ChunkedMessage, bool, *SegmentServerClient, *NicoHls) error
	isInitialCommentsReceiving bool
	hls                    *NicoHls
	onNetworkError         func() error
	mu                     sync.Mutex // For thread-safety
}

// NewSegmentServerClient コンストラクタ
func NewSegmentServerClient(uri string, processData func(*pb.ChunkedMessage, bool, *SegmentServerClient, *NicoHls) error, hls *NicoHls, onNetworkError func() error, isInitialCommentsReceiving bool) *SegmentServerClient {
	headers := make(map[string]string)
	client := &SegmentServerClient{
		uri:                    uri,
		headers:                headers,
		streamReceiver:         NewStreamReceiver(SegmentRawData),
		stream:                 NewBinaryStream(),
		processData:            processData,
		isInitialCommentsReceiving: isInitialCommentsReceiving,
		hls:                    hls,
		onNetworkError:         onNetworkError,
		isDisconnect:           false,
	}
	return client
}

func (s *SegmentServerClient) DoConnect() error {
	//fmt.Println("s.DoConnect")
	err := s.streamReceiver.Receive(s.uri, s.headers, s, s.hls)
	if err != nil {
		return err
	}
	if s.streamReceiver.UnexpectedDisconnect() {
		if err := s.onNetworkError(); err != nil {
			return err
		}
	}
	s.isDisconnect = true

	// GCしやすいように子クラスの参照を消す
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stream = nil
	s.processData = nil
	s.streamReceiver = nil

	return nil
}

func (s *SegmentServerClient) IsUnexpectedDisconnect() bool {
	return s.streamReceiver.UnexpectedDisconnect()
}

func (s *SegmentServerClient) IsDisconnect() bool {
	return s.isDisconnect
}

func (s *SegmentServerClient) Disconnect() bool {
	if s.streamReceiver != nil {
		s.streamReceiver.StopReceiving()
		s.isDisconnect = true
		log.Println("disconnect segment server.")
	}
	return true
}

func SegmentRawData(data []byte, s *SegmentServerClient, hls *NicoHls) error {
	if s.streamReceiver == nil {
		return nil
	}
	//fmt.Println("SegmentRawData()")
	//fmt.Printf("segment received %d bytes.\n", len(data))

	s.mu.Lock()
	defer s.mu.Unlock()

	s.stream.AddBuffer(data)

	for item := range s.stream.Read() {
		entry := &pb.ChunkedMessage{}
		if err := proto.Unmarshal(item, entry); err != nil {
			return err
		}
		//fmt.Println(entry)	//DEBUG
		if err := s.processData(entry, s.isInitialCommentsReceiving, s, s.hls); err != nil {
			return err
		}
	}

	s.stream.CheckClearBuffer()

	return nil
}

func ProcessSegmentData(entry *pb.ChunkedMessage, inicomment bool, s *SegmentServerClient, hls *NicoHls) error {
	//fmt.Println("ProcessSegmentData()")
	str := entry.String()
	if len(str) <= 0 {
		return nil
	}
	//fmt.Println(str)
	var e string
	if ma := regexp.MustCompile(`^([\w]+):`).FindStringSubmatch(str); len(ma) > 0 {
		e = ma[1]
	}
	switch e {
	case "signal":
		fmt.Println(str)
	case "meta":
		//fmt.Println(str)
		fmt.Println(str)
	default:
		fmt.Println("Unknown entry: "+str)
	}

	return nil
}

func SegmentNetworkError() error {
	log.Println("segmentServer Network error")

	return nil
}

