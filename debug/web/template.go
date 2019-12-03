package web

const dashboard = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta charset="utf-8">
    <title>Micro Web Debug.Stats</title>
    <meta name="application-name" content="netdata">
</head>
<script type="text/javascript" src="dashboard.js?v20190902-0"></script>
<body>

<div style="width: 100%; text-align: center;">
    <div data-netdata="netdata.server_cpu"
            data-dimensions="user"
            data-chart-library="gauge"
            data-width="150px"
            data-after="-60"
            data-points="60"
            data-title="Yes! Realtime!"
            data-units="I am alive!"
            data-colors="#FF5555"
            ></div>
    <br/>
    <div data-netdata="netdata.server_cpu"
            data-dimensions="user"
            data-chart-library="dygraph"
            data-dygraph-theme="sparkline"
            data-width="200px"
            data-height="30px"
            data-after="-60"
            data-points="60"
            data-colors="#FF5555"
            ></div>
</div>
</body>
</html>
`
