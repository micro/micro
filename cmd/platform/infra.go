package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/micro/platform/infra"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// infraCmd represents the infrastructure command
var infraCmd = &cobra.Command{
	Use:   "infra",
	Short: "Manage the platform's infrastructure'",
	Long: `Manage the platform's infra. Based on a configuration file,
a complete platform can be created across multiple cloud providers`,
}

func init() {
	cobra.OnInitialize(viperConfig)

	// TODO at some point: Resurrect the Infra command?
	// rootCmd.AddCommand(infraCmd)

	infraCmd.PersistentFlags().StringP(
		"config-file",
		"c",
		"",
		"Path to infrastructure definition file ($MICRO_CONFIG_FILE)",
	)
	viper.BindPFlag("config-file", infraCmd.PersistentFlags().Lookup("config-file"))
}

// viperConfig is run before every infra command, parsing config using viper
func viperConfig() {
	// Defaults - can be overwritten in the config file or env variables, but undocumented atm
	viper.SetDefault("state-store", "azure")
	// AWS Defaults
	viper.SetDefault("aws-region", "eu-west-2")
	viper.SetDefault("aws-s3-bucket", "micro-platform-terraform-state")
	viper.SetDefault("aws-dynamodb-table", "micro-platform-terraform-lock")
	// Azure defaults
	viper.SetDefault("azure-state-resource-group", "micro-terraform-states")
	viper.SetDefault("azure-storage-account", "microplatform")
	viper.SetDefault("azure-storage-container", "tfstate")

	// Handle env variables, e.g. --config-file flag can be set with MICRO_CONFIG_FILE
	viper.SetEnvPrefix("micro")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Read in config file
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Validate the configuration",
	Long: `Show what actions will be carried out to the platform

Instantiates various terraform modules, then runs terraform init, terraform validate`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, p := range validate() {
			s, err := p.Steps()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
			if err := infra.ExecutePlan(s); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
		}
		fmt.Printf("Plan Succeeded - run infra apply\n")
	},
}

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the configuration",
	Long: `Applies the configuration - this creates or modifies cloud resources

If you cancel this command, data loss may occur`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, p := range validate() {
			s, err := p.Steps()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
			if err := infra.ExecuteApply(s); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
		}
		fmt.Printf("Apply Succeeded\n")
	},
}

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the configuration",
	Long: `Destroys the configuration - this destroys or modifies cloud resources

If you cancel this command, data loss may occur`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, p := range validate() {
			s, err := p.Steps()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
			if err := infra.ExecuteDestroy(s); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
		}
		fmt.Printf("Destroy Succeeded\n")
	},
}

func validate() []infra.Platform {
	if viper.Get("platforms") == nil || len(viper.Get("platforms").([]interface{})) == 0 {
		fmt.Fprintf(os.Stderr, "No platforms defined in config file %s\n", viper.Get("config-file"))
		os.Exit(1)
	}
	var platforms []infra.Platform
	err := viper.UnmarshalKey("platforms", &platforms)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	return platforms
}

func init() {
	infraCmd.AddCommand(planCmd)
	infraCmd.AddCommand(applyCmd)
	infraCmd.AddCommand(destroyCmd)
}
