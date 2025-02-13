# gh-game

A GitHub CLI extension that allows you to play games through the GitHub CLI.

## Installation

```sh
gh extension install chrisreddington/gh-game
```

## Commands

### Coin Toss

Play a coin toss game where you try to guess whether the coin will land on heads or tails. Keep your streak going by guessing correctly!

```sh
gh game cointoss heads  # or tails
```

The game will continue as long as you keep guessing correctly, allowing you to build up a streak. You can quit at any time by selecting "Quit" when prompted for your next guess.

### Whoami

Display information about the currently authenticated GitHub user.

```sh
gh game whoami
```

## Development

### Prerequisites

- Go 1.23 or newer
- GitHub CLI installed

### Building from source

1. Clone the repository
2. Run `go build` to build the extension
3. Run `go test ./...` to run the tests

### Development Container

This project includes a [development container configuration](.devcontainer/devcontainer.json) for VS Code, which provides a consistent development environment with:
- Go 1.23
- GitHub CLI
- Several VS Code extensions for GitHub

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is available as open source under the terms of the [MIT License](LICENSE).