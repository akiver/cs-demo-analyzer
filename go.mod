module github.com/akiver/cs-demo-analyzer

go 1.27.0

require (
	github.com/golang/geo v0.0.0-20261106073233-5e8d6967632a
	github.com/klauspost/compress v1.18.0
	github.com/markus-wa/demoinfocs-golang/v4 v4.5.1
	github.com/markus-wa/gobitread v0.2.6
	github.com/pkg/errors v0.9.1
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/golang/snappy v1.0.0 // indirect
	github.com/markus-wa/go-unassert v0.1.3 // indirect
	github.com/markus-wa/godispatch v1.5.1 // indirect
	github.com/markus-wa/ice-cipher-go v1.0.0 // indirect
	github.com/markus-wa/quickhull-go/v2 v2.2.0 // indirect
)

replace github.com/markus-wa/demoinfocs-golang/v4 v4.5.1 => github.com/markus-wa/demoinfocs-golang/v4 v4.5.2-0.20261122234809-b7de0dd59922
