package g

// add struct : Marathon (for marathonctl)
type Marathon struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// add struct : base (for docker image)
type Base struct {
	BuildPath  string `json:"buildPath"`
	DeployJson string `json:"deployJson"`
	Gitlib     string `json:"gitlib"`
	Registry   string `json:"registry"`
	DockerPre  string `json:"dockerPre"`
}

type MarathonObj struct {
	Marathoninfo Marathon `json:"marathoninfo"`
	Baseinfo     Base     `json:"baseinfo"`
}
