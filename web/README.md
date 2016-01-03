# Micro Web

Micro web provides a visual point of entry for the micro environment and should replicate 
the features of the CLI.

It also includes a reverse proxy to route requests to micro web 
apps. /[name] will proxy to the service [namespace].[name]. The default namespace is 
go.micro.web.

## Run Web UI
```bash
$ go get github.com/micro/micro
$ micro web
```

Browse to localhost:8082

<img src="https://github.com/micro/micro/blob/master/web/web1.png">
-
<img src="https://github.com/micro/micro/blob/master/web/web2.png">
-
<img src="https://github.com/micro/micro/blob/master/web/web3.png">
