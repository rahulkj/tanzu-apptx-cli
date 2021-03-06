package services

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type VRNI struct {
	alias              string
	url                string
	username           string
	password           string
	saAlias            string
	serviceAccountType string
	vrniFqdn           string
	vcNames            string
	isSaaS             bool
	vrniApiToken       string
	operation          string
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

	var alias string
	var url string
	var username string
	var password string
	var vrniFqdn string
	var vcNames string
	var saAlias string
	var isSaaS bool
	var vrniAPIToken string
	var serviceAccountType string

	if operation == REGISTER {
		registerCmd.StringVar(&alias, "vrni-name", "", "vRNI Name")
		registerCmd.StringVar(&url, "fqdn", "", "Application Transformer FQDN / IP, ex: appliance.example.com")
		registerCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		registerCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		registerCmd.StringVar(&vrniFqdn, "vrni-fqdn", "", "vCenter FQDN")
		registerCmd.StringVar(&vcNames, "vc-names", "", "comma separated list of vCenter Name(s)")
		registerCmd.StringVar(&saAlias, "sa-alias", "", "vRNI service account alias")
		registerCmd.StringVar(&serviceAccountType, "sa-account-type", "", "vRNI service account type, ex: LOCAL or LDAP")
		registerCmd.BoolVar(&isSaaS, "isSaaS", false, "using a SaaS vRNI instance, default is false")
		registerCmd.StringVar(&vrniAPIToken, "vrni-api-token", "", "SaaS vRNI api token")

		registerCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(vrniFqdn) == 0 || len(vcNames) == 0) ||
			(strings.Contains(url, "https://")) ||
			(isSaaS && len(vrniAPIToken) == 0 && len(saAlias) != 0) ||
			(!isSaaS && len(saAlias) == 0 && len(vrniAPIToken) != 0 && len(serviceAccountType) == 0) || len(alias) == 0 {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, VRNI_CMD, REGISTER)
			fmt.Println("Available Flags:")
			registerCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == UNREGISTER {
		unregisterCmd.StringVar(&url, "fqdn", "", "Application Transformer FQDN / IP, ex: appliance.example.com")
		unregisterCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		unregisterCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		unregisterCmd.StringVar(&vrniFqdn, "vrni-fqdn", "", "vCenter FQDN")

		unregisterCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(vrniFqdn) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, VRNI_CMD, UNREGISTER)
			fmt.Println("Available Flags:")
			unregisterCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == UPDATE_CREDENTIALS {
		updateCredentialsCmd.StringVar(&alias, "vrni-name", "", "vRNI Name")
		updateCredentialsCmd.StringVar(&url, "fqdn", "", "Application Transformer FQDN / IP, ex: appliance.example.com")
		updateCredentialsCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		updateCredentialsCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		updateCredentialsCmd.StringVar(&vrniFqdn, "vrni-fqdn", "", "vCenter FQDN")
		updateCredentialsCmd.StringVar(&saAlias, "sa-alias", "", "vRNI service account alias")
		updateCredentialsCmd.StringVar(&serviceAccountType, "sa-account-type", "", "vRNI service account type, ex: LOCAL or LDAP")
		updateCredentialsCmd.StringVar(&vrniAPIToken, "vrni-api-token", "", "SaaS vRNI api token")

		updateCredentialsCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(vrniFqdn) == 0) ||
			(strings.Contains(url, "https://")) ||
			(len(vrniAPIToken) == 0 || len(saAlias) != 0) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, VRNI_CMD, UPDATE_CREDENTIALS)
			fmt.Println("Available Flags:")
			updateCredentialsCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == ADD_VCENTERS {
		addVcentersCmd.StringVar(&url, "fqdn", "", "Application Transformer FQDN / IP, ex: appliance.example.com")
		addVcentersCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		addVcentersCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		addVcentersCmd.StringVar(&vrniFqdn, "vrni-fqdn", "", "vCenter FQDN")
		addVcentersCmd.StringVar(&vcNames, "vc-names", "", "comma separated list of vCenter Name(s)")

		addVcentersCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(vrniFqdn) == 0 || len(vcNames) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, VRNI_CMD, ADD_VCENTERS)
			fmt.Println("Available Flags:")
			addVcentersCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == REMOVE_VCENTERS {
		removeVcentersCmd.StringVar(&url, "fqdn", "", "Application Transformer FQDN / IP, ex: appliance.example.com")
		removeVcentersCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		removeVcentersCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		removeVcentersCmd.StringVar(&vrniFqdn, "vrni-fqdn", "", "vCenter FQDN")
		removeVcentersCmd.StringVar(&vcNames, "vc-names", "", "comma separated list of vCenter Name(s)")

		removeVcentersCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(vrniFqdn) == 0 || len(vcNames) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, VRNI_CMD, REMOVE_VCENTERS)
			fmt.Println("Available Flags:")
			removeVcentersCmd.PrintDefaults()
			os.Exit(1)
		}
	} else {
		vRNI.printUsage()
	}

	vRNI = VRNI{alias, url, username, password, saAlias, serviceAccountType, vrniFqdn, vcNames, isSaaS, vrniAPIToken, operation}
	return vRNI
}

