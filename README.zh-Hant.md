# webpc
通過 http5 管理多台遠程計算機

[English](README.md) 中文README_zh.md

TeamViewer 等收費軟件幫你遠程控制多臺計算機，WebPC 完成類似功能並且完全免費且開源。同時 WebPC 充分利用現代瀏覽器只需要藉助瀏覽器即可遠程掌控你的計算機無需要安裝額外控制端。

WebPC 主要有下述特色：

* 所有操作通過瀏覽器無需要客戶端軟件
* 所有受控端都在瀏覽器上提供了模仿本地系統的檔案操作
* 通過遠程 shell 在受控端執行指令，例如運行 vim
* 在瀏覽器上通常 noVNC 實現到受控端的遠程桌面功能
* 受控端可自行配置提供的功能，避免 WebPC 被用作木馬程序
* 簡易的用戶和分組系統，以控制用戶訪問權限和範圍
* 支持多種系統平臺(linux windows mac ...)

## 內容列表

- [背景](#背景)
- [如何工作](#如何工作)
- [安裝](#安裝)
    - [linux-amd64](#linux-amd64)
    - [windows-amd64](#windows-amd64)
- [設定](#設定)
    - [設定master](#設定master)
    - [設定slave](#設定slave)
- [編譯](#編譯)
    - [編譯前端](#編譯前端)
    - [編譯WebPC](#編譯WebPC)

## 背景

計算機設備的普及以及網路的快速發展使遠程操作電腦變得越來越必要跟普遍，ssh rdp 都只能控制一臺有獨立外網ip的設備，並且它們各有優缺。提供遠程管控多臺設備的服務如 TeamViewer 通常都是收費不適合非商業用戶。並且通常只提供了遠程桌面而對於linux等設備遠程 shell 通常更有用且有效率，此外還需要安裝煩人的遠控操作端程式也令人厭惡。於是我整理了上述需要，並且將所有客戶功能實現到瀏覽器便是此項目 WebPC。

## 如何工作

1. 首先需要一臺服務器接收用戶請求以及受控端的註冊。我們稱這個服務器爲 master，所有受控端爲 slave。

2. slave 會向 master 註冊自己，並且維持一個到master的虛擬網路通道，salve 在此虛擬網路通道上爲 master 提供 grpc 服務以支持各種遠程控制功能。

3. 當 master 收到用戶請求時，尋找註冊的 slave，並將請求轉發到 slave。

## 安裝

對於 linux-amd64 和 windows-amd64 已經提供好了預編譯程序，請直接下載安裝，對於其它平臺需要自行編譯，並且參考下文進行安裝。

### linux-amd64

#### master

1. 下載最新程式並且解壓到 /opt/webpc/webpc

2. 複製 /opt/webpc/webpc-master.service 到 /etc/systemd/system/webpc-master.service 爲 systemd 安裝 master 服務

    ```
    sudo cp /opt/webpc/webpc-master.service /etc/systemd/system/webpc-master.service
    ```

    如果修改了安裝目錄，記得要修改 webpc-master.service 中對應的路徑

3. 創建 webpc 用戶

    ```
    sudo useradd webpc -Mrs /sbin/nologin
    ```

4. 運行服務 

    ```
    sudo systemctl start webpc-master.service
    ```

#### slave

1. 下載最新程式並且解壓到 /opt/webpc/webpc

2. 複製 /opt/webpc/webpc-slave.service 到 /etc/systemd/system/webpc-slave.service 爲 systemd 安裝 slave 服務

    ```
    sudo cp /opt/webpc/webpc-slave.service /etc/systemd/system/webpc-slave.service
    ```

    如果修改了安裝目錄，記得要修改 webpc-slave.service 中對應的路徑

3. 創建 webpc 用戶

    ```
    sudo useradd webpc -Mrs /sbin/nologin
    ```

4. 運行服務 

    ```
    sudo systemctl start webpc-slave.service
    ```

### windows-amd64

#### master

1. 下載最新程式並且解壓
2. 以管理員身份運行 controller.bat 後輸入指令 **1** 安裝 master 服務
3. 在 windows 服務管理器中運行 webpc-master-service 服務

#### slave

1. 下載最新程式並且解壓
2. 以管理員身份運行 controller.bat 後輸入指令 **2** 安裝 slave 服務
3. 在 windows 服務管理器中運行 webpc-slave-service 服務

## 設定

* WebPC 使用 jsonnet 作爲設定檔案
* **etc/master.jsonnet** 是 master 的主設定，其 import 了 **etc/master** 檔案夾下的子設定
* **etc/slave.jsonnet** 是 slave 的主設定，其 import 了 **etc/slave** 檔案夾下的子設定

### 設定master

http.libsonnet 定義了 master 如何提供網站服務

```
local def = import "def.libsonnet";
local size = def.Size;
local duration = def.Duration;
{
    // http addr
    Addr: ':9000',
    // if not empty use https
    CertFile: '',
    KeyFile: '',
    // enable swagger-ui on /document/
    Swagger: true,
    // grpc server option
    Option: {
        WriteBufferSize: 32*size.KB,
        ReadBufferSize: 32*size.KB,
        InitialWindowSize: 0*size.KB, // < 64k ignored
        InitialConnWindowSize: 0*size.KB, // < 64k ignored
        MaxRecvMsgSize: 0, // <1 6mb
        MaxSendMsgSize: 0, // <1 math.MaxInt32
        MaxConcurrentStreams: 0,
        ConnectionTimeout: 120 * duration.Second,
        Keepalive: {
            MaxConnectionIdle: 0,
            MaxConnectionAge: 0,
            MaxConnectionAgeGrace: 0,
            Time: 0,
            Timeout: 0,
        },
    },
}
```

system.libsonnet 定義將 master 本身作爲一個slave 註冊給自己以提供遠程控制，只有 Root 權限的用戶可以控制 master

### 設定slave

connect.libsonnet 設置了受控端到何處註冊自己，需要先在網頁服務中添加一個遠控設備，並複製註冊地址填寫爲此檔案的URL屬性

```
local def = import "def.libsonnet";
local size = def.Size;
local duration = def.Duration;
{
    // http addr
    URL: 'ws://127.0.0.1:9000/api/v1/dialer/64048031f73a11eba3890242ac120064',
    // if true allow insecure server connections when using SSL
    // Insecure: true, 
    // grpc server option
    Option: {
        WriteBufferSize: 32*size.KB,
        ReadBufferSize: 32*size.KB,
        InitialWindowSize: 0*size.KB, // < 64k ignored
        InitialConnWindowSize: 0*size.KB, // < 64k ignored
        MaxRecvMsgSize: 0, // <1 6mb
        MaxSendMsgSize: 0, // <1 math.MaxInt32
        MaxConcurrentStreams: 0,
        ConnectionTimeout: 120 * duration.Second,
        Keepalive: {
            MaxConnectionIdle: 0,
            MaxConnectionAge: 0,
            MaxConnectionAgeGrace: 0,
            Time: 0,
            Timeout: 0,
        },
    },
}
```

system.libsonnet 設置了slave 將提供哪些功能

```
{
    //Shell : "shell-linux", // if empty, linux default value shell-linux.sh
    //Shell : "shell-windows.bat", // if empty, windows default value shell-windows.bat
    // vnc server address
    VNC: "127.0.0.1:5900",
    // mount path to web
    Mount: [
        {
            // web display name
            Name: "s_movie",
            // local filesystem path
            Root: "/home/dev/movie",
            // Set the directory to be readable. Users with read/write permissions can read files
            Read: true,
            // Set the directory to be writable. Users with write permission can write files
            // If Write is true, Read will be forcibly set to true
            Write: true,
            // Set as a shared directory to allow anyone to read the file
            // If Shared is true, Read will be forcibly set to true
            Shared: true,
        },
        {
            Name: "s_home",
            Root: "/home/dev",
            Write: true,
            Read: true,
            Shared: false,
        },
        {
            Name: "s_root",
            Root: "/",
            Write: false,
            Read: true,
            Shared: false,
        },
        {
            Name: "s_media",
            Root: "/media/dev/",
            Write: false,
            Read: true,
            Shared: false,
        },
    ],
}
```

## 編譯

WebPC 後端使用 golang 和 grpc 開發，前端使用 angular 開發需要分別編譯。

### 編譯前端

1. 安裝好必要的開發環境 node typescript yarn

2. 下載項目並且切換工作目錄到 webpc/view

    ```
    git clone git@github.com:powerpuffpenguin/webpc.git && cd webpc/view
    ```

3. 安裝項目依賴

    ```
    yarn install
    ```

    or

    ```
    npm install
    ```

4. 編譯前端代碼

    ```
    ../build.sh view
    ```

### 編譯WebPC

1. 安裝好必要的開發環境 gcc golang proto3 grpc protoc-gen-go protoc-gen-grpc-gateway protoc-gen-openapiv2 

2. 安裝 golang 代碼嵌入工具

    ```
    go install github.com/rakyll/statik
    ```

3. 下載項目並且切換工作目錄到 webpc

    ```
    git clone git@github.com:powerpuffpenguin/webpc.git && cd webpc
    ```

4. 生成 grpc 代碼

    ```
    ./build.sh grpc
    ```

5. 將前端代碼和靜態檔案嵌入到 golang 代碼

    ```
    ./build.sh document
    ./build.sh view -s
    ```

6. 編譯 go 代碼

    ```
    go build -o bin/webpc
    ```