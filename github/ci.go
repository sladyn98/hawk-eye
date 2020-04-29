package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getBranchCIStatus(githubURLEndpoint string, projectName string, owner string, branch string) (string, error) {
	return getCIStatus(githubURLEndpoint, projectName, owner, branch, "", "")
}

func getTagCIStatus(githubURLEndpoint string, projectName string, owner string, tag string) (string, error) {
	return getCIStatus(githubURLEndpoint, projectName, owner, "", tag, "")
}

func getCIStatus(githubURLEndpoint string, projectName string, owner string, branch string, tag string, prNumber string) (string, error) {

	var url string
	if tag != "" {
		url = fmt.Sprintf("%s/repos/%s/%s/commits/%s/check-runs", githubURLEndpoint, owner, projectName, tag)
	}

	if branch != "" {
		url = fmt.Sprintf("%s/repos/%s/%s/commits/%s/check-runs", githubURLEndpoint, owner, projectName, branch)
	}

	if prNumber != "" {
		sha, err := getSHAfromPR(owner, projectName, prNumber)
		if err != nil {
			return "", err
		}
		url = fmt.Sprintf("%s/repos/%s/%s/commits/%s/check-runs", githubURLEndpoint, owner, projectName, sha)
	}

	fmt.Println("request url is", url)

	req, err := http.NewRequest("GET", url, nil)
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

	// Display individual CI statuses on the terminal
	for i := 0; i < len(ciS.CiStatuses); i++ {
		conclusion := ciS.CiStatuses[i].Conclusion
		appName := ciS.CiStatuses[i].App.Slug
		fmt.Printf("%s: %s\n", appName, conclusion)
		if conclusion == "action_required" || conclusion == "canceled" || conclusion == "timed_out" ||
			conclusion == "failed" {
			failCount++
		} else {
			successCount++
		}
	}

	if failCount == 0 {
		return "true", nil
	}

	return "false", nil
}

func getSHAfromPR(owner string, projectName string, prNumber string) (string, error) {

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
