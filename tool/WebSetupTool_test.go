package tool

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"html/template"
	"fmt"
)

func TestSearchInstitutions(t *testing.T) {

	temp, err := template.ParseFiles("../index.ejs")

	fmt.Print(err)

	assert.Nil(t, err)
	assert.NotNil(t, temp)
}
