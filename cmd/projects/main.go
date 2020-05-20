package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	limit = 20
	size  = 20
)

type repo struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Desc    string   `json:"desc"`
	Lang    string   `json:"lang"`
	Forks   int      `json:"forks"`
	Stars   int      `json:"stars"`
	Tags    []string `json:"tags"`
	Type    string   `json:"type"`
	Updated int64    `json:"updated"`
}

var (
	lock sync.RWMutex
	// projects index
	index = make(map[string]*repo)
	// latest updates
	latest = []*repo{}
	// top projects
	top = []*repo{}

	token = ""

	f = []func(context.Context, *github.Client) ([]*github.Repository, error){
		get("go-micro"),
		get("go-plugins"),
		get("micro"),
		getMicro,
	}
)

func init() {
	if len(token) == 0 {
		token = os.Getenv("GITHUB_ACCESS_TOKEN")
	}

	// load index
	load()
	// index latest
	indexLatest()
	// index top
	indexTop()
}

func indexLatest() {
	lock.RLock()
	defer lock.RUnlock()

	var latestIndex []*repo

	for _, v := range index {
		latestIndex = append(latestIndex, v)
	}

	// sort by updated
	sort.Slice(latestIndex, func(i, j int) bool { return latestIndex[i].Updated > latestIndex[j].Updated })

	latest = latestIndex
}

func indexTop() {
	lock.RLock()
	defer lock.RUnlock()

	var topIndex []*repo

	for _, v := range index {
		topIndex = append(topIndex, v)
	}

	// sort by stars
	sort.Slice(topIndex, func(i, j int) bool { return topIndex[i].Stars > topIndex[j].Stars })

	top = topIndex
}

func load() {
	lock.Lock()
	defer lock.Unlock()

	// load the index
	b, err := ioutil.ReadFile("projects.idx")
	if err != nil {
		fmt.Println("error reading index", err)
		return
	}

	if err := json.Unmarshal(b, &index); err != nil {
		fmt.Println("error unmarshaling index", err)
	}
}

func save() {
	lock.RLock()
	defer lock.RUnlock()

	b, err := json.Marshal(index)
	if err != nil {
		fmt.Println("error marshaling index", err)
		return
	}

	if err := ioutil.WriteFile("projects.idx", b, 0666); err != nil {
		fmt.Println("error writing index", err)
		return
	}
}

func update(ctx context.Context, c *github.Client) map[string]bool {
	updated := make(map[string]bool)
	indexCopy := make(map[string]string)

	lock.RLock()

	for _, v := range index {
		indexCopy[v.Name] = v.URL
	}

	lock.RUnlock()

	// process all the things
	for name, url := range indexCopy {
		fmt.Println("Checking updates", name, url)
		rsp, err := http.Get(url)
		if err != nil {
			continue
		}
		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()

		// delete on 404
		if rsp.StatusCode == 404 {
			fmt.Printf("deleting %s\n", name)
			lock.Lock()
			delete(index, name)
			lock.Unlock()
			updated[name] = true
			continue
		}

		// update
		parts := strings.Split(name, "/")
		if len(parts) < 2 {
			continue
		}

		// get and update
		fmt.Printf("updating %s\n", name)
		rep, _, err := c.Repositories.Get(ctx, parts[0], parts[1])
		if err != nil {
			if r, ok := err.(*github.RateLimitError); ok {
				time.Sleep(r.Rate.Reset.Sub(time.Now()))
			}
			continue
		}

		store([]*github.Repository{rep})
		updated[name] = true
		time.Sleep(time.Second * 3)
	}

	return updated
}

func get(term string) func(ctx context.Context, c *github.Client) ([]*github.Repository, error) {
	results := make(map[string]*github.Repository)

	return func(ctx context.Context, c *github.Client) ([]*github.Repository, error) {
		page := 1

		for {
			so := &github.SearchOptions{
				Sort: "indexed",
			}
			so.ListOptions = github.ListOptions{Page: page}

			search := fmt.Sprintf(`"github.com/micro/%s"`, term)
			fmt.Println("Searching for", search, page)
			res, rsp, err := c.Search.Code(ctx, search, so)
			if err != nil {
				fmt.Printf("[%s] error searching %v\n", term, err)
				break
			}

			for _, r := range res.CodeResults {
				// only process go files
				if !strings.HasSuffix(*r.Path, ".go") {
					continue
				}
				// skip tests
				if strings.HasSuffix(*r.Path, "_test.go") {
					continue
				}
				results[*r.Repository.FullName] = r.Repository
			}

			if page == rsp.LastPage || rsp.NextPage == 0 || rsp.NextPage >= 100 {
				break
			}

			page = rsp.NextPage
			time.Sleep(time.Second * 5)
		}

		var repos []*github.Repository

		for _, v := range results {
			if v.PushedAt != nil {
				repos = append(repos, v)
				continue
			}

			fmt.Println("Retrieving repo", *v.Owner.Login, *v.Name)
			rep, _, err := c.Repositories.Get(ctx, *v.Owner.Login, *v.Name)
			if err != nil {
				if r, ok := err.(*github.RateLimitError); ok {
					time.Sleep(r.Rate.Reset.Sub(time.Now()))
				}
				continue
			}
			repos = append(repos, rep)
			time.Sleep(time.Second * 10)
		}

		return repos, nil
	}
}

