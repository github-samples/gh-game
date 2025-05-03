# GitHub Game - Copilot Instructions

## Project Overview
This is a GitHub CLI Extension written in Go that allows users to play various games (coin toss, rock paper scissors, tic tac toe) through the GitHub CLI. It's built using Go and the GitHub CLI extension framework, leveraging the Cobra command library for CLI functionality and Charm libraries for terminal UI rendering.

## Code Standards

### Required Before Commit
- All tests must pass: `go test ./...`
- Code must be properly formatted: `go fmt ./...`
- Linting must pass: `golangci-lint run`
- Ensure documentation is up-to-date for any new commands or features
- Verify that any new game implementation follows the established patterns

### Go Patterns
- Follow standard Go idioms and best practices
- Use GoDoc comments for all exported functions, types, and packages:
  ```go
  // FunctionName does something specific
  // with these parameters and returns this result
  func FunctionName() {}
  ```
- Error handling should follow Go conventions (return errors rather than using exceptions)
- Use meaningful variable and function names that describe their purpose
- Keep functions small and focused on a single responsibility
- Separate interface definitions from implementations where appropriate

## Development Flow

- Build: `go build`
- Test: `go test ./...`
- Lint: `golangci-lint run`
- Format: `go fmt ./...`
- Run: `gh game <subcommand>` (e.g., `gh game cointoss`)

## Repository Structure
- `/cmd`: Main command implementations and CLI structure
  - Command files for each game and core CLI functionality
- `/internal`: Internal packages not intended for external use
  - Game logic implementations
  - Utility functions
- `/pkg`: Public libraries that could potentially be used by other projects
- `/test`: Test utilities and fixtures
- `main.go`: Entry point for the application
- `README.md`: Project documentation
- `LICENSE`: MIT license file
- `go.mod` & `go.sum`: Go module declarations and dependency tracking

## Key Guidelines

1. **User Experience Focus**:
   - Games should be intuitive and easy to play
   - Provide clear instructions and feedback to the user
   - Handle errors gracefully with helpful messages

2. **Command Structure**:
   - New games should be added as subcommands to the main `gh game` command
   - Maintain consistency in command structure and naming

3. **Terminal UI**:
   - Use the Charm libraries (lipgloss, etc.) for consistent UI rendering
   - Ensure games are playable in different terminal environments
   - Consider accessibility in UI design (e.g., color contrast)

4. **Testing**:
   - Write unit tests for game logic
   - Include integration tests for command behavior
   - Mock external dependencies when testing

5. **Documentation**:
   - Update README when adding new games or features
   - Include usage examples for each game
   - Document any complex algorithms or design decisions