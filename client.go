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

func promptYesNo() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to continue [y/N]? :")
	yn, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("prompt: %w", err)
	}
	return len(yn) > 0 && strings.ToLower(yn)[0] == 'y', nil
}

func (c *Client) DeleteVariables(ctx context.Context, vars []string) error {
	pv, err := c.ci.Projects.ListVariables(ctx, c.projectSlug)
	if err != nil {
		return fmt.Errorf("delete vars: %w", err)
	}
	if pv.NextPageToken != "" {
		logrus.Warn("Warning! Not all variables are listed.")
	}
	curs := make(map[string]string, len(pv.Items))
	for _, v := range pv.Items {
		curs[v.Name] = v.Value
	}

	dels := make([]string, 0)
	for _, v := range vars {
		val, prs := curs[v]
		if !prs {
			continue
		}
		dels = append(dels, v)
		fmt.Printf("%s=%s\n", v, val)
	}
	if len(dels) == 0 {
		fmt.Println("There are no deleted variables.")
		return nil
	}
	fmt.Println("These variables will be removed.")
	yes, err := promptYesNo()
	if err != nil {
		return fmt.Errorf("delete vars: %w", err)
	}
	if !yes {
		fmt.Println("Cancelled.")
		return nil
	}

	for _, v := range dels {
		if err := c.ci.Projects.DeleteVariable(ctx, c.projectSlug, v); err != nil {
			logrus.WithField("key", v).Errorf("Failed to delete: %w\n", err)
		} else {
			fmt.Printf("Deleted: %s\n", v)
		}
	}
	return nil
}

func (c *Client) UpdateOrCreateVariable(ctx context.Context, key string, val string) error {
	pv, err := c.ci.Projects.CreateVariable(ctx, c.projectSlug, circleci.ProjectCreateVariableOptions{
		Name:  &key,
		Value: &val,
	})
	if err != nil {
		return fmt.Errorf("update or create variable for key=%s: %w", key, err)
	}
	logrus.WithFields(logrus.Fields{
		"key":   pv.Name,
		"value": pv.Value,
	}).Info("created")
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
	for _, v := range vars.Items {
		fmt.Printf("%s=%s\n", v.Name, v.Value)
	}
	if vars.NextPageToken != "" {
		logrus.WithField("NextPageToken", vars.NextPageToken).Warn("Not all values are displayed")
	}
	return nil
}
