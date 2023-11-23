package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// process - класс процесса
type process struct {
	processUUID    uuid.UUID
	parentProcess  *process
	childProcesses []uuid.UUID
}

// newProcess - конструктор класса process
func newProcess(parent *process) *process {
	processUUID := uuid.New()
	fmt.Println("Created a new process with uuid", processUUID)
	return &process{
		processUUID:    processUUID,
		parentProcess:  parent,
		childProcesses: make([]uuid.UUID, 0),
	}
}

// createChild - метод для создания подпроцесса
func (p *process) createChild() *process {
	childProcess := newProcess(p)
	p.childProcesses = append(p.childProcesses, childProcess.processUUID)
	return childProcess
}

// processTable - класс контроллера процессов
type processTable struct {
	processes   map[uuid.UUID]*process
	mainProcess *process
	directory   string
}

// newProcessTable - конструктор класса processTable
func newProcessTable(mainProcess *process, startingDirectory string) *processTable {
	processMap := make(map[uuid.UUID]*process)
	processMap[mainProcess.processUUID] = mainProcess
	return &processTable{
		processes:   processMap,
		mainProcess: mainProcess,
		directory:   startingDirectory,
	}
}

// addProcess - метод для добавления процесса в контроллер
func (pt *processTable) addProcess(processToAdd *process) *process {
	pt.processes[processToAdd.processUUID] = processToAdd
	return processToAdd
}

// terminateProcess - метод для уничтожения процесса и его подпроцессов
func (pt *processTable) terminateProcess(processUUID uuid.UUID) {
	if _, ok := pt.processes[processUUID]; !ok {
		fmt.Println("Process not found")
	} else {
		for _, proc := range pt.processes {
			if proc.parentProcess != nil {
				if proc.parentProcess.processUUID == processUUID {
					delete(pt.processes, proc.processUUID)
				}
			}
		}
		delete(pt.processes, processUUID)
	}
}

// printCurrent - метод для вывода текущей командной строки консоли
func (pt *processTable) printCurrent() {
	fmt.Print(pt.directory + "> ")
}

// changeDirectory - метод для смены директории
func (pt *processTable) changeDirectory(targetDir string) {
	newPath := filepath.Clean(filepath.Join(pt.directory, targetDir))
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		fmt.Println("Invalid directory")
		return
	}
	pt.directory = newPath
}

// printProcesses - метод для вывода списка активных процессов
func (pt *processTable) printProcesses() {
	fmt.Println("Currently active processes:")
	for _, pr := range pt.processes {
		fmt.Println(pr.processUUID)
	}
}

// printPath - метод для вывода пути к текущей директории контроллера
func (pt *processTable) printPath() {
	fmt.Println(pt.directory)
}

// interpretCommand - метод для обработки команды процессом
func (p *process) interpretCommand(splitCommand []string) {
	if len(splitCommand) < 1 {
		fmt.Println("Invalid command arguments")
		return
	}
	switch splitCommand[0] {
	case "echo":
		if len(splitCommand) < 2 {
			fmt.Println("Invalid command arguments")
			return
		}
		fmt.Println(splitCommand[1])
	default:
		fmt.Println("Unknown command")
	}
}

// interpretComplexCommand - метод для обработки команды контроллером
func (pt *processTable) interpretComplexCommand(rawCommand string) {
	switch {
	case strings.Contains(rawCommand, "&"):
		splitCommand := strings.Split(rawCommand, "&")
		if len(splitCommand) >= 2 {
			var parentProcess *process
			if len(splitCommand[1]) == 0 {
				if len(splitCommand[0]) == 0 {
					fmt.Println("Invalid command")
					return
				}
				parentProcess = pt.mainProcess
			} else {
				parentProcess = newProcess(nil)
			}
			pt.addProcess(parentProcess).interpretCommand(strings.Fields(splitCommand[1]))
			pt.addProcess(parentProcess.createChild()).interpretCommand(strings.Fields(splitCommand[0]))

		} else {
			fmt.Println("Invalid command")
		}
	case strings.Contains(rawCommand, "exec"):
		splitCommand := strings.Split(rawCommand, "exec")
		pt.terminateProcess(pt.mainProcess.processUUID)
		pt.mainProcess = pt.addProcess(newProcess(nil))
		pt.interpretComplexCommand(splitCommand[1])
	case strings.Contains(rawCommand, "|"):
		splitCommand := strings.Split(rawCommand, "|")
		for i := range splitCommand {
			splitCommand[i] = strings.TrimSuffix(splitCommand[i], " ")
			splitCommand[i] = strings.TrimPrefix(splitCommand[i], " ")
			pt.interpretComplexCommand(splitCommand[i])
		}
	default:
		splitCommand := strings.Fields(rawCommand)
		if len(splitCommand) < 1 {
			fmt.Println("Invalid command arguments")
			return
		}
		switch splitCommand[0] {
		case "cd":
			if len(splitCommand) < 2 {
				fmt.Println("Invalid command arguments")
				return
			}
			pt.changeDirectory(splitCommand[1])
		case "pwd":
			pt.printPath()
		case "echo":
			if len(splitCommand) < 2 {
				fmt.Println("Invalid command arguments")
				return
			}
			pt.addProcess(newProcess(nil)).interpretCommand(splitCommand)
		case "kill":
			if len(splitCommand) < 2 {
				fmt.Println("Invalid command arguments")
				return
			}
			processUUID, err := uuid.Parse(splitCommand[1])
			if err != nil {
				fmt.Println("Invalid UUID")
				return
			}
			pt.terminateProcess(processUUID)
		case "ps":
			pt.printProcesses()
		default:
			fmt.Println("Unknown command")
		}
	}

}

func main() {

	// определение изначальной директории
	startingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	// создание главного процесса и контроллера
	startingProcess := newProcess(nil)
	console := newProcessTable(startingProcess, startingDirectory)

	reader := bufio.NewReader(os.Stdin) // создание ридера для чтения из консоли

	for {
		console.printCurrent() // вывод текущей директории контроллера

		// чтение введенной из консоли строки
		rawCommand, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error while interpreting command:", err)
		}

		rawCommand = strings.TrimSuffix(rawCommand, "\n")
		rawCommand = strings.TrimSuffix(rawCommand, "\r")

		// если введена команда quit - выход
		if rawCommand == "quit" {
			break
		}

		// обработка полученной команды контроллером
		console.interpretComplexCommand(rawCommand)

	}

}
