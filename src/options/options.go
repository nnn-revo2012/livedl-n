package options

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/nnn-revo2012/livedl/buildno"
	"github.com/nnn-revo2012/livedl/cryptoconf"
	"github.com/nnn-revo2012/livedl/files"
	"github.com/nnn-revo2012/livedl/httpbase"
	"golang.org/x/crypto/sha3"
)

const MinimumHttpTimeout = 30

var DefaultTcasRetryTimeoutMinute = 5 // TcasRetryTimeoutMinute
var DefaultTcasRetryInterval = 60     // TcasRetryInterval

type Option struct {
	Command                string
	NicoLiveId             string
	NicoStatusHTTPS        bool
	NicoSession            string
	NicoLoginAlias         string
	NicoLoginOnly          bool
	NicoCookies            string
	NicoTestTimeout        int
	TcasId                 string
	TcasRetry              bool
	TcasRetryTimeoutMinute int // 再試行を終了する時間(初回終了または録画終了からの時間「分」)
	TcasRetryInterval      int // 再試行を行うまでの待ち時間
	YoutubeId              string
	ConfFile               string // deprecated
	ConfPass               string // deprecated
	ZipFile                string
	DBFile                 string
	NicoHlsPort            int
	NicoLimitBw            string
	//NicoTsStart            float64
	NicoTsStart            int64
	NicoTsStop             int64
	NicoFormat             string
	NicoFastTs             bool
	NicoUltraFastTs        bool
	NicoAutoConvert        bool
	NicoAutoDeleteDBMode   int  // 0:削除しない 1:mp4が分割されなかったら削除 2:分割されても削除
	NicoDebug              bool // デバッグ情報の記録
	ConvExt                string
	ExtractChunks          bool
	NicoConvForceConcat    bool
	NicoConvSeqnoStart     int64
	NicoConvSeqnoEnd       int64
	NicoForceResv          bool // 終了番組の上書きタイムシフト予約
	YtNoStreamlink         bool
	YtCommentStart         float64
	YtNoYoutubeDl          bool
	YtEmoji                bool
	NicoSkipHb             bool // コメント出力時に/hbコマンドを出さない
	NicoAdjustVpos         bool // コメント出力時にvposを補正する
	HttpRootCA             string
	HttpSkipVerify         bool
	HttpProxy              string
	NoChdir                bool
	HttpTimeout            int
	NicoNoStreamlink       bool
	NicoNoYtdlp            bool
	NicoCommentOnly        bool
	NicoExecBw             string
}

