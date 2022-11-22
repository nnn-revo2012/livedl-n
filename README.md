# livedl  

Nicolive recording tool that supports new distribution (HTML5).  

## This version is in module-aware mode.   
See [Migrating to Go Modules](https://blog.golang.org/migrating-to-go-modules)  
### **golang 1.15.x or lower**  
If you use **golang 1.15.x or lower**, set the environment variable **`GO111MODULE=on`** before `go build`.  
Otherwise you will get an error at compile.  
Or use golang 1.16.x or higher.  

```  
cd src
GO111MODULE=on go build -o ../livedl livedl.go
cd ..
```  
Or  
```
export GO111MODULE=on
cd src
go build -o ../livedl livedl.go
cd ..
``` 

## Get the latest livedl source program  

` git clone https://github.com/nnn-revo2012/livedl.git`  
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
  $ cd src; go build -o ../livedl livedl.go; cd ..
  ```  
- Windows:
  ```
  > cd src
  > go build -o ../livedl.exe livedl.go
  > cd ..
  ```  
- Windows (and use powershell):
  ```
  > ./build.ps1
  ```  
- Docker:

  - Get a livedl source
    ```
    git clone https://github.com/nnn_revo2012/livedl.git
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
