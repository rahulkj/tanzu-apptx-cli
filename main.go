package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rahulkj/iris-cli/services"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("subcommand is required")
		fmt.Println("- " + services.SERVICE_ACCOUNT_CMD)
		fmt.Println("- " + services.GLOBAL_DEFAULT_CMD)
		os.Exit(1)
	}

	switch os.Args[1] {
	case services.SERVICE_ACCOUNT_CMD:
		sa := services.ServiceAccounts{}
		sa.Execute()
	case services.GLOBAL_DEFAULT_CMD:
		ga := services.GlobalDefaults{}
		ga.Execute()
	default:
		fmt.Println("supported commands are:")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
