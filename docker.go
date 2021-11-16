package docker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"time"
)

type (
	Login struct {
		Registry string // Docker registry address
		Username string // Docker registry username
		Password string // Docker registry password
		Config   string // Docker Auth Config
	}

	Build struct {
		Name       string   // Docker build Name
		Dockerfile string   // Docker build Dockerfile
		Context    string   // Docker build context
		Tags       []string // Docker build tags
		Repo       string   // Docker build repository
	}

	Plugin struct {
		Login Login
		Build Build
	}
)

func formatCmdError(cmd *exec.Cmd, message string, err error) string {
	return fmt.Sprintf("[command]:%s,[message]:%s,error:%s", cmd.Args[:2], message, err)
}

func (p *Plugin) Exec() (err error) {
	for i := 0; ; i++ {
		cmd := commandInfo()
		logrus.Info(cmd)
		err := cmd.Run()
		if err == nil {
			break
		}
		if i == 15 {
			logrus.Error("Unable to reach Docker Daemon after 15 attempts.")
			break
		}
		time.Sleep(time.Second * 1)
	}

	if p.Login.Password != "" {
		cmd := commandLogin(p.Login)
		logrus.Info(cmd)
		raw, err := cmd.CombinedOutput()
		if err != nil {
			out := string(raw)
			out = strings.Replace(out, "WARNING! Using --password via the CLI is insecure. Use --password-stdin.", "", -1)
			fmt.Println(out)
			return fmt.Errorf("Error authenticating: exit status 1")
		}
	}

	// build
	buildCmd := commandBuild(p.Build)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	logrus.Info(buildCmd)
	err = buildCmd.Run()
	if err != nil {
		print(formatCmdError(buildCmd, "build failed", err))
		return
	}

	// tag
	for _, tag := range p.Build.Tags {
		tagCmd := commandTag(p.Build, tag)
		tagCmd.Stdout = os.Stdout
		tagCmd.Stderr = os.Stderr
		logrus.Info(tagCmd)
		err = tagCmd.Run()
		if err != nil {
			print(formatCmdError(tagCmd, "build failed", err))
			return
		}
	}

	// push
	if p.Login.Registry != "" {
		for _, tag := range p.Build.Tags {
			pushCmd := commandPush(p.Build, tag)
			pushCmd.Stdout = os.Stdout
			pushCmd.Stderr = os.Stderr
			logrus.Info(pushCmd)
			err = pushCmd.Run()
			if err != nil {
				print(formatCmdError(pushCmd, "push failed", err))
				return
			}
		}
	}
	return
}

func commandInfo() *exec.Cmd {
	return exec.Command(dockerExe, "info")
}

func commandLogin(login Login) *exec.Cmd {
	return exec.Command(dockerExe, "login",
		"-u", login.Username,
		"-p", login.Password,
		login.Registry,
	)
}

func commandBuild(build Build) *exec.Cmd {
	args := []string{
		"build",
		"-f", build.Dockerfile,
		"-t", build.Name,
	}

	args = append(args, build.Context)

	return exec.Command(dockerExe, args...)
}

func commandTag(build Build, tag string) *exec.Cmd {
	target := fmt.Sprintf("%s:%s", build.Repo, tag)
	return exec.Command(dockerExe, "tag", build.Name, target)
}

func commandPush(build Build, tag string) *exec.Cmd {
	target := fmt.Sprintf("%s:%s", build.Repo, tag)
	return exec.Command(dockerExe, "push", target)
}

func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
