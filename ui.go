package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func PromptYesNo(msg string) (bool, error) {
	ans := false
	pmt := &survey.Confirm{
		Message: msg,
	}
	err := survey.AskOne(pmt, &ans)
	if err != nil {
		return false, fmt.Errorf("prompt: %w", err)
	}
	return ans, nil
}

func SelectFromList(msg string, ls []string) ([]string, error) {
	ans := []string{}
	pmt := &survey.MultiSelect{
		Message: msg,
		Options: ls,
	}
	err := survey.AskOne(pmt, &ans)
	if err != nil {
		return nil, fmt.Errorf("select from list: %w", err)
	}
	return ans, nil
}

func ReadSecret(msg string) (string, error) {
	ans := ""
	pmt := &survey.Password{
		Message: msg,
	}
	err := survey.AskOne(pmt, &ans)
	if err != nil {
		return "", fmt.Errorf("read secret: %w", err)
	}
	return ans, nil
}

func ReadInput(msg string) (string, error) {
	ans := ""
	pmt := &survey.Input{
		Message: msg,
	}
	err := survey.AskOne(pmt, &ans)
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}
	return ans, nil
}
