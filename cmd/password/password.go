package main

import (
	"flag"
	"fmt"
	"github.com/howeyc/gopass"
	pa "keyadmin"
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

	var update, get, getall, generate, domain, user bool
	flag.BoolVar(&update, "update", false, "update a user password")
	flag.BoolVar(&get, "get", false, "get password of one user")
	flag.BoolVar(&getall, "getall", false, "get all of a user's password")
	flag.BoolVar(&generate, "generate", false, "create new key and new storage")
	flag.BoolVar(&domain, "domain", false, "get all domain ")
	flag.BoolVar(&user, "user", false, "get all user of one domain")

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

	} else if generate {
		//warn and ensure
		fmt.Printf("generate will drop old rsa key, are you continue(y/n)")
		var y_n string
		_, _ = fmt.Scanln(&y_n)

		if y_n = strings.ToLower(y_n); y_n != "y" {
			os.Exit(0)
		}
		//delete exist data
		updater, err := pa.GetUpdater(config.DataUri)
		if err != nil {
			log.Fatal(parseerr)
		}
		if updater.IsDataExist() {
			fmt.Printf("old data file will be remove ,are you continue(y/n)")
			_, _ = fmt.Scanln(&y_n)
			if y_n = strings.ToLower(y_n); y_n != "y" {
				fmt.Printf("%v,%v", string(y_n), y_n)
				os.Exit(0)
			}
			err = updater.Drop()
			if err != nil {
				log.Fatal(parseerr)
				os.Exit(0)
			}
		}
		//input only one passwordpharse
		var passpharse, passpharse_again string
		fmt.Printf("input private key passpharse:")
		passpharse = string(gopass.GetPasswd())
		fmt.Printf("input private key passpharse again:")
		passpharse_again = string(gopass.GetPasswd())
		if !strings.EqualFold(passpharse, passpharse_again) {
			fmt.Printf("your passpharse do not match\n")
			os.Exit(0)
		}
		generr := pa.GenerateKey([]byte(passpharse), config)
		if generr != nil {
			log.Fatal(generr)
			os.Exit(2)
		}
	} else if domain {
		domains, err := pa.Domains(config)
		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}
		for _, s := range domains {
			fmt.Printf("%s\n", s)
		}
	} else if user {
		var domain string
		fmt.Printf("input domain:")
		fmt.Scanln(&domain)
		usernames, err := pa.UserNames(domain, config)
		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}
		for _, s := range usernames {
			fmt.Printf("%s\n", s)
		}
	}
}