func getCmd() (cmd string) {
	cmd = filepath.Base(os.Args[0])
	ext := filepath.Ext(cmd)
	cmd = strings.TrimSuffix(cmd, ext)
	return
}
func versionStr() string {
	cmd := filepath.Base(os.Args[0])
	ext := filepath.Ext(cmd)
	cmd = strings.TrimSuffix(cmd, ext)
	return fmt.Sprintf(`%s%s (%s)`, cmd, buildno.GetBuildLite(), buildno.GetBuildNo())
}
func version() {
	fmt.Println(versionStr())
	os.Exit(0)
}
func Help(verbose ...bool) {
	cmd := filepath.Base(os.Args[0])
	ext := filepath.Ext(cmd)
	cmd = strings.TrimSuffix(cmd, ext)

	format := `%s%s (%s)
Usage:
%s [COMMAND] options... [--] FILE

COMMAND:
  -nico    ニコニコ生放送の録画
  -tcas    ツイキャスの録画
  -yt      YouTube Liveの録画
  -d2m     録画済みのdb(.sqlite3)をmp4に変換する(-db-to-mp4)
  -dbinfo  録画済みのdb(.sqlite3)の各種情報を表示する
           e.g. $ livedl -dbinfo -- 'C:/home/hogehoge/livedl/rec/lvxxxxxxxx.sqlite3'
  -d2h     [実験的] 録画済みのdb(.sqlite3)を視聴するためのHLSサーバを立てる(-db-to-hls)
           開始シーケンス番号は（変換ではないが） -nico-conv-seqno-start で指定
           使用例：$ livedl lvXXXXXXXXX.sqlite3 -d2h -nico-hls-port 12345 -nico-conv-seqno-start 2780

オプション/option:
  -h         ヘルプを表示
  -vh        全てのオプションを表示
  -v         バージョンを表示
  -no-chdir  起動する時chdirしない(conf.dbは起動したディレクトリに作成されます)
  --         後にオプションが無いことを指定

ニコニコ生放送録画用オプション:
  -nico-login <id>,<password>    (+) ニコニコのIDとパスワードを指定する
                                 2段階認証(MFA)に対応しています
  -nico-session <session>        Cookie[user_session]を指定する
  -nico-login-only=on            (+) 必ずログイン状態で録画する
  -nico-login-only=off           (+) 非ログインでも録画可能とする(デフォルト)
  -nico-cookies firefox[:profile|cookiefile]
                                 firefoxのcookieを使用する(デフォルトはdefault-release)
                                 profileまたはcookiefileを直接指定も可能
                                 スペースが入る場合はquoteで囲む
  -nico-hls-port <portnum>       [実験的] ローカルなHLSサーバのポート番号
  -nico-limit-bw <bandwidth>     (+) HLSのBANDWIDTHの上限値を指定する。0=制限なし
                                 audio_high or audio_only = 音声のみ
  -nico-exec-bw "FORMAT"         (+) Streamlink/yt-dlpの場合のBANDWIDTHを指定する
                                 フォーマットはStreamlink/yt-dlpで指定するものと同じ
                                 ※-nico-limit-bwとは連動していないので注意
  -nico-format "FORMAT"          (+) 保存時のファイル名を指定する
  -nico-fast-ts                  倍速タイムシフト録画を行う(新配信タイムシフト)
  -nico-fast-ts=on               (+) 上記を有効に設定
  -nico-fast-ts=off              (+) 上記を無効に設定(デフォルト)
  -nico-auto-convert=on          (+) 録画終了後自動的にMP4に変換するように設定
  -nico-auto-convert=off         (+) 上記を無効に設定
  -nico-auto-delete-mode 0       (+) 自動変換後にデータベースファイルを削除しないように設定(デフォルト)
  -nico-auto-delete-mode 1       (+) 自動変換でMP4が分割されなかった場合のみ削除するように設定
  -nico-auto-delete-mode 2       (+) 自動変換でMP4が分割されても削除するように設定
  -nico-force-reservation=on     (+) 視聴にタイムシフト予約が必要な場合に自動的に上書きする
  -nico-force-reservation=off    (+) 自動的にタイムシフト予約しない(デフォルト)
  -nico-skip-hb=on               (+) コメント書き出し時に/hbコマンドを出さない
  -nico-skip-hb=off              (+) コメント書き出し時に/hbコマンドも出す(デフォルト)
  -nico-adjust-vpos=on           (+) コメント書き出し時にvposの値を補正する(デフォルト)
                                 vposの値が-1000より小さい場合はコメント出力しない
  -nico-adjust-vpos=off          (+) コメント書き出し時にvposの値をそのまま出力する
  -nico-ts-start <num>           タイムシフトの録画を指定した再生時間（秒）から開始する
  -nico-ts-stop <num>            タイムシフトの録画を指定した再生時間（秒）で停止する
                                 上記2つは ＜分＞:＜秒＞ | ＜時＞:＜分＞:＜秒＞ の形式でも指定可能
  -nico-ts-start-min <num>       タイムシフトの録画を指定した再生時間（分）から開始する
  -nico-ts-stop-min <num>        タイムシフトの録画を指定した再生時間（分）で停止する
                                 上記2つは ＜時＞:＜分＞ の形式でも指定可能
  -nico-conv-seqno-start <num>   MP4への変換を指定したセグメント番号から開始する
  -nico-conv-seqno-end <num>     MP4への変換を指定したセグメント番号で終了する
  -nico-conv-force-concat        MP4への変換で画質変更または抜けがあっても分割しないように設定
  -nico-conv-force-concat=on     (+) 上記を有効に設定
  -nico-conv-force-concat=off    (+) 上記を無効に設定(デフォルト)
  -nico-no-streamlink=on         (+) Streamlinkを使用しない(デフォルト)
  -nico-no-streamlink=off        (+) Streamlinkを使用する
  -nico-no-ytdlp=on              (+) yt-dlpを使用しない(デフォルト)
  -nico-no-ytdlp=off             (+) yt-dlpを使用する
  -nico-comment-only=on          (+) コメントのみダウンロードする
  -nico-comment-only=off         (+) 動画とコメントをダウンロードする(デフォルト)

ツイキャス録画用オプション:
  -tcas-retry=on                 (+) 録画終了後に再試行を行う
  -tcas-retry=off                (+) 録画終了後に再試行を行わない
  -tcas-retry-timeout            (+) 再試行を開始してから終了するまでの時間（分)
                                     -1で無限ループ。デフォルト: 5分
  -tcas-retry-interval           (+) 再試行を行う間隔（秒）デフォルト: 60秒

Youtube live録画用オプション:
  -yt-api-key <key>              (+) YouTube Data API v3 keyを設定する(未使用)
  -yt-no-streamlink=on           (+) Streamlinkを使用しない
  -yt-no-streamlink=off          (+) Streamlinkを使用する(デフォルト)
  -yt-no-youtube-dl=on           (+) yt-dlpを使用しない
  -yt-no-youtube-dl=off          (+) yt-dlpを使用する(デフォルト)
  -yt-comment-start              YouTube Liveアーカイブでコメント取得開始時間（秒）を指定
                                 ＜分＞:＜秒＞ | ＜時＞:＜分＞:＜秒＞ の形式でも指定可能
                                 0：続きからコメント取得  1：最初からコメント取得
  -yt-emoji=on                   (+) コメントにAlternate emojisを表示する(デフォルト)
  -yt-emoji=off                  (+) コメントにAlternate emojisを表示しない

変換オプション:
  -extract-chunks=off            (+) -d2mで動画ファイルに書き出す(デフォルト)
  -extract-chunks=on             (+) [上級者向] 各々のフラグメントを書き出す(大量のファイルが生成される)
  -conv-ext=mp4                  (+) -d2mで出力の拡張子を.mp4とする(デフォルト)
  -conv-ext=ts                   (+) -d2mで出力の拡張子を.tsとする

HTTP関連
  -http-skip-verify=on           (+) TLS証明書の認証をスキップする (32bit版対策)
  -http-skip-verify=off          (+) TLS証明書の認証をスキップしない (デフォルト)
  -http-timeout <num>            (+) タイムアウト時間（秒）デフォルト: 30秒（最低値）


(+)のついたオプションは、次回も同じ設定が使用されることを示す。

FILE:
  ニコニコ生放送/nicolive:
    https://live.nicovideo.jp/watch/lvXXXXXXXXX
    lvXXXXXXXXX
  ツイキャス/twitcasting:
    https://twitcasting.tv/XXXXX
`
	fmt.Printf(format, cmd, buildno.GetBuildLite(), buildno.GetBuildNo(), cmd)

	for _, b := range verbose {
		if b {
			fmt.Print(`
旧オプション:
  -conf-pass <password> [廃止] 設定ファイルのパスワード
  -z2m                  録画済みのzipをmp4に変換する(-zip-to-mp4)
  -nico-status-https    -

デバッグ用オプション:
  -nico-test-run           ニコ生テストラン
  -nico-test-timeout <num> ニコ生テストランでの各放送のタイムアウト
  -nico-test-format        フォーマット、保存しない
  -nico-ufast-ts           TS保存にウェイトを入れない
  -nico-debug              デバッグ用ログ出力する

HTTP関連
  -http-root-ca <file>    ルート証明書ファイルを指定(pem/der)
  -http-skip-verify       TLS証明書の認証をスキップする
  -http-proxy <proxy url> [警告] proxyを設定する
[警告] 情報流出に注意。信頼できるproxy serverのみに使用すること。

`)
			break
		}
	}

	os.Exit(0)
}

