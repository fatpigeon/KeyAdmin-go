package keyadmin

import (
	"net/url"
	"testing"
)

func TestJsonUpdate(t *testing.T) {
	u, err := url.Parse(testJsonDir)
	if err != nil {
		t.Error(err)
	}
	updater := JsonUpdater{*u}
	for _, info := range testPassword {
		err = updater.Update(info.domain, info.username, info.password)
		if err != nil {
			t.Error(err)
		}
	}

}

func TestJsonGet(t *testing.T) {
	if testJsonDir != "" {
		u, err := url.Parse(testJsonDir)
		if err != nil {
			t.Error(err)
		}
		updater := JsonUpdater{*u}
		domainList, err := updater.Domains()

		if err != nil {
			t.Error(err)
		}
		t.Log(domainList)
		for _, domain := range domainList {
			usernameList, err := updater.UserNames(domain)
			if err != nil {
				t.Error(err)
			}
			for _, username := range usernameList {
				password, err := updater.Get(domain, username, -1)
				if err != nil {
					t.Error(err)
				}
				t.Logf("domain:%s,username:%s, password:%s", domain, username, password)
			}
		}
	} else {
		testFileByte := []byte(testJsonData)
		domainList, err := __jsonDomains(testFileByte)

		if err != nil {
			t.Error(err)
		}
		t.Log(domainList)
		for _, domain := range domainList {
			usernameList, err := __jsonUserNames(domain, testFileByte)
			if err != nil {
				t.Error(err)
			}
			for _, username := range usernameList {
				password, err := __jsonGet(domain, username, -1, testFileByte)
				if err != nil {
					t.Error(err)
				}
				t.Logf("domain:%s,username:%s, password:%s", domain, username, password)
			}
		}
	}

}
