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
// Original source: github.com/micro/go-micro/v3/runtime/local/source/git/git_test.go

package git

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/onsi/gomega/types"

	. "github.com/onsi/gomega"
)

type parseCase struct {
	source   string
	expected *Source
}

func TestParseSource(t *testing.T) {
	cases := []parseCase{
		{
			source: "github.com/micro/services/helloworld",
			expected: &Source{
				Repo:   "github.com/micro/services",
				Folder: "helloworld",
				Ref:    "latest",
			},
		},
		{
			source: "github.com/micro/services/helloworld",
			expected: &Source{
				Repo:   "github.com/micro/services",
				Folder: "helloworld",
				Ref:    "latest",
			},
		},
		{
			source: "github.com/micro/services/helloworld@v1.12.1",
			expected: &Source{
				Repo:   "github.com/micro/services",
				Folder: "helloworld",
				Ref:    "v1.12.1",
			},
		},
		{
			source: "github.com/micro/services/helloworld@branchname",
			expected: &Source{
				Repo:   "github.com/micro/services",
				Folder: "helloworld",
				Ref:    "branchname",
			},
		},
		{
			source: "github.com/crufter/reponame/helloworld@branchname",
			expected: &Source{
				Repo:   "github.com/crufter/reponame",
				Folder: "helloworld",
				Ref:    "branchname",
			},
		},
	}
	for i, c := range cases {
		result, err := ParseSource(c.source)
		if err != nil {
			t.Fatalf("Failed case %v: %v", i, err)
		}
		if result.Folder != c.expected.Folder {
			t.Fatalf("Folder does not match for '%v', expected '%v', got '%v'", i, c.expected.Folder, result.Folder)
		}
		if result.Repo != c.expected.Repo {
			t.Fatalf("Repo address does not match for '%v', expected '%v', got '%v'", i, c.expected.Repo, result.Repo)
		}
		if result.Ref != c.expected.Ref {
			t.Fatalf("Ref does not match for '%v', expected '%v', got '%v'", i, c.expected.Ref, result.Ref)
		}
	}
}

type localParseCase struct {
	source     string
	expected   *Source
	workDir    string
	pathExists bool
}

func TestLocalParseSource(t *testing.T) {
	cases := []localParseCase{
		{
			source: ".",
			expected: &Source{
				Folder: "folder2",
				Ref:    "latest",
			},
			workDir:    "/folder1/folder2",
			pathExists: true,
		},
	}
	for i, c := range cases {
		result, err := ParseSourceLocal(c.workDir, c.source, func(s string) (bool, error) {
			return c.pathExists, nil
		})
		if err != nil {
			t.Fatalf("Failed case %v: %v", i, err)
		}
		if result.Folder != c.expected.Folder {
			t.Fatalf("Folder does not match for '%v', expected '%v', got '%v'", i, c.expected.Folder, result.Folder)
		}
		if result.Repo != c.expected.Repo {
			t.Fatalf("Repo address does not match for '%v', expected '%v', got '%v'", i, c.expected.Repo, result.Repo)
		}
		if result.Ref != c.expected.Ref {
			t.Fatalf("Ref does not match for '%v', expected '%v', got '%v'", i, c.expected.Ref, result.Ref)
		}
	}
}

type nameCase struct {
	fileContent string
	expected    string
}

func TestServiceNameExtract(t *testing.T) {
	cases := []nameCase{
		{
			fileContent: `func main() {
			// New Service
			service := micro.NewService(
				micro.Name("go.micro.service.helloworld"),
				micro.Version("latest"),
			)`,
			expected: "go.micro.service.helloworld",
		},
	}
	for i, c := range cases {
		result := extractServiceName([]byte(c.fileContent))
		if result != c.expected {
			t.Fatalf("Case %v, expected: %v, got: %v", i, c.expected, result)
		}
	}
}

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestDefaultBranch(t *testing.T) {
	tcs := []struct {
		name           string
		repo           string
		branchOrCommit string
		remoteBranch   string
		errMatcher     types.GomegaMatcher
	}{
		{name: "github-latest", repo: "https://github.com/micro/services", branchOrCommit: "latest", remoteBranch: "latest"},
		{name: "github-master", repo: "https://github.com/micro/services", branchOrCommit: "latest", remoteBranch: "master"},
		{name: "github-main", repo: "https://github.com/micro/services", branchOrCommit: "latest", remoteBranch: "main"},
		{name: "github-error", repo: "https://github.com/micro/services", branchOrCommit: "latest", remoteBranch: "someotherdefault", errMatcher: HaveOccurred()},
		{name: "gitlab-latest", repo: "https://gitlab.com/micro-test/basic-micro-service", branchOrCommit: "latest", remoteBranch: "latest"},
		{name: "gitlab-master", repo: "https://gitlab.com/micro-test/basic-micro-service", branchOrCommit: "latest", remoteBranch: "master"},
		{name: "gitlab-main", repo: "https://gitlab.com/micro-test/basic-micro-service", branchOrCommit: "latest", remoteBranch: "main"},
		{name: "gitlab-error", repo: "https://gitlab.com/micro-test/basic-micro-service", branchOrCommit: "latest", remoteBranch: "someotherdefault", errMatcher: HaveOccurred()},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			gInt := NewGitter(nil)
			gitter := gInt.(*binaryGitter)
			gitter.client = &http.Client{
				Transport: roundTripFunc(func(req *http.Request) *http.Response {
					if !strings.Contains(req.URL.String(), tc.remoteBranch) {
						return &http.Response{
							StatusCode: 404,
							Body:       ioutil.NopCloser(new(bytes.Buffer)),
							Header:     make(http.Header),
						}
					}
					if strings.HasSuffix(req.URL.String(), ".zip") {
						buf := new(bytes.Buffer)
						zipw := zip.NewWriter(buf)
						w, _ := zipw.Create("foo/bar")
						w.Write([]byte("foobar"))
						zipw.Close()

						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(buf),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					}

					if strings.HasSuffix(req.URL.String(), "tar.gz") {
						buf := new(bytes.Buffer)
						gz := gzip.NewWriter(buf)
						tw := tar.NewWriter(gz)
						hdr := &tar.Header{
							Name:     "foo",
							Mode:     0600,
							Typeflag: tar.TypeDir,
						}
						tw.WriteHeader(hdr)
						hdr = &tar.Header{
							Name:     "foo/bar",
							Mode:     0600,
							Size:     int64(len([]byte("foobar"))),
							Typeflag: tar.TypeReg,
						}
						tw.WriteHeader(hdr)
						tw.Write([]byte("foobar"))
						tw.Close()
						gz.Close()
						return &http.Response{
							StatusCode: 200,
							// Send response to be tested
							Body: ioutil.NopCloser(buf),
							// Must be set to non-nil value or it panics
							Header: make(http.Header),
						}
					}
					return &http.Response{
						StatusCode: 404,
						Body:       ioutil.NopCloser(new(bytes.Buffer)),
						Header:     make(http.Header),
					}
				}),
			}

			g := NewWithT(t)
			err := gitter.Checkout(tc.repo, tc.branchOrCommit)
			if tc.errMatcher != nil {
				g.Expect(err).To(tc.errMatcher)
			} else {
				g.Expect(err).To(BeNil())
			}

		})
	}

}
