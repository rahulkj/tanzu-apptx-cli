package main

import (
	"fmt"
	"os"
	"strings"

	"gitlab.eng.vmware.com/vmware-navigator-practice/tooling/iris-cli/services"
)

func main() {

	if len(os.Args) < 2 {
		printUsage()
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
	case strings.ToLower(services.VRNI_CMD):
		vr := services.VRNI{}
		vr.Execute()
	case strings.ToLower(services.VIRTUAL_MACHINES):
		vm := services.VirtualMachines{}
		vm.Execute()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: 'iris-cli [command]' \n")

	fmt.Println("Available Commands:")
	fmt.Printf("  %s \t\t%s \n", services.SERVICE_ACCOUNT_CMD, "Service Accounts operations")
	fmt.Printf("  %s \t\t%s \n", services.GLOBAL_DEFAULT_CMD, "Global Defaults operations")
	fmt.Printf("  %s \t\t\t%s \n", services.VCENTER_CMD, "vCenter operations")
	fmt.Printf("  %s \t\t\t\t%s \n", services.VRNI_CMD, "vRNI operations")
	fmt.Printf("  %s \t\t%s \n", services.VIRTUAL_MACHINES, "Virtual Machines operations")
	os.Exit(1)
}
