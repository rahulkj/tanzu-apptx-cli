package services

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type VCenters struct {
	url       string
	username  string
	password  string
	saAlias   string
	vcFqdn    string
	vcName    string
	operation string
}

func (vCenters VCenters) Execute() {
	vCenters = vCenters.validate()

	request := Request{vCenters.url, vCenters.username, vCenters.password}
	authResponse := Authenticate(request)

	switch vCenters.operation {
	case REGISTER:
		vCenters.register(authResponse.Token, request)
	case UNREGISTER:
		vCenters.unregister(authResponse.Token, request)
	case SYNC_VCENTERS:
		vCenters.syncVcenters(authResponse.Token, request)
	case SCAN_VIRTUAL_MACHINES:
		vCenters.scanVirtualMachines(authResponse.Token, request)
	case SCAN_COMPONENTS:
		vCenters.scanComponents(authResponse.Token, request)
	case DISCOVER_TOPOLOGY:
		vCenters.discoverTopology(authResponse.Token, request)
	default:
		fmt.Println("Operation not supported \n")
		vCenters.printUsage()
		os.Exit(1)
	}
}

func (vCenters VCenters) printUsage() {
	fmt.Printf("Usage: '%s %s [command]' \n", CLI_NAME, VCENTER_CMD)
	fmt.Println("Available Commands:")
	fmt.Printf("  %s \t\t\t%s \n", REGISTER, "Register vCenter instance")
	fmt.Printf("  %s \t\t\t%s \n", UNREGISTER, "Remove vCenter instance")
	fmt.Printf("  %s \t\t\t\t%s \n", SYNC_VCENTERS, "Sync vCenter inventory")
	fmt.Printf("  %s \t%s \n", SCAN_VIRTUAL_MACHINES, "Scan for virtual machines managed by a vCenter")
	fmt.Printf("  %s \t\t%s \n", SCAN_COMPONENTS, "Scan for components running on the virtual machines managed by a vCenter")
	fmt.Printf("  %s \t\t%s \n", DISCOVER_TOPOLOGY, "Discover topology for the components running on the virtual machines managed by a vCenter")
	os.Exit(1)
}

