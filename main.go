package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rahulkj/iris-cli/services"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("subcommand is required")
		fmt.Println("- " + services.SERVICE_ACCOUNT_CMD)
		fmt.Println("- " + services.GLOBAL_DEFAULT_CMD)
		fmt.Println("- " + services.VCENTER_CMD)
		os.Exit(1)
	}

	switch strings.ToLower(os.Args[1]) {
	case strings.ToLower(services.SERVICE_ACCOUNT_CMD):
		sa := services.ServiceAccounts{}
		sa.Execute()
	case strings.ToLower(services.GLOBAL_DEFAULT_CMD):
		ga := services.GlobalDefaults{}
		ga.Execute()
	case strings.ToLower(services.VCENTER_CMD):
		vc := services.VCenters{}
		vc.Execute()
	default:
		fmt.Println("supported commands are:")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
