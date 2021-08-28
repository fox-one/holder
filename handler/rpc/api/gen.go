package api

//go:generate protoc -I . holder.proto --twirp_out=. --go_out=.
//go:generate protoc-go-inject-tag -input=./holder.pb.go
