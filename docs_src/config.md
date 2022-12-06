# Configuration
- UTDocs uses [`INI`](https://en.wikipedia.org/wiki/INI_file) file format for configuration.
- The configuration file is located at `utdocs.ini` in the root of the project.

## Options

### Site Name
- The name of site. Used in the title of the HTML pages to generate titles such as `Home - My Site`.

```ini
[Site]
Name = My Site
```

### Disabling Default Search
- By default, UTDocs comes with a built-in full text search functionality. If you want to disable it, you can do so by setting the `DefaultSearch` option to `false`.

```ini
[Site]
DefaultSearch = false
```

### Theme ID
- The ID of the theme to use. The default theme is `default`. You can find more themes in the [Themeing](/themeing) section.

```ini
[Site]
Theme = default
```

### Theme Directory
- The directory where the theme is located. By default, UTDocs will look for the theme in the `themes` directory in the root of the project.

```bash
export THEME_DIR=/path/to/theme
```

### Custom Head Tags

- You can add custom tags to the `<head>` of the HTML pages by adding them to the `HeadTags` option.

```ini
[Site]
HeadTags = <link rel="stylesheet" href="https://example.com/style.css">
```

### Custom Font
- You can add a custom font to the site by adding the link to the `CustomFont` option.

<ut-alert type="warn" message="Make sure to have the font in the head by using custom head tags."></ut-alert>

```ini
[Site]
CustomFont = Roboto
```



