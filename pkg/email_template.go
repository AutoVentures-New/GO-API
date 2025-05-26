package pkg

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
)

type EmailValidateCode struct {
	Code string
}

type EmailClickLink struct {
	Link string
}

func ParseTemplate(
	ctx context.Context,
	templateName string,
	data any,
) (string, error) {
	content, err := os.ReadFile(fmt.Sprintf("./email-templates/%s.html", templateName))
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("email").Parse(string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
