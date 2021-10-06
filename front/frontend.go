package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"

	"github.com/kataras/iris/v12"
)

var (
	endpoint *string
)

// Task type for json unmashaling
type Task struct {
	ID       string
	TaskName string
}

// TaskList type which hold Tasks for unmarshalling
type TaskList map[string]Task

// BackendPayload type for having a consistent structure for sending messages to the backend server
type BackendPayload struct {
	Mode     string
	ID       string
	TaskName string
}

func main() {
	app := iris.New()

	// set server listen to port
	port := flag.String("listen", "8080", "port number here")
	// attempt connect to the backend
	endpoint = flag.String("backend", ":8090", "port number here")

	// Parse all templates from the "./views" folder
	// where extension is ".html" and parse them
	// using the standard `html/template` package.
	tmpl := iris.HTML("./views", ".html")
	// Set custom delimeters.
	tmpl.Delims("{{", "}}")
	// Register the view engine to the views,
	// this will load the templates.
	app.RegisterView(tmpl)

	// GET endpoint for loading the front page
	app.Get("/", getFromBackEnd, getTasks)
	// POST used for creating or deleting a task
	app.Post("/create-or-update-task", createUpdateTask)
	// DELETE used for deleting
	app.Post("/delete-task", deleteTask)

	app.Listen(":" + *port)
}

// gets the actual data from the back end server
func getFromBackEnd(ctx iris.Context) {
	connection, err := net.Dial("tcp", *endpoint)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()
	decoder := json.NewDecoder(connection)
	encoder := json.NewEncoder(connection)

	message := BackendPayload{"getTasks", "", ""}
	// sending message to connection via the json encoder
	err = encoder.Encode(&message)
	if err != nil {
		fmt.Println(err)
		return
	}
	var response TaskList
	// recieveing the response from the backend through the json decoder
	err = decoder.Decode(&response)
	if err != nil {
		fmt.Println(err)
		return
	}
	// saving thr response for the next handler in the series
	ctx.Values().Set("backendMessage", response)
	ctx.Next()
}

// serves the client the current list of tasks
func getTasks(ctx iris.Context) {
	// gets the data that was stored in context.values from the backend
	TodoList, _ := ctx.Values().Get("backendMessage").(TaskList)
	// Bind: {{range .Todos}} with the list of todos
	ctx.ViewData("Todos", TodoList)
	// Render template file: ./views/index.html
	ctx.View("index.html")
}

// this function is in charge of creating and updating tasks
func createUpdateTask(ctx iris.Context) {
	connection, err := net.Dial("tcp", *endpoint)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()
	encoder := json.NewEncoder(connection)

	var message BackendPayload
	form := ctx.FormValues()
	incomingTask := form["taskname"][0] // get value of form input with name "taskname"
	if incomingTask == "" {             // check if got a task name, empty task name in update taskname does nothing! theres a delete button right there!
		ctx.Redirect("/", iris.StatusFound) // reload page for the client
		return
	}
	if len(form["id"]) == 0 { // check if id field exists, if no value was recieved then its a create task
		message = BackendPayload{"createTask", "", incomingTask} // setup appropriate paylooad to send to the backend server

	} else {
		incomingTaskID := form["id"][0]                                      // attempt to retreive the id of the task in the case its an update POST
		message = BackendPayload{"updateTask", incomingTaskID, incomingTask} // setup appropriate paylooad to send to the backend server

	}
	// sending message to connection via the json encoder
	err = encoder.Encode(&message)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Redirect("/", iris.StatusFound) // reload page for the client
	return
}

// deleteTask takes the iris context, which provides the incoming form data, this function deletes the task if it exists and reloads
func deleteTask(ctx iris.Context) {
	connection, err := net.Dial("tcp", *endpoint)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()
	encoder := json.NewEncoder(connection)

	form := ctx.FormValues()
	incomingTaskID := form["id"][0] // attempt to retreive the id of the task in the case its an update POST
	if incomingTaskID == "" {       // if no ID was specifed then undefined error, this does not happen unless client edits the HTML
		ctx.Redirect("/", iris.StatusFound) // reload page for the client
		return
	}
	message := BackendPayload{"deleteTask", incomingTaskID, ""}
	err = encoder.Encode(&message)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Redirect("/", iris.StatusFound) // reload page for the client
	return
}
