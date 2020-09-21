package runtime

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/micro/go-micro/v3/runtime/local/source/git"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/internal/config"
	"github.com/urfave/cli/v2"
)

// timeAgo returns the time passed
func timeAgo(v string) string {
	if len(v) == 0 {
		return "unknown"
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return v
	}

	return fmt.Sprintf("%v ago", fmtDuration(time.Since(t)))
}

func fmtDuration(d time.Duration) string {
	// round to secs
	d = d.Round(time.Second)

	var resStr string
	days := d / (time.Hour * 24)
	if days > 0 {
		d -= days * time.Hour * 24
		resStr = fmt.Sprintf("%dd", days)
	}
	h := d / time.Hour
	if len(resStr) > 0 || h > 0 {
		d -= h * time.Hour
		resStr = fmt.Sprintf("%s%dh", resStr, h)
	}
	m := d / time.Minute
	if len(resStr) > 0 || m > 0 {
		d -= m * time.Minute
		resStr = fmt.Sprintf("%s%dm", resStr, m)
	}
	s := d / time.Second
	resStr = fmt.Sprintf("%s%ds", resStr, s)
	return resStr
}

// exists returns whether the given file or directory exists
func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func sourceExists(source *git.Source) error {
	ref := source.Ref
	if ref == "" || ref == "latest" {
		ref = "master"
	}

	sourceExistsAt := func(url string, source *git.Source) error {
		req, _ := http.NewRequest("GET", url, nil)

		// add the git credentials if set
		if creds, ok := getGitCredentials(source.Repo); ok {
			req.Header.Set("Authorization", "token "+creds)
		}

		client := new(http.Client)
		resp, err := client.Do(req)

		// @todo gracefully degrade?
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return fmt.Errorf("service at %v@%v not found", source.Repo, ref)
		}
		return nil
	}

	if strings.Contains(source.Repo, "github") {
		// Github specific existence checs
		repo := strings.ReplaceAll(source.Repo, "github.com/", "")
		url := fmt.Sprintf("https://api.github.com/repos/%v/contents/%v?ref=%v", repo, source.Folder, ref)
		return sourceExistsAt(url, source)
	} else if strings.Contains(source.Repo, "gitlab") {
		// Gitlab specific existence checks

		// @todo better check for gitlab
		url := fmt.Sprintf("https://%v", source.Repo)
		return sourceExistsAt(url, source)
	}
	return nil
}

func appendSourceBase(ctx *cli.Context, workDir, source string) string {
	isLocal, _ := git.IsLocal(workDir, source)
	// @todo add list of supported hosts here or do this check better
	if !isLocal && !strings.Contains(source, ".com") && !strings.Contains(source, ".org") && !strings.Contains(source, ".net") {
		baseURL, _ := config.Get("git", util.GetEnv(ctx).Name, "baseurl")
		if len(baseURL) == 0 {
			baseURL, _ = config.Get("git", "baseurl")
		}
		if len(baseURL) == 0 {
			return path.Join("github.com/micro/services", source)
		}
		return path.Join(baseURL, source)
	}
	return source
}

func getGitCredentials(repo string) (string, bool) {
	repo = strings.Split(repo, "/")[0]

	for _, org := range GitOrgs {
		if !strings.Contains(repo, org) {
			continue
		}

		// check the creds for the org
		creds, err := config.Get("git", "credentials", org)
		if err == nil && len(creds) > 0 {
			return creds, true
		}
	}

	return "", false
}

// todo: remove this
func grepMain(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".go") {
			continue
		}
		file := filepath.Join(path, f.Name())
		b, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		if strings.Contains(string(b), "package main") {
			return nil
		}
	}
	return fmt.Errorf("Directory does not contain a main package")
}
