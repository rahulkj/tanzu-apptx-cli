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
	SERVICE_ACCOUNT_CMD = "serviceAccount"
	GLOBAL_DEFAULT_CMD  = "globalDefault"
	VCENTER_CMD         = "vCenter"
	VRNI_CMD            = "vRNI"
)
