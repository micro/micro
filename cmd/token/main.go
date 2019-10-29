package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/handlers"
	"github.com/micro/micro/internal/token"
	"github.com/patrickmn/go-cache"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gopkg.in/gomail.v2"
)

var (
	// otp cache
	c = cache.New(5*time.Minute, 10*time.Minute)
	// db
	db *bolt.DB

	// domain for the sender of the otp
	domain = os.Getenv("EMAIL_DOMAIN")
	// boltdb
	dbase = "token.db"
	// storage encryption key used for state
	dkey = os.Getenv("STORE_KEY")
	// outbound encryption key used for tokens
	key = os.Getenv("TOKEN_KEY")

	// Only used where validUser code is uncommented
	// admin token used to call the account service
	aToken = os.Getenv("ADMIN_TOKEN")
	// basic auth for account service
	user = os.Getenv("USERNAME")
	pass = os.Getenv("PASSWORD")
)

func setup() {
	// setup db
	d, err := bolt.Open(dbase, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	db = d

	// create buckets
	if err := db.Update(func(tx *bolt.Tx) error {
		for _, b := range []string{"token"} {
			if _, err := tx.CreateBucketIfNotExists([]byte(b)); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func sendEmail(key, email string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "Micro <support@"+domain+">")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your One Time Pass (OTP)")
	m.SetBody("text/html", "Hey,<p>Your one time pass is <b>"+key+"</b><p>The pass is valid for 5 minutes")

	d := gomail.NewDialer("smtp-relay.gmail.com", 587, "", "")
	d.LocalName = domain

	log.Println("sending email to:", email)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

func saveToken(email string, t *token.Token) error {
	k := []byte(fmt.Sprintf("%s-%s", email, t.Id))
	v, err := t.Encode(dkey)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(`token`)).Put(k, []byte(v))
	})
}

func delToken(email string, t *token.Token) error {
	k := []byte(fmt.Sprintf("%s-%s", email, t.Id))
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(`token`)).Delete(k)
	})
}

