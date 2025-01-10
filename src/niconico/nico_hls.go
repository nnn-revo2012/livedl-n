package niconico

import (
	"fmt"
	"html"
	"io/ioutil"
	"bufio"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/nnn-revo2012/livedl/files"
	"github.com/nnn-revo2012/livedl/objs"
	"github.com/nnn-revo2012/livedl/options"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/sha3"

	"context"
	"github.com/gin-gonic/gin"
	"github.com/nnn-revo2012/livedl/gorman"
	"github.com/nnn-revo2012/livedl/httpbase"
	"github.com/nnn-revo2012/livedl/procs/streamlink"
	"github.com/nnn-revo2012/livedl/procs/youtube_dl"
	"github.com/nnn-revo2012/livedl/procs"
	"log"
	"math"
	_ "net/http/pprof"

	pb "github.com/nnn-revo2012/livedl/proto"

    "google.golang.org/protobuf/encoding/protojson"
    //"google.golang.org/protobuf/proto"
)

type OBJ = map[string]interface{}

type syncData struct {
	seqNo              int
	date               int64
}
type playlist struct {
	uri                *url.URL
	uriMaster          *url.URL
	uriTimeshiftMaster *url.URL
	bandwidth          int
	nextTime           time.Time
	format             string
	withoutFormat      bool
	seqNo              int
	position           float64
}
type NicoHls struct {
	startDelay int
	playlist   playlist

	webSocketUrl string
	myUserId     string
	proxy        string

	nicoLiveId       string
	nicoNoStreamlink bool
	nicoNoYtdlp      bool

	commentStarted    bool
	mtxCommentStarted sync.Mutex

	chInterrupt  chan os.Signal
	nInterrupt   int
	mtxInterrupt sync.Mutex

	mtxRestart  sync.Mutex
	restartMain bool
	quality     string

	errNumChunk   int
	errRestartCnt int

	dbName     string
	db         *sql.DB
	dbMtx      sync.Mutex
	lastCommit time.Time

	isTimeshift        bool
	timeshiftStart     float64
	tsStart            int64
	tsStop             int64
	fastTimeshift      bool
	ultrafastTimeshift bool

	fastTimeshiftOrig      bool
	ultrafastTimeshiftOrig bool

	finish      bool
	commentDone bool

	NicoSession string
	limitBw     int
	limitBwOrig int
	syncData    []string

	nicoDebug     bool
	msgErrorCount int
	msgErrorSeqNo int
	memdb         *sql.DB
	memdbMtx      sync.Mutex
	seqNo500      int
	cnt500        int
	bw500         int

	mtxWg sync.Mutex

	gmPlst *gorman.GoroutineManager
	gmCmnt *gorman.GoroutineManager
	gmDB   *gorman.GoroutineManager
	gmMain *gorman.GoroutineManager
}

