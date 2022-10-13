var namespace = "micro";
var url = "https://api.m3o.com";
var cookie = "micro-token";

function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
}

function setCookie(name, value, expiry) {
  const d = new Date();
  d.setTime(d.getTime() + (expiry*1000));
  let expires = "expires="+ d.toUTCString();
  document.cookie = name + "=" + value + ";" + expires + ";path=/";
}

async function call(url = '', data = {}) {
  var token = getCookie(cookie);

  // Default options are marked with *
  const response = await fetch(url, {
    method: 'POST', // *GET, POST, PUT, DELETE, etc.
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + token,
      'Micro-Namespace': namespace,
    },
    body: JSON.stringify(data) // body data type must match "Content-Type" header
  });
  return response.json(); // parses JSON response into native JavaScript objects
}

async function login(username, password) {
	return call(url + '/auth/Token', {
		"id": username,
		"secret": password,
		"token_expiry": 30 * 86400,
	}).then(function(rsp) {
		setCookie(cookie, rsp.token.access_token, rsp.token.expiry);
	});
}

async function logout() {
	setCookie(cookie, "", 0);
}

async function setURL(u) {
	url = u;	
}

function submitLogin(form) {
	var username = document.getElementById("username").value;
	var password = document.getElementById("password").value;
	login(username, password);
}

function submitLogout(form) {
	logout()
}
