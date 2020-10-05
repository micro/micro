module example

go 1.13

require (
	github.com/bwmarrin/discordgo v0.20.2 // indirect
	github.com/forestgiant/sliceutil v0.0.0-20160425183142-94783f95db6c // indirect
	github.com/go-git/go-git/v5 v5.1.0 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v3 v3.0.0-beta.0.20200821101742-6cda6ef92e50
	github.com/micro/micro/v3 v3.0.0-20200728090928-ad22505562c9
	github.com/nlopes/slack v0.6.1-0.20191106133607-d06c2a2b3249 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/telegram-bot-api.v4 v4.6.4 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../../..
