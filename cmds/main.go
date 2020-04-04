package main

import (
	"flag"
	"github.com/egeneralov/cloud-init-tool/internal/t"
	"github.com/egeneralov/cloud-init-tool/internal/utils"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	hostname = utils.Hostname()
	MetaData = t.MetaData{}.New(hostname)
	UserData = t.UserData{}.New()

	RootEnabled = true
	password    = "ok"

	CloudInitDirectory string
	OutputDirectory    string
)

func main() {
	dir, err := os.Getwd()

	if err != nil {
		dir = "/"
	}

	flag.StringVar(
		&CloudInitDirectory,
		"cloud-init-directory",
		dir+"/cloud-init",
		"directory for cloud-init",
	)
	flag.StringVar(
		&OutputDirectory,
		"output-directory",
		dir,
		"output directory",
	)
	flag.StringVar(
		&password,
		"password",
		password,
		"password",
	)
	flag.BoolVar(
		&RootEnabled,
		"root-enabled",
		false,
		"enable root access",
	)
	/*
		flag.BoolVar(
			&UserData.AptPreserveSourcesList,
			"apt-preserve-sources-list",
			false,
			"apt-preserve-sources-list",
		)
	*/
	flag.Parse()

	for _, path := range []string{CloudInitDirectory, OutputDirectory} {
		os.Mkdir(path, os.ModePerm)
	}

	prepareUserData()
	writeFiles()
	err = utils.DirectoryToISO(
		OutputDirectory+"/init.iso",
		CloudInitDirectory,
	)

	if err != nil {
		panic(err)
	}

}

func prepareUserData() {
	myUser, err := t.User{}.DefaultUser()
	if err == nil {
		myUser.PlainTextPasswd = password
		UserData.Users = append(UserData.Users, myUser)
	}

	if RootEnabled {
		root := t.User{}.New("root", password)
		root.Sudo = "ALL=(ALL) NOPASSWD:ALL"
		sshKey, err := utils.GetIDRSA()
		if err == nil {
			root.SSHAuthorizedKeys = append(root.SSHAuthorizedKeys, sshKey)
		}
		UserData.Users = append(UserData.Users, root)
	}

	UserData.Chpasswd = t.Chpasswd{}.FromUsers(UserData.Users)
}

func writeFiles() {
	// meta-data
	MetaDataYaml, err := yaml.Marshal(MetaData)
	if err != nil {
		panic(err)
	}
	err = utils.WriteStringToFile(
		CloudInitDirectory+"/meta-data",
		string(MetaDataYaml),
	)
	if err != nil {
		panic(err)
	}

	// user-data
	UserDataYaml, err := yaml.Marshal(UserData)
	if err != nil {
		panic(err)
	}
	err = utils.WriteStringToFile(
		CloudInitDirectory+"/user-data",
		"#cloud-config\n"+string(UserDataYaml),
	)
	if err != nil {
		panic(err)
	}
}
