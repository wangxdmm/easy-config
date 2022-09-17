package main

import (
	"fmt"
	"os"
	"time"

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

func UpdateRepo(name string, pkgMap map[string]string, exctPath string, done func()) {
	gitUrl := fmt.Sprintf("http://oauth2:svazkrYzqkMZeWaKx86b@runafe.cn:8088/v4/%s.git", name)
	// gitUrl := fmt.Sprintf("http://oauth2:d6n9LvaWsoZazzQFx4hV@wangxd.cn:8088/wxdtest/%s.git", name)

	sh("git", "clone", gitUrl, "--branch", "dev", "--single-branch")
	jsonObj, err := gabs.ParseJSONFile(fmt.Sprintf("%s/package.json", exctPath))
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
	os.WriteFile(fmt.Sprintf("%s/package.json", exctPath), jsonObj.EncodeJSON(gabs.EncodeOptHTMLEscape(false), gabs.EncodeOptIndent("", "  ")), 7770)

	os.Chdir(exctPath)
	sh("git", "commit", "-am", "'chore: update package.json'")
	// sh("git", "push")
	// sh("echo Done")

	fmt.Printf("%s done\n", name)
	done()
}

func main() {
	defer shell.ErrExit()
	start := time.Now()
	repos := []string{"web-monitor", "web-maintain", "web-ai", "web-charge", "web-AIScheduler", "web-analysis", "web-balance", "web-dashboard", "web-dataCenter", "web-device", "web-dispatchCommand", "web-dispatchplatform", "web-emerg-conduct", "web-gis", "web-giscience", "web-meter", "web-powermonitor", "web-temAnalysis", "web-user", "web-wechat", "web-wechat2", "web-wechatCharge"}

	updatedPkgMap := map[string]string{
		"@runafe/runa-system": "2.0.5-beta.10",
		"dayjs":               "1.0.1",
	}

	pwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	os.Chdir(pwd)
	repo := "repositories"

	if ok, _ := PathExists(repo); !ok {
		os.Mkdir(repo, 0777)
	}

	os.Chdir(repo)

	for _, name := range repos {
		if _, err := os.Stat(name); err == nil {
			if err := os.RemoveAll(name); err != nil {
				fmt.Print(err)
				panic(err)
			}
		}
		exactPath := fmt.Sprintf("%s/%s/%s", pwd, repo, name)
		fmt.Println(exactPath)
		UpdateRepo(name, updatedPkgMap, exactPath, func() {
		})
	}

	// var wg sync.WaitGroup
	// wg.Add(len(repos))
	// for _, name := range repos {
	// 	if _, err := os.Stat(name); err == nil {
	// 		if err := os.RemoveAll(name); err != nil {
	// 			fmt.Print(err)
	// 			panic(err)
	// 		}
	// 	}
	// 	exactPath := fmt.Sprintf("%s/%s/%s", pwd, repo, name)
	// 	fmt.Println(exactPath)
	// 	go UpdateRepo(name, updatedPkgMap, exactPath, func() {
	// 		wg.Done()
	// 	})
	// }
	// wg.Wait()

	fmt.Println("Verify:", time.Since(start))
}
