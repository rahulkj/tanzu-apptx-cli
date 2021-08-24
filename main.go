package main

import (
	"github.com/rahulkj/iris-cli/services"
)

func main() {
	sa := services.ServiceAccounts{}
	sa.Execute()
}
