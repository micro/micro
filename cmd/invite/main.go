package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var (
	inviteURL = os.Getenv("SLACK_INVITE_URL")
	secret    = os.Getenv("GOOGLE_RECAPTCHA_SECRET")
)

func invite(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	c := r.Form.Get("g-recaptcha-response")

	rsp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
		"secret":   {secret},
		"response": {c},
	})
	if err != nil {
		return
	}
	if rsp.StatusCode != 200 {
		return
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return
	}
	success, ok := resp["success"].(bool)
	if !ok {
		return
	}

	if !success {
		return
	}

	// success
	http.Redirect(w, r, inviteURL, 302)
}

func main() {
	http.HandleFunc("/join", invite)
	http.ListenAndServe(":8090", nil)
}
