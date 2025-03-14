package niconico

import (
	"fmt"
	"time"
	"os"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"database/sql"

	"github.com/nnn-revo2012/livedl/files"
)

var SelMedia = `SELECT
	seqno, bandwidth, size, data FROM media
	WHERE IFNULL(notfound, 0) == 0 AND data IS NOT NULL
	ORDER BY seqno`

var SelComment = `SELECT
	vpos,
	date,
	date_usec,
	IFNULL(no, -1) AS no,
	IFNULL(anonymity, 0) AS anonymity,
	user_id,
	content,
	IFNULL(mail, "") AS mail,
	%s
	IFNULL(premium, 0) AS premium,
	IFNULL(score, 0) AS score,
	thread,
	IFNULL(origin, "") AS origin,
	IFNULL(locale, "") AS locale
	FROM comment
	ORDER BY date2`

func SelMediaF(seqnoStart, seqnoEnd int64) (ret string) {
	ret = `SELECT
	seqno, bandwidth, size, data FROM media
	WHERE IFNULL(notfound, 0) == 0 AND data IS NOT NULL`
	ret += ` AND seqno >= ` + fmt.Sprint(seqnoStart)
	ret += ` AND seqno <= ` + fmt.Sprint(seqnoEnd)
	ret += ` ORDER BY seqno`

	return
}

func (hls *NicoHls) dbOpen() (err error) {
	db, err := sql.Open("sqlite3", hls.dbName)
	if err != nil {
		return
	}

	hls.db = db

	_, err = hls.db.Exec(`
		PRAGMA synchronous = OFF;
		PRAGMA journal_mode = WAL;
	`)
	if err != nil {
		return
	}

	err = hls.dbCreate()
	if err != nil {
		hls.db.Close()
	}
	return
}

func (hls *NicoHls) dbCreate() (err error) {
	hls.dbMtx.Lock()
	defer hls.dbMtx.Unlock()

	// table media

	_, err = hls.db.Exec(`
	CREATE TABLE IF NOT EXISTS media (
		seqno     INTEGER PRIMARY KEY NOT NULL UNIQUE,
		current   INTEGER,
		position  REAL,
		notfound  INTEGER,
		bandwidth INTEGER,
		size      INTEGER,
		data      BLOB
	)
	`)
	if err != nil {
		return
	}

	_, err = hls.db.Exec(`
	CREATE UNIQUE INDEX IF NOT EXISTS media0 ON media(seqno);
	CREATE INDEX IF NOT EXISTS media1 ON media(position);
	---- for debug ----
	CREATE INDEX IF NOT EXISTS media100 ON media(size);
	CREATE INDEX IF NOT EXISTS media101 ON media(notfound);
	`)
	if err != nil {
		return
	}

	// table comment

	_, err = hls.db.Exec(`
	CREATE TABLE IF NOT EXISTS comment (
		vpos      INTEGER NOT NULL,
		date      INTEGER NOT NULL,
		date_usec INTEGER NOT NULL,
		date2     INTEGER NOT NULL,
		no        INTEGER,
		anonymity INTEGER,
		user_id   TEXT NOT NULL,
		content   TEXT NOT NULL,
		mail      TEXT,
		name      TEXT,
		premium   INTEGER,
		score     INTEGER,
		thread    TEXT,
		origin    TEXT,
		locale    TEXT,
		hash      TEXT UNIQUE NOT NULL
	)`)
	if err != nil {
		return
	}

	_, err = hls.db.Exec(`
	CREATE UNIQUE INDEX IF NOT EXISTS comment0 ON comment(hash);
	---- for debug ----
	CREATE INDEX IF NOT EXISTS comment100 ON comment(date2);
	CREATE INDEX IF NOT EXISTS comment101 ON comment(no);
	`)
	if err != nil {
		return
	}


	// kvs media

	_, err = hls.db.Exec(`
	CREATE TABLE IF NOT EXISTS kvs (
		k TEXT PRIMARY KEY NOT NULL UNIQUE,
		v BLOB
	)
	`)
	if err != nil {
		return
	}
	_, err = hls.db.Exec(`
	CREATE UNIQUE INDEX IF NOT EXISTS kvs0 ON kvs(k);
	`)
	if err != nil {
		return
	}

	// syncData

	_, err = hls.db.Exec(`
	CREATE TABLE IF NOT EXISTS sync (
		seqno     INTEGER PRIMARY KEY NOT NULL UNIQUE,
		date      INTEGER NOT NULL
	)
	`)
	if err != nil {
		return
	}
	_, err = hls.db.Exec(`
	CREATE UNIQUE INDEX IF NOT EXISTS sync0 ON sync(seqno);
	`)
	if err != nil {
		return
	}

	//hls.__dbBegin()

	return
}

