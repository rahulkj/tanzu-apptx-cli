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

type Components struct {
	url          string
	username     string
	password     string
	outputFormat string
	operation    string
}

func (components Components) Execute() {
	components = components.validate()

	request := Request{components.url, components.username, components.password}
	authResponse := Authenticate(request)

	switch components.operation {
	case LIST:
		list(authResponse.Token, components)
	default:
		fmt.Println("Operation not supported")
		components.printUsage()
		os.Exit(1)
	}
}

func list(token string, components Components) {
	componentsList := components.list(token)

	if components.outputFormat == "table" {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
		fmt.Fprintln(w, "Component Name\tProcess Name\tComponent Type\tVM Name\tVM UUID\tService Type\tIs Containerizable")
		for _, component := range componentsList.Embedded.Components {
			fmt.Fprintln(w, component.CompName,
				"\t", component.ProcessName, "\t", component.Type, "\t", component.VMName,
				"\t", component.VMUUID, "\t", component.ServiceType, "\t", component.IsContainerizable)

		}
		w.Flush()
	} else if components.outputFormat == "json" {
		prettyJSON, err := json.MarshalIndent(componentsList, "", "    ")
		if err != nil {
			log.Fatal("Failed to generate json", err)
		}
		fmt.Printf("%s\n", string(prettyJSON))
	} else if components.outputFormat == "csv" {
		fmt.Println("Component Name,Process Name,Component Type,VM Name,VM UUID,Service Type,Is Containerizable")
		for _, component := range componentsList.Embedded.Components {
			fmt.Println(component.CompName,
				",", component.ProcessName, ",", component.Type, ",", component.VMName,
				",", component.VMUUID, ",", component.ServiceType, ",", component.IsContainerizable)

		}
	}
}

func (components Components) validate() Components {
	listCmd := flag.NewFlagSet(LIST, flag.ExitOnError)

	if len(os.Args) < 3 {
		components.printUsage()
	}

	operation := os.Args[2]

	var url string
	var username string
	var password string
	var format string

	if operation == LIST {
		listCmd.StringVar(&url, "fqdn", "", "Application Transformer FQDN / IP, ex: appliance.example.com")
		listCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		listCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		listCmd.StringVar(&format, "output-format", "table", "Output format - (json,csv,table) (Default: table)")

		listCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, COMPONENTS_CMD, LIST)
			fmt.Println("Available Flags:")
			listCmd.PrintDefaults()
			os.Exit(1)
		}
	} else {
		components.printUsage()
	}

	components = Components{url, username, password, format, operation}
	return components
}

func (components Components) printUsage() {
	fmt.Printf("Usage: '%s %s [command]' \n", CLI_NAME, COMPONENTS_CMD)
	fmt.Println("Available Commands:")
	fmt.Printf("  %s \t\t\t%s \n", LIST, "List all components")
	os.Exit(1)
}

func (components Components) list(token string) (response ComponentsListResponse) {
	url := PROTOCOL + "://" + components.url + "/" + PREFIX + "/" + COMPONENTS

	body, responseCode := processRequest(token, url, "GET", nil)

	if responseCode == 200 {
		log.Println("Successfully fetched the list of components \n")
	} else {
		log.Println("Failed to fetch the list of components. Response code:", responseCode)
	}

	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Failed to parse the response body.\n[ERROR] -", err)
		os.Exit(1)
	}

	return response
}
