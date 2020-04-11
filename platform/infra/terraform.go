package infra

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// TerraformModule is a task that fetches and applies a terraform module
type TerraformModule struct {
	// ID is a persistent unique ID for the name of the stored state
	ID string
	// Name is the name of the module - for logging purposes
	Name string
	// Path is the path to the module. It's set to working directory for terraform
	Path string
	// Source is a net.URL to the module
	Source string
	// Any environment variables to pass to terraform
	Env map[string]string
	// Any terraform variables
	Variables map[string]string
	// Any remote states to import key = state name, value = remote state ID
	RemoteStates map[string]string
	// Dry-run
	DryRun bool
}

// Validate attempts to fetch terraform code then runs terraform init and terraform validate
func (t *TerraformModule) Validate() error {
	if err := os.MkdirAll(t.Path, 0o777); err != nil {
		return err
	}
	if err := os.MkdirAll("/tmp/micro-platform-plugin-cache", 0o700); err != nil {
		return err
	}

	u, err := url.Parse(t.Source)
	if err != nil {
		return err
	}
	switch u.Scheme {
	case "http":
		fallthrough
	case "https":
		return errors.New("TODO: Download and extract " + u.String() + " to " + t.Path)
	case "git":
		return errors.New("TODO: Clone " + u.String() + " to " + t.Path)
	default:
		if len(u.Scheme) == 0 {
			fmt.Fprintf(os.Stderr, "[%s] No source scheme provided, assuming path to directory\n", t.Name)
			if _, err := os.Stat(u.Path); err != nil {
				return err
			}
			if err := filepath.Walk(u.Path, t.filecopy); err != nil {
				return errors.Wrap(err, "filepath.Walk failed")
			}
		} else {
			return errors.New("Module " + t.Name + " Scheme " + u.Scheme + " not supported")
		}
	}

	// Set up remote state storage
	if err := t.generateBackendConfig(); err != nil {
		return err
	}

	// import any remote states
	if err := t.generateRemoteStateDataSources(); err != nil {
		return err
	}

	// Initialise terraform and validate the syntax is correct
	if err := t.execTerraform(context.Background(), "init"); err != nil {
		return err
	}
	return t.execTerraform(context.Background(), "validate")
}

// Plan runs terraform plan
func (t *TerraformModule) Plan() error {
	return t.execTerraform(context.Background(), "plan")
}

// Apply runs terraform apply
func (t *TerraformModule) Apply() error {
	if t.DryRun {
		_, err := fmt.Fprintf(os.Stderr, "[%s] Dry run enabled, skipping apply\n", t.Name)
		return err
	}
	return t.execTerraform(context.Background(), "apply", "-auto-approve")
}

// Destroy runs terraform apply
func (t *TerraformModule) Destroy() error {
	if t.DryRun {
		_, err := fmt.Fprintf(os.Stderr, "[%s] Dry run enabled, skipping destroy\n", t.Name)
		return err
	}
	return t.execTerraform(context.Background(), "destroy", "-auto-approve")
}

// Finalise removes the directory
func (t *TerraformModule) Finalise() error {
	return os.RemoveAll(t.Path)
}

