package opt

import (
	"net/smtp"
	"path"
	"strconv"
	"strings"
	"time"

	mhn "github.com/gambol99/go-marathon"
	"github.com/layneYoo/mCtl/check"
)

type SendMail struct {
}

func (s SendMail) Apply(args []string) {
	check.Check(len(args) == 9, "mail : nine arguments needed")
	appTpl := args[0]
	appName := strings.Split(path.Base(appTpl), ".")[0] + "-" + args[1]
	config := mhn.NewDefaultConfig()
	config.URL = args[6]
	config.HTTPBasicAuthUser = args[7]
	config.HTTPBasicPassword = args[8]

	mCtrl, err := mhn.NewClient(config)
	check.Check(err == nil, "create marathon client error")
	// wait for deploy complete
	mDeploys := []*mhn.Deployment{}
getDeploy:
	mDeploys, err = mCtrl.Deployments()
	check.Check(err == nil, "get marathon deploys error")
	for _, deploy := range mDeploys {
		if deploy.AffectedApps[0] == "/"+appName {
			time.Sleep(time.Second * 5)
			goto getDeploy
		}
	}
	mApp, err := mCtrl.Application(appName)
	check.Check(err == nil, "get marathon app error")
	cont := ""
	for num, appTask := range mApp.Tasks {
		cont += "<br>" + strconv.Itoa(num+1) + ". " + appTask.Host + ":" + strconv.Itoa(appTask.Ports[0]) + "</br>"
	}
	server := args[2]
	sender := args[3] + "@" + args[2]
	passwd := args[4]
	reciver := args[5]
	to := strings.Split(reciver, ",")
	subject := "Marathon App Info : " + appName
	body := "<html><body><h5>" + "apps测试实例如下:<hr><br>" + cont + "</br></h5></body></html>"
	auth := smtp.PlainAuth("", sender, passwd, server)
	msg := []byte("To: " + reciver + "\r\nFrom: " + sender + " \r\nSubject: " + subject + "\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n" + body)
	err = smtp.SendMail(server+":25", auth, sender, to, msg)
	check.Check(err == nil, "send mail error")
	return
}
