package github

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGetCIStatus(t *testing.T) {
	got, err := getCIStatus("hawk-eye", "sladyn98","","450f8716451a212ef62b6ce3c66b10d4129271a3","")
	if(err != nil) {
		fmt.Println(err)
	}
	assert.Equal(t, got, true)
}