# v1.3.1

* xterm: use xterm-addon-webgl instead of xterm-addon-canvas
* server: resize rows-1/rows on shell attach

# v1.3.0

* upgrade golang to 1.21
* fix jsonnet import  abs path on windows
* upgrade github.com/gorilla/websocket to v1.5.1
* upgrade github.com/gin-gonic/gin to v1.9.1
* "attachment;filename="+url.PathEscape(name)
* upgrade angular to 14
* angular target to es2020
* rm require.min.js
* upgrade noVNC to 1.4
* upgrade xterm to 5.3.0


# v1.2.8
* cache session

# v1.2.7
* fix shared

# v1.2.6

* fix forward/socks5 RefreshToken error
* new session module
* movie list

# v1.1.2

* fix fs module download not set header Content-Disposition for filename

# v1.1.1

* Use 'embed' added in golang 1.6 instead of 'statik'
* Upgrade golang to 1.17
* GRPC download not use chunked,used context-length

# v1.1.0

* web view auto reset session on token not exists 
* port forwarding
* socks5 socks4 socks4a
* automatic upgrades

# v1.0.8

* fix view navigation-bar server check auth.Server 
* fix run-linux not set $HOME
* auth Slave

# v1.0.5

* fix shared not work
* fix slave websocket not check group
* fix some fs api not forward
* function list on view's navigation-bar

# v1.0.0

First release

* filesystem
* filesystem shared 
* shell
* vnc
* logger