// timeshift
func (hls *NicoHls) dbSetPosition() {
	hls.dbExec(`UPDATE media SET position = ? WHERE seqno=?`,
		hls.playlist.position,
		hls.playlist.seqNo,
	)
}

// timeshift
func (hls *NicoHls) dbGetLastPosition() (res float64) {
	hls.dbMtx.Lock()
	defer hls.dbMtx.Unlock()

	hls.db.QueryRow("SELECT position FROM media ORDER BY POSITION DESC LIMIT 1").Scan(&res)
	return
}

//func (hls *NicoHls) __dbBegin() {
//	return
	///////////////////////////////////////////
	//hls.db.Exec(`BEGIN TRANSACTION`)
//}
//func (hls *NicoHls) __dbCommit(t time.Time) {
//	return
	///////////////////////////////////////////

	//// Never hls.dbMtx.Lock()
	//var start int64
	//hls.db.Exec(`COMMIT; BEGIN TRANSACTION`)
	//if t.UnixNano() - hls.lastCommit.UnixNano() > 500000000 {
	//	log.Printf("Commit: %s\n", hls.dbName)
	//}
	//hls.lastCommit = t
//}
func (hls *NicoHls) dbCommit() {
//	hls.dbMtx.Lock()
//	defer hls.dbMtx.Unlock()

//	hls.__dbCommit(time.Now())
}
func (hls *NicoHls) dbExec(query string, args ...interface{}) {
	hls.dbMtx.Lock()
	defer hls.dbMtx.Unlock()

	if hls.nicoDebug {
		start := time.Now().UnixNano()
		defer func() {
			t := (time.Now().UnixNano() - start) / (1000 * 1000)
			if t > 100 {
				fmt.Fprintf(os.Stderr, "%s:[WARN]dbExec: %d(ms):%s\n", debug_Now(), t, query)
			}
		}()
	}

	if _, err := hls.db.Exec(query, args...); err != nil {
		fmt.Printf("dbExec %#v\n", err)
		//hls.db.Exec("COMMIT")
		hls.db.Close()
		os.Exit(1)
	}
}

func (hls *NicoHls) dbKVSet(k string, v interface{}) {
	query := `INSERT OR REPLACE INTO kvs (k,v) VALUES (?,?)`
	hls.startDBGoroutine(func(sig <-chan struct{}) int {
		hls.dbExec(query, k, v)
		return OK
	})
}

func (hls *NicoHls) dbKVExist(k string) (res int){
	hls.dbMtx.Lock()
	defer hls.dbMtx.Unlock()
	query := `SELECT COUNT(*) FROM kvs WHERE k = ?`
	hls.db.QueryRow(query, k).Scan(&res)
	return
}

func DbKVGet(db *sql.DB) (data map[string]interface{}) {
	data = make(map[string]interface{})
	rows, err := db.Query(`SELECT k,v FROM kvs`)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var k	string
		var v	interface{}
		err := rows.Scan(&k, &v)
		if err != nil {
			log.Println(err)
		}
		data[k] = v
	}

	return
}

func (hls *NicoHls) dbSyncSet(data string) {
	var seqno, date int64
	if ma := regexp.MustCompile(`"beginning_timestamp"\:(\d+)\,"sequence"\:(\d+)`).FindStringSubmatch(data); len(ma) > 0 {
		//fmt.Printf("syncData %s=%s\n", ma[2], ma[1])
		seqno, _ = strconv.ParseInt(ma[2], 10, 64)
		date, _ = strconv.ParseInt(ma[1], 10, 64)
	}
	query := `INSERT OR REPLACE INTO sync (seqno,date) VALUES (?,?)`
	hls.startDBGoroutine(func(sig <-chan struct{}) int {
		hls.dbExec(query, seqno, date)
		return OK
	})
}

