package command

import "fmt"

type ProjectCmd struct {
	Show ProjectShowCmd `cmd:"" help:"Show the project information for the repository."`
}

type ProjectShowCmd struct {
}

func (p *ProjectShowCmd) Run(c *Context) error {
	client, err := c.ClientGenerator()
	if err != nil {
		return fmt.Errorf("show project: %w", err)
	}
	return client.ShowProject(c.Ctx)
}
