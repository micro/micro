package main

import (
	"fmt"
	"os"

	"github.com/micro/platform/infra"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"math/rand"
)

var (
	kubeCommand = &cobra.Command{
		Use:   "kubernetes",
		Short: "Provision Kubernetes clusters",
		Long:  `Provision Kubernetes clusters`,
	}

	kubeCreateCommand = &cobra.Command{
		Use:   "create",
		Short: "Create a Kubernetes cluster",
		Long:  "Create a Kubernetes cluster",
		Run: func(cmd *cobra.Command, args []string) {
			k, err := makeKube()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			err = infra.ExecuteApply(k)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
		},
	}

	kubeDestroyCommand = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy a Kubernetes cluster",
		Long:  "Destroy a Kubernetes cluster",
		Run: func(cmd *cobra.Command, args []string) {
			k, err := makeKube()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			err = infra.ExecuteDestroy(k)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
		},
	}

	kubeConfigCommand = &cobra.Command{
		Use:   "get-config",
		Short: "Get Kube config for a created cluster",
		Long:  "Get Kube config for a created cluster",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := makeKubeConfig(viper.GetString("kube-config-path"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			err = infra.ExecuteApply(c)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
		},
	}
)

func makeKube() ([]infra.Step, error) {
	k := &infra.Kubernetes{
		Name:     viper.GetString("cluster-name"),
		Provider: viper.GetString("cloud-provider"),
		Region:   viper.GetString("cluster-region"),
	}
	return k.Steps(rand.Int31())
}

func makeKubeConfig(path string) ([]infra.Step, error) {
	k := &infra.Kubernetes{
		Name:     viper.GetString("cluster-name"),
		Provider: viper.GetString("cloud-provider"),
		Region:   viper.GetString("cluster-region"),
	}
	return k.Config(rand.Int31(), path)
}

func init() {
	rootCmd.AddCommand(kubeCommand)
	kubeCommand.AddCommand(kubeCreateCommand)
	kubeCommand.AddCommand(kubeDestroyCommand)
	kubeCommand.AddCommand(kubeConfigCommand)
	kubeCommand.PersistentFlags().StringP("name", "n", "microdev", "Cluster name")
	viper.BindPFlag("cluster-name", kubeCommand.PersistentFlags().Lookup("name"))
	kubeCommand.PersistentFlags().StringP("region", "r", "westeurope", "Cluster Region")
	viper.BindPFlag("cluster-region", kubeCommand.PersistentFlags().Lookup("region"))
}
