package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CircleciApiToken string `split_words:"true"`
	ProjectSlug      string `split_words:"true"`
}

func setConfigFromEnv() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func listEnv(token string, projectSlug string) {
	url := fmt.Sprintf("https://circleci.com/api/v2/project/%s/envvar", projectSlug)
	fmt.Println(url)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Circle-Token", token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

func deleteEnv(token string, projectSlug string, target string) {
	url := fmt.Sprintf("https://circleci.com/api/v2/project/%s/envvar/%s", projectSlug, target)
	fmt.Println(url)

	req, _ := http.NewRequest("DELETE", url, nil)

	req.Header.Add("Circle-Token", token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	cfg, err := setConfigFromEnv()
	handleErr(err)
	listEnv(cfg.CircleciApiToken, cfg.ProjectSlug)
	// deleteEnv(cfg.CircleciApiToken, cfg.ProjectSlug, "TEST_ENV")
}
