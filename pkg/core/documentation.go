package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
	"utdocs/diagnostics"
	"utdocs/manifest"
	"utdocs/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/microcosm-cc/bluemonday"
)

type navNode struct {
	Name     string
	Url      string
	Children []*navNode
	Active   bool
}

type tocEntry struct {
	Id    string
	Name  string
	Level int
}

type generatorInfo struct {
	Name    string
	Version string
}

type pageInfo struct {
	FilePath string
	FileName string
	Title    string
	Body     string
	Toc      []tocEntry
}

type SearchIndexEntry struct {
	Title   string
	Url     string
	Content string
}

type pageContext struct {
	Page      pageInfo
	Generator generatorInfo
	Now       string
	Site      manifest.SiteManifest
	Nav       []*navNode
	RootPath  string
	Url       string
}

func isIndexFile(filePath string) bool {
	return utils.GetFileName(filePath) == IndexFileName
}

func findDirForPage(page pageInfo, siteManifest manifest.SiteManifest) string {
	// Build the path for the output directory
	relativePath := page.FilePath[len(siteManifest.InputPath)+1:]
	outputDirPath := filepath.Dir(relativePath)

	// If it is not the index file, put it in its own subdirectory so that we get /subdir/index.html
	// which we can open in the browser as just /subdir
	if !isIndexFile(page.FilePath) {
		outputDirPath = filepath.Join(outputDirPath, page.FileName)
	}

	return outputDirPath
}

func createPageContext(mdFile string, rootPath string, siteManifest manifest.SiteManifest, themeManifest manifest.ThemeManifest) (pageContext, error) {
	page, err := renderMarkdownPage(mdFile, themeManifest, siteManifest)
	if err != nil {
		return pageContext{}, err
	}
	return pageContext{
		Page: page,
		Site: siteManifest,
		Generator: generatorInfo{
			Name:    ProgramName,
			Version: ProgramVersion,
		},
		Now:      time.Now().Format("2006-01-02 15:04:05.000"),
		RootPath: rootPath,
		Url:      filepath.ToSlash(findDirForPage(page, siteManifest)),
	}, nil
}

func openOutputFileForPage(pageContext *pageContext, siteManifest manifest.SiteManifest) (*os.File, error) {
	outPath := filepath.Join(siteManifest.OutputPath, pageContext.Url)

	err := os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(filepath.Join(outPath, "index.html"))
	return file, err
}

func generateThemedHtmlForPage(pageContext *pageContext, siteManifest manifest.SiteManifest, themeTemplate *template.Template) {
	mdFile := pageContext.Page.FileName
	writer, err := openOutputFileForPage(pageContext, siteManifest)
	if err != nil {
		diagnostics.PrintError(err, "failed to open output file for "+mdFile)
		return
	}

	htmlBuf := strings.Builder{}
	err = themeTemplate.ExecuteTemplate(&htmlBuf, RootTemplateName, pageContext)
	if err != nil {
		diagnostics.PrintError(err, "failed to execute template for "+mdFile)
		return
	}

	err = processHtml(strings.NewReader(htmlBuf.String()), writer, pageContext)
	if err != nil {
		diagnostics.PrintError(err, "failed to run HTML postproc for "+mdFile)
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBuf.String()))
	if err != nil {
		diagnostics.PrintError(err, "failed to parse HTML for "+mdFile)
		return
	}

	if siteManifest.DefaultSearch {
		html, e := doc.Find(".main-content").Html()
		if e != nil {
			diagnostics.PrintError(e, "failed to find main content for "+mdFile)
			return
		}

		stripper := bluemonday.StrictPolicy()

		// Remove all the HTML tags
		html = stripper.Sanitize(html)

		html = strings.ReplaceAll(html, "\n", " ")

		AddToSearchIndex(siteManifest, SearchIndexEntry{
			Title:   pageContext.Page.Title,
			Url:     pageContext.Url,
			Content: html,
		})
	}

	err = writer.Close()
	if err != nil {
		diagnostics.PrintError(err, "failed to close output file for "+mdFile)
		return
	}
}

func sortTree(node *navNode) {
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].Url == "." {
			return true
		}
		return len(node.Children[i].Children) < len(node.Children[j].Children)
	})
	for _, c := range node.Children {
		sortTree(c)
	}
}

