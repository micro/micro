package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	secret  = os.Getenv("GITHUB_WEBHOOK_SECRET")
	command = "./update.sh"
	lock    sync.RWMutex
	// the latest update
	update = new(Update)
)

type Update struct {
	Commit  string `json:"commit"`
	Image   string `json:"image"`
	Release string `json:"release"`
}

// get the latest commit
func getLatestCommit() (string, error) {
	rsp, err := http.Get("https://api.github.com/repos/micro/micro/commits")
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	// unmarshal commits
	var commits []map[string]interface{}
	err = json.Unmarshal(b, &commits)
	if err != nil {
		return "", err
	}
	// get the commits
	if len(commits) == 0 {
		return "", err
	}
	// the latest commit
	commit := commits[0]["sha"].(string)
	return commit, nil
}

func getLatestRelease() (string, error) {
	rsp, err := http.Get("https://api.github.com/repos/micro/micro/releases")
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	// unmarshal commits
	var releases []map[string]interface{}
	err = json.Unmarshal(b, &releases)
	if err != nil {
		return "", err
	}
	// get the commits
	if len(releases) == 0 {
		return "", err
	}
	// the latest commit
	release := releases[0]["tag_name"].(string)
	return release, nil
}

func getLatestImage() (string, error) {
	rsp, err := http.Get("https://hub.docker.com/v2/repositories/micro/micro/tags/latest")
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	// unmarshal commits
	var images map[string]interface{}
	err = json.Unmarshal(b, &images)
	if err != nil {
		return "", err
	}
	// get the commits
	updated := images["last_updated"].(string)
	return updated, nil
}

// set new update
func getUpdates() {
	commit, _ := getLatestCommit()
	release, _ := getLatestRelease()
	image, _ := getLatestImage()

	// update commit and release
	lock.Lock()
	defer lock.Unlock()

	update.Commit = commit
	update.Release = release
	update.Image = image
}

func main() {
	// get the latest updates
	getUpdates()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			lock.RLock()
			defer lock.RUnlock()

			b, _ := json.Marshal(update)
			w.Write(b)
			return
		}

		// if signature is blank assume its the docker webhook
		if len(r.Header.Get("X-Hub-Signature")) == 0 {
			// check if its the docker webhook
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return
			}
			var request map[string]interface{}
			if err := json.Unmarshal(b, &request); err != nil {
				return
			}
			// check a callback url exists
			v := request["callback_url"]
			if v == nil {
				return
			}
			callback, ok := v.(string)
			if ok && strings.HasPrefix(callback, "https://registry.hub.docker.com") {
				image, err := getLatestImage()
				if err != nil {
					return
				}
				lock.Lock()
				update.Image = image
				lock.Unlock()
			}
		}

		// assume github push
		parts := strings.Split(r.Header.Get("X-Hub-Signature"), "=")
		if len(parts) < 2 {
			log.Print("not enough parts in X-Hub-Signature")
			return
		}

		sha, _ := hex.DecodeString(parts[1])
		b, _ := ioutil.ReadAll(r.Body)
		mac := hmac.New(sha1.New, []byte(secret))
		mac.Write(b)
		expect := mac.Sum(nil)
		equals := hmac.Equal(sha, expect)

		if !equals {
			log.Print("hmac not equal expected")
			return
		}

		// update the latest values based on what type of event was received
		switch r.Header.Get("X-GitHub-Event") {
		case "push":
			commit, err := getLatestCommit()
			if err != nil {
				return
			}
			lock.Lock()
			update.Commit = commit
			lock.Unlock()
		case "release":
			release, err := getLatestRelease()
			if err != nil {
				return
			}
			lock.Lock()
			update.Release = release
			lock.Unlock()

			// no updates on release
			return
		default:
			log.Print("received unknown git event", r.Header.Get("X-GitHub-Event"))
			return
		}

		// exec the update
		go func() {
			lock.Lock()
			defer lock.Unlock()
			// run the command
			log.Print("update micro...error:", exec.Command(command).Run())
		}()
	})

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
