package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	elastic "gopkg.in/olivere/elastic.v5"
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
	e *elastic.Client
	t = os.Getenv("GITHUB_ACCESS_TOKEN")
	f = []func(context.Context, *github.Client) ([]*github.Repository, error){
		getMicro,
		get("go-micro"),
		get("go-plugins"),
		get("micro"),
	}
)

func init() {
	el, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}
	_, err = el.CreateIndex("micro").Do(context.Background())
	if err != nil && !strings.Contains(err.Error(), "type=index_already_exists_exception") {
		panic(err)
	}
	e = el
}

func update(ctx context.Context, c *github.Client) map[string]bool {
	from := 0
	size := 100
	updated := make(map[string]bool)

	for {
		ctx := context.Background()
		// process all the things
		s, err := e.Search("micro").Type("repo").From(from).Size(size).Sort("updated", false).Do(ctx)
		if err != nil {
			break
		}

		for _, h := range s.Hits.Hits {
			var r repo
			if err := json.Unmarshal(*h.Source, &r); err != nil {
				continue
			}

			rsp, err := http.Get(r.URL)
			if err != nil {
				continue
			}
			io.Copy(ioutil.Discard, rsp.Body)
			rsp.Body.Close()

			// delete on 404
			if rsp.StatusCode == 404 {
				fmt.Printf("deleting %s\n", r.Name)
				e.Delete().Index("micro").Type("repo").Refresh("true").Id(r.Name).Do(ctx)
				updated[r.Name] = true
				continue
			}

			// update
			parts := strings.Split(r.Name, "/")
			if len(parts) < 2 {
				continue
			}

			// get and update
			fmt.Printf("updating %s\n", r.Name)
			rep, _, err := c.Repositories.Get(ctx, parts[0], parts[1])
			if err != nil {
				if r, ok := err.(*github.RateLimitError); ok {
					time.Sleep(r.Rate.Reset.Sub(time.Now()))
				}
				continue
			}

			store([]*github.Repository{rep})
			updated[r.Name] = true
			time.Sleep(time.Second * 3)
		}

		from += size
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

			if page == rsp.LastPage || rsp.NextPage == 0 || rsp.NextPage >= 5 {
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
func run(onStart bool) {
	if !onStart {
		time.Sleep(time.Hour * 3)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

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
			time.Sleep(time.Minute)
		}

		time.Sleep(time.Hour)
	}
}

// store the things
func store(repos []*github.Repository) {
	ctx := context.Background()

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
		_, err := e.Index().Index("micro").Type("repo").Id(re.Name).BodyJson(re).Refresh("true").Do(ctx)
		if err != nil {
			fmt.Printf("failed to store %s\n", re.Name)
			continue
		}

		fmt.Printf("stored %s\n", re.Name)
	}
}

func browseHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	page, from := getPageOffset(r.Form)

	sort := r.Form.Get("s")

	switch sort {
	case "stars":
		sort = "stars"
	default:
		sort = "updated"
	}

	ctx := context.Background()
	s, err := e.Search("micro").Type("repo").From(from).Size(size).Sort(sort, false).Do(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var re []interface{}

	for _, h := range s.Hits.Hits {
		re = append(re, h.Source)
	}

	b, err := json.Marshal(map[string]interface{}{
		"repos": re,
		"page":  page,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func recentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	s, err := e.Search("micro").Type("repo").Size(size).Sort("updated", false).Do(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var re []interface{}

	for _, h := range s.Hits.Hits {
		re = append(re, h.Source)
	}

	b, err := json.Marshal(map[string]interface{}{
		"repos": re,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := r.Form.Get("q")

	if len(q) == 0 || len(q) > 256 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
		return
	}

	page, from := getPageOffset(r.Form)
	qe := elastic.NewQueryStringQuery(q).Escape(true)

	ctx := context.Background()
	s, err := e.Search("micro").Type("repo").From(from).Size(size).Sort("updated", false).Query(qe).Do(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var re []interface{}

	for _, h := range s.Hits.Hits {
		re = append(re, h.Source)
	}

	b, err := json.Marshal(map[string]interface{}{
		"repos": re,
		"page":  page,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	go run(true)

	// api
	http.HandleFunc("/explore/api/browse", browseHandler)
	http.HandleFunc("/explore/api/recent", recentHandler)
	http.HandleFunc("/explore/api/search", searchHandler)

	http.Handle("/explore/", http.StripPrefix("/explore/", http.FileServer(http.Dir("html/explore"))))

	http.ListenAndServe(":8089", nil)
}
