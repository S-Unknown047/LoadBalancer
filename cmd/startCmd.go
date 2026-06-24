package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	proxy "github.com/S-Unknown047/LoadBalancer/internal/proxy/l4"
	"github.com/S-Unknown047/LoadBalancer/internal/router"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var start = &cobra.Command{
	Use:   "To start Aap",
	Short: "This is to start the server and  to let the app know  which mode to use Nat-mode for l4 or l7",
	Long: `This is to start the  app and the server must be running in NAT mode or DNS mode : 
    
	start -s l7/l4
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
}

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

	fmt.Println("Mode: ", mode)

	if mode == "l4" {
		go proxy.Init()
	} else if mode == "l6" {
		go startServer()
	} else {
		fmt.Println("Invalid mode")
		return
	}

}

func startServer() {
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
}
