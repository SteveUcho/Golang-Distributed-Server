package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"sync"
)

// Task struct
type Task struct {
	ID       string
	TaskName string
}

// TodoList is a global variable that hold the current set of tasks
var TodoList sync.Map

// BackendPayload type for having a consistent structure for sending messages to the backend server
type BackendPayload struct {
	Mode     string
	ID       string
	TaskName string
}

// TodoCounter is a global variable that holds a counter for use as global id for task
var TodoCounter = 4

func init() {
	TodoList.Store("1", Task{"1", "Mastering Concurrency in Go"})
	TodoList.Store("2", Task{"2", "Go Design Patterns"})
	TodoList.Store("3", Task{"3", "Black Hat Go"})
}

func main() {
	// retrieve the port number if flag was set, else use the default value
	port := flag.String("listen", "8090", "port number here")
	listener, err := net.Listen("tcp", ":"+*port) // start listening on the port
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	// continuosly search for new connections, when connection accepted spawn a goroutine to handle it
	fmt.Println("Started Server")
	for {
		connection, err := listener.Accept() // accept connection from frontend server
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Accepted connection")
		go handleConnection(connection) // handle the connection in a goroutine to allow for multiple connections to the backend
	}
}

// function for concurrently handeling connection to the server, chooses between choices for the specific connection
func handleConnection(conn net.Conn) {
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	var incomingMsg BackendPayload
	// recieveing the response from the backend through the json decoder
	err := decoder.Decode(&incomingMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch incomingMsg.Mode { // choose function based on the mode sent by front end server
	case "getTasks":
		getTasks(encoder)
	case "createTask":
		createTask(incomingMsg)
	case "updateTask":
		updateTask(incomingMsg)
	case "deleteTask":
		deleteTask(incomingMsg)
	}
}

// serves the client the current list of tasks
func getTasks(encoder *json.Encoder) {
	var realMap = map[string]Task{}
	// this function takes a key value pair of the sync.Map and adds it to the realMap map Variable
	// the convert function returns a bool because sync.Map.range() requires it in order to stop on error
	convert := func(key interface{}, value interface{}) bool {
		keyString, err := key.(string) // verifying the value is of the required type with type assertion
		if !err {
			fmt.Println("Key errror, key is not string")
			return false
		}
		valueTask, err := value.(Task) // verifying the value is of the required type with type assertion
		if !err {
			fmt.Println("Key errror, key is not string")
			return false
		}
		realMap[keyString] = valueTask // insert the value to the new map
		return true
	}
	TodoList.Range(convert)
	err := encoder.Encode(realMap)
	if err != nil {
		fmt.Println(err)
	}
}

// this function is in charge of creating and updating tasks
func createTask(msg BackendPayload) {
	strCount := strconv.Itoa(TodoCounter) // convert integer counter to string
	// TodoList[strCount] = Task{strCount, msg.TaskName} // create new task
	TodoList.Store(strCount, Task{strCount, msg.TaskName})
	TodoCounter = TodoCounter + 1
	return
}

// this function is in charge of creating and updating tasks
func updateTask(msg BackendPayload) {
	incomingTaskID := msg.ID // attempt to retreive the id of the task in the case its an update POST
	// TodoList[incomingTaskID] = Task{incomingTaskID, msg.TaskName} // update task with id provided
	TodoList.Store(incomingTaskID, Task{incomingTaskID, msg.TaskName})
	return
}

// deleteTask takes the iris context, which provides the incoming form data, this function deletes the task if it exists and reloads
func deleteTask(msg BackendPayload) {
	incomingTaskID := msg.ID // attempt to retreive the id of the task in the case its an update POST
	// _, exists := TodoList[incomingTaskID]
	// if exists { // only attempt to delete if it exists
	// 	delete(TodoList, incomingTaskID)
	// }
	TodoList.Delete(incomingTaskID)
	return
}
