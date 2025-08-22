package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"time"
)

const (
	taskJson   = "tasks.json"
	charset    = "abcdefghijklmnopqrstuvwxyz0123456789"
	eventsJson = "events.json"
)

var commands = []string{"add", "list", "done", "delete", "close", "help", "showLogs"}

type Task struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAT   time.Time `json:"created_at"`
	CompleteAT  time.Time `json:"update_at"`
}

func help() {
	fmt.Println("chose the action: add, list, done, delete, close, showLogs")
}

func printTasks(title string, descrTasks []string) {
	fmt.Println(title + ":")
	for _, val := range descrTasks {
		fmt.Printf("%v", val)
	}
	fmt.Println()
}

func idGen() string {
	lenght := 5
	b := make([]byte, lenght)
	for i := range b {
		randInd := rand.Intn(len(charset))
		b[i] = charset[randInd]
	}
	return string(b)
}

func add(str string) error {
	tasks, err := loadItems[Task](taskJson)
	if err != nil {
		return fmt.Errorf("parsing error in add: %w", err)
	}

	id := idGen()
	descr := str
	now := time.Now()
	compl := false

	t := Task{
		ID:          id,
		Description: descr,
		Completed:   compl,
		CreatedAT:   now,
		CompleteAT:  now,
	}

	tasks = append(tasks, t)

	err = saveItems(taskJson, tasks)
	if err != nil {
		return fmt.Errorf("save error in add: %w", err)
	}
	fmt.Printf("task `%s` added with ID: %s in %v\n", str, t.ID, now.Format("2006-01-02 15:04"))
	return nil
}

func list(scanner *bufio.Scanner) error {
	ans := 0
	fmt.Printf("show UNCOMPLETE tasks (1)\nshow COMPLETE tasks (2)\nshow ALL tasks (3)\n")

	scanner.Scan()
	ans, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("convertion error in list: %w", err)
	}
	tasks, err := loadItems[Task](taskJson) //parseFile(taskJson)
	if err != nil {
		return fmt.Errorf("parsing error in list: %w", err)
	}
	var compleTasks []string
	var unCompleTasks []string

	for _, val := range tasks {
		if val.Completed {
			compleTasks = append(compleTasks, val.Description+fmt.Sprintf(" (ID: %s)", val.ID)+"\n")
		} else {
			unCompleTasks = append(unCompleTasks, val.Description+fmt.Sprintf(" (ID: %s)", val.ID)+"\n")
		}
	}

	switch ans {
	case 1:
		if len(unCompleTasks) != 0 {
			printTasks("uncomplete tasks: ", unCompleTasks)
		} else {
			fmt.Printf("no complete tasks\n\n")
		}

	case 2:
		if len(compleTasks) != 0 {
			printTasks("complete tasks: ", compleTasks)
		} else {
			fmt.Printf("no complete tasks\n\n")
		}

	case 3:
		if len(compleTasks) != 0 && len(unCompleTasks) != 0 {
			printTasks("complete tasks", compleTasks)
			printTasks("uncomplete tasks", unCompleTasks)
		} else if len(compleTasks) == 0 && len(unCompleTasks) != 0 {
			printTasks("uncomplete tasks", unCompleTasks)
			fmt.Printf("no complete tasks\n\n")
		} else if len(compleTasks) != 0 && len(unCompleTasks) == 0 {
			printTasks("complete tasks", compleTasks)
			fmt.Printf("no uncomplete tasks\n\n")
		} else {
			fmt.Println("no tasks")
		}
	default:
		fmt.Println("no such action, retry now")
	}
	return nil
}

