package opt

import (
	"fmt"
	mhn "github.com/gambol99/go-marathon"
	"github.com/layneYoo/mCtl/check"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"text/template"
)

type DeployItem struct {
	Version     string
	Appname     string
	Instance    string
	Constraints string
}

type ImageBuild struct {
	//gitDstPath string
}

func (m ImageBuild) Apply(args []string) {
	check.Check(len(args) == 12, "seven arguments needed")
	if args[0] == "" || args[1] == "" || args[2] == "" {
		log.Fatal("argument null")
		return
	}
	buildPath := args[0]
	registryPath := args[1]
	gitUrl := args[2]
	gitBranch := args[3]
	if gitBranch == "" {
		gitBranch = "master"
	}
	deployTpl := args[4]
	TlpnamePre := strings.Split(path.Base(deployTpl), ".")[0]
	TlpPath := path.Dir(args[4])
	dockerConstraintsOffline := args[5]
	dockerConstraintsOnline := args[6]
	dockerOffInstances := args[7]
	dockerOnInstances := args[8]
	marathonUrl := args[9]
	marathonUser := args[10]
	marathonPasswd := args[11]
	// get app instance
	offInstance := getAppInstance(marathonUrl, marathonUser, marathonPasswd, TlpnamePre+"-test")
	if offInstance != "" {
		dockerOffInstances = offInstance
	}
	onInstance := getAppInstance(marathonUrl, marathonUser, marathonPasswd, TlpnamePre+"-online")
	if onInstance != "" {
		dockerOnInstances = onInstance
	}
	testCts := strings.Split(dockerConstraintsOffline, ",")
	proCts := strings.Split(dockerConstraintsOnline, ",")
	// testing the path
	_, err := os.Stat(buildPath)
	if err != nil {
		existOr := os.IsExist(err)
		//Check(existOr, "error : ["+buildPath+"] No such directory")
		// not exist, git clone
		if existOr == false {
			out, err := exec.Command("bash", "-c", "/usr/local/bin/git clone -b "+gitBranch+" "+gitUrl+" "+buildPath).Output()
			check.Check(err == nil, "git clone error")
			fmt.Println("\n git clone " + string(out))
		}
	} else {
		// exist, git pull
		out, err := exec.Command("bash", "-c", "cd "+buildPath+" && /usr/local/bin/git pull origin "+gitBranch).Output()
		check.Check(err == nil, "git pull error")
		fmt.Println("\n " + string(out))
	}
	// testing registryPath
	// testing the json
	buildCmdHead := "cd "
	buildCmdGitV := ` && /usr/local/bin/git log -1 | head -1 | awk -F" " '{print $2}'`

	// get the commit version
	out, err := exec.Command("bash", "-c", buildCmdHead+buildPath+buildCmdGitV).Output()
	check.Check(err == nil, "get git version error")
	//gitVersion := string(out[0 : len(out)-2])
	gitVersion := string(out[0:9])

	// build docker image
	buildCmdBuild := ` && docker build -t `
	out, err = exec.Command("bash", "-c", buildCmdHead+buildPath+buildCmdBuild+registryPath+":"+gitVersion+" .").Output()
	check.Check(err == nil, "build command error:")
	fmt.Println("\n" + string(out))

	// create the marathon's json for deploying
	tlp, err := template.ParseFiles(deployTpl)
	check.Check(err == nil, "template parsefile error")
	deployTestCts := "\""
	for i := 0; i < len(testCts); i++ {
		if i < len(testCts)-1 {
			deployTestCts += testCts[i] + "\", \""
		} else {
			deployTestCts += testCts[i]
		}
	}
	deployTestCts += "\""
	deployOnlineCts := "\""
	for i := 0; i < len(proCts); i++ {
		if i < len(proCts)-1 {
			deployOnlineCts += proCts[i] + "\", \""
		} else {
			deployOnlineCts += proCts[i]
		}
	}
	deployOnlineCts += "\""
	deployTest := DeployItem{Version: gitVersion, Instance: dockerOffInstances, Appname: TlpnamePre + "-offline", Constraints: deployTestCts}
	deployOnline := DeployItem{Version: gitVersion, Instance: dockerOnInstances, Appname: TlpnamePre + "-online", Constraints: deployOnlineCts}
	deployNameTest := TlpPath + "/" + TlpnamePre + "_offline.json"
	deployNameOnline := TlpPath + "/" + TlpnamePre + "_online.json"
	ofpTest, err := os.OpenFile(deployNameTest, os.O_WRONLY|os.O_CREATE, 0666)
	ofpOnline, err := os.OpenFile(deployNameOnline, os.O_WRONLY|os.O_CREATE, 0666)
	check.Check(err == nil, "create file error")
	defer ofpTest.Close()
	defer ofpOnline.Close()
	err = tlp.Execute(ofpOnline, deployOnline)
	check.Check(err == nil, "template pro execute error")
	err = tlp.Execute(ofpTest, deployTest)
	check.Check(err == nil, "template test execute error")
}

func getAppInstance(marathon, user, passwd, appName string) string {
	config := mhn.NewDefaultConfig()
	config.URL = marathon
	config.HTTPBasicAuthUser = user
	config.HTTPBasicPassword = passwd

	mCtrl, err := mhn.NewClient(config)
	check.Check(err == nil, "create marathon client error")
	mApp, err := mCtrl.Application(appName)
	if err != nil {
		//panic(err)
		return ""
	}
	inst := *mApp.Instances
	return strconv.Itoa(inst)
}

type ImageUpload struct {
}

func (m ImageUpload) Apply(args []string) {
	check.Check(len(args) == 2, "two arguments needed")
	if args[0] == "" || args[1] == "" {
		log.Fatal("argument null")
		return
	}
	buildPath := args[0]
	registryPath := args[1]
	buildCmdHead := "cd "
	buildCmdGitV := ` && /usr/local/bin/git log -1 | head -1 | awk -F" " '{print $2}'`

	// get the commit version
	out, err := exec.Command("bash", "-c", buildCmdHead+buildPath+buildCmdGitV).Output()
	check.Check(err == nil, "get git version error")
	gitVersion := string(out[0:9])

	// push the image[ registryPath:gitVersion ]
	out, err = exec.Command("bash", "-c", "docker push "+registryPath+":"+gitVersion).Output()
	check.Check(err == nil, "git push error")
	fmt.Println(string(out))
}
