package web

import (
	"bytes"
	"html/template"
	"testing"
)

// Compile the template
func TestTemplate(t *testing.T) {
	dashboardTemplate, err := template.New("dashboard").Parse(dashboardHTML)
	if err != nil {
		t.Error(err.Error())
	}
	err = dashboardTemplate.Execute(&bytes.Buffer{}, nil)
	if err != nil {
		t.Error(err.Error())
	}
}
