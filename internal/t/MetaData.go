package t

type MetaData struct {
	// 	DSMode        string `yaml:"dsmode",default:"auto"`
	InstanceID    string `yaml:"instance-id"`
	LocalHostname string `yaml:"local-hostname"`
}

func (self MetaData) New(hostname string) MetaData {
	return MetaData{
		// 		DSMode:        "local",
		InstanceID:    "iid-local01",
		LocalHostname: hostname,
// 		LocalHostname: "cloudimg",
	}
}
