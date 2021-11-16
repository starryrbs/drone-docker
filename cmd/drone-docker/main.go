package main

import (
	"github.com/sirupsen/logrus"
	docker "github.com/starryrbs/drone-docker"
	"github.com/urfave/cli"
	"os"
)

func main() {
	logrus.Println("start app...")

	app := cli.NewApp()
	app.Name = "drone docker plugin"
	app.Usage = "drone docker plugin"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "docker.registry",
			EnvVar: "PLUGIN_REGISTRY",
		},
		cli.StringFlag{
			Name:   "docker.username",
			EnvVar: "PLUGIN_USERNAME",
		},
		cli.StringFlag{
			Name:   "docker.password",
			EnvVar: "PLUGIN_PASSWORD",
		},
		cli.StringFlag{
			Name:   "docker.config",
			EnvVar: "PLUGIN_CONFIG",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
			Value:  "00000000",
		},
		cli.StringFlag{
			Name:   "dockerfile",
			EnvVar: "PLUGIN_DOCKERFILE",
		},
		cli.StringFlag{
			Name:   "context",
			EnvVar: "PLUGIN_CONTEXT",
		},
		cli.StringSliceFlag{
			Name:   "tags",
			Usage:  "build tags",
			Value:  &cli.StringSlice{"latest"},
			EnvVar: "PLUGIN_TAGS",
		},
		cli.StringFlag{
			Name:   "repo",
			EnvVar: "PLUGIN_REPO",
		},
	}

	logrus.Println("PLUGIN_REGISTRY", os.Getenv("PLUGIN_REGISTRY"))

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := docker.Plugin{
		Login: docker.Login{
			Registry: c.String("docker.registry"),
			Username: c.String("docker.username"),
			Password: c.String("docker.password"),
			Config:   c.String("docker.config"),
		},
		Build: docker.Build{
			Name:       c.String("commit.sha"),
			Dockerfile: c.String("dockerfile"),
			Context:    c.String("context"),
			Tags:       c.StringSlice("tags"),
			Repo:       c.String("repo"),
		},
	}
	return plugin.Exec()
}
