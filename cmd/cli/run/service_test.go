package runtime

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/micro/micro/v3/service/runtime/source/git"
	"github.com/onsi/gomega/types"

	. "github.com/onsi/gomega"

	"github.com/micro/micro/v3/test/fakes"
)

func TestFmtDuration(t *testing.T) {
	tcs := []struct {
		seconds  int64
		expected string
	}{
		{seconds: 15, expected: "15s"},
		{seconds: 0, expected: "0s"},
		{seconds: 60, expected: "1m0s"},
		{seconds: 75, expected: "1m15s"},
		{seconds: 903, expected: "15m3s"},
		{seconds: 4532, expected: "1h15m32s"},
		{seconds: 82808, expected: "23h0m8s"},
		{seconds: 86400, expected: "1d0h0m0s"},
		{seconds: 1006400, expected: "11d15h33m20s"},
		{seconds: 111006360, expected: "1284d19h6m0s"},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			res := fmtDuration(time.Duration(tc.seconds) * time.Second)
			if res != tc.expected {
				t.Errorf("Expected %s but got %s", tc.expected, res)
			}
		})
	}
}

func TestSourceExists(t *testing.T) {
	tcs := []struct {
		name        string
		source      string
		errMatcher  types.GomegaMatcher
		expectedURL string
		callCount   int
	}{
		{name: "github-monorepo-main", source: "github.com/micro/services/helloworld", callCount: 3, expectedURL: "https://api.github.com/repos/micro/services/contents/helloworld?ref=main"},
		{name: "github-monorepo-specific", source: "github.com/micro/services/hello/world@foobar", callCount: 1, expectedURL: "https://api.github.com/repos/micro/services/contents/hello/world?ref=foobar"},
		{name: "github-multirepo-master", source: "github.com/micro/services", callCount: 2, expectedURL: "https://api.github.com/repos/micro/services/contents/?ref=master"},
		{name: "github-multirepo-specific", source: "github.com/micro/services@foobar", callCount: 1, expectedURL: "https://api.github.com/repos/micro/services/contents/?ref=foobar"},
		{name: "github-multirepo-specific-error", source: "github.com/micro/services@foobar", callCount: 1, expectedURL: "", errMatcher: HaveOccurred()},
		{name: "gitlab-monorepo", source: "gitlab.com/micro-test/basic-micro-services", callCount: 1, expectedURL: "https://gitlab.com/micro-test/basic-micro-services"},
		{name: "gitlab-multirepo", source: "gitlab.com/micro-test/basic-micro-services/foobar", callCount: 1, expectedURL: "https://gitlab.com/micro-test/basic-micro-services"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			fakeTripper := &fakes.FakeRoundTripper{}
			httpClient.Transport = fakeTripper
			fakeTripper.RoundTripStub = func(request *http.Request) (*http.Response, error) {
				if request.URL.String() != tc.expectedURL {
					return &http.Response{
						StatusCode: 404,
						Body:       ioutil.NopCloser(new(bytes.Buffer)),
						Header:     make(http.Header),
					}, nil
				}
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(new(bytes.Buffer)),
					Header:     make(http.Header),
				}, nil
			}
			g := NewWithT(t)
			src, err := git.ParseSource(tc.source)
			g.Expect(err).To(BeNil())
			err = sourceExists(src)
			if tc.errMatcher != nil {
				g.Expect(err).To(tc.errMatcher)
			} else {
				g.Expect(err).To(BeNil())
			}
			g.Expect(fakeTripper.RoundTripCallCount()).To(Equal(tc.callCount))
		})

	}
}
