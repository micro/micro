module github.com/micro/micro/v3

go 1.15

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/caddyserver/certmagic v0.10.6
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dustin/go-humanize v1.0.0
	github.com/evanphx/json-patch/v5 v5.0.0
	github.com/getkin/kin-openapi v0.26.0
	github.com/go-acme/lego/v3 v3.4.0
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/go-version v1.2.1
	github.com/hpcloud/tail v1.0.0
	github.com/improbable-eng/grpc-web v0.13.0
	github.com/kr/pretty v0.2.0
	github.com/m3o/platform v0.0.0-20210512105552-e6700d708e74 // indirect
	github.com/m3o/platform/profile/ci v0.0.0-20210512103138-4a540cec553c // indirect
	github.com/m3o/platform/profile/platform v0.0.0-20210512103138-4a540cec553c
	github.com/miekg/dns v1.1.27
	github.com/nightlyone/lockfile v1.0.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/onsi/gomega v1.7.1
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/rhysd/go-github-selfupdate v1.2.2
	github.com/serenize/snaker v0.0.0-20171204205717-a683aaf2d516
	github.com/stoewer/go-strcase v1.2.0
	github.com/stretchr/objx v0.1.1
	github.com/stretchr/testify v1.7.0
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	github.com/urfave/cli/v2 v2.3.0
	github.com/xanzy/go-gitlab v0.35.1
	github.com/xlab/treeprint v0.0.0-20181112141820-a009c3971eca
	go.etcd.io/bbolt v1.3.5
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
