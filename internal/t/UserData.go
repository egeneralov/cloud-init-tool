package t

type UserData struct {
	AptPreserveSourcesList bool     `yaml:"apt_preserve_sources_list"`
	Chpasswd               Chpasswd `yaml:"chpasswd"`
	Bootcmd                []string `yaml:"bootcmd"`
	Users                  []User   `yaml:"users"`
	// 	GrowPart GrowPart `yaml:"growpart"`
}

func (self UserData) New() UserData {
	return UserData{
		AptPreserveSourcesList: true,
		Bootcmd: []string{
			"touch /var/lib/cloud/instance/locale-check.skip",
			"apt-get update -q",
		},
		/*
			GrowPart: GrowPart{
				Mode:    "auto",
				Devices: []string{"/"},
			},
		*/
	}
}