func (t *TerraformModule) execTerraform(ctx context.Context, args ...string) error {
	// Set up terraform command
	tf := exec.CommandContext(ctx, "terraform", args...)
	tf.Dir = t.Path
	tf.Env = os.Environ()
	for k, v := range t.Env {
		tf.Env = append(tf.Env, fmt.Sprintf("%s=%s", k, v))
	}
	for k, v := range t.Variables {
		tf.Env = append(tf.Env, fmt.Sprintf("TF_VAR_%s=%s", k, v))
	}
	tf.Env = append(tf.Env, "TF_PLUGIN_CACHE_DIR=/tmp/micro-platform-plugin-cache")
	stdout, err := tf.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "StdoutPipe failed")
	}
	stderr, err := tf.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "StderrPipe failed")
	}

	// Wait so we don't truncate output from the underlying terraform binary
	ioWait := make(chan struct{})
	defer func() {
		// wait for the buffered readers/writers to finish
		<-ioWait
		<-ioWait
	}()

	for _, ioPair := range []struct {
		in  io.ReadCloser
		out *os.File
	}{
		{in: stdout, out: os.Stdout},
		{in: stderr, out: os.Stderr},
	} {
		go func(name string, in io.ReadCloser, out *os.File, done chan<- struct{}) {
			r := bufio.NewReader(in)
			defer func() { done <- struct{}{} }()
			defer in.Close()
			for {
				s, err := r.ReadString('\n')
				if err == nil || err == io.EOF {
					if len(strings.TrimSpace(s)) != 0 {
						fmt.Fprintf(out, "[%s] %s", name, s)
					}
					if err == io.EOF {
						return
					}
				} else {
					fmt.Fprintf(out, "[%s] Error: %s\n", name, err.Error())
					return
				}
			}
		}(t.Name, ioPair.in, ioPair.out, ioWait)
	}
	if err := tf.Start(); err != nil {
		return errors.Wrap(err, "Couldn't execute terraform")
	}

	return tf.Wait()
}

func (t *TerraformModule) filecopy(path string, fi os.FileInfo, err error) error {
	if strings.HasPrefix(path, "./") ||
		strings.Contains(path, "tfstate") ||
		strings.Contains(path, ".terraform") ||
		strings.Contains(path, ".git") {
		// skip
	} else if fi.IsDir() {
		if err := os.MkdirAll(filepath.Join(t.Path, t.cleanPath(path)), fi.Mode()); err != nil {
			return err
		}
	} else if fi.Mode().IsRegular() {
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		dest, err := os.OpenFile(filepath.Join(t.Path, t.cleanPath(path)), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fi.Mode())
		if err != nil {
			return err
		}
		if _, err := io.Copy(dest, src); err != nil {
			return err
		}
		// Explicitly check dest.Close(), as the OS can sometimes error on closing
		// a writable file (eg EBADF, EINTR, EIO) and defer() would swallow it
		if err := dest.Close(); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(os.Stderr, "[%s] Encountered non regular file or directory: %s\n", t.Name, path)
	}
	return nil
}

func (t *TerraformModule) cleanPath(in string) string {
	prefix := strings.TrimPrefix(t.Source, "."+string([]rune{filepath.Separator}))
	return strings.TrimPrefix(in, prefix+string([]rune{filepath.Separator}))
}

func (t *TerraformModule) generateBackendConfig() error {
	stateStore := viper.GetString("state-store")
	if len(stateStore) == 0 {
		stateStore = viper.GetString("cloud-provider")
	}
	switch stateStore {
	case "aws":
		return t.generateBackendConfigAWS()
	case "azure":
		return t.generateBackendConfigAzure()
	default:
		return errors.New(stateStore + " is not a supported remote state store")
	}
}

