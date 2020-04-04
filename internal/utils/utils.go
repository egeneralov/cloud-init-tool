package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"
)

func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "localhost"
	}
	return hostname

}

func DirectoryToISO(output string, directory string) error {
	os.Remove(output)

	var command *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// hdiutil makehybrid -o init.iso -hfs -joliet -iso -default-volume-name cidata config/
		command = exec.Command("hdiutil", "makehybrid", "-o", output, "-hfs", "-joliet", "-iso", "-default-volume-name", "cidata", directory)
	default:
		return errors.New("Not implemented for " + runtime.GOOS)
	}

	_, err := command.Output()
	/*
		out, err := command.Output()
		fmt.Println(string(out))
	*/

	if err != nil {
		return fmt.Errorf("Command: %+v\nError: %+v", command, err)
	}

	return nil
}

func GetIDRSA() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	dat, err := ioutil.ReadFile(usr.HomeDir + "/.ssh/id_rsa.pub")
	if err != nil {
		return "", err
	}

	return string(dat), nil
}

func createFile(filepath string) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		var file, err = os.Create(filepath)
		if err != nil {
			return
		}
		defer file.Close()
	}
}

func WriteStringToFile(filepath string, content string) error {
	os.Remove(filepath)
	createFile(filepath)

	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}
