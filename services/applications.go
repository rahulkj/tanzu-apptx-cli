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

type Applications struct {
	url          string
	username     string
	password     string
	outputFormat string
	operation    string
}

type ApplicationsListResponse struct {
	Embedded struct {
		Applications []struct {
			ID                     string `json:"id"`
			Name                   string `json:"name"`
			ComponentsGroupedByVMs []struct {
				VMName     string `json:"vmName"`
				Components []struct {
					ID                string `json:"id"`
					VMName            string `json:"vmName"`
					VMUUID            string `json:"vmUUID"`
					Type              string `json:"type"`
					ProcessName       string `json:"processName`
					IsContainerizable bool   `json:"isContainerizable"`
					ServiceType       string `json:"serviceType"`
					CompName          string `json:"compName"`
				} `json:"components"`
			} `json:"componentsGroupedByVMs"`
		} `json:"applications"`
	} `json:"_embedded"`
}

func (applications Applications) Execute() {
	applications = applications.validate()

	request := Request{applications.url, applications.username, applications.password}
	authResponse := Authenticate(request)

	switch applications.operation {
	case LIST:
		applicationsList := applications.list(authResponse.Token)

		if applications.outputFormat == "table" {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
			fmt.Fprintln(w, "Application Name\tID\tComponent Name\tProcess Name\tComponent Type\tVM Name\tVM UUID\tService Type\tIs Containerizable")
			for _, application := range applicationsList.Embedded.Applications {
				for _, componentsGroupedByVM := range application.ComponentsGroupedByVMs {
					for _, component := range componentsGroupedByVM.Components {
						fmt.Fprintln(w, application.Name, "\t", application.ID, "\t", component.CompName,
							"\t", component.ProcessName, "\t", component.Type, "\t", component.VMName,
							"\t", component.VMUUID, "\t", component.ServiceType, "\t", component.IsContainerizable)
					}
				}
			}
			w.Flush()
		} else if applications.outputFormat == "json" {
			prettyJSON, err := json.MarshalIndent(applicationsList, "", "    ")
			if err != nil {
				log.Fatal("Failed to generate json", err)
			}
			fmt.Printf("%s\n", string(prettyJSON))
		} else if applications.outputFormat == "csv" {
			fmt.Println("Application Name,ID,Component Name,Process Name,Component Type,VM Name,VM UUID,Service Type,Is Containerizable")
			for _, application := range applicationsList.Embedded.Applications {
				for _, componentsGroupedByVM := range application.ComponentsGroupedByVMs {
					for _, component := range componentsGroupedByVM.Components {
						fmt.Println(application.Name, ",", application.ID, ",", component.CompName,
							",", component.ProcessName, ",", component.Type, ",", component.VMName,
							",", component.VMUUID, ",", component.ServiceType, ",", component.IsContainerizable)
					}
				}
			}
		}
	default:
		fmt.Println("Operation not supported")
		applications.printUsage()
		os.Exit(1)
	}
}

func (applications Applications) validate() Applications {
	listCmd := flag.NewFlagSet(LIST, flag.ExitOnError)

	if len(os.Args) < 3 {
		applications.printUsage()
	}

	operation := os.Args[2]

	var url string
	var username string
	var password string
	var format string

	if operation == LIST {
		listCmd.StringVar(&url, "url", "", "Application Transformer URL, ex: appliance.example.com")
		listCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		listCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		listCmd.StringVar(&format, "output-format", "table", "Output format - json,csv,table")

		listCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, APPLICATIONS_CMD, LIST)
			fmt.Println("Available Flags:")
			listCmd.PrintDefaults()
			os.Exit(1)
		}
	} else {
		applications.printUsage()
	}

	applications = Applications{url, username, password, format, operation}
	return applications
}

func (applications Applications) printUsage() {
	fmt.Printf("Usage: '%s %s [command]' \n", CLI_NAME, APPLICATIONS_CMD)
	fmt.Println("Available Commands:")
	fmt.Printf("  %s \t\t\t%s \n", LIST, "List all applications")
	os.Exit(1)
}

func (applications Applications) list(token string) (response ApplicationsListResponse) {
	url := PROTOCOL + "://" + applications.url + "/" + PREFIX + "/" + APPLICATIONS

	body, responseCode := processRequest(token, url, "GET", nil)

	if responseCode == 200 {
		log.Println("Successfully fetched the list of applications \n")
	} else {
		log.Println("Failed to fetch the list of application. Response code:", responseCode)
	}

	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Failed to parse the response body.\n[ERROR] -", err)
		os.Exit(1)
	}

	return response
}
