package cmd

import (
	"fmt"

	helper "github.com/S-Unknown047/LoadBalancer/internal/helper"
	model "github.com/S-Unknown047/LoadBalancer/internal/model"
	"github.com/spf13/cobra"
)

var setup = &cobra.Command{
	Use:   "setup for Backend",
	Short: "set up for the backend with Algorithm",
	Long: `This is setup command for the backend.
	
Example:
  setup setup -a/algo "rr(round robin)/lc(for least connection)"`,
	Run: SetupCommand,
}

func init() {
	rootCmd.AddCommand(setup)

	// StringSlice allows the user to provide the --server (or -s) flag multiple times
	setup.Flags().StringP("algo", "a", "", "Server details in format algo -a lc/rr")

	// Mark the flag as required so the user doesn't pass empty values
	setup.MarkFlagRequired("algo")
}

func SetupCommand(cmd *cobra.Command, args []string) {
	algorithm, err := cmd.Flags().GetString("algo")

	if err != nil {
		fmt.Println(err)
		return
	}

	if algorithm != "rr" && algorithm != "lc" && algorithm != "Round Robin" && algorithm != "Least Connection" {
		fmt.Println("Invalid algorithm")
		return
	}

	var backendObj model.ReqSetup
	backendObj.Algo = algorithm
	helper.HandelSetup(&backendObj)
}
