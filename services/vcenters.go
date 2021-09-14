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
	sa_alias  string
	vc_fqdn   string
	vc_name   string
	operation string
}

type VCenterRequest struct {
	Fqdn                 string `json:"fqdn"`
	VCName               string `json:"vcName"`
	VCServiceAccountUUID string `json:"vcServiceAccountUUID"`
}

type VCenter struct {
	Fqdn        string       `json:"fqdn"`
	VCenterUUID string       `json:"irisVcenterUUID"`
	VCName      string       `json:"vcName"`
	Datacenters []Datacenter `json:"dataCenters"`
}

type Datacenter struct {
	ModId    string `json:"modId"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Clusters []struct {
		ModId         string `json:"modId"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		ResourcePools []struct {
			ModId string `json:"modId"`
			Name  string `json:"name"`
			Type  string `json:"type"`
		} `json:"resourcePools"`
	} `json:"clusters"`
	Folders []struct {
		ModId string `json:"modId"`
		Name  string `json:"name"`
		Type  string `json:"type"`
	} `json:"folders"`
}

type VCenterListResponse struct {
	Embedded struct {
		VCenters []VCenter `json:"vcenters"`
	} `json:"_embedded"`
}

type VCenterScanVMRequest struct {
	ComponentScan         bool       `json:"componentScan"`
	ApplyCredentialPolicy bool       `json:"applyCredentialPolicy"`
	BinaryAnalysis        bool       `json:"binaryAnalysis"`
	Filters               Datacenter `json:"filters"`
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
	default:
		fmt.Println("Operation not supported \n")
		vCenters.printUsage()
		os.Exit(1)
	}
}

func (vCenters VCenters) printUsage() {
	fmt.Println("Usage: 'iris-cli vCenter [command]' \n")
	fmt.Println("Available Commands:")
	fmt.Printf("  %s \t\t\t%s \n", REGISTER, "Register vCenter instance")
	fmt.Printf("  %s \t\t\t%s \n", UNREGISTER, "Remove vCenter instance")
	fmt.Printf("  %s \t\t\t\t%s \n", SYNC_VCENTERS, "Sync vCenter inventory")
	fmt.Printf("  %s \t%s \n", SCAN_VIRTUAL_MACHINES, "Scan for virtual machines managed by a vCenter")
	fmt.Printf("  %s \t\t%s \n", SCAN_COMPONENTS, "Scan for components running on the virtual machines managed by a vCenter")
	os.Exit(1)
}

