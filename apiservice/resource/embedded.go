// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package resource

import (
	"embed"
)

//go:embed repome/*
var RepoMeResourceFS embed.FS

// GetRepoMeHTMLTemplate returns the embedded HTML template content
func GetRepoMeHTMLTemplate() ([]byte, error) {
	return RepoMeResourceFS.ReadFile("repome/template.html")
}

// GetRepoMeCSSContent returns the embedded CSS content
func GetRepoMeCSSContent() ([]byte, error) {
	return RepoMeResourceFS.ReadFile("repome/styles.css")
}

// GetRepoMeJSContent returns the embedded JS content
func GetRepoMeJSContent() ([]byte, error) {
	return RepoMeResourceFS.ReadFile("repome/scripts.js")
}

// GetRepoMeMarkdownContent returns the embedded markdown content
func GetRepoMeMarkdownContent() ([]byte, error) {
	return RepoMeResourceFS.ReadFile("repome/streams_to_river_repome.md")
}

// GetRepoMeStaticFile returns the content of a static file by filename
func GetRepoMeStaticFile(filename string) ([]byte, error) {
	return RepoMeResourceFS.ReadFile("repome/" + filename)
}
