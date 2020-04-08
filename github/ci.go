package github

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func  getBranchCIStatus(projectName string, owner string, branch string) (string, error) {
	return getCIStatus(projectName, owner, branch, "", "")
}

func  getTagCIStatus(projectName string, owner string, tag string) (string, error) {
	return getCIStatus(projectName, owner, "", tag, "")
}

func  getCIStatus(projectName string, owner string, branch string, tag string, prNumber string) (string, error) {

	var url string
	if tag != "" {
		url = fmt.Sprintf("%s/repos/%s/%s/commits/%s/status", githubV3Url, owner, projectName, tag)
	}

	if branch != "" {
		url = fmt.Sprintf("%s/repos/%s/%s/commits/%s/status", githubV3Url, owner, projectName, branch)
	}

	if prNumber != "" {
		sha, err:= getSHAfromPR(owner,projectName,prNumber)
		if err != nil {
			return "", err
		}
		url = fmt.Sprintf("%s/repos/%s/%s/commits/%s/status", githubV3Url, owner, projectName, sha)
	}

	fmt.Println("request url is", url)
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

func getSHAfromPR(owner string, projectName string, prNumber string)(string, error) {
	
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%s", githubV3Url, owner, projectName, prNumber)
	fmt.Println("Request url is", url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github.sailor-v-preview+json")
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

	fmt.Println("Data is", string(data))

	pullRequestSHA := struct {
	 	SHA string `json:"merge_commit_sha"`
	}{}

	err = json.Unmarshal(data, &pullRequestSHA)
	if err != nil {
		return "", err
	}
	fmt.Println("SHA is............", pullRequestSHA.SHA)
	return pullRequestSHA.SHA, nil
}