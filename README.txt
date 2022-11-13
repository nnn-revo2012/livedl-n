livedl (20221108.51-windows-amd64)
Usage:
livedl [COMMAND] options... [--] FILE

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
  -nico-hls-only                 録画時にHLSのみを試す
  -nico-hls-only=on              (+) 上記を有効に設定
  -nico-hls-only=off             (+) 上記を無効に設定(デフォルト)
  -nico-rtmp-only                録画時にRTMPのみを試す
  -nico-rtmp-only=on             (+) 上記を有効に設定
  -nico-rtmp-only=off            (+) 上記を無効に設定(デフォルト)
  -nico-rtmp-max-conn <num>      RTMPの同時接続数を設定
  -nico-rtmp-index <num>[,<num>] RTMP録画を行うメディアファイルの番号を指定
  -nico-hls-port <portnum>       [実験的] ローカルなHLSサーバのポート番号
  -nico-limit-bw <bandwidth>     (+) HLSのBANDWIDTHの上限値を指定する。0=制限なし
                                 audio_high or audio_only = 音声のみ
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
  -nico-adjust-vpos=on           (+) コメント書き出し時にvposの値を補正する
                                 vposの値が-1000より小さい場合はコメント出力しない
  -nico-adjust-vpos=off          (+) コメント書き出し時にvposの値をそのまま出力する(デフォルト)
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
  -yt-no-youtube-dl=on           (+) youtube-dlを使用しない
  -yt-no-youtube-dl=off          (+) youtube-dlを使用する(デフォルト)
  -yt-comment-start              YouTube Liveアーカイブでコメント取得開始時間（秒）を指定
                                 ＜分＞:＜秒＞ | ＜時＞:＜分＞:＜秒＞ の形式でも指定可能
                                 0：続きからコメント取得  1：最初からコメント取得
  -yt-emoji=on                   (+) コメントにemojiを表示する(デフォルト)
  -yt-emoji=off                  (+) コメントにemojiを表示しない

変換オプション:
  -extract-chunks=off            (+) -d2mで動画ファイルに書き出す(デフォルト)
  -extract-chunks=on             (+) [上級者向] 各々のフラグメントを書き出す(大量のファイルが生成される)
  -conv-ext=mp4                  (+) -d2mで出力の拡張子を.mp4とする(デフォルト)
  -conv-ext=ts                   (+) -d2mで出力の拡張子を.tsとする

HTTP関連
  -http-skip-verify=on           (+) TLS証明書の認証をスキップする (32bit版対策)
  -http-skip-verify=off          (+) TLS証明書の認証をスキップしない (デフォルト)
  -http-timeout <num>            (+) タイムアウト時間（秒）デフォルト: 5秒（最低値）


(+)のついたオプションは、次回も同じ設定が使用されることを示す。

FILE:
  ニコニコ生放送/nicolive:
    https://live.nicovideo.jp/watch/lvXXXXXXXXX
    lvXXXXXXXXX
  ツイキャス/twitcasting:
    https://twitcasting.tv/XXXXX


﻿更新履歴
202xxxxx.52
・録画済みのデータベース(sqlite3)の各種情報を表示するコマンド(-dbinfo)追加
    ./livedl -dbinfo -- 'データーベースのファイル名をフルパスで'
  - youtubeのデーターベースはcomment情報のみ表示
  - データベース情報表示、データベースextractの際DBをreadonlyで開くように修正
  - データベースファイルの存在チェックを追加

20221108.51
・直接ログインの２段階認証(MFA)対応
・上記に伴うlogin APIのendpoint、cookie取得方法の変更
・firefoxからのcookie取得機能追加
  -nico-cookies firefox[:profile|cookiefile]
  e.g.  
  - profile default-release のcookieを取得
      ./livedl -nico-cookies firefox
  - profile NicoTaro のcookieを取得
      ./livedl -nico-cookies firefox:NicoTaro 
  - 直接cookiefileを指定
      ./livedl -nico-cookies firefox:'C:/Users/*******/AppData/Roaming/Mozilla/Firefox/Profiles/*****/cookies.sqlite' 
※Mac/Linuxで `cookies from browser failed: firefox profiles not found`が 表示される場合は報告おねがいします   
※直接cookiefile指定の場合は必ず'か"で囲ってください  
※プロファイルにspaceを含む場合は'か"で囲ってください  

