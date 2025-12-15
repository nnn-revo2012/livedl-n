# livedl  
~~livedl-nは開発終了いたしました~~  
1週間ぐらい限定で再公開 **(ソースファイルのみ、最新の実行ファイルはありません)**  

## 動画のダウンロードについて  
**ニコニコ生放送（ニコ生）は2025/02/05から新動画サーバー(dlive)の配信に移行しています。  
dliveで配信する動画はすべてAES128暗号化されており、これを解除する方法の公開やツール作成は日本の著作権法に違反する可能性があります。  
現在の作者(nnn-revo2012)は日本在住なのでこのツールを対応させることができません。  
※アメリカ、EU、中国、韓国を含むほとんどの国でもDRM暗号化の解除は違法なのでAES暗号化動画の解除も違法になる可能性があります。  
今後の動画についてはご自分で情報を探すなりして対応してください。**  

## 使い方
**以下のwikiを参照してください**  
https://github.com/nnn-revo2012/livedl-n/wiki  

## This version is require golang 1.20.x or higher and module-aware mode.   
See [Migrating to Go Modules](https://blog.golang.org/migrating-to-go-modules)  
### **golang 1.19.x or lower**  
If you use **golang 1.19.x or lower**, use golang 1.20.x or higher.  

## Get the latest livedl source program  

` git clone https://github.com/nnn-revo2012/livedl-n.git`  
(Pull requests other than #51 on himananiito/livedl and 5ch's patches have been applied)  

**This repository has been rebuilt. If you cloned it before December 13, 2020, please do the following:**  
```
git fetch origin master
git reset --hard origin/master
```

## How to build a program  

See https://github.com/himananiito/livedl  

go build:  
- Linux and Mac:
  ```
  $ go build -C src -o ../livedl livedl.go
  ```  
- Windows:
  ```
  > go build -C src -o ../livedl.exe livedl.go
  ```  
- Windows (and use powershell):
  ```
  > ./build.ps1
  ```  
- Docker:

  - Get a livedl source
    ```
    git clone https://github.com/nnn-revo2012/livedl-n.git
    cd livedl
    git checkout master # Or another version that supports docker (contains Dockerfile)
    ```

  - Make Imagefile
    ```
    docker build -t livedl .
    ```

  - excution docker image
    - mount an output directry on /livedl
    ```
    docker run --rm -it -v "$(pwd):/livedl" livedl "https://live.nicovideo.jp/watch/..."
    ```

## 参考にしたサイト  
### 新コメントサーバー関連  

- nicolive-comment-protobuf  
  https://github.com/n-air-app/nicolive-comment-protobuf  

- 帰ってきたニコニコのニコ生コメントサーバーからのコメント取得備忘録  
  https://qiita.com/DaisukeDaisuke/items/3938f245caec1e99d51e  

- NDGRClient  
  https://github.com/tsukumijima/NDGRClient  

- protobuf  
  https://github.com/protocolbuffers/protobuf/releases  

- protobuf-go  
  https://github.com/protocolbuffers/protobuf-go  

