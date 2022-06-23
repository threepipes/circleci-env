package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
	cli "github.com/threepipes/circleci-env"
)

type Context struct {
	ctx    context.Context
	client *cli.Client
}

type RmCmd struct {
	Envs        []string `arg:"" optional:"" name:"env_name" help:"Environment variable names to remove."`
	Interactive bool     `optional:"" name:"interactive" short:"i" help:"Launch interactive removal mode."`
}

func (r *RmCmd) Run(c *Context) error {
	if len(r.Envs) > 0 && r.Interactive {
		fmt.Println("InvalidArgumentError: Do not specify both args `envs` and `-i, --interactive` in `rm` command.")
		return nil
	}
	if r.Interactive {
		return c.client.DeleteVariablesInteractive(c.ctx)
	} else {
		if len(r.Envs) == 0 {
			fmt.Println("InvalidArgumentError: Please specify at least one environment variable or set `-i`.")
			return nil
		}
		return c.client.DeleteVariables(c.ctx, r.Envs)
	}
}

type LsCmd struct {
}

func (l *LsCmd) Run(c *Context) error {
	return c.client.ListVariables(c.ctx)
}

type AddCmd struct {
	Name  string `arg:"" name:"name" help:"An environment variable name to be added."`
	Value string `arg:"" name:"value" help:"An environment variable value to be added."`
}

func (l *AddCmd) Run(c *Context) error {
	return c.client.UpdateOrCreateVariable(c.ctx, l.Name, l.Value)
}

type InitConfig struct {
}

func (l *InitConfig) Run(c *Context) error {
	org, err := cli.ReadInput("Please set your default GitHub organization: ")
	if err != nil {
		return err
	}
	token, err := cli.ReadSecret("Please set your personal API token: ")
	if err != nil {
		return err
	}
	cfg := cli.Config{
		OrganizationName: org,
		ApiToken:         token,
	}
	return cli.WriteConfig(&cfg)
}

var cmd struct {
	Org  string `short:"o" help:"Set your CircleCI organization name. If not specified, the default value is used."`
	Repo string `short:"r" help:"Set your target repository name. If not specified, the current directory name is used."`

	Rm     RmCmd      `cmd:"" help:"Remove environment variables. Either environment variables or the interactive flag must be specified."`
	Ls     LsCmd      `cmd:"" help:"List environment variables."`
	Add    AddCmd     `cmd:"" help:"Add an environment variable."`
	Config InitConfig `cmd:"" help:"Initialize ccienv configurations"`
}

func handleErr(err error) {
	// FIXME: introduce error type
	if err != nil {
		logrus.WithField("error", err).Error("Internal error occured.")
		os.Exit(1)
	}
}

func getDefaultRepoName() string {
	path, err := os.Getwd()
	handleErr(err)
	return filepath.Base(path)
}

func constructProjectSlug(org string, repo string) string {
	return fmt.Sprintf("gh/%s/%s", org, repo)
}

func main() {
	cfg, err := cli.SetConfig()
	handleErr(err)
	kc := kong.Parse(&cmd)

	repo := cmd.Repo
	if repo == "" {
		repo = getDefaultRepoName()
	}
	org := cmd.Org
	if org == "" {
		org = cfg.OrganizationName
	}

	slug := constructProjectSlug(org, repo)
	client, err := cli.NewClient(cfg, slug)
	handleErr(err)

	ctx := context.Background()
	err = kc.Run(&Context{ctx, client})
	handleErr(err)
}
