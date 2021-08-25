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
)

type Request struct {
	URL      string
	Username string
	Password string
}

const (
	SERVICE_ACCOUNT_CMD = "serviceAccount"
	GLOBAL_DEFAULT_CMD  = "globalDefault"
)
