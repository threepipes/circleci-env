package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	cli "github.com/threepipes/circleci-env"
)

type Context struct {
	ctx    context.Context
	client *cli.Client
}

type RmCmd struct {
	Envs []string `arg:"" name:"env" help:"Environment variables to remove."`
}

func (r *RmCmd) Run(c *Context) error {
	c.client.DeleteVariables(c.ctx, r.Envs)
	return nil
}

type LsCmd struct {
}

func (l *LsCmd) Run(c *Context) error {
	c.client.ListVariables(c.ctx)
	return nil
}

type AddCmd struct {
	Name  string `arg:"" name:"name" help:"Environment variable name to be add."`
	Value string `arg:"" name:"value" help:"Environment variable value to be add."`
}

func (l *AddCmd) Run(c *Context) error {
	c.client.UpdateOrCreateVariable(c.ctx, l.Name, l.Value)
	return nil
}

var cmd struct {
	Repo string `required:"" help:"Set your target repository name."`

	Rm  RmCmd  `cmd:"" help:"Remove environment variables."`
	Ls  LsCmd  `cmd:"" help:"List environment variables."`
	Add AddCmd `cmd:"" help:"Add an environment variable."`
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
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
	kc.FatalIfErrorf(err)
}
