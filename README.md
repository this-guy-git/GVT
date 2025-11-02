<div align="center">

![GVT logo (1)](https://github.com/user-attachments/assets/8aa72a9b-abff-4d87-a04d-cca80669358e)




# GVT (Guy's Versioning Tool)

[![License](https://img.shields.io/github/license/this-guy-git/GVT?style=flat&color=7c3aed&labelColor=222222)](https://www.gnu.org/licenses/gpl-3.0)
[![Stars](https://img.shields.io/github/stars/this-guy-git/GVT?style=flat&color=7c3aed&labelColor=222222)](https://github.com/this-guy-git/GVT/stargazers)
[![Forks](https://img.shields.io/github/forks/this-guy-git/GVT?style=flat&color=7c3aed&labelColor=222222)](https://github.com/this-guy-git/GVT/forks)
[![Issues](https://img.shields.io/github/issues/this-guy-git/GVT?style=flat&color=7c3aed&labelColor=222222)](https://github.com/this-guy-git/GVT/issues)
[![Pull Requests](https://img.shields.io/github/issues-pr/this-guy-git/GVT?style=flat&color=7c3aed&labelColor=222222)](https://github.com/this-guy-git/GVT/pulls)
[![Release](https://img.shields.io/github/v/release/this-guy-git/GVT?style=flat&color=7c3aed&labelColor=222222)](https://github.com/this-guy-git/GVT/releases/)

**A lightweight version control system designed to help you track changes in your projects, manage branches, and maintain a complete history of your work.**

[Features](#core-commands) • [Installation](#installation) • [Documentation](#table-of-contents) • [Contributing](#support)

</div>

---

## Table of Contents

1. [Installation](#installation)
2. [Getting Started](#getting-started)
3. [Core Commands](#core-commands)
4. [Branching & Merging](#branching--merging)
5. [Advanced Features](#advanced-features)
6. [Configuration](#configuration)
7. [File Ignoring](#file-ignoring)
8. [Best Practices](#best-practices)

---

## Installation

### Linux Installation

```bash
# Download the latest release
wget https://github.com/this-guy-git/GVT/releases/latest/download/gvt-linux.tar.gz

# Make folder for GVT
sudo mkdir -p /usr/local/bin/gvt

# Extract to /usr/local/bin
sudo tar -xzf gvt-linux.tar.gz -C /usr/local/bin/gvt

# Make executable
sudo chmod +x /usr/local/bin/gvt/bin/linux/gvt

# Add to PATH
export PATH="/usr/local/bin/gvt/bin/linux:$PATH"
```

To make the PATH change permanent, add the export line to your `~/.bashrc` or `~/.zshrc` file.

### Windows Installation

Download the installer from [here](https://github.com/this-guy-git/GVT/releases/latest/download/GVTinstaller.exe).

### Build from source
Requires:
- [Make](https://www.gnu.org/software/make/)
- [Go (1.24+)](https://go.dev/)

```bash
# Clone repo
git clone https://github.com/this-guy-git/GVT.git
cd GVT
# Install Dependancies
go install github.com/spf13/cobra@latest
go install github.com/akavel/rsrc@latest

# Build with make
make PLATFORM=linux/win
```

---

## Getting Started

### Initialize a Repository

Create a new GVT repository in your current directory:

```bash
gvt init
```

This creates a `.gvt` directory that stores all version control data.

### Configure Your Identity

Set your name and email for commit tracking:

```bash
gvt config user.name "Your Name"
gvt config user.email "you@example.com"
```

### Basic Workflow

1. **Stage files** you want to track:
   ```bash
   gvt add file.txt
   gvt add .              # Stage all files
   ```

2. **Commit** your changes:
   ```bash
   gvt commit -m "Your commit message"
   ```

3. **Check status** to see what's changed:
   ```bash
   gvt status
   ```

---

## Core Commands

### `gvt init`
Initialize a new GVT repository in the current directory.

```bash
gvt init
```

### `gvt add [files]`
Stage files or directories for the next commit.

```bash
gvt add file.txt           # Stage single file
gvt add src/              # Stage entire directory
gvt add .                 # Stage all files
```

Files matching patterns in `.gvtignore` will be automatically excluded.

### `gvt commit`
Record staged changes as a new commit.

```bash
gvt commit -m "Add new feature"     # With message
gvt commit                          # Auto-generated message
```

If no message is provided, GVT generates one based on the number of files changed.

### `gvt status`
Display the current state of your working directory.

```bash
gvt status
```

Shows:
- Current branch
- Staged files
- Modified files
- Untracked files

### `gvt log`
View the commit history.

```bash
gvt log
```

Displays:
- Commit ID
- Author
- Timestamp
- Commit message

### `gvt diff`
Show changes between the last two commits.

```bash
gvt diff              # Show changed files only
gvt diff --all        # Show all files
```

Displays a Git-like colored diff with:
- Green `+` for added lines
- Red `-` for removed lines
- Context lines around changes

### `gvt remove [files]`
Unstage files without deleting them.

```bash
gvt remove file.txt
gvt remove src/
gvt remove .              # Unstage all
```

---

## Branching & Merging

### `gvt branch [name]`
List or create branches.

```bash
gvt branch                # List all branches
gvt branch feature-x      # Create new branch
```

The current branch is marked with an asterisk (`*`).

### `gvt switch <branch>`
Switch to a different branch.

```bash
gvt switch feature-x          # Switch to existing branch
gvt switch -c new-feature     # Create and switch to new branch
```

**Flags:**
- `-c, --create`: Create the branch if it doesn't exist

### `gvt merge <branch>`
Merge another branch into the current branch.

```bash
gvt merge feature-x
```

**Note:** Only fast-forward merges are currently supported. If conflicts are detected, the merge will be aborted and you'll need to resolve them manually.

---

## Advanced Features

### `gvt revert [commit-id]`
Restore your working directory to a previous commit.

```bash
gvt revert                    # Revert to last commit
gvt revert 20250101-120000    # Revert to specific commit
```

This replaces all tracked files with their versions from the specified commit.

### `gvt restore <file>`
Restore a single file to its last committed state.

```bash
gvt restore file.txt
```

Useful for discarding unwanted changes to specific files.

### `gvt delete [commit-id]`
Remove a commit from history.

```bash
gvt delete 20250101-120000        # Delete specific commit
gvt delete 20250101-120000 -f     # Force delete latest commit
gvt delete --all                  # Delete all commits (with confirmation)
```

**Flags:**
- `-f, --force`: Allow deletion of the latest commit
- `--all`: Delete all commits (requires confirmation)

### `gvt clone <source> [directory]`
Clone an existing GVT repository.

```bash
gvt clone /path/to/repo              # Clone to repo name
gvt clone /path/to/repo my-copy      # Clone to specific directory
```

Respects `.gvtignore` patterns during cloning.

---

## Configuration

### Configuration Levels

GVT supports two configuration levels:

1. **Global** (default): `~/.gvt/config.json`
2. **Local** (repository-specific): `.gvt/config.json`

### View Configuration

```bash
gvt config user.name              # Get username
gvt config user.email             # Get email
```

### Set Configuration

```bash
# Global configuration (default)
gvt config user.name "Your Name"
gvt config user.email "you@example.com"

# Local configuration (repository-specific)
gvt config -l user.name "Project Name"
gvt config -l user.email "project@example.com"
```

**Flags:**
- `-l, --local`: Use local repository configuration

---

## File Ignoring

Create a `.gvtignore` file in your repository root to exclude files from version control.

### `.gvtignore` Syntax

```
# Comments start with #
*.log
*.tmp
node_modules/
.env
build/
dist/
```

**Pattern Rules:**
- Use `*` as a wildcard
- End directories with `/`
- One pattern per line
- Lines starting with `#` are comments

### Common Patterns

```
# Dependencies
node_modules/
vendor/

# Build outputs
build/
dist/
*.exe
*.o

# Environment files
.env
.env.local

# OS files
.DS_Store
Thumbs.db

# IDE files
.vscode/
.idea/
*.swp

# Logs
*.log
logs/
```

---

## Best Practices

### Commit Messages

Write clear, descriptive commit messages:

```bash
# Good
gvt commit -m "Add user authentication feature"
gvt commit -m "Fix login bug for mobile devices"
gvt commit -m "Update documentation for API endpoints"

# Avoid
gvt commit -m "updates"
gvt commit -m "fix"
gvt commit -m "wip"
```

### Staging Strategy

Stage related changes together:

```bash
# Stage specific features
gvt add src/auth/
gvt commit -m "Implement authentication"

gvt add src/api/
gvt commit -m "Add API endpoints"
```

### Branching Strategy

Use branches for features and experiments:

```bash
# Create a feature branch
gvt switch -c feature-new-ui

# Make changes and commit
gvt add .
gvt commit -m "Redesign user interface"

# Switch back and merge
gvt switch main
gvt merge feature-new-ui
```

### Regular Commits

Commit frequently with logical changes:
- After completing a feature
- After fixing a bug
- Before making major changes
- At the end of each work session

### Before Major Changes

Always commit or revert before making significant changes:

```bash
# Check status
gvt status

# Commit current work
gvt add .
gvt commit -m "Save progress before refactoring"

# Or revert unwanted changes
gvt revert
```

---

## Command Reference

| Command | Description |
|---------|-------------|
| `gvt init` | Initialize a new repository |
| `gvt add [files]` | Stage files for commit |
| `gvt commit [-m "message"]` | Create a new commit |
| `gvt status` | Show working directory status |
| `gvt log` | View commit history |
| `gvt diff [--all]` | Show changes between commits |
| `gvt remove [files]` | Unstage files |
| `gvt branch [name]` | List or create branches |
| `gvt switch <branch>` | Switch to another branch |
| `gvt merge <branch>` | Merge a branch |
| `gvt revert [commit-id]` | Restore to a previous commit |
| `gvt restore <file>` | Restore a single file |
| `gvt delete [commit-id]` | Delete a commit |
| `gvt clone <source>` | Clone a repository |
| `gvt config <key> [value]` | Get or set configuration |
| `gvt version` | Show GVT version |

---

## Technical Details

### Storage Format

- Commits are stored in `.gvt/commits/<branch>/<commit-id>/`
- Files are compressed using zlib compression
- Each commit has a `meta.json` with commit metadata
- File hashes (MD5) detect changes

### Commit ID Format

Commits use timestamp-based IDs: `YYYYMMDD-HHMMSS`

Example: `20250101-143022`

### Repository Structure

```
.gvt/
├── HEAD                    # Points to current branch
├── config.json            # Local configuration
├── history.json           # Commit history
├── stage.json            # Staged files
├── commits/              # Commit storage
│   └── main/            # Branch commits
│       └── 20250101-120000/
│           ├── meta.json
│           └── file.txt.zlib
└── refs/
    └── heads/            # Branch references
        └── main
```

---

## Troubleshooting

### "Not a GVT repository"

**Problem:** Running GVT commands outside a repository.

**Solution:** Run `gvt init` or navigate to a directory with `.gvt`.

### "Nothing staged to commit"

**Problem:** No files are staged.

**Solution:** Use `gvt add` to stage files before committing.

### "Branch already exists"

**Problem:** Trying to create a branch that already exists.

**Solution:** Use `gvt switch <branch>` to switch to it, or choose a different name.

### Files Not Being Staged

**Problem:** Files appear untracked or won't stage.

**Solution:** Check if they match patterns in `.gvtignore`.

---

## License

GVT is free software licensed under the GNU General Public License v3.0.

You are free to:
- Use GVT for any purpose
- Study and modify the source code
- Distribute copies
- Distribute modified versions

For more details, see <https://www.gnu.org/licenses/>.

---

## Support

For bug reports, feature requests, or contributions, visit:
- GitHub: `github.com/this-guy-git/GVT`
- Email: `thisguy@thisguylabs.com`

---

**Copyright © 2025 this guy Labs**
