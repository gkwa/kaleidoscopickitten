# kaleidoscopickitten

Query and modify YAML frontmatter in markdown files using yq expressions.

Uses yqlib syntax from https://mikefarah.gitbook.io/yq/v/v4.x/.

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
