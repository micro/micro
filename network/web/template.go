package web

var (
	templateFile = `
<html>
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
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
  </style>
</head>
<body>
  <div id="head">
    <ul class="nav">
      <li><a href="/network/nodes">Nodes</a></li>
      <li><a href="/network/routes">Routes</a></li>
      <li><a href="/network/services">Services</a></li>
    </ul>
  </div>
  <div id="content">
    {{.}}
  </div>
</body>
</html>
`
)