func debug_Now() string {
	return time.Now().Format("2006/01/02-15:04:05")
}
func NewHls(opt options.Option, prop map[string]interface{}) (hls *NicoHls, err error) {

	webSocketUrl, ok := prop["//webSocketUrl"].(string)
	if !ok {
		err = fmt.Errorf("webSocketUrl is not string")
		return
	}

	myUserId, _ := prop["//myId"].(string)
	if myUserId == "" {
		myUserId = "NaN"
	}

	var timeshift bool
	if status, ok := prop["status"].(string); ok && status == "ENDED" {
		timeshift = true
	}

	var pid string
	if nicoliveProgramId, ok := prop["nicoliveProgramId"]; ok {
		pid, _ = nicoliveProgramId.(string)
	}

	var uname string // ユーザ名
	var uid string   // ユーザID
	var cname string // コミュ名 or チャンネル名
	var cid string   // コミュID or チャンネルID

	var pt string
	if providerType, ok := prop["providerType"]; ok {
		if pt, ok = providerType.(string); ok {
			if pt == "official" {
				uname = "official"
				uid = "official"
				cname = "official"
				cid = "official"
			}
		}
	}

	// ユーザ名
	if userName, ok := prop["userName"]; ok {
		uname, _ = userName.(string)
	}

	// ユーザID
	if userPageUrl, ok := prop["userPageUrl"]; ok {
		if u, ok := userPageUrl.(string); ok {
			if m := regexp.MustCompile(`/user/(\d+)`).FindStringSubmatch(u); len(m) > 0 {
				uid = m[1]
				prop["userId"] = uid
			}
		}
	}
	if uid == "" && pt == "channel" {
		uid = "channel"
	}

	// コミュ名
	if socName, ok := prop["socName"]; ok {
		cname, _ = socName.(string)
	}

	// コミュID
	if comId, ok := prop["comId"]; ok {
		cid, _ = comId.(string)
	}
	if cid == "" {
		if socId, ok := prop["socId"]; ok {
			cid, _ = socId.(string)
		}
	}

	var title string
	if t, ok := prop["title"]; ok {
		title, _ = t.(string)
	}

	var beginTime int64
	if t, ok := prop["beginTime"]; ok {
		if bt, ok := t.(float64); ok {
			beginTime = int64(bt)
		}
	}
	tBegin := time.Unix(beginTime, 0)
	sYear := fmt.Sprintf("%04d", tBegin.Year())
	sMonth := fmt.Sprintf("%02d", tBegin.Month())
	sDay := fmt.Sprintf("%02d", tBegin.Day())
	sDay8 := fmt.Sprintf("%04d%02d%02d", tBegin.Year(), tBegin.Month(), tBegin.Day())
	sDay6 := fmt.Sprintf("%02d%02d%02d", tBegin.Year()%100, tBegin.Month(), tBegin.Day())
	sHour := fmt.Sprintf("%02d", tBegin.Hour())
	sMinute := fmt.Sprintf("%02d", tBegin.Minute())
	sSecond := fmt.Sprintf("%02d", tBegin.Second())
	sTime6 := fmt.Sprintf("%02d%02d%02d", tBegin.Hour(), tBegin.Minute(), tBegin.Second())
	sTime4 := fmt.Sprintf("%02d%02d", tBegin.Hour(), tBegin.Minute())

	// "${PID}-${UNAME}-${TITLE}"
	dbName := opt.NicoFormat
	dbName = strings.Replace(dbName, "?PID?", files.ReplaceForbidden(pid), -1)
	dbName = strings.Replace(dbName, "?UNAME?", files.ReplaceForbidden(uname), -1)
	dbName = strings.Replace(dbName, "?UID?", files.ReplaceForbidden(uid), -1)
	dbName = strings.Replace(dbName, "?CNAME?", files.ReplaceForbidden(cname), -1)
	dbName = strings.Replace(dbName, "?CID?", files.ReplaceForbidden(cid), -1)
	dbName = strings.Replace(dbName, "?TITLE?", files.ReplaceForbidden(title), -1)
	// date,time
	dbName = strings.Replace(dbName, "?YEAR?", sYear, -1)
	dbName = strings.Replace(dbName, "?MONTH?", sMonth, -1)
	dbName = strings.Replace(dbName, "?DAY?", sDay, -1)
	dbName = strings.Replace(dbName, "?DAY8?", sDay8, -1)
	dbName = strings.Replace(dbName, "?DAY6?", sDay6, -1)
	dbName = strings.Replace(dbName, "?HOUR?", sHour, -1)
	dbName = strings.Replace(dbName, "?MINUTE?", sMinute, -1)
	dbName = strings.Replace(dbName, "?SECOND?", sSecond, -1)
	dbName = strings.Replace(dbName, "?TIME6?", sTime6, -1)
	dbName = strings.Replace(dbName, "?TIME4?", sTime4, -1)

	if timeshift {
		dbName = dbName + "(TS)"
	}
	dbName = dbName + ".sqlite3"

	files.MkdirByFileName(dbName)

	var quality string
	var limitBw int
	if m := regexp.MustCompile(`^(audio_only|audio_high)$`).FindStringSubmatch(opt.NicoLimitBw); len(m) > 0 {
		quality = m[0]
		limitBw = 0
	} else {
		quality = "abr"
		limitBw, _ = strconv.Atoi(opt.NicoLimitBw)
	}

	hls = &NicoHls{
		webSocketUrl: webSocketUrl,
		myUserId:     myUserId,
		proxy:        opt.HttpProxy,

		quality: quality,
		dbName:  dbName,

		nicoLiveId: pid,
		nicoNoStreamlink: opt.NicoNoStreamlink,
		nicoNoYtdlp: opt.NicoNoYtdlp,

		isTimeshift:        timeshift,
		fastTimeshift:      opt.NicoFastTs || opt.NicoUltraFastTs,
		ultrafastTimeshift: false,

		NicoSession: opt.NicoSession,
		limitBw:     limitBw,
		limitBwOrig: limitBw,
		nicoDebug:   opt.NicoDebug,

		gmPlst: gorman.WithChecker(func(c int) { hls.checkReturnCode(c) }),
		gmCmnt: gorman.WithChecker(func(c int) { hls.checkReturnCode(c) }),
		gmDB:   gorman.WithChecker(func(c int) { hls.checkReturnCode(c) }),
		gmMain: gorman.WithChecker(func(c int) { hls.checkReturnCode(c) }),

		timeshiftStart: float64(opt.NicoTsStart),
		tsStart:  opt.NicoTsStart,
		tsStop:  opt.NicoTsStop,
	}

	hls.fastTimeshiftOrig = hls.fastTimeshift
	hls.ultrafastTimeshiftOrig = hls.ultrafastTimeshift

	for i := 0; i < 2; i++ {
		err := hls.dbOpen()
		if err != nil {
			if !strings.Contains(err.Error(), "able to open") {
				log.Fatalln(err)
			}
		} else if _, err := os.Stat(hls.dbName); err == nil {
			break
		}

		fmt.Printf("can't open: %s\n", hls.dbName)
		hls.dbName = fmt.Sprintf("%s.sqlite3", pid)
	}

	if err := hls.memdbOpen(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 放送情報をdbに入れる。自身のユーザ情報は入れない
	// dbに入れたくないデータはキーの先頭を//としている
	// 生放送の場合はデータが既にあればupdateしない
	if hls.isTimeshift || hls.dbKVExist("serverTime") == 0 {
		for k, v := range prop {
			if !strings.HasPrefix(k, "//") {
				hls.dbKVSet(k, v)
			}
		}
		if !hls.nicoNoStreamlink || !hls.nicoNoYtdlp {
			sno := hls.tsStart / 5
			hls.dbKVSet("seqStart", sno)
			eno := hls.tsStop / 5
			hls.dbKVSet("seqEnd", eno)
		}
		//fmt.Println("Write dbKVSet")
	}

	return
}
func (hls *NicoHls) Close() {
	hls.dbCommit()
	if hls.db != nil {
		hls.db.Close()
	}
	if hls.memdb != nil {
		hls.memdb.Close()
	}
}

func getModifier(attrmap map[string]interface{}) (mail []string, istranslucent bool) {
	mail = nil
	istranslucent = false

	modMap := attrmap["modifier"].(map[string]interface{})
	if modMap != nil {
		if _, ok := modMap["namedColor"].(string); ok {
			mail = append(mail, modMap["namedColor"].(string))
		} else if fcol, ok := modMap["fullColor"].(map[string]interface{}); ok {
			//fullColor:{"r":34."b":44,"g":11} //may be
			var red, blue, green int32 = 0, 0, 0
			if _, ok := fcol["r"].(float64); ok {
				red = int32(fcol["r"].(float64))
 			}
			if _, ok := fcol["g"].(float64); ok {
				green = int32(fcol["g"].(float64))
 			}
			if _, ok := fcol["b"].(float64); ok {
				blue = int32(fcol["b"].(float64))
 			}
			mail = append(mail, fmt.Sprintf("#%02x%02x%02x", red, green, blue))
		}
		if _, ok := modMap["pos"].(string); ok {
			mail = append(mail, modMap["pos"].(string))
		}
		if _, ok := modMap["size"].(string); ok {
			mail = append(mail, modMap["size"].(string))
		}
		if _, ok := modMap["font"].(string); ok {
			mail = append(mail, modMap["font"].(string))
		}
		if opacity, ok := modMap["opacity"].(string); ok {
			if opacity == "Translucent" {
				mail = append(mail, opacity)
				istranslucent = true
			}
		}
	}
	return
}

// Comment method
//func (hls *NicoHls) commentHandler(tag string, attr interface{}) (err error) {
func (hls *NicoHls) commentHandler(tag string, entry *pb.ChunkedMessage) (err error) {

	var e, s string
	if entry.GetMessage() != nil {
		s = entry.GetMessage().String()
	} else if entry.GetState() != nil {
		s = entry.GetState().String()
	//} else if entry.GetSignal() != nil {
	//	s = entry.GetSignal().String()
	}
	//fmt.Println(s)
	if ma := regexp.MustCompile(`^([\w]+):`).FindStringSubmatch(s); len(ma) > 0 {
		e = ma[1]
	}

	var jsond []byte
	switch e {
	case "chat":
		jsond, err = protojson.Marshal(entry.GetMessage().GetChat())
		if err != nil {
			return
		}
	case "simple_notification":
		jsond, err = protojson.Marshal(entry.GetMessage().GetSimpleNotification())
		if err != nil {
			return
		}
	case "gift":
		jsond, err = protojson.Marshal(entry.GetMessage().GetGift())
		if err != nil {
			return
		}
		//fmt.Println(string(jsond))
	case "nicoad":
		jsond, err = protojson.Marshal(entry.GetMessage().GetNicoad())
		if err != nil {
			return
		}
		//fmt.Println(string(jsond))
	case "game_update":
		//jsond, err = protojson.Marshal(entry.GetMessage().GetGameUpdate())
		//if err != nil {
		//	return
		//}
		fmt.Println(string(jsond))
		return
	case "tag_updated":
		//jsond, err = protojson.Marshal(entry.GetMessage().GetTagUpdated())
		//if err != nil {
		//	return
		//}
		//fmt.Println(string(jsond))
		return
	case "moderator_updated":
		jsond, err = protojson.Marshal(entry.GetMessage().GetModeratorUpdated())
		if err != nil {
			return
		}
		fmt.Println(string(jsond))
	case "ssng_updated":
		jsond, err = protojson.Marshal(entry.GetMessage().GetSsngUpdated())
		if err != nil {
			return
		}
		fmt.Println(string(jsond))
	case "overflowed_chat":
		err = fmt.Errorf("commentHandler: Recieved OverflowedChat\n")
		return
		//jsond, err = protojson.Marshal(entry.GetMessage().GetOverflowedChat())
		//if err != nil {
		//	return
		//}
	//ここからstate
	case "marquee":
		jsond, err = protojson.Marshal(entry.GetState().GetMarquee())
		if err != nil {
			return
		}
		if string(jsond) == "{}" {
			//fmt.Println("marquee: clear "+string(jsond))
			return
		}
		jsond, err = protojson.Marshal(entry.GetState().GetMarquee().GetDisplay().GetOperatorComment())
		if err != nil {
			return
		}
		//fmt.Println("marquee: "+string(jsond))
	case "enquete":
		jsond, err = protojson.Marshal(entry.GetState().GetEnquete())
		if err != nil {
			return
		}
		//fmt.Println("enquate: "+string(jsond))
	case "statistics":
		//jsond, err = protojson.Marshal(entry.GetState().GetStatistics())
		//if err != nil {
		//	return
		//}
		//fmt.Println("statistics: "+string(jsond))
		return
	case "trial_panel":
		jsond, err = protojson.Marshal(entry.GetState().GetTrialPanel())
		if err != nil {
			return
		}
		//fmt.Println("trial_panel: "+string(jsond))
		//fmt.Println("trial_panel: "+entry.GetState().String())
	case "program_status":
		jsond, err = protojson.Marshal(entry.GetState().GetProgramStatus())
		if err != nil {
			return
		}
		//fmt.Println("program_status: "+string(jsond))
	case "move_order":
		jsond, err = protojson.Marshal(entry.GetState().GetMoveOrder())
		if err != nil {
			return
		}
		//fmt.Println("move_order: "+string(jsond))
	default:
		fmt.Printf("Unknown data: %s\n",entry.String())
		return
	}

	var attrMap map[string]interface{}
	if err := json.Unmarshal(jsond, &attrMap); err != nil {
		err = fmt.Errorf("[FIXME] commentHandler: not a map: %#v", jsond)
		return err
	}
	//fmt.Printf("%#v\n", attrMap)

	if tag == "chat" {
		var vpos int64
		if d, ok := attrMap["vpos"].(float64); ok {
			vpos = int64(d)
		}
		var date, date_usec int64
		if ma := regexp.MustCompile(`at:{seconds:([\d]+)`).FindStringSubmatch(entry.GetMeta().String()); len(ma) > 0 {
			date, _ = strconv.ParseInt(ma[1], 10, 64)
		}
		if ma := regexp.MustCompile(`nanos:([\d]+)}`).FindStringSubmatch(entry.GetMeta().String()); len(ma) > 0 {
			date_usec, _ = strconv.ParseInt(ma[1], 10, 64)
			date_usec /= 1000	//ナノ秒->マイクロ秒
		}
		date2 := (date * 1000 * 1000) + date_usec
		var user_id string
		var mail []string
		if s, ok := attrMap["hashedUserId"].(string); ok {
			user_id = s
			attrMap["anonymity"] = 1
			mail = append(mail, "184")
		} else if s, ok := attrMap["rawUserId"].(string); ok {
			user_id = s
		}
		if _, ok := attrMap["accountStatus"].(string); ok {
			attrMap["premium"] = 1
		}
		var content string
		switch e {
		case "chat":
			modifier, istranslucent := getModifier(attrMap)
			if modifier != nil {
				mail = append(mail, modifier...)
			}
			if istranslucent {
				//normalアカウントで混雑時 or チャンネル非会員の場合
				if attrMap["premium"] == 1 {
					attrMap["premium"] = 25
				} else {
					attrMap["premium"] = 24
				}
			}
			if mail != nil {
				attrMap["mail"] = strings.Join(mail, " ")
			}
			if s, ok := attrMap["content"].(string); ok {
				content = s
			}
		case "simple_notification":
			attrMap["premium"] = 3
			if s, ok := attrMap["emotion"].(string); ok {
				content = "/emotion " + s
			} else if s, ok := attrMap["cruise"].(string); ok {
				content = "/cruise \"" + s + "\""
			} else if s, ok := attrMap["quote"].(string); ok {
				content = "/quote \"" + s + "\""
			} else if s, ok := attrMap["programExtended"].(string); ok {
				content = "/info 3 " + s	//3秒
			} else if s, ok := attrMap["rankingIn"].(string); ok {
				content = "/info 8 " + s	//8秒
			} else if s, ok := attrMap["rankingUpdated"].(string); ok {
				content = "/info 8 " + s	//8秒
			} else if s, ok := attrMap["visited"].(string); ok {
				content = "/info 10 " + s
			} else if s, ok := attrMap["ichiba"].(string); ok {
				content = "/info 10 " + s
			} else {
				content = "/info 10 " + string(jsond)
				fmt.Printf("%s\n",content)
			}
		case "gift":
			attrMap["premium"] = 3
			if ttt, ok := attrMap["itemId"].(string); ok {
				content = "/gift " + ttt
			}
			if ttt, ok := attrMap["advertiserUserId"].(string); ok {
				content += " " + ttt
			} else {
				content += " NULL"
			}
			if ttt, ok := attrMap["advertiserName"].(string); ok {
				content += " \"" + ttt + "\""
			}
			if ttt, ok := attrMap["point"].(string); ok {
				content += " " + ttt
			}
			if ttt, ok := attrMap["itemName"].(string); ok {
				content += " \"\" \"" + ttt + "\""
			}
			if rank, ok := attrMap["contributionRank"].(float64); ok {
				content += " " + strconv.Itoa(int(rank))
			}
			if len(content) <= 0 { 
				content = "/gift " + string(jsond)
				fmt.Printf("[FIXME]: %s\n",content)
			}
		case "nicoad":
			attrMap["premium"] = 3
			if ma := regexp.MustCompile(`^{"v1":{(.+)}}$`).FindStringSubmatch(string(jsond)); len(ma) > 0 {
				content = "/nicoad {\"version:\":\"1\"," + ma[1] + "}"
			} else if ma := regexp.MustCompile(`^{"v0":{(.+)}}$`).FindStringSubmatch(string(jsond)); len(ma) > 0 {
				content = "/nicoad {\"version:\":\"0\"," + ma[1] + "}"
			} else {
				content = "/nicoad " + string(jsond)
				fmt.Printf("[FIXME]: %s\n",content)
			}
		case "marquee":
			if attrMap["modifier"] != nil {
				mail, _ := getModifier(attrMap)
				if mail != nil {
					attrMap["mail"] = strings.Join(mail, " ")
				}
			}
			attrMap["premium"] = 3
			if s, ok := attrMap["content"].(string); ok {
				content = s
			}
			if s, ok := attrMap["link"].(string); ok {
				content += "(\"" + s + "\")"
			}
		case "enquete":
			attrMap["premium"] = 3
			isResult := false
			if len(entry.GetState().GetEnquete().String()) <= 0 {
				content = "/vote stop"
			} else {
				if entry.GetState().GetEnquete().GetChoices() != nil {
					choicesMap := attrMap["choices"].([]interface{})
					if choicesMap[0].(map[string]interface{})["perMille"] != nil {
						isResult = true
						content = "/vote showresult per "
					} else {
						content = "/vote start \"" + attrMap["question"].(string) + "\" "
					}
					for _, item := range choicesMap {
						descMap := item.(map[string]interface{})
						if isResult {
							if permille, ok := descMap["perMille"].(float64); ok {
								content += fmt.Sprintf("%d ", int32(permille))
							}
						}else {
							content += "\"" + descMap["description"].(string) + "\" "
						}
					}
				}
			}
		case "trial_panel":
			attrMap["premium"] = 3
			if ttt, ok := attrMap["panel"].(string); ok {
				content = "/trial_panel panel " + ttt
			} else if ttt, ok := attrMap["unqualifiedUser"].(string); ok {
				content = "/trial_panel unqualified_user " + ttt
			} else {
				content = "/trial_panel " + string(jsond)
			}
		case "program_status":
			attrMap["premium"] = 3
			if ttt, ok := attrMap["state"].(string); ok {
				content = "/status " + ttt
			} else {
				content = "/program_status " + string(jsond)
			}
		case "move_order":
			attrMap["premium"] = 3
			if jumpMap, ok := attrMap["jump"].(map[string]interface{}); ok {
				content = "/move_order " + jumpMap["message"].(string) +
					"(https://live.nicovideo.jp/watch/" + jumpMap["content"].(string) + ")"
			} else if redirectMap, ok := attrMap["redirect"].(map[string]interface{}); ok {
				content = "/move_order " + redirectMap["message"].(string) +
					"(" + redirectMap["uri"].(string) + ")"
			} else {
				content = "/move_order " + string(jsond)
			}
		default:
			attrMap["premium"] = 3
			content = "/??? " + string(jsond)
			fmt.Printf("[FIXME]: %s\n",content)
		}

		calc_s := fmt.Sprintf("%d,%d,%d,%s,%s", vpos, date, date_usec, user_id, content)
		hash := fmt.Sprintf("%x", sha3.Sum256([]byte(calc_s)))

		if ma := regexp.MustCompile(`meta:{id:"([^"]*)"`).FindStringSubmatch(entry.String()); len(ma) > 0 {
			attrMap["thread"] = ma[1]
		}

		//fmt.Printf("date=\"%d\" date_usec=\"%d\" date2=\"%d\" vpos=\"%d\">%s</chat>\n",
		//	date, date_usec, date2, vpos, attrMap["content"])

		hls.dbInsert("comment", map[string]interface{}{
			"vpos":      vpos,
			"date":      date,
			"date_usec": date_usec,
			"date2":     date2,
			"no":        attrMap["no"],
			"anonymity": attrMap["anonymity"],
			"user_id":   user_id,
			"content":   content,
			"mail":      attrMap["mail"],
			"name":      attrMap["name"],
			"premium":   attrMap["premium"],
			"score":     0,
			"thread":    attrMap["thread"],
			"origin":    "",
			"locale":    "",
			"hash":      hash,
		})
	} else {
		// 生放送の場合はデータが既にあればupdateしない
		if hls.isTimeshift || hls.dbKVExist("comment/thread") == 0 {
			if d, ok := attrMap["thread"].(float64); ok {
				hls.dbKVSet("comment/thread", fmt.Sprintf("%.f", d))
			} else if s, ok := attrMap["thread"].(string); ok {
				hls.dbKVSet("comment/thread", s)
			}
			//fmt.Println("Write dbKVSet(command/thread)")
		}
	}

	return
}

// return code
const (
	OK = iota
	INTERRUPT
	MAIN_WS_ERROR
	MAIN_DISCONNECT
	MAIN_END_PROGRAM
	MAIN_INVALID_STREAM_QUALITY
	MAIN_TEMPORARILY_ERROR
	PLAYLIST_END
	TIMESHIFT_STOP
	PLAYLIST_403
	PLAYLIST_ERROR
	DELAY
	COMMENT_WS_ERROR
	COMMENT_SAVE_ERROR
	COMMENT_DONE
	GOT_SIGNAL
	ERROR_SHUTDOWN
	NETWORK_ERROR
)

func (hls *NicoHls) stopPCGoroutines() {
	hls.stopPGoroutines()
	hls.stopCGoroutines()
}
func (hls *NicoHls) stopAllGoroutines() {
	hls.stopPGoroutines()
	hls.stopCGoroutines()
	hls.stopMGoroutines()
}
func (hls *NicoHls) stopPGoroutines() {
	hls.gmPlst.Cancel()
}
func (hls *NicoHls) stopCGoroutines() {
	hls.gmCmnt.Cancel()
}
func (hls *NicoHls) stopMGoroutines() {
	hls.gmMain.Cancel()
}
func (hls *NicoHls) working() bool {
	return hls.gmPlst.Count() > 0 || hls.gmCmnt.Count() > 0 || hls.gmDB.Count() > 0
}

func (hls *NicoHls) stopInterrupt() {
	if hls.chInterrupt != nil {
		signal.Stop(hls.chInterrupt)
	}
}
func (hls *NicoHls) startInterrupt() {
	if hls.chInterrupt == nil {
		hls.chInterrupt = make(chan os.Signal, 10)
		signal.Notify(hls.chInterrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	}

	hls.startMGoroutine(func(sig <-chan struct{}) int {
		select {
		case <-hls.chInterrupt:
			hls.IncrInterrupt()
			fmt.Printf("Interrupt count: %d\n", hls.nInterrupt)
			go func() {
				hls.dbCommit()
			}()
			if hls.nInterrupt >= 2 {
				os.Exit(0)
			}
			return INTERRUPT
		case <-sig:
			return GOT_SIGNAL
		}
	})
}
func (hls *NicoHls) IncrInterrupt() {
	hls.mtxInterrupt.Lock()
	defer hls.mtxInterrupt.Unlock()
	hls.nInterrupt++
}
func (hls *NicoHls) interrupted() bool {
	hls.mtxInterrupt.Lock()
	defer hls.mtxInterrupt.Unlock()
	return hls.nInterrupt != 0
}

func (hls *NicoHls) getStartDelay() int {
	hls.mtxRestart.Lock()
	defer hls.mtxRestart.Unlock()

	return hls.startDelay
}
func (hls *NicoHls) markRestartMain(delay int) {
	hls.mtxRestart.Lock()
	defer hls.mtxRestart.Unlock()

	if (!hls.restartMain) && (!hls.finish) {
		hls.startDelay = delay
		hls.restartMain = true
	}
}
func (hls *NicoHls) checkReturnCode(code int) {
	// NEVER restart goroutines here except interrupt handler
	var fPlaylistEnd = func() {
		hls.finish = true
		if hls.isTimeshift {
			if hls.commentDone {
				hls.stopPCGoroutines()
			} else if !hls.getCommentStarted() {
				hls.stopPCGoroutines()
			} else {
				fmt.Println("waiting comment")
			}
		} else {
			hls.stopPCGoroutines()
		}
	}
	switch code {
	case NETWORK_ERROR, MAIN_TEMPORARILY_ERROR:
		delay := hls.getStartDelay()
		if delay < 1 {
			hls.markRestartMain(1)
		} else if delay < 3 {
			hls.markRestartMain(3)
		} else if delay < 13 {
			// if 3,4,5..12
			hls.markRestartMain(delay + 1)
		} else {
			hls.markRestartMain(60)
		}
		hls.stopPCGoroutines()

	case DELAY:
		//log.Println("delay")
	case PLAYLIST_403:
		// 番組終了時、websocketでEND_PROGRAMが来るよりも先にこうなるが、
		// END_PROGRAMを受信するにはwebsocketの再接続が必要
		//log.Println("403")
		if !hls.interrupted() {
			hls.markRestartMain(0)
		}
		hls.stopPGoroutines()

	case PLAYLIST_END:
		fmt.Println("playlist end.")
		fPlaylistEnd()

	case TIMESHIFT_STOP:
		fmt.Println("timeshift stop.")
		fPlaylistEnd()

	case MAIN_WS_ERROR:
		hls.stopPGoroutines()

	case MAIN_DISCONNECT:
		hls.stopPCGoroutines()

	case MAIN_END_PROGRAM:
		hls.finish = true
		hls.stopPCGoroutines()

	case MAIN_INVALID_STREAM_QUALITY:
		hls.markRestartMain(0)
		hls.stopPGoroutines()

	case PLAYLIST_ERROR:
		hls.stopPCGoroutines()

	case COMMENT_WS_ERROR:
		//log.Println("comment websocket error")
		hls.stopCGoroutines()

	case COMMENT_SAVE_ERROR:
		//log.Println("comment save error")
		hls.stopCGoroutines()

	case INTERRUPT:
		hls.startInterrupt()
		hls.stopPCGoroutines()

	case ERROR_SHUTDOWN:
		hls.stopPCGoroutines()

	case COMMENT_DONE:
		hls.commentDone = true
		if hls.finish {
			hls.stopPCGoroutines()
		}

	case OK:
	}
}

// Of playlist
func (hls *NicoHls) startPGoroutine(f func(<-chan struct{}) int) {
	if !hls.interrupted() {
		hls.gmPlst.Go(f)
	}
}

// Of comment
func (hls *NicoHls) startCGoroutine(f func(<-chan struct{}) int) {
	if !hls.interrupted() {
		hls.gmCmnt.Go(f)
	}
}

// Of DB
func (hls *NicoHls) startDBGoroutine(f func(<-chan struct{}) int) {
	if !hls.interrupted() {
		hls.gmDB.Go(f)
	}
}

// Of main
func (hls *NicoHls) startMGoroutine(f func(<-chan struct{}) int) {
	hls.gmMain.Go(f)
}

func (hls *NicoHls) waitRestartMain() bool {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		fn := runtime.FuncForPC(pc)
		if !strings.HasSuffix(fn.Name(), ".Wait") {
			log.Printf("[FIXME] Don't call waitRestartMain from %s\n", fn.Name())
		}
	}

	hls.waitPGoroutines()

	hls.mtxRestart.Lock()
	defer hls.mtxRestart.Unlock()
	if hls.restartMain {
		hls.restartMain = false
		//hls.wgPlaylist = &sync.WaitGroup{}
		hls.startMain()
		return true
	}
	return false
}

func (hls *NicoHls) waitPGoroutines() {
	hls.gmPlst.Wait()
}
func (hls *NicoHls) waitCGoroutines() {
	hls.gmCmnt.Wait()
}
func (hls *NicoHls) waitDBGoroutines() {
	hls.gmDB.Wait()
}
func (hls *NicoHls) waitMGoroutines() {
	hls.gmMain.Wait()
}
func (hls *NicoHls) waitAllGoroutines() {
	hls.waitPGoroutines()
	hls.waitCGoroutines()
	hls.waitDBGoroutines()
	hls.waitMGoroutines()
}

func (hls *NicoHls) getTsCommentFromWhen() (when int64) {
	return hls.dbGetFromWhen()
}

func (hls *NicoHls) setCommentStarted(val bool) {
	hls.mtxCommentStarted.Lock()
	defer hls.mtxCommentStarted.Unlock()
	hls.commentStarted = val
}
func (hls *NicoHls) getCommentStarted() bool {
	hls.mtxCommentStarted.Lock()
	defer hls.mtxCommentStarted.Unlock()
	return hls.commentStarted
}

var split = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data) ; i++ {
		if data[i] == '\n' {
			return i + 1, data[:i + 1], nil
		}
		if data[i] == '\r' {
			if (i + 1) == len(data) {
				return 0, nil, nil
			}
			if data[i + 1] == '\n' {
				return i + 2, data[:i + 2], nil
			}
			return i + 1, data[:i + 1], nil
		}
	}

	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}

	return 0, nil, nil
}

