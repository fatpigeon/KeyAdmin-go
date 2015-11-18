package passwordadmin

var testJsonDir = ""
var testPassword = []struct {
	domain   string
	username string
	password string
}{
	{"www.baidu.com", "jia", "nopassword"},
	{"www.baidu.com", "jia", "haspassword"},
	{"www.baidu.com", "pigeon", "pigeon"},
	{"www.baidu.com", "pigeon", "bigpigeon"},
	{"www.google.com", "pigeon", "456"},
	{"www.acfun.com", "jia", "123"},
}

var testJsonData = `
{
	"www.acfun.com":{
		"jia":["123"]
	},
	"www.baidu.com":{
		"jia":["nopassword","haspassword","nopassword"],
		"pigeon":["pigeon","bigpigeon","pigeon","bigpigeon","pigeon","bigpigeon","pigeon","bigpigeon"]
	},
	"www.google.com":{
		"pigeon":["456"]
	}
}
`
var testconfigstr = `{
    "way": ["md5", "sha512"],
	"private_key_dir": "",
	"public_key_dir": "",
	"data_dir": "",
	"key_length": 2048
}
`
