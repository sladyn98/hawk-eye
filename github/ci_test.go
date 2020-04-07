package github

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGetCIStatus(t *testing.T) {
	got, err := getCIStatus("configuration-as-code", "jenkinsci")
	if(err != nil) {
		fmt.Println(err)
	}
	assert.Equal(t, got, true)
}