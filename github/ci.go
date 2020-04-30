package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetBranchCIStatus gets the CI status using branch name
func GetBranchCIStatus(githubURLEndpoint string, projectName string, owner string, branch string) (string, error) {
	return getCIStatus(githubURLEndpoint, projectName, owner, branch, "", "")
}
// GetTagCIStatus gets the CI status using tag name
func GetTagCIStatus(githubURLEndpoint string, projectName string, owner string, tag string) (string, error) {
	return getCIStatus(githubURLEndpoint, projectName, owner, "", tag, "")
}

// GetPRCIStatus gets the CI status using pull request number 
func GetPRCIStatus(githubURLEndpoint string, projectName string, owner string, prNumber string) (string, error) {
	return getCIStatus(githubURLEndpoint, projectName, owner, "", "", prNumber)
}

func getCIStatus(githubURLEndpoint string, projectName string, owner string, branch string, tag string, prNumber string) (string, error) {

	var URL string
	if tag != "" {
		URL = fmt.Sprintf("%s/repos/%s/%s/commits/%s/check-runs", githubURLEndpoint, owner, projectName, tag)
	}

	if branch != "" {
		URL = fmt.Sprintf("%s/repos/%s/%s/commits/%s/check-runs", githubURLEndpoint, owner, projectName, branch)
	}

	if prNumber != "" {
		sha, err := getSHAfromPR(owner, projectName, prNumber)
		if err != nil {
			return "", err
		}
		URL = fmt.Sprintf("%s/repos/%s/%s/commits/%s/check-runs", githubURLEndpoint, owner, projectName, sha)
	}

	// fmt.Println("request url is", url)

	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Accept", "application/vnd.github.antiope-preview+json")
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

	// Printing string data
	// fmt.Println(string(data))

	type CiStatus struct {
		ID         int    `json:"id"`
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
		App        struct {
			Slug string `json:"name"`
		}
	}

	type CheckRuns struct {
		CiStatuses []CiStatus `json:"check_runs"`
	}

	ciS := &CheckRuns{}
	err = json.Unmarshal(data, ciS)
	if err != nil {
		return "", err
	}

	var failCount int
	var successCount int

	if len(ciS.CiStatuses) == 0 {
		fmt.Println("No CI statuses for this request.")
	}
	

	// Display individual CI statuses on the terminal
	for i := 0; i < len(ciS.CiStatuses); i++ {
		conclusion := ciS.CiStatuses[i].Conclusion
		appName := ciS.CiStatuses[i].App.Slug
		fmt.Printf("%s: %s\n", appName, conclusion)
		if conclusion == "success" {
			successCount++
		} else {
			failCount++
		}
	}

	if failCount == 0 {
		return "true", nil
	}

	return "false", nil
}

func getSHAfromPR(owner string, projectName string, prNumber string) (string, error) {

	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%s", githubV3Url, owner, projectName, prNumber)
	// fmt.Println("Request url is", url)
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

	// fmt.Println(string(data))
	pullRequestSHA := struct {
		Head struct {
			SHA string `json:"sha"`
		}
	}{}

	err = json.Unmarshal(data, &pullRequestSHA)
	if err != nil {
		return "", err
	}
	// fmt.Println("SHA is............", pullRequestSHA.Head.SHA)
	return pullRequestSHA.Head.SHA, nil
}