//func execStreamlink(gm *gorman.GoroutineManager, uri, name string) (notSupport bool, err error) {
func (hls *NicoHls) execStreamlink(uri, name, session string) (notSupport bool, err error) {

	cmd, stdout, stderr, err := streamlink.Open(uri, "--http-cookie=user_session="+session, "best", "--retry-max", "10", "-o", name)
	if err != nil {
		return
	}
	defer stdout.Close()
	defer stderr.Close()

	chStdout := make(chan string, 10)
	chStderr := make(chan string, 10)
	chEof := make(chan struct{}, 2)

	// stdout
	//gm.Go(func(c <-chan struct{}) int {
	hls.startCGoroutine(func(sig <-chan struct{}) int {
		defer func(){
			chEof <- struct{}{}
		}()
		scanner := bufio.NewScanner(stdout)
		scanner.Split(split)

		for scanner.Scan() {
			chStdout <- scanner.Text()
		}

		return 0
	})

	// stderr
	//gm.Go(func(c <-chan struct{}) int {
	hls.startCGoroutine(func(sig <-chan struct{}) int {
		defer func(){
			chEof <- struct{}{}
		}()
		scanner := bufio.NewScanner(stderr)
		scanner.Split(split)

		for scanner.Scan() {
			chStderr <- scanner.Text()
		}

		return 0
	})


	// outputs
	//gm.Go(func(c <-chan struct{}) int {
	hls.startCGoroutine(func(sig <-chan struct{}) int {
		for {
			var s string
			select {
			case s = <-chStdout:
			case s = <-chStderr:
			case <-chEof:
				return 0
			}

			if strings.HasPrefix(s, "[cli][error]") {
				fmt.Print(s)

				notSupport = true
				procs.Kill(cmd.Process.Pid)
				break
			} else if strings.HasPrefix(s, "Traceback (most recent call last):") {
				fmt.Print(s)

				notSupport = true
				//procs.Kill(cmd.Process.Pid)
				//break
			} else {
				fmt.Print(s)
			}
		}
		return 0
	})

	cmd.Wait()

	return
}

//func execYoutube_dl(gm *gorman.GoroutineManager, uri, name string) (err error) {
func (hls *NicoHls) execYoutube_dl(uri, name, session string) (err error) {
	defer func() {
		part := name + ".part"
		if _, test := os.Stat(part); test == nil {
			if _, test := os.Stat(name); test != nil {
				os.Rename(part, name)
			}
		}
	}()

	cmd, stdout, stderr, err := youtube_dl.Open("--add-headers", "Cookie: user_session="+session, "-no-progress", "-o", name, uri)
	if err != nil {
		return
	}
	defer stdout.Close()
	defer stderr.Close()

	chStdout := make(chan string, 10)
	chStderr := make(chan string, 10)
	chEof := make(chan struct{}, 2)

	// stdout
	//gm.Go(func(c <-chan struct{}) int {
	hls.startCGoroutine(func(sig <-chan struct{}) int {
		defer func(){
			chEof <- struct{}{}
		}()
		scanner := bufio.NewScanner(stdout)
		scanner.Split(split)

		for scanner.Scan() {
			chStdout <- scanner.Text()
		}

		return 0
	})

	// stderr
	//gm.Go(func(c <-chan struct{}) int {
	hls.startCGoroutine(func(sig <-chan struct{}) int {
		defer func(){
			chEof <- struct{}{}
		}()
		scanner := bufio.NewScanner(stderr)
		scanner.Split(split)

		for scanner.Scan() {
			chStderr <- scanner.Text()
		}

		return 0
	})

	// outputs
	//gm.Go(func(c <-chan struct{}) int {
	hls.startCGoroutine(func(sig <-chan struct{}) int {
		var old int64
		for {
			var s string
			select {
			case s = <-chStdout:
			case s = <-chStderr:
			case <-chEof:
				return 0
			}

			if strings.HasPrefix(s, "[hls @ ") || strings.HasPrefix(s, "[https @ ") {
				// ffmpeg unwanted logs
			} else {
				if strings.HasPrefix(s, "[download]") {
					var now = time.Now().UnixNano()
					if now - old > 2 * 1000 * 1000 * 1000 {
						old = now
					} else {
						continue
					}
				}
				fmt.Print(s)
			}
		}
		return 0
	})

	cmd.Wait()
	return
}