func getMicro(ctx context.Context, c *github.Client) ([]*github.Repository, error) {
	page := 1
	var repos []*github.Repository

	for {
		ro := &github.RepositoryListOptions{}
		ro.ListOptions = github.ListOptions{Page: page}
		reps, rsp, err := c.Repositories.List(ctx, "micro", ro)
		if err != nil {
			break
		}

		repos = append(repos, reps...)

		if page == rsp.LastPage || rsp.NextPage == 0 {
			break
		}

		page = rsp.NextPage
		time.Sleep(time.Second)
	}

	return repos, nil
}

func getPageOffset(vars url.Values) (int, int) {
	page, err := strconv.Atoi(vars.Get("p"))
	if err != nil {
		page = 1
	}

	if page > limit {
		page = limit
	}

	next := page - 1
	if page == 1 {
		next = 0
	}

	offset := next * size
	return page, offset
}

// get all the things
func run(token string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	// create http and github client
	ctx := context.Background()
	cli := oauth2.NewClient(ctx, ts)
	client := github.NewClient(cli)

	fmt.Println("Starting runner")

	for {
		// clean up or refresh first
		updated := update(ctx, client)

		for _, v := range f {
			repos, err := v(ctx, client)
			if err != nil {
				fmt.Println(err)
				continue
			}

			var upRepos []*github.Repository

			for _, r := range repos {
				// only process what has not been updated
				if !updated[*r.Name] {
					upRepos = append(upRepos, r)
				}
			}

			fmt.Printf("storing %d\n", len(upRepos))
			store(upRepos)

			// save and index
			indexLatest()
			indexTop()
			save()

			time.Sleep(time.Minute)
		}

		time.Sleep(time.Hour)
	}
}

// store the things
func store(repos []*github.Repository) {
	for _, r := range repos {
		if *r.Private {
			continue
		}

		re := &repo{
			Name:    r.GetFullName(),
			Desc:    r.GetDescription(),
			Lang:    r.GetLanguage(),
			Stars:   r.GetStargazersCount(),
			Forks:   r.GetForksCount(),
			URL:     r.GetHTMLURL(),
			Updated: r.GetPushedAt().Unix(),
		}

		switch {
		case strings.HasPrefix(*r.Name, "go-"):
			re.Type = "lib"
		case strings.HasSuffix(re.Name, "-api"):
			re.Type = "api"
		case strings.HasSuffix(re.Name, "-bot"):
			re.Type = "bot"
		case strings.HasSuffix(re.Name, "-srv"):
			re.Type = "srv"
		case strings.HasSuffix(re.Name, "-web"):
			re.Type = "web"
		}

		lock.Lock()
		index[re.Name] = re
		lock.Unlock()

		fmt.Printf("stored %s\n", re.Name)
	}
}

func getRepos(repos []*repo, from, to int) []*repo {
	if len(repos) < from {
		return nil
	}

	check := repos[from:]

	if len(check) <= to {
		return check
	}

	var cp []*repo

	for i := 0; i < to; i++ {
		cp = append(cp, check[i])
	}

	return cp
}

func browseHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	page, from := getPageOffset(r.Form)
	sort := r.Form.Get("s")

	lock.RLock()

	var repos []*repo

	switch sort {
	case "stars":
		repos = getRepos(top, from, size)
	default:
		repos = getRepos(latest, from, size)
	}

	b, err := json.Marshal(map[string]interface{}{
		"repos": repos,
		"page":  page,
	})

	lock.RUnlock()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func recentHandler(w http.ResponseWriter, r *http.Request) {
	lock.RLock()

	repos := latest
	if len(repos) > 20 {
		repos = repos[:20]
	}

	b, err := json.Marshal(map[string]interface{}{
		"repos": repos,
	})

	lock.RUnlock()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func reposHandler(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	count := fmt.Sprintf("%d", len(index))
	lock.RUnlock()
	w.Write([]byte(count))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := r.Form.Get("q")
	page, from := getPageOffset(r.Form)

	if len(q) == 0 || len(q) > 256 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
		return
	}

	lock.RLock()

	repos := []*repo{}
	tokens := strings.Split(q, " ")

	for _, repo := range latest {
		for _, token := range tokens {
			if strings.Contains(repo.Name, token) {
				repos = append(repos, repo)
				break
			}

			if strings.Contains(repo.Desc, token) {
				repos = append(repos, repo)
				break
			}
		}
	}

	// strip it down
	repos = getRepos(repos, from, size)

	b, err := json.Marshal(map[string]interface{}{
		"repos": repos,
		"page":  page,
	})

	lock.RUnlock()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	// api
	http.HandleFunc("/api/browse", browseHandler)
	http.HandleFunc("/api/recent", recentHandler)
	http.HandleFunc("/api/search", searchHandler)
	http.HandleFunc("/api/repos", reposHandler)

	go run(token)

	http.ListenAndServe(":7090", nil)
}
