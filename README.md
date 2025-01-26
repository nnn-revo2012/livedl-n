# livedl  

Nicolive recording tool that supports new distribution (HTML5).  

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