func DbSyncGet(db *sql.DB) (data []string) {
	rows, err := db.Query(`SELECT seqno, date FROM sync ORDER by seqno`)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var seqno	int64
		var date	int64
		err := rows.Scan(&seqno, &date)
		if err != nil {
			log.Println(err)
		}
		//fmt.Printf("data: %d,%d\n", seqno, date)
		data = append(data, fmt.Sprintf("%d,%d",seqno, date))
	}

	return
}

func (hls *NicoHls) dbInsertReplaceOrIgnore(table string, data map[string]interface{}, replace bool) {
	var keys []string
	var qs []string
	var args []interface{}

	for k, v := range data {
		keys = append(keys, k)
		qs = append(qs, "?")
		args = append(args, v)
	}

	var replaceOrIgnore string
	if replace {
		replaceOrIgnore = "REPLACE"
	} else {
		replaceOrIgnore = "IGNORE"
	}

	query := fmt.Sprintf(
		`INSERT OR %s INTO %s (%s) VALUES (%s)`,
		replaceOrIgnore,
		table,
		strings.Join(keys, ","),
		strings.Join(qs, ","),
	)

	hls.startDBGoroutine(func(sig <-chan struct{}) int {
		hls.dbExec(query, args...)
		return OK
	})
}

func (hls *NicoHls) dbInsert(table string, data map[string]interface{}) {
	hls.dbInsertReplaceOrIgnore(table, data, false)
}
func (hls *NicoHls) dbReplace(table string, data map[string]interface{}) {
	hls.dbInsertReplaceOrIgnore(table, data, true)
}

// timeshift
func (hls *NicoHls) dbGetFromWhen() (when int64) {
	hls.dbMtx.Lock()
	defer hls.dbMtx.Unlock()
	var date2 int64

	hls.db.QueryRow("SELECT date2, FROM comment ORDER BY date2 ASC LIMIT 1").Scan(&date2)
	if date2 == 0 {
		var endTime float64
		hls.db.QueryRow(`SELECT v FROM kvs WHERE k = "endTime"`).Scan(&endTime)
		when = int64(endTime) + 360
		dt := time.Now()
		unix := dt.Unix()
		if unix < when {
			when = unix
		}
	} else {
		when = int64(date2) / (1000 * 1000) + 1
	}

	return
}

//新コメントサーバー用
func dbadjustVpos(opentime, vposbasetime, offset, date, vpos, premium int64) (ret int64) {
	ret = vpos
	if premium == 3 {
		ret = (date - opentime) * 100 - offset
	} else {
		ret = vpos + ((vposbasetime - opentime) * 100) - offset
	}
	return ret
}

func dbGetCommentRevision(db *sql.DB) (commentRevision int) {
	commentRevision = 0
	var nameCount int64
	db.QueryRow(`SELECT COUNT(name) FROM pragma_table_info('comment') WHERE name = 'name'`).Scan(&nameCount)
	if nameCount > 0 {
		commentRevision = 1
	}
	db.QueryRow(`SELECT COUNT(k) FROM 'kvs' WHERE k = 'vposBaseTime'`).Scan(&nameCount)
	if nameCount > 0 {
		commentRevision = 2
	}
	db.QueryRow(`SELECT COUNT(k) FROM 'kvs' WHERE k = 'streamType'`).Scan(&nameCount)
	if nameCount > 0 {
		commentRevision = 3
	}
	return
}

