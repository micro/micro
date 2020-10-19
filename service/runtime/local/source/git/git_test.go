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
	"testing"
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
