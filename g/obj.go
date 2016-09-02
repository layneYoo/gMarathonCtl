package g

// add struct : Marathon (for marathonctl)
type Marathon struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// add struct : base (for docker image)
type GitlabInfo struct {
	Host   string `json:"githost"`
	Branch string `json:"gitbranch"`
}

type AppContraints struct {
	Offline string `json:"offline"`
	Online  string `json:"online"`
}

type AppIntances struct {
	Offline string `json:"insoffline"`
	Online  string `json:"insonline"`
}
type DockerAppInfo struct {
	Constraints AppContraints `json:"constraints"`
	Instances   AppIntances   `json:"instances"`
}

type Base struct {
	BuildPath  string        `json:"buildPath"`
	DeployJson string        `json:"deployJson"`
	Gitlib     GitlabInfo    `json:"gitlib"`
	Registry   string        `json:"registry"`
	DockerPre  DockerAppInfo `json:"dockerPre"`
}

type Mail struct {
	Server string `json:"mailserver"`
	User   string `json:"mailuser"`
	Passwd string `json:"mailpassword"`
	Sendto string `json:"mailsend"`
}

type MarathonObj struct {
	Marathoninfo Marathon `json:"marathoninfo"`
	Baseinfo     Base     `json:"baseinfo"`
	Mailinfo     Mail     `json:"mailinfo"`
}