func (vCenters VCenters) validate() VCenters {
	registerCmd := flag.NewFlagSet(REGISTER, flag.ExitOnError)
	unregisterCmd := flag.NewFlagSet(UNREGISTER, flag.ExitOnError)
	syncVCenterCmd := flag.NewFlagSet(SYNC_VCENTERS, flag.ExitOnError)
	scanVirtualMachinesCmd := flag.NewFlagSet(SCAN_VIRTUAL_MACHINES, flag.ExitOnError)
	scanComponentsCmd := flag.NewFlagSet(SCAN_COMPONENTS, flag.ExitOnError)

	if len(os.Args) < 3 {
		vCenters.printUsage()
	}

	operation := os.Args[2]

	var url *string
	var username *string
	var password *string
	var vc_fqdn *string
	var vc_name *string
	var sa_alias *string

	if operation == REGISTER {
		url = registerCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = registerCmd.String("username", "", "Iris admin username")
		password = registerCmd.String("password", "", "Iris admin password")
		vc_fqdn = registerCmd.String("vc-fqdn", "", "vCenter FQDN")
		vc_name = registerCmd.String("vc-name", "", "vCenter Name")
		sa_alias = registerCmd.String("sa-alias", "", "service account alias")

		registerCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*vc_fqdn) == 0 || len(*vc_name) == 0 || len(*sa_alias) == 0) ||
			(strings.Contains(*url, "https://")) {
			fmt.Println("Usage: 'iris-cli vCenter register [flags]' \n")
			fmt.Println("Flags:")
			registerCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == UNREGISTER {
		url = unregisterCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = unregisterCmd.String("username", "", "Iris admin username")
		password = unregisterCmd.String("password", "", "Iris admin password")
		vc_fqdn = unregisterCmd.String("vc-fqdn", "", "vCenter FQDN")
		vc_name = unregisterCmd.String("vc-name", "", "vCenter Name")
		sa_alias = new(string)

		unregisterCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*vc_fqdn) == 0 && len(*vc_name) == 0) ||
			(strings.Contains(*url, "https://")) {
			fmt.Println("Usage: 'iris-cli vCenter unregister [flags]' \n")
			fmt.Println("Flags:")
			unregisterCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == SYNC_VCENTERS {
		url = syncVCenterCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = syncVCenterCmd.String("username", "", "Iris admin username")
		password = syncVCenterCmd.String("password", "", "Iris admin password")
		vc_fqdn = syncVCenterCmd.String("vc-fqdn", "", "vCenter FQDN")
		vc_name = syncVCenterCmd.String("vc-name", "", "vCenter Name")
		sa_alias = new(string)

		syncVCenterCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*vc_fqdn) == 0 && len(*vc_name) == 0) ||
			(strings.Contains(*url, "https://")) {
			fmt.Println("Usage: 'iris-cli vCenter sync [flags]' \n")
			fmt.Println("Flags:")
			syncVCenterCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == SCAN_VIRTUAL_MACHINES {
		url = scanVirtualMachinesCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = scanVirtualMachinesCmd.String("username", "", "Iris admin username")
		password = scanVirtualMachinesCmd.String("password", "", "Iris admin password")
		vc_fqdn = scanVirtualMachinesCmd.String("vc-fqdn", "", "vCenter FQDN")
		vc_name = scanVirtualMachinesCmd.String("vc-name", "", "vCenter Name")
		sa_alias = new(string)

		scanVirtualMachinesCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*vc_fqdn) == 0 && len(*vc_name) == 0) ||
			(strings.Contains(*url, "https://")) {
			fmt.Println("Usage: 'iris-cli vCenter scan-virtual-machines [flags]' \n")
			fmt.Println("Flags:")
			scanVirtualMachinesCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == SCAN_COMPONENTS {
		url = scanComponentsCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = scanComponentsCmd.String("username", "", "Iris admin username")
		password = scanComponentsCmd.String("password", "", "Iris admin password")
		vc_fqdn = scanComponentsCmd.String("vc-fqdn", "", "vCenter FQDN")
		vc_name = scanComponentsCmd.String("vc-name", "", "vCenter Name")
		sa_alias = new(string)

		scanComponentsCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*vc_fqdn) == 0 && len(*vc_name) == 0) ||
			(strings.Contains(*url, "https://")) {
			fmt.Println("Usage: 'iris-cli vCenter scan-components [flags]' \n")
			fmt.Println("Flags:")
			scanComponentsCmd.PrintDefaults()
			os.Exit(1)
		}
	} else {
		vCenters.printUsage()
	}

	vCenters = VCenters{*url, *username, *password, *sa_alias, *vc_fqdn, *vc_name, operation}
	return vCenters
}

func (vCenters VCenters) register(token string, request Request) {
	serviceAccounts := ServiceAccounts{}
	response := serviceAccounts.findServiceAccount(vCenters.sa_alias, token, request)

	if len(response.Embedded.ServiceAccounts) > 0 {

		for _, serviceAccount := range response.Embedded.ServiceAccounts {
			if serviceAccount.Alias == vCenters.sa_alias {
				url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VCENTERS + "?action=register"

				vcRequest := VCenterRequest{vCenters.vc_fqdn, vCenters.vc_name, serviceAccount.UUID}
				body, responseCode := processRequest(token, url, "POST", vcRequest)

				if responseCode == 202 {
					tasks := Tasks{}
					err := json.Unmarshal(body, &tasks)
					if err != nil {
						log.Println("Failed to parse the response body.\n[ERROR] -", err)
						os.Exit(1)
					}

					status := tasks.MonitorTask(token, tasks.TaskId, request)
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

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VCENTERS + "?page=0&size=10"

	if len(vCenters.vc_name) > 0 {
		url = url + "&vcName=" + vCenters.vc_name
	}

	if len(vCenters.vc_fqdn) > 0 {
		url = url + "&fqdn=" + vCenters.vc_fqdn
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
			if (vCenter.VCName == vCenters.vc_name) || (vCenter.Fqdn == vCenters.vc_fqdn) {
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
		vCenter := VCenters{vc_name: vcenter}
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

		status := tasks.MonitorTask(token, tasks.TaskId, request)
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

		status := tasks.MonitorTask(token, tasks.TaskId, request)
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

		status := tasks.MonitorTask(token, tasks.TaskId, request)
		if status != "SUCCESS" {
			log.Println("Failed to scan components running on the virtual machines managed by the provided vCenter")
		} else {
			log.Println("Successfully scanned components running on the virtual machines managed by the provided vCenter")
		}
	} else {
		log.Println("Failed to scan components running on the virtual machines managed by the provided vCenter. Response Code:", responseCode)
	}
}
