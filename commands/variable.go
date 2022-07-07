package command

import "fmt"

type RmCmd struct {
	Envs        []string `arg:"" optional:"" name:"env_name" help:"Environment variable names to remove."`
	Interactive bool     `optional:"" name:"interactive" short:"i" help:"Launch interactive removal mode."`
}

func (r *RmCmd) Run(c *Context) error {
	client, err := c.ClientGenerator()
	if err != nil {
		return fmt.Errorf("rm command: %w", err)
	}
	if len(r.Envs) > 0 && r.Interactive {
		fmt.Println("InvalidArgumentError: Do not specify both args `envs` and `-i, --interactive` in `rm` command.")
		return nil
	}
	if r.Interactive {
		return client.DeleteVariablesInteractive(c.Ctx)
	} else {
		if len(r.Envs) == 0 {
			fmt.Println("InvalidArgumentError: Please specify at least one environment variable or set `-i`.")
			return nil
		}
		return client.DeleteVariables(c.Ctx, r.Envs)
	}
}

type LsCmd struct {
}

func (l *LsCmd) Run(c *Context) error {
	client, err := c.ClientGenerator()
	if err != nil {
		return fmt.Errorf("ls command: %w", err)
	}
	return client.ListVariables(c.Ctx)
}

type AddCmd struct {
	Name  string `arg:"" name:"name" help:"An environment variable name to be added."`
	Value string `arg:"" name:"value" help:"An environment variable value to be added."`
}

func (l *AddCmd) Run(c *Context) error {
	client, err := c.ClientGenerator()
	if err != nil {
		return fmt.Errorf("add command: %w", err)
	}
	return client.UpdateOrCreateVariable(c.Ctx, l.Name, l.Value)
}
