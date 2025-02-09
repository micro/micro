// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/util/qson/qson_test.go

// Package qson implmenets decoding of URL query params
// into JSON and Go values (using JSON struct tags).
//
// See https://golang.org/pkg/encoding/json/ for more
// details on JSON struct tags.
package qson

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var (
	// ErrInvalidParam is returned when invalid data is provided to the ToJSON or Unmarshal function.
	// Specifically, this will be returned when there is no equals sign present in the URL query parameter.
	ErrInvalidParam error = errors.New("qson: invalid url query param provided")

	bracketSplitter *regexp.Regexp
)

func init() {
	bracketSplitter = regexp.MustCompile("\\[|\\]")
}

func btSplitter(str string) []string {
	r := bracketSplitter.Split(str, -1)
	for idx, s := range r {
		if len(s) == 0 {
			if len(r) > idx+1 {
				copy(r[idx:], r[idx+1:])
				r = r[:len(r)-1]
			}
		}
	}
	return r
}

// Unmarshal will take a dest along with URL
// query params and attempt to first turn the query params
// into JSON and then unmarshal those into the dest variable
//
// BUG(joncalhoun): If a URL query param value is something
// like 123 but is expected to be parsed into a string this
// will currently result in an error because the JSON
// transformation will assume this is intended to be an int.
// This should only affect the Unmarshal function and
// could likely be fixed, but someone will need to submit a
// PR if they want that fixed.
func Unmarshal(dst interface{}, query string) error {
	b, err := ToJSON(query)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

// ToJSON will turn a query string like:
//
//	cat=1&bar%5Bone%5D%5Btwo%5D=2&bar[one][red]=112
//
// Into a JSON object with all the data merged as nicely as
// possible. Eg the example above would output:
//
//	{"bar":{"one":{"two":2,"red":112}}}
func ToJSON(query string) ([]byte, error) {
	var (
		builder interface{} = make(map[string]interface{})
	)
	params := strings.Split(query, "&")
	for _, part := range params {
		tempMap, err := queryToMap(part)
		if err != nil {
			return nil, err
		}
		builder = merge(builder, tempMap)
	}
	return json.Marshal(builder)
}

// queryToMap turns something like a[b][c]=4 into
//
//	  map[string]interface{}{
//	    "a": map[string]interface{}{
//			  "b": map[string]interface{}{
//				  "c": 4,
//			  },
//		  },
//	  }
func queryToMap(param string) (map[string]interface{}, error) {
	rawKey, rawValue, err := splitKeyAndValue(param)
	if err != nil {
		return nil, err
	}
	rawValue, err = url.QueryUnescape(rawValue)
	if err != nil {
		return nil, err
	}
	rawKey, err = url.QueryUnescape(rawKey)
	if err != nil {
		return nil, err
	}

	pieces := btSplitter(rawKey)
	key := pieces[0]

	// If len==1 then rawKey has no [] chars and we can just
	// decode this as key=value into {key: value}
	if len(pieces) == 1 {
		var value interface{}
		// First we try parsing it as an int, bool, null, etc
		err = json.Unmarshal([]byte(rawValue), &value)
		if err != nil {
			// If we got an error we try wrapping the value in
			// quotes and processing it as a string
			err = json.Unmarshal([]byte("\""+rawValue+"\""), &value)
			if err != nil {
				// If we can't decode as a string we return the err
				return nil, err
			}
		}
		return map[string]interface{}{
			key: value,
		}, nil
	}

	// If len > 1 then we have something like a[b][c]=2
	// so we need to turn this into {"a": {"b": {"c": 2}}}
	// To do this we break our key into two pieces:
	//   a and b[c]
	// and then we set {"a": queryToMap("b[c]", value)}
	ret := make(map[string]interface{}, 0)
	ret[key], err = queryToMap(buildNewKey(rawKey) + "=" + rawValue)
	if err != nil {
		return nil, err
	}

	// When URL params have a set of empty brackets (eg a[]=1)
	// it is assumed to be an array. This will get us the
	// correct value for the array item and return it as an
	// []interface{} so that it can be merged properly.
	if pieces[1] == "" {
		temp := ret[key].(map[string]interface{})
		ret[key] = []interface{}{temp[""]}
	}
	return ret, nil
}

// buildNewKey will take something like:
// origKey = "bar[one][two]"
// pieces = [bar one two ]
// and return "one[two]"
func buildNewKey(origKey string) string {
	pieces := btSplitter(origKey)

	ret := origKey[len(pieces[0])+1:]
	ret = ret[:len(pieces[1])] + ret[len(pieces[1])+1:]
	return ret
}

// splitKeyAndValue splits a URL param at the last equal
// sign and returns the two strings. If no equal sign is
// found, the ErrInvalidParam error is returned.
func splitKeyAndValue(param string) (string, string, error) {
	li := strings.LastIndex(param, "=")
	if li == -1 {
		return "", "", ErrInvalidParam
	}
	return param[:li], param[li+1:], nil
}

// merge merges a with b if they are either both slices
// or map[string]interface{} types. Otherwise it returns b.
func merge(a interface{}, b interface{}) interface{} {
	switch aT := a.(type) {
	case map[string]interface{}:
		return mergeMap(aT, b.(map[string]interface{}))
	case []interface{}:
		return mergeSlice(aT, b.([]interface{}))
	default:
		return b
	}
}

// mergeMap merges a with b, attempting to merge any nested
// values in nested maps but eventually overwriting anything
// in a that can't be merged with whatever is in b.
func mergeMap(a map[string]interface{}, b map[string]interface{}) map[string]interface{} {
	for bK, bV := range b {
		if _, ok := a[bK]; ok {
			a[bK] = merge(a[bK], bV)
		} else {
			a[bK] = bV
		}
	}
	return a
}

// mergeSlice merges a with b and returns the result.
func mergeSlice(a []interface{}, b []interface{}) []interface{} {
	a = append(a, b...)
	return a
}
