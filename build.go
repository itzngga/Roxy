package main

import (
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

func main() {
	os := runtime.GOOS
	switch os {
	case "windows":
		RunExec("go build ./main.go")
		RunExec("pm2 delete goRoxy")
		RunExec("pm2 start ./main.exe --name BOT")
	case "darwin":
		if !isRoot() {
			panic("Should be root runtime")
		}
		RunExec(`go build main.go`)
		RunExec("sudo chmod +x main")
		RunExec("pm2 delete goRoxy")
		RunExec("pm2 start ./main --name BOT")
	case "linux":
		if !isRoot() {
			panic("Should be root runtime")
		}
		RunExec(`go build main.go`)
		RunExec("sudo chmod +x main")
		RunExec("pm2 delete goRoxy")
		RunExec("pm2 start ./main --name BOT")
	default:
		panic("unsupported platform")
	}
}

func RunExec(str string) {
	args := strings.Split(str, " ")
	cmd := exec.Command(args[0], args[1:]...)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(string(stdout))
}

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}
