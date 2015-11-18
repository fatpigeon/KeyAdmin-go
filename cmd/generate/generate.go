package main

import (
	"bufio"
	"flag"
	"fmt"
	pa "github.com/fatpigeon/KeyAdmin-go/src/passwordadmin"
	"github.com/howeyc/gopass"
	"log"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func main() {
	var configfile, cmdWay string
	var cmdconfig pa.ConfigType
	flag.StringVar(&configfile, "json", "config.json", "config file name")
	flag.StringVar(&cmdconfig.DataUri, "pdb", "", "password database name")
	flag.StringVar(&cmdconfig.PublicKeyDir, "pub", "", "public file name")
	flag.StringVar(&cmdconfig.PrivateKeyDir, "pri", "", "private file name")
	flag.IntVar(&cmdconfig.KeyLength, "kl", 0, "rsa key length")
	flag.StringVar(&cmdWay, "way", "", "hash way with password,split with ','")

	flag.Parse()
	cmdconfig.Way = strings.Split(cmdWay, ",")

	var config pa.ConfigType
	parseerr := config.ParseJsonConfig(configfile)
	if parseerr != nil {
		log.Fatal(parseerr)
		os.Exit(2)
	}
	config.UpdateConfig(&cmdconfig)
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
}
