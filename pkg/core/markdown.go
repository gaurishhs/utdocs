package core

import (
	"bytes"
	"os"
	"path/filepath"
	"utdocs/manifest"
	"utdocs/utils"

	hhtml "github.com/alecthomas/chroma/v2/formatters/html"
	headingid "github.com/jkboxomine/goldmark-headingid"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

func scanTree(node ast.Node, consumer func(node ast.Node)) {
	consumer(node)
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		scanTree(child, consumer)
	}
}

func analyzeDocument(astRoot ast.Node, source []byte, pageInfo *pageInfo) {
	ids := headingid.NewIDs()
	scanTree(astRoot, func(node ast.Node) {
		switch node.Kind() {
		case ast.KindHeading:
			heading := node.(*ast.Heading)

			// Find the page title
			if heading.Level == 1 && pageInfo.Title == "" {
				pageInfo.Title = string(heading.Text(source))
			}

			// Build the table of contents
			if heading.Level < 3 {
				headingName := heading.Text(source)
				pageInfo.Toc = append(pageInfo.Toc, tocEntry{
					Id:    string(ids.Generate(headingName, ast.KindHeading)),
					Name:  string(headingName),
					Level: heading.Level,
				})
			}
		}
	})
}

func renderMarkdownPage(mdFile string, theme manifest.ThemeManifest, siteManifest manifest.SiteManifest) (pageInfo, error) {
	result := pageInfo{
		FilePath: filepath.Clean(mdFile),
		FileName: utils.GetFileName(mdFile),
	}

	source, err := os.ReadFile(mdFile)
	if err != nil {
		return result, err
	}

	// Create Markdown parser
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle(theme.Highlighting.Style),
				highlighting.WithFormatOptions(
					hhtml.WithLineNumbers(theme.Highlighting.LineNumbers),
				),
			),
			emoji.New(
				emoji.WithRenderingMethod(emoji.Twemoji),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	// Parse Markdown
	reader := text.NewReader(source)
	context := parser.NewContext(parser.WithIDs(headingid.NewIDs()))
	astRoot := md.Parser().Parse(reader, parser.WithContext(context))
	analyzeDocument(astRoot, source, &result)
	if result.Title == "" {
		result.Title = utils.PrettifyTitle(mdFile)
	}

	// Render to HTML
	var buf bytes.Buffer
	err = md.Renderer().Render(&buf, source, astRoot)
	if err == nil {
		result.Body = buf.String()
	}

	return result, err
}
