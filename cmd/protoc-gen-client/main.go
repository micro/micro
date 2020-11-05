// Apache 2.0 => https://github.com/microhq/clients
package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	protocURL        = "https://github.com/protocolbuffers/protobuf/releases/download/v%s/protoc-%s-linux-x86_64.zip"
	protocGenJavaURL = "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/%s/protoc-gen-grpc-java-%s-linux-x86_64.exe"
	protocGenGoURL   = "https://github.com/golang/protobuf/archive/v%s.zip"
	protocGenJsURL   = "https://registry.npmjs.org/grpc-tools/-/grpc-tools-%v.tgz"
	//protocGenJsURL       = "https://registry.npmjs.org/grpc/-/grpc-%v.tgz"
	protocBin            = "protoc"
	protocGenJavaVersion = "1.28.0"
	protocGenJavaBin     = "protoc-gen-grpc-java"
	protocVersion        = "3.11.4"
	protocGenGoVersion   = "1.3.5"
	protocGenGoBin       = "protoc-gen-go"
	protocGenJsBin       = "grpc_tools_node_protoc_plugin"
	protocGenJsVersion   = "1.8.1"
	protocGenRubyBin     = "grpc_tools_ruby_protoc"
	//protocGenJsVersion = "1.24.2"
)

func download(path string, lang string) error {
	var url string

	toolPath := filepath.Join(path, fmt.Sprintf("%s-tool", lang))

	switch lang {
	case "go":
		if p, err := exec.LookPath(filepath.Join(toolPath, protocGenGoBin)); err == nil && len(p) > 0 {
			log.Printf("%s tool already installed", protocGenGoBin)
			return nil
		}
		var cmd *exec.Cmd
		if gomod, ok := os.LookupEnv("GOMOD"); ok && len(gomod) > 0 {
			cmd = exec.Command("go", "get", "-v", fmt.Sprintf("github.com/golang/protobuf/protoc-gen-go@v%s", protocGenGoVersion))
		} else {
			cmd = exec.Command("go", "get", "-v", "github.com/golang/protobuf/protoc-gen-go")
		}

		cmd.Env = append(cmd.Env, []string{
			fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
			fmt.Sprintf("GOBIN=%s", toolPath),
			fmt.Sprintf("XDG_CACHE_HOME=%s/cache", path),
			fmt.Sprintf("GOPATH=%s/cache/go-path", path),
		}...)
		if buf, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%s", buf)
		} else {
			log.Printf("%s", buf)
		}
	case "protoc":
		if p, err := exec.LookPath(filepath.Join(toolPath, "bin", protocBin)); err == nil && len(p) > 0 {
			log.Printf("%s tool already installed", protocBin)
			return nil
		}
		url = fmt.Sprintf(protocURL, protocVersion, protocVersion)
	case "java":
		if p, err := exec.LookPath(filepath.Join(toolPath, protocGenJavaBin)); err == nil && len(p) > 0 {
			log.Printf("%s tool already installed", protocGenJavaBin)
			return nil
		}
		url = fmt.Sprintf(protocGenJavaURL, protocGenJavaVersion, protocGenJavaVersion)
	case "node":
		if _, err := os.Stat(filepath.Join(toolPath, protocGenJsBin)); err == nil {
			log.Printf("%s tool already installed", protocGenJsBin)
			return nil
		}
		url = fmt.Sprintf(protocGenJsURL, protocGenJsVersion)
	}

	if err := os.MkdirAll(filepath.Join(path, fmt.Sprintf("%s-tool", lang)), os.FileMode(0755)); err != nil {
		return err
	}

	if len(url) == 0 {
		return nil
	}

	rsp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	// Create the file
	out, err := os.Create(filepath.Join(path, "tmpfile"))
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, rsp.Body); err != nil {
		return err
	}

	switch lang {
	case "node":
		if err = untar(filepath.Join(path, "tmpfile"), filepath.Join(path, fmt.Sprintf("%s-tool", lang))); err != nil {
			return err
		}
		if err = os.Symlink("bin/protoc_plugin.js", filepath.Join(path, fmt.Sprintf("%s-tool", lang), "grpc_tools_node_protoc_plugin")); err != nil {
			return err
		}
		if err = os.Remove(filepath.Join(path, "tmpfile")); err != nil {
			return err
		}
	case "java":
		if err = os.Rename(
			filepath.Join(path, "tmpfile"),
			filepath.Join(path, fmt.Sprintf("%s-tool", lang), "protoc-gen-grpc-java")); err != nil {
			return err
		}
		if err = os.Chmod(filepath.Join(path, fmt.Sprintf("%s-tool", lang), "protoc-gen-grpc-java"), os.FileMode(0755)); err != nil {
			return err
		}
	case "protoc":
		if err = unzip(filepath.Join(path, "tmpfile"), filepath.Join(path, fmt.Sprintf("%s-tool", lang))); err != nil {
			return err
		}
		if err = os.Remove(filepath.Join(path, "tmpfile")); err != nil {
			return err
		}
	}

	return nil
}

