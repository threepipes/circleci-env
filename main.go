package main

import (
	"context"
	"log"

	"github.com/alecthomas/kong"
)

type Context struct {
	ctx    context.Context
	client *Client
}

type RmCmd struct {
	Envs []string `arg:"" name:"env" help:"Environment variables to remove." type:"env"`
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

var cmd struct {
	Rm RmCmd `cmd:"" help:"Remove environment variables."`
	Ls LsCmd `cmd:"" help:"List environment variables."`
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	cfg, err := SetConfigFromEnv()
	handleErr(err)
	client, err := NewClient(cfg, cfg.ProjectSlug)
	handleErr(err)

	ctx := context.Background()
	kc := kong.Parse(&cmd)
	err = kc.Run(&Context{ctx, client})
	kc.FatalIfErrorf(err)
}
