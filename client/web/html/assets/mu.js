var namespace = "micro";
var url = "http://localhost:8080";
var cookie = "micro-token";
var services = {};

Object.defineProperty(String.prototype, 'capitalize', {
    value: function() {
        return this.charAt(0).toUpperCase() + this.slice(1);
    },
    enumerable: false
});

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

function setCookie(name, value, expiry) {
    const d = new Date();
    d.setTime(d.getTime() + (expiry * 1000));
    let expires = "expires=" + d.toUTCString();
    document.cookie = name + "=" + value + ";" + expires + ";path=/";
}

async function call(service = '', endpoint = '', method = '', data = {}) {
    var token = getCookie(cookie);

    // Default options are marked with *
    const response = await fetch(`${url}/${service}/${endpoint}/${method}`, {
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
    return call('auth', 'Auth', 'Token', {
        "id": username,
        "secret": password,
        "token_expiry": 30 * 86400,
    }).then(function(rsp) {
        setCookie(cookie, rsp.token.access_token, rsp.token.expiry);
        window.location.href = "/";
    });
}

async function logout() {
    setCookie(cookie, "", 0);
}

async function listServices() {
    return call('registry', 'Registry', 'ListServices')
        .then(function(rsp) {
            return rsp.services;
        });
}

async function getService(name) {
    return call('registry', 'Registry', 'GetService', {
        "service": name
    });
}

async function setURL(u) {
    url = u;
}

function renderLogin() {
    var div = document.getElementById("login");

    var h3 = document.createElement("h3");
    h3.innerText = "Login";

    var form = document.createElement("form")
    form.onsubmit = submitLogin;

    var label1 = document.createElement("label");
    label1.innerText = "Username";
    var label2 = document.createElement("label");
    label2.innerText = "Password";

    var username = document.createElement("input");
    username.id = "username";
    username.type = "username";
    username.name = "username";
    username.required = true;

    var password = document.createElement("input");
    password.id = "password";
    password.type = "password";
    password.name = "password";
    password.required = true;

    var submit = document.createElement("button")
    submit.innerText = "Submit";

    div.appendChild(h3);
    form.appendChild(label1);
    form.appendChild(username);
    form.appendChild(label2);
    form.appendChild(password);
    form.appendChild(submit);
    div.appendChild(form);
}

function renderServices(fn) {
    var search = function() {
        var refs = $('a[data-filter]');
        $('.search').on('keyup', function() {
            var val = $.trim(this.value);
            refs.hide();
            refs.filter(function() {
                return $(this).data('filter').search(val) >= 0
            }).show();
        });

        $('.search').on('keypress', function(e) {
            if (e.which != 13) {
                return;
            }
            $('.service').each(function() {
                if ($(this).css('display') == "none") {
                    return;
                }
                window.location.href = $(this).attr('href');
            })
        });
    };

    // load into the #services div
    var heading = document.getElementById("heading");
    var content = document.getElementById("content");
    var service = document.createElement("div");
    service.id = "services";

    var render = function(rsp) {
        rsp.forEach(function(srv) {
            var a = document.createElement("a");
            a.href = "/" + srv.name;
            a.setAttribute("data-filter", srv.name);
            a.setAttribute("class", "service");
            a.innerText = srv.name;
            service.appendChild(a);
        });

        // setup search filtering
        search();

        // execute user defined function
        if (fn != undefined) {
            fn();
        }
    }

    // the search box
    // <h4><input class="input-lg search" type=text placeholder="Search" autofocus></h4>
    var input = document.createElement("input")
    input.setAttribute("class", "input-lg search");
    input.type = "text"
    input.placeholder = "Search"
    input.autofocus = true;

    // render from the cache
    if (services.length > 0) {
        return render(services);
    }

    // call the backend
    listServices().then(function(rsp) {
        // reset content
        heading.innerHTML = "";
        content.innerHTML = "";

        // append the search box
        heading.appendChild(input);

        // append services to content
        content.appendChild(service);

        // cache the list for next time
        rsp.forEach(function(srv) {
            services[srv.name] = srv;
        });

        // render the content
        render(rsp);
    });
}

function renderService(service) {
    getService(service)
        .then(function(rsp) {
            console.log("rendering service", service);
            var heading = document.getElementById("heading");
            heading.innerText = service;
            var content = document.getElementById("content");
            content.innerHTML = "";
            var div = document.createElement("div");
            div.id = "service";
            content.appendChild(div);

            var eps = {};

            rsp.services[0].endpoints.forEach(function(endpoint) {
                var parts = endpoint.name.split(".");
                var name = parts[1];

                // eg auth != accounts
                if (service != parts[0].toLowerCase()) {
                    name = parts[0] + "/" + parts[1];
                }
                // define a new div for the endpoint
                var ep = document.createElement("div")
                ep.setAttribute("class", "endpoint");

                // create the endpoint link
                var a = document.createElement("a");
                a.href = "/" + service + "/" + name;
                a.innerText = name;

                // set the content
                ep.appendChild(a);
                div.appendChild(ep);
            })

        });
}

function renderEndpoint(service, endpoint, method) {
    getService(service)
        .then(function(rsp) {
            console.log("rendering", service, endpoint, method);
            var heading = document.getElementById("heading");
            heading.innerText = service + " " + endpoint;
            var content = document.getElementById("content");
            content.innerHTML = "";
            var request = document.createElement("div");
            request.id = "request";
            var response = document.createElement("div");
            response.id = "response";

            content.appendChild(request);
            content.appendChild(response);

            // construct the endpoint
            var name = service.capitalize() + "." + endpoint;
            if (method != undefined) {
                name = endpoint + "." + method;
                heading.innerText += " " + method;
            } else {
                method = endpoint;
                endpoint = service.capitalize();
            }

            rsp.services[0].endpoints.forEach(function(ep) {
                // render the form
                if (name == ep.name) {
                    console.log("rending endpoint", ep.name);
                    // get request info
                    // render a form
                    // render output
                    var form = document.createElement("form");

                    form.onsubmit = function(ev) {
                        ev.preventDefault();

                        // build request
                        var request = {};

                        for (i = 0; i < form.elements.length; i++) {
                            var entry = form.elements[i];
                            if (entry.name.length == 0) {
                                continue
                            }
                            if (entry.value.length == 0) {
                                continue
                            }
                            request[entry.name] = entry.value;
                        }

                        // renders the output recursively
                        var render = function(key, val, depth) {
                            // print the value if it's not an object
                            var print = function(key, val) {
                                var value = document.createElement("div");
                                value.setAttribute("class", "field");
                                key = key.capitalize();
                                value.innerText = `${key}: ${val}`
                                return value;
                            }
                            // not an object, just print it
                            if (typeof val != "object") {
                                return print(key, val)
                            }
                            // if it's an array then check types and print
                            if (val.constructor == Array) {
                                if (typeof val[0] != "object") {
                                    return print(key, val);
                                }
                            }

                            // is actually a number
                            if (val.constructor == Number) {
                                return print(key, val);
                            }

                            var output = document.createElement("div");
                            output.setAttribute("class", "response");

                            // iterate the objects as needed
                            for (const [key, value] of Object.entries(val)) {
                                // append the next output value
                                depth++
                                output.appendChild(render(key, value, depth));
                            }

                            // return the entire rendered value
                            return output;
                        }

                        call(service, endpoint, method, request)
                            .then(function(rsp) {
                                // prepend to response
                                response.innerText = '';
                                // set a title
                                var h4 = document.createElement("h4");
                                h4.setAttribute("class", "title");
                                h4.innerText = "Response";
                                response.appendChild(h4);
                                // render values
                                response.appendChild(render("response", rsp, 0));
                            });
                    };

                    ep.request.values.forEach(function(value) {
                        var input = document.createElement("input");
                        input.id = value.name
                        input.type = "text";
                        input.name = value.name;
                        input.placeholder = value.name;
                        input.autocomplete = "off";
                        form.appendChild(input);
                    });

                    // generate the button
                    var submit = document.createElement("button")
                    submit.innerText = "Submit";
                    form.appendChild(submit);
                    request.appendChild(form);
                }
                // end forEach
            })
            // end Promise
        });
}

function submitLogin(form) {
    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;
    login(username, password);
    return false;
}

function submitLogout(form) {
    logout()
}

function main() {
    // parse the url
    if (window.location.pathname == "/") {
        console.log("render services");
        return renderServices();
    }

    var parts = window.location.pathname.split("/")

    // process service
    if (parts.length == 2) {
        console.log("render service", parts[1]);
        renderService(parts[1]);
    }

    if (parts.length == 3) {
        console.log("render endpoint", parts[1], parts[2]);
        renderEndpoint(parts[1], parts[2]);
    }

    if (parts.length == 4) {
        console.log("render endpoint", parts[1], parts[2], parts[3]);
        renderEndpoint(parts[1], parts[2], parts[3]);
    }
}
