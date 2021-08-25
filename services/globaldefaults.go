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
	sa_alias  string
	sa_type   string
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
	case "assign":
		globalDefaults.assign(authResponse.Token, request)
	case "reset":
		globalDefaults.reset(authResponse.Token)
	default:
		fmt.Println("Operation not supported")
		os.Exit(1)
	}
}

func (globalDefaults GlobalDefaults) validate() GlobalDefaults {
	inputCmd := flag.NewFlagSet(GLOBAL_DEFAULT_CMD, flag.ExitOnError)
	url := inputCmd.String("url", "", "Iris URL, ex: appliance.example.com")
	username := inputCmd.String("username", "", "Iris admin username")
	password := inputCmd.String("password", "", "Iris admin password")
	sa_alias := inputCmd.String("sa-alias", "", "service account alias")
	sa_type := inputCmd.String("service-account-type", "", "service account type, ex: VCs, VRNIs, LINUX_VMs")
	operation := inputCmd.String("operation", "", "assign, reset")

	inputCmd.Parse(os.Args[2:])

	if (*url == "" || *username == "" || *password == "") ||
		(*operation == "" || *sa_type == "") ||
		(strings.Contains(*url, "https://")) {
		fmt.Println("subcommand 'globalDefault'")
		inputCmd.PrintDefaults()
		os.Exit(1)
	}

	if *operation == "assign" && *sa_alias == "" {
		fmt.Println("subcommand 'globalDefault'")
		inputCmd.PrintDefaults()
		os.Exit(1)
	}

	globalDefaults = GlobalDefaults{*url, *username, *password, *sa_alias, *sa_type, *operation}
	return globalDefaults
}

func (globalDefaults GlobalDefaults) assign(token string, request Request) {
	serviceAccounts := ServiceAccounts{}
	response := serviceAccounts.findServiceAccount(globalDefaults.sa_alias, token, request)

	if len(response.Embedded.ServiceAccounts) > 0 {

		for _, serviceAccount := range response.Embedded.ServiceAccounts {
			if serviceAccount.Alias == globalDefaults.sa_alias {
				url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "/defaults/" + globalDefaults.sa_type

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
	url := PROTOCOL + "://" + globalDefaults.url + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "/defaults/" + globalDefaults.sa_type
	_, responseCode := processRequest(token, url, "DELETE", nil)

	if responseCode == 200 {
		log.Println("Successfully reset the global default")
	} else {
		log.Println("Failed to reset the global default. Response code:", responseCode)
	}
}
