package main

import (
	"flag"
	"fmt"
	pa "github.com/fatpigeon/KeyAdmin-go/src/passwordadmin"
	"github.com/howeyc/gopass"
	"log"
	"os"
	"strings"
)

func main() {
	var configfile, cmdWay string
	var cmdconfig pa.ConfigType
	flag.StringVar(&configfile, "json", "config.json", "config file name")
	flag.StringVar(&cmdconfig.DataUri, "pdb", "", "password database name")
	flag.StringVar(&cmdconfig.PublicKeyDir, "pub", "", "public file name")
	flag.StringVar(&cmdconfig.PrivateKeyDir, "pri", "", "private file name")
	flag.IntVar(&cmdconfig.KeyLength, "kl", 0, "rsa key length")
	flag.StringVar(&cmdWay, "way", "", "hash way with password,split with ','")

	var update, get, getall bool
	flag.BoolVar(&update, "update", false, "update or get or get all")
	flag.BoolVar(&get, "get", false, "update or get or get all")
	flag.BoolVar(&getall, "getall", false, "update or get or get all")

	flag.Parse()
	cmdconfig.Way = strings.Split(cmdWay, ",")

	var config pa.ConfigType
	parseerr := config.ParseJsonConfig(configfile)
	if parseerr != nil {
		log.Fatal(parseerr)
		os.Exit(2)
	}
	config.UpdateConfig(&cmdconfig)
	//
	if update {
		var domain, username, password, passwordAgain string
		fmt.Printf("input domain:")
		_, _ = fmt.Scanln(&domain)
		fmt.Printf("input user name:")
		_, _ = fmt.Scanln(&username)
		fmt.Printf("input password:")
		password = string(gopass.GetPasswd())
		fmt.Printf("input password agein:")
		passwordAgain = string(gopass.GetPasswd())
		if password != passwordAgain {
			fmt.Printf("your passpharse do not match\n")
			os.Exit(0)
		}
		err := pa.SavePassword(domain, username, password, config)
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}
	} else if get {
		var domain, username, passpharse string
		fmt.Printf("input domain:")
		_, _ = fmt.Scanln(&domain)
		fmt.Printf("input user name:")
		_, _ = fmt.Scanln(&username)
		fmt.Printf("input decode passpharse:")
		passpharse = string(gopass.GetPasswd())
		password, err := pa.GetPassword(domain, username, passpharse, config)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(0)
		}
		fmt.Printf("domain: %s\nuser name: %s\npassword: %s\n", domain, username, password)
	} else if getall {
		var domain, username, passpharse string
		fmt.Printf("input domain:")
		_, _ = fmt.Scanln(&domain)
		fmt.Printf("input user name:")
		_, _ = fmt.Scanln(&username)
		fmt.Printf("input decode passpharse:")
		passpharse = string(gopass.GetPasswd())
		password, err := pa.GetAllPassword(domain, username, passpharse, config)
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}
		fmt.Printf("domain: %s\nuser name: %s\npassword: %v\n", domain, username, password)

	}
}
