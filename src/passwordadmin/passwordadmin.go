package passwordadmin

import (
	//"flag"
	//"go/parser"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	//"log"
	"encoding/base64"
	"net/url"
	"os"
	"strings"
)

func _passpharseHash(passpharse []byte, encodingWay []string) (encodedpasspharse []byte) {
	encodepasspharse := passpharse
	for _, hashalgor := range encodingWay {
		switch hashalgor {
		case "md5":
			sum := md5.Sum(encodepasspharse)
			encodepasspharse = sum[:]
		case "sha1":
			sum := sha1.Sum(encodepasspharse)
			encodepasspharse = sum[:]
		case "sha224":
			sum := sha256.Sum224(encodepasspharse)
			encodepasspharse = sum[:]
		case "sha256":
			sum := sha256.Sum256(encodepasspharse)
			encodepasspharse = sum[:]
		case "sha384":
			sum := sha512.Sum384(encodepasspharse)
			encodepasspharse = sum[:]
		case "sha512":
			sum := sha512.Sum512(encodepasspharse)
			encodepasspharse = sum[:]
		}
	}
	//issue if return with not args,the return value will be null
	return encodepasspharse
}

func _generateKey(passpharse []byte, config ConfigType) (pubBlock, priBlock *pem.Block, err error) {
	encodepasspharse := _passpharseHash(passpharse, config.Way)
	pri, err := rsa.GenerateKey(rand.Reader, config.KeyLength)
	if err != nil {
		return
	}
	//public key encoding
	pubbyte, err := x509.MarshalPKIXPublicKey(pri.Public())
	if err != nil {
		return
	}
	pubBlock, err = x509.EncryptPEMBlock(rand.Reader, "RSA PUBLIC KEY", pubbyte, []byte{}, x509.PEMCipherAES256)
	if err != nil {
		return
	}
	//private key encoding

	pribyte := x509.MarshalPKCS1PrivateKey(pri)
	priBlock, err = x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", pribyte, encodepasspharse, x509.PEMCipherAES256)

	return
}

func GenerateKey(passpharse []byte, config ConfigType) (err error) {
	pubBlock, priBlock, err := _generateKey(passpharse, config)
	if err != nil {
		return
	}

	pubkeyOut, err := os.OpenFile(config.PublicKeyDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		if strings.EqualFold(config.PublicKeyDir, "") {
			pubkeyOut = os.Stdout
		} else {
			return
		}

	}
	prikeyOut, err := os.OpenFile(config.PrivateKeyDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		if strings.EqualFold(config.PrivateKeyDir, "") {
			prikeyOut = os.Stdout
		} else {
			return
		}
	}
	err = pem.Encode(pubkeyOut, pubBlock)

	if err != nil {
		return
	}
	err = pem.Encode(prikeyOut, priBlock)
	return
}

func SavePassword(domain, username, password string, config ConfigType) (err error) {
	pubBlockByte, err := ioutil.ReadFile(config.PublicKeyDir)
	if err != nil {
		return
	}
	pubBlock, _ := pem.Decode(pubBlockByte)
	encodePassword, err := _passwordEncoding([]byte(password), []byte{}, pubBlock)
	if err != nil {
		return
	}
	updater, err := GetUpdater(config.DataUri)
	if err != nil {
		return
	}
	err = updater.Update(domain, username, string(encodePassword))
	if err != nil {
		return
	}
	return
}

func GetUpdater(datauri string) (updater Updater, err error) {
	dataurl, err := url.Parse(datauri)
	if err != nil {
		return
	}
	switch dataurl.Scheme {
	case "":
		fallthrough
	case "file":
		updater = &JsonUpdater{*dataurl}
	}
	return
}

func _passwordListEncoding(
	password [][]byte,
	passpharse []byte,
	pubblock *pem.Block,
) (
	base64encodepassword [][]byte,
	err error,
) {
	pubBite, err := x509.DecryptPEMBlock(pubblock, passpharse)

	if err != nil {
		panic(err)
	}
	pub, err := x509.ParsePKIXPublicKey(pubBite)

	if err != nil {
		panic(err)
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		panic("counld not convert to rsa public key")
	}
	encodedpassword := make([][]byte, len(password))
	for i, e := range password {
		encodedpassword[i], err = rsa.EncryptPKCS1v15(rand.Reader, rsaPub, e)
	}

	if err != nil {
		panic(err)
	}
	base64encodepassword = make([][]byte, len(password))
	for i, e := range encodedpassword {
		base64encodepassword[i] = []byte(base64.StdEncoding.EncodeToString(e))
	}

	return

}

