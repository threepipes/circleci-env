package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/grezar/go-circleci"
	"github.com/sirupsen/logrus"
)

type Client struct {
	ci          *circleci.Client
	projectSlug string
}

func NewClient(cfg *Config, prj string) (*Client, error) {
	config := circleci.DefaultConfig()
	config.Token = cfg.CircleciApiToken
	ci, err := circleci.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}
	return &Client{
		ci:          ci,
		projectSlug: prj,
	}, nil
}

func promptYesNo(msg string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [Y/n]? :", msg)
	yn, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("prompt: %w", err)
	}
	return len(yn) > 0 && strings.ToLower(yn)[0] == 'y', nil
}

func dumpVariables(items []*circleci.ProjectVariable) {
	maxlen := 0
	for _, v := range items {
		if len(v.Name) > maxlen {
			maxlen = len(v.Name)
		}
	}
	for _, v := range items {
		fmt.Printf("%-*s %s\n", maxlen, v.Name, v.Value)
	}
}

func getIntersection(vars []string, items []*circleci.ProjectVariable) []*circleci.ProjectVariable {
	curs := make(map[string]interface{}, 0)
	var t interface{}
	for _, v := range vars {
		curs[v] = t
	}

	dels := make([]*circleci.ProjectVariable, 0)
	for _, v := range items {
		_, prs := curs[v.Name]
		if !prs {
			continue
		}
		dels = append(dels, v)
	}
	return dels
}

func (c *Client) DeleteVariables(ctx context.Context, vars []string) error {
	pv, err := c.ci.Projects.ListVariables(ctx, c.projectSlug)
	if err != nil {
		return fmt.Errorf("delete vars: %w", err)
	}
	if pv.NextPageToken != "" {
		logrus.Warn("Warning! Not all variables are listed.")
	}

	dels := getIntersection(vars, pv.Items)
	if len(dels) == 0 {
		fmt.Println("There are no deleted variables.")
		return nil
	}

	fmt.Println("These variables will be removed.")
	fmt.Println()
	dumpVariables(dels)
	fmt.Println()

	yes, err := promptYesNo("Do you want to continue?")
	if err != nil {
		return fmt.Errorf("delete vars: %w", err)
	}
	if !yes {
		fmt.Println("Cancelled.")
		return nil
	}

	for _, v := range dels {
		if err := c.ci.Projects.DeleteVariable(ctx, c.projectSlug, v.Name); err != nil {
			logrus.WithField("key", v).Errorf("Failed to delete: %v\n", err)
		} else {
			fmt.Printf("Deleted: %s\n", v)
		}
	}
	return nil
}

func (c *Client) UpdateOrCreateVariable(ctx context.Context, key string, val string) error {
	v, err := c.ci.Projects.GetVariable(ctx, c.projectSlug, key)
	if err != nil {
		return fmt.Errorf("update or create variable for key=%s: %w", key, err)
	}
	if v != nil {
		fmt.Printf("key:%s already exists as value=%s\n", v.Name, v.Value)
		yes, err := promptYesNo("Do you want to overwrite?")
		if err != nil {
			return err
		}
		if !yes {
			fmt.Println("Cancelled.")
			return nil
		}
	}
	pv, err := c.ci.Projects.CreateVariable(ctx, c.projectSlug, circleci.ProjectCreateVariableOptions{
		Name:  &key,
		Value: &val,
	})
	if err != nil {
		return fmt.Errorf("update or create variable for key=%s: %w", key, err)
	}
	fmt.Printf("%s=%s is created\n", pv.Name, pv.Value)
	return nil
}

func (c *Client) ListVariables(ctx context.Context) error {
	/*
		TODO: use nextPageToken
		Currently, number of variables on the project is less than 50.
	*/
	vars, err := c.ci.Projects.ListVariables(ctx, c.projectSlug)
	if err != nil {
		return fmt.Errorf("list vars: %w", err)
	}
	dumpVariables(vars.Items)
	if vars.NextPageToken != "" {
		logrus.WithField("NextPageToken", vars.NextPageToken).Warn("Not all values are displayed")
	}
	return nil
}
