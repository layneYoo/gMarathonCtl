package cfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Masterminds/glide/msg"
	"github.com/gMarathonCtl/g"
	"github.com/layneYoo/mCtl/check"
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

	var marathonObj g.MarathonObj

	if err != nil {
		msg.Info("no config file found, using argument: -h ...\n")
		/*
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
		*/
		goto init
	} else {
		jsonParse := json.NewDecoder(config)
		check.Check(jsonParse != nil, "json config decode error...")
		if err = jsonParse.Decode(&marathonObj); err != nil {
			fmt.Println(err.Error())
		}
		goto init
	}
init:
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
