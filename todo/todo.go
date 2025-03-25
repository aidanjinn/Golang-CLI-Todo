package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type task struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Date         string `json:"date"`
	Notes        string `json:"notes"`
	Completed    bool   `json:"completed"`
	Dependencies []int  `json:"dependencies"`
}

func loadTasks() ([]task, error) {
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []task{}, nil
		}
		return nil, err
	}

	var tasks []task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func saveTasks(tasks []task) error {
	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("tasks.json", jsonData, 0644)
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func displayTasks() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	clearScreen()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Title", "Done", "Dependencies", "Date"})

	for _, currTask := range tasks {
		t.AppendRows([]table.Row{
			{
				currTask.ID,
				currTask.Title,
				func(completed bool) string {
					if completed {
						return "✅"
					}
					return "❌"
				}(currTask.Completed),
				len(currTask.Dependencies),
				currTask.Date,
			},
		})
		t.AppendSeparator()
	}
	t.Render()
}

func postTask(title string, notes string, dependencies []string) int {

	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return -1
	}

	id := len(tasks)
	var dependenciesList []int

	for _, depIDStr := range dependencies {
		depID, err := strconv.Atoi(depIDStr)
		if err != nil {
			fmt.Println("Invalid dependency ID:", depIDStr)
			continue
		}

		if depID < 0 || depID >= len(tasks) {
			fmt.Printf("Dependency ID %d doesn't exist\n", depID)
			continue
		}

		dependenciesList = append(dependenciesList, depID)
	}

	newTask := task{
		ID:           id,
		Title:        title,
		Date:         time.Now().Format(time.RFC3339),
		Notes:        notes,
		Completed:    false,
		Dependencies: dependenciesList,
	}

	tasks = append(tasks, newTask)

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Error saving tasks:", err)
		return -1
	}

	displayTasks()
	return 0
}

func updateDependencies(tasks []task, updatedTask task) {
	for i := range tasks {
		for j, depID := range tasks[i].Dependencies {
			if depID == updatedTask.ID {
				tasks[i].Dependencies[j] = updatedTask.ID
			}
		}
	}
}

func removeDependency(tasks []task, taskID int) {
	for i := range tasks {
		var newDeps []int
		for _, depID := range tasks[i].Dependencies {
			if depID != taskID {
				newDeps = append(newDeps, depID)
			}
		}
		tasks[i].Dependencies = newDeps
	}
}

func markDone(id int) int {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return -1
	}

	if id < 0 || id >= len(tasks) {
		fmt.Println("Invalid task ID")
		return -1
	}

	tasks[id].Completed = !tasks[id].Completed

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Error saving tasks:", err)
		return -1
	}

	displayTasks()
	return 0
}

func updateTitle(id int, newTitle string) int {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return -1
	}

	if id < 0 || id >= len(tasks) {
		fmt.Println("Invalid task ID")
		return -1
	}

	tasks[id].Title = newTitle
	updateDependencies(tasks, tasks[id])

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Error saving tasks:", err)
		return -1
	}

	return 0
}

func updateNotes(id int, newNotes string) int {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return -1
	}

	if id < 0 || id >= len(tasks) {
		fmt.Println("Invalid task ID")
		return -1
	}

	tasks[id].Notes = newNotes
	updateDependencies(tasks, tasks[id])

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Error saving tasks:", err)
		return -1
	}

	return 0
}

func deleteIndex(id int) int {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return -1
	}

	if id < 0 || id >= len(tasks) {
		fmt.Println("Invalid task ID")
		return -1
	}

	removeDependency(tasks, id)

	tasks = append(tasks[:id], tasks[id+1:]...)

	for i := range tasks {
		tasks[i].ID = i

		for j, depID := range tasks[i].Dependencies {
			if depID > id {
				tasks[i].Dependencies[j] = depID - 1
			}
		}
	}

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Error saving tasks:", err)
		return -1
	}

	displayTasks()
	return 0
}

