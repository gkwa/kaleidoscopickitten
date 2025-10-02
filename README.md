# kaleidoscopickitten

Query and modify YAML frontmatter in markdown files using yq expressions.

Uses yqlib syntax from https://mikefarah.gitbook.io/yq/v/v4.x/.

## Usage

```bash
# Read fields (outputs frontmatter only)
kaleidoscopickitten frontmatter view '.title' file.md
kaleidoscopickitten fm v '.metadata.category' file.md

# View all frontmatter (expression defaults to '.')
kaleidoscopickitten fm v file.md

# Query arrays
kaleidoscopickitten fm v '.tags' file.md
kaleidoscopickitten fm v '.tags | length' file.md
kaleidoscopickitten fm v '.tags[0]' file.md

# Modify values (in-place editing)
kaleidoscopickitten frontmatter edit '.draft = true' file.md
kaleidoscopickitten fm e '.published = "2025-10-02"' file.md
kaleidoscopickitten fm e '.tags += ["new-tag"]' file.md

# Create frontmatter in files that have none
kaleidoscopickitten fm e '.title = "New Title"' file-without-frontmatter.md

# Clear all frontmatter
kaleidoscopickitten fm e '{}' file.md
```

## Commands

- `frontmatter` (alias: `fm`) - Work with YAML frontmatter
  - `view` (alias: `v`) `[expression] <file.md>` - Query frontmatter without modifying the file. If no expression is provided, defaults to `.` (show all).
  - `edit` (alias: `e`) `<expression> <file.md>` - Modify frontmatter and write changes back to the file
