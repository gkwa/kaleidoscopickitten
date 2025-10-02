# kaleidoscopickitten

Query and modify YAML frontmatter in markdown files using yq expressions.

Uses yqlib syntax from https://mikefarah.gitbook.io/yq/v/v4.x/.

## Usage

```bash
# Read fields (outputs frontmatter only)
kaleidoscopickitten frontmatter view '.title' file.md
kaleidoscopickitten frontmatter view '.metadata.category' file.md

# Query arrays
kaleidoscopickitten frontmatter view '.tags' file.md
kaleidoscopickitten frontmatter view '.tags | length' file.md
kaleidoscopickitten frontmatter view '.tags[0]' file.md

# Modify values (in-place editing)
kaleidoscopickitten frontmatter edit '.draft = true' file.md
kaleidoscopickitten frontmatter edit '.published = "2025-10-02"' file.md
kaleidoscopickitten frontmatter edit '.tags += ["new-tag"]' file.md

# Create frontmatter in files that have none
kaleidoscopickitten frontmatter edit '.title = "New Title"' file-without-frontmatter.md

# Clear all frontmatter
kaleidoscopickitten frontmatter edit '{}' file.md
```

## Commands

- `frontmatter view <expression> <file.md>` - Query frontmatter without modifying the file
- `frontmatter edit <expression> <file.md>` - Modify frontmatter and write changes back to the file
