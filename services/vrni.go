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
	case REGISTER:
		vRNI.register(authResponse.Token, request)
	case UNREGISTER:
		vRNI.unregister(authResponse.Token, request)
	case UPDATE_CREDENTIALS:
		vRNI.update(authResponse.Token, request)
	case ADD_VCENTERS:
		vRNI.addVcenters(authResponse.Token, request)
	case REMOVE_VCENTERS:
		vRNI.deleteVcenters(authResponse.Token, request)
	default:
		fmt.Println("Operation not supported \n")
		vRNI.printUsage()
		os.Exit(1)
	}
}

func (vRNI VRNI) validate() VRNI {
	registerCmd := flag.NewFlagSet(REGISTER, flag.ExitOnError)
	unregisterCmd := flag.NewFlagSet(UNREGISTER, flag.ExitOnError)
	updateCredentialsCmd := flag.NewFlagSet(UPDATE_CREDENTIALS, flag.ExitOnError)
	addVcentersCmd := flag.NewFlagSet(ADD_VCENTERS, flag.ExitOnError)
	removeVcentersCmd := flag.NewFlagSet(REMOVE_VCENTERS, flag.ExitOnError)

	if len(os.Args) < 3 {
		vRNI.printUsage()
	}

	operation := os.Args[2]

	var url *string
	var username *string
	var password *string
	var vrni_fqdn *string
	var vc_names *string
	var sa_alias *string
	var is_SaaS *bool
	var vrni_api_token *string

	if operation == REGISTER {
		url = registerCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = registerCmd.String("username", "", "Iris admin username")
		password = registerCmd.String("password", "", "Iris admin password")
		vrni_fqdn = registerCmd.String("vrni-fqdn", "", "vCenter FQDN")
		vc_names = registerCmd.String("vc-names", "", "comma separated list of vCenter Name(s)")
		sa_alias = registerCmd.String("sa-alias", "", "vRNI service account alias")
		is_SaaS = registerCmd.Bool("isSaaS", false, "using a SaaS vRNI instance, default is false")
		vrni_api_token = registerCmd.String("vrni-api-token", "", "SaaS vRNI api token")

		registerCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*vrni_fqdn) == 0 || len(*vc_names) == 0) ||
			(strings.Contains(*url, "https://")) ||
			(*is_SaaS && len(*vrni_api_token) == 0 && len(*sa_alias) != 0) ||
			(!*is_SaaS && len(*sa_alias) == 0 && len(*vrni_api_token) != 0) {
			fmt.Println("Usage: 'iris-cli vRNI register [flags]' \n")
			fmt.Println("Flags:")
			registerCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == UNREGISTER {
		url = unregisterCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = unregisterCmd.String("username", "", "Iris admin username")
		password = unregisterCmd.String("password", "", "Iris admin password")
		vrni_fqdn = unregisterCmd.String("vrni-fqdn", "", "vCenter FQDN")

		unregisterCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*vrni_fqdn) == 0) ||
			(strings.Contains(*url, "https://")) {
			fmt.Println("Usage: 'iris-cli vRNI unregister [flags]' \n")
			fmt.Println("Flags:")
			unregisterCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == UPDATE_CREDENTIALS {
		url = updateCredentialsCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = updateCredentialsCmd.String("username", "", "Iris admin username")
		password = updateCredentialsCmd.String("password", "", "Iris admin password")
		vrni_fqdn = updateCredentialsCmd.String("vrni-fqdn", "", "vCenter FQDN")
		sa_alias = updateCredentialsCmd.String("sa-alias", "", "vRNI service account alias")
		vrni_api_token = updateCredentialsCmd.String("vrni-api-token", "", "SaaS vRNI api token")

		updateCredentialsCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*vrni_fqdn) == 0) ||
			(strings.Contains(*url, "https://")) ||
			(len(*vrni_api_token) == 0 || len(*sa_alias) != 0) {
			fmt.Println("Usage: 'iris-cli vRNI update-credentials [flags]' \n")
			fmt.Println("Flags:")
			updateCredentialsCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == ADD_VCENTERS {
		url = addVcentersCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = addVcentersCmd.String("username", "", "Iris admin username")
		password = addVcentersCmd.String("password", "", "Iris admin password")
		vrni_fqdn = addVcentersCmd.String("vrni-fqdn", "", "vCenter FQDN")
		vc_names = addVcentersCmd.String("vc-names", "", "comma separated list of vCenter Name(s)")

		addVcentersCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*vrni_fqdn) == 0 || len(*vc_names) == 0) ||
			(strings.Contains(*url, "https://")) {
			fmt.Println("Usage: 'iris-cli vRNI add-vcenters [flags]' \n")
			fmt.Println("Flags:")
			addVcentersCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == REMOVE_VCENTERS {
		url = removeVcentersCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = removeVcentersCmd.String("username", "", "Iris admin username")
		password = removeVcentersCmd.String("password", "", "Iris admin password")
		vrni_fqdn = removeVcentersCmd.String("vrni-fqdn", "", "vCenter FQDN")
		vc_names = removeVcentersCmd.String("vc-names", "", "comma separated list of vCenter Name(s)")

		removeVcentersCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*vrni_fqdn) == 0 || len(*vc_names) == 0) ||
			(strings.Contains(*url, "https://")) {
			fmt.Println("Usage: 'iris-cli vRNI remove-vcenters [flags]' \n")
			fmt.Println("Flags:")
			removeVcentersCmd.PrintDefaults()
			os.Exit(1)
		}
	} else {
		vRNI.printUsage()
	}

	vRNI = VRNI{*url, *username, *password, *sa_alias, *vrni_fqdn, *vc_names, *is_SaaS, *vrni_api_token, operation}
	return vRNI
}

func (vRNI VRNI) printUsage() {
	fmt.Println("Usage: 'iris-cli vRNI [command]' \n")
	fmt.Println("Available Commands:")
	fmt.Printf("  %s \t\t\t%s \n", REGISTER, "Register vRNI instance")
	fmt.Printf("  %s \t\t\t%s \n", UNREGISTER, "Remove vRNI instance")
	fmt.Printf("  %s \t\t%s \n", UPDATE_CREDENTIALS, "Update credentials for the vRNI instance")
	fmt.Printf("  %s \t\t\t%s \n", ADD_VCENTERS, "Add vCenters to the vRNI instance")
	fmt.Printf("  %s \t\t%s \n", REMOVE_VCENTERS, "Remove vCenters from the vRNI instance")
	os.Exit(1)
}

func (vRNI VRNI) register(token string, request Request) {
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

func (vRNI VRNI) unregister(token string, request Request) {
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