func (vRNI VRNI) printUsage() {
	fmt.Printf("Usage: '%s %s [Command]' \n", CLI_NAME, VRNI_CMD)
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
	vCenterUUIDs := vCenters.findAll(token, request, vRNI.vcNames)

	vrniResponses := vRNI.findAll(token, request)
	for _, vrniResponse := range vrniResponses {
		if vrniResponse.IP == vRNI.vrniFqdn {
			fmt.Println("vRNI is already registered")
			os.Exit(1)
		}
	}

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS
	certificateThumbprint := getCertificateThumbprint(vRNI.vrniFqdn, HTTPS_PORT, "sha1")

	var vrniRequest VRNIRequest

	if vRNI.isSaaS {
		vrniRequest = VRNIRequest{vRNI.alias, vRNI.vrniFqdn, vRNI.vrniApiToken, vRNI.isSaaS, vCenterUUIDs, "", certificateThumbprint, ""}
	} else {
		serviceAccounts := ServiceAccounts{}
		response := serviceAccounts.findServiceAccount(vRNI.saAlias, token, request)

		if len(response.Embedded.ServiceAccounts) > 0 {
			for _, serviceAccount := range response.Embedded.ServiceAccounts {
				if serviceAccount.Alias == vRNI.saAlias {
					vrniRequest = VRNIRequest{vRNI.alias, vRNI.vrniFqdn, vRNI.vrniApiToken, vRNI.isSaaS, vCenterUUIDs, serviceAccount.UUID, certificateThumbprint, vRNI.serviceAccountType}
				} else {
					fmt.Println("Cannot complete the operation as the Service Account does not exist")
				}
			}
		} else {
			fmt.Println("Cannot complete the operation as the Service Account does not exist")
		}
	}

	_, responseCode := processRequest(token, url, "POST", vrniRequest)

	if responseCode == 200 {
		fmt.Println("Successfully registered vRNI with the provided information")
	} else {
		fmt.Println("Failed to register vRNI with the provided information. Response Code:", responseCode)
	}
}

func (vRNI VRNI) unregister(token string, request Request) {
	vrniResponses := vRNI.findAll(token, request)
	for _, vrniResponse := range vrniResponses {
		if vrniResponse.IP == vRNI.vrniFqdn {
			url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS + "/" + vrniResponse.Id
			_, responseCode := processRequest(token, url, "DELETE", nil)

			if responseCode == 204 {
				fmt.Println("Successfully deleted vRNI with the provided information")
			} else {
				fmt.Println("Failed to delete vRNI with the provided information. Response Code:", responseCode)
			}
		} else {
			fmt.Println("Could not find the vRNI instance provided")
		}
	}
}

