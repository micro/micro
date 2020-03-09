package web

const templates = `
	{{define "basic"}}
		<html>
		<head>
			<style>
				.inner {
					position: absolute;
					left: 50%;
					top: 50%;
					transform: translate(-50%, -50%);
					max-width: 100vw;
					width: 400px;
				}

				form  {
					display: flex;
					flex-direction: column;
				}

				input {
					margin-top:  5px;
					margin-bottom: 20px;
					outline: none;
					height: 25px;
				}

				input[type=submit] {
					
				}
			</style>
		</head>
		<body>
			<div class='inner'>
				<h1>Login</h1>
				<form method='post'>
					<label for='email'>Email</label>
					<input type='email' name='email' required />

					<label for='password'>Password</label>
					<input type='password' name='password' required />

					<input type='submit' value='Submit' />
				</form>
			</div>
		</body>
		</html>
	{{end}}
`
