# Formatting Configuration

## .markdownlint.json

Updated markdown linting rules to be more practical for technical documentation:

- **MD013**: Increased line length to 100 characters (from 80)
- **MD022**: Disabled heading blank line requirements (too strict for lists)
- **MD029**: Already disabled (ordered list numbering)
- **MD032**: Disabled list blank line requirements (too strict for nested content)
- **MD033**: Allow specific HTML elements (`<details>`, `<summary>`, `<br>`)

## .prettierrc

Enhanced Prettier configuration with project-specific settings:

- **General**: 100 character line width, 2-space indentation
- **JavaScript/TypeScript**: Single quotes, trailing commas, arrow function parens
- **Markdown**: Always wrap prose, 100 character width
- **YAML**: Double quotes, 2-space indentation
- **JSON**: Double quotes, 2-space indentation

These settings ensure consistent formatting across all file types while being practical for technical documentation and code.
