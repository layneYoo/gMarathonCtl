// Author Seth Hoenig 2015

// Command marathonctl is a CLI tool for Marathon
package main

import (
	"flag"
	"fmt"
	//"os"

	"github.com/layneYoo/mCtl/check"
	mctl "github.com/layneYoo/mCtl/marathon"
	"xlei/maractl/cfg"
	"xlei/maractl/g"
	"xlei/maractl/opt"
)

func main() {
	//host, login, format, e := mctl.Config()
	e, formatter := cfg.Config()

	//if e != nil {
	//	fmt.Printf("config error: %s\n\n", e)
	//	check.Usage()
	//}
	//fmt.Println(e)
	f := mctl.NewFormatter(formatter)
	l := mctl.NewLogin(e.Marathoninfo.Host, e.Marathoninfo.User+":"+e.Marathoninfo.Password)
	c := mctl.NewClient(l)
	app := &mctl.Category{
		Actions: map[string]mctl.Action{
			"list":     mctl.AppList{c, f},
			"versions": mctl.AppVersions{c, f},
			"show":     mctl.AppShow{c, f},
			"create":   mctl.AppCreate{c, f},
			"update":   mctl.AppUpdate{c, f},
			"restart":  mctl.AppRestart{c, f},
			"destroy":  mctl.AppDestroy{c, f},
		},
	}
	task := &mctl.Category{
		Actions: map[string]mctl.Action{
			"list":  mctl.TaskList{c, f},
			"kill":  mctl.TaskKill{c, f},
			"queue": mctl.TaskQueue{c, f},
		},
	}
	group := &mctl.Category{
		Actions: map[string]mctl.Action{
			"list":    mctl.GroupList{c, f},
			"create":  mctl.GroupCreate{c, f},
			"update":  mctl.GroupUpdate{c, f},
			"destroy": mctl.GroupDestroy{c, f},
		},
	}
	deploy := &mctl.Category{
		Actions: map[string]mctl.Action{
			"list":   mctl.DeployList{c, f},
			"cancel": mctl.DeployCancel{c, f},
		},
	}
	marathon := &mctl.Category{
		Actions: map[string]mctl.Action{
			"leader":   mctl.MarathonLeader{c, f},
			"abdicate": mctl.MarathonAbdicate{c, f},
			"ping":     mctl.MarathonPing{c, f},
		},
	}
	artifact := &mctl.Category{
		Actions: map[string]mctl.Action{
			"upload": mctl.ArtifactUpload{c, f},
			"get":    mctl.ArtifactGet{c, f},
			"delete": mctl.ArtifactDelete{c, f},
		},
	}
	image := &mctl.Category{
		Actions: map[string]mctl.Action{
			"build":  opt.ImageBuild{},
			"upload": opt.ImageUpload{},
		},
	}
	t := &mctl.Tool{
		Selections: map[string]mctl.Selector{
			"app":      app,
			"task":     task,
			"group":    group,
			"deploy":   deploy,
			"marathon": marathon,
			"artifact": artifact,
			"image":    image,
		},
	}

	// add client action
	clientAciton(t, e, flag.Args())
}

// client actions
func clientAciton(t *mctl.Tool, e g.MarathonObj, args []string) {
	image := e.Baseinfo
	if len(args) == 0 {
		check.Usage()
	}
	act := args[0]
	if selection, ok := t.Selections[act]; !ok {
		check.Usage()
	} else {
		if act == "image" {
			if args[1] == "build" {
				if len(args[2:]) == 4 {
					image.BuildPath = args[2]
					image.Registry = args[3]
					image.Gitlib = args[4]
					image.DeployJson = args[5]
					buildImage := []string{args[1], image.BuildPath, image.Registry, image.Gitlib, image.DeployJson}
					selection.Select(buildImage)
				} else if len(args[2:]) == 0 {
					buildImage := []string{args[1], image.BuildPath, image.Registry, image.Gitlib, image.DeployJson}
					selection.Select(buildImage)
				} else {
					fmt.Println("image build : arguments error")
					return
				}
			} else if args[1] == "upload" {
				if len(args[2:]) == 2 {
					image.BuildPath = args[2]
					image.Registry = args[3]
					uploadImage := []string{args[1], image.BuildPath, image.Registry}
					selection.Select(uploadImage)
				} else if len(args[2:]) == 0 {
					uploadImage := []string{args[1], image.BuildPath, image.Registry}
					selection.Select(uploadImage)
				} else {
					fmt.Println("image upload : arguments error")
					return
				}
				//fmt.Println(image)
			} else {
				fmt.Println("not support : image " + args[1])
			}
		} else {
			selection.Select(args[1:])
		}
	}
}