func dbConfSet(db *sql.DB, k string, v interface{}) {
	query := `INSERT OR REPLACE INTO conf (k,v) VALUES (?,?)`

	if _, err := db.Exec(query, k, v); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func SetNicoLogin(hash, user, pass string) (err error) {
	db, err := dbAccountOpen()
	if err != nil {
		if db != nil {
			db.Close()
		}
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT OR IGNORE INTO niconico (alias, user, pass) VALUES(?, ?, ?);
		UPDATE niconico SET user = ?, pass = ? WHERE alias = ?
	`, hash, user, pass, user, pass, hash)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("niconico account saved.\n")
	return
}

func SetNicoSession(hash, session string) (err error) {
	db, err := dbAccountOpen()
	if err != nil {
		if db != nil {
			db.Close()
		}
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT OR IGNORE INTO niconico (alias, session) VALUES(?, ?);
		UPDATE niconico SET session = ? WHERE alias = ?
	`, hash, session, session, hash)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
func LoadNicoAccount(alias string) (user, pass, session string, err error) {
	db, err := dbAccountOpen()
	if err != nil {
		if db != nil {
			db.Close()
		}
		return
	}
	defer db.Close()

	db.QueryRow(`SELECT user, pass, IFNULL(session, "") FROM niconico WHERE alias = ?`, alias).Scan(&user, &pass, &session)
	return
}
func SetYoutubeApiKey(key string) (err error) {
	db, err := dbAccountOpen()
	if err != nil {
		if db != nil {
			db.Close()
		}
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT OR IGNORE INTO youtubeapikey (id, key) VALUES(1, ?);
		UPDATE youtubeapikey SET key = ? WHERE id = 1
	`, key, key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Youtube API KEY saved.\n")
	return
}
func LoadYoutubeApiKey() (key string, err error) {
	db, err := dbAccountOpen()
	if err != nil {
		if db != nil {
			db.Close()
		}
		return
	}
	defer db.Close()

	db.QueryRow(`SELECT IFNULL(key, "") FROM youtubeapikey WHERE id = 1`).Scan(&key)
	if key == "" {
		err = fmt.Errorf("apikey not found")
	}
	return
}
func dbAccountOpen() (db *sql.DB, err error) {

	base := func() string {
		if b := os.Getenv("LIVEDL_DIR"); b != "" {
			return b
		}
		if b := os.Getenv("APPDATA"); b != "" {
			return fmt.Sprintf("%s/livedl", b)
		}
		if b := os.Getenv("HOME"); b != "" {
			return fmt.Sprintf("%s/.livedl", b)
		}
		return ""
	}()
	if base == "" {
		log.Fatalln("basedir for account not defined")
	}

	name := fmt.Sprintf("%s/account.db", base)
	files.MkdirByFileName(name)
	db, err = sql.Open("sqlite3", name)
	if err != nil {
		log.Println(err)
		return
	}

	// niconico
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS niconico (
		alias TEXT PRIMARY KEY NOT NULL UNIQUE,
		user TEXT NOT NULL,
		pass TEXT NOT NULL,
		session TEXT
	)
	`)
	if err != nil {
		return
	}

	_, err = db.Exec(`
	CREATE UNIQUE INDEX IF NOT EXISTS niconico0 ON niconico(alias);
	CREATE UNIQUE INDEX IF NOT EXISTS niconico1 ON niconico(user);
	`)
	if err != nil {
		return
	}

	// youtube API key
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS youtubeapikey (
		id PRIMARY KEY NOT NULL UNIQUE,
		key TEXT
	)
	`)
	if err != nil {
		return
	}

	_, err = db.Exec(`
	CREATE UNIQUE INDEX IF NOT EXISTS youtubeapikey0 ON youtubeapikey(id);
	`)
	if err != nil {
		return
	}

	return
}

func dbOpen() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", "conf.db")
	if err != nil {
		return
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS conf (
		k TEXT PRIMARY KEY NOT NULL UNIQUE,
		v BLOB
	)
	`)
	if err != nil {
		return
	}

	_, err = db.Exec(`
	CREATE UNIQUE INDEX IF NOT EXISTS conf0 ON conf(k);
	`)
	if err != nil {
		return
	}
	return
}

func GetBrowserName(str string) (name string) {
	name = "error"
	if len(str) <= 0 {
		return
	}
	if m := regexp.MustCompile(`^(firefox:?['\"]?.*['\"]?)`).FindStringSubmatch(str); len(m) > 0 {
		name = m[1]
	}
	return
}

func parseTime(arg string) (ret int64, err error) {
	var hour, min, sec int

	if m := regexp.MustCompile(`^(\d+):(\d+):(\d+)$`).FindStringSubmatch(arg); len(m) > 0 {
		hour, err = strconv.Atoi(m[1])
		if err != nil {
			return
		}
		min, err = strconv.Atoi(m[2])
		if err != nil {
			return
		}
		sec, err = strconv.Atoi(m[3])
		if err != nil {
			return
		}
	} else if m := regexp.MustCompile(`^(\d+):(\d+)$`).FindStringSubmatch(arg); len(m) > 0 {
		min, err = strconv.Atoi(m[1])
		if err != nil {
			return
		}
		sec, err = strconv.Atoi(m[2])
		if err != nil {
			return
		}
	} else if m := regexp.MustCompile(`^(\d+)$`).FindStringSubmatch(arg); len(m) > 0 {
		sec, err = strconv.Atoi(m[1])
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("regexp not matched")
	}

	ret = int64(hour * 3600 + min * 60 + sec)

	return
}
func SecondsToHHMMSS(seconds int64) string {
	// 時間、分、秒を計算
	if seconds <= 0 {
		return ""
	}
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	// フォーマットされた文字列を返す
	return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
}

func ParseArgs() (opt Option) {
	//dbAccountOpen()
	db, err := dbOpen()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.QueryRow(`
		SELECT
		IFNULL((SELECT v FROM conf WHERE k == "NicoFormat"), ""),
		IFNULL((SELECT v FROM conf WHERE k == "NicoLimitBw"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "NicoLoginOnly"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "NicoFastTs"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "NicoLoginAlias"), ""),
		IFNULL((SELECT v FROM conf WHERE k == "NicoAutoConvert"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "NicoAutoDeleteDBMode"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "TcasRetry"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "TcasRetryTimeoutMinute"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "TcasRetryInterval"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "ConvExt"), ""),
		IFNULL((SELECT v FROM conf WHERE k == "ExtractChunks"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "NicoConvForceConcat"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "NicoForceResv"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "YtNoStreamlink"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "YtNoYoutubeDl"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "YtEmoji"), 1),
		IFNULL((SELECT v FROM conf WHERE k == "NicoSkipHb"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "NicoAdjustVpos"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "HttpSkipVerify"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "HttpTimeout"), 30),
		IFNULL((SELECT v FROM conf WHERE k == "NicoNoStreamlink"), 1),
		IFNULL((SELECT v FROM conf WHERE k == "NicoNoYtdlp"), 1),
		IFNULL((SELECT v FROM conf WHERE k == "NicoCommentOnly"), 0),
		IFNULL((SELECT v FROM conf WHERE k == "NicoExecBw"), "");
	`).Scan(
		&opt.NicoFormat,
		&opt.NicoLimitBw,
		&opt.NicoLoginOnly,
		&opt.NicoFastTs,
		&opt.NicoLoginAlias,
		&opt.NicoAutoConvert,
		&opt.NicoAutoDeleteDBMode,
		&opt.TcasRetry,
		&opt.TcasRetryTimeoutMinute,
		&opt.TcasRetryInterval,
		&opt.ConvExt,
		&opt.ExtractChunks,
		&opt.NicoConvForceConcat,
		&opt.NicoForceResv,
		&opt.YtNoStreamlink,
		&opt.YtNoYoutubeDl,
		&opt.YtEmoji,
		&opt.NicoSkipHb,
		&opt.NicoAdjustVpos,
		&opt.HttpSkipVerify,
		&opt.HttpTimeout,
		&opt.NicoNoStreamlink,
		&opt.NicoNoYtdlp,
		&opt.NicoCommentOnly,
		&opt.NicoExecBw,
	)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	args := os.Args[1:]
	var match []string

	type Parser struct {
		re *regexp.Regexp
		cb func() error
	}

	nextArg := func() (str string, err error) {
		if len(args) <= 0 {
			if len(match[0]) > 0 {
				err = fmt.Errorf("%v: value required", match[0])
			} else {
				err = fmt.Errorf("value required")
			}
		} else {
			str = args[0]
			args = args[1:]
		}

		return
	}

	parseList := []Parser{
		Parser{regexp.MustCompile(`\A(?i)(?:--?|/)(?:\?|h|help)\z`), func() error {
			Help()
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)(?:--?|/)v(?:\?|h|help)\z`), func() error {
			Help(true)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?(?:v|version)\z`), func() error {
			version()
			return nil
		}},
		Parser{regexp.MustCompile(`\A(https?://(?:[^/]*@)?(?:[^/]*\.)*nicovideo\.jp(?::[^/]*)?/(?:[^/]*?/)*)?(lv\d+)(?:\?.*)?\z`), func() error {
			switch opt.Command {
			default:
				fmt.Printf("Use \"--\" option for FILE for %s\n", opt.Command)
				Help()
			case "", "NICOLIVE":
				opt.NicoLiveId = match[2]
				opt.Command = "NICOLIVE"
			case "NICOLIVE_TEST":
				opt.NicoLiveId = match[2]
			}
			return nil
		}},
		Parser{regexp.MustCompile(`\A--?conf-?pass\z`), func() (err error) {
			str, err := nextArg()
			if err != nil {
				return
			}
			opt.ConfPass = str
			return
		}},
		Parser{regexp.MustCompile(`\Ahttps?://twitcasting\.tv/([^/]+)(?:/.*)?\z`), func() error {
			opt.TcasId = match[1]
			opt.Command = "TWITCAS"
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?tcas-?retry(?:=(on|off))\z`), func() error {
			if strings.EqualFold(match[1], "on") {
				opt.TcasRetry = true
			} else if strings.EqualFold(match[1], "off") {
				opt.TcasRetry = false
			}
			dbConfSet(db, "TcasRetry", opt.TcasRetry)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?tcas-?retry-?timeout(?:-?minutes?)?\z`), func() error {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("--tcas-retry-timeout: Not a number: %s\n", s)
			}
			opt.TcasRetryTimeoutMinute = num
			dbConfSet(db, "TcasRetryTimeoutMinute", opt.TcasRetryTimeoutMinute)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?tcas-?retry-?interval\z`), func() error {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("--tcas-retry-interval: Not a number: %s\n", s)
			}
			if num <= 0 {
				return fmt.Errorf("--tcas-retry-interval: Invalid: %d: greater than 1\n", num)
			}

			opt.TcasRetryInterval = num
			dbConfSet(db, "TcasRetryInterval", opt.TcasRetryInterval)
			return nil
		}},
		Parser{regexp.MustCompile(`\Ahttps?://(?:[^/]*\.)*youtube\.com/(?:.*\W)?v=([\w-]+)(?:[^\w-].*)?\z`), func() error {
			opt.YoutubeId = match[1]
			opt.Command = "YOUTUBE"
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico\z`), func() error {
			opt.Command = "NICOLIVE"
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?test-?run\z`), func() error {
			opt.Command = "NICOLIVE_TEST"
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?test-?timeout\z`), func() error {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("--nico-test-timeout: Not a number: %s\n", s)
			}
			if num <= 0 {
				return fmt.Errorf("--nico-test-timeout: Invalid: %d: must be greater than or equal to 1\n", num)
			}
			opt.NicoTestTimeout = num
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?tcas\z`), func() error {
			opt.Command = "TWITCAS"
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?(?:yt|youtube|youtube-live)\z`), func() error {
			opt.Command = "YOUTUBE"
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?(?:z|zip)-?(?:2|to)-?(?:m|mp4)\z`), func() error {
			opt.Command = "ZIP2MP4"
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?(?:d|db|sqlite3?)-?(?:2|to)-?(?:m|mp4)\z`), func() error {
			opt.Command = "DB2MP4"
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?(?:d|db|sqlite3?)-?(?:2|to)-?(?:h|hls)\z`), func() error {
			opt.Command = "DB2HLS"
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?(?:d|db|sqlite3?)-?(?:i|info)\z`), func() error {
			opt.Command = "DBINFO"
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?login-?only(?:=(on|off))?\z`), func() error {
			if strings.EqualFold(match[1], "on") {
				opt.NicoLoginOnly = true
				dbConfSet(db, "NicoLoginOnly", opt.NicoLoginOnly)
			} else if strings.EqualFold(match[1], "off") {
				opt.NicoLoginOnly = false
				dbConfSet(db, "NicoLoginOnly", opt.NicoLoginOnly)
			} else {
				opt.NicoLoginOnly = true
			}
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?fast-?ts(?:=(on|off))?\z`), func() error {
			if strings.EqualFold(match[1], "on") {
				opt.NicoFastTs = true
				dbConfSet(db, "NicoFastTs", opt.NicoFastTs)
			} else if strings.EqualFold(match[1], "off") {
				opt.NicoFastTs = false
				dbConfSet(db, "NicoFastTs", opt.NicoFastTs)
			} else {
				opt.NicoFastTs = true
			}
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?auto-?convert(?:=(on|off))?\z`), func() error {
			if strings.EqualFold(match[1], "on") {
				opt.NicoAutoConvert = true
				dbConfSet(db, "NicoAutoConvert", opt.NicoAutoConvert)
			} else if strings.EqualFold(match[1], "off") {
				opt.NicoAutoConvert = false
				dbConfSet(db, "NicoAutoConvert", opt.NicoAutoConvert)
			} else {
				opt.NicoAutoConvert = true
			}
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?auto-?delete-?mode\z`), func() error {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("--nico-auto-delete-mode: Not a number: %s\n", s)
			}
			if num < 0 || 2 < num {
				return fmt.Errorf("--nico-auto-delete-mode: Invalid: %d: one of 0, 1, 2\n", num)
			}

			opt.NicoAutoDeleteDBMode = num
			dbConfSet(db, "NicoAutoDeleteDBMode", opt.NicoAutoDeleteDBMode)

			return nil
		}},

		Parser{regexp.MustCompile(`\A(?i)--?nico-?(?:u|ultra)fast-?ts\z`), func() error {
			opt.NicoUltraFastTs = true
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?status-?https\z`), func() error {
			// experimental
			opt.NicoStatusHTTPS = true
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?hls-?port\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("--nico-hls-port: Not a number: %s\n", s)
			}
			if num <= 0 {
				return fmt.Errorf("--nico-hls-port: Invalid: %d: must be greater than or equal to 1\n", num)
			}
			opt.NicoHlsPort = num
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?limit-?bw\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			if m := regexp.MustCompile(`^(audio_only|audio_high)$`).FindStringSubmatch(s); len(m) > 0 {
				opt.NicoLimitBw = m[0]
				dbConfSet(db, "NicoLimitBw", opt.NicoLimitBw)
				return nil
			}
			_, err = strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("--nico-limit-bw: Not a number: %s\n", s)
			}
			opt.NicoLimitBw = s
			dbConfSet(db, "NicoLimitBw", opt.NicoLimitBw)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?ts-?start\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := parseTime(s)
			if err != nil {
				return fmt.Errorf("--nico-ts-start: Not a number %s\n", s)
			}
			//opt.NicoTsStart = float64(num)
			opt.NicoTsStart = num
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?ts-?start-?min\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := parseTime(s + ":0")
			if err != nil {
				return fmt.Errorf("--nico-ts-start-min: Not a number %s\n", s)
			}
			//opt.NicoTsStart = float64(num)
			opt.NicoTsStart = num
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?ts-?stop\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := parseTime(s)
			if err != nil {
				return fmt.Errorf("--nico-ts-stop: Not a number %s\n", s)
			}
			opt.NicoTsStop = num
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?ts-?stop-?min\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := parseTime(s + ":0")
			if err != nil {
				return fmt.Errorf("--nico-ts-stop-min: Not a number %s\n", s)
			}
			opt.NicoTsStop = num
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?(?:format|fmt)\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			if s == "" {
				return fmt.Errorf("--nico-format: null string not allowed\n", s)
			}
			opt.NicoFormat = s
			dbConfSet(db, "NicoFormat", opt.NicoFormat)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?exec-?bw\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			opt.NicoExecBw = s
			dbConfSet(db, "NicoExecBw", opt.NicoExecBw)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?test-?(?:format|fmt)\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			if s == "" {
				return fmt.Errorf("--nico-test-format: null string not allowed\n", s)
			}
			opt.NicoFormat = s
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?login\z`), func() (err error) {
			str, err := nextArg()
			if err != nil {
				return
			}
			ar := strings.SplitN(str, ",", 2)
			if len(ar) >= 2 && ar[0] != "" {
				loginId := ar[0]
				loginPass := ar[1]
				opt.NicoLoginAlias = fmt.Sprintf("%x", sha3.Sum256([]byte(loginId)))
				SetNicoLogin(opt.NicoLoginAlias, loginId, loginPass)
				dbConfSet(db, "NicoLoginAlias", opt.NicoLoginAlias)
				opt.NicoLoginOnly = true

			} else {
				return fmt.Errorf("--nico-login: <id>,<password>")
			}
			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?cookies?\z`), func() (err error) {
			str, err := nextArg()
			if err != nil {
				return
			}
			str = GetBrowserName(str)
			if str != "error" {
				opt.NicoCookies = str
			} else {
				return fmt.Errorf("--nico-cookies: invalid browser name")
			}
			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?session\z`), func() (err error) {
			str, err := nextArg()
			if err != nil {
				return
			}
			opt.NicoSession = str
			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?load-?session\z`), func() (err error) {
			name, err := nextArg()
			if err != nil {
				return
			}
			b, err := ioutil.ReadFile(name)
			if err != nil {
				return
			}
			if ma := regexp.MustCompile(`(\S+)`).FindSubmatch(b); len(ma) > 0 {
				opt.NicoSession = string(ma[1])
			} else {
				err = fmt.Errorf("--nico-load-session: load failured")
			}

			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?debug\z`), func() error {
			opt.NicoDebug = true
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i).+\.zip\z`), func() (err error) {
			switch opt.Command {
			case "", "ZIP2MP4":
				opt.Command = "ZIP2MP4"
				opt.ZipFile = match[0]
			default:
				return fmt.Errorf("%s: Use -- option before \"%s\"", opt.Command, match[0])
			}
			return
		}},
		Parser{regexp.MustCompile(`\A(?i).+\.sqlite3\z`), func() (err error) {
			switch opt.Command {
			case "", "DB2MP4":
				opt.Command = "DB2MP4"
				opt.DBFile = match[0]
			case "DB2HLS":
				opt.DBFile = match[0]
			default:
				return fmt.Errorf("%s: Use -- option before \"%s\"", opt.Command, match[0])
			}
			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?conv-?ext(?:=(mp4|ts))\z`), func() error {
			if strings.EqualFold(match[1], "mp4") {
				opt.ConvExt = "mp4"
			} else if strings.EqualFold(match[1], "ts") {
				opt.ConvExt = "ts"
			}
			dbConfSet(db, "ConvExt", opt.ConvExt)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?extract(?:-?chunks)?(?:=(on|off))\z`), func() error {
			if strings.EqualFold(match[1], "on") {
				opt.ExtractChunks = true
			} else if strings.EqualFold(match[1], "off") {
				opt.ExtractChunks = false
			}
			dbConfSet(db, "ExtractChunks", opt.ExtractChunks)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?conv-?force-?concat(?:=(on|off))?\z`), func() error {
			if strings.EqualFold(match[1], "on") {
				opt.NicoConvForceConcat = true
				dbConfSet(db, "NicoConvForceConcat", opt.NicoConvForceConcat)
			} else if strings.EqualFold(match[1], "off") {
				opt.NicoConvForceConcat = false
				dbConfSet(db, "NicoConvForceConcat", opt.NicoConvForceConcat)
			} else {
				opt.NicoConvForceConcat = true
			}
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?conv-?seqno-?start\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("--nico-conv-seqno-start: Not a number: %s\n", s)
			}
			if num < 0 {
				return fmt.Errorf("--nico-conv-seqno-start: Invalid: %d: must be greater than or equal to 0\n", num)
			}
			opt.NicoConvSeqnoStart = int64(num)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?conv-?seqno-?end\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("--nico-conv-seqno-end: Not a number: %s\n", s)
			}
			if num < 0 {
				return fmt.Errorf("--nico-conv-seqno-end: Invalid: %d: must be greater than or equal to 0\n", num)
			}
			opt.NicoConvSeqnoEnd = int64(num)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?force-?(?:re?sv|reservation)(?:=(on|off))\z`), func() error {
			if strings.EqualFold(match[1], "on") {
				opt.NicoForceResv = true
			} else if strings.EqualFold(match[1], "off") {
				opt.NicoForceResv = false
			}
			dbConfSet(db, "NicoForceResv", opt.NicoForceResv)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?yt-?api-?key\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return
			}
			if s == "" {
				return fmt.Errorf("--yt-api-key: null string not allowed\n", s)
			}
			err = SetYoutubeApiKey(s)
			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?yt-?no-?streamlink(?:=(on|off))?\z`), func() (err error) {
			if strings.EqualFold(match[1], "on") {
				opt.YtNoStreamlink = true
				dbConfSet(db, "YtNoStreamlink", opt.YtNoStreamlink)
			} else if strings.EqualFold(match[1], "off") {
				opt.YtNoStreamlink = false
				dbConfSet(db, "YtNoStreamlink", opt.YtNoStreamlink)
			} else {
				opt.YtNoStreamlink = true
			}
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?yt-?no-?youtube-?dl(?:=(on|off))?\z`), func() (err error) {
			if strings.EqualFold(match[1], "on") {
				opt.YtNoYoutubeDl = true
				dbConfSet(db, "YtNoYoutubeDl", opt.YtNoYoutubeDl)
			} else if strings.EqualFold(match[1], "off") {
				opt.YtNoYoutubeDl = false
				dbConfSet(db, "YtNoYoutubeDl", opt.YtNoYoutubeDl)
			} else {
				opt.YtNoYoutubeDl = true
			}
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?yt-?emoji(?:=(on|off))?\z`), func() (err error) {
			if strings.EqualFold(match[1], "on") {
				opt.YtEmoji = true
				dbConfSet(db, "YtEmoji", opt.YtEmoji)
			} else if strings.EqualFold(match[1], "off") {
				opt.YtEmoji = false
				dbConfSet(db, "YtEmoji", opt.YtEmoji)
			} else {
				opt.YtEmoji = true
			}
			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?yt-?comment-?start\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := parseTime(s)
			if err != nil {
				return fmt.Errorf("--yt-comment-start: Not a number %s\n", s)
			}
			opt.YtCommentStart = float64(num)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?skip-?hb(?:=(on|off))?\z`), func() (err error) {
			if strings.EqualFold(match[1], "on") {
				opt.NicoSkipHb = true
				dbConfSet(db, "NicoSkipHb", opt.NicoSkipHb)
			} else if strings.EqualFold(match[1], "off") {
				opt.NicoSkipHb = false
				dbConfSet(db, "NicoSkipHb", opt.NicoSkipHb)
			} else {
				opt.NicoSkipHb = true
			}
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?http-?root-?ca\z`), func() (err error) {
			str, err := nextArg()
			if err != nil {
				return
			}
			opt.HttpRootCA = str
			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?http-?skip-?verify(?:=(on|off))?\z`), func() (err error) {
			if strings.EqualFold(match[1], "on") {
				opt.HttpSkipVerify = true
				dbConfSet(db, "HttpSkipVerify", opt.HttpSkipVerify)
			} else if strings.EqualFold(match[1], "off") {
				opt.HttpSkipVerify = false
				dbConfSet(db, "HttpSkipVerify", opt.HttpSkipVerify)
			} else {
				opt.HttpSkipVerify = true
			}

			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?http-?proxy\z`), func() (err error) {
			str, err := nextArg()
			if err != nil {
				return
			}
			if !strings.Contains(str, "://") {
				str = "http://" + str
			}
			opt.HttpProxy = str
			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?no-?chdir\z`), func() (err error) {
			opt.NoChdir = true
			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?http-?timeout\z`), func() (err error) {
			s, err := nextArg()
			if err != nil {
				return err
			}
			num, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("--http-timeout: Not a number: %s\n", s)
			}
			if num < MinimumHttpTimeout {
				return fmt.Errorf("--http-timeout: Invalid: %d: must be greater than or equal to %#v\n", num, MinimumHttpTimeout)
			}
			opt.HttpTimeout = num
			dbConfSet(db, "HttpTimeout", opt.HttpTimeout)
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?adjust-?vpos(?:=(on|off))?\z`), func() (err error) {
			if strings.EqualFold(match[1], "on") {
				opt.NicoAdjustVpos = true
				dbConfSet(db, "NicoAdjustVpos", opt.NicoAdjustVpos)
			} else if strings.EqualFold(match[1], "off") {
				opt.NicoAdjustVpos = false
				dbConfSet(db, "NicoAdjustVpos", opt.NicoAdjustVpos)
			} else {
				opt.NicoAdjustVpos = true
			}
			return
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?no-?streamlink(?:=(on|off))?\z`), func() (err error) {
			if strings.EqualFold(match[1], "on") {
				opt.NicoNoStreamlink = true
				dbConfSet(db, "NicoNoStreamlink", opt.NicoNoStreamlink)
			} else if strings.EqualFold(match[1], "off") {
				opt.NicoNoStreamlink = false
				dbConfSet(db, "NicoNoStreamlink", opt.NicoNoStreamlink)
			} else {
				opt.NicoNoStreamlink = false
			}
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?no-?ytdlp(?:=(on|off))?\z`), func() (err error) {
			if strings.EqualFold(match[1], "on") {
				opt.NicoNoYtdlp = true
				dbConfSet(db, "NicoNoYtdlp", opt.NicoNoYtdlp)
			} else if strings.EqualFold(match[1], "off") {
				opt.NicoNoYtdlp = false
				dbConfSet(db, "NicoNoYtdlp", opt.NicoNoYtdlp)
			} else {
				opt.NicoNoYtdlp = false
			}
			return nil
		}},
		Parser{regexp.MustCompile(`\A(?i)--?nico-?comment-?only(?:=(on|off))?\z`), func() (err error) {
			if strings.EqualFold(match[1], "on") {
				opt.NicoCommentOnly = true
				dbConfSet(db, "NicoCommentOnly", opt.NicoCommentOnly)
			} else if strings.EqualFold(match[1], "off") {
				opt.NicoCommentOnly = false
				dbConfSet(db, "NicoCommentOnly", opt.NicoCommentOnly)
			} else {
				opt.NicoCommentOnly = false
			}
			return nil
		}},
	}

	checkFILE := func(arg string) bool {
		switch opt.Command {
		default:
			//fmt.Printf("command not specified: -- \"%s\"\n", arg)
			//os.Exit(1)
		case "YOUTUBE":
			if ma := regexp.MustCompile(`v=([\w-]+)`).FindStringSubmatch(arg); len(ma) > 0 {
				opt.YoutubeId = ma[1]
				return true
			} else if ma := regexp.MustCompile(`\A([\w-]+)\z`).FindStringSubmatch(arg); len(ma) > 0 {
				opt.YoutubeId = ma[1]
				return true
			} else {
				fmt.Printf("Not YouTube id: %s\n", arg)
				os.Exit(1)
			}
		case "NICOLIVE":
			if ma := regexp.MustCompile(`(lv\d+)`).FindStringSubmatch(arg); len(ma) > 0 {
				opt.NicoLiveId = ma[1]
				return true
			}
		case "TWITCAS":
			if opt.TcasId != "" {
				fmt.Printf("Unknown option: %s\n", arg)
				Help()
			}
			if ma := regexp.MustCompile(`(?:.*/)?([^/]+)\z`).FindStringSubmatch(arg); len(ma) > 0 {
				opt.TcasId = ma[1]
				return true
			}
		case "ZIP2MP4":
			if ma := regexp.MustCompile(`(?i)\.zip`).FindStringSubmatch(arg); len(ma) > 0 {
				opt.ZipFile = arg
				return true
			}
		case "DB2MP4", "DB2HLS", "DBINFO":
			if ma := regexp.MustCompile(`(?i)\.sqlite3`).FindStringSubmatch(arg); len(ma) > 0 {
				opt.DBFile = arg
				return true
			}
			return false
		} // end switch
		return false
	}

