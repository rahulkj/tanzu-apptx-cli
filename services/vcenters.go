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

type VCenterResponse struct {
	Embedded struct {
		VCenters []struct {
			Fqdn        string `json:"fqdn"`
			VCenterUUID string `json:"irisVcenterUUID"`
			VCName      string `json:"vcName"`
		} `json:"vcenters"`
	} `json:"_embedded"`
}

func (vCenters VCenters) Execute() {
	vCenters = vCenters.validate()

	request := Request{vCenters.url, vCenters.username, vCenters.password}
	authResponse := Authenticate(request)

	switch vCenters.operation {
	case "create":
		vCenters.create(authResponse.Token, request)
	case "delete":
		vCenters.delete(authResponse.Token, request)
	default:
		fmt.Println("Operation not supported")
		os.Exit(1)
	}
}

func (vCenters VCenters) validate() VCenters {
	inputCmd := flag.NewFlagSet(VCENTER_CMD, flag.ExitOnError)
	url := inputCmd.String("url", "", "Iris URL, ex: appliance.example.com")
	username := inputCmd.String("username", "", "Iris admin username")
	password := inputCmd.String("password", "", "Iris admin password")
	vc_fqdn := inputCmd.String("vc-fqdn", "", "vCenter FQDN")
	vc_name := inputCmd.String("vc-name", "", "vCenter Name")
	sa_alias := inputCmd.String("sa-alias", "", "service account alias")
	operation := inputCmd.String("operation", "", "create, delete")

	inputCmd.Parse(os.Args[2:])

	if (*url == "" || *username == "" || *password == "") ||
		(*operation == "" || *vc_fqdn == "" || *vc_name == "") ||
		(strings.Contains(*url, "https://")) {
		fmt.Println("subcommand 'vCenter'")
		inputCmd.PrintDefaults()
		os.Exit(1)
	}

	if *operation == "create" && *sa_alias == "" {
		fmt.Println("subcommand 'vCenter'")
		inputCmd.PrintDefaults()
		os.Exit(1)
	}

	vCenters = VCenters{*url, *username, *password, *sa_alias, *vc_fqdn, *vc_name, *operation}
	return vCenters
}

func (vCenters VCenters) create(token string, request Request) {
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

func (vCenters VCenters) delete(token string, request Request) {
	vCenterUUID := vCenters.findVCenter(token, request)

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VCENTERS + "/" + vCenterUUID
	_, responseCode := processRequest(token, url, "DELETE", nil)

	if responseCode == 200 {
		log.Println("Successfully deleted vCenter")
	} else {
		log.Println("Failed to delete vCenter. Response Code: ", responseCode)
	}
}

func (vCenters VCenters) findVCenter(token string, request Request) (vCenterUUID string) {

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VCENTERS + "?page=0&size=10"

	if len(vCenters.vc_name) > 0 {
		url = url + "&vcName=" + vCenters.vc_name
	}

	if len(vCenters.vc_fqdn) > 0 {
		url = url + "&fqdn=" + vCenters.vc_fqdn
	}

	body, _ := processRequest(token, url, "GET", nil)

	response := VCenterResponse{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Failed to parse the response body.\n[ERROR] -", err)
		os.Exit(1)
	}

	if len(response.Embedded.VCenters) > 0 {
		for _, vCenter := range response.Embedded.VCenters {
			if (vCenter.VCName == vCenters.vc_name) || (vCenter.Fqdn == vCenters.vc_fqdn) {
				vCenterUUID = vCenter.VCenterUUID
			} else {
				log.Println("Cannot delete vCenter as it does not exist")
			}
		}
	} else {
		log.Println("Could not find the vCenter")
	}

	return vCenterUUID
}

func (vCenters VCenters) findAll(token string, request Request, vCentersCSV string) (vCenterUUIDs []string) {

	vCentersArray := strings.Split(vCentersCSV, ",")

	vCentersUUIDsArray := []string{}

	for _, vcenter := range vCentersArray {
		vCenter := VCenters{vc_name: vcenter}
		vCenterUUID := vCenter.findVCenter(token, request)
		vCentersUUIDsArray = append(vCentersUUIDsArray, vCenterUUID)
	}

	return vCentersUUIDsArray
}
