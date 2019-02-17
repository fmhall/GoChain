# GoChain
This repo contains two implementations of a simple blockchain in Go

The core logic for both the networking and HTML implementations is the same: simple block generation, validation, and adding blocks to the canonical chain.
In the HTML folder, you will find the main logic for displaying the chain in its current form on an HTML server using Mux. You can send POST and GET requests to the server using Postman or curl, and a block with that data will be added to the chain.

In the Networking folder, the HTML interface has been removed, and all information is displayed in terminals. The networking implementation uses TCP servers to simulate the nodes that broadcast and recieve blocks.
Each TCP server (new terminal) can read data, create a block, validate that block, and send it to the Go channel. Go routines are used to handle this process concurrently, and mutexes are used to prevent race conditions.
The conn interface from the net package allows elegant connection handling, and the channel enables the various concurrent routines to reference a single canonical chain.

## Deployment:

```
git clone https://github.com/fmhall/GoChain.git
```
### For the HTML version:

* Navigate to the HTML directory and rename the example file `mv example.env .env`
```
go run main.go
```
* Open a web browser and visit http://localhost:8080/
* To write new blocks, send a POST request using Postman or curl to http://localhost:8080/ with a JSON payload with data as the key and an integer as the value. For example:
```{"data":1001}```
* Send requests, and refresh the browser to watch the chain grow over time

### For the TCP/Networking version:

* Navigate to the Networking directory and rename the example file `mv example.env .env`
```
go run main.go
```
* Open a few terminals and type `nc localhost 9000` to initialize the TCP server (our 'node') that talks to the port specified in the .env file.
* To write new blocks, input an integer data value in any of the terminals. The node will validate the data, create the block and relevant hashes, validate the block, and send it to the channel. Watch the blockchain grow as it is updated and sent to each terminal.




