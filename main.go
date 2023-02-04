package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type TodoList struct {
	Todo []TodoItem `json:"todo_list"`
}

type TodoItem struct {
	Id     int    `json:"id"`
	Task   string `json:"task"`
	Status string `json:"status"`
}

func ListIncompleteTasks(todo_list *TodoList) {
	fmt.Println("TODO List:")
	for i := 0; i < len(todo_list.Todo); i++ {
		if todo_list.Todo[i].Status == "incomplete" {
			fmt.Println("	-" + strconv.Itoa(todo_list.Todo[i].Id) + "- " + todo_list.Todo[i].Task)
		}
	}
}

func ListCompleteTasks(todo_list *TodoList) {
	fmt.Println("TODO List:")
	for i := 0; i < len(todo_list.Todo); i++ {
		if todo_list.Todo[i].Status == "complete" {
			fmt.Println("	-" + strconv.Itoa(todo_list.Todo[i].Id) + "- " + todo_list.Todo[i].Task + " ✓")
		}
	}
}

func reorderTaskId(todo_list *TodoList) *TodoList {
	for i := 0; i < len(todo_list.Todo); i++ {
		todo_list.Todo[i].Id = i + 1
	}
	return todo_list
}

const helpMessage string = "Current Program Arguments: \n\thelp -- Displays this message\n\tlist -- Lists complete or incomplete tasks ex: \"list complete\"\n\tadd -- adds a task to your list ex: \"add wash dishes\"\n\tremove -- Removes a task based on its number on the list \"remove 4\""

func main() {
	jsonFile, err := os.Open("task.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var todo_list TodoList
	json.Unmarshal(byteValue, &todo_list)

	if len(os.Args[1:]) == 0 {
		ListIncompleteTasks(&todo_list)
	} else {
		if strings.EqualFold(os.Args[1], "list") {
			if len(os.Args[2:]) == 0 {
				ListIncompleteTasks(&todo_list)
			} else if strings.EqualFold(os.Args[2], "incomplete") {
				ListIncompleteTasks(&todo_list)
			} else if strings.EqualFold(os.Args[2], "complete") {
				ListCompleteTasks(&todo_list)
			}
		} else if strings.EqualFold(os.Args[1], "help") {
			fmt.Println(helpMessage)
		} else if strings.EqualFold(os.Args[1], "add") {
			compundArgs := ""
			for _, arg := range os.Args[2:] {
				compundArgs = compundArgs + arg + " "
			}
			newTask := TodoItem{}
			newTask.Id = len(todo_list.Todo) + 1
			newTask.Task = compundArgs
			newTask.Status = "incomplete"
			todo_list.Todo = append(todo_list.Todo, newTask)
			updatedJSON, err := json.MarshalIndent(todo_list, "", "    ")
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}
			err = ioutil.WriteFile("task.json", updatedJSON, 0644)
			if err != nil {
				fmt.Println("Error writing JSON file:", err)
				return
			}
			fmt.Println("\"" + compundArgs + "\" added to your todo list")
		} else if strings.EqualFold(os.Args[1], "remove") {
			indexToRemove, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Println("Error converting string to int:", err)
			}
			if indexToRemove == len(todo_list.Todo) {
				todo_list.Todo = todo_list.Todo[:indexToRemove-1]
			} else {
				todo_list.Todo = append(todo_list.Todo[:indexToRemove], todo_list.Todo[indexToRemove+1:]...)
			}
			updatedJSON, err := json.MarshalIndent(reorderTaskId(&todo_list), "", "    ")
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}
			err = ioutil.WriteFile("task.json", updatedJSON, 0644)
			if err != nil {
				fmt.Println("Error writing JSON file:", err)
				return
			}
			fmt.Println(strconv.Itoa(indexToRemove) + " removed from task.json file successfully!")
		} else if strings.EqualFold(os.Args[1], "mark") {
			indexToMark, err := strconv.Atoi(os.Args[2])
			tempvar := todo_list.Todo[indexToMark-1]
			if err != nil {
				fmt.Println("Error converting string to int:", err)
			}
			for i := 0; i <= len(todo_list.Todo); i++ {
				if i == indexToMark {
					if strings.EqualFold(os.Args[3], "complete") {
						todo_list.Todo[i-1].Status = "complete"
						tempvar := todo_list.Todo[indexToMark-1]
						todo_list.Todo = append(todo_list.Todo[:indexToMark-1], todo_list.Todo[indexToMark:]...)
						todo_list.Todo = append(todo_list.Todo, tempvar)
					} else if strings.EqualFold(os.Args[3], "incomplete") {
						todo_list.Todo[i-1].Status = "incomplete"
						tempvar := todo_list.Todo[indexToMark-1]
						todo_list.Todo = append(todo_list.Todo[:indexToMark-1], todo_list.Todo[indexToMark:]...)
						todo_list.Todo = append([]TodoItem{tempvar}, todo_list.Todo...)
					} else {
						fmt.Println("Please use command as \"mark 4 complete/incomplete\"")
						os.Exit(0)
					}
				}
			}
			updatedJSON, err := json.MarshalIndent(reorderTaskId(&todo_list), "", "    ")
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}
			err = ioutil.WriteFile("task.json", updatedJSON, 0644)
			if err != nil {
				fmt.Println("Error writing JSON file:", err)
				return
			}
			fmt.Println("Seccesfuly changed the status of \"" + tempvar.Task + "\"")
			// fmt.Println("Seccesfuly marked \"" + todo_list.Todo[indexToMark-1].Task + "\" as " + todo_list.Todo[indexToMark-1].Status)
		} else if strings.EqualFold(os.Args[1], "reorder") {
			updatedJSON, err := json.MarshalIndent(reorderTaskId(&todo_list), "", "    ")
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}
			err = ioutil.WriteFile("task.json", updatedJSON, 0644)
			if err != nil {
				fmt.Println("Error writing JSON file:", err)
				return
			}
			fmt.Println("Successfully re-ordered the whole task list")
		}
	}
}

// arg := os.Args[1:]
// fmt.Println(arg)
// fmt.Println("Id: " + strconv.Itoa(todo_list.Todo[i].Id))
// fmt.Println("Task: " + todo_list.Todo[i].Task)
// fmt.Println("Status: " + todo_list.Todo[i].Status + "\n")
// content, err := json.Marshal(newTask)
// if err != nil {
// 	fmt.Println(err)
// }
// err = ioutil.WriteFile("task.json", content, 0644)
// if err != nil {
// 	log.Fatal(err)
// }