func (hls *NicoHls) startExec(nicoNoStreamlink, nicoNoYoutube_dl bool) (err error) {
	hls.startPGoroutine(func(sig <-chan struct{}) int {
		var retry bool
		var name string
		var err error
		uri := fmt.Sprintf("https://live.nicovideo.jp/watch/%s", hls.nicoLiveId)
		name = files.ChangeExtention(hls.dbName, "mp4")
		name, err = files.GetFileNameNext(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if !nicoNoStreamlink {
			retry, err = hls.execStreamlink(uri, name, hls.NicoSession)
		}
		if !hls.interrupted() {
			if err != nil || retry || (nicoNoStreamlink && (!nicoNoYoutube_dl)) {
				hls.execYoutube_dl(uri, name, hls.NicoSession)
			}
		}
		if hls.interrupted() {
			return GOT_SIGNAL
		}
		return OK
	})
	return
}

/*
func (hls *NicoHls) startComment(messageServerUri, threadId, waybackkey string) {
	if (!hls.getCommentStarted()) && (!hls.commentDone) {
		hls.setCommentStarted(true)

		hls.startCGoroutine(func(sig <-chan struct{}) int {
			defer func() {
				hls.setCommentStarted(false)
			}()

			var err error

			// here blocks several seconds
			conn, _, err := websocket.DefaultDialer.Dial(
				messageServerUri,
				map[string][]string{
					"Origin":                 []string{"https://live.nicovideo.jp"},
					"Sec-WebSocket-Protocol": []string{"msg.nicovideo.jp#json"},
					"User-Agent":             []string{httpbase.GetUserAgent()},
				},
			)
			if err != nil {
				if !hls.interrupted() {
					log.Println("comment connect:", err)
				}
				return COMMENT_WS_ERROR
			}
			var wsMtx sync.Mutex
			writeJson := func(d interface{}) error {
				wsMtx.Lock()
				defer wsMtx.Unlock()
				return conn.WriteJSON(d)
			}

			hls.startCGoroutine(func(sig <-chan struct{}) int {
				<-sig
				if conn != nil {
					conn.Close()
				}
				return OK
			})

			hls.startCGoroutine(func(sig <-chan struct{}) int {
				for !hls.interrupted() {
					select {
					case <-time.After(60 * time.Second):
						if conn != nil {
							if err := writeJson(""); err != nil {
								if !hls.interrupted() {
									log.Println("comment send null:", err)
								}
								return COMMENT_WS_ERROR
							}
						} else {
							return OK
						}
					case <-sig:
						return GOT_SIGNAL
					}
				}
				return OK
			})

			var mtxChatTime sync.Mutex
			var _chatCount int64
			incChatCount := func() {
				mtxChatTime.Lock()
				defer mtxChatTime.Unlock()
				_chatCount++
			}
			getChatCount := func() int64 {
				mtxChatTime.Lock()
				defer mtxChatTime.Unlock()
				return _chatCount
			}

			if hls.isTimeshift {

				hls.startCGoroutine(func(sig <-chan struct{}) int {
					defer func() {
						fmt.Println("Comment done.")
					}()

					var pre int64
					var finishHint int
					for !hls.interrupted() {
						select {
						case <-time.After(1 * time.Second):
							c := getChatCount()
							if c == 0 || c == pre {

								_, when := hls.getTsCommentFromWhen()

								//fmt.Printf("getTsCommentFromWhen %f %d\n", when, res_from)

								err = writeJson([]OBJ{
									OBJ{"ping": OBJ{"content": "rs:1"}},
									OBJ{"ping": OBJ{"content": "ps:5"}},
									OBJ{"thread": OBJ{
										"fork":        0,
										"nicoru":      0,
										"res_from":    -1000,
										"scores":      1,
										"thread":      threadId,
										"user_id":     hls.myUserId,
										"version":     "20061206",
										"waybackkey":  waybackkey,
										"when":        when + 1,
										"with_global": 1,
									}},
									OBJ{"ping": OBJ{"content": "pf:5"}},
									OBJ{"ping": OBJ{"content": "rf:1"}},
								})
								if err != nil {
									return NETWORK_ERROR
								}

							} else if c < pre+100 {
								// 通常,1000カウント弱増えるが、少ししか増えない場合
								finishHint++
								if finishHint > 2 {
									return COMMENT_DONE
								}

							} else {
								finishHint = 0
							}
							pre = c

						case <-sig:
							return GOT_SIGNAL
						}
					}
					return COMMENT_DONE
				})

			} else {
				err = writeJson([]OBJ{
					OBJ{"ping": OBJ{"content": "rs:0"}},
					OBJ{"ping": OBJ{"content": "ps:0"}},
					OBJ{"thread": OBJ{
						"fork":        0,
						"nicoru":      0,
						"res_from":    -100,
						"scores":      1,
						"thread":      threadId,
						"user_id":     hls.myUserId,
						"version":     "20061206",
						"with_global": 1,
					}},
					OBJ{"ping": OBJ{"content": "pf:0"}},
					OBJ{"ping": OBJ{"content": "rf:0"}},
				})
				if err != nil {
					if !hls.interrupted() {
						log.Println("comment send first:", err)
					}
					return COMMENT_WS_ERROR
				}
			}

			for !hls.interrupted() {
				select {
				case <-sig:
					return GOT_SIGNAL
				default:
					var res interface{}
					// Blocks here
					if err = conn.ReadJSON(&res); err != nil {
						return COMMENT_WS_ERROR
					}

					//fmt.Printf("debug %#v\n", res)

					if data, ok := objs.Find(res, "chat"); ok {
						if err := hls.commentHandler("chat", data); err != nil {
							return COMMENT_SAVE_ERROR
						}
						incChatCount()

					} else if data, ok := objs.Find(res, "thread"); ok {
						if err := hls.commentHandler("thread", data); err != nil {
							return COMMENT_SAVE_ERROR
						}

					} else if _, ok := objs.Find(res, "ping"); ok {
						// nop
					} else {
						fmt.Printf("[FIXME] Unknown Message: %#v\n", res)
					}
				}
			}
			return OK
		})
	}
}
*/

func (hls *NicoHls) ConnectPackedServer(uri string, isgetall bool) {
	hls.startCGoroutine(func(sig <-chan struct{}) int {
		var err error

		segment := make(chan *pb.PackedSegment, 3)
		psc := NewPackedSegment(uri, segment)
		var packedSegment *pb.PackedSegment
		defer func() {
			close(segment)
		}()
		// パックドセグメントサーバーに接続
		for !hls.interrupted() {
			hls.startCGoroutine(func(sig <-chan struct{}) int {
				log.Println("connect packedServer")
				err = psc.Connect()
				if err != nil {
					log.Println("packedServer connect error:", err)
					psc.Disconnect()
					return COMMENT_WS_ERROR
				}
				return OK
			})
LoopPacked:
			select {
			case <-sig:
				psc.Disconnect()
				return GOT_SIGNAL
			case <-time.After(5 * time.Second):
				if !hls.interrupted() {
					log.Println("packedServer disconnect")
					psc.Disconnect()
					return COMMENT_WS_ERROR
				}
			case  packedSegment= <-segment:
				//fmt.Println(packedSegment)
				for _, message := range packedSegment.Messages {
					if !hls.interrupted() {
						//fmt.Println(message)
						if err := hls.commentHandler("chat", message); err != nil {
							log.Println("Comment save errer:")
							return COMMENT_SAVE_ERROR
						}
					} else {
						psc.Disconnect()
						return GOT_SIGNAL
					}
				}
				if hls.interrupted() {
					psc.Disconnect()
					return GOT_SIGNAL
				}
				if !isgetall {
					return OK
				}
				if packedSegment.GetNext() != nil &&  !hls.interrupted() {
					nexturi := packedSegment.GetNext().GetUri()
					//fmt.Println("next uri: "+nexturi)
					if psc.GetNextUri() != nexturi {
						psc.SetNextUri(nexturi)
						time.Sleep(500 * time.Millisecond)
						break LoopPacked
					} else {
						fmt.Println("Comment done.")
						return COMMENT_DONE
					}
				}
				fmt.Println("Comment done.")
				return COMMENT_DONE
			}
		}
		return COMMENT_DONE
	})
}

func (hls *NicoHls) ConnectSegmentServer(uri string, issegment bool) {
	hls.startCGoroutine(func(sig <-chan struct{}) int {
		var err error
		servername := "segment"
		if !issegment {
			servername = "previous"
		}

		message := make(chan *pb.ChunkedMessage, 3)
		ssc := NewSegmentServer(uri, servername, message)
		var chunkedMessage *pb.ChunkedMessage
		defer func() {
			close(message)
		}()
		// セグメントサーバーに接続
		if !hls.interrupted() {
			hls.startCGoroutine(func(sig <-chan struct{}) int {
				log.Println("connect "+servername+"Server")
				err = ssc.Connect()
				if err != nil {
					log.Println(servername+"Server connect error:", err)
					ssc.Disconnect()
					return COMMENT_WS_ERROR
				}
				return OK
			})
//LoopSegment:
			for !hls.interrupted() {
				select {
				case <-sig:
					ssc.Disconnect()
					return GOT_SIGNAL
				case <-time.After(10 * time.Second):
					log.Println(servername+"Server timeout(10Sec).")
					if ssc.IsUnexpectedDisconnect() {
						log.Println(servername+"Server disconnect")
						return COMMENT_WS_ERROR
					} else if ssc.IsDisconnect() {
						return OK
					}
				case chunkedMessage = <-message:
					//fmt.Println(chunkedMessage)
					s := chunkedMessage.String()
					//fmt.Println(s)
					var e string
					if ma := regexp.MustCompile(`^([\w]+):`).FindStringSubmatch(s); len(ma) > 0 {
						e = ma[1]
					}
					switch e {
					case "signal":
						//fmt.Println(s)
					case "meta":
						//fmt.Println(s)
						if err := hls.commentHandler("chat", chunkedMessage); err != nil {
							log.Println("Comment save errer:")
							return COMMENT_SAVE_ERROR
						}
					default:
						fmt.Println("Unknown entry: "+s)
					}
				}
			}
		}
		return OK
	})
}

func (hls *NicoHls) startComment(messageServerUri, threadId, waybackkey string) {

	if (!hls.getCommentStarted()) && (!hls.commentDone) {
		hls.setCommentStarted(true)

		hls.startCGoroutine(func(sig <-chan struct{}) int {
			defer func() {
				hls.setCommentStarted(false)
			}()

			var err error

			var mtxChatTime sync.Mutex
			var _chatPrevious int = 2
			decChatPrevious := func() {
				mtxChatTime.Lock()
				defer mtxChatTime.Unlock()
				_chatPrevious--
			}
			getChatPrevious := func() int {
				mtxChatTime.Lock()
				defer mtxChatTime.Unlock()
				return _chatPrevious
			}

			entry := make(chan *pb.ChunkedEntry, 3)
			msc := NewMessageServer(messageServerUri, hls.proxy, entry)
			var chunkedEntry *pb.ChunkedEntry
			defer func() {
				close(entry)
			}()
			if hls.isTimeshift {
				when := hls.getTsCommentFromWhen()
				msc.SetNextStreamAt(strconv.FormatInt(when, 10))
			}
			// メッセージサーバーに接続
			for !hls.interrupted() {
				hls.startCGoroutine(func(sig <-chan struct{}) int {
					log.Println("connect messageServer at:", msc.GetNextStreamAt())
					err = msc.Connect()
					if err != nil {
						log.Println("messageServer connect error:", err)
						msc.Disconnect()
						return COMMENT_WS_ERROR
					}
					return OK
				})
LoopMessage:
				for !hls.interrupted() {
					select {
					case <-sig:
						msc.Disconnect()
						return GOT_SIGNAL
					case <-time.After(30 * time.Second):
						if !hls.interrupted() {
							log.Println("messageServer timeout(30Sec).")
							if msc.IsDisconnect() {
								log.Println("messageServer disconnect")
								return COMMENT_WS_ERROR
							}
						}
					case chunkedEntry = <-entry:
						//fmt.Println(chunkedEntry)
						s := chunkedEntry.String()
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
								return COMMENT_WS_ERROR
							}
							if msc.beforeNextStreamAt == msc.nextStreamAt {
								msc.mu.Unlock()
								fmt.Println("messageServer Network error")
								return COMMENT_WS_ERROR
							}
							msc.mu.Unlock()
							//log.Println("nextAt: ", msc.nextStreamAt)
							//log.Println("beforenextAt: ", msc.beforeNextStreamAt)
							chunkedEntry = nil
							if !hls.isTimeshift {
								break LoopMessage
							} else {
								return OK
							}
						case "backward":
							//backward:{until:{seconds:1723789900}  segment:{uri:"https://mpn.live.nicovideo.jp/data/backward/v4/BBxEfXcPJuFVyZ97aTmoSSLC4mVIjNHLXX6cMHpoJSjj5Pqqp4odv_9O_2dYB6oiaO-SuaVX34RJTDToKZNwr5gBWks"}  snapshot:{uri:"https://mpn.live.nicovideo.jp/data/snapshot/v4/BByuTtvHa5vSWxnGEbDrPivYTDLuPGR2W1WXoiCRISeTQwgw-T27nbvwovofl3rKo3heRUkha5Mb42vsPvw4Qw"}}
							if hls.isTimeshift {
								if ma := regexp.MustCompile(`segment:{uri:"([^"]+)"}`).FindStringSubmatch(s); len(ma) > 0 {
									//fmt.Println("backword uri: "+ma[1])
									hls.ConnectPackedServer(ma[1], true)	//get allcomments
								}
							}
						case "previous":
							//previous:{from:{seconds:1723789916}  until:{seconds:1723789932}  uri:"https://mpn.live.nicovideo.jp/data/segment/v4/BBzuEZXfmsvy4vfcCoBFmp0sjQJX13dqzTxyrxhNIw_2kLl1Jsc6tllJh93dITT5CKj7_U16-MvwtIt-DKIFmr2k"}
							if !hls.isTimeshift && getChatPrevious() > 0 {
								decChatPrevious()
								if ma := regexp.MustCompile(`uri:"([^"]+)"}`).FindStringSubmatch(s); len(ma) > 0 {
									//fmt.Println("previous uri: "+ma[1])
									hls.ConnectSegmentServer(ma[1], false)	//false=previous
								}
							}
						case "segment":
							//segment:{from:{seconds:1723789932}  until:{seconds:1723789948}  uri:"https://mpn.live.nicovideo.jp/data/segment/v4/BBwWCLcROYRA-MqsINQ8cjWLXsAqzVNfiMfFlT-UI6CxOQweAhdxlC305oHkdckSTggbyDbPgEzO-1BIbFrP-WpF"}
							if !hls.isTimeshift {
								if ma := regexp.MustCompile(`uri:"([^"]+)"}`).FindStringSubmatch(s); len(ma) > 0 {
									//fmt.Println("segment uri: "+ma[1])
									hls.ConnectSegmentServer(ma[1], true)	//true=segment
								}
							}
						default:
							fmt.Println("Unknown entry: "+s)
						}
					}
				}
			}
			return OK
		})
	}
}