LB_ARG:
	for len(args) > 0 {
		arg, _ := nextArg()

		if arg == "--" {
			switch len(args) {
			case 0:
				fmt.Printf("argument not specified after \"--\"\n")
				os.Exit(1)
			default:
				fmt.Printf("too many arguments after \"--\": %v\n", args)
				os.Exit(1)
			case 1:
				arg, _ := nextArg()
				checkFILE(arg)
			}

		} else {
			for _, p := range parseList {
				if match = p.re.FindStringSubmatch(arg); len(match) > 0 {
					if e := p.cb(); e != nil {
						fmt.Println(e)
						os.Exit(1)
					}
					continue LB_ARG
				}
			}
			if ok := checkFILE(arg); !ok {
				fmt.Printf("Unknown option: %v\n", arg)
				Help()
			}
		}
	}

	if opt.ConfFile == "" {
		opt.ConfFile = fmt.Sprintf("%s.conf", getCmd())
	}

	if opt.HttpTimeout == 0 {
		opt.HttpTimeout = MinimumHttpTimeout
	}
	httpbase.SetTimeout(opt.HttpTimeout)

	// [deprecated]
	// load session info
	if data, e := cryptoconf.Load(opt.ConfFile, opt.ConfPass); e != nil {
		err = e
		return
	} else {
		loginId, _ := data["NicoLoginId"].(string)
		if loginId != "" {
			loginPass, _ := data["NicoLoginPass"].(string)
			hash := fmt.Sprintf("%x", sha3.Sum256([]byte(loginId)))
			SetNicoLogin(hash, loginId, loginPass)
			if opt.NicoLoginAlias == "" {
				opt.NicoLoginAlias = hash
				dbConfSet(db, "NicoLoginAlias", opt.NicoLoginAlias)
			}
			os.Remove(opt.ConfFile)
		}
	}

	// prints
	switch opt.Command {
	case "NICOLIVE":
		fmt.Printf("Conf(NicoLoginOnly): %#v\n", opt.NicoLoginOnly)
		fmt.Printf("Conf(NicoFormat): %#v\n", opt.NicoFormat)
		fmt.Printf("Conf(NicoLimitBw): %#v\n", opt.NicoLimitBw)
		fmt.Printf("Conf(NicoExecBw): %#v\n", opt.NicoExecBw)
		fmt.Printf("Conf(NicoFastTs): %#v\n", opt.NicoFastTs)
		fmt.Printf("Conf(NicoAutoConvert): %#v\n", opt.NicoAutoConvert)
		if opt.NicoAutoConvert {
			fmt.Printf("Conf(NicoAutoDeleteDBMode): %#v\n", opt.NicoAutoDeleteDBMode)
			fmt.Printf("Conf(ExtractChunks): %#v\n", opt.ExtractChunks)
			fmt.Printf("Conf(NicoConvForceConcat): %#v\n", opt.NicoConvForceConcat)
			fmt.Printf("Conf(ConvExt): %#v\n", opt.ConvExt)
		}
		fmt.Printf("Conf(NicoForceResv): %#v\n", opt.NicoForceResv)
		//fmt.Printf("Conf(NicoSkipHb): %#v\n", opt.NicoSkipHb)
		fmt.Printf("Conf(NicoAdjustVpos): %#v\n", opt.NicoAdjustVpos)
		fmt.Printf("Conf(NicoNoStreamlink): %#v\n", opt.NicoNoStreamlink)
		fmt.Printf("Conf(NicoNoYtdlp): %#v\n", opt.NicoNoYtdlp)
		fmt.Printf("Conf(NicoCommentOnly): %#v\n", opt.NicoCommentOnly)

	case "YOUTUBE":
		fmt.Printf("Conf(YtNoStreamlink): %#v\n", opt.YtNoStreamlink)
		fmt.Printf("Conf(YtNoYoutubeDl): %#v\n", opt.YtNoYoutubeDl)
		fmt.Printf("Conf(YtEmoji): %#v\n", opt.YtEmoji)

	case "TWITCAS":
		fmt.Printf("Conf(TcasRetry): %#v\n", opt.TcasRetry)
		fmt.Printf("Conf(TcasRetryTimeoutMinute): %#v\n", opt.TcasRetryTimeoutMinute)
		fmt.Printf("Conf(TcasRetryInterval): %#v\n", opt.TcasRetryInterval)
	case "DB2MP4", "DBINFO":
		fmt.Printf("Conf(ExtractChunks): %#v\n", opt.ExtractChunks)
		fmt.Printf("Conf(NicoConvForceConcat): %#v\n", opt.NicoConvForceConcat)
		fmt.Printf("Conf(ConvExt): %#v\n", opt.ConvExt)
		//fmt.Printf("Conf(NicoSkipHb): %#v\n", opt.NicoSkipHb)
		fmt.Printf("Conf(NicoAdjustVpos): %#v\n", opt.NicoAdjustVpos)
		fmt.Printf("Conf(YtEmoji): %#v\n", opt.YtEmoji)
	case "DB2HLS":
		fmt.Printf("Conf(NicoHlsPort): %#v\n", opt.NicoHlsPort)
		fmt.Printf("Conf(NicoConvSeqnoStart): %#v\n", opt.NicoConvSeqnoStart)
	}
	fmt.Printf("Conf(HttpSkipVerify): %#v\n", opt.HttpSkipVerify)
	fmt.Printf("Conf(HttpTimeout): %#v\n", opt.HttpTimeout)

	// check
	switch opt.Command {
	case "":
		fmt.Printf("Command not specified\n")
		Help()
	case "YOUTUBE":
		if opt.YoutubeId == "" {
			Help()
		}
	case "NICOLIVE":
		if opt.NicoLiveId == "" {
			Help()
		}
	case "NICOLIVE_TEST":
	case "TWITCAS":
		if opt.TcasId == "" {
			Help()
		}
	case "ZIP2MP4":
		if opt.ZipFile == "" {
			Help()
		}
	case "DB2MP4", "DB2HLS", "DBINFO":
		if opt.DBFile == "" {
			Help()
		}
	default:
		fmt.Printf("[FIXME] options.go/argcheck for %s\n", opt.Command)
		os.Exit(1)
	}

	return
}
