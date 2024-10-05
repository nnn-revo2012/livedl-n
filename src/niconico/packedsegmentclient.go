package niconico

import (
	"fmt"
	"log"
	//"sync"
	//"regexp"

	pb "github.com/nnn-revo2012/livedl/proto"

    "google.golang.org/protobuf/proto"
)

type PackedSegmentClient struct {
	uri             string
	headers         map[string]string
	onDisconnect    func(*pb.PackedSegment, *PackedSegmentClient) error
	streamReceiver  *StreamReceiver[PackedSegmentClient]
	isDisconnect    bool
	onNetworkError  func() error
	segment         chan<- *pb.PackedSegment
}

// NewPackedSegmentClient コンストラクタ
func NewPackedSegmentClient(uri string, onDisconnect func(*pb.PackedSegment, *PackedSegmentClient) error, onNetworkError func() error, segment chan<- *pb.PackedSegment) *PackedSegmentClient {
	headers := make(map[string]string)
	client := &PackedSegmentClient{
		uri:            uri,
		headers:        headers,
		onDisconnect:   onDisconnect,
		streamReceiver: NewStreamReceiver(PackedRawData),
		onNetworkError: onNetworkError,
		isDisconnect:   false,
		segment:        segment,
	}
	return client
}

func (p *PackedSegmentClient) DoConnect() error {
	//fmt.Println("p.DoConnect")
	err := p.streamReceiver.Receive(p.uri, p.headers, p)
	if err != nil {
		return err
	}
	p.isDisconnect = true

	if p.streamReceiver.UnexpectedDisconnect() {
		if err := p.onNetworkError(); err != nil {
			return err
		}
		return nil
	}

	// PackedSegmentの解析とonDisconnectコールバックの呼び出し
	buffers := p.streamReceiver.GetBuffers()
	segment := &pb.PackedSegment{}
	if err := proto.Unmarshal(buffers, segment); err != nil {
		return err
	}
	//if err := p.onDisconnect(segment); err != nil {
	//	return err
	//}
	p.segment <- segment

	return nil
}

// PackedRawData handles the received data for the PackedSegmentClient
func PackedRawData(data []byte, p *PackedSegmentClient) error {
	fmt.Printf("PackedSegmentClient: Received %d bytes.\n", len(data))
	return nil
}

func (p *PackedSegmentClient) GetNextUri() string {
	return p.uri
}

func (p *PackedSegmentClient) SetNextUri(nexturi string) error {
	if len(nexturi) > 0 {
		p.uri = nexturi
	}
	return nil
}

func PackedDisconnect(packed *pb.PackedSegment, p *PackedSegmentClient) error {
	if p.streamReceiver != nil {
		p.streamReceiver.StopReceiving()
		p.isDisconnect = true
		log.Println("disconnect packedsegment server.")
	}
	return nil
}

func (p *PackedSegmentClient) Disconnect() error {
	if p.streamReceiver != nil {
		p.streamReceiver.StopReceiving()
		p.isDisconnect = true
		log.Println("disconnect packedsegment server.")
	}
	return nil
}

func PackedNetworkError() error {
	log.Println("packedSegmentServer Network error")

	return nil
}