func getToken(email string, t *token.Token) (map[string]*token.Token, error) {
	tokens := make(map[string]*token.Token)

	err := db.View(func(tx *bolt.Tx) error {
		// get token
		if t != nil {
			key := fmt.Sprintf("%s-%s", email, t.Id)
			v := tx.Bucket([]byte(`token`)).Get([]byte(key))
			tk := new(token.Token)
			if err := tk.Decode(dkey, v); err != nil {
				return err
			}
			tokens[tk.Id] = tk
			return nil
		}

		// list tokens
		c := tx.Bucket([]byte(`token`)).Cursor()
		prefix := []byte(email)

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			t := new(token.Token)
			if err := t.Decode(dkey, v); err != nil {
				return err
			}
			if t.Claims["email"] != email {
				continue
			}
			tokens[t.Id] = t
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func requestToken(r *http.Request) (*token.Token, error) {
	tk := r.Header.Get("X-Micro-Token")
	tk = strings.TrimSpace(tk)
	if len(tk) == 0 {
		return nil, fmt.Errorf("invalid token")
	}
	// decode token
	t := new(token.Token)
	if err := t.Decode(key, []byte(tk)); err != nil {
		return nil, err
	}

	if err := t.Valid(); err != nil {
		return nil, err
	}

	return t, nil
}

func validateOTP(email, pass string) error {
	// the otp to verify we can do this
	if len(pass) == 0 {
		return fmt.Errorf("invalid pass")
	}

	// get otp from the cache
	cmail, ok := c.Get("otp:" + pass)
	if !ok {
		return fmt.Errorf("pass expired")
	}

	// check email of otp and request email match
	if email != cmail {
		return fmt.Errorf("invalid email")
	}

	// validate the otp
	ok, err := totp.ValidateCustom(pass, key, time.Now(), totp.ValidateOpts{
		Period:    300,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if !ok || err != nil {
		return fmt.Errorf("invalid otp")
	}

	// delete the otp
	c.Delete(pass)
	return nil
}

func validUser(w http.ResponseWriter, r *http.Request) error {
	email := r.Form.Get("email")
	if len(email) == 0 {
		return fmt.Errorf("email is blank")
	}

/*
	TODO: replace validation of the user 
	uri, err := url.Parse("http://localhost:9091/verify/email?email=" + email)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", uri.String(), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(user, pass)
	req.Header.Set("X-Micro-Token", aToken)

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.StatusCode != 200 {
		return fmt.Errorf(strings.TrimSpace(string(b)))
	}
*/

	return nil
}

func main() {
	setup()

	// requires email
	http.HandleFunc("/pass", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		// TODO: check if the user is valid
		// Previously verified against chargebee
		if err := validUser(w, r); err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		email := r.Form.Get("email")

		// gen otp
		k, err := totp.GenerateCodeCustom(key, time.Now(), totp.ValidateOpts{
			Period:    300,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
		if err != nil {
			http.Error(w, "failed to generate code "+err.Error(), 500)
			return
		}

		// email otp
		if err := sendEmail(k, email); err != nil {
			http.Error(w, "failed to send email: "+err.Error(), 500)
			return
		}

		// save otp:email
		c.Set("otp:"+k, email, cache.DefaultExpiration)
	})

	// get token
	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		// email to generate token for
		email := r.Form.Get("email")
		if len(email) == 0 {
			http.Error(w, "invalid email", 500)
			return
		}

		// check if we have a pass
		pass := r.Form.Get("pass")
		if err := validateOTP(email, pass); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// passed validation

		// gen token
		t := token.New()
		t.Claims["email"] = email
		str, err := t.Encode(key)
		if err != nil {
			http.Error(w, "failed to encode token", 500)
			return
		}

		// make request
		b, err := json.Marshal(map[string]interface{}{
			"token": str,
		})
		if err != nil {
			http.Error(w, "failed to encode response", 500)
			return
		}

		// save token
		saveToken(email, t)

		// send response
		w.Write(b)
	})

	// check token is valid
	// open to all
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		var t *token.Token

		// verify header token or token field?
		tv := r.Form.Get("token")

		// got token field
		if len(tv) > 0 {
			tr := new(token.Token)

			// decode token
			if err := tr.Decode(key, []byte(tv)); err != nil {
				http.Error(w, "token is invalid", 500)
				return
			}

			t = tr
		} else {
			// get header token
			tr, err := requestToken(r)
			if err != nil {
				http.Error(w, "invalid token", 401)
				return
			}
			t = tr
		}

		// is the token valid?
		if err := t.Valid(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// check db
		tk, err := getToken(t.Claims["email"], t)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// got token
		if _, ok := tk[t.Id]; !ok {
			http.Error(w, "invalid token", 401)
			return
		}

		b, err := json.Marshal(t)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(b)
	})

	// get list of tokens
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		// get header token
		t, err := requestToken(r)
		if err != nil {
			http.Error(w, "invalid token", 401)
			return
		}

		// list of tokens
		mtokens, err := getToken(t.Claims["email"], nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// check the X-Micro-Token hasn't been revoked
		if _, ok := mtokens[t.Id]; !ok {
			http.Error(w, "invalid token", 401)
			return
		}

		var tokens []*token.Token
		for _, v := range mtokens {
			tokens = append(tokens, v)
		}

		// encode the response
		b, err := json.Marshal(map[string]interface{}{
			"tokens": tokens,
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Write(b)
	})

	http.HandleFunc("/revoke", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		t, err := requestToken(r)
		if err != nil {
			http.Error(w, "invalid token", 401)
			return
		}

		// get token
		mtokens, err := getToken(t.Claims["email"], nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// check the X-Micro-Token hasn't been revoked
		if _, ok := mtokens[t.Id]; !ok {
			http.Error(w, "invalid token", 401)
			return
		}

		// get the token to revoke
		trevoke := r.Form.Get("token")
		id := r.Form.Get("id")

		if len(trevoke) == 0 && len(id) == 0 {
			http.Error(w, "token/id is blank", 500)
			return
		}

		var tr *token.Token

		if len(id) > 0 {
			t, ok := mtokens[id]
			if !ok {
				return
			}
			tr = t
		} else {
			tr = new(token.Token)

			// decode token
			if err := tr.Decode(key, []byte(trevoke)); err != nil {
				http.Error(w, "token is invalid", 500)
				return
			}

			if err := tr.Valid(); err != nil {
				http.Error(w, err.Error(), 401)
				return
			}
		}

		// not the same
		if t.Claims["email"] != tr.Claims["email"] {
			http.Error(w, "not allowed to revoke token", 401)
			return
		}

		// revoke the token
		if err := delToken(tr.Claims["email"], tr); err != nil {
			http.Error(w, err.Error(), 401)
		}
	})

	lh := handlers.LoggingHandler(os.Stdout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/token") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/token")
		}
		http.DefaultServeMux.ServeHTTP(w, r)
	}))

	hd := handlers.ProxyHeaders(lh)

	if err := http.ListenAndServe(":10001", hd); err != nil {
		log.Fatal(err)
	}
}