func deleteCompleted() int {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return -1
	}

	var newTasks []task
	var completedIDs []int

	for _, t := range tasks {
		if t.Completed {
			completedIDs = append(completedIDs, t.ID)
		} else {
			newTasks = append(newTasks, t)
		}
	}

	for i := range newTasks {
		var newDeps []int

		for _, depID := range newTasks[i].Dependencies {
			keep := true
			for _, compID := range completedIDs {
				if depID == compID {
					keep = false
					break
				}
			}
			if keep {
				adjustedID := depID
				for _, compID := range completedIDs {
					if depID > compID {
						adjustedID--
					}
				}
				newDeps = append(newDeps, adjustedID)
			}
		}
		newTasks[i].Dependencies = newDeps
	}

	for i := range newTasks {
		newTasks[i].ID = i
	}

	if err := saveTasks(newTasks); err != nil {
		fmt.Println("Error saving tasks:", err)
		return -1
	}

	displayTasks()
	return 0
}

func deleteAll() int {
	if err := saveTasks([]task{}); err != nil {
		fmt.Println("Error saving tasks:", err)
		return -1
	}
	displayTasks()
	return 0
}

func review(id int) int {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return -1
	}

	if id < 0 || id >= len(tasks) {
		fmt.Println("Invalid task ID")
		return -1
	}

	clearScreen()

	tmp := tasks[id]

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Title", "Done", "Date", "Notes"})
	t.AppendRow(table.Row{
		tmp.ID,
		tmp.Title,
		func(completed bool) string {
			if completed {
				return "✅"
			}
			return "❌"
		}(tmp.Completed),
		tmp.Date,
		tmp.Notes,
	})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Dependencies"})
	t.AppendSeparator()

	if len(tmp.Dependencies) == 0 {
		t.AppendRow(table.Row{"None"})
	} else {
		for _, depID := range tmp.Dependencies {
			if depID >= 0 && depID < len(tasks) {
				depTask := tasks[depID]
				t.AppendRow(table.Row{
					depTask.ID,
					depTask.Title,
					func(completed bool) string {
						if completed {
							return "✅"
						}
						return "❌"
					}(depTask.Completed),
					depTask.Date,
					depTask.Notes,
				})
				t.AppendSeparator()
			}
		}
	}
	t.AppendFooter(table.Row{})
	t.Render()
	return 0
}

func main() {
	for {
		fmt.Print("Enter Command (Use 'Help' for List): ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		command := strings.ToLower(input.Text())
		command = strings.TrimSpace(command)

		switch command {
		case "help":
			fmt.Println("Options: Post, Review, Update, Delete, Display, Help")

		case "post":
			fmt.Print("Enter Task Title: ")
			input.Scan()
			title := input.Text()
			fmt.Print("Enter Any Notes: ")
			input.Scan()
			notes := input.Text()
			fmt.Print("Enter ID for any dependencies (comma separated): ")
			input.Scan()
			id := input.Text()
			idList := strings.Split(id, ",")
			postTask(title, notes, idList)

		case "update":
			fmt.Print("Enter Task ID: ")
			input.Scan()
			id, err := strconv.Atoi(input.Text())
			if err != nil {
				fmt.Println("Invalid ID:", err)
				continue
			}

			if review(id) == -1 {
				continue
			}

			fmt.Print("Done, Title, or Notes: ")
			input.Scan()
			selection := strings.ToLower(input.Text())
			selection = strings.TrimSpace(selection)

			switch selection {
			case "done":
				markDone(id)
			case "title":
				fmt.Print("Enter New Task Title: ")
				input.Scan()
				newTitle := input.Text()
				updateTitle(id, newTitle)
			case "notes":
				fmt.Print("Enter New Task Notes: ")
				input.Scan()
				newNotes := input.Text()
				updateNotes(id, newNotes)
			}
			displayTasks()

		case "display":
			displayTasks()

		case "review":
			fmt.Print("Enter Task ID: ")
			input.Scan()
			id, err := strconv.Atoi(input.Text())
			if err != nil {
				fmt.Println("Invalid ID:", err)
				continue
			}
			review(id)

		case "delete":
			fmt.Print("Delete: Done, ID, All: ")
			input.Scan()
			selection := strings.ToLower(input.Text())
			selection = strings.TrimSpace(selection)

			switch selection {
			case "id":
				fmt.Print("Enter Task ID: ")
				input.Scan()
				id, err := strconv.Atoi(input.Text())
				if err != nil {
					fmt.Println("Invalid ID:", err)
					continue
				}
				deleteIndex(id)
			case "done":
				deleteCompleted()
			case "all":
				deleteAll()
			}

		case "exit":
			return

		default:
			fmt.Println("Unknown command. Available commands: post, update, display, review, delete, exit")
		}
	}
}
