package cmd

var serviceTemplate = `
---
weight: 11
title: {{ .serviceName }}
---
# {{ .serviceName }}
{{ range $rpc := .rpcs }}
## {{ (parentService $rpc).Name }}.{{ $rpc.Name }}
` + "```" + `go
package main
import (
  "github.com/micro/clients/go/client"
  {{ packageNamify $.serviceName }}_proto "{{ $.goImportPrefix }}/{{ $.serviceName }}/proto"
)
func main() {
  c := client.NewClient(nil)
  req := {{ packageNamify $.serviceName }}_proto.{{ $rpc.RequestType }}{}
  rsp := {{ packageNamify $.serviceName }}_proto.{{ $rpc.ReturnsType }}{}
  if err := c.Call("go.micro.srv.{{ $.serviceName }}", "{{ (parentService $rpc).Name }}.{{ $rpc.Name }}", req, &rsp); err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(rsp)
}
` + "```" + `
` + "```" + `javascript
// To install "npm install --save @microhq/ng-client"
import { Component, OnInit } from "@angular/core";
import { ClientService } from "@microhq/ng-client";
@Component({
  selector: "app-example",
  templateUrl: "./example.component.html",
  styleUrls: ["./example.component.css"]
})
export class ExampleComponent implements OnInit {
  constructor(private mc: ClientService) {}
  ngOnInit() {
    this.mc
      .call("go.micro.srv.{{ $.serviceName }}", "{{ (parentService $rpc).Name }}.{{ $rpc.Name }}", {})
      .then((response: any) => {
        console.log(response)
      });
  }
}
` + "```" + `
{{ commentLines $rpc.Comment }}
### Request Parameters
Name |  Type | Description
--------- | --------- | ---------
{{ range $field := (getNormalFields $rpc.RequestType) }}{{ $field.Name }} | {{ $field.Type }} | {{ commentLines $field.Comment }}
{{ end }}
### Response Parameters
Name |  Type | Description
--------- | --------- | ---------
{{ range $field := (getNormalFields $rpc.ReturnsType) }}{{ $field.Name }} | {{ $field.Type }} | {{ commentLines $field.Comment }}
{{ end }}
{{ range $messageName := (messagesUsedByReqRsp $rpc) }}
### Message {{ $messageName }}
Name |  Type | Description
--------- | --------- | ---------
{{ range $field := (getNormalFields $messageName) }}{{ $field.Name }} | {{ $field.Type }} | {{ commentLines $field.Comment }}
{{ end }}
{{ end }}
### 
<aside class="success">
Remember — a happy kitten is an authenticated kitten!
</aside>
{{ end }}
`

// > The above command returns JSON structured like this:
//
//` + "```" + `json
//{ sadasas }
//` + "```" + `
