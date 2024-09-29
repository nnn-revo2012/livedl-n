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
	p.Disconnect()

	return nil
}

// PackedRawData handles the received data for the PackedSegmentClient
func PackedRawData(data []byte, p *PackedSegmentClient) error {
	fmt.Printf("PackedSegmentClient: Received %d bytes.\n", len(data))
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

/*
func ProcessPackedData(entry *pb.ChunkedEntry, msc *MessageServerClient) error {
	//fmt.Println("ProcessMessageData")
	s := entry.String()
	if len(s) <= 0 {
		return nil
	}
	//fmt.Println(s)
	var e string
	if ma := regexp.MustCompile(`^([\w]+):`).FindStringSubmatch(s); len(ma) > 0 {
		e = ma[1]
	}
	switch e {
	case "next":
		//next:{at:1723789941}
		msc.mu.Lock()
		if ma := regexp.MustCompile(`{at:([\d]+)}`).FindStringSubmatch(s); len(ma) > 0 {
			//fmt.Println(ma[1])
			if msc.nextStreamAt != "now" {
				msc.beforeNextStreamAt = msc.nextStreamAt
			}
			msc.nextStreamAt = ma[1]
		} else {
			msc.mu.Unlock()
			fmt.Println("entry next error: "+s)
			return nil
		}
		if msc.beforeNextStreamAt == msc.nextStreamAt {
			msc.mu.Unlock()
			return msc.onNetworkError()
		}
		msc.mu.Unlock()
		//log.Println("nextAt: ", msc.nextStreamAt)
		//log.Println("beforenextAt: ", msc.beforeNextStreamAt)
		msc.cmd <- e
	case "backward":
		//backward:{until:{seconds:1723789900}  segment:{uri:"https://mpn.live.nicovideo.jp/data/backward/v4/BBxEfXcPJuFVyZ97aTmoSSLC4mVIjNHLXX6cMHpoJSjj5Pqqp4odv_9O_2dYB6oiaO-SuaVX34RJTDToKZNwr5gBWks"}  snapshot:{uri:"https://mpn.live.nicovideo.jp/data/snapshot/v4/BByuTtvHa5vSWxnGEbDrPivYTDLuPGR2W1WXoiCRISeTQwgw-T27nbvwovofl3rKo3heRUkha5Mb42vsPvw4Qw"}}
		//fmt.Println(s)
		if ma := regexp.MustCompile(`{segment:{uri:"([^"]+)"}`).FindStringSubmatch(s); len(ma) > 0 {
			//fmt.Println("backword uri: "+ma[1])
		}
		msc.cmd <- e
	case "previous":
		//previous:{from:{seconds:1723789916}  until:{seconds:1723789932}  uri:"https://mpn.live.nicovideo.jp/data/segment/v4/BBzuEZXfmsvy4vfcCoBFmp0sjQJX13dqzTxyrxhNIw_2kLl1Jsc6tllJh93dITT5CKj7_U16-MvwtIt-DKIFmr2k"}
		//fmt.Println(s)
		if ma := regexp.MustCompile(`uri:"([^"]+)"}`).FindStringSubmatch(s); len(ma) > 0 {
			//fmt.Println("previous uri: "+ma[1])
		}
		msc.cmd <- e
	case "segment":
		//segment:{from:{seconds:1723789932}  until:{seconds:1723789948}  uri:"https://mpn.live.nicovideo.jp/data/segment/v4/BBwWCLcROYRA-MqsINQ8cjWLXsAqzVNfiMfFlT-UI6CxOQweAhdxlC305oHkdckSTggbyDbPgEzO-1BIbFrP-WpF"}
		//fmt.Println(s)
		if ma := regexp.MustCompile(`uri:"([^"]+)"}`).FindStringSubmatch(s); len(ma) > 0 {
			fmt.Println("segment uri: "+ma[1])
			//hls.ConnectSegmentServer(ma[1], false)
		}
		msc.cmd <- e
	default:
		fmt.Println("Unknown entry: "+s)
	}

	return nil
}
*/

func PackedNetworkError() error {
	log.Println("packedSegmentServer Network error")

	return nil
}
