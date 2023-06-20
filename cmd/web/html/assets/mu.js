var namespace = "micro";
var api = "http://localhost:8080";
var cookie = "micro-token";
var services = {};

Object.defineProperty(String.prototype, 'capitalize', {
    value: function() {
        return this.charAt(0).toUpperCase() + this.slice(1);
    },
    enumerable: false
});

String.prototype.parseURL = function(embed) {
    return this.replace(/[A-Za-z]+:\/\/[A-Za-z0-9-_]+\.[A-Za-z0-9-_:%&~\?\/.=]+/g, function(url) {
        if (embed == true) {
            var match = url.match(/^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=)([^#\&\?]*).*/);
            if (match && match[2].length == 11) {
                return '<div class="iframe">' +
                    '<iframe src="//www.youtube.com/embed/' + match[2] +
                    '" frameborder="0" allowfullscreen></iframe>' + '</div>';
            };
            if (url.match(/^.*giphy.com\/media\/[a-zA-Z0-9]+\/[a-zA-Z0-9]+.gif$/)) {
                return '<div class="animation"><img src="' + url + '"></div>';
            }
        };
        // var pretty = url.replace(/^http(s)?:\/\/(www\.)?/, '');
        return url.link(url);
    }).replace(/[\@\#]([a-zA-z0-9_]+)/g, function(m, m1) {
        var t = '<a href="http://twitter.com/';
        if (m.charAt(0) == '#')
            t += 'hashtag/';
        return t + encodeURI(m1) + '" target="_blank">' + m + '</a>';
    });
};

function generateURL(service = '', endpoint = '', query = []) {
    var u = `/${service}/${endpoint}`;
    var q = "";

    query.forEach(function(val) {
        if (q.length == 0) {
            q = "?" + val;
            return
        }
        q += "&" + val;
    })

    return u + q;
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

function setCookie(name, value, expiry) {
    const d = new Date();
    d.setTime(expiry * 1000);
    let expires = "expires=" + d.toUTCString();
    document.cookie = name + "=" + value + ";" + expires + ";path=/";
}

async function call(service = '', endpoint = '', method = '', data = {}) {
    var token = getCookie(cookie);

    var headers = {
        'Content-Type': 'application/json',
        'Micro-Namespace': namespace,
    }

    if (token != undefined) {
	headers['Authorization'] = 'Bearer ' + token;
    }

    // Default options are marked with *
    const response = await fetch(`${api}/${service}/${endpoint}/${method}`, {
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
        headers: headers,
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
        if (rsp.token == undefined) {
            var error = document.getElementById("error");
            error.innerText = rsp.detail;
        }

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

async function setAPI(u) {
    api = u;
}

function renderLogin() {
    var content = document.getElementById("content");
    content.innerHTML = '';

    var div = document.createElement("div");
    div.id = "login";

    content.appendChild(div);

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

function renderIndex(fn) {
    var search = function() {
        var refs = $('a[data-filter]');
        $('.search').on('keyup', function() {
            var val = $.trim(this.value.toLowerCase());
            if (val.length == 0) {
                $('#services').hide();
            } else {
                $('#services').show();
            }
            refs.hide();
            refs.filter(function() {
                return $(this).data('filter').search(val) >= 0
            }).show();
        });

        $('.search').on('keypress', function(e) {
            if (e.which != 13) {
                return;
            };
            var val = $.trim(this.value);
            var parts = val.split(" ");

            // assuming it's some full query
            if (parts.length > 1) {
                service = parts[0];
                endpoint = parts[1];
                request = [];

                // assemble a request
                parts.slice(2).forEach(function(val) {
                    if (val.split("=").length == 2) {
                        request.push(val);
                    }
                });

                window.location.href = generateURL(service, endpoint, request);
            }

            // partial string
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
    service.style.display = 'none';

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
    }

    // the search box
    // <h4><input class="input-lg search" type=text placeholder="Search" autofocus></h4>
    var input = document.createElement("input")
    input.setAttribute("class", "search");
    input.type = "text"
    input.placeholder = "Search"
    input.autofocus = true;

    // render from the cache
    if (services.length > 0) {
        return render(services);
    }

    // reset content
    heading.innerHTML = "";
    content.innerHTML = "";

    // append the search box
    heading.appendChild(input);

    // append services to content
    content.appendChild(service);


    // execute user defined function
    renderQueries();

    // call the backend
    listServices().then(function(rsp) {
        // cache the list for next time
        rsp.forEach(function(srv) {
            services[srv.name] = srv;
        });

        // render the content
        render(rsp);
    });
}

function renderServices() {
    var search = function() {
        var refs = $('a[data-filter]');
        $('.search').on('keyup', function() {
            var val = $.trim(this.value.toLowerCase());
            refs.hide();
            refs.filter(function() {
                return $(this).data('filter').search(val) >= 0
            }).show();
        });

        $('.search').on('keypress', function(e) {
            if (e.which != 13) {
                return;
            };
            var val = $.trim(this.value);
            var parts = val.split(" ");

            // assuming it's some full query
            if (parts.length > 1) {
                service = parts[0];
                endpoint = parts[1];
                request = [];

                // assemble a request
                parts.slice(2).forEach(function(val) {
                    if (val.split("=").length == 2) {
                        request.push(val);
                    }
                });

                window.location.href = generateURL(service, endpoint, request);
            }

            // partial string
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
    //service.id = "services";

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
    }

    // the search box
    // <h4><input class="input-lg search" type=text placeholder="Search" autofocus></h4>
    var input = document.createElement("input")
    input.setAttribute("class", "search");
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
        heading.innerHTML = "Services";
        content.innerHTML = "";

        // append the search box
        //heading.appendChild(input);

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
            heading.innerText = service + " / " + endpoint;
            var content = document.getElementById("content");
            content.innerHTML = "";
            var request = document.createElement("div");
            request.id = "request";
            var response = document.createElement("div");
            response.id = "response";

            content.appendChild(request);
            content.appendChild(response);

            // construct the endpoint
            var name = service.capitalize() + "." + endpoint.capitalize();
            var query = service + " " + endpoint;

            if (method != undefined) {
                name = endpoint.capitalize() + "." + method.capitalize();
                heading.innerText += " / " + method.capitalize();
                service += " " + method;
            } else {
                method = endpoint.capitalize();
                endpoint = service.capitalize();
            }

            const urlSearchParams = new URLSearchParams(window.location.search);
            const params = Object.fromEntries(urlSearchParams.entries());

            rsp.services[0].endpoints.forEach(function(ep) {
                // render the form
                if (name == ep.name) {
                    console.log("rending endpoint", ep.name);
                    // get request info
                    // render a form
                    // render output
                    var form = document.createElement("form");
                    form.id = name;

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

                        call(service, endpoint, method, request)
                            .then(function(rsp) {
                                renderResponse(response, rsp);
                            });
                    };

                    var submitForm = false;

                    ep.request.values.forEach(function(value, idx) {
                        // create a label
                        var label = document.createElement("label");
                        label.innerText = value.name.split("_").join(" ").capitalize();
                        // create the input
                        var input = document.createElement("input");
                        input.id = value.name
                        input.type = "text";
                        input.name = value.name;
                        input.placeholder = value.name;
                        input.autocomplete = "off";

                        if (idx == 0) {
                            input.autofocus = true;
                        }

                        if (params[value.name] != undefined) {
                            input.value = params[value.name];
                            submitForm = true;
                        }

                        form.appendChild(label);
                        form.appendChild(input);
                    });

                    // generate the button
                    var submit = document.createElement("button")
                    submit.innerText = "Submit";
                    form.appendChild(submit);
                    request.appendChild(form);

                    // auto-submit when we have form values
                    if (submitForm) {
                        $(form).submit();
                    }

                    // save the query
                    saveQuery(query);
                }
                // end forEach
            })
            // end Promise
        });
}

// renders the output recursively as a set of divs
function renderOutput(key, val, depth) {
    // print the value if it's not an object
    var print = function(key, val) {
        var value = document.createElement("div");
        value.setAttribute("class", "field");
        key = key.replaceAll("_", " ").capitalize();

        // parse a URL if its a string
        if (val.parseURL != undefined) {
            val = val.parseURL();
        }

        value.innerHTML = `<span class="key">${key}</span>&nbsp;<span class="value">${val}</span>`
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
        output.appendChild(renderOutput(key, value, depth));
    }

    // return the entire rendered value
    return output;
}

// returns JSON formatted in a pre object
function renderJSON(val) {
    var json = document.createElement("pre");
    json.innerText = JSON.stringify(val, null, "\t");
    return json;
}

// renders the recent queries
function renderQueries() {
    var queries = JSON.parse(localStorage.getItem('recent'));
    if (queries == null) {
        return
    }

    var recent = document.createElement("div");
    recent.id = "recent";
    recent.innerHTML = '<h4>Recent</h4>';
    var content = document.getElementById("content");
    content.appendChild(recent);

    Object.entries(queries).forEach(([key, value]) => {
        var a = document.createElement("a");
        a.href = value;
        a.innerText = key;
        a.setAttribute("class", "query");
        recent.appendChild(a);
    });
}

// render the response output
function renderResponse(response, rsp) {
    response.innerText = '';

    // set a title
    var h4 = document.createElement("h4");
    h4.setAttribute("class", "title");
    h4.innerText = "Response";

    // text based output
    var textOutput = renderOutput("response", rsp, 0);
    // json based output
    var jsonOutput = renderJSON(rsp);

    var output = document.createElement("div");
    output.id = "output";

    var links = document.createElement("span");
    links.id = "links";

    // create the output links
    var textBtn = document.createElement("a");
    var divider = document.createElement("span");
    var jsonBtn = document.createElement("a");
    divider.innerText = ' | ';
    textBtn.innerText = "Text";
    jsonBtn.innerText = "JSON";

    // execute on text
    textBtn.onclick = function(ev) {
        ev.preventDefault();
        output.innerText = '';
        output.appendChild(textOutput);
    };

    // execute on json
    jsonBtn.onclick = function(ev) {
        ev.preventDefault();
        output.innerText = '';
        output.appendChild(jsonOutput);
    };

    textBtn.href = '#text';
    jsonBtn.href = '#json';
    // create the links
    links.appendChild(textBtn);
    links.appendChild(divider);
    links.appendChild(jsonBtn);
    // add links to header
    h4.appendChild(links);

    if (window.location.hash == "#json") {
        output.appendChild(jsonOutput);
    } else {
        output.appendChild(textOutput);
    }
    // set the response output;
    response.appendChild(h4);
    response.appendChild(output);
}

function submitLogin(form) {
    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;
    login(username, password);
    return false;
}

function saveQuery(q) {
    var recent = JSON.parse(localStorage.getItem('recent'));
    if (recent == null) {
        recent = {};
    }
    var key = q + window.location.search.replace(/^\?/, " ");
    console.log("saving", key, window.location.href);
    recent[key] = window.location.href;
    localStorage.setItem('recent', JSON.stringify(recent));
}

function submitLogout(form) {
    logout()
}

function main() {
    // check cookie
    var token = getCookie(cookie);

    if (token == undefined || token == "") {
        return renderLogin();
    }

    // parse the url
    if (window.location.pathname == "/") {
        console.log("render index");
        return renderIndex();
    }

    if (window.location.pathname == "/services") {
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
