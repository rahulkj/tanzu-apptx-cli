package services

const (
	PROTOCOL = "https"
	PREFIX   = "discovery"

	APPLICATIONS     = "applications"
	COMPONENTS       = "components"
	SESSION          = "session"
	SERVICE_ACCOUNTS = "serviceaccounts"
	QUESTIONS        = "questions"
	VCENTERS         = "vcenters"
	VRNIS            = "vrni"
	VIRTUAL_MACHINES = "virtualmachines"
	TASKS            = "tasks"
)

type Request struct {
	URL      string
	Username string
	Password string
}

const (
	SERVICE_ACCOUNT_CMD  = "service-account"
	GLOBAL_DEFAULT_CMD   = "global-default"
	VCENTER_CMD          = "vcenter"
	VRNI_CMD             = "vrni"
	QUESTIONS_CMD        = "questions"
	VIRTUAL_MACHINES_CMD = "virtual-machines"
	APPLICATIONS_CMD     = "applications"
	COMPONENTS_CMD       = "components"
)

const (
	ASSIGN                = "assign"
	CREATE                = "create"
	DELETE                = "delete"
	GET                   = "get"
	REGISTER              = "register"
	RESET                 = "reset"
	LIST                  = "list"
	UNREGISTER            = "unregister"
	UPDATE_CREDENTIALS    = "update-credentials"
	ADD_VCENTERS          = "add-vcenters"
	REMOVE_VCENTERS       = "remove-vcenters"
	SYNC_VCENTERS         = "sync"
	SCAN_VIRTUAL_MACHINES = "scan-virtual-machines"
	SCAN_COMPONENTS       = "scan-components"
	INTROSPECT            = "introspect"
	DISCOVER_TOPOLOGY     = "discover-topology"
)