func generate(lang string, tool string, arg string, src string, dst string, args ...string) error {
	var err error

	//	pwd, err := os.Getwd()
	//	if err != nil {
	//		return err
	//	}

	//	if err = download(pwd, lang); err != nil {
	//		return err
	//	}

	chFile, chErr := findProto(src)

	for {
		select {
		case err := <-chErr:
			return err
		case proto := <-chFile:
			fmt.Printf("proto: %s\n", proto)
			if len(proto) == 0 {
				return nil
			}

			dstpath := filepath.Join(dst, lang)
			if err = os.MkdirAll(dstpath, os.FileMode(0755)); err != nil {
				return err
			}

			var cmdargs []string
			switch lang {
			case "rust":
				cmdargs = append(cmdargs,
					"protoc",
					fmt.Sprintf("-I%s", src),
					fmt.Sprintf("-I%s", filepath.Dir(proto)),
					fmt.Sprintf("-I%s", filepath.Join(src, "vendor")),
					fmt.Sprintf("--%s_out=%s:%s", tool, arg, dstpath),
					fmt.Sprintf("--rust-grpc_out=%s", dstpath),
					fmt.Sprintf("%s", proto),
					fmt.Sprintf("%s", strings.Join(args, " ")),
				)
			case "ruby":
				cmdargs = append(cmdargs,
					"protoc",
					fmt.Sprintf("-I%s", src),
					fmt.Sprintf("-I%s", filepath.Dir(proto)),
					fmt.Sprintf("-I%s", filepath.Join(src, "vendor")),
					fmt.Sprintf("--%s_out=%s:%s", tool, arg, dstpath),
					fmt.Sprintf("--grpc_out=%s", dstpath),
					fmt.Sprintf("%s", proto),
					fmt.Sprintf("%s", strings.Join(args, " ")),
				)
			case "java":
				cmdargs = append(cmdargs,
					"protoc",
					fmt.Sprintf("-I%s", src),
					fmt.Sprintf("-I%s", filepath.Dir(proto)),
					fmt.Sprintf("-I%s", filepath.Join(src, "vendor")),
					fmt.Sprintf("--%s_out=%s:%s", tool, arg, dstpath),
					fmt.Sprintf("--grpc-java_out=%s", dstpath),
					fmt.Sprintf("%s", proto),
					fmt.Sprintf("%s", strings.Join(args, " ")),
				)
			case "go":
				cmdargs = append(cmdargs,
					"protoc",
					fmt.Sprintf("-I%s", src),
					fmt.Sprintf("-I%s", filepath.Dir(proto)),
					fmt.Sprintf("-I%s", filepath.Join(src, "vendor")),
					fmt.Sprintf("--%s_out=%s:%s", tool, arg, dstpath),
					fmt.Sprintf("%s", proto),
				)
			case "python":
				cmdargs = append(cmdargs,
					"python3",
					"-m",
					"grpc.tools.protoc",
					fmt.Sprintf("-I%s", src),
					fmt.Sprintf("-I%s", filepath.Dir(proto)),
					fmt.Sprintf("-I%s", filepath.Join(src, "vendor")),
					fmt.Sprintf("--%s_out=%s:%s", tool, arg, dstpath),
					fmt.Sprintf("--grpc_python_out=%s", dstpath),
					fmt.Sprintf("%s", proto),
				)
			case "node":
				cmdargs = append(cmdargs,
					"protoc",
					fmt.Sprintf("-I%s", src),
					fmt.Sprintf("-I%s", filepath.Dir(proto)),
					fmt.Sprintf("-I%s", filepath.Join(src, "vendor")),
					fmt.Sprintf("--js_out=%s:%s", arg, dstpath),
					fmt.Sprintf("--grpc_out=%s", dstpath),
					fmt.Sprintf("%s", proto),
					fmt.Sprintf("%s", strings.Join(args, " ")),
				)
			}
			fmt.Printf("%s\n", strings.Join(cmdargs, " "))
			if out, err := exec.Command(cmdargs[0], cmdargs[1:]...).CombinedOutput(); err != nil {
				log.Fatalf("%s\n", out)
			}
		}
	}
}