func urlJoin(base *url.URL, uri string) (res *url.URL, err error) {
	u, e := url.Parse(uri)
	if e != nil {
		err = e
		return
	}
	res = base.ResolveReference(u)
	return
}

func getStringBase(uri string, header map[string]string) (s string, code int, t int64, err, neterr error) {
	start := time.Now().UnixNano()
	defer func() {
		t = (time.Now().UnixNano() - start) / (1000 * 1000)
	}()

	resp, err, neterr := httpbase.Get(uri, header, nil)
	if err != nil {
		return
	}
	if neterr != nil {
		return
	}
	defer resp.Body.Close()

	bs, neterr := ioutil.ReadAll(resp.Body)
	if neterr != nil {
		return
	}

	s = string(bs)

	code = resp.StatusCode

	return
}
func getString(uri string) (s string, code int, t int64, err, neterr error) {
	return getStringBase(uri, nil)
}
func getStringHeader(uri string, header map[string]string) (s string, code int, t int64, err, neterr error) {
	return getStringBase(uri, header)
}
func postStringHeader(uri string, header map[string]string, val url.Values) (s string, code int, t int64, err, neterr error) {
	start := time.Now().UnixNano()
	defer func() {
		t = (time.Now().UnixNano() - start) / (1000 * 1000)
	}()

	resp, err, neterr := httpbase.PostForm(uri, header, nil, val)
	if err != nil {
		return
	}
	if neterr != nil {
		return
	}
	defer resp.Body.Close()

	bs, neterr := ioutil.ReadAll(resp.Body)
	if neterr != nil {
		return
	}

	s = string(bs)

	code = resp.StatusCode

	return
}

func getBytes(uri string) (code int, buff []byte, t int64, err, neterr error) {
	start := time.Now().UnixNano()
	defer func() {
		t = (time.Now().UnixNano() - start) / (1000 * 1000)
	}()

	resp, err, neterr := httpbase.Get(uri, nil, nil)
	if err != nil {
		return
	}
	if neterr != nil {
		return
	}
	defer resp.Body.Close()

	buff, neterr = ioutil.ReadAll(resp.Body)
	if neterr != nil {
		return
	}

	code = resp.StatusCode

	return
}

func (hls *NicoHls) saveMedia(seqno int, uri string) (is403, is404, is500 bool, neterr, err error) {

	var timePassed []int64
	if hls.nicoDebug {
		timePassed = append(timePassed, time.Now().UnixNano())

		start := time.Now().UnixNano()
		defer func() {
			now := time.Now().UnixNano()
			timePassed = append(timePassed, now)
			t := (now - start) / (1000 * 1000)
			fmt.Fprintf(os.Stderr, "%s:saveMedia: seqno=%d, total %d(ms) %v\n", debug_Now(), seqno, t, timePassed)
		}()
	}

	code, buff, millisec, err, neterr := getBytes(uri)
	if hls.nicoDebug {
		fmt.Fprintf(os.Stderr, "%s:getBytes@saveMedia: seqno=%d, code=%v, err=%v, neterr=%v, %v(ms), len=%v\n",
			debug_Now(), seqno, code, err, neterr, millisec, len(buff))
	}
	if err != nil || neterr != nil {
		return
	}

	switch code {
	case 403:
		is403 = true
		return
	case 404:
		data := map[string]interface{}{
			"seqno":    seqno,
			"current":  hls.playlist.seqNo,
			"notfound": 1,
		}
		if hls.nicoDebug {
			timePassed = append(timePassed, time.Now().UnixNano())
		}
		hls.dbInsert("media", data)
		if hls.nicoDebug {
			timePassed = append(timePassed, time.Now().UnixNano())
		}
		hls.memdbSet404(seqno)
		is404 = true
		return
	case 500:
		is500 = true
		return
	case 200:
		// OK
	}

	data := map[string]interface{}{
		"seqno":     seqno,
		"current":   hls.playlist.seqNo,
		"size":      len(buff),
		"bandwidth": hls.playlist.bandwidth,
		"data":      buff,
	}

	if seqno == hls.playlist.seqNo {
		if hls.isTimeshift {
			data["position"] = hls.playlist.position
		}
	}

	if hls.nicoDebug {
		timePassed = append(timePassed, time.Now().UnixNano())
	}
	hls.dbReplace("media", data)
	if hls.nicoDebug {
		timePassed = append(timePassed, time.Now().UnixNano())
	}
	hls.memdbSet200(seqno)

	return
}