func (vCenters VCenters) validate() VCenters {
	registerCmd := flag.NewFlagSet(REGISTER, flag.ExitOnError)
	unregisterCmd := flag.NewFlagSet(UNREGISTER, flag.ExitOnError)
	syncVCenterCmd := flag.NewFlagSet(SYNC_VCENTERS, flag.ExitOnError)
	scanVirtualMachinesCmd := flag.NewFlagSet(SCAN_VIRTUAL_MACHINES, flag.ExitOnError)
	scanComponentsCmd := flag.NewFlagSet(SCAN_COMPONENTS, flag.ExitOnError)
	discoverTopologyCmd := flag.NewFlagSet(DISCOVER_TOPOLOGY, flag.ExitOnError)

	if len(os.Args) < 3 {
		vCenters.printUsage()
	}

	operation := os.Args[2]

	var url string
	var username string
	var password string
	var vcFqdn string
	var vcName string
	var saAlias string

	if operation == REGISTER {
		registerCmd.StringVar(&url, "url", "", "Application Transformer URL, ex: appliance.example.com")
		registerCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		registerCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		registerCmd.StringVar(&vcFqdn, "vc-fqdn", "", "vCenter FQDN")
		registerCmd.StringVar(&vcName, "vc-name", "", "vCenter Name")
		registerCmd.StringVar(&saAlias, "sa-alias", "", "service account alias")

		registerCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(vcFqdn) == 0 || len(vcName) == 0 || len(saAlias) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, VCENTER_CMD, REGISTER)
			fmt.Println("Available Flags:")
			registerCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == UNREGISTER {
		unregisterCmd.StringVar(&url, "url", "", "Application Transformer URL, ex: appliance.example.com")
		unregisterCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		unregisterCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		unregisterCmd.StringVar(&vcFqdn, "vc-fqdn", "", "vCenter FQDN")
		unregisterCmd.StringVar(&vcName, "vc-name", "", "vCenter Name")
		// saAlias = new(string)

		unregisterCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(vcFqdn) == 0 && len(vcName) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, VCENTER_CMD, UNREGISTER)
			fmt.Println("Available Flags:")
			unregisterCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == SYNC_VCENTERS {
		syncVCenterCmd.StringVar(&url, "url", "", "Application Transformer URL, ex: appliance.example.com")
		syncVCenterCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		syncVCenterCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		syncVCenterCmd.StringVar(&vcFqdn, "vc-fqdn", "", "vCenter FQDN")
		syncVCenterCmd.StringVar(&vcName, "vc-name", "", "vCenter Name")

		syncVCenterCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(vcFqdn) == 0 && len(vcName) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, VCENTER_CMD, SYNC_VCENTERS)
			fmt.Println("Available Flags:")
			syncVCenterCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == SCAN_VIRTUAL_MACHINES {
		scanVirtualMachinesCmd.StringVar(&url, "url", "", "Application Transformer URL, ex: appliance.example.com")
		scanVirtualMachinesCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		scanVirtualMachinesCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		scanVirtualMachinesCmd.StringVar(&vcFqdn, "vc-fqdn", "", "vCenter FQDN")
		scanVirtualMachinesCmd.StringVar(&vcName, "vc-name", "", "vCenter Name")

		scanVirtualMachinesCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(vcFqdn) == 0 && len(vcName) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, VCENTER_CMD, SCAN_VIRTUAL_MACHINES)
			fmt.Println("Available Flags:")
			scanVirtualMachinesCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == SCAN_COMPONENTS {
		scanComponentsCmd.StringVar(&url, "url", "", "Application Transformer URL, ex: appliance.example.com")
		scanComponentsCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		scanComponentsCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		scanComponentsCmd.StringVar(&vcFqdn, "vc-fqdn", "", "vCenter FQDN")
		scanComponentsCmd.StringVar(&vcName, "vc-name", "", "vCenter Name")

		scanComponentsCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(vcFqdn) == 0 && len(vcName) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, VCENTER_CMD, SCAN_COMPONENTS)
			fmt.Println("Available Flags:")
			scanComponentsCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == DISCOVER_TOPOLOGY {
		discoverTopologyCmd.StringVar(&url, "url", "", "Application Transformer URL, ex: appliance.example.com")
		discoverTopologyCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		discoverTopologyCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		discoverTopologyCmd.StringVar(&vcFqdn, "vc-fqdn", "", "vCenter FQDN")
		discoverTopologyCmd.StringVar(&vcName, "vc-name", "", "vCenter Name")

		discoverTopologyCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(vcFqdn) == 0 && len(vcName) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, VCENTER_CMD, DISCOVER_TOPOLOGY)
			fmt.Println("Available Flags:")
			discoverTopologyCmd.PrintDefaults()
			os.Exit(1)
		}
	} else {
		vCenters.printUsage()
	}

	vCenters = VCenters{url, username, password, saAlias, vcFqdn, vcName, operation}
	return vCenters
}

func (vCenters VCenters) register(token string, request Request) {
	serviceAccounts := ServiceAccounts{}
	response := serviceAccounts.findServiceAccount(vCenters.saAlias, token, request)

	if len(response.Embedded.ServiceAccounts) > 0 {

		for _, serviceAccount := range response.Embedded.ServiceAccounts {
			if serviceAccount.Alias == vCenters.saAlias {
				url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VCENTERS + "?action=register"

				vcRequest := VCenterRequest{vCenters.vcFqdn, vCenters.vcName, serviceAccount.UUID}
				body, responseCode := processRequest(token, url, "POST", vcRequest)

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
						log.Println("Failed to register vCenter with the provided information")
					} else {
						log.Println("Successfully registered vCenter with the provided information")
					}
				} else {
					log.Println("Failed to register vCenter with the provided information. Response Code:", responseCode)
				}
			} else {
				log.Println("Cannot complete the operation as the Service Account does not exist")
			}
		}

	} else {
		log.Println("Cannot complete the operation as the Service Account does not exist")
	}
}

