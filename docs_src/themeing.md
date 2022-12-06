# Themeing

Themes are built using [Go's templating sytem](https://pkg.go.dev/text/template). You can check out the source code for the default theme [here](https://github.com/gaurishhs/utdocs/tree/main/themes/default), which should give you a start into making your own theme.

To allow UTDocs to recognize, you need a manifest at the theme directory root called `theme.ini`:

```ini
[Theme]
Name = My UTDocs Theme
Description = A new theme for UTDocs
Repository = https://github.com/gaurishhs/my-theme
Version = 1.0.0
Author = You
License = MIT

[Highlighting]
Style = bw
LineNumbers = false
```

The only properties that are mandatory is the `Name` and `Version` key. The highlighter will assume the default settings shown above, and all other information will default to empty.

Cutedoc uses Goldmark for Markdown processing, which internally uses `chroma` as a code highlighter. An overview of available highlighting styles can be found [here](https://xyproto.github.io/splash/docs/all.html).

## Customizing the theme

### Directory Structure

The theme directory should have the following structure:

```bash
.
├── themes
│   └── my-theme (theme id)
│       ├── js
│       │   └── main.js
│       ├── css
│       │   └── main.css
│       ├── index.html
│       ├── 404.html
│       ├── nav.html
│       ├── main.html
│       ├── body.html
│       ├── theme.ini (Required)
```

Take a look at the [default theme](https://github.com/gaurishhs/utdocs/tree/main/themes/default) for reference.

### Templates

The `index.html` and `404.html` files are the templates for the home page and the 404 page respectively. The `main.html` file is the template for the main page, which is the page that contains the documentation. The `body.html` file is the template for the body of the page, which is the part that contains the documentation.

The `main.html` and `body.html` files are the only files that are required. The other files are optional.

