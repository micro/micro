module github.com/micro/micro/v2

go 1.13

require (
	github.com/boltdb/bolt v1.3.1
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/dustin/go-humanize v1.0.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-acme/lego/v3 v3.4.0
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/go-version v1.2.1
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/netdata/go-orchestrator v0.0.0-20190905093727-c793edba0e8f
	github.com/olekukonko/tablewriter v0.0.4
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/serenize/snaker v0.0.0-20171204205717-a683aaf2d516
	github.com/stretchr/testify v1.5.1
	github.com/stripe/stripe-go/v71 v71.28.0
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/xlab/treeprint v0.0.0-20181112141820-a009c3971eca
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2
	golang.org/x/tools v0.0.0-20191216173652-a0e659d51361
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/square/go-jose.v2 v2.4.1 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

//replace github.com/micro/go-micro/v2 => /home/crufter/code/src/github.com/micro/go-micro

replace github.com/micro/go-micro/v2 => github.com/micro/go-micro/v2 k8supload