func WriteComment(db *sql.DB, fileName string, skipHb, adjustVpos bool, seqnoStart, seqnoEnd int64) {

	var fSelComment = func(revision int) string {
		var selAppend string
		if revision >= 1 {
			selAppend += `IFNULL(name, "") AS name,`
		}
		return fmt.Sprintf(SelComment, selAppend)
	}

	commentRevision :=  dbGetCommentRevision(db)
	fmt.Println("commentRevision: ", commentRevision)
	if commentRevision < 2 {
		fmt.Println("DBfile is old. Can't output comments this program.")
		return
	}

	//kvsテーブルから読み込み
	var openTime, vposBaseTime int64
	var providerType string
	var streamType string
	var offset int64
	kvs := DbKVGet(db)

	//syncテーブルから読み込み
	sync := DbSyncGet(db)
	var sync_seqno, sync_date int64
	if len(sync) > 0 {
		//fmt.Printf("sync: %v\n", sync)
		data := strings.Split(sync[0], ",")
		sync_seqno, _ = strconv.ParseInt(data[0], 10, 64)
		sync_date, _ = strconv.ParseInt(data[1], 10, 64)
	}
	fmt.Printf("sync: %d=%d\n", sync_seqno, sync_date)

	var t float64
	var sts string
	t = kvs["openTime"].(float64)
	openTime = int64(t)
	t = kvs["vposBaseTime"].(float64)
	vposBaseTime = int64(t)
	if commentRevision > 2 {
		streamType = kvs["streamType"].(string)
	}
	sts = kvs["status"].(string)
	if sts == "ENDED" {
		if streamType == "dlive" {
			offset = seqnoStart * 600 //timeshift
		} else {
			offset = seqnoStart * 500 //timeshift
		}
	} else {
		if streamType == "dlive" {
			offset = (sync_date/10) - (openTime*100) + (seqnoStart-sync_seqno)*300 //on_air
		} else {
			offset = (sync_date/10) - (openTime*100) + (seqnoStart-sync_seqno)*150 //on_air
		}
	}
	providerType = kvs["providerType"].(string)
	fmt.Println("status: ", sts)

	fmt.Println("adjustVpos: ", adjustVpos)
	fmt.Println("providerType: ", providerType)
	fmt.Println("openTime: ", openTime)
	fmt.Println("vposBaseTime: ", vposBaseTime)
	fmt.Println("offset: ", offset)

	rows, err := db.Query(fSelComment(commentRevision))
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	fileName = files.ChangeExtention(fileName, "xml")
	fileName, err = files.GetFileNameNext(fileName)
	fmt.Println("xml file: ", fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	fmt.Fprintf(f, "%s\r\n", `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintf(f, "%s\r\n", `<packet>`)

	for rows.Next() {
		var vpos      int64
		var date      int64
		var date_usec int64
		var no        int64
		var anonymity int64
		var user_id   string
		var content   string
		var mail      string
		var name      string
		var premium   int64
		var score     int64
		var thread    string
		var origin    string
		var locale    string
		var dest0 = []interface{} {
			&vpos      ,
			&date      ,
			&date_usec ,
			&no        ,
			&anonymity ,
			&user_id   ,
			&content   ,
			&mail      ,
		}
		var dest1 = []interface{} {
			&premium   ,
			&score     ,
			&thread    ,
			&origin    ,
			&locale    ,
		}
		if commentRevision >= 1 {
			dest0 = append(dest0, &name)
		}
		var dest = append(dest0, dest1...)
		err = rows.Scan(dest...)
		if err != nil {
			log.Println(err)
			return
		}

		// skip /hb
		if (premium > 1) && skipHb && strings.HasPrefix(content, "/hb ") {
			continue
		}

		// vpos計算
		if adjustVpos == true {
			vpos = dbadjustVpos(openTime, vposBaseTime, offset, date, vpos, premium)
			// vposが-1000(-10秒)より小さい場合は出力しない
			if vpos <= -1000 {
				continue
			}
		} else {
			//premium=3のみvpos計算
			if premium == 3 {
				vpos = dbadjustVpos(openTime, vposBaseTime, 0, date, vpos, premium)
			}
		}

		line := fmt.Sprintf(
			`<chat thread="%s" vpos="%d" date="%d" date_usec="%d" user_id="%s"`,
			thread,
			vpos,
			date,
			date_usec,
			user_id,
		)

		if no >= 0 {
			line += fmt.Sprintf(` no="%d"`, no)
		}
		if anonymity != 0 {
			line += fmt.Sprintf(` anonymity="%d"`, anonymity)
		}
		if mail != "" {
			mail = strings.Replace(mail, `"`, "&quot;", -1)
			mail = strings.Replace(mail, "&", "&amp;", -1)
			mail = strings.Replace(mail, "<", "&lt;", -1)
			mail = strings.Replace(mail, ">", "&gt;", -1)
			line += fmt.Sprintf(` mail="%s"`, mail)
		}
		if name != "" {
			name = strings.Replace(name, `"`, "&quot;", -1)
			name = strings.Replace(name, "&", "&amp;", -1)
			name = strings.Replace(name, "<", "&lt;", -1)
			name = strings.Replace(name, ">", "&gt;", -1)
			line += fmt.Sprintf(` name="%s"`, name)
		}
		if origin != "" {
			origin = strings.Replace(origin, `"`, "&quot;", -1)
			origin = strings.Replace(origin, "&", "&amp;", -1)
			origin = strings.Replace(origin, "<", "&lt;", -1)
			origin = strings.Replace(origin, ">", "&gt;", -1)
			line += fmt.Sprintf(` origin="%s"`, origin)
		}
		if premium != 0 {
			line += fmt.Sprintf(` premium="%d"`, premium)
		}
		if score != 0 {
			line += fmt.Sprintf(` score="%d"`, score)
		}
		if locale != "" {
			locale = strings.Replace(locale, `"`, "&quot;", -1)
			locale = strings.Replace(locale, "&", "&amp;", -1)
			locale = strings.Replace(locale, "<", "&lt;", -1)
			locale = strings.Replace(locale, ">", "&gt;", -1)
			line += fmt.Sprintf(` locale="%s"`, locale)
		}
		line += ">"
		if premium == 3 {
			content = strings.Replace(content, "&", "&amp;", -1)
			content = strings.Replace(content, "<", "&lt;", -1)
		} else {
			content = strings.Replace(content, `"`, "&quot;", -1)
			content = strings.Replace(content, "&", "&amp;", -1)
			content = strings.Replace(content, "<", "&lt;", -1)
			content = strings.Replace(content, ">", "&gt;", -1)
		}
		line += content
		line += "</chat>"
		fmt.Fprintf(f, "%s\r\n", line)
	}
	fmt.Fprintf(f, "%s\r\n", `</packet>`)
}

