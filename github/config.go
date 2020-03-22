package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sladyn98/hawk-eye/input"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	ErrBadProjectURL = errors.New("bad project url")
)

func PromptLogin() (string, error) {
	var login string

	validator := func(_ string, value string) (string, error) {
		ok, fixed, err := validateUsername(value)
		if err != nil {
			return "", err
		}
		if !ok {
			return "invalid login", nil
		}
		login = fixed
		return "", nil
	}

	_, err := input.Prompt("Github login", "login", input.Required, validator)
	if err != nil {
		return "", err
	}

	return login, nil
}

func validateUsername(username string) (bool, string, error) {
	url := fmt.Sprintf("%s/users/%s", githubV3Url, username)

	client := &http.Client{
		Timeout: defaultTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, "", err
	}

	if resp.StatusCode != http.StatusOK {
		return false, "", nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	err = resp.Body.Close()
	if err != nil {
		return false, "", err
	}

	var decoded struct {
		Login string `json:"login"`
	}
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		return false, "", err
	}

	if decoded.Login == "" {
		return false, "", fmt.Errorf("validateUsername: missing login in the response")
	}

	return true, decoded.Login, nil
}

func LoginAndRequestToken(login, owner, project string) (string, error) {
	fmt.Println("hawk-eye will now generate an access token in your Github profile. Your credential are not stored and are only used to generate the token. The token is stored in the global git config.")
	fmt.Println()
	fmt.Println("The access scope depend on the type of repository.")
	fmt.Println("Public:")
	fmt.Println("  - 'public_repo': to be able to read public repositories")
	fmt.Println("Private:")
	fmt.Println("  - 'repo'       : to be able to read private repositories")
	fmt.Println()

	// prompt project visibility to know the token scope needed for the repository
	index, err := input.PromptChoice("repository visibility", []string{"public", "private"})
	if err != nil {
		return "", err
	}
	scope := []string{"public_repo", "repo"}[index]

	password, err := input.PromptPassword("Password", "password", input.Required)
	if err != nil {
		return "", err
	}

	// Attempt to authenticate and create a token
	note := fmt.Sprintf("hawk-eye - %s/%s", owner, project)
	resp, err := requestToken(note, login, password, scope)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// Handle 2FA is needed
	OTPHeader := resp.Header.Get("X-GitHub-OTP")
	if resp.StatusCode == http.StatusUnauthorized && OTPHeader != "" {
		otpCode, err := input.PromptPassword("Two-factor authentication code", "code", input.Required)
		if err != nil {
			return "", err
		}

		resp, err = requestTokenWith2FA(note, login, password, otpCode, scope)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusCreated {
		return decodeBody(resp.Body)
	}

	b, _ := ioutil.ReadAll(resp.Body)
	return "", fmt.Errorf("error creating token %v: %v", resp.StatusCode, string(b))
}

func randomFingerprint() string {
	// Doesn't have to be crypto secure, it's just to avoid token collision
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 32)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func requestToken(note, login, password string, scope string) (*http.Response, error) {
	return requestTokenWith2FA(note, login, password, "", scope)
}

func requestTokenWith2FA(note, login, password, otpCode string, scope string) (*http.Response, error) {
	url := fmt.Sprintf("%s/authorizations", githubV3Url)
	params := struct {
		Scopes      []string `json:"scopes"`
		Note        string   `json:"note"`
		Fingerprint string   `json:"fingerprint"`
	}{
		Scopes:      []string{scope},
		Note:        note,
		Fingerprint: randomFingerprint(),
	}

	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(login, password)
	req.Header.Set("Content-Type", "application/json")

	if otpCode != "" {
		req.Header.Set("X-GitHub-OTP", otpCode)
	}

	client := &http.Client{
		Timeout: defaultTimeout,
	}

	return client.Do(req)
}

func decodeBody(body io.ReadCloser) (string, error) {
	data, _ := ioutil.ReadAll(body)

	aux := struct {
		Token string `json:"token"`
	}{}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return "", err
	}

	if aux.Token == "" {
		return "", fmt.Errorf("no token found in response: %s", string(data))
	}

	return aux.Token, nil
}

func PromptURL() (string, string, error) {
	validRemotes := make([]string, 0)

	validator := func(name, value string) (string, error) {
		_, _, err := splitURL(value)
		if err != nil {
			return err.Error(), nil
		}
		return "", nil
	}

	url, err := input.PromptURLWithRemote("Github project URL", "URL", validRemotes, input.Required, validator)
	if err != nil {
		return "", "", err
	}

	return splitURL(url)
}

// splitURL extract the owner and project from a github repository URL. It will remove the
// '.git' extension from the URL before parsing it.
// Note that Github removes the '.git' extension from projects names at their creation
func splitURL(url string) (owner string, project string, err error) {
	cleanURL := strings.TrimSuffix(url, ".git")

	re := regexp.MustCompile(`github\.com[/:]([a-zA-Z0-9\-_]+)/([a-zA-Z0-9\-_.]+)`)

	res := re.FindStringSubmatch(cleanURL)
	if res == nil {
		return "", "", ErrBadProjectURL
	}

	owner = res[1]
	project = res[2]
	return
}
