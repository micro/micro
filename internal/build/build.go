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
// Original source: github.com/micro/go-micro/v3/build/build.go

// Package build is for building source into a package
package build

// Source is the source of a build
type Source struct {
	// Path to the source if local
	Path string
	// Language is the language of code
	Language string
	// Location of the source
	Repository string
}

// Package is packaged format for source
type Package struct {
	// Name of the package
	Name string
	// Location of the package
	Path string
	// Type of package e.g tarball, binary, docker
	Type string
	// Source of the package
	Source *Source
}