func (vRNI VRNI) findAll(token string, request Request) (response []VRNIResponse) {
	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS

	body, responseCode := processRequest(token, url, "GET", nil)

	if responseCode == 200 {
		err := json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Failed to parse the response body.\n[ERROR] -", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Failed to register vRNI with the provided information")
	}

	return response
}
func (vRNI VRNI) update(token string, request Request) {
	vrniResponses := vRNI.findAll(token, request)
	for _, vrniResponse := range vrniResponses {
		if vrniResponse.IP == vRNI.vrniFqdn {
			url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS + "/" + vrniResponse.Id
			certificateThumbprint := getCertificateThumbprint(vRNI.vrniFqdn, HTTPS_PORT, "sha1")

			var vCenterUUIDs []string
			for _, vCenter := range vrniResponse.VCenters {
				vCenterUUIDs = append(vCenterUUIDs, vCenter.VCenterUUID)
			}

			var vrniRequest VRNIRequest

			if vrniResponse.IsSaaS {
				vrniRequest = VRNIRequest{vrniResponse.Alias, vrniResponse.IP, vRNI.vrniApiToken, true, vCenterUUIDs, "", certificateThumbprint, ""}
			} else {
				serviceAccounts := ServiceAccounts{}
				response := serviceAccounts.findServiceAccount(vRNI.saAlias, token, request)

				if len(response.Embedded.ServiceAccounts) > 0 {
					for _, serviceAccount := range response.Embedded.ServiceAccounts {
						if serviceAccount.Alias == vRNI.saAlias {
							vrniRequest = VRNIRequest{vRNI.alias, vrniResponse.IP, "", false, vCenterUUIDs, serviceAccount.UUID, certificateThumbprint, vRNI.serviceAccountType}
						} else {
							fmt.Println("Cannot complete the operation as the Service Account does not exist")
						}
					}
				} else {
					fmt.Println("Cannot complete the operation as the Service Account does not exist")
				}
			}

			_, responseCode := processRequest(token, url, "PUT", vrniRequest)

			if responseCode == 200 {
				fmt.Println("Successfully updated vRNI credentials with the provided information")
			} else {
				fmt.Println("Failed to update vRNI with the provided information. Response Code:", responseCode)
			}
		} else {
			fmt.Println("Could not find the vRNI instance provided")
		}
	}
}
func (vRNI VRNI) addVcenters(token string, request Request) {
	vrniResponses := vRNI.findAll(token, request)
	for _, vrniResponse := range vrniResponses {
		if vrniResponse.IP == vRNI.vrniFqdn {
			url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS + "/" + vrniResponse.Id
			certificateThumbprint := getCertificateThumbprint(vRNI.vrniFqdn, HTTPS_PORT, "md5")

			var vCenterUUIDs []string
			for _, vCenter := range vrniResponse.VCenters {
				vCenterUUIDs = append(vCenterUUIDs, vCenter.VCenterUUID)
			}

			vCenters := VCenters{}
			newVCenterUUIDs := vCenters.findAll(token, request, vRNI.vcNames)

			vCenterUUIDs = append(vCenterUUIDs, newVCenterUUIDs...)

			var vrniRequest VRNIRequest
			if vrniResponse.IsSaaS {
				vrniRequest = VRNIRequest{vrniResponse.Alias, vrniResponse.IP, vrniResponse.ApiToken, vrniResponse.IsSaaS, vCenterUUIDs, "", certificateThumbprint, ""}
			} else {
				vrniRequest = VRNIRequest{vrniResponse.Alias, vrniResponse.IP, "", vrniResponse.IsSaaS, vCenterUUIDs, vrniResponse.ServiceAccount.UUID, certificateThumbprint, vrniResponse.ServiceAccountType}
			}

			_, responseCode := processRequest(token, url, "PUT", vrniRequest)

			if responseCode == 200 {
				fmt.Println("Successfully added the vCenters to vRNI with the provided information")
			} else {
				fmt.Println("Failed to add vCenters to vRNI with the provided information. Response Code:", responseCode)
			}
		} else {
			fmt.Println("Could not find the vRNI instance provided")
		}

	}
}

func (vRNI VRNI) deleteVcenters(token string, request Request) {
	vrniResponses := vRNI.findAll(token, request)
	for _, vrniResponse := range vrniResponses {
		if vrniResponse.IP == vRNI.vrniFqdn {
			url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + VRNIS + "/" + vrniResponse.Id
			certificateThumbprint := getCertificateThumbprint(vRNI.vrniFqdn, HTTPS_PORT, "sha1")

			var vCenterUUIDs []string
			for _, vCenter := range vrniResponse.VCenters {
				vCenterUUIDs = append(vCenterUUIDs, vCenter.VCenterUUID)
			}

			vCenters := VCenters{}
			toDeleteVCenterUUIDs := vCenters.findAll(token, request, vRNI.vcNames)

			for i, vCenterUUID := range vCenterUUIDs {
				for _, toDeleteVCenterUUID := range toDeleteVCenterUUIDs {
					if vCenterUUID == toDeleteVCenterUUID {
						vCenterUUIDs = append(vCenterUUIDs[:i], vCenterUUIDs[i+1:]...)
					}
				}
			}

			var vrniRequest VRNIRequest
			if vrniResponse.IsSaaS {
				vrniRequest = VRNIRequest{vrniResponse.Alias, vrniResponse.IP, vrniResponse.ApiToken, vrniResponse.IsSaaS, vCenterUUIDs, "", certificateThumbprint, ""}
			} else {
				vrniRequest = VRNIRequest{vrniResponse.Alias, vrniResponse.IP, "", vrniResponse.IsSaaS, vCenterUUIDs, vrniResponse.ServiceAccount.UUID, certificateThumbprint, vrniResponse.ServiceAccountType}
			}

			_, responseCode := processRequest(token, url, "PUT", vrniRequest)

			if responseCode == 200 {
				fmt.Println("Successfully removed the vCenters from vRNI with the provided information")
			} else {
				fmt.Println("Failed to remove vCenters from vRNI with the provided information. Response Code:", responseCode)
			}
		} else {
			fmt.Println("Could not find the vRNI instance provided")
		}

	}
}
