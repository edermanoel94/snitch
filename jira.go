package main

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
)

// y3ZbzrLzmIbFd598bWwa8C82

type JiraCredentials struct {
	Host          string
	Username      string
	PersonalToken string
}

func (cred JiraCredentials) getIssue(repo string, todo Todo) (map[string]interface{}, error) {

	request, err := http.NewRequest(http.MethodGet, "", nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	request.SetBasicAuth(cred.Username, cred.PersonalToken)

	return QueryHTTP(request)
}

func (cred JiraCredentials) postIssue(repo string, todo Todo, body string) (Todo, error) {

	request, err := http.NewRequest(http.MethodGet, "", nil)

	if err != nil {
		return Todo{}, err
	}

	request.Header.Set("Content-Type", "application/json")

	request.SetBasicAuth(cred.Username, cred.PersonalToken)

	return Todo{}, nil
}

func (cred JiraCredentials) getHost() string {
	return cred.Host
}

/*func JiraCredentialsFromSession(ctx context.Context, username, password string) (JiraCredentials, error) {

	reader := strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "", reader)

	if err != nil {
		return JiraCredentials{}, err
	}

	client := http.DefaultClient

	response, err := client.Do(req)

	if err != nil {
		return JiraCredentials{}, err
	}

	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)

	return JiraCredentials{}, nil
}*/

func JiraCredentialsFromToken(credentials string) (JiraCredentials, error) {

	return JiraCredentials{}, nil
}

func getJiraCredentials(creds []IssueAPI) []IssueAPI {

	tokenEnvar := os.Getenv("JIRA_PERSONAL_TOKEN")
	xdgEnvar := os.Getenv("XDG_CONFIG_HOME")
	usr, err := user.Current()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(tokenEnvar) != 0 {
		for _, credential := range strings.Split(tokenEnvar, ",") {
			parsedCredentials, err := JiraCredentialsFromToken(credential)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			creds = append(creds, parsedCredentials)
		}
	}

	// custom XDG_CONFIG_HOME
	if len(xdgEnvar) != 0 {
		filePath := path.Join(xdgEnvar, "snitch/gitlab.ini")
		if _, err := os.Stat(filePath); err == nil {
			for _, cred := range GitlabCredentialsFromFile(filePath) {
				creds = append(creds, cred)
			}
		}
	}

	// default XDG_CONFIG_HOME
	if len(xdgEnvar) == 0 {
		filePath := path.Join(usr.HomeDir, ".config/snitch/gitlab.ini")
		if _, err := os.Stat(filePath); err == nil {
			for _, cred := range GitlabCredentialsFromFile(filePath) {
				creds = append(creds, cred)
			}
		}
	}

	filePath := path.Join(usr.HomeDir, ".snitch/gitlab.ini")
	if _, err := os.Stat(filePath); err == nil {
		for _, cred := range GitlabCredentialsFromFile(filePath) {
			creds = append(creds, cred)
		}
	}

	return creds
}
