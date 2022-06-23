package command

import (
	cli "github.com/threepipes/circleci-env"
)

type ConfigCmd struct {
	Init ConfigInitCmd `cmd:"" help:"Initialize ccienv configurations."`
}

type ConfigInitCmd struct {
}

func (l *ConfigInitCmd) Run(c *Context) error {
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
