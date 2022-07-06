package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

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
	client, err := c.ClientGenerator()
	handleErr(err)
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

func (l *LsCmd) Run(c *command.Context) error {
	client, err := c.ClientGenerator()
	handleErr(err)
	return client.ListVariables(c.Ctx)
}

type AddCmd struct {
	Name  string `arg:"" name:"name" help:"An environment variable name to be added."`
	Value string `arg:"" name:"value" help:"An environment variable value to be added."`
}

func (l *AddCmd) Run(c *command.Context) error {
	client, err := c.ClientGenerator()
	handleErr(err)
	return client.UpdateOrCreateVariable(c.Ctx, l.Name, l.Value)
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

func extractRepoName(repo string) (string, string, error) {
	r := regexp.MustCompile(`.+github\.com[:/]([^/]+)/(.+)\.git/?$`)
	match := r.FindStringSubmatch(repo)
	if len(match) < 3 {
		return "", "", fmt.Errorf("failed to parse repo name: %v", repo)
	}
	return match[2], match[1], nil
}

func getDefaultRepoName() (string, string) {
	var stderr bytes.Buffer
	cmd := exec.Command("git", strings.Split("config --get remote.origin.url", " ")...)
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		logrus.Error("Failed to read git remote repository with git command. Please specify the repository by the `-r` option or go to the directory where .git is with git command.")
		handleErr(err)
	}
	repo := strings.TrimSpace(string(out))
	rn, org, err := extractRepoName(repo)
	if err != nil {
		handleErr(err)
	}
	return rn, org
}

func constructProjectSlug(org string, repo string) string {
	return fmt.Sprintf("gh/%s/%s", org, repo)
}

func getClient() (*cli.Client, error) {
	cfg, err := cli.ReadConfig()
	if err != nil {
		return nil, err
	}

	org := cmd.Org
	if org == "" {
		org = cfg.OrganizationName
	}
	repo := cmd.Repo
	if repo == "" {
		repo, org = getDefaultRepoName()
	}

	slug := constructProjectSlug(org, repo)
	client, err := cli.NewClient(cfg, slug)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func mainRun() {
	kc := kong.Parse(&cmd)

	ctx := context.Background()
	err := kc.Run(&command.Context{
		Ctx:             ctx,
		ClientGenerator: getClient,
	})
	handleErr(err)
}

func main() {
	mainRun()
}
