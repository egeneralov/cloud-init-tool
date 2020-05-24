package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"

	"io"
	// 	"io/ioutil"
	"net"
	"net/http"
	"time"

	"archive/zip"
	"path/filepath"
	"strings"

	"github.com/egeneralov/cloud-init-tool/internal/t"
)

func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "localhost"
	}
	return hostname

}

func DirectoryToISO(output string, directory string) error {
// 	fmt.Println(output)
	os.Remove(output)

	var command *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// hdiutil makehybrid -o init.iso -hfs -joliet -iso -default-volume-name cidata config/
		command = exec.Command("hdiutil", "makehybrid", "-o", output, "-hfs", "-joliet", "-iso", "-default-volume-name", "cidata", directory)
	case "linux":
		// genisoimage -J -r -V cidata -o /root/init.iso /root/cloud-init
		command = exec.Command("genisoimage", "-J", "-r", "-V", "cidata", "-o", output, directory)
	default:
		return errors.New("Not implemented for " + runtime.GOOS)
	}

	_, err := command.Output()

	if err != nil {
		return fmt.Errorf("Command: %+v\nError: %+v", command, err)
	}

	return nil
}

func RunCmd(cmd string) ([]byte, error) {
	var command *exec.Cmd
	command = exec.Command("/bin/bash", "-xec", cmd)
	return command.Output()
}

func RunParallelsVM(opts t.VMOptions) error {
	fmt.Printf("%+q\n", opts)

	script := []string{
		fmt.Sprintf("prlctl list | grep %v && prlctl stop %v --kill", opts.Name, opts.Name),
		fmt.Sprintf("prlctl list -a | grep %v && prlctl delete %v", opts.Name, opts.Name),
		// 		fmt.Sprintf("prlctl delete %v", opts.Name),
		fmt.Sprintf("prlctl create %v --ostype linux --distribution debian --location /tmp/", opts.Name),
		fmt.Sprintf(
			"prlctl set %v --cpus %d --memsize %d --autostart off --autostop stop --bios-type legacy --faster-vm on --resource-quota unlimited",
			opts.Name, opts.CPUs, opts.Memory,
		),
		fmt.Sprintf("prlctl set %v --device-del hdd0 --destroy-image", opts.Name),
		fmt.Sprintf("prlctl set %v --device-del usb", opts.Name),
		fmt.Sprintf("prlctl set %v --device-del sound0", opts.Name),
		fmt.Sprintf("prlctl set %v --device-set cdrom0 --image %v --connect", opts.Name, opts.BootISO),
		fmt.Sprintf("prlctl set %v --device-add cdrom --image %v --connect", opts.Name, opts.InitISO),
	}
	after_script := []string{
		fmt.Sprintf("prlctl start %v", opts.Name),
		/*
		   		fmt.Sprintf("while [ \"$(prlctl list | grep '%v')\" ]; do sleep 1; done;", opts.Name),
		   // 		fmt.Sprintf("prlctl stop %v --kill", opts.Name),
		   		fmt.Sprintf("prlctl delete %v", opts.Name),
		*/
	}

	for _, disk := range opts.Disks {
		//   	fmt.Println(disk)
		script = append(
			script,
			fmt.Sprintf(
				"prlctl set %v --device-add hdd --type expand --size %d %v",
				opts.Name, disk.Size, disk.CreateOptions,
			),
		)
	}

	for _, el := range after_script {
		script = append(script, el)
	}

	for _, el := range script {
		fmt.Println(el)
		RunCmd(el)

		/*
		   		_, err := RunCmd(el)
		   		if err != nil {
		   			fmt.Println(el)
		   			fmt.Println(err)
		   // 			return err
		   		}
		*/

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

func GenerateQemuCommand(opts t.VMOptions) string {
	return fmt.Sprintf(
		"/usr/local/bin/qemu-system-x86_64 -smp %d -m %d -drive file='%v' -drive file='%v' %+v %+v\n",
		// 		"qemu-system-x86_64 -smp %d -m %d -drive file='%v' -drive file='%v' %+v %+v\n",
		opts.CPUs, opts.Memory, opts.BootISO, opts.InitISO, opts.Network, opts.Args,
	)
}

/*
func RunQemu(opts t.VMOptions) error {
	command := exec.Command("bash", "-xec", "\"" + cmd + "\"")
	_, err = command.Output()

	if err != nil {
		panic(err)
	}
}
*/

/*
func DownloadFile(filepath string, url string, timeoutDial int, timeoutTLS int, timeoutHTTP int) error {
	netClient := &http.Client{
		Timeout: time.Second * time.Duration(timeoutHTTP),
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Second * time.Duration(timeoutDial),
			}).Dial,
			TLSHandshakeTimeout: time.Second * time.Duration(timeoutTLS),
		},
	}
	resp, err := netClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
*/

func DownloadFile(o t.DownloadOptions) error {
	netClient := &http.Client{
		Timeout: time.Second * time.Duration(o.TimeoutHTTP),
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Second * time.Duration(o.TimeoutDial),
			}).Dial,
			TLSHandshakeTimeout: time.Second * time.Duration(o.TimeoutTLS),
		},
	}
	resp, err := netClient.Get(o.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(o.Filename)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func RemoveRecursive(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	os.Remove(dir)
	return nil
}
