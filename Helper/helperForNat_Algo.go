package helper

import (
	"fmt"

	model "github.com/S-Unknown047/LoadBalancer/Model"
	natmode "github.com/S-Unknown047/LoadBalancer/NAT_MODE"
)

func startScan(b *model.Backend) {
	if b.Mode == "L4" {
		fmt.Println("L4 mode running")
		natmode.Test()
	}
}
