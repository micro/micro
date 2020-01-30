package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/handlers"
	pb "github.com/micro/micro/v2/cmd/usage/proto"
)

var (
	db *bolt.DB
	fd = "usage.db"

	mtx  sync.RWMutex
	seen = map[string]uint64{}
)

func setup() {
	// setup db
	d, err := bolt.Open(fd, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	db = d

	if err := db.Update(func(tx *bolt.Tx) error {
		for _, b := range []string{"usage", "metrics"} {
			if _, err := tx.CreateBucketIfNotExists([]byte(b)); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	go flush()
}

func flush() {
	for {
		time.Sleep(time.Hour)
		now := time.Now().UnixNano()
		mtx.Lock()
		for k, v := range seen {
			d := uint64(now) - v
			// 48 hours
			if d > 1.728e14 {
				delete(seen, k)
			}
		}
		seen = make(map[string]uint64)
		mtx.Unlock()
	}
}

func process(w http.ResponseWriter, r *http.Request, u *pb.Usage) {
	today := time.Now().Format("20060102")
	key := fmt.Sprintf("%s-%s", u.Service, u.Id)
	now := uint64(time.Now().UnixNano())

	mtx.Lock()
	last := seen[key]
	lastSeen := now - last
	seen[key] = now
	mtx.Unlock()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`usage`))
		buf, err := proto.Marshal(u)
		if err != nil {
			return err
		}
		k := fmt.Sprintf("%d-%s", u.Timestamp, key)
		// save this usage
		if err := b.Put([]byte(k), buf); err != nil {
			return err
		}

		// save daily usage
		b = tx.Bucket([]byte(`metrics`))
		dailyKey := fmt.Sprintf("%s-%s", today, u.Service)

		// get usage
		v := b.Get([]byte(dailyKey))
		if v == nil {
			// todo: don't overwrite this
			u.Metrics.Count["services"] = uint64(1)
			m, _ := proto.Marshal(u.Metrics)
			return b.Put([]byte(dailyKey), m)
		}

		m := new(pb.Metrics)
		if err := proto.Unmarshal(v, m); err != nil {
			return err
		}

		// update request count
		m.Count["requests"] += u.Metrics.Count["requests"]
		m.Count["services"] += u.Metrics.Count["services"]

		// not seen today add it
		if lastSeen == 0 || lastSeen > 7.2e13 {
			c := m.Count["instances"]
			c++
			m.Count["instances"] = c
		}

		buf, err = proto.Marshal(m)
		if err != nil {
			return err
		}

		// store today-micro.api/new/cli/proxy
		return b.Put([]byte(dailyKey), buf)
	})
}

func metrics(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	prefix := time.Now().Add(time.Hour * -24).Format("20060102")
	metrics := map[string]interface{}{}

	if date := r.Form.Get("date"); len(date) >= 4 && len(date) <= 8 {
		prefix = date
	}

	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(`metrics`)).Cursor()

		for k, v := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, []byte(prefix)); k, v = c.Next() {
			m := new(pb.Metrics)
			proto.Unmarshal(v, m)
			key := strings.TrimPrefix(string(k), prefix+"-")
			metrics[key] = m
		}
		return nil
	})

	var buf []byte
	ct := r.Header.Get("Content-Type")

	if v := r.Form.Get("pretty"); len(v) > 0 || ct != "application/json" {
		buf, _ = json.MarshalIndent(metrics, "", "\t")
	} else {
		buf, _ = json.Marshal(metrics)
	}

	if len(buf) == 0 {
		buf = []byte(`{}`)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// return metrics
	if r.Method == "GET" {
		metrics(w, r)
		return
	}

	// require post for updates
	if r.Method != "POST" {
		return
	}
	if r.Header.Get("Content-Type") != "application/protobuf" {
		return
	}

	if r.UserAgent() != "micro/usage" {
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	u := new(pb.Usage)
	if err := proto.Unmarshal(b, u); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	go process(w, r, u)
}

func main() {
	setup()
	http.HandleFunc("/", handler)

	lh := handlers.LoggingHandler(os.Stdout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/usage") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/usage")
		}
		http.DefaultServeMux.ServeHTTP(w, r)
	}))

	if err := http.ListenAndServe(":8091", lh); err != nil {
		log.Fatal(err)
	}
}
