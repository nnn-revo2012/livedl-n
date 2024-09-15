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

type MessageServerClient struct {
	nextStreamAt        string
	processData         func(*pb.ChunkedEntry, *MessageServerClient, *NicoHls) error
	streamReceiver      *StreamReceiver[MessageServerClient]
	uri                 string
	headers             map[string]string
	isDisconnect        bool
	stream              *BinaryStream
	hls                 *NicoHls
	onNetworkError      func() error
	beforeNextStreamAt  string
	mu                  sync.Mutex
}

func NewMessageServerClient(uri string, processData func(*pb.ChunkedEntry, *MessageServerClient, *NicoHls) error, hls *NicoHls, onNetworkError func() error) *MessageServerClient {
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
		hls:                hls,
		onNetworkError:     onNetworkError,
		beforeNextStreamAt: "",
	}
}

func (msc *MessageServerClient) DoConnect() error {
	//fmt.Println("msc.DoConnect")
	//for !msc.isDisconnect && !msc.IsUnexpectedDisconnect() {
		//err := msc.streamReceiver.Receive(msc.uri + "?at=" + msc.nextStreamAt, msc.headers, msc.stream)
		err := msc.streamReceiver.Receive(msc.uri + "?at=" + msc.nextStreamAt, msc.headers, msc, msc.hls)
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

func ProcessRawData(data []byte, msc *MessageServerClient, hls *NicoHls) error {
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
		if err := msc.processData(entry, msc, msc.hls); err != nil {
			return err
		}
	}

	msc.stream.CheckClearBuffer()
	return nil
}

func ProcessMessageData(entry *pb.ChunkedEntry, msc *MessageServerClient, hls *NicoHls) error {
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
			return nil
			return msc.onNetworkError()
		}
		msc.mu.Unlock()
		//log.Println("nextAt: ", msc.nextStreamAt)
		//log.Println("beforenextAt: ", msc.beforeNextStreamAt)
	case "backward":
		//backward:{until:{seconds:1723789900}  segment:{uri:"https://mpn.live.nicovideo.jp/data/backward/v4/BBxEfXcPJuFVyZ97aTmoSSLC4mVIjNHLXX6cMHpoJSjj5Pqqp4odv_9O_2dYB6oiaO-SuaVX34RJTDToKZNwr5gBWks"}  snapshot:{uri:"https://mpn.live.nicovideo.jp/data/snapshot/v4/BByuTtvHa5vSWxnGEbDrPivYTDLuPGR2W1WXoiCRISeTQwgw-T27nbvwovofl3rKo3heRUkha5Mb42vsPvw4Qw"}}
		//fmt.Println(s)
		if ma := regexp.MustCompile(`{segment:{uri:"([^"]+)"}`).FindStringSubmatch(s); len(ma) > 0 {
			//fmt.Println("backword uri: "+ma[1])
		}
	case "previous":
		//previous:{from:{seconds:1723789916}  until:{seconds:1723789932}  uri:"https://mpn.live.nicovideo.jp/data/segment/v4/BBzuEZXfmsvy4vfcCoBFmp0sjQJX13dqzTxyrxhNIw_2kLl1Jsc6tllJh93dITT5CKj7_U16-MvwtIt-DKIFmr2k"}
		//fmt.Println(s)
		if ma := regexp.MustCompile(`uri:"([^"]+)"}`).FindStringSubmatch(s); len(ma) > 0 {
			//fmt.Println("previous uri: "+ma[1])
		}
	case "segment":
		//segment:{from:{seconds:1723789932}  until:{seconds:1723789948}  uri:"https://mpn.live.nicovideo.jp/data/segment/v4/BBwWCLcROYRA-MqsINQ8cjWLXsAqzVNfiMfFlT-UI6CxOQweAhdxlC305oHkdckSTggbyDbPgEzO-1BIbFrP-WpF"}
		//fmt.Println(s)
		if ma := regexp.MustCompile(`uri:"([^"]+)"}`).FindStringSubmatch(s); len(ma) > 0 {
			//fmt.Println("segment uri: "+ma[1])
			go hls.ConnectSegmentServer(ma[1], false)
		}
	default:
		fmt.Println("Unknown entry: "+s)
	}

	return nil
}

func NetworkError() error {
	log.Println("messageServer Network error")

	return nil
}