func ShowDbInfo(fileName, ext string) (done bool, err error) {
	_, err = os.Stat(fileName)
	if err != nil {
		fmt.Println("sqlite3 file not found:")
		return
	}
	db, err := sql.Open("sqlite3", "file:"+url.PathEscape(fileName)+"?mode=ro&immutable=1")
	if err != nil {
		return
	}
	defer db.Close()

	fmt.Println("----- DATABASE info. -----")
	fmt.Println("sqlite3 file :", fileName)
	for _, tbl := range []string{"kvs", "media", "comment"} {
		if !dbIsExistTable(db, tbl) {
			fmt.Println("table", tbl, "not found")
		} else {
			fmt.Println("table", tbl, "exist")
		}
	}

	fmt.Println("----- broadcast info. -----")
	kvs := DbKVGet(db)
	if len(kvs) > 0 {
		id := kvs["nicoliveProgramId"].(string)
		title :=  kvs["title"].(string)
		sts :=  kvs["status"].(string)
		ptype :=  kvs["providerType"].(string)
		open :=  int64(kvs["openTime"].(float64))
		begin :=  int64(kvs["beginTime"].(float64))
		end :=  int64(kvs["endTime"].(float64))
		vpos :=  int64(kvs["vposBaseTime"].(float64))
		username :=  kvs["userName"].(string)

		fmt.Println("id: ", id)
		fmt.Println("title: ", title)
		fmt.Println("username: ", username)
		fmt.Println("providerType: ", ptype)
		fmt.Println("status: ", sts)
		fmt.Println("openTime: ", time.Unix(open, 0))
		if ptype == "official" {
			fmt.Println("beginTime: ", time.Unix(begin, 0))
		}
		fmt.Println("endTime: ", time.Unix(end, 0))
		fmt.Println("vposBaseTime: ", time.Unix(vpos, 0))
	} else {
		fmt.Println("kvs data not found")
	}
	//syncテーブルから読み込み
	sync := DbSyncGet(db)
	if len(sync) > 0 {
		fmt.Printf("sync: %v\n", sync)
	}

	commentRevision :=  dbGetCommentRevision(db)
	fmt.Println("commentRevision: ", commentRevision)

	media_all  := DbGetCountMedia(db , 0)
	media_err  := DbGetCountMedia(db , 2)
	media_sseq := DbGetFirstSeqNo(db , 0)
	media_eseq := DbGetLastSeqNo(db , 0)
	comm_data := DbGetCountComment(db)

	fmt.Println("----- media info. -----")
	fmt.Println("start seqno: ", media_sseq)
	fmt.Println("end seqno: ", media_eseq)
	fmt.Println("data: ", media_all, "(media:", media_all - media_err, "err:", media_err, ")")

	fmt.Println("----- comment info. -----")
	fmt.Println("data: ", comm_data)

	done = true

	return
}

