package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
	cli "github.com/threepipes/circleci-env"
	command "github.com/threepipes/circleci-env/commands"
)

type RmCmd struct {
	Envs        []string `arg:"" optional:"" name:"env_name" help:"Environment variable names to remove."`
	Interactive bool     `optional:"" name:"interactive" short:"i" help:"Launch interactive removal mode."`
}

func (r *RmCmd) Run(c *command.Context) error {
	if len(r.Envs) > 0 && r.Interactive {
		fmt.Println("InvalidArgumentError: Do not specify both args `envs` and `-i, --interactive` in `rm` command.")
		return nil
	}
	if r.Interactive {
		return c.Client.DeleteVariablesInteractive(c.Ctx)
	} else {
		if len(r.Envs) == 0 {
			fmt.Println("InvalidArgumentError: Please specify at least one environment variable or set `-i`.")
			return nil
		}
		return c.Client.DeleteVariables(c.Ctx, r.Envs)
	}
}

type LsCmd struct {
}

func (l *LsCmd) Run(c *command.Context) error {
	return c.Client.ListVariables(c.Ctx)
}

type AddCmd struct {
	Name  string `arg:"" name:"name" help:"An environment variable name to be added."`
	Value string `arg:"" name:"value" help:"An environment variable value to be added."`
}

func (l *AddCmd) Run(c *command.Context) error {
	return c.Client.UpdateOrCreateVariable(c.Ctx, l.Name, l.Value)
}

var cmd struct {
	Org  string `short:"o" help:"Set your CircleCI organization name. If not specified, the default value is used."`
	Repo string `short:"r" help:"Set your target repository name. If not specified, the current directory name is used."`

	Rm  RmCmd  `cmd:"" help:"Remove environment variables. Either environment variables or the interactive flag must be specified."`
	Ls  LsCmd  `cmd:"" help:"List environment variables."`
	Add AddCmd `cmd:"" help:"Add an environment variable."`

	Config  command.ConfigCmd  `cmd:"" help:"Commands for ccienv configurations."`
	Project command.ProjectCmd `cmd:"" help:"Commands for CircleCI projects."`
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
	kc := kong.Parse(&cmd)
	cfg, err := cli.SetConfig()
	handleErr(err)

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
	err = kc.Run(&command.Context{
		Ctx:    ctx,
		Client: client,
	})
	handleErr(err)
}
