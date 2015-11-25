package keyadmin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

type ConfigType struct {
	DataUri       string   `json:"data_uri"`
	PublicKeyDir  string   `json:"public_key_dir"`
	PrivateKeyDir string   `json:"private_key_dir"`
	Way           []string `json:"way"`
	KeyLength     int      `json:"key_length"`
}

//config part
func (config *ConfigType) ParseJsonConfig(configfile string) (err error) {
	if strings.EqualFold(configfile, "") {
		return
	}
	_, err = os.Stat(configfile)
	if err != nil {

		return
	}

	filebyte, err := ioutil.ReadFile(configfile)
	if err != nil {

		return
	}
	if err = json.Unmarshal(filebyte, config); err != nil {
		return
	}
	return
}

func (config *ConfigType) UpdateConfig(newconfig *ConfigType) {
	// Using reflection here is not necessary, but it's a good exercise.
	// For more information on reflections in Go, read "The Laws of Reflection"
	// http://golang.org/doc/articles/laws_of_reflection.html
	newVal := reflect.ValueOf(newconfig).Elem()
	oldVal := reflect.ValueOf(config).Elem()

	// typeOfT := newVal.Type()
	for i := 0; i < newVal.NumField(); i++ {
		newField := newVal.Field(i)
		oldField := oldVal.Field(i)
		// log.Printf("%d: %s %s = %v\n", i,
		// typeOfT.Field(i).Name, newField.Type(), newField.Interface())
		switch newField.Kind() {
		case reflect.Interface:
			if fmt.Sprintf("%v", newField.Interface()) != "" {
				oldField.Set(newField)
			}
		case reflect.String:
			s := newField.String()
			if s != "" {
				oldField.SetString(s)
			}
		case reflect.Int:
			i := newField.Int()
			if i != 0 {
				oldField.SetInt(i)
			}
		}
	}
}