// ts
func (hls *NicoHls) dbGetLastMedia(i int) (res []byte) {
	hls.dbMtx.Lock()
	defer hls.dbMtx.Unlock()
	hls.db.QueryRow("SELECT data FROM media WHERE seqno = ?", i).Scan(&res)
	return
}
//
func (hls *NicoHls) dbGetLastSeqNo(flg int) (res int64) {
	hls.dbMtx.Lock()
	defer hls.dbMtx.Unlock()
	var sql string
	if flg == 1 {
		sql = "SELECT seqno FROM media WHERE IFNULL(notfound, 0) == 0 AND data IS NOT NULL ORDER BY seqno DESC LIMIT 1"
	} else {
		sql = "SELECT seqno FROM media ORDER BY seqno DESC LIMIT 1"
	}
	hls.db.QueryRow(sql).Scan(&res)
	return
}
func DbGetLastSeqNo(db *sql.DB, flg int) (res int64) {
	var sql string
	if flg == 1 {
		sql = "SELECT seqno FROM media WHERE IFNULL(notfound, 0) == 0 AND data IS NOT NULL ORDER BY seqno DESC LIMIT 1"
	} else {
		sql = "SELECT seqno FROM media ORDER BY seqno DESC LIMIT 1"
	}
	db.QueryRow(sql).Scan(&res)
	return
}
func DbGetFirstSeqNo(db *sql.DB, flg int) (res int64) {
	var sql string
	if flg == 1 {
		sql = "SELECT seqno FROM media WHERE IFNULL(notfound, 0) == 0 AND data IS NOT NULL ORDER BY seqno ASC LIMIT 1"
	} else {
		sql = "SELECT seqno FROM media ORDER BY seqno ASC LIMIT 1"
	}
	db.QueryRow(sql).Scan(&res)
	return
}
func DbKVGetSeqNo(db *sql.DB, k string) (res int64) {
	query := `SELECT v FROM kvs WHERE k = ?`
	db.QueryRow(query, k).Scan(&res)
	return
}
func DbGetCountMedia(db *sql.DB, flg int) (res int64) {
	var sql string
	if flg == 1 {
		sql = "SELECT COUNT(seqno) FROM media WHERE IFNULL(notfound, 0) == 0 AND data IS NOT NULL"
	} else if flg == 2 {
		sql = "SELECT COUNT(seqno) FROM media WHERE IFNULL(notfound, 0) != 0 OR data IS NULL"
	} else {
		sql = "SELECT COUNT(seqno) FROM media"
	}
	db.QueryRow(sql).Scan(&res)
	return
}
func DbGetCountComment(db *sql.DB) (res int64) {
	db.QueryRow("SELECT COUNT(date) FROM comment").Scan(&res)
	return
}
func dbIsExistTable(db *sql.DB, table_name string) (ret bool) {
	var res int
	ret = false
	if len(table_name) > 0 {
		db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE TYPE='table' AND name=?", table_name).Scan(&res)
		if res > 0 {
			ret = true
		}
	}
	return
}
