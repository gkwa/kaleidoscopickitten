# kaleidoscopickitten

Query and modify YAML frontmatter in markdown files using yq expressions.

Uses yqlib syntax from https://mikefarah.gitbook.io/yq/v/v4.x/.

## Usage

```bash
# View all frontmatter
kaleidoscopickitten fm file.md
kaleidoscopickitten fm v file.md
kaleidoscopickitten fm v '.' file.md
kaleidoscopickitten frontmatter view '.' file.md

# View specific fields
kaleidoscopickitten fm '.title' file.md
kaleidoscopickitten fm v '.title' file.md
kaleidoscopickitten fm v '.metadata.category' file.md

# Query arrays
kaleidoscopickitten fm '.tags' file.md
kaleidoscopickitten fm '.tags | length' file.md
kaleidoscopickitten fm '.tags[0]' file.md

# Modify values (in-place editing)
kaleidoscopickitten fm e '.draft = true' file.md
kaleidoscopickitten fm e '.published = "2025-10-02"' file.md
kaleidoscopickitten fm e '.tags += ["new-tag"]' file.md
kaleidoscopickitten frontmatter edit '.draft = true' file.md

# Create frontmatter in files that have none
kaleidoscopickitten fm e '.title = "New Title"' file-without-frontmatter.md

# Clear all frontmatter
kaleidoscopickitten fm e '{}' file.md
```

## Commands

- `frontmatter` (alias: `fm`) `[expression] <file.md>` - Defaults to view mode when no subcommand specified
  - `view` (alias: `v`) `[expression] <file.md>` - Query frontmatter without modifying the file. If no expression is provided, defaults to `.` (show all).
  - `edit` (alias: `e`) `<expression> <file.md>` - Modify frontmatter and write changes back to the file
