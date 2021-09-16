package services

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type GlobalDefaults struct {
	url       string
	username  string
	password  string
	saAlias   string
	saType    string
	operation string
}

type GlobalDefaultRequest struct {
	ServiceAccountUUID string `json:"serviceAccountUUID"`
}

func (globalDefaults GlobalDefaults) Execute() {
	globalDefaults = globalDefaults.validate()

	request := Request{globalDefaults.url, globalDefaults.username, globalDefaults.password}
	authResponse := Authenticate(request)

	switch globalDefaults.operation {
	case ASSIGN:
		globalDefaults.assign(authResponse.Token, request)
	case RESET:
		globalDefaults.reset(authResponse.Token)
	default:
		fmt.Println("Operation not supported")
		globalDefaults.printUsage()
		os.Exit(1)
	}
}

func (globalDefaults GlobalDefaults) validate() GlobalDefaults {
	assignCmd := flag.NewFlagSet(ASSIGN, flag.ExitOnError)
	resetCmd := flag.NewFlagSet(RESET, flag.ExitOnError)

	if len(os.Args) < 3 {
		globalDefaults.printUsage()
	}

	operation := os.Args[2]

	var url string
	var username string
	var password string
	var saAlias string
	var saType string

	if operation == ASSIGN {
		assignCmd.StringVar(&url, "url", "", "Iris URL, ex: appliance.example.com")
		assignCmd.StringVar(&username, "username", "", "Iris admin username")
		assignCmd.StringVar(&password, "password", "", "Iris admin password")
		assignCmd.StringVar(&saType, "service-account-type", "", "service account type, ex: VCs, VRNIs, LINUX_VMs")
		assignCmd.StringVar(&saAlias, "sa-alias", "", "service account alias")

		assignCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(saType) == 0 || len(saAlias) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Println("Usage: 'iris-cli global-default assign [flags]' \n")
			fmt.Println("Flags:")
			assignCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == RESET {
		resetCmd.StringVar(&url, "url", "", "Iris URL, ex: appliance.example.com")
		resetCmd.StringVar(&username, "username", "", "Iris admin username")
		resetCmd.StringVar(&password, "password", "", "Iris admin password")
		resetCmd.StringVar(&saType, "service-account-type", "", "service account type, ex: VCs, VRNIs, LINUX_VMs")

		resetCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(saType) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Println("Usage: 'iris-cli global-default reset [flags]' \n")
			fmt.Println("Flags:")
			resetCmd.PrintDefaults()
			os.Exit(1)
		}
	} else {
		globalDefaults.printUsage()
	}

	globalDefaults = GlobalDefaults{url, username, password, saAlias, saType, operation}
	return globalDefaults
}

func (globalDefaults GlobalDefaults) printUsage() {
	fmt.Println("Usage: 'iris-cli global-default [command]' \n")
	fmt.Println("Available Commands:")
	fmt.Printf("  %s \t\t\t%s \n", ASSIGN, "Set service account as a global default")
	fmt.Printf("  %s \t\t\t%s \n", RESET, "Reset the global default")
	os.Exit(1)
}

func (globalDefaults GlobalDefaults) assign(token string, request Request) {
	serviceAccounts := ServiceAccounts{}
	response := serviceAccounts.findServiceAccount(globalDefaults.saAlias, token, request)

	if len(response.Embedded.ServiceAccounts) > 0 {

		for _, serviceAccount := range response.Embedded.ServiceAccounts {
			if serviceAccount.Alias == globalDefaults.saAlias {
				url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "/defaults/" + globalDefaults.saType

				request := GlobalDefaultRequest{serviceAccount.UUID}
				_, responseCode := processRequest(token, url, "POST", request)

				if responseCode == 200 {
					log.Println("Successfully assigned the service credential to the global default")
				} else {
					log.Println("Failed to assign the service credential to the global default. Response code:", responseCode)
				}
			} else {
				log.Println("Cannot complete the operation as the Service Account does not exist")
			}
		}

	} else {
		log.Println("Cannot complete the operation as the Service Account does not exist")
	}
}

func (globalDefaults GlobalDefaults) reset(token string) {
	url := PROTOCOL + "://" + globalDefaults.url + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "/defaults/" + globalDefaults.saType
	_, responseCode := processRequest(token, url, "DELETE", nil)

	if responseCode == 200 {
		log.Println("Successfully reset the global default")
	} else {
		log.Println("Failed to reset the global default. Response code:", responseCode)
	}
}
