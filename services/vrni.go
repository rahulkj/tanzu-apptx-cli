package services

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type VRNI struct {
	url            string
	username       string
	password       string
	sa_alias       string
	vrni_fqdn      string
	vc_names       string
	is_SaaS        bool
	vrni_api_token string
	operation      string
}

type VRNIRequest struct {
	Fqdn               string   `json:"ip"`
	ApiToken           string   `json:"apiToken"`
	isSaaS             bool     `json:"isSaas"`
	VCenterUUIDs       []string `json:"vcUuids"`
	ServiceAccountUUID string   `json:"serviceAccountUUID"`
}

type VRNIResponse struct {
	Id       string `json:"id"`
	IP       string `json:"ip"`
	IsSaaS   bool   `json:"isSaaS`
	ApiToken string `json:"apiToken`
	VCenters []struct {
		Fqdn        string `json:"fqdn"`
		VCenterUUID string `json:"irisVcenterUUID"`
		VCName      string `json:"vcName"`
	} `json:"vcenters"`
	ServiceAccount struct {
		UUID  string `json:"uuid"`
		Alias string `json:"alias"`
	} `json:"serviceAccount"`
}

func (vRNI VRNI) Execute() {
	vRNI = vRNI.validate()

	request := Request{vRNI.url, vRNI.username, vRNI.password}
	authResponse := Authenticate(request)

	switch vRNI.operation {
	case "create":
		vRNI.create(authResponse.Token, request)
	case "delete":
		vRNI.delete(authResponse.Token, request)
	case "update-credentials":
		vRNI.update(authResponse.Token, request)
	case "add-vcenters":
		vRNI.addVcenters(authResponse.Token, request)
	case "remove-vcenters":
		vRNI.deleteVcenters(authResponse.Token, request)
	default:
		fmt.Println("Operation not supported")
		os.Exit(1)
	}
}

func (vRNI VRNI) validate() VRNI {
	vrniCmd := flag.NewFlagSet(VRNI_CMD, flag.ExitOnError)
	url := vrniCmd.String("url", "", "Iris URL, ex: appliance.example.com")
	username := vrniCmd.String("username", "", "Iris admin username")
	password := vrniCmd.String("password", "", "Iris admin password")

	vrni_fqdn := vrniCmd.String("vrni-fqdn", "", "vCenter FQDN")
	vc_names := vrniCmd.String("vc-names", "", "comma separated list of vCenter Name(s)")
	sa_alias := vrniCmd.String("sa-alias", "", "vRNI service account alias")
	is_SaaS := vrniCmd.Bool("isSaaS", false, "using a SaaS vRNI instance, default is false")
	vrni_api_token := vrniCmd.String("vrni-api-token", "", "SaaS vRNI api token")

	operation := vrniCmd.String("operation", "", "create, delete, update-credentials, add-vcenters, remove-vcenters")

	vrniCmd.Parse(os.Args[2:])

	if (*url == "" || *username == "" || *password == "") ||
		(*operation == "" || *vrni_fqdn == "") ||
		(strings.Contains(*url, "https://")) ||
		(*is_SaaS && *vrni_api_token == "" && *sa_alias != "") ||
		(!*is_SaaS && *sa_alias == "" && *vrni_api_token != "") {
		fmt.Println("subcommand 'vRNI'")
		vrniCmd.PrintDefaults()
		os.Exit(1)
	}

	if *operation == "create" && (*sa_alias == "" || *vc_names == "") {
		fmt.Println("subcommand 'vRNI'")
		vrniCmd.PrintDefaults()
		os.Exit(1)
	} else if *operation == "update-credentials" && (*sa_alias == "" || *vrni_api_token != "") {
		fmt.Println("subcommand 'vRNI'")
		vrniCmd.PrintDefaults()
		os.Exit(1)
	} else if *operation == "add-vcenters" && (*vc_names == "") {
		fmt.Println("subcommand 'vRNI'")
		vrniCmd.PrintDefaults()
		os.Exit(1)
	} else if *operation == "remove-vcenters" && (*vc_names == "") {
		fmt.Println("subcommand 'vRNI'")
		vrniCmd.PrintDefaults()
		os.Exit(1)
	}

	vRNI = VRNI{*url, *username, *password, *sa_alias, *vrni_fqdn, *vc_names, *is_SaaS, *vrni_api_token, *operation}
	return vRNI
}

