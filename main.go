package main

import (
	"net/http"
	"os"
	"io"
	"archive/zip"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"flag"
	"strings"
	"time"
)

// ko微服务项目cli

// - 帮助信息
// - 安装项目：网关，服务
// - 增加网关服务
// - 增加服务接口

func main() {

	installCommand := flag.NewFlagSet("install", flag.ExitOnError)

	if os.Args[1] == "help" || os.Args[1] == "--help" || os.Args[1] == "--h" || os.Args[1] == "-h" {
		printHelpMenu("main")
		return
	}

	if os.Args[1] == "version" || os.Args[1] == "-V" || os.Args[1] == "--v" || os.Args[1] == "--version"{
		printVersion()
		return
	}

	switch os.Args[1] {
	case "install":
		if len(os.Args) > 2 && os.Args[2] == "help" {
			printHelpMenu("install")
			return
		} else {
			var name string
			var pjtype string
			installCommand.StringVar(&name, "name", "ko", "project name")
			installCommand.StringVar(&pjtype, "type", "0", "project type")
			installCommand.Parse(os.Args[2:])
			install(pjtype, name)
			return
		}
	default:
		fmt.Printf("%s", "Invaild command. \n\n")
		printHelpMenu("main")
		return
	}
}

func printVersion()  {
	fmt.Printf("%s", "ko version: 0.0.1\n")
	fmt.Printf("%s", "ko-cli version: 0.0.1\n")
}

func printHelpMenu(command string) {
	if command == "main" {
		fmt.Printf("%s", "Ko 0.0.1\n\n")

		fmt.Printf("%s", "Usage:\n")
		fmt.Printf("%s", "	command [options] [arguments]\n")

		fmt.Printf("%s", "Options:\n")
		fmt.Printf("%s", "	-h, --help            Display this help message\n")
		fmt.Printf("%s", "	-q, --quiet           Do not output any message\n")
		fmt.Printf("%s", "	-V, --version         Display this application version\n")

		fmt.Printf("%s", "Available commands:\n")
		fmt.Printf("%s", "	addsvc               Remove the compiled class file\n")
		fmt.Printf("%s", "	install              Put the application into maintenance mode\n")
	}
	if command == "install" {
		fmt.Printf("%s", "Ko 0.0.1\n\n")

		fmt.Printf("%s", "Usage:\n")
		fmt.Printf("%s", "	install [options] [arguments]\n")

		fmt.Printf("%s", "Options:\n")
		fmt.Printf("%s", "	--name            project name\n")
		fmt.Printf("%s", "	--type            project type, 0 gateway, 1 service\n")
	}

}

func install(pjtype string, name string) {
	go func() {
		fmt.Printf("%s", "[")
		consoleStr := "█"
		for i := 0; i != 10; i = i + 1 {
			//log.Println(consoleStr)
			fmt.Printf("%s", consoleStr)
			time.Sleep(time.Second * 1)
		}
	}()

	url := "https://api.github.com/repos/chenhg5/ko/zipball/master"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/vnd.github.v3+json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	file, _ := os.Create("tmp.zip")
	io.Copy(file, res.Body)

	unzipDir("tmp.zip", "tmp")

	files, _ := ioutil.ReadDir("./tmp")

	if pjtype == "0" {
		// 网关
		os.Rename("./tmp/" + files[0].Name() + "/gateway", "./" + name)
		os.RemoveAll("tmp")
		os.Remove("tmp.zip")
		renameProject( "./" + name, name, "ko/gateway")
	} else {
		// 服务
		os.Rename("./tmp/" + files[0].Name() + "/services", "./" + name)
		os.RemoveAll("tmp")
		os.Remove("tmp.zip")
		os.Rename("./" + files[0].Name(), "./" + name)

		renameProject( "./" + name, name, "ko/services")
	}

	fmt.Printf("%s", "] 100% \ninstall ok!\n\n")
}

func unzipDir(zipFile, dir string) {

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		//log.Fatalf("Open zip file failed: %s\n", err.Error())
	}
	defer r.Close()

	for _, f := range r.File {
		func() {
			path := dir + string(filepath.Separator) + f.Name
			os.MkdirAll(filepath.Dir(path), 0755)
			fDest, err := os.Create(path)
			if err != nil {
				//log.Printf("Create failed: %s\n", err.Error())
				return
			}
			defer fDest.Close()

			fSrc, err := f.Open()
			if err != nil {
				//log.Printf("Open failed: %s\n", err.Error())
				return
			}
			defer fSrc.Close()

			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				//log.Printf("Copy failed: %s\n", err.Error())
				return
			}
		}()
	}
}

func renameProject(fileDir string, projectName string, oldpath string)  {
	//fmt.Println("path: " +  fileDir)
	files, _ := ioutil.ReadDir(fileDir)
	for _,file := range files{
		if file.IsDir(){
			renameProject(fileDir + "/" + file.Name(), projectName, oldpath)
		} else {
			path := fileDir + "/" + file.Name()
			//fmt.Println("replace path: " +  path)
			buf, _ := ioutil.ReadFile(path)
			content := string(buf)

			//替换
			newContent := strings.Replace(content, oldpath + "/", projectName + "/", -1)
			newContent = strings.Replace(newContent, "package services", "package main", -1)
			newContent = strings.Replace(newContent, "services.", projectName + ".", -1)
			newContent = strings.Replace(newContent, oldpath, projectName, -1)

			//重新写入
			ioutil.WriteFile(path, []byte(newContent), 0)
		}
	}
}