# kaleidoscopickitten

Query and modify YAML frontmatter in markdown files using yq expressions.

Uses yqlib syntax from https://mikefarah.gitbook.io/yq/v/v4.x/.

## Usage

```bash
# Print all frontmatter
kaleidoscopickitten frontmatter print . file.md
kaleidoscopickitten fm file.md
kaleidoscopickitten fm p file.md
kaleidoscopickitten fm p . file.md

# Print specific fields
kaleidoscopickitten fm .title file.md
kaleidoscopickitten fm p .title file.md
kaleidoscopickitten fm p .metadata.category file.md

# Query arrays
kaleidoscopickitten fm print .tags file.md
kaleidoscopickitten fm p '.tags | length' file.md
kaleidoscopickitten fm '.tags[0]' file.md

# Modify values (in-place editing)
kaleidoscopickitten frontmatter edit '.draft = true' file.md
kaleidoscopickitten fm e '.draft = true' file.md
kaleidoscopickitten fm e '.published = "2025-10-02"' file.md
kaleidoscopickitten fm e '.tags += ["new-tag"]' file.md

# Create frontmatter in files that have none
kaleidoscopickitten fm e '.title = "New Title"' file-without-frontmatter.md

# Clear all frontmatter
kaleidoscopickitten fm e '{}' file.md
```