20220905.50
・ニコ生のコメントのvposを補正
  -nico-adjust-vpos=on
     コメント書き出し時にvposの値を補正する
     vposの値が-1000より小さい場合はコメント出力しない
  -nico-adjust-vpos=off
     コメント書き出し時にvposの値をそのまま出力する(デフォルト)
  ※ExtractChunks()もコメントvposを補正するように修正
  ※ニコ生の生放送を録画する際、再接続してもkvsテーブルは更新されません
・Youtubeのコメントにemojiを出力する/しない
  -yt-emoji=on
     コメントにemojiを表示する(デフォルト)
  -yt-emoji=off
  コメントにemojiを表示しない
・音声のみ録画対応
　-nico-limit-bw に audio_high または audio_only を指定してください
・-http-timeout の設定を保存するように修正
・live2.* -> live.* に修正
・その他主にコメント関連の修正
  - livedl.exeとsqlite3ファイルが別のフォルダーにある場合、コメント出力時にxmlファイルに
    -数字が付かなかったのを修正
  - ニコ生の生放送の最初に取得するコメント数を1000から100に変更した(サーバー側の仕様による)

20211017.49
・livedlのあるディレクトリ以外から実行する時カレントディレクトリにconf.dbが作成されるのを修正
https://egg.5ch.net/test/read.cgi/software/1595715643/922
例: C:\bin\livedl\livedl.exe を D:\home\tmp をカレントディレクトリとして実行した場合、conf.dbは D:\home\tmp に作成されてしまう

仕様：conf.dbは実行するlivedlと同じディレクトリに作成する
      ただし、オプション -no-chdir が指定された場合はカレントディレクトリにconf.dbを作成する
      (livedl実行ファイルがユーザ書き込み権限のないディレクトリにある場合を想定)

20210607.48
・livedl で YouTube Live のアーカイブコメントの取得開始時刻を指定するオプション
https://egg.5ch.net/test/read.cgi/software/1595715643/789

使用例：
　livedl -yt-comment-start 3:21:06 https://～
  特殊例 0：続きからコメント取得  1：最初からコメント取得

20210202.47
・livedl で-yt-no-streamlink=on -yt-no-youtube-dl=on が指定されたとき、YouTube Live のコメントを永久に取得し続けるパッチ
https://egg.5ch.net/test/read.cgi/software/1595715643/567

・livedl を YouTube Live の直近の仕様変更に対応
https://egg.5ch.net/test/read.cgi/software/1595715643/559

20210128.46
・金額のフォーマットの要望ないみたいだからこっちで勝手に決めさせてもらったよ
https://egg.5ch.net/test/read.cgi/software/1595715643/543

Youtube Liveのコメントにamount属性を追加

・livedl で YouTubeLive リプレイのコメントが取れるよう直したよ
https://egg.5ch.net/test/read.cgi/software/1595715643/523

