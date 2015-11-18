package passwordadmin

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

type Updater interface {
	Uri() url.URL
	Update(domain, username, password string) error
	Get(domain, username string, index int) (string, error)
	UserNames(domain string) ([]string, error)
	GetAll(domain, username string) ([]string, error)
	Domains() ([]string, error)
	IsDataExist() bool
	Drop() error
}

type JsonUpdater struct {
	uri url.URL
}

func (updater *JsonUpdater) Uri() url.URL {
	return updater.uri
}

func (updater *JsonUpdater) Update(domain, username, password string) (err error) {
	fullPath := updater.uri.Host + updater.uri.Path
	fileByte, err := ioutil.ReadFile(fullPath)
	if err != nil {
		fileByte = []byte("{}")
	}

	newFileBute, err := __jsonUpdate(domain, username, password, fileByte)
	if err != nil {
		return
	}
	var Writer io.Writer
	if fullPath == "" {
		Writer = os.Stdout
	} else {
		Writer, err = os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return
		}
	}
	_, err = io.WriteString(Writer, string(newFileBute))

	return
}

func __jsonUpdate(domain, username, password string, oldData []byte) (newData []byte, err error) {
	data := make(map[string]map[string][]string)
	err = json.Unmarshal(oldData, &data)
	domainUsers, ok := data[domain]
	if !ok {
		data[domain] = make(map[string][]string)
		domainUsers, _ = data[domain]
	}
	_, ok = domainUsers[username]
	if !ok {
		domainUsers[username] = []string{}
	}
	// prevent set same password with last time
	passwordList_len := len(domainUsers[username])
	if passwordList_len > 0 && strings.EqualFold(domainUsers[username][passwordList_len-1], password) {
		return
	}
	domainUsers[username] = append(domainUsers[username], password)
	newData, err = json.Marshal(data)
	if err != nil {
		return
	}
	return
}

func (updater *JsonUpdater) Get(domain, username string, index int) (password string, err error) {
	fullPath := updater.uri.Host + updater.uri.Path
	var reader io.Reader
	if fullPath == "" {
		reader = os.Stdin
	} else {
		reader, err = os.Open(fullPath)
		if err != nil {
			return
		}
	}
	fileByte, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	return __jsonGet(domain, username, index, fileByte)
}

func __jsonGet(domain, username string, index int, filebyte []byte) (password string, err error) {
	data := make(map[string]map[string][]string)
	json.Unmarshal(filebyte, &data)
	users, ok := data[domain]
	if !ok {
		err = fmt.Errorf("domain %s not found", domain)
		return
	}
	passwordList, ok := users[username]
	if !ok {
		err = fmt.Errorf("user name %s not in %s list", username, domain)
		return
	}
	passwordList_len := len(passwordList)
	if passwordList_len < 1 {
		err = fmt.Errorf("user name %s not in %s list", username, domain)
		return
	}
	index = index % passwordList_len
	password = passwordList[index]
	return
}

func (updater *JsonUpdater) UserNames(domain string) (usernameList []string, err error) {
	fullPath := updater.uri.Host + updater.uri.Path
	var reader io.Reader
	if fullPath == "" {
		reader = os.Stdin
	} else {
		reader, err = os.Open(fullPath)
		if err != nil {
			return
		}
	}
	fileByte, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	return __jsonUserNames(domain, fileByte)
}

func (updater *JsonUpdater) IsDataExist() bool {
	fullPath := updater.uri.Host + updater.uri.Path
	_, err := os.Stat(fullPath)
	isExist := os.IsNotExist(err)
	return !isExist
}

func (updater *JsonUpdater) Drop() (err error) {
	fullPath := updater.uri.Host + updater.uri.Path
	err = os.Remove(fullPath)
	if err != nil {
		return
	}
	return
}

func __jsonUserNames(domain string, filebyte []byte) (usernameList []string, err error) {
	data := make(map[string]map[string][]string)
	json.Unmarshal(filebyte, &data)
	users, ok := data[domain]
	if !ok {
		err = fmt.Errorf("domain %s not found", domain)
		return
	}
	keyLen := len(users)
	usernameList = make([]string, 0, keyLen)
	for k := range users {
		usernameList = append(usernameList, k)
	}
	return
}

func (updater *JsonUpdater) GetAll(domain, username string) (passwordList []string, err error) {
	fullPath := updater.uri.Host + updater.uri.Path
	var reader io.Reader
	if fullPath == "" {
		reader = os.Stdin
	} else {
		reader, err = os.Open(fullPath)
		if err != nil {
			return
		}
	}
	fileByte, err := ioutil.ReadAll(reader)
	if err != nil {

		return
	}
	return __jsonGetALl(domain, username, fileByte)
}

func __jsonGetALl(domain, username string, filebyte []byte) (passwordList []string, err error) {
	data := make(map[string]map[string][]string)
	json.Unmarshal(filebyte, &data)
	users, ok := data[domain]
	if !ok {
		err = fmt.Errorf("domain %s not found", domain)
		return
	}
	passwordList, ok = users[username]
	if !ok {
		err = fmt.Errorf("user name %s not in %s list", username, domain)
	}
	return
}

func (updater *JsonUpdater) Domains() (domainList []string, err error) {
	fullPath := updater.uri.Host + updater.uri.Path
	var reader io.Reader
	if fullPath == "" {
		reader = os.Stdin
	} else {
		reader, err = os.Open(fullPath)
		if err != nil {
			return
		}
	}
	fileByte, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	return __jsonDomains(fileByte)
}

func __jsonDomains(filebyte []byte) (domainList []string, err error) {
	data := make(map[string]map[string][]string)
	json.Unmarshal(filebyte, &data)
	domainListLen := len(data)
	domainList = make([]string, 0, domainListLen)
	for k := range data {
		domainList = append(domainList, k)
	}
	return
}