func (hls *NicoHls) getPlaylist(argUri *url.URL) (is403, isEnd, isStop, is500 bool, neterr, err error) {
	u := argUri.String()
	m3u8, code, millisec, err, neterr := getString(u)
	if hls.nicoDebug {
		fmt.Fprintf(os.Stderr, "%s:getPlaylist: code=%v, err=%v, neterr=%v, %v(ms) >>>%s<<<\n",
			debug_Now(), code, err, neterr, millisec, m3u8)
	}
	if err != nil || neterr != nil {
		return
	}

	switch code {
	case 200:
	case 403:
		is403 = true
		return
	default:
		if 500 <= code && code <= 599 {
			if strings.Contains(u, "playlist.m3u8") || !strings.Contains(u, "master.m3u8") {
				if hls.seqNo500 == hls.playlist.seqNo {
					hls.cnt500++
					if hls.cnt500 >= 3 {
						if hls.bw500 == hls.playlist.bandwidth {
							err = fmt.Errorf("# playlist code=%v, hls.bw500=%v, hls.playlist.bandwidth=%v",
								code, hls.bw500, hls.playlist.bandwidth,
							)
							return
						} else {
							hls.bw500 = hls.playlist.bandwidth
							fmt.Printf("Changing limitBw: %v -> %v\n", hls.limitBw, hls.playlist.bandwidth-1)
							hls.limitBw = hls.playlist.bandwidth - 1
						}
					}
				} else {
					hls.seqNo500 = hls.playlist.seqNo
					hls.cnt500 = 1
				}
			} else {
				// master.m3u8が500
				hls.seqNo500 = -1
				hls.cnt500 = 0
				hls.bw500 = -1
				hls.limitBw = hls.limitBwOrig
			}

			is500 = true
			return
		}
		fmt.Printf("#### playlist code: %d: %s\n", code, argUri.String())
		err = fmt.Errorf("playlist code: %d: %s", code, argUri.String())
		return
	}

	re := regexp.MustCompile(`#EXT-X-MEDIA-SEQUENCE:(\d+)`)
	ma := re.FindStringSubmatch(m3u8)
	if len(ma) > 0 {

		// Index m3u8

		// #CURRENT-POSITION:0.0
		// #DMC-CURRENT-POSITION:0.0
		var currentPos float64
		if ma := regexp.MustCompile(`#(?:DMC-)?CURRENT-POSITION:([\+\-]?\d+(?:\.\d+)?(?:[eE][\+\-]?\d+)?)`).
			FindStringSubmatch(m3u8); len(ma) > 0 {
			if hls.isTimeshift {
				n, e := strconv.ParseFloat(ma[1], 64)
				if e != nil {
					err = e
					return
				}
				currentPos = n
				hls.playlist.position = currentPos
			} else {
				// timeshiftじゃないのにCURRENT-POSITIONがあれば終了
				isEnd = true
				return
			}

		} else {
			if hls.isTimeshift {
				currentPos = hls.timeshiftStart
			}
		}

		// 総時間
		var streamDuration float64
		if hls.isTimeshift {
			if ma := regexp.MustCompile(`#(?:DMC-)?STREAM-DURATION:([\+\-]?\d+(?:\.\d+)?(?:[eE][\+\-]?\d+)?)`).
				FindStringSubmatch(m3u8); len(ma) > 0 {
				n, e := strconv.ParseFloat(ma[1], 64)
				if e != nil {
					err = e
					return
				}
				streamDuration = n
			}
		}

		var seqStart int

		seqStart, err = strconv.Atoi(ma[1])
		if err != nil {
			log.Fatal(err)
		}
		hls.playlist.seqNo = seqStart

		re := regexp.MustCompile(`#EXTINF:([\+\-]?\d+(?:\.\d+)?(?:[eE][\+\-]?\d+)?)[^\n]*\n(\S+)`)
		ma := re.FindAllStringSubmatch(m3u8, -1)

		if len(ma) == 0 {
			log.Println("No medias in playlist")
			hls.playlist.nextTime = time.Now().Add(time.Second)
			return
		}

		type seq_t struct {
			seqno    int
			duration float64
			uri      string
		}
		var seqlist []seq_t

		var seqMax int
		var totalDuration float64
		for i, a := range ma {
			var duration float64
			seqno := i + hls.playlist.seqNo
			if seqno > seqMax {
				seqMax = seqno
			}

			if hls.isTimeshift || i == 0 {
				d, e := strconv.ParseFloat(a[1], 64)
				if e != nil {
					err = e
					return
				}

				if hls.isTimeshift {
					duration = d
					totalDuration += d
				} else {
					if i == 0 {
						if d > 3 {
							fmt.Printf("debug: found EXTINF=%v\n", d)
							d = 2.0
						} else {
							d = d + 0.5
						}
						t := time.Duration(float64(time.Second) * d)
						hls.playlist.nextTime = time.Now().Add(t)
					}
				}
			}

			uri, e := urlJoin(argUri, a[2])
			if e != nil {
				err = e
				return
			}

			seqlist = append(seqlist, seq_t{
				seqno:    seqno,
				duration: duration,
				uri:      uri.String(),
			})

			// メディアのURLがシーケンス番号の部分だけが変わる形式かどうか
			if (!hls.isTimeshift) && (!hls.playlist.withoutFormat) {
				f := strings.Replace(
					strings.Replace(uri.String(), "%", "%%", -1),
					fmt.Sprintf("%d.ts?", seqno),
					"%d.ts?",
					1,
				)

				if hls.playlist.format == "" {
					hls.playlist.format = f

				} else if hls.playlist.format != f {
					fmt.Println(m3u8)
					fmt.Println("[FIXME] media format changed")
					hls.playlist.withoutFormat = true
				}
			}
		}

		if hls.isTimeshift {
				td := seqlist[0].duration * float64(time.Second)
				hls.playlist.nextTime = time.Now().Add(time.Duration(td))
		}

		// prints Current SeqNo
		if hls.isTimeshift {
			sec := int(hls.playlist.position)
			if hls.tsStop != 0 && int64(sec) >= hls.tsStop {
				isStop = true
				return
			}
			var pos string
			if sec >= 3600 {
				pos += fmt.Sprintf("%02d:%02d:%02d", sec/3600, (sec%3600)/60, sec%60)
			} else {
				pos += fmt.Sprintf("%02d:%02d", sec/60, sec%60)
			}
			fmt.Printf("Current SeqNo: %d, Pos: %s\n", hls.playlist.seqNo, pos)

		} else {
			fmt.Printf("Current SeqNo: %d\n", hls.playlist.seqNo)
		}

		minSeq := math.MaxInt32
		maxSeq := -1
		if (!hls.isTimeshift) && (!hls.playlist.withoutFormat) {
			// 404になるまで後ろに戻ってチャンクを取得する
			if hls.nicoDebug {
				fmt.Fprintf(os.Stderr, "%s:start chunks(back)\n", debug_Now())
			}
			for i := hls.playlist.seqNo - 1; i >= 0; i-- {
				if hls.memdbGetStopBack(i) {
					break
				}

				u := fmt.Sprintf(hls.playlist.format, i)
				var is404 bool
				is403, is404, _, neterr, err = hls.saveMedia(i, u)
				if neterr != nil || err != nil {
					return
				}
				if is403 {
					return
				}

				if i > maxSeq {
					maxSeq = i
				}
				if i < minSeq {
					minSeq = i
				}

				if is404 {
					break
				}
			}
		}

		// m3u8の通りにチャンクを取得する
		if hls.nicoDebug {
			fmt.Fprintf(os.Stderr, "%s:start chunks(normal)\n", debug_Now())
		}

		// 一時的に倍速モードを切っているかもしれないので戻す
		if hls.isTimeshift && (0 < hls.playlist.seqNo && hls.playlist.seqNo < 10) {
			hls.fastTimeshift = hls.fastTimeshiftOrig
			hls.ultrafastTimeshift = hls.ultrafastTimeshiftOrig
		}

		if hls.isTimeshift {
			hls.timeshiftStart = currentPos - 0.49
		}

		var found404 bool
		for _, seq := range seqlist {
			if hls.isTimeshift {
				hls.timeshiftStart += seq.duration
			}

			if hls.memdbCheck200(seq.seqno) {
				if seq.seqno == hls.playlist.seqNo {
					if hls.isTimeshift {
						hls.dbSetPosition()
					}
				}
				continue
			}

			var is404 bool
			is403, is404, is500, neterr, err = hls.saveMedia(seq.seqno, seq.uri)
			if neterr != nil || err != nil {
				return
			}
			if is404 {
				fmt.Printf("sequence 404: %d\n", seq.seqno)
				found404 = true
			}
			if is403 {
				return
			}

			// TS時、先頭(SeqNo=0)で500となる時があるが
			// Seekしなければ次回に取得可能なので一時的に倍速モードを切る
			if is500 && hls.fastTimeshift && (seq.seqno == 0) {
				fmt.Println("[WARN] disabled fastTimeshift")

				hls.fastTimeshift = false
				hls.ultrafastTimeshift = false
				return
			}

			if seq.seqno < minSeq {
				minSeq = seq.seqno
			}
			if !found404 {
				maxSeq = seq.seqno
			}
		}

		if minSeq != math.MaxInt32 && maxSeq > 0 {
			for i := minSeq; i <= maxSeq; i++ {
				hls.memdbSetStopBack(i)
			}
			hls.memdbDelete(hls.playlist.seqNo)
		}

		if strings.Contains(m3u8, "#EXT-X-ENDLIST") {
			isEnd = true
			return
		}

		if hls.isTimeshift {
			d := streamDuration - (currentPos + totalDuration)
			if d < 1.0 {
				isEnd = true
				return
			}
		}

	} else {
		// Master m3u8
		re := regexp.MustCompile(`#EXT-X-STREAM-INF:(?:[^\n]*[^\n\w-])?BANDWIDTH=(\d+)[^\n]*\n(\S+)`)
		ma := re.FindAllStringSubmatch(m3u8, -1)
		if len(ma) > 0 {
			var maxBw, idxBw int
			var uri *url.URL
			for i, a := range ma {
				bw, err := strconv.Atoi(a[1])
				if err != nil {
					log.Fatal(err)
				}

				set := func() {
					maxBw = bw
					idxBw = i
					uri, err = urlJoin(argUri, a[2])
					if err != nil {
						log.Println(err)
					}
				}

				if maxBw == 0 {
					set()

				} else if hls.limitBw > 0 {
					// with limit
					// もし現在値が制限を超えていたら、現在値より小さければセット。
					if hls.limitBw < maxBw && bw < maxBw {
						set()

						// 現在値が制限以下で、制限を超えないかつ現在値より大きければセット。
					} else if maxBw <= hls.limitBw && bw <= hls.limitBw && maxBw < bw {
						set()
					}

				} else {
					// without limit
					if maxBw < bw {
						set()
					}
				}
			}
			if uri == nil {
				err = fmt.Errorf("playlist uri not defined")
				return
			}

			if !hls.isTimeshift {
				if idxBw < 1 {
					idxBw = 1
				}
				fmt.Printf("BANDWIDTH: %d (index=%d)\n", maxBw, idxBw)
				//DBにsegnoとdateを書き込み処理
				//fmt.Printf("syncData Index=%d: %v\n", idxBw, hls.syncData[idxBw-1])
				hls.dbSyncSet(hls.syncData[idxBw-1])
			} else {
				fmt.Printf("BANDWIDTH: %d\n", maxBw)
			}
			hls.playlist.bandwidth = maxBw
			if hls.isTimeshift && hls.fastTimeshift {

			} else {
				hls.playlist.uriMaster = argUri
				hls.playlist.uri = uri
			}
			return hls.getPlaylist(uri)

		} else {
			log.Println("playlist error")
		}
	}
	return
}

func (hls *NicoHls) streamSync(uri string) {
	hls.startPGoroutine(func(sig <-chan struct{}) int {
		data, code, millisec, err, neterr := getString(uri)
		if hls.nicoDebug {
			fmt.Fprintf(os.Stderr, "%s:streamSync: code=%v, err=%v, neterr=%v, %v(ms) \n>>>%s<<<\n",
				debug_Now(), code, err, neterr, millisec, data)
		}

		//{"encoding_id":"1","beginning_timestamp":1730604446500,"sequence":2295},
		if ma := regexp.MustCompile(`"segments_metadata"\:\[\{([^\]]+)\}\]`).FindStringSubmatch(data); len(ma) > 0 {
			hls.syncData = strings.Split(ma[1], "},{")
			//idx := len(hls.syncData) - 1
			//fmt.Printf("syncData index=%d: %s\n", idx, hls.syncData[idx])
		}
		if err != nil || neterr != nil {
			return NETWORK_ERROR
		}

		if !hls.nicoNoStreamlink || !hls.nicoNoYtdlp {
			//DBにsegnoとdateを書き込み処理
			hls.dbSyncSet(hls.syncData[0])
			if ma := regexp.MustCompile(`"beginning_timestamp"\:(\d+)\,"sequence"\:(\d+)`).FindStringSubmatch(data); len(ma) > 0 {
				sno, _ := strconv.ParseInt(ma[2], 10, 64)
				hls.dbKVSet("seqStart", sno)
			}
		}

		return OK
	})
}