Youtube Liveのコメントが途中で切れるのを修正
Youtube liveのアーカイブダウンロード中に`json decode error'となり中断するのを修正
(404エラーになる場合は少しwaitする)

20210102.45
・livedl で一部コメントが保存されないのを修正するパッチ
https://egg.5ch.net/test/read.cgi/software/1595715643/457

・Remove VPOS > 0 (Commit 03417972d920cce0af92221583fc42bc559ef469)
　（VPOS<0のコメントが抜けるため削除した）

20201221.44
・livedl で waybackkey の取得方法を変更するパッチ
https://egg.5ch.net/test/read.cgi/software/1595715643/424

20201213.43
・livedl で YouTube Live を扱えるようにするためのパッチ（リビジョン1）
patch は livedl.youtube-r1.patch のみ適用
https://egg.5ch.net/test/read.cgi/software/1595715643/402
https://egg.5ch.net/test/read.cgi/software/1595715643/406

20201115.42
・livedl で HTTP のタイムアウト時間を変更できるようにするパッチ
https://egg.5ch.net/test/read.cgi/software/1595715643/272

20201026.41
・旧配信のタイムシフトを録画できるようにするパッチ
https://egg.5ch.net/test/read.cgi/software/1595715643/228

20201008.40
・XMLコメントのname属性(出演者が名前付きのコメントする時に使用)を保存するように修正
https://egg.5ch.net/test/read.cgi/software/1595715643/174
patch は livedl.comment-name-attribute-r1.patch.gz のみ適用
https://egg.5ch.net/test/read.cgi/software/1595715643/194

20201008
・指定時間でタイムシフト録画を停止するためのパッチ（＋α）
https://egg.5ch.net/test/read.cgi/software/1595715643/163

オプション
　-nico-ts-start ＜num＞
　　タイムシフトの録画を指定した再生時間（秒）から開始する
　-nico-ts-stop ＜num＞
　　タイムシフトの録画を指定した再生時間（秒）で停止する
　上記2つは ＜分＞:＜秒＞ | ＜時＞:＜分＞:＜秒＞ の形式でも指定可能

　-nico-ts-start-min ＜num＞
　　タイムシフトの録画を指定した再生時間（分）から開始する
　-nico-ts-stop-min ＜num＞
　　タイムシフトの録画を指定した再生時間（分）で停止する
　上記2つは ＜時＞:＜分＞ の形式でも指定可能

20200903.39
https://egg.5ch.net/test/read.cgi/software/1595715643/57
・セルフ追っかけ再生
　例：http://127.0.0.1:12345/m3u8/2/1200/index.m3u8
　　現在のシーケンス番号から1200セグメント（リアルタイムの場合30分）戻ったところを再生

・追加オプション
　-nico-conv-seqno-start ＜num＞
　　MP4への変換を指定したセグメント番号から開始する
　-nico-conv-seqno-end ＜num＞
　　MP4への変換を指定したセグメント番号で終了する
　-nico-conv-force-concat
　　MP4への変換で画質変更または抜けがあっても分割しないように設定
　-nico-conv-force-concat=on
　　(+) 上記を有効に設定
　-nico-conv-force-concat=off
　　(+) 上記を無効に設定(デフォルト)

・　-d2h
　　[実験的] 録画済みのdb(.sqlite3)を視聴するためのHLSサーバを立てる(-db-to-hls)
　　　開始シーケンス番号は（変換ではないが） -nico-conv-seqno-start で指定
　　　　使用例：$ livedl lvXXXXXXXXX.sqlite3 -d2h -nico-hls-port 12345 -nico-conv-seqno-start 2780

20200828.38
ニコ生仕様変更に対応
https://egg.5ch.net/test/read.cgi/software/1595715643/116

・コメントサーバー仕様変更に対応(threadId、waybackkey廃止など)(2020/07/27)
　ID:jM/9Q+5+0作成のpatchを適用
https://egg.5ch.net/test/read.cgi/software/1570634489/932

・'broadcastId not found'エラー表示されるのを修正(2020/06/25)
　ID:jM/9Q+5+0作成のpatchを適用
https://egg.5ch.net/test/read.cgi/software/1570634489/744

・放送情報取得時のwsapiがv2に変わって録画できなくなったのを修正(2020/06/02)
　ID:jM/9Q+5+0作成のpatchを適用
http://egg.5ch.net/test/read.cgi/software/1570634489/535

・20181215.35以降の修正を追加
・TS録画時にセグメント抜けが起こるのを修正 (PR#47)
・http -> httpsに修正 (PR#39)

20181215.35
・-nico-ts-start-minオプションの追加
・win32bit版のビルドを追加
・-http-skip-verifyオプションを保存できるようにした
・ライセンスをMITにした

20181107.34
・[ニコ生] (暫定)TEMPORARILY_CROWDEDで録画終了するようにした
・ファイル名が半角ドットで終わる場合に全角ドットにした
・[YouTubeLive] コメントの改行をCRLFにした
・[ニコ生TS] タイムシフトの録画を指定した再生時間(秒)から開始するオプション追加(merged)
・[ニコ生TS] 32bitで終了しない問題を修正(merged)

20181008.33
・[Youtube] チャットが取得できない問題を修正
・[Youtube] Streamlinkでダウンロードできない場合にyoutube-dlを使うようにした
・[Youtube] コメントファイルを書き出せるようにした。
・#15 [ニコ生コメント] 出力をCRLFにした。/hbコマンドを出さないオプションを追加

20181003.32
・#14 ★緊急 [ニコ生] 新配信録画のプレイリスト取得にウェイトが入らない問題を修正
・#9 [ニコ生TS] プレイリストの最後で無限ループしてしまう問題を修正
・YoutubeLiveコメント対応中(未完了)
・[実験的] -yt-api-key オプションの追加（未使用）

20180925.31
・#8 [ツイキャス] 「c:」から始まるユーザ名が録画できない問題を修正
・#11 [ツイキャス] 実行直後またはリトライ中にエラーで終了する問題を修正
・#10 [ツイキャス] -tcas-retry-intervalが効かない問題を修正
・#12 [ニコ生] タイムシフトで先頭のセグメント(seqno=0)が取得できない問題を修正

