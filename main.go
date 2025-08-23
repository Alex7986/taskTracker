package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/fatih/color"

	l "github.com/Alex7986/taskTracker/logger"
	s "github.com/Alex7986/taskTracker/store"
)

var commands = []string{"add", "list", "done", "delete", "close", "help", "showLogs"}

var c = color.New(color.FgCyan)
var r = color.New(color.FgRed)
var g = color.New(color.FgGreen)
var b = color.New(color.FgBlue)

func help() {
	c.Println("chose the action: add, list, done, delete, close, showLogs, help")
}

func printTasks(title string, descrTasks []string) {
	b.Println(title + ":")
	for _, val := range descrTasks {
		b.Printf("%v", val)
	}
	fmt.Println()
}

func add(str string) error {
	tasks, err := s.LoadItems[s.Task](s.TaskJson)
	if err != nil {
		return fmt.Errorf("parsing error in add: %w", err)
	}

	id := s.IdGen()
	descr := str
	now := time.Now()
	compl := false

	t := s.Task{
		ID:          id,
		Description: descr,
		Completed:   compl,
		CreatedAT:   now,
		UpdateAT:    now,
	}

	tasks = append(tasks, t)

	err = s.SaveItems(s.TaskJson, tasks)
	if err != nil {
		return fmt.Errorf("save error in add: %w", err)
	}
	g.Printf("task `%s` added with ID: %s in %v\n", str, t.ID, now.Format("2006-01-02 15:04"))
	return nil
}

func list(scanner *bufio.Scanner) error {
	ans := 0
	c.Printf("show UNCOMPLETE tasks (1)\nshow COMPLETE tasks (2)\nshow ALL tasks (3)\n")

	scanner.Scan()
	ans, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("convertion error in list: %w", err)
	}
	tasks, err := s.LoadItems[s.Task](s.TaskJson) //parseFile(taskJson)
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
			b.Printf("no complete tasks\n\n")
		}

	case 2:
		if len(compleTasks) != 0 {
			printTasks("complete tasks: ", compleTasks)
		} else {
			b.Printf("no complete tasks\n\n")
		}

	case 3:
		if len(compleTasks) != 0 && len(unCompleTasks) != 0 {
			printTasks("complete tasks", compleTasks)
			printTasks("uncomplete tasks", unCompleTasks)
		} else if len(compleTasks) == 0 && len(unCompleTasks) != 0 {
			printTasks("uncomplete tasks", unCompleTasks)
			b.Printf("no complete tasks\n\n")
		} else if len(compleTasks) != 0 && len(unCompleTasks) == 0 {
			printTasks("complete tasks", compleTasks)
			b.Printf("no uncomplete tasks\n\n")
		} else {
			b.Println("no tasks")
		}
	default:
		r.Println("no such action, try again")
	}
	return nil
}

func done(ID string) error {
	tasks, err := s.LoadItems[s.Task](s.TaskJson) //parseFile(taskJson)
	if err != nil {
		return fmt.Errorf("parse error in complete %w", err)

	}
	for ind, val := range tasks {
		if val.ID == ID {
			if !tasks[ind].Completed {
				tasks[ind].Completed = true
				tasks[ind].UpdateAT = time.Now()
				g.Printf("task %#v complete \n", tasks[ind].Description)
			} else {
				r.Printf("task %#v already complete\n", val.Description)
			}
			err = s.SaveItems(s.TaskJson, tasks)
			if err != nil {
				return fmt.Errorf("save error in complete: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("no task that such ID: %s", ID)
}

func del(ID string) error {
	tasks, err := s.LoadItems[s.Task](s.TaskJson)
	if err != nil {
		return fmt.Errorf("parse error in delete: %w", err)
	}
	for ind, val := range tasks {
		if val.ID == ID {
			tasks = append(tasks[:ind], tasks[ind+1:]...)
			if err = s.SaveItems(s.TaskJson, tasks); err != nil {
				return fmt.Errorf("error when saving after deletion: %w", err)
			} // error
			g.Printf("task %#v sucesesfully delete\n", val.Description)
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
			r.Println("no such command, try again")
			if logErr := l.Log(s.EventsJson, action, fmt.Errorf("invalid command")); logErr != nil {
				r.Println("unsucesessful save event")
			}
			continue
		}
		return action
	}

}

func handleUserInput(scanner *bufio.Scanner, commName string, prompt string, actionFunc func(string) error) error {
	g.Println(prompt)
	scanner.Scan()
	userInput := scanner.Text()
	fullInput := fmt.Sprintf("%s %s", commName, userInput)

	if err := actionFunc(userInput); err != nil {
		r.Printf("%s error, try again\n", commName)
		logErr := l.Log(s.EventsJson, fullInput, err)
		if logErr != nil {
			r.Println("unsucesessful save event + penis")
			r.Println(logErr)
		}
		return err
	}
	logErr := l.Log(s.EventsJson, fullInput, nil)
	if logErr != nil {
		r.Println("unsucesessful save event")
	}
	return nil
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	c.Println("chose the action: add, list, done, delete, close, showLogs, help")
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
				if logErr := l.Log(s.EventsJson, "list", err); logErr != nil {
					fmt.Println("unsucesessful save event")
				}
				r.Println("no such command, try again")
				continue
			}
			if logErr := l.Log(s.EventsJson, "list", nil); logErr != nil {
				r.Println("unsucesessful save event")
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
			if logErr := l.Log(s.EventsJson, "close", nil); logErr != nil {
				r.Println("unsucesessful save event")
			}
			os.Exit(0)
		case "help":
			if logErr := l.Log(s.EventsJson, "help", nil); logErr != nil {
				r.Println("unsucesessful save event")
			}
			help()
		case "showLogs":
			if err := l.ShowEvents(s.EventsJson); err != nil {
				r.Println("showLogs error, try again")
				if logErr := l.Log(s.EventsJson, "showLogs", err); logErr != nil {
					r.Println("unsucesessful save event")
				}
				continue
			}
			if logErr := l.Log(s.EventsJson, "showLogs", nil); logErr != nil {
				r.Println("unsucesessful save event")
			}
		default:
			continue
		}
	}

}