func _passwordEncoding(password, passpharse []byte, pubblock *pem.Block) (enencodedpassword []byte, err error) {
	pubBite, err := x509.DecryptPEMBlock(pubblock, passpharse)

	if err != nil {
		panic(err)
	}
	pub, err := x509.ParsePKIXPublicKey(pubBite)

	if err != nil {
		panic(err)
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		panic("counld not convert to rsa public key")
	}
	encodedpassword, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, password)
	if err != nil {
		panic(err)
	}
	enencodedpassword = []byte(base64.StdEncoding.EncodeToString(encodedpassword))

	return

}

//
func GetPassword(domain, username, passpharse string, config ConfigType) (password string, err error) {
	encodedpasspharse := _passpharseHash([]byte(passpharse), config.Way)
	updater, err := GetUpdater(config.DataUri)
	encodedPassword, err := updater.Get(domain, username, -1)
	if err != nil {
		return
	}

	priBlockByte, err := ioutil.ReadFile(config.PrivateKeyDir)
	if err != nil {
		return
	}
	priBlock, _ := pem.Decode(priBlockByte)
	if err != nil {
		return
	}
	passwordbyte, err := _passwordDecoding([]byte(encodedPassword), encodedpasspharse, priBlock)
	if err != nil {
		return
	}
	password = string(passwordbyte)
	return
}

func GetAllPassword(domain, username, passpharse string, config ConfigType) (password []string, err error) {
	encodedpasspharse := _passpharseHash([]byte(passpharse), config.Way)
	updater, err := GetUpdater(config.DataUri)
	encodedPasswordList, err := updater.GetAll(domain, username)
	if err != nil {
		return
	}

	priBlockByte, err := ioutil.ReadFile(config.PrivateKeyDir)
	if err != nil {
		return
	}
	priBlock, _ := pem.Decode(priBlockByte)
	if err != nil {
		return
	}
	passwordByteList := make([][]byte, len(encodedPasswordList))
	for i, s := range encodedPasswordList {
		passwordByteList[i] = []byte(s)
	}
	passwordbyte, err := _passwordListDecoding(passwordByteList, encodedpasspharse, priBlock)
	if err != nil {
		return
	}
	password = make([]string, len(encodedPasswordList))
	for i, e := range passwordbyte {
		password[i] = string(e)
	}
	return
}

func _passwordListDecoding(encodedpassword [][]byte, passpharse []byte, priblock *pem.Block) (password [][]byte, err error) {
	//decode with base64
	encodedpasswordList := make([][]byte, len(encodedpassword))
	for i, e := range encodedpassword {
		encodedpasswordList[i], err = base64.StdEncoding.DecodeString(string(e))
		if err != nil {
			return
		}
	}

	//compare to private key
	priByte, err := x509.DecryptPEMBlock(priblock, passpharse)
	if err != nil {
		return
	}

	pri, err := x509.ParsePKCS1PrivateKey(priByte)
	if err != nil {
		return
	}
	//decode with rsa
	password = make([][]byte, len(encodedpasswordList))
	for i, e := range encodedpasswordList {
		password[i], err = rsa.DecryptPKCS1v15(rand.Reader, pri, e)
	}

	if err != nil {
		return
	}

	return
}

func _passwordDecoding(encodedpassword, passpharse []byte, priblock *pem.Block) (password []byte, err error) {
	//decode with base64
	_encodedpassword, err := base64.StdEncoding.DecodeString(string(encodedpassword))
	if err != nil {
		return
	}

	//compare to private key
	priByte, err := x509.DecryptPEMBlock(priblock, passpharse)
	if err != nil {
		return
	}

	pri, err := x509.ParsePKCS1PrivateKey(priByte)
	if err != nil {
		return
	}
	//decode with rsa
	password, err = rsa.DecryptPKCS1v15(rand.Reader, pri, _encodedpassword)

	if err != nil {
		return
	}

	return
}

func Domains(config ConfigType) (domainlist []string, err error) {
	updater, err := GetUpdater(config.DataUri)
	if err != nil {
		return
	}
	domainlist, err = updater.Domains()
	if err != nil {
		return
	}
	return
}

func UserNames(domain string, config ConfigType) (usernamelist []string, err error) {
	updater, err := GetUpdater(config.DataUri)
	if err != nil {
		return
	}
	usernamelist, err = updater.UserNames(domain)
	if err != nil {
		return
	}
	return
}
