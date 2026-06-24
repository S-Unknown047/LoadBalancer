package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	helper "github.com/S-Unknown047/LoadBalancer/internal/helper"
	model "github.com/S-Unknown047/LoadBalancer/internal/model"
	proxy "github.com/S-Unknown047/LoadBalancer/internal/proxy/l4"
	"github.com/S-Unknown047/LoadBalancer/internal/router"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var start = &cobra.Command{
	Use:   "start",
	Short: "This is to start the server and  to let the app know  which mode to use Nat-mode for l4 or l7",
	Long: `This is to start the  app and the server must be running in NAT mode or DNS mode : 
    
	start -s l4 -a rr -r 127.0.0.1:90:1 -r 127.0.0.1:100:2
	.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Executing the server logic directly inside Cobra's Run hook
		StartCmd(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(start)
	start.Flags().StringP("mode", "s", "l7", "Mode to use Nat-mode for l4 or l7")
	start.MarkFlagRequired("mode")

	start.Flags().StringSliceP("server", "r", []string{}, "Server details in format IP:PORT:WEIGHT")
	start.MarkFlagRequired("server")

	start.Flags().StringP("algo", "a", "", "Algorithm to use (rr/lc/RoundRobin/LeastConnection)")
	start.MarkFlagRequired("algo")
}

var grp sync.WaitGroup

func StartCmd(cmd *cobra.Command, args []string) {

	if !isRoot() {
		fmt.Println("Error: This command requires administrative privileges.")
		fmt.Println("Please run this application using 'sudo'.")
		os.Exit(1)
	}

	previliagedCmdRun()

	mode, err := cmd.Flags().GetString("mode")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse server details
	serverStrings, err := cmd.Flags().GetStringSlice("server")
	if err != nil {
		fmt.Println("Error reading server flag:", err)
		return
	}

	var servers []model.ReqServer
	for _, sStr := range serverStrings {
		parts := strings.Split(sStr, ":")
		if len(parts) != 3 {
			fmt.Printf("Error: Invalid server format '%s'. Expected IP:PORT:WEIGHT\n", sStr)
			os.Exit(1)
		}

		ip := parts[0]
		port := parts[1]
		weightStr := parts[2]

		weight, err := strconv.Atoi(weightStr)
		if err != nil {
			fmt.Printf("Error: Invalid weight '%s' for server %s:%s. Must be an integer.\n", weightStr, ip, port)
			os.Exit(1)
		}

		servers = append(servers, model.ReqServer{
			IP:     ip,
			Port:   port,
			Weight: uint64(weight),
		})
	}

	// Parse algorithm
	algorithm, err := cmd.Flags().GetString("algo")
	if err != nil {
		fmt.Println("Error reading algo flag:", err)
		return
	}

	if algorithm != "rr" && algorithm != "lc" && algorithm != "Round Robin" && algorithm != "Least Connection" && algorithm != "RoundRobin" && algorithm != "LeastConnection" {
		fmt.Printf("Error: Invalid algorithm '%s'. Allowed: rr, lc, RoundRobin, LeastConnection\n", algorithm)
		os.Exit(1)
	}

	// Store servers in memory first
	helper.HandelServer(&servers)

	// Setup backend config
	var backendObj model.ReqSetup
	backendObj.Algo = algorithm
	helper.HandelSetup(&backendObj)

	fmt.Println("Mode: ", mode)

	if mode == "l4" {
		grp.Add(1)
		go proxy.Init(&grp)
	} else if mode == "l7" {
		grp.Add(1)
		go startServer(&grp)
	} else {
		fmt.Println("Invalid mode")
		return
	}

	grp.Wait()
}

func startServer(grp *sync.WaitGroup) {
	defer grp.Done()
	godotenv.Load()
	server := router.Routing()
	port := os.Getenv("APP_PORT")
	fmt.Println("running on server ", port)
	log.Fatal(http.ListenAndServe(port, server))
}

func isRoot() bool {
	return os.Geteuid() == 0
}

func previliagedCmdRun() {
	godotenv.Load()

	translation_ip := os.Getenv("TRANSLATION_IP")
	vip := os.Getenv("VIP")

	exec.Command("sudo", "ip", "addr", "add", vip, "dev", "wlp0s20f3").Run()
	exec.Command("sudo", "iptables", "-A", "INPUT", "-p", "tcp", "-d", translation_ip, "--dport", "49152:65535", "-j", "DROP").Run()
	exec.Command("sudo", "iptables", "-A", "INPUT", "-p", "tcp", "-d", vip, "--dport", "80", "-j", "DROP").Run()
}
