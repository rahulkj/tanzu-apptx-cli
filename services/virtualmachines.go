package services

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

type VirtualMachines struct {
	url            string
	username       string
	password       string
	vcFqdn         string
	vcDatacenter   string
	vcCluster      string
	vcResourcePool string
	vcFolder       string
	vmName         string
	vmIP           string
	outputFormat   string
	operation      string
}

type VirtualMachinesResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Network      string   `json:"network"`
	Hostname     string   `json:"hostname"`
	Datastore    string   `json:"datastore"`
	IP           string   `json:"ip"`
	NumCPU       int      `json:"numCPU"`
	MemoryMB     string   `json:"memoryMB"`
	Services     []string `json:"services"`
	VcenterFqdn  string   `json:"vcenterFqdn"`
	DataCenter   string   `json:"dataCenter"`
	Cluster      string   `json:"cluster"`
	ResourcePool string   `json:"resourcePool'`
	Folder       string   `json:"folder"`
	NumOfDisks   int      `json:"numOfDisks"`
	SizeOfDisks  string   `json:"sizeOfDisks"`
}

type VirtualMachinesListResponse struct {
	Embedded struct {
		VirtualMachinesResponse []VirtualMachinesResponse `json:"virtualmachines"`
	} `json:"_embedded"`
}

func (virtualMachines VirtualMachines) Execute() {
	virtualMachines = virtualMachines.validate()

	request := Request{virtualMachines.url, virtualMachines.username, virtualMachines.password}
	authResponse := Authenticate(request)

	switch virtualMachines.operation {
	case LIST:
		virtualMachinesList := virtualMachines.list(authResponse.Token)

		if virtualMachines.outputFormat == "table" {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
			fmt.Fprintln(w, "VM ID\tNAME\tvCenter\tDataCenter\tCluster\tResource Pool\tFolder\tNetwork\tDatastore\tIP\tCPU\tMemory (in MB)\tDisk Size\tServices")
			for _, virtualMachine := range virtualMachinesList.Embedded.VirtualMachinesResponse {
				fmt.Fprintln(w, virtualMachine.ID, "\t", virtualMachine.Name, "\t", virtualMachine.VcenterFqdn,
					"\t", virtualMachine.DataCenter, "\t", virtualMachine.Cluster, "\t", virtualMachine.ResourcePool,
					"\t", virtualMachine.Folder, "\t", virtualMachine.Network, "\t", virtualMachine.Datastore,
					"\t", virtualMachine.IP, "\t", virtualMachine.NumCPU, "\t", virtualMachine.MemoryMB,
					"\t", virtualMachine.SizeOfDisks, "\t", strings.Join(virtualMachine.Services, ","))
			}
			w.Flush()
		} else if virtualMachines.outputFormat == "json" {
			prettyJSON, err := json.MarshalIndent(virtualMachinesList, "", "    ")
			if err != nil {
				log.Fatal("Failed to generate json", err)
			}
			fmt.Printf("%s\n", string(prettyJSON))
		} else if virtualMachines.outputFormat == "csv" {
			fmt.Println("VM ID,NAME,vCenter,DataCenter,Cluster,Resource Pool,Folder,Network,Datastore,IP,CPU,Memory (in MB),Disk Size,Services")
			for _, virtualMachine := range virtualMachinesList.Embedded.VirtualMachinesResponse {
				fmt.Println(virtualMachine.ID, ",", virtualMachine.Name, ",", virtualMachine.VcenterFqdn,
					",", virtualMachine.DataCenter, ",", virtualMachine.Cluster, ",", virtualMachine.ResourcePool,
					",", virtualMachine.Folder, ",", virtualMachine.Network, ",", virtualMachine.Datastore,
					",", virtualMachine.IP, ",", virtualMachine.NumCPU, ",", virtualMachine.MemoryMB,
					",", virtualMachine.SizeOfDisks, ",", virtualMachine.Services)
			}
		}

	case INTROSPECT:
		virtualMachines.introspect(authResponse.Token, request)
	default:
		fmt.Println("Operation not supported")
		virtualMachines.printUsage()
		os.Exit(1)
	}
}

