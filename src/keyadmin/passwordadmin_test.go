package keyadmin

import (
	"encoding/json"
	"encoding/pem"
	"testing"
)

func TestGenerate(t *testing.T) {
	var config ConfigType
	json.Unmarshal([]byte(testconfigstr), &config)
	pubBlock, priBlock, err := _generateKey([]byte("123456"), config)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	t.Logf("pub:\n%s", string(pem.EncodeToMemory(pubBlock)))
	t.Logf("pri:\n%s", string(pem.EncodeToMemory(priBlock)))
}

func TestJsonStorageExtractPassword(t *testing.T) {
	var config ConfigType
	json.Unmarshal([]byte(testconfigstr), &config)
	passpharseByte := []byte("123456")
	pubBlock, priBlock, err := _generateKey(passpharseByte, config)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	jsonData := []byte("{}")
	for _, info := range testPassword {
		encodedpassword, err := _passwordEncoding([]byte(info.password), []byte{}, pubBlock)
		if err != nil {
			t.Errorf("%v", err)
			return
		}
		jsonNewData, err := __jsonUpdate(info.domain, info.username, string(encodedpassword), jsonData)
		if err != nil {
			t.Errorf("%v", err)
			return
		}
		jsonData = jsonNewData
	}
	domainList, err := __jsonDomains(jsonData)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	for _, domain := range domainList {
		usernameList, err := __jsonUserNames(domain, jsonData)
		if err != nil {
			t.Errorf("%v", err)
			return
		}
		for _, username := range usernameList {
			encodedpassword, err := __jsonGet(domain, username, -1, jsonData)
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			encodedpasspharse := _passpharseHash(passpharseByte, config.Way)
			passowrdbyte, err := _passwordDecoding([]byte(encodedpassword), encodedpasspharse, priBlock)

			if err != nil {
				t.Errorf("%v", err)
				return
			}
			t.Logf("domain:%s,username:%s,password:%s", domain, username, passowrdbyte)
		}
	}
}
