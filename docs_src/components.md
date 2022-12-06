# Components
- List of components in the default theme

## Alert
- Alerts are used to display important information to the user.

### Usage
- The alert component can be used as follows:

```html
<ut-alert type="info" message="This is an info alert"></ut-alert>
```

### Props
- `type` - The type of the alert. Can be one of `info`, `success`, `warn`, `error`.
- `message` - The message to display in the alert.

### Example

<ut-alert type="info" message="This is an info alert"></ut-alert>
<ut-alert type="success" message="This is a success alert"></ut-alert>
<ut-alert type="warn" message="This is a warning alert"></ut-alert>
<ut-alert type="error" message="This is an error alert"></ut-alert>