func (virtualMachines VirtualMachines) validate() VirtualMachines {
	listCmd := flag.NewFlagSet(LIST, flag.ExitOnError)
	introspectCmd := flag.NewFlagSet(INTROSPECT, flag.ExitOnError)

	if len(os.Args) < 3 {
		virtualMachines.printUsage()
	}

	operation := os.Args[2]

	var url string
	var username string
	var password string
	var vcFqdn string
	var vcDatacenter string
	var vcCluster string
	var vcResourcePool string
	var vcFolder string
	var vmName string
	var vmIP string
	var format string

	if operation == LIST {
		listCmd.StringVar(&url, "url", "", "Iris URL, ex: appliance.example.com")
		listCmd.StringVar(&username, "username", "", "Iris admin username")
		listCmd.StringVar(&password, "password", "", "Iris admin password")
		listCmd.StringVar(&vcFqdn, "vc-fqdn", "", "vCenter FQDN")
		listCmd.StringVar(&vcDatacenter, "vc-datacenter", "", "vCenter Datacenter")
		listCmd.StringVar(&vcCluster, "vc-cluster", "", "vCenter Cluster Name")
		listCmd.StringVar(&vcResourcePool, "vc-resource-pool", "", "vCenter Resource Pool Name")
		listCmd.StringVar(&vcFolder, "vc-folder", "", "vCenter Folder Name")
		listCmd.StringVar(&vmName, "vm-name", "", "Virtual Machine Name")
		listCmd.StringVar(&vmIP, "vm-ip", "", "Virtual Machine IP")
		listCmd.StringVar(&format, "output-format", "table", "Output format - (json,csv,table) (Default: table)")

		listCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Println("Usage: 'iris-cli virtual-machines list [flags]' \n")
			fmt.Println("Flags:")
			listCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == INTROSPECT {
		introspectCmd.StringVar(&url, "url", "", "Iris URL, ex: appliance.example.com")
		introspectCmd.StringVar(&username, "username", "", "Iris admin username")
		introspectCmd.StringVar(&password, "password", "", "Iris admin password")
		introspectCmd.StringVar(&vmName, "vm-name", "", "Virtual Machine Name")

		introspectCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			len(vmName) == 0 ||
			(strings.Contains(url, "https://")) {
			fmt.Println("Usage: 'iris-cli virtual-machines get [flags]' \n")
			fmt.Println("Flags:")
			introspectCmd.PrintDefaults()
			os.Exit(1)
		}
	} else {
		virtualMachines.printUsage()
	}

	virtualMachines = VirtualMachines{url, username, password, vcFqdn, vcDatacenter, vcCluster, vcResourcePool, vcFolder, vmName, vmIP, format, operation}
	return virtualMachines
}

func (virtualMachines VirtualMachines) printUsage() {
	fmt.Println("Usage: 'iris-cli virtual-machines [command]' \n")
	fmt.Println("Available Commands:")
	fmt.Printf("  %s \t\t\t%s \n", LIST, "List all virtual machines")
	fmt.Printf("  %s \t\t%s \n", INTROSPECT, "Introspect a virtual machine")
	os.Exit(1)
}

func (virtualMachines VirtualMachines) list(token string) (response VirtualMachinesListResponse) {
	url := PROTOCOL + "://" + virtualMachines.url + "/" + PREFIX + "/" + VIRTUAL_MACHINES + "?"

	if len(virtualMachines.vcFqdn) > 0 {
		url = url + "&vcenterFqdn=" + virtualMachines.vcFqdn
	}

	if len(virtualMachines.vcDatacenter) > 0 {
		url = url + "&dataCenter=" + virtualMachines.vcDatacenter
	}

	if len(virtualMachines.vcFolder) > 0 {
		url = url + "&folder=" + virtualMachines.vcFolder
	}

	if len(virtualMachines.vcCluster) > 0 {
		url = url + "&cluster=" + virtualMachines.vcCluster
	}

	if len(virtualMachines.vcResourcePool) > 0 {
		url = url + "&resourcePool=" + virtualMachines.vcResourcePool
	}

	if len(virtualMachines.vmName) > 0 {
		url = url + "&name=" + virtualMachines.vmName
	}

	if len(virtualMachines.vmIP) > 0 {
		url = url + "&ip=" + virtualMachines.vmIP
	}

	body, responseCode := processRequest(token, url, "GET", nil)

	if responseCode == 200 {
		log.Println("Successfully fetched the list of virtual machines \n")
	} else {
		log.Println("Failed to fetch the list of virtual machines. Response code:", responseCode)
	}

	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Failed to parse the response body.\n[ERROR] -", err)
		os.Exit(1)
	}

	return response
}

func (virtualMachines VirtualMachines) introspect(token string, request Request) {
	virtualMachinesListResponse := virtualMachines.list(token)

	for _, virtualMachine := range virtualMachinesListResponse.Embedded.VirtualMachinesResponse {
		url := PROTOCOL + "://" + virtualMachines.url + "/" + PREFIX + "/" + VIRTUAL_MACHINES + "/" + virtualMachine.ID + "/components"
		body, responseCode := processRequest(token, url, "POST", nil)

		if responseCode == 202 {
			tasks := Tasks{}
			err := json.Unmarshal(body, &tasks)
			if err != nil {
				log.Println("Failed to parse the response body.\n[ERROR] -", err)
				os.Exit(1)
			}

			log.Println("Submitted the request and the taskID is:", tasks.TaskID)

			status := tasks.MonitorTask(token, tasks.TaskID, request)
			if status != "SUCCESS" {
				log.Println("Failed to execute sync on the vCenter provided")
			} else {
				log.Println("Successfully executed sync on the vCenter provided")
			}
		} else {
			log.Println("Failed to execute sync on the vCenter provided. Response Code:", responseCode)
		}
	}
}
