package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var todosFile string

func init() {
	home, _ := os.UserHomeDir()
	todosFile = filepath.Join(home, ".todos.json")
}

func loadTodos() ([]Todo, error) {
	data, err := os.ReadFile(todosFile)
	if os.IsNotExist(err) {
		return []Todo{}, nil
	}
	if err != nil {
		return nil, err
	}
	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

func saveTodos(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(todosFile, data, 0644)
}

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A simple todo CLI",
}

var addCmd = &cobra.Command{
	Use:   "add [text]",
	Short: "Add a new todo",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		todos, err := loadTodos()
		if err != nil {
			return err
		}
		id := 1
		if len(todos) > 0 {
			id = todos[len(todos)-1].ID + 1
		}
		todos = append(todos, Todo{ID: id, Text: args[0], Done: false})
		if err := saveTodos(todos); err != nil {
			return err
		}
		fmt.Printf("Added: [%d] %s\n", id, args[0])
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all todos",
	RunE: func(cmd *cobra.Command, args []string) error {
		todos, err := loadTodos()
		if err != nil {
			return err
		}
		if len(todos) == 0 {
			fmt.Println("No todos")
			return nil
		}
		for _, t := range todos {
			status := " "
			if t.Done {
				status = "x"
			}
			fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Text)
		}
		return nil
	},
}

var doneCmd = &cobra.Command{
	Use:   "done [id]",
	Short: "Mark a todo as done",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var id int
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
			return fmt.Errorf("invalid id: %s", args[0])
		}
		todos, err := loadTodos()
		if err != nil {
			return err
		}
		for i, t := range todos {
			if t.ID == id {
				todos[i].Done = true
				if err := saveTodos(todos); err != nil {
					return err
				}
				fmt.Printf("Completed: [%d] %s\n", id, t.Text)
				return nil
			}
		}
		return fmt.Errorf("todo %d not found", id)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doneCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
