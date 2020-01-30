// Package token is for api token management
package token

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hako/branca"
	p "github.com/micro/micro/v2/internal/token/proto"
	"github.com/pborman/uuid"
)

type Token struct {
	*p.Token
}

var (
	// agent
	a = "me0"

	// api token
	t = os.Getenv("MICRO_TOKEN_KEY")

	// token api
	u = "https://micro.mu/token/"
)

func init() {
	if uri := os.Getenv("MICRO_TOKEN_API"); len(uri) > 0 {
		u = uri
	}
}

func (t *Token) Encode(key string) (string, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	if len(key) > 32 {
		key = key[:32]
	}
	br := branca.NewBranca(key)
	str, err := br.EncodeToString(string(b))
	if err != nil {
		return "", err
	}
	t.Key = str
	return str, nil
}

func (t *Token) Decode(key string, b []byte) error {
	if len(key) > 32 {
		key = key[:32]
	}
	br := branca.NewBranca(key)
	str, err := br.DecodeToString(string(b))
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(str), t); err != nil {
		return err
	}
	return nil
}

func (t *Token) Valid() error {
	// check id
	if len(t.Id) == 0 {
		return fmt.Errorf("token id invalid")
	}

	// check token expiry
	if u := time.Now().Unix(); (t.Expires - uint64(u)) < 0 {
		return fmt.Errorf("token expired")
	}

	// no claims
	if t.Claims == nil || len(t.Claims["email"]) == 0 {
		return fmt.Errorf("token claims invalid")
	}

	return nil
}

// SendPass sends a one time pass
func SendPass(email string) error {
	uri, err := url.Parse(u + "pass?email=" + email)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Micro-Agent", a)
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
	return nil
}

// Generate generates the token
func Generate(email, pass string) (string, error) {
	rsp, err := http.PostForm(u+"generate", url.Values{
		"email": {email},
		"pass":  {pass},
	})
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	if rsp.StatusCode != 200 {
		return "", fmt.Errorf(string(b))
	}
	var res map[string]interface{}
	if err := json.Unmarshal(b, &res); err != nil {
		return "", err
	}
	token, _ := res["token"].(string)
	return token, nil
}

// Revoke revokes a token
func Revoke(tk string) error {
	data := url.Values{
		"token": {tk},
	}
	req, err := http.NewRequest("POST", u+"revoke", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Micro-Token", t)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.StatusCode == 401 {
		return fmt.Errorf("Api error: %s (require MICRO_TOKEN_KEY)", strings.TrimSpace(string(b)))
	}
	if rsp.StatusCode != 200 {
		return fmt.Errorf("API error: %s", strings.TrimSpace(string(b)))
	}
	return nil
}

// List lists the tokens
func List() ([]*Token, error) {
	if len(t) == 0 {
		return nil, fmt.Errorf("Require MICRO_TOKEN_KEY")
	}
	req, err := http.NewRequest("GET", u+"list", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Micro-Token", t)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode == 401 {
		return nil, fmt.Errorf("Api error: %s (require MICRO_TOKEN_KEY)", strings.TrimSpace(string(b)))
	}
	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s", strings.TrimSpace(string(b)))
	}
	var list map[string][]*Token
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}
	return list["tokens"], nil
}

// SetToken sets the api token
func SetToken(tk string) {
	t = tk
}

// Verify a token is valid
func Verify(tk string) error {
	data := url.Values{
		"token": {tk},
	}
	req, err := http.NewRequest("POST", u+"verify", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Micro-Token", t)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.StatusCode == 401 {
		return fmt.Errorf(strings.TrimSpace(string(b)))
	}
	if rsp.StatusCode != 200 {
		return fmt.Errorf(strings.TrimSpace(string(b)))
	}
	return nil
}

// New returns a new token
func New() *Token {
	return &Token{&p.Token{
		Id:      uuid.NewUUID().String(),
		Expires: uint64(time.Now().Add(time.Hour * 24 * 7).Unix()),
		Claims:  make(map[string]string),
	}}
}
