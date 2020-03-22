// Package github contains the Github bridge implementation
package github

import (
	"time"
)

const (
	target             = "github"
	metaKeyGithubLogin = "github-login"
	githubV3Url        = "https://api.github.com"
	defaultTimeout     = 60 * time.Second
)

type Github struct{}

func (*Github) Target() string {
	return target
}

func (g *Github) LoginMetaKey() string {
	return metaKeyGithubLogin
}
