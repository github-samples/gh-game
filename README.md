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

### Rock Paper Scissors

Play Rock Paper Scissors against the computer. Best of 3, 5, 7, or 9 rounds.

```sh
gh game rockpaperscissors
```

Choose your move (rock, paper, or scissors) in each round, and the computer will randomly select its move. The game follows standard Rock Paper Scissors rules:
- Rock crushes Scissors
- Scissors cuts Paper
- Paper covers Rock

### Tic Tac Toe

Play the classic game of Tic Tac Toe against another player. Players take turns placing X's and O's on a 3x3 grid, trying to get three in a row horizontally, vertically, or diagonally.

```sh
gh game tictactoe
```

The game provides an interactive interface where you can select positions on the board using numbers 1-9, corresponding to the grid positions from left to right, top to bottom.

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