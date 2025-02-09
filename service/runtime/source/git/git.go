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
// Original source: github.com/micro/go-micro/v3/runtime/local/source/git/git.go

package git

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/teris-io/shortid"
)

const credentialsKey = "GIT_CREDENTIALS"

type Gitter interface {
	Checkout(repo, branchOrCommit string) error
	RepoDir() string
}

type binaryGitter struct {
	folder  string
	secrets map[string]string
	client  *http.Client
}

func (g *binaryGitter) Checkout(repo, branchOrCommit string) error {
	// The implementation of this method is questionable.
	// We use archives from github/gitlab etc which doesnt require the user to have got
	// and probably is faster than downloading the whole repo history,
	// but it comes with a bit of custom code for EACH host.
	// @todo probably we should fall back to git in case the archives are not available.
	doCheckout := func(repo, branchOrCommit string) error {
		if strings.HasPrefix(repo, "https://github.com") {
			return g.checkoutGithub(repo, branchOrCommit)
		} else if strings.HasPrefix(repo, "https://gitlab.com") {
			err := g.checkoutGitLabPublic(repo, branchOrCommit)
			if err != nil && len(g.secrets[credentialsKey]) > 0 {
				// If the public download fails, try getting it with tokens.
				// Private downloads needs a token for api project listing, hence
				// the weird structure of this code.
				return g.checkoutGitLabPrivate(repo, branchOrCommit)
			}
			return err
		}
		if len(g.secrets[credentialsKey]) > 0 {
			return g.checkoutAnyRemote(repo, branchOrCommit, true)
		}
		return g.checkoutAnyRemote(repo, branchOrCommit, false)
	}

	if branchOrCommit != "latest" {
		return doCheckout(repo, branchOrCommit)
	}
	// default branches
	defaults := []string{"latest", "master", "main", "trunk"}
	var err error
	for _, ref := range defaults {
		err = doCheckout(repo, ref)
		if err == nil {
			return nil
		}
	}
	return err
}

// This aims to be a generic checkout method. Currently only tested for bitbucket,
// see tests
func (g *binaryGitter) checkoutAnyRemote(repo, branchOrCommit string, useCredentials bool) error {
	repoFolder := strings.ReplaceAll(strings.ReplaceAll(repo, "/", "-"), "https:--", "")
	g.folder = filepath.Join(os.TempDir(),
		repoFolder+"-"+shortid.MustGenerate())
	err := os.MkdirAll(g.folder, 0755)
	if err != nil {
		return err
	}

	// Assumes remote address format is git@gitlab.com:micro-test/monorepo-test.git
	remoteAddr := fmt.Sprintf("https://%v", strings.TrimPrefix(repo, "https://"))
	if useCredentials {
		remoteAddr = fmt.Sprintf("https://%v@%v", g.secrets[credentialsKey], repo)
	}

	cmd := exec.Command("git", "clone", remoteAddr, "--depth=1", ".")
	cmd.Dir = g.folder
	outp, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Git clone failed: %v", string(outp))
	}

	cmd = exec.Command("git", "fetch", "origin", branchOrCommit, "--depth=1")
	cmd.Dir = g.folder
	outp, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Git fetch failed: %v", string(outp))
	}

	cmd = exec.Command("git", "checkout", "FETCH_HEAD")
	cmd.Dir = g.folder
	outp, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Git  checkout failed: %v", string(outp))
	}
	return nil
}

