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
    .nav {
      margin-bottom: 10px;
      padding: 0;
    }
    .nav li {
      display: inline-block;
    }
    table td {
      padding-right: 5px;
    }
    .graphs {
      position: relative;
      margin: 0 auto;
      max-width: 1000px;
    }
    .graph {
      width: 500px;
      display: inline-block;
      margin: 20px;
    }
  </style>
</head>
<body>
  <h1>Debug</h1>

  <div id="head">
    <ul class="nav">
      <li><a href="/debug/">Stats</a></li>
    </ul>
  </div>
  <div id="content">
	<div data-netdata="system.cpu" data-chart-library="sparkline" data-height="30" data-after="-600" data-sparkline-linecolor="#888"></div>
      <div id="graphs">
        <p>&nbsp;</p>
	<div class="graph">
	    <div data-netdata="go_micro_services.micro_service_memory"
		data-chart-library="dygraph"
		data-width="100%"
		data-height="300"
		data-after="-600"
		></div>
	</div>

	<div class="graph">
	    <div data-netdata="go_micro_services.micro_service_threads"
		data-chart-library="dygraph"
		data-width="100%"
		data-height="300"
		data-after="-600"
		></div>
	</div>

	<div class="graph">
	    <div data-netdata="go_micro_services.micro_service_gcrate"
		data-chart-library="dygraph"
		data-width="100%"
		data-height="300"
		data-after="-600"
		></div>
	</div>

	<div class="graph">
	    <div data-netdata="go_micro_services.micro_service_uptime"
		data-chart-library="dygraph"
		data-width="100%"
		data-height="300"
		data-after="-600"
		></div>
	</div>
    </div>
  </div>
  <script type="text/javascript" src="dashboard.js?v20190902-0"></script>
</body>
</html>
`
)
