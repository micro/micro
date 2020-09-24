module logger

go 1.13

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/micro/v3 => ../../..

require (
	github.com/bwmarrin/discordgo v0.20.2 // indirect
	github.com/go-git/go-git/v5 v5.1.0 // indirect
	github.com/lucas-clemente/quic-go v0.14.1 // indirect
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200917131714-7750f542b4c2
	github.com/micro/micro/v3 v3.0.0-00010101000000-000000000000
	github.com/nlopes/slack v0.6.1-0.20191106133607-d06c2a2b3249 // indirect
)
