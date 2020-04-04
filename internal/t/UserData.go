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
		/*
			GrowPart: GrowPart{
				Mode:    "auto",
				Devices: []string{"/"},
			},
		*/
	}
}
