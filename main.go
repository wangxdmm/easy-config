package main

import (
	"fmt"
	"os"

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

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func UpdateRepo(name string, pkgMap map[string]string, pwd string) {
	os.Chdir(pwd)
	// gitUrl := fmt.Sprintf("http://oauth2:CdVcbeg21xv8PuJ48exN@runafe.cn:8088/wangxd/%s.git", name)
	gitUrl := fmt.Sprintf("http://oauth2:d6n9LvaWsoZazzQFx4hV@wangxd.cn:8088/wxdtest/%s.git", name)
	repo := "./repositories"

	if ok, _ := PathExists(repo); !ok {
		os.Mkdir(repo, 0777)
	}

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
		for k, v := range pkgMap {
			if key == k {
				jsonObj.Set(v, "dependencies", key)
			}
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
	repos := []string{"test1", "test2", "test3"}
	updatedPkgMap := map[string]string{
		"@runafe/runa-system": "2.0.5-beta.1",
		"dayjs":               "1.0.1",
	}

	pwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	for _, name := range repos {
		UpdateRepo(name, updatedPkgMap, pwd)
	}
}
