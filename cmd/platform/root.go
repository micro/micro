package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "platform",
	Short: "The Micro platform binary",
	Long: `The Micro platform binary.

All features of the micro platform can be started with this command.`,
}

// Global Flags
func init() {
	rootCmd.PersistentFlags().StringP("cloud-provider", "p", "azure", "Cloud provider (azure, do, local)")
	viper.BindPFlag("cloud-provider", rootCmd.PersistentFlags().Lookup("cloud-provider"))
	dir, err := homedir.Dir()
	if err != nil {
		dir = ""
	}
	rootCmd.PersistentFlags().StringP("kubeconfig", "k", dir+"/.kube/config", "Path to Kube Config")
	viper.BindPFlag("kube-config-path", rootCmd.PersistentFlags().Lookup("kubeconfig"))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
