package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"utdocs/core"
	"utdocs/diagnostics"
	"utdocs/manifest"
	"utdocs/utils"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

func DefaultNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Not Found"))
}

func FileServerWithCustom404(fs http.FileSystem, port int) http.Handler {
	color.Green("Serving on http://localhost:%d", port)
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if err == nil {
			fsh.ServeHTTP(w, r)
			return
		} else {
			DefaultNotFound(w, r)
			return
		}
	})
}

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
					if !utils.ShouldRebuild(event.Name, event.Op) {
						continue
					}
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

	err = http.ListenAndServe(":"+strconv.Itoa(port), FileServerWithCustom404(http.Dir(siteManifest.OutputPath), port))
	if err != nil {
		return err
	}

	return nil
}
