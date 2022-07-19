package services

type Request struct {
	URL      string
	Username string
	Password string
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token        string `json:"Set-Cookie"`
	RefreshToken string `json:"token"`
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

type ComponentsListResponse struct {
	Embedded struct {
		Components []struct {
			ID                string `json:"id"`
			VMName            string `json:"vmName"`
			VMUUID            string `json:"vmUUID"`
			Type              string `json:"type"`
			ProcessName       string `json:"processName`
			IsContainerizable bool   `json:"isContainerizable"`
			ServiceType       string `json:"serviceType"`
			CompName          string `json:"compName"`
			Owner             string `json:"owner"`
			LastIntrospect    string `json:"lastIntrospect"`
		} `json:"components"`
	} `json:"_embedded"`
}

type GlobalDefaultRequest struct {
	ServiceAccountUUID string `json:"serviceAccountUUID"`
}

type serviceAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

type serviceAccount struct {
	UUID     string `json:"uuid"`
	Alias    string `json:"alias"`
	Username string `json:"username"`
}

type response struct {
	Embedded struct {
		ServiceAccounts []struct {
			UUID     string `json:"uuid"`
			Alias    string `json:"alias"`
			Username string `json:"username"`
		} `json:"serviceAccounts"`
	} `json:"_embedded"`
}

type VCenterRequest struct {
	Fqdn                  string `json:"fqdn"`
	VCName                string `json:"vcName"`
	VCServiceAccountUUID  string `json:"vcServiceAccountUUID"`
	CertificateThumbprint string `json:"certificateThumbprint"`
}

type VCenter struct {
	Fqdn        string       `json:"fqdn"`
	VCenterUUID string       `json:"irisVcenterUUID"`
	VCName      string       `json:"vcName"`
	Datacenters []Datacenter `json:"dataCenters"`
}

type Datacenter struct {
	ModID    string `json:"modId"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Clusters []struct {
		ModID         string `json:"modId"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		ResourcePools []struct {
			ModID string `json:"modId"`
			Name  string `json:"name"`
			Type  string `json:"type"`
		} `json:"resourcePools"`
	} `json:"clusters"`
	Folders []struct {
		ModID string `json:"modId"`
		Name  string `json:"name"`
		Type  string `json:"type"`
	} `json:"folders"`
}

type VCenterListResponse struct {
	Embedded struct {
		VCenters []VCenter `json:"vcenters"`
	} `json:"_embedded"`
}

type VCenterScanVMRequest struct {
	ComponentScan         bool       `json:"componentScan"`
	ApplyCredentialPolicy bool       `json:"applyCredentialPolicy"`
	BinaryAnalysis        bool       `json:"binaryAnalysis"`
	Filters               Datacenter `json:"filters"`
}

type DiscoverTopologyRequest struct {
	Filters Datacenter `json:"filters"`
}

type VirtualMachinesResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Network      string   `json:"network"`
	Hostname     string   `json:"hostname"`
	Datastore    string   `json:"datastore"`
	IP           string   `json:"ip"`
	NumCPU       int      `json:"numCPU"`
	MemoryMB     string   `json:"memoryMB"`
	Services     []string `json:"services"`
	VcenterFqdn  string   `json:"vcenterFqdn"`
	DataCenter   string   `json:"dataCenter"`
	Cluster      string   `json:"cluster"`
	ResourcePool string   `json:"resourcePool'`
	Folder       string   `json:"folder"`
	NumOfDisks   int      `json:"numOfDisks"`
	SizeOfDisks  string   `json:"sizeOfDisks"`
}

type VirtualMachinesListResponse struct {
	Embedded struct {
		VirtualMachinesResponse []VirtualMachinesResponse `json:"virtualmachines"`
	} `json:"_embedded"`
}

type VRNIRequest struct {
	Alias                 string   `json:"alias"`
	Fqdn                  string   `json:"ip"`
	ApiToken              string   `json:"apiToken"`
	IsSaaS                bool     `json:"isSaas"`
	VCenterUUIDs          []string `json:"vcUuids"`
	ServiceAccountUUID    string   `json:"serviceAccountUUID"`
	CertificateThumbprint string   `json:"certificateThumbprint"`
	ServiceAccountType    string   `json:"vrniType"`
}

type VRNIResponse struct {
	Alias              string `json:"alias"`
	Id                 string `json:"id"`
	IP                 string `json:"ip"`
	ServiceAccountType string `json:"vrniType"`
	IsSaaS             bool   `json:"isSaaS`
	ApiToken           string `json:"apiToken`
	VCenters           []struct {
		Fqdn        string `json:"fqdn"`
		VCenterUUID string `json:"irisVcenterUUID"`
		VCName      string `json:"vcName"`
	} `json:"vcenters"`
	ServiceAccount struct {
		UUID  string `json:"uuid"`
		Alias string `json:"alias"`
	} `json:"serviceAccount"`
}
