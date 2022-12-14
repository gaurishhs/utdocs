package core

import (
	"text/template"
	"utdocs/utils"
)

func generateTemplate(themeDir string) (*template.Template, error) {
	files, err := utils.ScanDir(themeDir, ".html")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
