package command

import (
	"context"

	cli "github.com/threepipes/circleci-env"
)

type Context struct {
	Ctx             context.Context
	ClientGenerator func() (*cli.Client, error)
}