func (g *binaryGitter) checkoutGithub(repo, branchOrCommit string) error {
	// @todo if it's a commit it must not be checked out all the time
	repoFolder := strings.ReplaceAll(strings.ReplaceAll(repo, "/", "-"), "https:--", "")
	g.folder = filepath.Join(os.TempDir(),
		repoFolder+"-"+shortid.MustGenerate())

	url := fmt.Sprintf("%v/archive/%v.zip", repo, branchOrCommit)
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	req, _ := http.NewRequest("GET", url, nil)
	if len(g.secrets[credentialsKey]) > 0 {
		req.Header.Set("Authorization", "token "+g.secrets[credentialsKey])
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("Can't get zip: %v", err)
	}

	defer resp.Body.Close()
	// Github returns 404 for tar.gz files...
	// but still gives back a proper file so ignoring status code
	// for now.
	//if resp.StatusCode != 200 {
	//	return errors.New("Status code was not 200")
	//}

	src := g.folder + ".zip"
	// Create the file
	out, err := os.Create(src)
	if err != nil {
		return fmt.Errorf("Can't create source file %v src: %v", src, err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return unzip(src, g.folder, true)
}

func (g *binaryGitter) checkoutGitLabPublic(repo, branchOrCommit string) error {
	// Example: https://gitlab.com/micro-test/basic-micro-service/-/archive/master/basic-micro-service-master.tar.gz
	// @todo if it's a commit it must not be checked out all the time
	repoFolder := strings.ReplaceAll(strings.ReplaceAll(repo, "/", "-"), "https:--", "")
	g.folder = filepath.Join(os.TempDir(),
		repoFolder+"-"+shortid.MustGenerate())

	tarName := strings.ReplaceAll(strings.ReplaceAll(repo, "gitlab.com/", ""), "/", "-")
	url := fmt.Sprintf("%v/-/archive/%v/%v.tar.gz", repo, branchOrCommit, tarName)
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("Can't get zip: %v", err)
	}

	defer resp.Body.Close()

	src := g.folder + ".tar.gz"
	// Create the file
	out, err := os.Create(src)
	if err != nil {
		return fmt.Errorf("Can't create source file %v src: %v", src, err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	err = Uncompress(src, g.folder)
	if err != nil {
		return err
	}
	// Gitlab zip/tar has contents inside a folder
	// It has the format of eg. basic-micro-service-master-314b4a494ed472793e0a8bce8babbc69359aed7b
	// Since we don't have the commit at this point we must list the dir
	files, err := ioutil.ReadDir(g.folder)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("No contents in dir downloaded from gitlab: %v", g.folder)
	}
	g.folder = filepath.Join(g.folder, files[0].Name())
	return nil
}

func (g *binaryGitter) checkoutGitLabPrivate(repo, branchOrCommit string) error {

	repoFolder := strings.ReplaceAll(strings.ReplaceAll(repo, "/", "-"), "https:--", "")
	g.folder = filepath.Join(os.TempDir(),
		repoFolder+"-"+shortid.MustGenerate())
	tarName := strings.ReplaceAll(strings.ReplaceAll(repo, "gitlab.com/", ""), "/", "-")

	url := fmt.Sprintf("%v/-/archive/%v/%v.tar.gz?private_token=%v", repo, branchOrCommit, tarName, g.secrets[credentialsKey])

	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("Can't get zip: %v", err)
	}

	defer resp.Body.Close()

	src := g.folder + ".tar.gz"
	// Create the file
	out, err := os.Create(src)
	if err != nil {
		return fmt.Errorf("Can't create source file %v src: %v", src, err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	err = Uncompress(src, g.folder)
	if err != nil {
		return err
	}
	// Gitlab zip/tar has contents inside a folder
	// It has the format of eg. basic-micro-service-master-314b4a494ed472793e0a8bce8babbc69359aed7b
	// Since we don't have the commit at this point we must list the dir
	files, err := ioutil.ReadDir(g.folder)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("No contents in dir downloaded from gitlab: %v", g.folder)
	}
	g.folder = filepath.Join(g.folder, files[0].Name())
	return nil
}

func (g *binaryGitter) RepoDir() string {
	return g.folder
}

func NewGitter(secrets map[string]string) Gitter {
	tmpdir, _ := ioutil.TempDir(os.TempDir(), "git-src-*")

	return &binaryGitter{
		folder:  tmpdir,
		secrets: secrets,
		client:  &http.Client{},
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func dirifyRepo(s string) string {
	s = strings.ReplaceAll(s, "https://", "")
	s = strings.ReplaceAll(s, "/", "-")
	return s
}

// exists returns whether the given file or directory exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// GetRepoRoot determines the repo root from a full path.
// Returns empty string and no error if not found
func GetRepoRoot(fullPath string) (string, error) {
	// traverse parent directories
	prev := fullPath
	for {
		current := prev

		// check for a go.mod
		goExists, err := pathExists(filepath.Join(current, "go.mod"))
		if err != nil {
			return "", err
		}
		if goExists {
			return current, nil
		}

		// check for .git
		gitExists, err := pathExists(filepath.Join(current, ".git"))
		if err != nil {
			return "", err
		}
		if gitExists {
			return current, nil
		}

		prev = filepath.Dir(current)
		// reached top level, see:
		// https://play.golang.org/p/rDgVdk3suzb
		if current == prev {
			break
		}
	}
	return "", nil
}

// Source is not just git related @todo move
type Source struct {
	// is it a local folder intended for a local runtime?
	Local bool
	// absolute path to service folder in local mode
	FullPath string
	// path of folder to repo root
	// be it local or github repo
	Folder string
	// github ref
	Ref string
	// for cloning purposes
	// blank for local
	Repo string
	// dir to repo root
	// blank for non local
	LocalRepoRoot string
}

// Name to be passed to RPC call runtime.Create Update Delete
// eg: `helloworld/api`, `crufter/myrepo/helloworld/api`, `localfolder`
func (s *Source) RuntimeName() string {
	if len(s.Folder) == 0 {
		// This is the case for top level url source ie. gitlab.com/micro-test/basic-micro-service
		return path.Base(s.Repo)
	}
	return path.Base(s.Folder)
}

// Source to be passed to RPC call runtime.Create Update Delete
// eg: `helloworld`, `github.com/crufter/myrepo/helloworld`, `/path/to/localrepo/localfolder`
func (s *Source) RuntimeSource() string {
	if s.Local && s.LocalRepoRoot != s.FullPath {
		relpath, _ := filepath.Rel(s.LocalRepoRoot, s.FullPath)
		return relpath
	}
	if s.Local {
		return s.FullPath
	}
	if len(s.Folder) == 0 {
		return s.Repo
	}
	return fmt.Sprintf("%v/%v", s.Repo, s.Folder)
}

// ParseSource parses a `micro run/update/kill` source.
func ParseSource(source string) (*Source, error) {
	if !strings.Contains(source, "@") {
		source += "@latest"
	}
	ret := &Source{}
	refs := strings.Split(source, "@")
	ret.Ref = refs[1]
	parts := strings.Split(refs[0], "/")

	max := 3
	if len(parts) < 3 {
		max = len(parts)
	}
	ret.Repo = strings.Join(parts[0:max], "/")

	if len(parts) > 1 {
		ret.Folder = strings.Join(parts[3:], "/")
	}

	return ret, nil
}

// ParseSourceLocal a version of ParseSource that detects and handles local paths.
// Workdir should be used only from the CLI @todo better interface for this function.
// PathExistsFunc exists only for testing purposes, to make the function side effect free.
func ParseSourceLocal(workDir, source string, pathExistsFunc ...func(path string) (bool, error)) (*Source, error) {
	var pexists func(string) (bool, error)
	if len(pathExistsFunc) == 0 {
		pexists = pathExists
	} else {
		pexists = pathExistsFunc[0]
	}
	isLocal, localFullPath := IsLocal(workDir, source, pexists)
	if isLocal {
		localRepoRoot, err := GetRepoRoot(localFullPath)
		if err != nil {
			return nil, err
		}
		var folder string
		// If the local repo root is a top level folder, we are not in a git repo.
		// In this case, we should take the last folder as folder name.
		if localRepoRoot == "" {
			folder = filepath.Base(localFullPath)
		} else {
			folder = strings.ReplaceAll(localFullPath, localRepoRoot+string(filepath.Separator), "")
		}

		return &Source{
			Local:         true,
			Folder:        folder,
			FullPath:      localFullPath,
			LocalRepoRoot: localRepoRoot,
			Ref:           "latest", // @todo consider extracting branch from git here
		}, nil
	}
	return ParseSource(source)
}

// IsLocal tries returns true and full path of directory if the path is a local one, and
// false and empty string if not.
func IsLocal(workDir, source string, pathExistsFunc ...func(path string) (bool, error)) (bool, string) {
	var pexists func(string) (bool, error)
	if len(pathExistsFunc) == 0 {
		pexists = pathExists
	} else {
		pexists = pathExistsFunc[0]
	}
	// Check for absolute path
	// @todo "/" won't work for Windows
	if exists, err := pexists(source); strings.HasPrefix(source, "/") && err == nil && exists {
		return true, source
		// Check for path relative to workdir
	} else if exists, err := pexists(filepath.Join(workDir, source)); err == nil && exists {
		return true, filepath.Join(workDir, source)
	}
	return false, ""
}

// CheckoutSource checks out a git repo (source) into a local temp directory. It will return the
// source of the local repo an error if one occured. Secrets can optionally be passed if the repo
// is private.
func CheckoutSource(source *Source, secrets map[string]string) (string, error) {
	gitter := NewGitter(secrets)
	repo := source.Repo
	if !strings.Contains(repo, "https://") {
		repo = "https://" + repo
	}
	if err := gitter.Checkout(repo, source.Ref); err != nil {
		return "", err
	}
	return gitter.RepoDir(), nil
}

// code below is not used yet

var nameExtractRegexp = regexp.MustCompile(`((micro|web)\.Name\(")(.*)("\))`)

func extractServiceName(fileContent []byte) string {
	hits := nameExtractRegexp.FindAll(fileContent, 1)
	if len(hits) == 0 {
		return ""
	}
	hit := string(hits[0])
	return strings.Split(hit, "\"")[1]
}

// Uncompress is a modified version of: https://gist.github.com/mimoo/25fc9716e0f1353791f5908f94d6e726
func Uncompress(src string, dst string) error {
	file, err := os.OpenFile(src, os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	// ungzip
	zr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	// untar
	tr := tar.NewReader(zr)

	// uncompress each element
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		target := header.Name

		// validate name against path traversal
		if !validRelPath(header.Name) {
			return fmt.Errorf("tar contained invalid name error %q\n", target)
		}

		// add dst + re-format slashes according to system
		target = filepath.Join(dst, header.Name)
		// if no join is needed, replace with ToSlash:
		// target = filepath.ToSlash(header.Name)

		// check the type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it (with 0755 permission)
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				// @todo think about this:
				// if we don't nuke the folder, we might end up with files from
				// the previous decompress.
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// if it's a file create it (with same permission)
		case tar.TypeReg:
			// the truncating is probably unnecessary due to the `RemoveAll` of folders
			// above
			fileToWrite, err := os.OpenFile(target, os.O_TRUNC|os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			// copy over contents
			if _, err := io.Copy(fileToWrite, tr); err != nil {
				return err
			}
			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			fileToWrite.Close()
		}
	}
	return nil
}

// check for path traversal and correct forward slashes
func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}

// taken from https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
func unzip(src, dest string, skipTopFolder bool) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		r.Close()
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			rc.Close()
		}()
		if skipTopFolder {
			f.Name = strings.Join(strings.Split(f.Name, string(filepath.Separator))[1:], string(filepath.Separator))
		}
		// zip slip https://snyk.io/research/zip-slip-vulnerability
		destpath, err := zipSafeFilePath(dest, f.Name)
		if err != nil {
			return err
		}
		path := destpath
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				f.Close()
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

// zipSafeFilePath checks whether the file path is safe to use or is a zip slip attack https://snyk.io/research/zip-slip-vulnerability
func zipSafeFilePath(destination, filePath string) (string, error) {
	if len(filePath) == 0 {
		return filepath.Join(destination, filePath), nil
	}
	destination, _ = filepath.Abs(destination) //explicit the destination folder to prevent that 'string.HasPrefix' check can be 'bypassed' when no destination folder is supplied in input
	destpath := filepath.Join(destination, filePath)
	if !strings.HasPrefix(destpath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return "", fmt.Errorf("%s: illegal file path", filePath)
	}
	return destpath, nil
}
