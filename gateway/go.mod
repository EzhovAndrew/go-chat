module github.com/go-chat/gateway

go 1.24.0

toolchain go1.24.3

replace (
	github.com/go-chat/auth => ../auth
	github.com/go-chat/chat => ../chat
	github.com/go-chat/notifications => ../notifications
	github.com/go-chat/social => ../social
	github.com/go-chat/users => ../users
)

require (
	github.com/go-chat/auth v0.0.0-00010101000000-000000000000
	github.com/go-chat/chat v0.0.0-00010101000000-000000000000
	github.com/go-chat/notifications v0.0.0-00010101000000-000000000000
	github.com/go-chat/social v0.0.0-00010101000000-000000000000
	github.com/go-chat/users v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	google.golang.org/grpc v1.76.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.10-20250912141014-52f32327d4b0.1 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251014184007-4626949a642f // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
