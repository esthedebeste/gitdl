package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func checkGit() (string, error) {
	gitPath, err := exec.LookPath("git")
	if err == exec.ErrNotFound {
		return "git", err
	}
	return gitPath, nil
}

func fatal(a string) {
	os.Stderr.WriteString(a)
	os.Exit(1)
}

func register(url, execute string) {
	key, _, err := registry.CreateKey(registry.CLASSES_ROOT, url, registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		fatal(err.Error())
	}
	defer key.Close()
	key.SetStringValue("", "URL:"+url)
	key.SetStringValue("URL Protocol", "")
	shell, _, err := registry.CreateKey(key, "shell", registry.CREATE_SUB_KEY)
	if err != nil {
		fatal(err.Error())
	}
	defer shell.Close()
	open, _, err := registry.CreateKey(shell, "open", registry.CREATE_SUB_KEY)
	if err != nil {
		fatal(err.Error())
	}
	defer open.Close()
	command, _, err := registry.CreateKey(open, "command", registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		fatal(err.Error())
	}
	defer command.Close()
	command.SetStringValue("", execute)
}

func main() {
	if runtime.GOOS != "windows" {
		fatal("This tool has been made for (and depends on) windows.")
	}
	gitPath, err := checkGit()
	if err != nil {
		fatal("You either don't have git, or it isn't on your path. Please install git from https://git-scm.com/download/win.")
	}
	fmt.Println("Git found!")
	executable, _ := os.Executable()

	if len(os.Args) == 0 {
		fatal("how")
	} else if len(os.Args) == 1 {
		cwd, _ := os.Getwd()
		verbPtr, _ := syscall.UTF16PtrFromString("runas")
		exePtr, _ := syscall.UTF16PtrFromString(executable)
		cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
		argPtr, _ := syscall.UTF16PtrFromString("--admined")
		err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, 2)
		if err != nil {
			fatal(err.Error())
		}
		fmt.Println("Opening UAC prompt.")
	} else {
		if os.Args[1] == "--admined" {
			register("gitdl", "\""+executable+"\" \"%1\"")
			register("gitshallowdl", "\""+executable+"\" \"%1\" --depth 1")
			return
		}
		url, err := url.Parse(os.Args[1])
		if err != nil {
			fatal(err.Error())
		}
		custom_args := os.Args[2:]
		user, _ := user.Current()
		dir := user.HomeDir + "\\Desktop\\"
		var clone *exec.Cmd
		treegex := regexp.MustCompile(`(.*\/(.*))\/tree\/[\w-]+\/?$`)
		if treegex.MatchString(url.Path) {
			dir += treegex.FindStringSubmatch(url.Path)[2]
			branch := regexp.MustCompile(`([\w-]+)/?$`).FindStringSubmatch(url.Path)[1]
			path := treegex.FindStringSubmatch(url.Path)[1]
			args := []string{"clone", "-b", branch, "--recursive", "http://" + url.Host + path, dir + "-" + branch}
			args = append(args, custom_args...)
			fmt.Println("`git", strings.Join(args, " ")+"`")
			clone = exec.Command(gitPath, args...)
			defer exec.Command("explorer.exe", dir+"-"+branch).Run()
		} else {
			pathsepped := strings.Split(strings.TrimSuffix(url.Path, "/"), "/")
			dir += pathsepped[len(pathsepped)-1]
			args := []string{"clone", "--recursive", "http://" + url.Host + url.Path, dir}
			args = append(args, custom_args...)
			fmt.Println("`git", strings.Join(args, " ")+"`")
			clone = exec.Command(gitPath, args...)
			defer exec.Command("explorer.exe", dir).Run()
		}
		clone.Stdout = os.Stdout
		clone.Stderr = os.Stderr
		err = clone.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}
