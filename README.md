#Key Admin
A password protect product for [Go language](http://golang.org).

**Go 1.3+ is required.**

##Quick start
####Install command

      go get github.com/fatpigeon/KeyAdmin-go/cmd/generate
      go get github.com/fatpigeon/KeyAdmin-go/cmd/password
      go get github.com/fatpigeon/KeyAdmin-go/cmd/user

####command
- Windows

```
      # generate
      mkdir keyadmin
      cd keyadmin
      mkdir data
      echo {^
          "way": ["md5", "sha512"],^
          "private_key_dir": "private",^
          "public_key_dir": "public",^
          "data_uri": "data/password.json",^
          "key_length": 2048^
      } > config.json
      %GOPATH%\bin\generate
      ...
      # update
      %GOPATH%\bin\password -update
      ...
      # find
      %GOPATH%\bin\password -get
      ...
```

- Unix like

```
> generate
      mkdir keyadmin
      cd keyadmin
      mkdir data
      echo $'{\
        "way": ["md5", "sha512"],\
        "private_key_dir": "private",\
        "public_key_dir": "public",\
        "data_uri": "data/password.json",\
        "key_length": 2048\
      }' > config.json
      $GOPATH/bin/generate
      ...
      # update
      $GOPATH/bin/password -update
      ...
      # find
      $GOPATH/bin/password -get
      ...
```
##Usage

- common

```
      <command> [-json <config>] [-pdb <password data>] 
      [-pri <private file>] [-pub <public file>] [-kl <key length>] 
      [-way <md5,sha1,...>]
```
- generate
```
      generate 
```
- password
```
      password [-get] [-getall] [-update]
```
- use
```
      use [-domain] [-username]
```


##ps
- keep you private file
