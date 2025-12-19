# 5 Building and Consuming RESTful APIs with Go
## Setting up a basic API server with Go
We will use Go’s built-in `net/http` package and `gorilla/mux` for routing.
### Setting up a project
```bash
$ mkdir go-url-shortener
$ cd go-url-shortener
$ go mod init github.com/<yourusername>/go-url-shortener
```
Next, install the `gorilla/mux` package for routing:
```bash
$ go get -u github.com/gorilla/mux
```
### Creating the API server
Now, create a [`main.go`](./0501-go-url-shortener/main.go) file and set up a basic HTTP server.
Run the server:
```bash
$ go run .
# or 
$ go run main.go
```
Open your browser and go to http://localhost:8080. You should see the message **Welcome to the URL shortener API!**.
### Creating the URL shortening logic
To store short URLs, we will use an in-memory map. Create a file named [handlers.go](./0501-go-url-shortener/handlers.go).

With these functions in place, we now have the basic logic for a working URL shortener. The `createShortURLHandler` function accepts a long URL, generates a short version, and stores the mapping. The `getOriginalURLHandler` function looks up the short code and redirects users to the original link. The use of a mutex ensures that multiple users can safely access and update the in-memory store at the same time. Together, these pieces demonstrate how Go can be used to build efficient and reliable web services.


### Adding routes
Now, update [main.go](./0501-go-url-shortener/main.go) to include these routes
```bash
$ go run .
```
Alternatively, we can specify the exact files:
```bash
$ go run main.go handlers.go
```
### Testing the API
Use curl or Postman to shorten a URL:
```bash
$ curl -X POST "http://localhost:8080/shorten" -H "Content-Type: application/json" -d '{"long_url": "https://example.com"}'
```
The following is an example response:
```json
{
  "short_url": "http://localhost:8080/abc123",
  "long_url": "https://example.com"
}
```
Visit http://localhost:8080/abc123 in your browser, and it should redirect you to https://example.com.

## Handling requests and responses
We will use Go’s `net/http` package along with `gorilla/mux` for better request handling.

### Understanding HTTP methods and JSON handling
Before we dive into handling data, it is important to understand how web APIs work with different types of HTTP methods. These methods define what kind of action the client wants the server to perform:

* GET: Retrieve data from the server
* POST: Send data to the server to create a resource
* PUT: Update an existing resource
* DELETE: Remove a resource from the server


Go provides the `encoding/json` package to encode and decode JSON data in API requests and responses.

### Validating incoming data
The `createShortURLHandler` function processes incoming `POST` requests to create a shortened URL. It includes three main steps:
1. Decode the JSON request body.
2. Validate that the long_url field is present.
3. Generate and store the shortened URL (short_url), and send a JSON response back to the client.


Here’s how createShortURLHandler is [implemented in Go](./0502-go-url-shortener/handlers.go)

### Improving error responses
Instead of [using `http.Error`](./0503-go-url-shortener/handlers.go), we can create a helper function to format errors as JSON responses.
Modify the `createShortURLHandler()` function to use [this helper](./0503-go-url-shortener/handlers.go).

### Handling GET requests and redirects
### Enhancing response structure
For better API responses, let’s create a function to format successful responses:
Modify[ `createShortURLHandler()` to use `respondWithJSON()`](./0505-go-url-shortener/handlers.go)

## Managing routing and middleware
For example, in a URL shortener API, middleware can do the following:
* Log incoming requests for debugging
* Check API keys to restrict access
* Handle **Cross-Origin Resource Sharing (CORS)** to allow requests from different domains
* Rate limit requests to prevent abuse
* Standardize error responses

### Logging middleware
```bash
$ go run .
```
Alternatively, we can specify the exact files:
```bash
$ go run main.go handlers.go middleware.go
```
Start the server and make some requests:
```bash
$ curl http://localhost:8080/
```
Check the logs:
```bash
GET / took 200ms
```
### [Adding authentication middleware](./0507-go-url-shortener/auth.go)
Now, let’s apply [middleware](./0507-go-url-shortener/main.go).

```bash
$ go run .
```
Alternatively, we can specify the exact files:
```bash
$ go run main.go handlers.go middleware.go auth.go
```
Make a request without an API key:
```bash
$ curl http://localhost:8080/shorten
```
The following is the response:
```json
{
  "error": "Forbidden"
}
```
Now, make a request with the API key:

```bash
$ curl -H "X-API-Key: key-admin-123" http://localhost:8080/shorten
```
### [Rate-limiting middleware](./0508-go-url-shortener/rate-limit.go)
Rate limiting prevents abuse by restricting the number of requests for a certain amount of time from a user.
Apply Middleware:
```golang
router.Use(RateLimitMiddleware)
```
Make multiple requests quickly:
```bash
for i in {1..6}; do curl http://localhost:8080/shorten; done
```
The sixth request should return the following:
```json
{
  "error": "Too many requests"
}
```

### Combining middleware
Middleware functions can be chained together by applying multiple .Use() calls:
```golang
router := mux.NewRouter()

router.Use(RateLimitMiddleware)
router.Use(LoggingMiddleware)        // Apply logging middleware globally
router.Use(AuthenticationMiddleware) // Apply authentication middleware
```
### Testing middleware in Go
Testing middleware is very important to ensure it behaves as expected. We can use the `net/http/httptest` package to simulate requests.
Run the tests:
```bash
go test -v
```
## Consuming external APIs and integrating data
We’ll use the `httpstatus.io` API, a service that provides URL status codes, to demonstrate how you can make API calls, parse responses, and integrate Prometheus metrics to track URL validation.
```bash
$ curl -X POST http://localhost:8080/shorten \
-H "Content-Type: application/json" \
-d '{"requestUrl": "https://example.com"}'
```
#### Vscode Debuging
* Configure VS Code Debug: press `Ctrl+Shift+D`  to create  `.vscode\launch.json`
  ```json
  {
        // Use IntelliSense to learn about possible attributes.
        // Hover to view descriptions of existing attributes.
        // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
        "version": "0.2.0",
        "configurations": [
        {
            "name": "go-url-shortener",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/ch05/go-url-shortener", 
            "cwd": "${workspaceFolder}",
            "env": {},
            "args": []
        }
     ]
  }
  ```
  * `${workspaceFolder}`: `E:\go_workspaces\packt-mastering-go-devops-automating-scaling-2025`
  * Install Delve (Go Debugger)
    ```bash
    go install github.com/go-delve/delve/cmd/dlv@latest
    ```
* In `VSCode` press `Ctrl+Shift+P` ,and then type `Go: Locate Configured Go Tools`
* Clean old compiling cache
  ```bash
  go clean -cache -modcache -testcache -fuzzcache
  ```
* Clean old pkg
  ```bash
  rmdir /s /q "%GOROOT%\pkg"
  ```   
### Best practice: Don’t hardcode secrets!
### Creating the URL shortener with metrics integration


### Make sure to use a valid or invalid URL to see how the middleware handles it:
```bash
$ curl -X POST http://localhost:8080/shorten \
-H "Content-Type: application/json" \
-d '{"requestUrl": "https://example.com"}'
```
Visit `/metrics` to view Prometheus metrics.
```bash
$  curl -X POST http://localhost:8080/metrics -H "Content-Type: application/json" -d '{}'
```