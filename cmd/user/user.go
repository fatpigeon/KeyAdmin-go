package main

import (
	"flag"
	"fmt"
	pa "github.com/fatpigeon/KeyAdmin-go/src/passwordadmin"
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

	var domain, username bool
	flag.BoolVar(&domain, "domain", false, "domain")
	flag.BoolVar(&username, "username", false, "username")

	flag.Parse()
	cmdconfig.Way = strings.Split(cmdWay, ",")

	var config pa.ConfigType
	parseerr := config.ParseJsonConfig(configfile)
	if parseerr != nil {
		log.Fatal(parseerr)
		os.Exit(2)
	}
	config.UpdateConfig(&cmdconfig)

	if domain {
		domains, err := pa.Domains(config)
		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}
		for _, s := range domains {
			fmt.Printf("%s\n", s)
		}
	} else if username {
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
