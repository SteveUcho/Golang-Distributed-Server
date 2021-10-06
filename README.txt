Steve Ucho

How to build:
      go build

How to run(ubuntu linux):
      ./frontend --listen <PORT> --backend <HOST:PORT>
      ./backend --listen <PORT>

State of assignment:
      All requirements completed.
      The servers can handle multiple requests, data is locked with sync.Map and there is failure detection in place, more on these points under design decisions

Resources used:
      https://www.iris-go.com/docs
      https://blog.golang.org/json
      https://golang.org/pkg/flag
      https://golang.org/pkg/sync/#Map

Design Decision:
      I have decided to use the sync.Map implementation of Map because it is safe for concurrent use by multiple goroutines without locking or coordination.
      There was no explicit or implicit requirement not to use the sync package or its map implementation.
      By doing this, the backend is free to edit and read the data of individual items, the map handles the locking and unlocking of each element in the map.
      Because of the fine granularity locking approach done by the map, its individual operations performace suffers slightly compared to a regular map using a read write mutex over the entire map.
      Beacuse sync.map is not a traditional map, both the key and value loose their type once they are inserted, I do type assertion in order to verify that any data exiting the map is of the right shape
      I did consider only creating one mutex for the entire map but it does not work very well with more than one person and different people editing different items. It does not create a good user experience

      I implemented the failure detection based on whether a connection can be established with the backend from the frontend.
      For any CRUD action performed a new connection is made between the frontend and backend.
      If the connection cannot be established then the frontend prints the error to the terminal.
      The user can then continue to retry the action until it succeds, printing the error each time it fails.

      I used vegeta to test using 50 worker for 30 seconds on the / (root) operation because this is the most work intense path.
      For each call to root, the back end must convert the entire sync.Map to a regular map in order to send the data to the frontend

Additional Thoughts:
      Instructions were clear.
      I was given more than enough support to develop my solution as always.
      This assignment was very interesting and made me think.