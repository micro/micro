package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Result struct {
	Count map[string]int64 `json:"count"`
}

func printKey(k string, v int64) {
	if v == 0 {
		return
	}

	var u string
	var c float64

	switch {
	case v > 1e9:
		c = float64(v) / 1e9
		u = "b"
	case v > 1e6:
		c = float64(v) / 1e6
		u = "m"
	case v > 1e4:
		c = float64(v) / 1e3
		u = "k"
	default:
		c = float64(v)
	}

	fmt.Printf("micro %s:\t%.2f%s\n", k, c, u)
}

func years() [][]byte {
	var years [][]byte

	for _, year := range []string{"2019", "2020"} {
		rsp, err := http.Get("https://micro.mu/usage?date=" + year)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer rsp.Body.Close()

		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		years = append(years, b)
	}

	return years
}

func main() {
	var cKey string
	if len(os.Args) > 1 {
		cKey = os.Args[1]
	}

	var results map[string]Result

	counts := map[string]int64{}
	highest := map[string]int64{}
	//	daily := map[string]int64{}
	monthly := map[string]int64{}

	// get all the results
	for _, year := range years() {
		var res map[string]Result

		if err := json.Unmarshal(year, &results); err != nil {
			fmt.Println(err)
			return
		}

		for k, v := range res {
			results[k] = v
		}
	}

	for k, v := range results {
		// 20190520-micro.new
		parts := strings.Split(k, ".")
		if len(parts) < 2 {
			continue
		}

		// micro.new
		key := parts[len(parts)-1]

		if len(cKey) > 0 && key != cKey {
			continue
		}

		// counts[micro.new] += requests
		c := counts[key]
		c += v.Count["requests"]
		// save
		counts[key] = c

		// set highest
		if i := highest[key]; v.Count["requests"] > i {
			highest[key] = v.Count["requests"]
		}

		// set monthly
		month := parts[0][:6]
		mkey := key + " (" + month + ")"
		c = monthly[mkey]
		c += v.Count["requests"]
		monthly[mkey] = c
	}

	fmt.Println("Total requests:")

	for k, v := range counts {
		printKey(k, v)
	}

	fmt.Println("\nHighest requests:")

	for k, v := range highest {
		printKey(k, v)
	}

	fmt.Println("\nMonthly requests:")
	var keys []string
	for k, _ := range monthly {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		printKey(k, monthly[k])
	}
}
