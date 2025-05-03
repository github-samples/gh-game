# Contributing to gh-game

Thank you for your interest in contributing to gh-game! We love your input and welcome contributions from everyone, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Implementing a new game
- Improving documentation

## Table of Contents

- [Getting Started](#getting-started)
  - [Fork the Repository](#fork-the-repository)
  - [Clone Your Fork](#clone-your-fork)
  - [Set Up GitHub Codespaces](#set-up-github-codespaces)
  - [Local Development Setup](#local-development-setup)
- [Development Workflow](#development-workflow)
  - [Creating a Branch](#creating-a-branch)
  - [Make Your Changes](#make-your-changes)
  - [Development Commands](#development-commands)
  - [Adding a New Game](#adding-a-new-game)
- [Pull Request Process](#pull-request-process)
  - [Create a Pull Request](#create-a-pull-request)
  - [Code Review Process](#code-review-process)
- [Contribution Standards](#contribution-standards)
  - [Code Style](#code-style)
  - [Documentation](#documentation)
  - [Testing](#testing)
- [Non-Code Contributions](#non-code-contributions)
- [Community](#community)
- [License](#license)

## Getting Started

### Fork the Repository

Start by forking the repository on GitHub by clicking the "Fork" button at the top-right of the repository page. This creates a copy of the repository in your own GitHub account.

### Clone Your Fork

Clone your fork to your local machine or continue with GitHub Codespaces:

```bash
git clone https://github.com/YOUR-USERNAME/gh-game.git
cd gh-game
git remote add upstream https://github.com/ORIGINAL-OWNER/gh-game.git
```

### Set Up GitHub Codespaces

The easiest way to start contributing is by using GitHub Codespaces, which comes pre-configured with all the dependencies needed for gh-game:

1. Navigate to your fork on GitHub
2. Click the "Code" button
3. Select "Open with Codespaces"
4. Click "New codespace"

This will create a development environment in the cloud with all the necessary tools installed, including the GitHub CLI (`gh`).

### Local Development Setup

If you prefer to develop locally, make sure you have:

1. Go 1.23 or later installed
2. GitHub CLI (`gh`) installed
3. All project dependencies:
   ```bash
   go mod download
   ```
4. Any other tools:
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

## Development Workflow

### Creating a Branch

Create a branch for your contributions:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bugfix-name
```

### Make Your Changes

Now you're ready to make your changes to the codebase.

### Development Commands

Here are the important commands to use during development:

- **Build the project**:
  ```bash
  go build
  ```

- **Run the CLI extension**:
  ```bash
  ./gh-game <subcommand>
  # For example:
  ./gh-game cointoss
  ./gh-game tictactoe
  ```

- **Run tests**:
  ```bash
  go test ./...
  ```

- **Format code** (required before commit):
  ```bash
  go fmt ./...
  ```

- **Run linter** (required before commit):
  ```bash
  golangci-lint run
  ```

### Adding a New Game

If you're adding a new game:

1. Create game logic in `/internal/<gamename>/`
2. Add command implementation in `/cmd/<gamename>.go`
3. Follow the patterns established in existing games
4. Update documentation to include the new game
5. Add comprehensive tests

## Pull Request Process

### Create a Pull Request

1. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Go to the original repository and create a pull request:
   - Click "New Pull Request"
   - Select your branch from your fork
   - Fill out the PR template with details about your changes

### Code Review Process

1. Maintainers will review your PR
2. Address any feedback or requested changes
3. Make sure all CI checks pass
4. Once approved, a maintainer will merge your PR

## Contribution Standards

### Code Style

- Follow standard Go idioms and best practices
- Use meaningful variable and function names
- Keep functions small and focused
- Write GoDoc comments for all exported functions, types, and packages

### Documentation

- Update README.md when adding new features
- Document new commands and options
- Include usage examples
- Document complex algorithms or design decisions

### Testing

- Write unit tests for game logic
- Include integration tests for command behavior
- Mock external dependencies when testing
- All tests must pass before submitting a PR

## Non-Code Contributions

Your contributions don't have to be code! We also welcome:

- **Bug reports**: Create a detailed GitHub issue explaining the bug
- **Feature requests**: Submit ideas through GitHub issues
- **Documentation improvements**: Corrections or additions to docs
- **UI/UX suggestions**: Ideas to improve game interactions
- **User testing**: Provide feedback on game usability

When submitting issues, please use the provided issue templates and include as much detail as possible.

## Community

Our community is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

Thank you for taking the time to contribute to gh-game! We look forward to your contributions!