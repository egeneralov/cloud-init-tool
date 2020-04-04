package t

import (
	"fmt"
)

type Chpasswd struct {
	List   string `yaml:"list"`
	Expire bool   `yaml:"expire"`
}

func (self Chpasswd) FromUsers(users []User) Chpasswd {
	for _, u := range users {
		self.List += fmt.Sprintf("%+v:%+v\n", u.Name, u.PlainTextPasswd)
	}
	return self
}