func (hls *NicoHls) startPlaylist(uri string) {
	hls.startPGoroutine(func(sig <-chan struct{}) int {
		hls.playlist = playlist{}
		//hls.playlist.uri = uri
		u, e := url.Parse(uri)
		if e != nil {
			return PLAYLIST_ERROR
		}

		hls.playlist.uri = u
		if hls.isTimeshift {
			hls.playlist.uriTimeshiftMaster = u
		}

		if hls.isTimeshift {
			if hls.timeshiftStart == 0 {
				hls.timeshiftStart = hls.dbGetLastPosition()
			}
			u := hls.playlist.uriTimeshiftMaster.String()
			u = regexp.MustCompile(`&start=\d+(?:\.\d*)?`).ReplaceAllString(u, "")
			u += fmt.Sprintf("&start=%f", hls.timeshiftStart)
			uri, _ := url.Parse(u)
			hls.playlist.uri = uri
		}

		for !hls.interrupted() {
			var dur time.Duration
			if hls.playlist.nextTime.IsZero() {
				dur = 0
			} else {
				now := time.Now()
				dur = hls.playlist.nextTime.Sub(now)
			}

			// 181002
			if dur < time.Second {
				dur = time.Second
			}

			if hls.nicoDebug {
				fmt.Fprintf(os.Stderr, "%s:time.After()=%v(sec)\n", debug_Now(), float64(dur)/float64(time.Second))
			}

			select {
			case <-time.After(dur):
				var uri *url.URL
				if hls.isTimeshift && hls.fastTimeshift {
					u := hls.playlist.uriTimeshiftMaster.String()
					u = regexp.MustCompile(`&start=\d+(?:\.\d*)?`).ReplaceAllString(u, "")
					u += fmt.Sprintf("&start=%f", hls.timeshiftStart)
					uri, _ = url.Parse(u)
				} else {
					uri = hls.playlist.uri
				}

				//fmt.Println(uri)

				is403, isEnd, isStop, is500, neterr, err := hls.getPlaylist(uri)
				if neterr != nil {
					if !hls.interrupted() {
						log.Println("playlist:", e)
					}
					return NETWORK_ERROR
				}
				if is500 {
					if !hls.interrupted() {
						log.Println("playlist(500):", e)
					}
					return NETWORK_ERROR
				}
				if err != nil {
					if !hls.interrupted() {
						log.Println("playlist:", e)
					}
					return PLAYLIST_ERROR
				}
				if is403 {
					return PLAYLIST_403
				}
				if isEnd {
					return PLAYLIST_END
				}
				if isStop {
					return TIMESHIFT_STOP
				}

			case <-sig:
				return GOT_SIGNAL
			}
		}
		return OK
	})
}
func (hls *NicoHls) startMain() {
	// エラー時はMAIN_*を返すこと
	hls.startPGoroutine(func(sig <-chan struct{}) int {
		if hls.nicoDebug {
			fmt.Fprintf(os.Stderr, "%s:startMain: delay = %d(sec)\n", debug_Now(), hls.startDelay)
		}

		select {
		case <-time.After(time.Duration(hls.startDelay) * time.Second):
		case <-sig:
			return GOT_SIGNAL
		}

		if hls.nicoDebug {
			fmt.Fprintf(os.Stderr, "%s:start dial main(%s)\n", debug_Now(), hls.webSocketUrl)
		}
		conn, _, err := websocket.DefaultDialer.Dial(
			hls.webSocketUrl,
			map[string][]string{
				"User-Agent": []string{httpbase.GetUserAgent()},
			},
		)
		if err != nil {
			return NETWORK_ERROR
		}
		var wsMtx sync.Mutex
		writeJson := func(d interface{}) error {
			wsMtx.Lock()
			defer wsMtx.Unlock()
			return conn.WriteJSON(d)
		}

		// debug
		if false {
			log.Printf("start ws error tsst")
			hls.startPGoroutine(func(sig <-chan struct{}) int {
				select {
				case <-time.After(10 * time.Second):
					conn.Close()
					return OK
				case <-sig:
					return GOT_SIGNAL
				}
			})
		}

		hls.startPGoroutine(func(sig <-chan struct{}) int {
			<-sig
			if conn != nil {
				conn.Close()
			}
			return OK
		})

		err = writeJson(OBJ{
			"type": "startWatching",
			"data": OBJ{
				"stream": OBJ{
					"quality":  hls.quality, //"abr", // high
					"protocol": "hls",
					"latency":  "high",
				},
				"room": OBJ{
					"protocol":    "webSocket",
					"commentable": true,
				},
				"reconnect": false, //2024/09/27 ここがtrueだとtype:messageServerが来なくなった
			},
		})
		if err != nil {
			if !hls.interrupted() {
				log.Println("websocket getpermit write:", err)
			}
			return NETWORK_ERROR
		}

		var playlistStarted bool
		var watchingStarted bool
		var watchinginterval int
		for !hls.interrupted() {
			select {
			case <-sig:
				return GOT_SIGNAL
			default:
			}
			var res interface{}
			err = conn.ReadJSON(&res)
			if err != nil {
				if (!hls.interrupted()) && (!hls.finish) {
					log.Println("websocket read:", err)
				}
				return NETWORK_ERROR
			}
			if hls.nicoDebug {
				fmt.Fprintf(os.Stderr, "%s:ReadJSON => %v\n", debug_Now(), res)
			}

			_type, ok := objs.FindString(res, "type")
			if !ok {
				fmt.Printf("type not found\n")
				continue
			}
			switch _type {
			//case "watch":
			//if cmd, ok := objs.FindString(res, "body", "command"); ok {
			//switch cmd {
			case "seat":
				if _arr, ok := objs.FindFloat64(res, "data", "keepIntervalSec"); ok {
					arr := []interface{}{_arr}
					for _, intf := range arr {
						if str, ok := intf.(float64); ok {
							num := int(str)
							if num > 0 {
								//hls.SetInterval(num)
								watchinginterval = num
								break
							}
						}
					}
				}

				if (!watchingStarted) && watchinginterval > 0 {
					watchingStarted = true
					hls.startPGoroutine(func(sig <-chan struct{}) int {
						for {
							select {
							case <-time.After(time.Duration(watchinginterval) * time.Second):
								err := writeJson(OBJ{
									"type": "keepSeat",
								})
								if err != nil {
									if !hls.interrupted() {
										log.Println("websocket watching:", err)
									}
									return NETWORK_ERROR
								}
							case <-sig:
								return GOT_SIGNAL
							}
						}
					})
				}

			case "stream":
				if !hls.isTimeshift {
					if sync_uri, ok := objs.FindString(res, "data", "syncUri"); ok {
						if (!playlistStarted) && sync_uri != "" {
							hls.streamSync(sync_uri)
						}
					}
				}
				if hls.nicoNoStreamlink && hls.nicoNoYtdlp {
					if uri, ok := objs.FindString(res, "data", "uri"); ok {
						if (!playlistStarted) && uri != "" {
							playlistStarted = true
							hls.startPlaylist(uri)
						}
					}
				} else {
					//hls.dbKVSet("startseq", fmt.Sprintf("%.f", float64(hls.timeshiftstart)))
					hls.startExec(hls.nicoNoStreamlink, hls.nicoNoYtdlp)
				}

			case "disconnect":
				// print params
				if _arr, ok := objs.FindString(res, "data", "reason"); ok {
					arr := []interface{}{0, _arr}
					fmt.Printf("%v\n", arr)
					if len(arr) >= 2 {
						if s, ok := arr[1].(string); ok {
							switch s {
							case "END_PROGRAM":
								return MAIN_END_PROGRAM
							case "SERVICE_TEMPORARILY_UNAVAILABLE", "INTERNAL_SERVERERROR":
								return MAIN_TEMPORARILY_ERROR
							case "TOO_MANY_CONNECTIONS":
								return MAIN_DISCONNECT
							case "TEMPORARILY_CROWDED":
								return MAIN_END_PROGRAM
							}
						}
					}
				}
				return MAIN_DISCONNECT

			case "room":	//2024/08/05　新メッセージサーバーに変更された
			case "messageServer":	//2024/08/05　新メッセージサーバーに変更された
				messageServerUri, ok := objs.FindString(res, "data", "viewUri")
				if ok {
					hls.startComment(messageServerUri, "", "")
				}

			case "statistics":
			case "permit":
			case "serverTime":
			case "schedule":
			case "tagUpdated":
				// nop
				//default:
				//	fmt.Printf("%#v\n", res)
				//	fmt.Printf("unknown command: %s\n", cmd)
				//} // end switch "command"
				//}
				// "watch"

			case "ping":
				err := writeJson(OBJ{
					"type": "pong",
				})
				if err != nil {
					if !hls.interrupted() {
						log.Println("websocket watching:", err)
					}
					return NETWORK_ERROR
				}
			case "error":
				code, ok := objs.FindString(res, "data", "code")
				if !ok {
					log.Printf("Unknown error: %#v\n", res)
					return ERROR_SHUTDOWN
				}

				// https://nicolive.cdn.nimg.jp/relive/front_assets/scripts/nicolib.4bb8b62b35.js
				switch code {
				case "INVALID_STREAM_QUALITY":
					// webSocket自体を再接続しないと、コメントサーバが取得できない
					switch hls.quality {
					case "abr":
						hls.quality = "high"
						return MAIN_INVALID_STREAM_QUALITY
					default:
						return ERROR_SHUTDOWN
					}
				//case
				//	"INTERNAL_SERVERERROR",
				//	"CONTENT_NOT_READY", // 終了後に出ることがある
				//	"CONNECT_ERROR": // 終了後に出ることがある
				//	return NETWORK_ERROR
				//case
				//	"INVALID_BROADCAST_ID",
				//	"BROADCAST_NOT_FOUND",
				//	"NO_THREAD_AVAILABLE",
				//	"NO_ROOM_AVAILABLE",
				//	"NO_PERMISSION":
				//	return ERROR_SHUTDOWN
				case "INVALID_MESSAGE":
					// 公式のTSで送られてきた。単純に無視する。
				default:
					//	log.Printf("Unknown error: %s\n%#v\n", code, res)
					//	return ERROR_SHUTDOWN
					fmt.Printf("error code: %v\n", code)
					if hls.msgErrorSeqNo == hls.playlist.seqNo {
						hls.msgErrorCount++
					} else {
						hls.msgErrorSeqNo = hls.playlist.seqNo
						hls.msgErrorCount = 1
					}
					if hls.msgErrorCount >= 3 {
						return ERROR_SHUTDOWN
					} else {
						return NETWORK_ERROR
					}
				}

			default:
				log.Printf("Unknown type: %s\n%#v\n", _type, res)
			} // end switch "type"
		} // for ReadJSON
		return OK
	})
}
func (hls *NicoHls) startMainV1() {
	return // old startMain
}

func (hls *NicoHls) serve(hlsPort int) {
	hls.startMGoroutine(func(sig <-chan struct{}) int {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultErrorWriter = ioutil.Discard
		gin.DefaultWriter = ioutil.Discard
		router := gin.Default()

		router.GET("", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/m3u8/2/0/index.m3u8")
			c.Abort()
		})

		router.GET("/m3u8/:delay/:shift/index.m3u8", func(c *gin.Context) {
			targetDuration := "2"
			extInf := "1.5"
			if hls.isTimeshift {
				targetDuration = "3"
				extInf = "3.0"
			}
			shift, err := strconv.Atoi(c.Param("shift"))
			if err != nil {
				shift = 0
			}
			if shift < 0 {
				shift = 0
			}
			delay, err := strconv.Atoi(c.Param("delay"))
			if err != nil {
				delay = 0
			}
			if delay < 2 {
				delay = 2
			}
			if !hls.isTimeshift {
				if delay < 4 {
					delay = 4
				}
			}
			seqno := hls.dbGetLastSeqNo(0) - int64(shift)
			body := fmt.Sprintf(
				`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:%s
#EXT-X-MEDIA-SEQUENCE:%d

`, targetDuration, seqno)
			for i := int64(delay); i >= 0; i-- {
				body += fmt.Sprintf(
					`#EXTINF:%s,
/ts/%d/test.ts

`, extInf, seqno-i)
			}
			c.Data(http.StatusOK, "application/x-mpegURL", []byte(body))
			return
		})

		router.GET("/ts/:idx/test.ts", func(c *gin.Context) {
			i, _ := strconv.Atoi(c.Param("idx"))
			b := hls.dbGetLastMedia(i)
			c.Data(http.StatusOK, "video/MP2T", b)
			return
		})

		srv := &http.Server{
			Addr:           fmt.Sprintf("127.0.0.1:%d", hlsPort),
			Handler:        router,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		chLocal := make(chan struct{})
		idleConnsClosed := make(chan struct{})
		defer func() {
			close(chLocal)
		}()
		go func() {
			select {
			case <-chLocal:
			case <-sig:
			}
			if err := srv.Shutdown(context.Background()); err != nil {
				log.Printf("srv.Shutdown: %v\n", err)
			}
			close(idleConnsClosed)
		}()

		// クライアントはlocalhostでなく127.0.0.1で接続すること
		// localhostは遅いため
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("srv.ListenAndServe: %v\n", err)
		}

		<-idleConnsClosed
		return OK
	})
}

func (hls *NicoHls) Wait(testTimeout, hlsPort int) {

	hls.startInterrupt()
	defer hls.stopInterrupt()

	if testTimeout > 0 {
		hls.startMGoroutine(func(sig <-chan struct{}) int {
			select {
			case <-sig:
				return GOT_SIGNAL
			case <-time.After(time.Duration(testTimeout) * time.Second):
				hls.chInterrupt <- syscall.Signal(1000)
				return OK
			}
		})
	}

	if hlsPort > 0 {
		hls.serve(hlsPort)
	}

	hls.startMain()
	for hls.working() {
		if hls.waitRestartMain() {
			continue
		}
		hls.stopPCGoroutines()
		hls.waitCGoroutines()
	}

	hls.stopAllGoroutines()
	hls.waitAllGoroutines()

	return
}

