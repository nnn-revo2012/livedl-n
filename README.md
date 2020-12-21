# livedl  

Niconico live recording tool that supports new distribution (HTML5).  

## How to build a program  

See https://github.com/himananiito/livedl or other forked repositories (e.g. https://github.com/hanaonnao/livedl)  

## Get the latest livedl source program  (**RECOMMEND**)  

Get one of the following with `git clone`  
- ` git clone https://github.com/nnn-revo2012/livedl.git ` (Pull requests other than #51 on himananiito/livedl and 5ch's patches have been applied)  

**This repository has been rebuilt. If you cloned it before December 13, 2020, please do the following:**  
```
git fetch origin master
git reset --hard origin/master
```

- ` git clone https://github.com/hanaonnao/livedl.git ` (5ch's patches have been applied)

## Update the program from scratch  

1. Merge pull requests on himananiito/livedl 
2. Apply the 5ch forum patches by volunteers  

Follow the steps below to merge pull requests and apply patches on the 5ch forum by volunteers.  

## Merge the pull requests on himananiito/livedl  

Merge the pull requests on himananiito/livedl, merge the following.  

- **If you merge PR #51, conflicts will occur frequently when applying 5ch's patches**  
- **If you merge PR #43 and #44 , conflicts will occur when applying 5ch's patches**  
```
git clone https://github.com/himananiito/livedl.git
cd livedl
git checkout master
git fetch origin pull/39/head:dev
git checkout dev
git reset --hard HEAD^
git checkout master
git merge dev
git fetch origin pull/47/head:seqskipfix
git merge seqskipfix
```

## Apply the 5ch forum patches by volunteers

### Precautions when applying patches

- Unix(Linux) and Mac  
When applying patch, if line enfings of the source program and the patch are different, the following error may occur.  
```
Hunk #1 at 49(different line endings).
```

Because the source program using line endings CRLF.
In this case, change the line feed code of the source program from CRLF to LF as follows.
```
dos2unix (source filename)  
```
e.g.  
```
dos2unix src/niconico/nico_hls.go  
dos2unix src/niconico/nico_db.go  
```

- Windows  

If the source program using line endings CRLF, change the line feed code of the patch to CRLF.  
```
unix2dos (patch filename)  
```
Also, add `` --binary`` when applying the patch.  
```
patch -p1 --binary <(patch filename)  
```
### Patch lists

- Fixed niconico has upgraded their websocket API(v2) on June 2, 2020  
http://egg.5ch.net/test/read.cgi/software/1570634489/535  
- Fixed 'broadcastId not found' error on June 25, 2020  
https://egg.5ch.net/test/read.cgi/software/1570634489/744  
- Fixed comments cannot be saved on July 27, 2020  
https://egg.5ch.net/test/read.cgi/software/1570634489/932  

**These patches have been newly uploaded in one**  
https://egg.5ch.net/test/read.cgi/software/1595715643/116  
```
[Copy the patch to the livedl folder]
patch -p1 <livedl-p20200919.01.patch
```
The following error is displayed when applying the patch.  
```
patching file src/niconico/nico_hls.go  
Hunk #8 FAILED at 638.  
1 out of 20 hunks FAILED -- saving rejects to file src/niconico/nico_hls.go.rej  
```
Modify line 638 (or surroundings) of `src/niconico/nico_hls.go` with an editor as follows.  
(Change **http** in the original line to **https**.)  
```
uri := fmt.Sprintf("https://live.nicovideo.jp/api/getwaybackkey?thread=%s",url.QueryEscape(threadId))
```

- Add self-chasing playback function and others  
https://egg.5ch.net/test/read.cgi/software/1595715643/57  
```
[Copy the patch to the livedl folder]
patch -p1 <livedl-55.patch
```
- Add stop time-shift recording at specified time  
https://egg.5ch.net/test/read.cgi/software/1595715643/163  
```
[Copy the patch to the livedl folder]
patch -p1 <livedl.ts-start-stop.patch
```

- Fixed to save the name attribute of XML comments (used when performers make named comments)  
https://egg.5ch.net/test/read.cgi/software/1595715643/174  

Patch applies only livedl.comment-name-attribute-r1.patch.gz  
https://egg.5ch.net/test/read.cgi/software/1595715643/194  
```
[Copy the patch to the livedl folder]
patch -p0 <livedl.comment-name-attribute-r1.patch
```

- Fixed to record old time-shift (Flash Player)  
https://egg.5ch.net/test/read.cgi/software/1595715643/228  

```
[Copy the patch to the livedl folder]
patch -p1 <livedl.getedgestatus.patch
```

- Add -http-timeout option: http connect and read timeout (Default is: 5)  
https://egg.5ch.net/test/read.cgi/software/1595715643/272  

```
[Copy the patch to the livedl folder]
patch -p0 <livedl.http-timeout.patch
```
The following error is displayed when applying the patch.  
```
$ patch -p0 <livedl.http-timeout.patch
patching file src/options/options.go
Hunk #2 FAILED at 63.
```
Modify line 69 of `src/niconico/options.go` with an editor as follows.  
(Add **HttpTimeout** in the original line.)  
```
	HttpProxy              string
	NoChdir                bool
	HttpTimeout            int
```

- Fixed youtube live 'ytplayer parse error'.  
https://egg.5ch.net/test/read.cgi/software/1595715643/402  
https://egg.5ch.net/test/read.cgi/software/1595715643/406  
Patch applies only livedl.youtube-r1.patch  
```
[Copy the patch to the livedl folder]
patch -p0 <livedl.youtube-r1.patch
```

- Fix not to call getwaybackkey API  
https://egg.5ch.net/test/read.cgi/software/1595715643/424  
```
[Copy the patch to the livedl folder]
patch -p1 <livedl.waybackkey.patch
```
