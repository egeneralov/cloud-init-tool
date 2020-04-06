package t

import (
	//	"github.com/egeneralov/cloud-init-tool/internal/utils"
	"io/ioutil"

	"os/user"
)

type User struct {
	Name              string   `yaml:"name"`
	Gecos             string   `yaml:"gecos"`
	Sudo              string   `yaml:"sudo"`
	PlainTextPasswd   string   `yaml:"plain_text_passwd"`
	LockPasswd        bool     `yaml:"lock_passwd",default:false`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
}

func (self User) New(name string, password string) User {
	return User{
		Name:            name,
		Gecos:           name,
		PlainTextPasswd: password,
		LockPasswd:      false,
	}
}

func (self User) DefaultUser() (User, error) {
	usr, err := user.Current()
	if err != nil {
		return User{}, err
	}

	self.Name = usr.Username
	self.PlainTextPasswd = usr.Username
	self.Sudo = "ALL=(ALL) NOPASSWD:ALL"
	self.LockPasswd = false

	if usr.Name != "" {
		self.Gecos = usr.Name
	}

	//	key, err := utils.GetIDRSA()
	key, err := GetIDRSA()
	if err == nil {
		self.SSHAuthorizedKeys = append(self.SSHAuthorizedKeys, key)
	}

	return self, nil
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
