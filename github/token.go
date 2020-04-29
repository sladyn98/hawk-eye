package github

const (
	configKeyTokenValue = "value"
)

// Token holds an API access token data
type Token struct {
	Value string
}

// NewToken instantiate a new token
func NewToken(value string) *Token {
	return &Token{
		Value: value,
	}
}
