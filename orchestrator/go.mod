module github.com/vitelabs/vite-portal/orchestrator

go 1.18

require (
	github.com/gorilla/websocket v1.5.0
	github.com/spf13/cobra v1.5.0
	github.com/vitelabs/vite-portal/shared v0.0.0
)

require (
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/rs/zerolog v1.27.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace github.com/vitelabs/vite-portal/shared v0.0.0 => ../shared
