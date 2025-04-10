package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"github.com/nnn-revo2012/livedl/options"
	"github.com/nnn-revo2012/livedl/twitcas"
	"github.com/nnn-revo2012/livedl/niconico"
	"github.com/nnn-revo2012/livedl/youtube"
	"github.com/nnn-revo2012/livedl/zip2mp4"
	"time"
	"strings"
	"github.com/nnn-revo2012/livedl/httpbase"
)

func main() {
	var baseDir string
	if regexp.MustCompile(`\AC:\\.*\\Temp\\go-build[^\\]*\\[^\\]+\\exe\\[^\\]*\.exe\z`).MatchString(os.Args[0]) {
		// go runで起動時
		pwd, e := os.Getwd()
		if e != nil {
			fmt.Println(e)
			return
		}
		baseDir = pwd
	} else {
		//pa, e := filepath.Abs(os.Args[0])
		pa, e := os.Executable()
		if e != nil {
			fmt.Println(e)
			return
		}

		// symlinkを追跡する
		for {
			sl, e := os.Readlink(pa)
			if e != nil {
				break
			}
			pa = sl
		}
		baseDir = filepath.Dir(pa)
	}

	// check option only -no-chdir
	args := os.Args[1:]
	nochdir := false
	r := regexp.MustCompile(`\A(?i)--?no-?chdir\z`)
	for _, s := range args {
		if r.MatchString(s) {
			nochdir = true
			break
		}
	}

	// chdir if not disabled
	if !nochdir {
		fmt.Printf("chdir: %s\n", baseDir)
		if e := os.Chdir(baseDir); e != nil {
			fmt.Println(e)
			return
		}
	} else {
		fmt.Printf("no chdir\n")
		pwd, e := os.Getwd()
		if e != nil {
			fmt.Println(e)
			return
		}
		fmt.Printf("read %s\n", filepath.FromSlash(pwd+"/conf.db"))
	}

	opt := options.ParseArgs()

	// http
	httpbase.SetTransport()

	if opt.HttpRootCA != "" {
		if err := httpbase.SetRootCA(opt.HttpRootCA); err != nil {
			fmt.Println(err)
			return
		}
	}
	if opt.HttpSkipVerify {
		if err := httpbase.SetSkipVerify(true); err != nil {
			fmt.Println(err)
			return
		}
	}
	if opt.HttpProxy != "" {
		if err := httpbase.SetProxy(opt.HttpProxy); err != nil {
			fmt.Println(err)
			return
		}
	}

	switch opt.Command {
	default:
		fmt.Printf("Unknown command: %v\n", opt.Command)
		os.Exit(1)

	case "TWITCAS":
		var doneTime int64
		for {
			done, dbLocked := twitcas.TwitcasRecord(opt.TcasId, "")
			if dbLocked {
				break
			}
			if (! opt.TcasRetry) {
				break
			}

			if opt.TcasRetryTimeoutMinute < 0 {

			} else if done {
				doneTime = time.Now().Unix()

			} else {
				if doneTime == 0 {
					doneTime = time.Now().Unix()
				} else {
					delta := time.Now().Unix() - doneTime
					var minutes int
					if opt.TcasRetryTimeoutMinute == 0 {
						minutes = options.DefaultTcasRetryTimeoutMinute
					} else {
						minutes = opt.TcasRetryTimeoutMinute
					}

					if minutes > 0 {
						if delta > int64(minutes * 60) {
							break
						}
					}
				}
			}

			var interval int
			if opt.TcasRetryInterval <= 0 {
				interval = options.DefaultTcasRetryInterval
			} else {
				interval = opt.TcasRetryInterval
			}
			select {
			case <-time.After(time.Duration(interval) * time.Second):
			}
		}

	case "YOUTUBE":
		err := youtube.Record(opt.YoutubeId, opt.YtNoStreamlink, opt.YtNoYoutubeDl, opt.YtCommentStart)
		if err != nil {
			fmt.Println(err)
		}

	case "NICOLIVE":
		hlsPlaylistEnd, dbname, err := niconico.Record(opt);
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if  ((!opt.NicoNoStreamlink || !opt.NicoNoYtdlp) && !opt.NicoCommentOnly) || hlsPlaylistEnd && opt.NicoAutoConvert {
			done, nMp4s, skipped, err := zip2mp4.ConvertDB(dbname, opt.ConvExt, opt.NicoSkipHb, opt.NicoAdjustVpos, opt.NicoConvForceConcat, opt.NicoConvSeqnoStart, opt.NicoConvSeqnoEnd)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if done {
				if nMp4s == 1 && (! skipped) {
					if 1 <= opt.NicoAutoDeleteDBMode {
						os.Remove(dbname)
					}
				} else if 1 < nMp4s || (nMp4s == 1 && skipped) {
					if 2 <= opt.NicoAutoDeleteDBMode {
						os.Remove(dbname)
					}
				}
				if (!opt.NicoNoStreamlink || !opt.NicoNoYtdlp) && !opt.NicoCommentOnly {
					os.Remove(dbname)
				}
			}
		}
	case "NICOLIVE_TEST":
		if err := niconico.TestRun(opt); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "ZIP2MP4":
		if err := zip2mp4.Convert(opt.ZipFile); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "DB2MP4":
		if strings.HasSuffix(opt.DBFile, ".yt.sqlite3") {
			zip2mp4.YtComment(opt.DBFile, opt.YtEmoji)

		} else if opt.ExtractChunks {
			if _, err := zip2mp4.ExtractChunks(opt.DBFile, opt.NicoSkipHb, opt.NicoAdjustVpos, opt.NicoConvSeqnoStart, opt.NicoConvSeqnoEnd); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		} else {
			if _, _, _, err := zip2mp4.ConvertDB(opt.DBFile, opt.ConvExt, opt.NicoSkipHb, opt.NicoAdjustVpos, opt.NicoConvForceConcat, opt.NicoConvSeqnoStart, opt.NicoConvSeqnoEnd); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

	case "DB2HLS":
		if opt.NicoHlsPort == 0 {
			fmt.Println("HLS port not specified")
			os.Exit(1)
		}
		if err := zip2mp4.ReplayDB(opt.DBFile, opt.NicoHlsPort, opt.NicoConvSeqnoStart); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "DBINFO":
		if strings.HasSuffix(opt.DBFile, ".yt.sqlite3") {
			if _, err := youtube.ShowDbInfo(opt.DBFile); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		} else {
			if _, err := niconico.ShowDbInfo(opt.DBFile, opt.ConvExt); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

	}


	return
}