func (vRNI VRNI) create(token string, request Request) {
	vCenters := VCenters{}
	vCenterUUIDs := vCenters.findAll(token, request, vRNI.vc_names)

	vrniResponses := vRNI.findAll(token, request)
	for _, vrniResponse := range vrniResponses {
		if vrniResponse.IP == vRNI.vrni_fqdn {
			log.Println("vRNI is already registered")
			os.Exit(1)
		}
	}

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS

	var vrniRequest VRNIRequest

	if vRNI.is_SaaS {
		vrniRequest = VRNIRequest{vRNI.vrni_fqdn, vRNI.vrni_api_token, vRNI.is_SaaS, vCenterUUIDs, ""}
	} else {
		serviceAccounts := ServiceAccounts{}
		response := serviceAccounts.findServiceAccount(vRNI.sa_alias, token, request)

		if len(response.Embedded.ServiceAccounts) > 0 {
			for _, serviceAccount := range response.Embedded.ServiceAccounts {
				if serviceAccount.Alias == vRNI.sa_alias {
					vrniRequest = VRNIRequest{vRNI.vrni_fqdn, vRNI.vrni_api_token, vRNI.is_SaaS, vCenterUUIDs, serviceAccount.UUID}
				} else {
					log.Println("Cannot complete the operation as the Service Account does not exist")
				}
			}
		} else {
			log.Println("Cannot complete the operation as the Service Account does not exist")
		}
	}

	_, responseCode := processRequest(token, url, "POST", vrniRequest)

	if responseCode == 200 {
		log.Println("Successfully registered vRNI with the provided information")
	} else {
		log.Println("Failed to register vRNI with the provided information. Response Code:", responseCode)
	}
}

func (vRNI VRNI) delete(token string, request Request) {
	vrniResponses := vRNI.findAll(token, request)
	for _, vrniResponse := range vrniResponses {
		if vrniResponse.IP == vRNI.vrni_fqdn {
			url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS + "/" + vrniResponse.Id
			_, responseCode := processRequest(token, url, "DELETE", nil)

			if responseCode == 204 {
				log.Println("Successfully deleted vRNI with the provided information")
			} else {
				log.Println("Failed to delete vRNI with the provided information. Response Code:", responseCode)
			}
		} else {
			log.Println("Could not find the vRNI instance provided")
		}
	}
}

