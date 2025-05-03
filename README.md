# gh-game

A GitHub CLI extension that allows you to play games through the GitHub CLI.

## Features

- **Multiple Games**: Play coin toss, higher or lower, rock paper scissors, tic tac toe, and word guess games
- **Terminal-Based**: Fully playable through your terminal using GitHub CLI
- **Interactive UI**: User-friendly terminal interfaces for all games
- **Score Tracking**: Keep track of your scores and streaks in supported games

## Background

The gh-game extension serves as both a fun diversion and a showcase of GitHub CLI extension capabilities.

### Roadmap

- Additional games to be added in the future

## Installation

```sh
gh extension install github-samples/gh-game
```

## Requirements

- Go 1.23 or newer
- GitHub CLI installed
- Terminal with Unicode support for optimal experience

### Setting Up Development Environment

1. Clone the repository
2. Run `go build` to build the extension
3. Run `go test ./...` to run the tests

### Development Container

This project includes a [development container configuration](.devcontainer/devcontainer.json) for VS Code, which provides a consistent development environment with:
- Go 1.23
- GitHub CLI
- Several VS Code extensions for GitHub

## Commands

### Coin Toss

Play a coin toss game where you try to guess whether the coin will land on heads or tails. Keep your streak going by guessing correctly!

```sh
gh game cointoss heads  # or tails
```

The game will continue as long as you keep guessing correctly, allowing you to build up a streak. You can quit at any time by selecting "Quit" when prompted for your next guess.

### Higher or Lower

Play a number guessing game where you predict if the next random number will be higher or lower than the current one. See how long you can maintain your streak of correct guesses!

```sh
gh game higherlower
```

Optional flags:
- `--min` or `-m`: Set the minimum possible number (default: 1)
- `--max` or `-M`: Set the maximum possible number (default: 100)

Example with custom range:
```sh
gh game higherlower --min 1 --max 1000
```

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

### Word Guess

Play a word guessing game where you guess a GitHub-related term one letter at a time. Try to reveal the word before running out of guesses!

```sh
gh game wordguess
```

The game selects a random GitHub-related term, and you need to guess it by suggesting one letter at a time. Each correct letter is revealed in its position. Each incorrect guess reduces your remaining guesses. You win by guessing the complete word before making 6 incorrect guesses.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. Check out our [contributing guidelines](CONTRIBUTING.md) for more details on how to get involved.

## Maintainers

This project is maintained by the GitHub Developer Relations team. See [CODEOWNERS](CODEOWNERS) file for specific maintainers.

## License

This project is available as open source under the terms of the [MIT License](LICENSE).
