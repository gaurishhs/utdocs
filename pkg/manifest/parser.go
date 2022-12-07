package manifest

import (
	"errors"
	"path/filepath"

	"gopkg.in/ini.v1"
)

func ParseSiteManifest(path string) (SiteManifest, error) {
	manifest, err := ini.Load(path)
	if err != nil {
		return SiteManifest{}, err
	}

	result := SiteManifest{}

	siteSection := manifest.Section("Site")
	if siteSection != nil {
		result.Name = siteSection.Key("Name").String()
		result.ThemeId = siteSection.Key("Theme").MustString("default")
		result.HeadTags = append(result.HeadTags, siteSection.Key("HeadTags").Strings(",")...)
		result.DefaultSearch = siteSection.Key("DefaultSearch").MustBool(true)
		result.CustomFont = siteSection.Key("CustomFont").String()
		result.InputPath = siteSection.Key("Input").MustString("docs")
		result.OutputPath = siteSection.Key("Output").MustString("docs_gen")
		result.Logo = siteSection.Key("Logo").MustString("img/book.svg")
	}

	if !result.IsValid() {
		return result, errors.New("missing required parameters")
	}

	result.InputPath = filepath.Clean(result.InputPath)
	result.OutputPath = filepath.Clean(result.OutputPath)

	return result, nil
}

func ParseThemeManifest(path string) (ThemeManifest, error) {
	manifest, err := ini.Load(path)
	if err != nil {
		return ThemeManifest{}, err
	}

	result := ThemeManifest{}

	rootSection := manifest.Section("Theme")
	if rootSection != nil {
		result.Name = rootSection.Key("Name").String()
		result.Description = rootSection.Key("Description").String()
		result.Repository = rootSection.Key("Repository").String()
		result.Version = rootSection.Key("Version").String()
		result.Author = rootSection.Key("Author").String()
		result.License = rootSection.Key("License").String()
	}

	highlightingSection := manifest.Section("Highlighting")
	if highlightingSection != nil {
		result.Highlighting.Style = highlightingSection.Key("Style").MustString("bw")
		result.Highlighting.LineNumbers = highlightingSection.Key("LineNumbers").MustBool(false)
	}

	if !result.IsValid() {
		return result, errors.New("missing required parameters")
	}

	return result, nil
}
