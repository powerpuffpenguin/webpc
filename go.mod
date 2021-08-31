module github.com/powerpuffpenguin/webpc

go 1.16

require (
	github.com/andybalholm/brotli v1.0.3
	github.com/creack/pty v1.1.14 // indirect
	github.com/denisenkom/go-mssqldb v0.10.0
	github.com/gin-gonic/gin v1.7.3
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/go-redis/redis/v8 v8.11.2 // indirect
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-jsonnet v0.17.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/iamacarpet/go-winpty v1.0.2 // indirect
	github.com/lib/pq v1.10.2
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/powerpuffpenguin/sessionid v0.0.0-20210526073144-ca6285fc8dec
	github.com/powerpuffpenguin/sessionid_server v1.0.1
	github.com/powerpuffpenguin/vnet v0.0.0-20210802025916-207c679e83b8
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/ugorji/go v1.2.6 // indirect
	go.opentelemetry.io/otel/internal/metric v0.22.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.0
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210811021853-ddbe55d93216
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	xorm.io/xorm v1.2.2
)

replace github.com/andybalholm/brotli v1.0.3 => github.com/zuiwuchang/brotli v0.0.0-20210825022947-ec682fbe0a16