func (vCenters VCenters) unregister(token string, request Request) {
	vCenterUUID := vCenters.findVCenter(token, request).VCenterUUID

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VCENTERS + "/" + vCenterUUID
	_, responseCode := processRequest(token, url, "DELETE", nil)

	if responseCode == 200 {
		log.Println("Successfully deleted vCenter")
	} else {
		log.Println("Failed to delete vCenter. Response Code: ", responseCode)
	}
}

func (vCenters VCenters) findVCenter(token string, request Request) (response VCenter) {

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VCENTERS + "?"

	if len(vCenters.vcName) > 0 {
		url = url + "&vcName=" + vCenters.vcName
	}

	if len(vCenters.vcFqdn) > 0 {
		url = url + "&fqdn=" + vCenters.vcFqdn
	}

	body, _ := processRequest(token, url, "GET", nil)

	vCenterResponse := VCenterListResponse{}
	err := json.Unmarshal(body, &vCenterResponse)
	if err != nil {
		log.Println("Failed to parse the response body.\n[ERROR] -", err)
		os.Exit(1)
	}

	if len(vCenterResponse.Embedded.VCenters) > 0 {
		for _, vCenter := range vCenterResponse.Embedded.VCenters {
			if (vCenter.VCName == vCenters.vcName) || (vCenter.Fqdn == vCenters.vcFqdn) {
				response = vCenter
			} else {
				log.Println("Cannot find the vCenter as it does not exist")
			}
		}
	} else {
		log.Println("Could not find the vCenter")
	}

	return response
}

func (vCenters VCenters) findAll(token string, request Request, vCentersCSV string) (vCenterUUIDs []string) {

	vCentersArray := strings.Split(vCentersCSV, ",")

	vCentersUUIDsArray := []string{}

	for _, vcenter := range vCentersArray {
		vCenter := VCenters{vcName: vcenter}
		vCenterResponse := vCenter.findVCenter(token, request)
		vCentersUUIDsArray = append(vCentersUUIDsArray, vCenterResponse.VCenterUUID)
	}

	return vCentersUUIDsArray
}

func (vCenters VCenters) syncVcenters(token string, request Request) {
	vCenter := vCenters.findVCenter(token, request)

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VCENTERS + "/" + vCenter.VCenterUUID + "/sync"

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

func (vCenters VCenters) scanVirtualMachines(token string, request Request) {
	vCenter := vCenters.findVCenter(token, request)

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VCENTERS + "/" + vCenter.VCenterUUID + "/virtualmachines"

	dataCenter := new(Datacenter)
	vcRequest := VCenterScanVMRequest{false, false, false, *dataCenter}

	body, responseCode := processRequest(token, url, "POST", vcRequest)

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
			log.Println("Failed to scan virtual machines managed by the provided vCenter")
		} else {
			log.Println("Successfully scanned virtual machines managed by the provided vCenter")
		}
	} else {
		log.Println("Failed to scan virtual machines managed by the provided vCenter. Response Code:", responseCode)
	}
}

func (vCenters VCenters) scanComponents(token string, request Request) {
	vCenter := vCenters.findVCenter(token, request)

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VCENTERS + "/" + vCenter.VCenterUUID + "/components"

	dataCenter := new(Datacenter)
	vcRequest := VCenterScanVMRequest{true, false, false, *dataCenter}

	body, responseCode := processRequest(token, url, "POST", vcRequest)

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
			log.Println("Failed to scan components running on the virtual machines managed by the provided vCenter")
		} else {
			log.Println("Successfully scanned components running on the virtual machines managed by the provided vCenter")
		}
	} else {
		log.Println("Failed to scan components running on the virtual machines managed by the provided vCenter. Response Code:", responseCode)
	}
}

func (vCenters VCenters) discoverTopology(token string, request Request) {
	vCenter := vCenters.findVCenter(token, request)

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VCENTERS + "/" + vCenter.VCenterUUID + "/correlation"

	dataCenter := new(Datacenter)
	discoverTopologyRequest := DiscoverTopologyRequest{*dataCenter}

	body, responseCode := processRequest(token, url, "POST", discoverTopologyRequest)

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
			log.Println("Failed to discover topology for the provided vCenter")
		} else {
			log.Println("Successfully discovered topology for the provided vCenter")
		}
	} else {
		log.Println("Failed to discover topology for the provided vCenter. Response Code:", responseCode)
	}
}
