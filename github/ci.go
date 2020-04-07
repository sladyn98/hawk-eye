package github

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)


func  getBranchCIStatus(projectName string, owner string, branch string) (string, error) {
	return getCIStatus(projectName, owner, branch, "")
}

func  getTagCIStatus(projectName string, owner string, tag string) (string, error) {
	return getCIStatus(projectName, owner, "", tag)
}

func  getCIStatus(projectName string, owner string, branch string, tag string) (string, error) {


	var url string
	if branch == "" {
	url = fmt.Sprintf("%s/repos/%s/%s/commits/%s/status", githubV3Url, owner, projectName, tag)
	}

	if tag == "" {
		url = fmt.Sprintf("%s/repos/%s/%s/commits/%s/status", githubV3Url, owner, projectName, branch)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}


	client := &http.Client{
		Timeout: defaultTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	data, _ := ioutil.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	fmt.Println(string(data))
	CiStatus:= struct {
		State string `json:"state"`
	}{}

	err = json.Unmarshal(data, &CiStatus)
	if err != nil {
		return "", err
	}

	return CiStatus.State, nil
}

// TO-DO Support getPRCIStatus