package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/chrisreddington/gh-game/internal/memorygame"
	userPrompt "github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/spf13/cobra"
)

// memoryGameCmd represents the memory game command
var memoryGameCmd = &cobra.Command{
	Use:   "memorygame",
	Short: "Test your memory by remembering sequences of colors",
	Long: `A memory game where you need to remember and reproduce sequences of colors.
The sequence grows longer with each round. You can choose your number of lives.`,
	Run: func(cmd *cobra.Command, args []string) {
		runMemoryGame()
	},
}

func init() {
	rootCmd.AddCommand(memoryGameCmd)
}

// clearScreen clears the terminal screen in a cross-platform way
func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func runMemoryGame() {
	prompter := userPrompt.New(os.Stdin, os.Stdout, os.Stderr)

	// 1. Select number of lives
	lifeOptions := []string{"1 life (hardcore mode)", "2 lives", "3 lives"}
	lifeValues := []int{1, 2, 3}
	lifeIdx, err := prompter.Select("Choose number of lives:", "2 lives", lifeOptions)
	if err != nil {
		fmt.Printf("Error selecting lives: %v\n", err)
		return
	}
	lives := lifeValues[lifeIdx]

	game := memorygame.NewGame(lives)
	colorStyles := map[memorygame.Color]lipgloss.Style{
		memorygame.Red:    lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true),
		memorygame.Yellow: lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true),
		memorygame.Green:  lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true),
		memorygame.Blue:   lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),
	}
	colorOptions := []string{"Red", "Yellow", "Green", "Blue"}
	colorMap := map[string]memorygame.Color{
		"Red":    memorygame.Red,
		"Yellow": memorygame.Yellow,
		"Green":  memorygame.Green,
		"Blue":   memorygame.Blue,
	}

	for !game.IsGameOver() {
		clearScreen()
		fmt.Printf("Round %d - Lives: %d\n\n", game.CurrentRound, game.Lives)
		game.GenerateSequence()

		// Show the sequence
		fmt.Println(lipgloss.NewStyle().Bold(true).Render("Remember this sequence:"))
		for _, color := range game.Sequence {
			cstr := string(color)
			if len(cstr) > 0 {
				cstr = strings.ToUpper(cstr[:1]) + cstr[1:]
			}
			fmt.Print(colorStyles[color].Render(cstr) + " ")
		}
		fmt.Println("\n\nMemorize it! It will disappear in 3 seconds...")
		time.Sleep(3 * time.Second)

		// Fail-fast input loop
		for {
			clearScreen()
			fmt.Printf("Round %d - Lives: %d\n\n", game.CurrentRound, game.Lives)
			fmt.Println("Select the sequence in order using the menu.")
			userSequence := make([]memorygame.Color, 0, len(game.Sequence))
			failed := false
			for i := 0; i < len(game.Sequence); i++ {
				prompt := fmt.Sprintf("Color %d:", i+1)
				idx, err := prompter.Select(prompt, colorOptions[0], colorOptions)
				if err != nil {
					fmt.Printf("Error selecting color: %v\n", err)
					return
				}
				picked := colorMap[colorOptions[idx]]
				userSequence = append(userSequence, picked)
				if picked != game.Sequence[i] {
					game.DecrementLives()
					if game.IsGameOver() {
						fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).Render(
							fmt.Sprintf("Game Over! You reached round %d.", game.CurrentRound),
						))
						return
					}
					fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).Render("Wrong color! Try again from the start of this round..."))
					time.Sleep(2 * time.Second)
					// Show the sequence again before retrying
					clearScreen()
					fmt.Printf("Round %d - Lives: %d\n\n", game.CurrentRound, game.Lives)
					fmt.Println(lipgloss.NewStyle().Bold(true).Render("Remember this sequence:"))
					for _, color := range game.Sequence {
						cstr := string(color)
						if len(cstr) > 0 {
							cstr = strings.ToUpper(cstr[:1]) + cstr[1:]
						}
						fmt.Print(colorStyles[color].Render(cstr) + " ")
					}
					fmt.Println("\n\nMemorize it! It will disappear in 3 seconds...")
					time.Sleep(3 * time.Second)
					failed = true
					break
				}
			}
			if !failed {
				fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true).Render("Correct! Next round..."))
				game.NextRound()
				time.Sleep(1 * time.Second)
				break
			}
		}
	}

	fmt.Println(lipgloss.NewStyle().Bold(true).Render("Thanks for playing!"))
}
