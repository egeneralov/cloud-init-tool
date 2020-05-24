package main

import (
	"flag"
	"github.com/egeneralov/cloud-init-tool/internal/t"
	"github.com/egeneralov/cloud-init-tool/internal/utils"
	"gopkg.in/yaml.v2"
	"os"

	"errors"
	"fmt"
	// 	"io/ioutil"
	"runtime"
	"strings"
)

var (
	hostname = utils.Hostname()
	MetaData = t.MetaData{}.New(hostname)
	UserData = t.UserData{}.New()

	RootEnabled = true
	password    = "ok"

	CloudInitDirectory string
	OutputDirectory    string

	err error

	// r.go
	APIURL    = "https://gitlab.com/api/v4"
	projectID = 17932039
	jobName   = "debirf"
	refName   = "master"

	downloadURL   string
	forceDownload = false
	// 	OutputDirectory string

	// 	isofile   string

	disksCount int
	disksSize  int
	VMOptions  = t.VMOptions{
		// 		Name:    strings.Split(OutputDirectory, "/")[2],
		CPUs:   2,
		Memory: 2048,
		/*
			BootISO: isofile,
			InitISO: OutputDirectory + "/init.iso",
		*/
		Disks: []t.Disk{
			t.Disk{
				Size:          1024 * 1000,
// 				CreateOptions: "--iface nvme",
			},
/*
			t.Disk{
				Size:          1024 * 1000,
// 				CreateOptions: "--iface nvme",
			},
*/
		},
	}
)

func main() {
	flags()

	CloudInitDirectory = OutputDirectory + "/cloud-init"

	for _, path := range []string{CloudInitDirectory, OutputDirectory} {
		os.Mkdir(path, os.ModePerm)
	}

	DownloadArtifactISO()

	prepareUserData()
	writeFiles()

	VMOptions.InitISO = OutputDirectory + "/init.iso"
	// 	VMOptions.Name = strings.Split(OutputDirectory, "/")[2]

	err = utils.DirectoryToISO(
		VMOptions.InitISO,
		CloudInitDirectory,
	)

	if err != nil {
		panic(err)
	}

	/*
		opts := t.VMOptions{
			Name:    strings.Split(OutputDirectory, "/")[2],
			CPUs:    2,
			Memory:  2048,
			BootISO: isofile,
			InitISO: OutputDirectory+"/init.iso",
			Disks: []t.Disk{
				t.Disk{
					Size:          10240,
					CreateOptions: "--iface nvme",
				},
				t.Disk{
					Size:          10240,
					CreateOptions: "--iface nvme",
				},
			},
		}
	*/

	// 	VMOptions.BootISO = isofile

	switch runtime.GOOS {
	case "darwin":
		// 		err = utils.RunParallelsVM(opts)
		err = utils.RunParallelsVM(VMOptions)
	/*
	    case "linux":
	  		cmd := utils.GenerateQemuCommand(opts)

	  		fmt.Println(cmd)

	  		command := exec.Command("bash", "-xec", cmd)
	  		_, err = command.Output()
	*/
	default:
		panic(errors.New("Not implemented for " + runtime.GOOS))
	}

	if err != nil {
		panic(err)
	}

	// 	utils.RemoveRecursive(OutputDirectory)
}

func prepareUserData() {
	myUser, err := t.User{}.DefaultUser()
	if err == nil {
		myUser.PlainTextPasswd = password
		UserData.Users = append(UserData.Users, myUser)
	}

	if RootEnabled && UserData.Users[0].Name != "root" {
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

func flags() {
	dir, err := os.Getwd()

	if err != nil {
		dir = "/"
	}

	/*
		flag.StringVar(
			&CloudInitDirectory,
			"cloud-init-directory",
			dir+"/cloud-init",
			"directory for cloud-init",
		)
	*/
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
	flag.BoolVar(
		&forceDownload,
		"force-download",
		false,
		"re-download artifact",
	)

	flag.StringVar(
		&VMOptions.Name,
		"vm-name",
		"linux-in-ram",
		"virtual machine name",
	)

	flag.Parse()
}

func DownloadArtifactISO() {
	/*
		OutputDirectory, err = ioutil.TempDir("/tmp/", "cloud-init-tool-")
		if err != nil {
			panic(err)
		}
	*/

	downloadURL = fmt.Sprintf(
		"%+v/projects/%d/jobs/artifacts/%+v/download?job=%+v",
		APIURL, projectID, refName, jobName,
	)

	if _, err := os.Stat(OutputDirectory + "/artifacts.zip"); err == nil {
		if forceDownload {
			os.Remove(OutputDirectory + "/artifacts.zip")
		}
	}

	if _, err := os.Stat(OutputDirectory + "/artifacts.zip"); err != nil {
		err = utils.DownloadFile(
			t.DownloadOptions{}.New(
				downloadURL,
				OutputDirectory+"/artifacts.zip",
			),
		)
		if err != nil {
			panic(err)
		}
	}

// 	fmt.Println("# " + OutputDirectory)

	files, err := utils.Unzip(
		OutputDirectory+"/artifacts.zip",
		OutputDirectory,
	)
	if err != nil {
		panic(err)
	}

	for _, el := range files {
		if strings.HasSuffix(el, ".iso") {
			VMOptions.BootISO = el
			// 			isofile = el
			break
		}
		if strings.HasSuffix(el, ".cgz") {
			continue
		}
		if strings.HasSuffix(el, ".zip") {
			continue
		}
	}
}