func (vRNI VRNI) findAll(token string, request Request) (response []VRNIResponse) {
	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS

	body, responseCode := processRequest(token, url, "GET", nil)

	if responseCode == 200 {
		err := json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Failed to parse the response body.\n[ERROR] -", err)
			os.Exit(1)
		}
	} else {
		log.Println("Failed to register vRNI with the provided information")
	}

	return response
}
func (vRNI VRNI) update(token string, request Request) {
	vrniResponses := vRNI.findAll(token, request)
	for _, vrniResponse := range vrniResponses {
		if vrniResponse.IP == vRNI.vrni_fqdn {
			url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS + "/" + vrniResponse.Id

			var vCenterUUIDs []string
			for _, vCenter := range vrniResponse.VCenters {
				vCenterUUIDs = append(vCenterUUIDs, vCenter.VCenterUUID)
			}

			var vrniRequest VRNIRequest

			if vrniResponse.IsSaaS {
				vrniRequest = VRNIRequest{vrniResponse.IP, vRNI.vrni_api_token, true, vCenterUUIDs, ""}
			} else {
				serviceAccounts := ServiceAccounts{}
				response := serviceAccounts.findServiceAccount(vRNI.sa_alias, token, request)

				if len(response.Embedded.ServiceAccounts) > 0 {
					for _, serviceAccount := range response.Embedded.ServiceAccounts {
						if serviceAccount.Alias == vRNI.sa_alias {
							vrniRequest = VRNIRequest{vrniResponse.IP, "", false, vCenterUUIDs, serviceAccount.UUID}
						} else {
							log.Println("Cannot complete the operation as the Service Account does not exist")
						}
					}
				} else {
					log.Println("Cannot complete the operation as the Service Account does not exist")
				}
			}

			_, responseCode := processRequest(token, url, "PUT", vrniRequest)

			if responseCode == 200 {
				log.Println("Successfully updated vRNI credentials with the provided information")
			} else {
				log.Println("Failed to update vRNI with the provided information. Response Code:", responseCode)
			}
		} else {
			log.Println("Could not find the vRNI instance provided")
		}
	}
}
func (vRNI VRNI) addVcenters(token string, request Request) {
	vrniResponses := vRNI.findAll(token, request)
	for _, vrniResponse := range vrniResponses {
		if vrniResponse.IP == vRNI.vrni_fqdn {
			url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS + "/" + vrniResponse.Id

			var vCenterUUIDs []string
			for _, vCenter := range vrniResponse.VCenters {
				vCenterUUIDs = append(vCenterUUIDs, vCenter.VCenterUUID)
			}

			vCenters := VCenters{}
			newVCenterUUIDs := vCenters.findAll(token, request, vRNI.vc_names)

			vCenterUUIDs = append(vCenterUUIDs, newVCenterUUIDs...)

			var vrniRequest VRNIRequest
			if vrniResponse.IsSaaS {
				vrniRequest = VRNIRequest{vrniResponse.IP, vrniResponse.ApiToken, vrniResponse.IsSaaS, vCenterUUIDs, ""}
			} else {
				vrniRequest = VRNIRequest{vrniResponse.IP, "", vrniResponse.IsSaaS, vCenterUUIDs, vrniResponse.ServiceAccount.UUID}
			}

			_, responseCode := processRequest(token, url, "PUT", vrniRequest)

			if responseCode == 200 {
				log.Println("Successfully added the vCenters to vRNI with the provided information")
			} else {
				log.Println("Failed to add vCenters to vRNI with the provided information. Response Code:", responseCode)
			}
		} else {
			log.Println("Could not find the vRNI instance provided")
		}

	}
}

func (vRNI VRNI) deleteVcenters(token string, request Request) {
	vrniResponses := vRNI.findAll(token, request)
	for _, vrniResponse := range vrniResponses {
		if vrniResponse.IP == vRNI.vrni_fqdn {
			url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS + "/" + vrniResponse.Id

			var vCenterUUIDs []string
			for _, vCenter := range vrniResponse.VCenters {
				vCenterUUIDs = append(vCenterUUIDs, vCenter.VCenterUUID)
			}

			vCenters := VCenters{}
			toDeleteVCenterUUIDs := vCenters.findAll(token, request, vRNI.vc_names)

			for i, vCenterUUID := range vCenterUUIDs {
				for _, toDeleteVCenterUUID := range toDeleteVCenterUUIDs {
					if vCenterUUID == toDeleteVCenterUUID {
						vCenterUUIDs = append(vCenterUUIDs[:i], vCenterUUIDs[i+1:]...)
					}
				}
			}

			var vrniRequest VRNIRequest
			if vrniResponse.IsSaaS {
				vrniRequest = VRNIRequest{vrniResponse.IP, vrniResponse.ApiToken, vrniResponse.IsSaaS, vCenterUUIDs, ""}
			} else {
				vrniRequest = VRNIRequest{vrniResponse.IP, "", vrniResponse.IsSaaS, vCenterUUIDs, vrniResponse.ServiceAccount.UUID}
			}

			_, responseCode := processRequest(token, url, "PUT", vrniRequest)

			if responseCode == 200 {
				log.Println("Successfully removed the vCenters from vRNI with the provided information")
			} else {
				log.Println("Failed to remove vCenters from vRNI with the provided information. Response Code:", responseCode)
			}
		} else {
			log.Println("Could not find the vRNI instance provided")
		}

	}
}
