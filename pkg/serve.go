package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"utdocs/core"
	"utdocs/diagnostics"
	"utdocs/manifest"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

func runServer(port int) error {
	// Parse manifest
	siteManifest, err := manifest.ParseSiteManifest(core.SiteManifestName)
	if err != nil {
		return err
	}

	// Start watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Start watch-event handler
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					// File system has changed, generate new version
					color.Yellow("File system has changed, generating new version...")
					// Delete Previous Search Index
					if siteManifest.DefaultSearch {
						erro := os.Remove(filepath.Join(siteManifest.OutputPath, "search", "index.json"))
						if erro != nil {
							log.Printf("Error deleting search index: %s\n", err)
						}
					}
					err := runGenerator()
					if err != nil {
						diagnostics.PrintError(err, "failed to regenerate")
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Error: %s\n", err)
			}
		}
	}()

	err = watcher.Add(siteManifest.InputPath)
	if err != nil {
		return err
	}

	// Generate initial version
	err = runGenerator()
	if err != nil {
		return err
	}

	// Start server
	server := http.FileServer(http.Dir(siteManifest.OutputPath))
	color.Green("Serving documentation on http://localhost:%d\n", port)
	err = http.ListenAndServe(":"+strconv.Itoa(port), server)
	if err != nil {
		return err
	}

	return nil
}