func (t *TerraformModule) generateBackendConfigAWS() error {
	backend := template.Must(template.New(t.ID + "backend").Parse(tfS3BackendTemplate))
	f, err := os.OpenFile(filepath.Join(t.Path, "backend-config-micro-platform.tf"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if err := backend.Execute(f, struct {
		Key         string
		Region      string
		StateBucket string
		LockTable   string
	}{
		Key: t.ID,
		Region: func() string {
			if r := os.Getenv("AWS_REGION"); len(r) != 0 {
				return r
			}
			return "eu-west-2"
		}(),
		StateBucket: viper.GetString("aws-s3-bucket"),
		LockTable:   viper.GetString("aws-dynamodb-table"),
	}); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}

func (t *TerraformModule) generateBackendConfigAzure() error {
	backend := template.Must(template.New(t.ID + "backend").Parse(tfAzureRMBackendTemplate))
	f, err := os.OpenFile(filepath.Join(t.Path, "backend-config-micro-platform.tf"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if err := backend.Execute(f, struct {
		ResourceGroupName  string
		StorageAccountName string
		ContainerName      string
		Key                string
	}{
		ResourceGroupName:  viper.GetString("azure-state-resource-group"),
		StorageAccountName: viper.GetString("azure-storage-account"),
		ContainerName:      viper.GetString("azure-storage-container"),
		Key:                t.ID,
	}); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}

func (t *TerraformModule) generateRemoteStateDataSources() error {
	stateStore := viper.GetString("state-store")
	if len(stateStore) == 0 {
		stateStore = viper.GetString("cloud-provider")
	}
	switch stateStore {
	case "aws":
		return t.generateRemoteStateAws()
	case "azure":
		return t.generateRemoteStateAzure()
	default:
		return errors.New(stateStore + " is not a supported remote state store")
	}
}

func (t *TerraformModule) generateRemoteStateAws() error {
	remote := template.Must(template.New(t.ID + "remote").Parse(tfS3RemoteStateTemplate))
	f, err := os.OpenFile(filepath.Join(t.Path, "remote-state-data-sources-micro-platform.tf"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	for k, v := range t.RemoteStates {
		if err := remote.Execute(f, struct {
			RemoteStateName string
			StateBucket     string
			LockTable       string
			Key             string
			Region          string
		}{
			RemoteStateName: k,
			StateBucket:     viper.GetString("aws-s3-bucket"),
			LockTable:       viper.GetString("aws-dynamodb-table"),
			Key:             v,
			Region: func() string {
				if r := os.Getenv("AWS_REGION"); len(r) != 0 {
					return r
				}
				return "eu-west-2"
			}(),
		}); err != nil {
			return err
		}
	}
	return f.Close()
}

func (t *TerraformModule) generateRemoteStateAzure() error {
	remote := template.Must(template.New(t.ID + "remote").Parse(tfAzureRmStateTemplate))
	f, err := os.OpenFile(filepath.Join(t.Path, "remote-state-data-sources-micro-platform.tf"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	for k, v := range t.RemoteStates {
		if err := remote.Execute(f, struct {
			RemoteStateName    string
			ResourceGroupName  string
			StorageAccountName string
			ContainerName      string
			Key                string
		}{
			RemoteStateName:    k,
			ResourceGroupName:  viper.GetString("azure-state-resource-group"),
			StorageAccountName: viper.GetString("azure-storage-account"),
			ContainerName:      viper.GetString("azure-storage-container"),
			Key:                v,
		}); err != nil {
			return err
		}
	}
	return f.Close()
}

const tfS3BackendTemplate = `terraform {
  backend "s3" {
    bucket         = "{{.StateBucket}}"
    dynamodb_table = "{{.LockTable}}"
    key            = "{{.Key}}"
    region         = "{{.Region}}"
  }
}
`
const tfS3RemoteStateTemplate = `data "terraform_remote_state" "{{.RemoteStateName}}" {
  backend = "s3"

  config = {
    bucket         = "{{.StateBucket}}"
    dynamodb_table = "{{.LockTable}}"
    key            = "{{.Key}}"
    region         = "{{.Region}}"
  }
}

`

const tfAzureRMBackendTemplate = `terraform {
  backend "azurerm" {
    resource_group_name  = "{{.ResourceGroupName}}"
    storage_account_name = "{{.StorageAccountName}}"
    container_name       = "{{.ContainerName}}"
    key                  = "{{.Key}}"
  }
}
`

const tfAzureRmStateTemplate = `data "terraform_remote_state" "{{.RemoteStateName}}" {
  backend = "azurerm"

  config = {
    resource_group_name  = "{{.ResourceGroupName}}"
    storage_account_name = "{{.StorageAccountName}}"
    container_name       = "{{.ContainerName}}"
    key                  = "{{.Key}}"
  }
}

`
