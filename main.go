package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Jeffail/gabs/v2"
	"github.com/progrium/go-shell"
)

var sh = shell.Run

type PackageJson struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Keywords        []string          `json:"keywords"`
	Homepage        string            `json:"homepage"`
	License         string            `json:"license"`
	Files           []string          `json:"files"`
	Main            string            `json:"main"`
	Scripts         map[string]string `json:"scripts"`
	Os              []string          `json:"os"`
	Cpu             []string          `json:"cpu"`
	Private         bool              `json:"private"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func ConfigGit() {
	if _, err := exec.Command("git", "config", "--global", "user.name", "fe-pub-bot").Output(); err != nil {
		panic("Config user.name error")
	}

	if _, err := exec.Command("git", "config", "--global", "user.email", "1121292341@qq.com").Output(); err != nil {
		panic("Config user.email error")
	}
}

func UpdateRepo(name string) {
	// gitUrl := fmt.Sprintf("http://oauth2:CdVcbeg21xv8PuJ48exN@runafe.cn:8088/wangxd/%s.git", name)
	gitUrl := fmt.Sprintf("http://pub-bot:6Q4ybPpgdKPxS6azN_n1@runafe.cn:8088/wxdtest/%s.git", name)
	os.Chdir("./repositories")
	if info, err := os.Stat(name); err == nil {
		name := info.Name()
		if err := os.RemoveAll(name); err != nil {
			fmt.Print(err)
			panic(err)
		}
	}
	sh("git", "clone", gitUrl)
	jsonObj, err := gabs.ParseJSONFile(fmt.Sprintf("./%s/package.json", name))
	if err != nil {
		panic(err)
	}
	for key := range jsonObj.S("dependencies").ChildrenMap() {
		if key == "@runafe/runa-system" {
			jsonObj.Set("2.0.1", "dependencies", key)
		}
	}
	// TODO can not sort the key
	// https://github.com/golang/go/issues/27179
	// https://github.com/golang/go/issues/6244
	os.WriteFile(fmt.Sprintf("./%s/package.json", name), jsonObj.EncodeJSON(gabs.EncodeOptHTMLEscape(false), gabs.EncodeOptIndent("", "  ")), 7770)

	os.Chdir(name)
	sh("git", "commit", "-am", "'chore: update package.json'")
	sh("git", "push")
	sh("echo Done")
}

func main() {
	defer shell.ErrExit()
	// data, err := ioutil.ReadFile("./test/package.json")
	// if err != nil {
	// 	panic("Read package file failed")
	// }

	// var pkg PackageJson

	// err = json.Unmarshal(data, &pkg)

	// if err != nil {
	// 	panic("Parse JSON failed")
	// }

	// fmt.Println(pkg.Dependencies["lodash"])

	UpdateRepo("test")
}
