# Introduction

UTDocs is a powerful clean documentation generator written in Golang.

## Features

- Write documentation using Markdown
- Outputs static HTML that can be hosted anywhere (Github Pages, S3, etc)
- Supports syntax highlighting and emojis
- Easily extensible with custom themes using [Go templates](https://pkg.go.dev/text/template)
- In-built Search functionality (powered by [Fuse.js](https://fusejs.io/))



## Installation
- Head over to the [releases](https://github.com/gaurishhs/utdocs/releases) page and download the latest binary for your platform.

## Usage

### Create a new project

```bash
git clone -b starter https://github.com/gaurishhs/utdocs.git
```

### Start the development server

```bash
utdocs serve
```

### Build the static site

```bash
utdocs generate
```