func postTsRsv0(opt options.Option) (err error) {
	if ma := regexp.MustCompile(`lv(\d+)`).FindStringSubmatch(opt.NicoLiveId); len(ma) > 0 {
		if err = postTsRsvBase(0, ma[1], opt.NicoSession); err != nil {
			return
		}
		err = postTsRsvBase(1, ma[1], opt.NicoSession)
	}
	return
}
func postTsRsv1(opt options.Option) (err error) {
	if ma := regexp.MustCompile(`lv(\d+)`).FindStringSubmatch(opt.NicoLiveId); len(ma) > 0 {
		err = postTsRsvBase(1, ma[1], opt.NicoSession)
	}
	return
}
func postTsRsvBase(num int, vid, session string) (err error) {
	var uri string
	if num == 0 {
		uri = fmt.Sprintf("https://live2.nicovideo.jp/api/watchingreservation?mode=watch_num&vid=%s", vid)
	} else {
		uri = fmt.Sprintf("https://live2.nicovideo.jp/api/watchingreservation?mode=confirm_watch_my&vid=%s", vid)
	}

	header := map[string]string{
		"Cookie": "user_session=" + session,
		"Referer": "https://live.nicovideo.jp/watch/lv" + vid,
		"Origin": "https://live.nicovideo.jp",
	}
	dat0, _, _, err, neterr := getStringHeader(uri, header)
	if err != nil || neterr != nil {
		if err == nil {
			err = neterr
		}
		return
	}

	var token string
	if ma := regexp.MustCompile(`(ulck_\d+)`).FindStringSubmatch(dat0); len(ma) > 0 {
		token = ma[1]
	} else if strings.Contains(dat0, "システムエラーが発生しました") {
		err = fmt.Errorf("postTsRsv: system error try again in a few minutes")
		return
	} else if strings.Contains(dat0, "申し込み期限切れ") {
		err = fmt.Errorf("postTsRsv: deadline expired")
		return
	} else if strings.Contains(dat0, "視聴済み") {
		err = fmt.Errorf("postTsRsv: already watched")
		return
	} else {
		fmt.Printf("postTsRsv: token not found: >>>%s<<<\n", dat0)
		err = fmt.Errorf("postTsRsv: token not found")
		return
	}

	// "X-Requested-With": "XMLHttpRequest",
	// "Origin": "https://live.nicovideo.jp",
	//	"Referer": fmt.Sprintf("https://live.nicovideo.jp/watch/%s", opt.NicoLiveId),
	// "X-Prototype-Version": "1.6.0.3",

	var vals url.Values
	if num == 0 {
		vals = url.Values{
			"mode":       []string{"auto_register"},
			"vid":        []string{vid},
			"token":      []string{token},
			"rec_pos":    []string{""},
			"rec_engine": []string{""},
			"rec_id":     []string{""},
			"_":          []string{""},
		}
	} else {
		vals = url.Values{
			"accept": []string{"true"},
			"mode":   []string{"use"},
			"vid":    []string{vid},
			"token":  []string{token},
			"":      []string{""},
		}
	}

	dat1, _, _, err, neterr := postStringHeader("https://live2.nicovideo.jp/api/watchingreservation", header, vals)
	if err != nil || neterr != nil {
		if err == nil {
			err = neterr
		}
		return
	}
	if (!strings.Contains(dat1, "status=\"ok\"")) && (!strings.Contains(dat1, "\"regist_finished\"")) {
		fmt.Printf("postTsRsv: status not ok: >>>%s<<<\n", dat1)
		err = fmt.Errorf("postTsRsv: status not ok")
		return
	} else {
		if num == 0 {
			fmt.Println("postTsRsv0: status ok")
		} else {
			fmt.Println("postTsRsv1: status ok")
		}
	}

	return
}

func postRsvUseTs(isRsv bool, opt options.Option) (err error) {

	var vid string
	if ma := regexp.MustCompile(`lv(\d+)`).FindStringSubmatch(opt.NicoLiveId); len(ma) > 0 {
		vid = ma[1]
	}
	uri := "https://live2.nicovideo.jp/api/timeshift.ticket.use"
	RsvMsg := "Use"
	if isRsv {
		uri = "https://live2.nicovideo.jp/api/timeshift.reservations";
		RsvMsg = "Rsv"
	}
	header := map[string]string{
		"Cookie": "user_session=" + opt.NicoSession,
		"Accept": "application/json, text/plain */*",
		"Referer": fmt.Sprintf("https://live.nicovideo.jp/watch/%s", opt.NicoLiveId),
		"Origin": "https://live.nicovideo.jp",
	}
	vals := url.Values{"vid": []string{vid},}

	dat, _, _, err, neterr := postStringHeader(uri, header, vals)
	if err != nil || neterr != nil {
		if err == nil {
			err = neterr
		}
		return
	}
	if !strings.Contains(dat, "status\":200") {
		if ma := regexp.MustCompile(`\"description\":\"([^\"]+)\"`).FindStringSubmatch(dat); len(ma) > 0 {
			err = fmt.Errorf("postRsvUseTs: %s %s\n", RsvMsg, ma[1])
		} else {
			fmt.Printf("postRsvUseTs: status not ok: >>>%s<<<\n", dat)
			err = fmt.Errorf("postRsvUseTs: %s status not ok", RsvMsg)
			return
		}
	} else {
		fmt.Printf("postRsvUseTs: %s status ok\n", RsvMsg)
	}

	return
}

func getProps(opt options.Option) (props interface{}, notLogin, rsvTs, useTs bool, err error) {

	header := map[string]string{}
	if opt.NicoSession != "" {
		header["Cookie"] = "user_session=" + opt.NicoSession
	}

	//ログインチェック
	uri := "https://www.nicovideo.jp"
	dat, _, _, err, neterr := getStringHeader(uri, header)
	if err != nil || neterr != nil {
		if err == nil {
			err = neterr
		}
		return
	}

	// ログイン判定
	if opt.NicoSession == "" {
		notLogin = true
	} else if ma := regexp.MustCompile(`login_status['"]*\s*[=:]\s*['"](.*?)['"]`).FindStringSubmatch(dat); len(ma) > 0 {
		switch string(ma[1]) {
		case "not_login":
			notLogin = true
		case "login":
			notLogin = false
		default:
			fmt.Printf("[FIXME] login_status = %s\n", ma[1])
		}
	} else {
		notLogin = true
	}

	if notLogin && opt.NicoLoginOnly {
		return
	}

	// ログインアカウント種別取得
	if ma := regexp.MustCompile(`member_status['"]*\s*[=:]\s*['"](.*?)['"]`).FindStringSubmatch(dat); len(ma) > 0 {
		fmt.Println("account:", ma[1])
	}
	if notLogin {
		fmt.Println("account: not_login")
	}

	// 放送ページ取得
	uri = fmt.Sprintf("https://live.nicovideo.jp/watch/%s", opt.NicoLiveId)
	dat, _, _, err, neterr = getStringHeader(uri, header)
	if err != nil || neterr != nil {
		if err == nil {
			err = neterr
		}
		return
	}

	// 放送種別(official/channel/community)取得
	var providertype string
	// &quot;providerType&quot;: &quot;(.+)l&quot;,
	if ma := regexp.MustCompile(`&quot;providerType&quot;:&quot;(\w+)&quot;`).FindStringSubmatch(dat); len(ma) > 0 {
		providertype = ma[1]
	}
	fmt.Println("providertype:", providertype)

	// 新配信
	if regexp.MustCompile(`ご指定のページが見つかりませんでした`).MatchString(dat) {
		err = fmt.Errorf("getProps: page not found")
	} else if regexp.MustCompile(`(放送者により(削除されました|削除された可能性があります))`).MatchString(dat) {
		err = fmt.Errorf("getProps: page not found")
	} else if ma := regexp.MustCompile(`rejectedReasons&quot;:\[([^\]]+)\]`).FindStringSubmatch(dat); len(ma) > 0 {
		ttt := strings.ReplaceAll(html.UnescapeString(ma[1]), "\",\"", " ")
		if !notLogin && providertype == "official" && regexp.MustCompile(`notHaveTimeshiftTicket`).MatchString(ttt) {
			//チケット予約処理
			fmt.Println("notHaveTimeshiftTicket: timeshift reservation required")
			rsvTs = true
		} else if !notLogin && providertype == "official" && regexp.MustCompile(`notUseTimeshiftTicket`).MatchString(ttt) {
			//チケット使用処理
			fmt.Println("notUseTimeshiftTicket: timeshift reservation required")
			useTs = true
		} else {
			err = fmt.Errorf("getProps: %s", strings.Trim(ttt, "\""))
		}
	} else if regexp.MustCompile(`webSocketUrl&quot;:&quot;&quot;,`).MatchString(dat) {
		err = fmt.Errorf("getProps: webSocketUrl not found")
	} else if ma := regexp.MustCompile(`data-props="(.+?)"`).FindStringSubmatch(dat); len(ma) > 0 {
		str := html.UnescapeString(string(ma[1]))
		if err = json.Unmarshal([]byte(str), &props); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("getProps: unknown error")
	}

	return
}

func NicoRecHls(opt options.Option) (done, playlistEnd, notLogin, reserved bool, dbName string, err error) {

	//http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 32

	//var props interface{}
	//var rsvTs bool
	//var useTs bool
	props, notLogin, rsvTs, useTs, err := getProps(opt)
	if err != nil {
		//fmt.Println(err)
		return
	}

	if notLogin {
		if opt.NicoLoginOnly {
			// 要ログイン
			return
		} else {
			// 非ログインでも録画可能なら再ログイン不要
			notLogin = false
		}
	}

	// TS予約必要
	if (rsvTs || useTs) && opt.NicoForceResv {
		if rsvTs {
			err = postRsvUseTs(true, opt)
			if (err != nil) {
				return
			}
		}
		err = postRsvUseTs(false, opt)
		if err == nil {
			reserved = true
		}
		return
	}

	if false {
		objs.PrintAsJson(props)
		os.Exit(9)
	}

	proplist := map[string][]string{
		// "community"
		"comId": []string{"community", "id"}, // "co\d+"
		// "program"
		"beginTime":         []string{"program", "beginTime"},           // integer
		"description":       []string{"program", "description"},         // 放送説明
		"endTime":           []string{"program", "endTime"},             // integer
		"isFollowerOnly":    []string{"program", "isFollowerOnly"},      // bool
		"isPrivate":         []string{"program", "isPrivate"},           // bool
		"mediaServerType":   []string{"program", "mediaServerType"},     // "DMC"
		"nicoliveProgramId": []string{"program", "nicoliveProgramId"},   // "lv\d+"
		"openTime":          []string{"program", "openTime"},            // integer
		"vposBaseTime":      []string{"program", "vposBaseTime"},        // integer
		"providerType":      []string{"program", "providerType"},        // "community"
		"status":            []string{"program", "status"},              //
		"userName":          []string{"program", "supplier", "name"},    // ユーザ名
		"userPageUrl":       []string{"program", "supplier", "pageUrl"}, // "https://www.nicovideo.jp/user/\d+"
		"title":             []string{"program", "title"},               // title
		// "site"
		"nicocas":        []string{"site", "nicocas"},                //
		"//webSocketUrl": []string{"site", "relive", "webSocketUrl"}, // "ws://..."
		"serverTime":     []string{"site", "serverTime"},             // integer
		// "socialGroup"
		"socDescription": []string{"socialGroup", "description"}, // コミュ説明
		"socId":          []string{"socialGroup", "id"},          // "co\d+" or "ch\d+"
		"socLevel":       []string{"socialGroup", "level"},       // integer
		"socName":        []string{"socialGroup", "name"},        // community name
		"socType":        []string{"socialGroup", "type"},        // "community"
		// "user"
		"accountType":  []string{"user", "accountType"}, // "premium"
		"//myId":       []string{"user", "id"},          // "\d+"
		"isLoggedIn":   []string{"user", "isLoggedIn"},  // bool
		"//myNickname": []string{"user", "nickname"},    // string
	}

	kv := map[string]interface{}{}
	for k, a := range proplist {
		v, ok := objs.Find(props, a...)
		if ok {
			kv[k] = v
		}
	}

	for _, k := range []string{
		"//webSocketUrl",
		//"//myId",
	} {
		if _, ok := kv[k]; !ok {
			fmt.Printf("%v not found\n", k)
			return
		}
	}

	if opt.NicoFormat == "" {
		opt.NicoFormat = "?PID?-?UNAME?-?TITLE?"
	}

	hls, e := NewHls(opt, kv)
	if e != nil {
		err = e
		fmt.Println(err)
		return
	}
	defer hls.Close()

	hls.Wait(opt.NicoTestTimeout, opt.NicoHlsPort)

	dbName = hls.dbName
	playlistEnd = hls.finish
	done = true

	return
}