func prepareDocumentationTree(dirPath string, rootDirPrefix string, parentNode *navNode, siteManifest manifest.SiteManifest, themeManifest manifest.ThemeManifest, contexts *[]pageContext) error {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, dirent := range dir {
		childPath := filepath.Join(dirPath, dirent.Name())
		newNode := &navNode{
			Name: strings.ReplaceAll(dirent.Name(), "-", " "),
			Url:  "",
		}

		if dirent.IsDir() {
			err := prepareDocumentationTree(childPath, rootDirPrefix+"../", newNode, siteManifest, themeManifest, contexts)
			if err != nil {
				return err
			}
		} else if filepath.Ext(dirent.Name()) == ".md" {
			diagnostics.Debug(func() {
				log.Println("processing:\n", childPath)
			})

			currentRootDirPrefix := rootDirPrefix
			if !isIndexFile(childPath) {
				currentRootDirPrefix = "../" + currentRootDirPrefix
			}
			context, err := createPageContext(childPath, currentRootDirPrefix, siteManifest, themeManifest)
			if err != nil {
				return err
			}

			newNode.Name = context.Page.Title
			newNode.Url = context.Url
			*contexts = append(*contexts, context)
		}

		parentNode.Children = append(parentNode.Children, newNode)
	}

	sortTree(parentNode)

	return nil
}

func copyMediaFiles(siteManifest manifest.SiteManifest, themeDir string) error {
	predicate := func(ext string) bool {
		return ext != ".md" && ext != ".html" && ext != ".ini"
	}

	err := utils.CopyDirContents(siteManifest.InputPath, siteManifest.OutputPath, predicate)
	if err != nil {
		return err
	}

	err = utils.CopyDirContents(themeDir, siteManifest.OutputPath, predicate)
	if err != nil {
		return err
	}

	return nil
}

func AddToSearchIndex(siteManifest manifest.SiteManifest, ndata SearchIndexEntry) {
	// If search index.json does not exist, create it
	indexPath := filepath.Join(siteManifest.OutputPath, "search", "index.json")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		// Create the search directory if it does not exist
		searchPath := filepath.Join(siteManifest.OutputPath, "search")
		if _, err := os.Stat(searchPath); os.IsNotExist(err) {
			err := os.MkdirAll(searchPath, os.ModePerm)
			if err != nil {
				diagnostics.PrintError(err, "failed to create search directory")
				return
			}
		}

		// Create the index.json file
		file, err := os.Create(indexPath)
		if err != nil {
			diagnostics.PrintError(err, "failed to create search index file")
			return
		}
		defer file.Close()
		data, err := json.Marshal(ndata)
		if err != nil {
			diagnostics.PrintError(err, "failed to marshal search content")
			return
		}
		_, err = file.Write([]byte(`[ ` + string(data) + ` ]`))
		if err != nil {
			diagnostics.PrintError(err, "failed to write search content")
			return
		}
	} else {
		// Read the index.json file
		file, err := os.Open(indexPath)
		if err != nil {
			diagnostics.PrintError(err, "failed to open search index file")
			return
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			diagnostics.PrintError(err, "failed to read search index file")
			return
		}

		var existingData []SearchIndexEntry
		if err := json.Unmarshal(data, &existingData); err != nil {
			diagnostics.PrintError(err, "failed to unmarshal search index file")
			return
		}

		// Append the new data
		existingData = append(existingData, ndata)
		jsonData, err := json.Marshal(existingData)
		if err != nil {
			diagnostics.PrintError(err, "failed to marshal search content")
			return
		}
		// Write the new data
		ioutil.WriteFile(indexPath, jsonData, os.ModePerm)
	}
}

func GenerateDocumentation(siteManifest manifest.SiteManifest, themeManifest manifest.ThemeManifest, themeDir string) error {
	var stopwatch diagnostics.Stopwatch
	stopwatch.Reset()

	var navTreeRoot navNode
	var generatedPageContexts []pageContext
	err := prepareDocumentationTree(siteManifest.InputPath, "", &navTreeRoot, siteManifest, themeManifest, &generatedPageContexts)
	if err != nil {
		return err
	}

	themeTemplate, err := generateTemplate(themeDir)
	if err != nil {
		return err
	}

	for _, pageCtx := range generatedPageContexts {
		pageCtx.Nav = navTreeRoot.Children
		generateThemedHtmlForPage(&pageCtx, siteManifest, themeTemplate)
	}

	err = copyMediaFiles(siteManifest, themeDir)
	if err != nil {
		return err
	}

	color.Green("Documentation generated in %dms", stopwatch.Milliseconds())

	return nil
}
