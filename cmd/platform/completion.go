package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completion scripts",
	Long: `To output completion scripts, use:
	
./platform completion <shell>

For example, on bash, use:

source <(platform completion bash)`,
}

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion",
	Long: `Generates GNU Bash completion

To use, run:

source <(platform completion bash)`,
	Run: func(cmd *cobra.Command, args []string) { rootCmd.GenBashCompletion(os.Stdout) },
}

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generates zsh completion",
	Long: `Generates Z shell completion

To use, run:

source <(kubectl completion zsh)`,
	Run: func(cmd *cobra.Command, args []string) { rootCmd.GenZshCompletion(os.Stdout) },
}

var completionPwshCmd = &cobra.Command{
	Use:   "pwsh",
	Short: "Generates pwsh completion",
	Long:  "Generates Powershell completion",
	Run:   func(cmd *cobra.Command, args []string) { rootCmd.GenPowerShellCompletion(os.Stdout) },
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionZshCmd)
	completionCmd.AddCommand(completionPwshCmd)
}