func findProto(src string) (chan string, chan error) {
	chFile := make(chan string)
	chErr := make(chan error)

	go func() {
		err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			// exit on err
			if err != nil {
				return err
			}

			// skip non proto files
			if !info.Mode().IsRegular() || !strings.Contains(info.Name(), ".proto") {
				return nil
			}

			chFile <- path

			return nil
		})

		if err != nil {
			chErr <- err
		}

		close(chFile)
		close(chErr)
	}()

	return chFile, chErr
}

func untar(src string, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}

	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			name := strings.TrimPrefix(header.Name, "package/")
			if err := os.MkdirAll(filepath.Join(dst, name), os.FileMode(0755)); err != nil {
				return err
			}
		case tar.TypeSymlink:
			name := strings.TrimPrefix(header.Name, "package/")
			if err := os.Symlink(header.Linkname, filepath.Join(dst, name)); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			name := strings.TrimPrefix(header.Name, "package/")
			if err := os.MkdirAll(filepath.Join(dst, filepath.Dir(name)), os.FileMode(0755)); err != nil {
				return err
			}
			outFile, err := os.Create(filepath.Join(dst, name))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		default:
			panic(header.Name)
		}

	}
	return nil
}

func unzip(src string, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dst, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

var (
	srcDir  = flag.String("srcdir", "", "source dir")
	dstDir  = flag.String("dstdir", "", "target dir")
	langs   = flag.String("langs", "go,java,node,python", "languages to generate")
	toolDir = flag.String("tooldir", "/tmp", "tool dir")
)

func main() {
	flag.Parse()
	var err error

	if err = download(*toolDir, "protoc"); err != nil {
		log.Fatal(err)
	}

	for _, lang := range strings.Split(*langs, ",") {
		if err = download(*toolDir, lang); err != nil {
			log.Fatal(err)
		}

		switch lang {
		case "rust":
			ppath, err := exec.LookPath("protoc-gen-rust-grpc")
			if err != nil {
				log.Fatal(err)
			}
			if err = generate(lang, "rust", "", *srcDir, *dstDir, fmt.Sprintf("--plugin=protoc-gen-rust-grpc=%s", ppath)); err != nil {
				log.Fatal(err)
			}
		case "ruby":
			ppath, err := exec.LookPath("grpc_tools_ruby_protoc_plugin")
			if err != nil {
				log.Fatal(err)
			}

			if err = generate(lang, "ruby", "", *srcDir, *dstDir, fmt.Sprintf("--plugin=protoc-gen-grpc=%s", ppath)); err != nil {
				log.Fatal(err)
			}
		case "go":
			if err = generate(lang, "go", "plugins=grpc,paths=source_relative", *srcDir, *dstDir); err != nil {
				log.Fatal(err)
			}
			if err = generate(lang, "micro", "paths=source_relative", *srcDir, *dstDir); err != nil {
				log.Fatal(err)
			}
		case "java":
			if err = generate(lang, "java", "", *srcDir, *dstDir, fmt.Sprintf("--plugin=protoc-gen-grpc-java=%s", filepath.Join(*toolDir, "java-tool", "protoc-gen-grpc-java"))); err != nil {
				log.Fatal(err)
			}
		case "python":
			if err = generate(lang, "python", "", *srcDir, *dstDir); err != nil {
				log.Fatal(err)
			}
		case "node":
			ppath, err := exec.LookPath("grpc_tools_node_protoc_plugin")
			if err != nil {
				log.Fatal(err)
			}

			if err = generate(lang, "node", "import_style=commonjs,binary", *srcDir, *dstDir, fmt.Sprintf("--plugin=protoc-gen-grpc=%s", ppath)); err != nil {
				log.Fatal(err)
			}
		}
	}
}
