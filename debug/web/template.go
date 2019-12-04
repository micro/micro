package web

import "html/template"

var dashboardTemplate *template.Template

const dashboardText = `<!DOCTYPE html>
<html lang="en">
<head>
    <title>Micro Debug.Stats</title>
    <meta name="application-name" content="netdata">
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
</head>
<body>
<div class="container-fluid">
<h1>Micro Debug.Stats <div data-netdata="system.cpu" data-chart-library="sparkline" data-height="30" data-after="-600" data-sparkline-linecolor="#888"></div></h1>

<div style="width: 33%; display: inline-block;">
    <div data-netdata="go_micro_services.micro_service_memory"
        data-chart-library="dygraph"
        data-width="100%"
        data-height="300"
        data-after="-600"
        ></div>
</div>

<div style="width: 33%; display: inline-block;">
    <div data-netdata="go_micro_services.micro_service_threads"
        data-chart-library="dygraph"
        data-width="100%"
        data-height="300"
        data-after="-600"
        ></div>
</div>

<div style="width: 33%; display: inline-block;">
    <div data-netdata="go_micro_services.micro_service_gcrate"
        data-chart-library="dygraph"
        data-width="100%"
        data-height="300"
        data-after="-600"
        ></div>
</div>

<br />

<div style="width: 33%; display: inline-block;">
    <div data-netdata="go_micro_services.micro_service_uptime"
        data-chart-library="dygraph"
        data-width="100%"
        data-height="300"
        data-after="-600"
        ></div>
</div>

</div>
</body>
</html>
<script type="text/javascript" src="dashboard.js?v20190902-0"></script>
`
