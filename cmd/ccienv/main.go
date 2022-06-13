package main

import (
	"context"
	"fmt"
	"os"

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
	Name  string `arg:"" name:"name" help:"Environment variable name to be add."`
	Value string `arg:"" name:"value" help:"Environment variable value to be add."`
}

func (l *AddCmd) Run(c *Context) error {
	return c.client.UpdateOrCreateVariable(c.ctx, l.Name, l.Value)
}

var cmd struct {
	Repo string `required:"" short:"r" help:"Set your target repository name."`

	Rm  RmCmd  `cmd:"" help:"Remove environment variables. Either environment variables or the interactive flag must be specified."`
	Ls  LsCmd  `cmd:"" help:"List environment variables."`
	Add AddCmd `cmd:"" help:"Add an environment variable."`
}

func handleErr(err error) {
	// FIXME: introduce error type
	if err != nil {
		logrus.WithField("error", err).Error("Internal error occured.")
		os.Exit(1)
	}
}

func constructProjectSlug(org string, repo string) string {
	return fmt.Sprintf("gh/%s/%s", org, repo)
}

func main() {
	cfg, err := cli.SetConfigFromEnv()
	handleErr(err)
	kc := kong.Parse(&cmd)

	slug := constructProjectSlug(cfg.CircleciOrganizationName, cmd.Repo)
	client, err := cli.NewClient(cfg, slug)
	handleErr(err)

	ctx := context.Background()
	err = kc.Run(&Context{ctx, client})
	handleErr(err)
}
