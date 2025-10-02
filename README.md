# kaleidoscopickitten

Query and modify YAML frontmatter in markdown files using yq expressions.

## Usage

```bash
# Read fields (outputs frontmatter only)
kaleidoscopickitten '.title' file.md
kaleidoscopickitten '.metadata.category' file.md

# Query arrays
kaleidoscopickitten '.tags' file.md
kaleidoscopickitten '.tags | length' file.md
kaleidoscopickitten '.tags[0]' file.md

# Modify values (outputs processed frontmatter only)
kaleidoscopickitten '.draft = true' file.md
kaleidoscopickitten '.published = "2025-10-02"' file.md
kaleidoscopickitten '.tags += ["new-tag"]' file.md

# Create frontmatter in files that have none
kaleidoscopickitten -w '.title = "New Title"' file-without-frontmatter.md

# Clear all frontmatter
kaleidoscopickitten -w '{}' file.md

# Edit file in-place (writes complete file back)
kaleidoscopickitten -w '.draft = true' file.md
```

## Flags

- `-w, --write` - Write changes back to file (in-place editing)

Default behavior outputs only the processed frontmatter to stdout. Use `-w` to modify the file directly.

## Features

- Works with files that have no frontmatter - creates it automatically
- Clear all frontmatter with `{}`
- All yq expressions are supported
- Preserves markdown body content during in-place edits

