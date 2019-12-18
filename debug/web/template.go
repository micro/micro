package web

import (
	"html/template"
)

var (
	dashboardTemplate *template.Template

	dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <title>Micro Debug</title>
    <meta name="application-name" content="netdata">
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
  <link href="https://fonts.googleapis.com/css?family=Source+Code+Pro&display=swap" rel="stylesheet">
  <style>
    html {
      font-family: 'Source Code Pro', monospace;
    }
    table td {
      padding-right: 5px;
    }
    #graphs {
      text-align: center;
    }
    .graph {
      width: 500px;
      display: inline-block;
      margin: 20px;
    }
  </style>
</head>
<body style="font-family: 'Source Code Pro', monospace; margin: 10px;">
  <h1 style="vertical-align: middle;">
    <a href="/"><img src="https://micro.mu/logo.png" height=50px width=auto /></a> Debug
  </h1>
  <p>&nbsp;</p>
  <div id="content">
      <!--
      <div data-netdata="system.cpu" data-chart-library="sparkline" data-height="30" data-after="-600" data-sparkline-linecolor="#888"></div>
      -->
      <div id="graphs">
  <div class="graph">
      <div data-netdata="go_micro_services.micro_service_memory"
    data-chart-library="dygraph"
    data-width="100%"
    data-height="300"
    data-after="-600"{{ if .Service }}
    data-dimensions="{{.Service}}*"{{end}}
    data-title="Memory Usage"
    ></div>
  </div>

  <div class="graph">
      <div data-netdata="go_micro_services.micro_service_threads"
    data-chart-library="dygraph"
    data-width="100%"
    data-height="300"
    data-after="-600"{{ if .Service }}
    data-dimensions="{{.Service}}*"{{end}}
    data-title="Go Routines"
    ></div>
  </div>

  <div class="graph">
      <div data-netdata="go_micro_services.micro_service_gcrate"
    data-chart-library="dygraph"
    data-width="100%"
    data-height="300"
    data-after="-600"{{ if .Service }}
    data-dimensions="{{.Service}}*"{{end}}
    data-title="GC Pause"
    ></div>
  </div>

  <div class="graph">
      <div data-netdata="go_micro_services.micro_service_uptime"
    data-chart-library="dygraph"
    data-width="100%"
    data-height="300"
    data-after="-600"{{ if .Service }}
    data-dimensions="{{.Service}}*"{{end}}
    data-title="Uptime"
    ></div>
  </div>
    </div>
  </div>
  <script type="text/javascript" src="/debug/dashboard.js?v20190902-0"></script>
</body>
</html>
`

	logTemplate = `
<html lang="en">
<head>
    <title>Micro Debug | Log</title>
    <meta name="application-name" content="netdata">
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
  <link href="https://fonts.googleapis.com/css?family=Source+Code+Pro&display=swap" rel="stylesheet">
  <style>
    html {
      font-family: 'Source Code Pro', monospace;
    }
    
  </style>
</head>
<body style="font-family: 'Source Code Pro', monospace; margin: 10px;">
  <h1 style="vertical-align: middle; font-weight: 500;">
    <a href="/"><img src="https://micro.mu/logo.png" height=50px width=auto style="vertical-align: middle;"/></a> Debug Log
  </h1>
  <p>&nbsp;</p>
  <div id="content">
    {{ range $index, $el := .Records }}
    <div>{{.}}</div>
    {{end}}
  </div>
</body>
</html>
`
)