func done(ID string) error {
	tasks, err := loadItems[Task](taskJson) //parseFile(taskJson)
	if err != nil {
		return fmt.Errorf("parse error in complete %w", err)

	}
	for ind, val := range tasks {
		if val.ID == ID {
			if !tasks[ind].Completed {
				tasks[ind].Completed = true
				tasks[ind].CompleteAT = time.Now()
				fmt.Printf("task %#v complete \n", tasks[ind].Description)
			} else {
				fmt.Printf("task %#v already complete\n", val.Description)
			}
			err = saveItems(taskJson, tasks)
			if err != nil {
				return fmt.Errorf("save error in complete: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("no task that such ID: %s", ID)
}

func del(ID string) error {
	tasks, err := loadItems[Task](taskJson)
	if err != nil {
		return fmt.Errorf("parse error in delete: %w", err)
	}
	for ind, val := range tasks {
		if val.ID == ID {
			tasks = append(tasks[:ind], tasks[ind+1:]...)
			if err = saveItems(taskJson, tasks); err != nil {
				return fmt.Errorf("error when saving after deletion: %w", err)
			} // error
			fmt.Printf("task %#v sucesesfully delete\n", val.Description)
			return nil
		}

	}

	return fmt.Errorf("no task that such ID: %s", ID)
}

func forAction() string {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		action := scanner.Text()
		if !slices.Contains(commands, action) {
			fmt.Println("no such command, retry now")
			if logErr := eventLogger(eventsJson, action, fmt.Errorf("invalid command")); logErr != nil {
				fmt.Println("unsucesessful save event")
			}
			continue
		}
		return action
	}

}

func handleUserInput(scanner *bufio.Scanner, commName string, prompt string, actionFunc func(string) error) error {
	fmt.Println(prompt)
	scanner.Scan()
	userInput := scanner.Text()
	fullInput := fmt.Sprintf("%s %s", commName, userInput)

	if err := actionFunc(userInput); err != nil {
		fmt.Printf("%s error, try again\n", commName)
		logErr := eventLogger(eventsJson, fullInput, err)
		if logErr != nil {
			fmt.Println("unsucesessful save event + penis")
			fmt.Println(logErr)
		}
		return err
	}
	logErr := eventLogger(eventsJson, fullInput, nil)
	if logErr != nil {
		fmt.Println("unsucesessful save event")
	}
	return nil
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("chose the action: add, list, done, delete, close, showLogs")
	for {
		action := forAction()

		switch action {
		case "add":

			if err := handleUserInput(scanner, "add", "take the describe task for add: ", add); err != nil {
				continue
			}

		case "list":
			err := list(scanner)

			if err != nil {
				if logErr := eventLogger(eventsJson, "list", err); logErr != nil {
					fmt.Println("unsucesessful save event")
				}
				continue
			}
			if logErr := eventLogger(eventsJson, "list", nil); logErr != nil {
				fmt.Println("unsucesessful save event")
			}

		case "done":

			if err := handleUserInput(scanner, "done", "take the ID task for complete: ", done); err != nil {
				continue
			}

		case "delete":

			if err := handleUserInput(scanner, "delete", "take the ID task for delete: ", del); err != nil {
				continue
			}
		case "close":
			if logErr := eventLogger(eventsJson, "close", nil); logErr != nil {
				fmt.Println("unsucesessful save event")
			}
			os.Exit(0)
		case "help":
			if logErr := eventLogger(eventsJson, "help", nil); logErr != nil {
				fmt.Println("unsucesessful save event")
			}
			help()
		case "showLogs":
			if err := showEvents(eventsJson); err != nil {
				fmt.Println("showLogs error, try again")
				if logErr := eventLogger(eventsJson, "showLogs", err); logErr != nil {
					fmt.Println("unsucesessful save event")
				}
				continue
			}
			if logErr := eventLogger(eventsJson, "showLogs", nil); logErr != nil {
				fmt.Println("unsucesessful save event")
			}
		default:
			continue
		}
	}

}

type Event struct {
	CreatedAT string `json:"created_at"`
	UserInput string `json:"user_input"`
	ErrorText string `json:"error_text"`
}

func eventLogger(fileName string, text string, err error) error {
	var errText string
	if err != nil {
		errText = err.Error()
	} else {
		errText = ""
	}
	event := Event{
		CreatedAT: time.Now().Format("2006-01-02 15:04"),
		UserInput: text,
		ErrorText: errText,
	}
	events, err := loadItems[Event](fileName)
	if err != nil {
		return fmt.Errorf("parsing error in eventLogger: %w", err)
	}
	events = append(events, event)
	if err = saveItems(fileName, events); err != nil {
		return fmt.Errorf("saving error in saveEvents: %w", err)
	}
	return nil
}

func showEvents(filename string) error {
	events, err := loadItems[Event](filename)
	if err != nil {
		return fmt.Errorf("parsing error in showEvents: %w", err)
	}
	jsonData, err := json.MarshalIndent(events, "", "\t")
	if err != nil {
		return fmt.Errorf("converse error in showEvents")
	}
	fmt.Println(string(jsonData))
	return nil
}

func loadItems[T Event | Task](filename string) ([]T, error) {
	var items []T

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return items, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		items = []T{}
	} else {
		err = json.Unmarshal(data, &items)
		if err != nil {
			return nil, err
		}

	}
	return items, nil
}

func saveItems[T Event | Task](filename string, items []T) error {
	jsonData, err := json.MarshalIndent(items, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, jsonData, 0600)
	if err != nil {
		return err
	}
	return nil
}
