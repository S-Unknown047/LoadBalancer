package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/S-Unknown047/LoadBalancer/internal/helper"
	model "github.com/S-Unknown047/LoadBalancer/internal/model"
	"github.com/spf13/cobra"
)

var addServer = &cobra.Command{
	Use:   "add",
	Short: "Add one or multiple servers to the load balancer",
	Long: `Add servers by providing a list of servers in the format IP:PORT:WEIGHT.
	
Example:
  addServer --server "192.168.1.1:8080:5" --server "192.168.1.2:9090:10"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Executing the server logic directly inside Cobra's Run hook
		Server(cmd)
	},
}

func init() {
	rootCmd.AddCommand(addServer)

	// StringSlice allows the user to provide the --server (or -s) flag multiple times
	addServer.Flags().StringSliceP("server", "s", []string{}, "Server details in format IP:PORT:WEIGHT")

	// Mark the flag as required so the user doesn't pass empty values
	addServer.MarkFlagRequired("server")
}

func Server(cmd *cobra.Command) {
	serverStrings, err := cmd.Flags().GetStringSlice("server")
	if err != nil {
		fmt.Println(err)
		return
	}
	var servers []model.ReqServer
	for _, sStr := range serverStrings {
		parts := strings.Split(sStr, ":")
		if len(parts) != 3 {
			fmt.Printf("Invalid server format '%s'. Expected IP:PORT:WEIGHT\n", sStr)
			continue
		}

		ip := parts[0]
		port := parts[1]
		weightStr := parts[2]

		weight, err := strconv.Atoi(weightStr)
		if err != nil {
			fmt.Printf("Invalid weight '%s' for server %s:%s. Must be an integer.\n", weightStr, ip, port)
			continue
		}

		servers = append(servers, model.ReqServer{
			IP:     ip,
			Port:   port,
			Weight: uint64(weight),
		})
	}

	helper.HandelServer(&servers)
}
