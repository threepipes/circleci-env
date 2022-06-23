package command

type ProjectCmd struct {
	Show ProjectShowCmd `cmd:"" help:"Show the project information for the repository."`
}

type ProjectShowCmd struct {
}

func (p *ProjectShowCmd) Run(c *Context) error {
	return c.Client.ShowProject(c.Ctx)
}
