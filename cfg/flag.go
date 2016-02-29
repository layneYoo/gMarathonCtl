package cfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/layneYoo/mCtl/check"
	"xlei/maractl/g"
)

func parseargs() (config, host, username, password, format string) {
	flag.StringVar(&config, "c", "", "json config file")
	flag.StringVar(&host, "h", "", "marathon host with transport and port")
	flag.StringVar(&username, "u", "", "username")
	flag.StringVar(&password, "p", "", "password")
	flag.StringVar(&format, "f", "", "output format")
	flag.Parse()
	//check.Check(flag.NFlag() == 4, "argument error, need two args")
	return
}

func Config() (g.MarathonObj, string) {
	configFile, host, name, passwd, format := parseargs()
	// todo : test the config

	if format == "" {
		format = "human"
	}

	config, err := os.Open(configFile)
	defer config.Close()
	if err != nil {
		fmt.Println("Note : no config file found, using argument")
	}
	jsonParse := json.NewDecoder(config)
	check.Check(jsonParse != nil, "json config decode error...")
	var marathonObj g.MarathonObj
	if err = jsonParse.Decode(&marathonObj); err != nil {
		//fmt.Println(err.Error())
	}
	//marathonObj.Actioninfo.Act = "atc"

	if host != "" {
		marathonObj.Marathoninfo.Host = host
	}
	if name != "" {
		marathonObj.Marathoninfo.User = name
	}
	if passwd != "" {
		marathonObj.Marathoninfo.Password = passwd
	}

	return marathonObj, format
}
