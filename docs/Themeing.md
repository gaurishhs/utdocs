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

- TODO

