
## Protocol Buffer Compiler
* https://protobuf.dev/installation/
* https://protobuf.dev/getting-started/
### Installation
#### Install Pre-compiled Binaries (Any OS)
* Download Protocol Buffer Compiler: https://github.com/protocolbuffers/protobuf/releases
  In the downloads section of each release, you can find pre-built binaries in zip packages:` protoc-$VERSION-$PLATFORM.zip`. It contains the protoc binary as well as a set of standard .proto files distributed along with protobuf. 
  ```
  PB_REL="https://github.com/protocolbuffers/protobuf/releases"
  curl -LO $PB_REL/download/v33.2/protoc-33.2-win64.zip
  ```  
* Unzip the file under `$HOME/.local` or a directory of your choice. For example:
  ```
  unzip protoc-33.2-win64.zip -d $HOME/.local
  ```
* Update your environment’s path variable to include the path to the protoc executable. For example:
  ```
  export PATH="$PATH:$HOME/.local/bin"
  ```
#### Install Using a Package Manager
* Windows, using [Winget](https://learn.microsoft.com/en-us/windows/package-manager/winget/)
  ```bash
  > winget install protobuf
  > protoc --version # Ensure compiler version is 3+
  ```
* Windows, using [Chocolatey](https://community.chocolatey.org/packages/protoc)
  ```bash
  choco install -y protoc
  ```
* MacOS, using [Homebrew](https://brew.sh/):  
  ```bash
  brew install protobuf
  protoc --version  # Ensure compiler version is 3+
  ```
* Linux, using apt or apt-get, for example:  
  ```bash
  brew install protobuf
  apt install -y protobuf-compiler
  protoc --version  # Ensure compiler version is 3+
  ```
# ch06 Working with gRPC and Microservices Architecture
## Defining service contracts with Protocol Buffers
Before writing our `.proto` file, let’s set up a clean project layout so the rest of the examples work as expected:
```bash
url-shortener/
├── go.mod
├── proto/
│   └── url_shortener.proto
├── server/
│   └── server.go
└── client/
    └── client.go
```    
Let’s initialize the Go module:
```bash
go mod init url-shortener
```
A Protobuf definition file has a `.proto` extension that describes the following:
* Messages: Data structures exchanged between services
* Services: Functions (RPCs) that clients can call
  
To get started, create the [`url_shortener.proto`](./url-shortener/proto/url_shortener.proto) file in the `proto/` folder.
The package name groups related definitions and prevents name conflicts:
* The `URLShortener` service defines a gRPC service with two RPC methods:
  * `ShortenURL` to shorten a URL
  * `GetOriginalURL` to retrieve the original URL
* The `ShortenRequest` message carries the original URL
* The `ShortenResponse` message returns the shortened URL
* The `GetRequest` message carries the shortened URL to get the original one
* The `GetResponse` message returns the original URL
  
To generate Go code from `.proto` files, we need the **Protocol Buffer Compiler** (`protoc`) installed on the system:
```bash
# Linux
sudo apt update
sudo apt install -y protobuf-compiler

# Windows 
choco install -y protoc
winget install protobuf
```
Install the Go plugins for protoc:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
Verify the installation:
```bash
protoc  --version
protoc-gen-go  --version
protoc-gen-go-grpc  --version
```
Generate the Go code from the .proto files;
```bah
protoc --go_out=. --go-grpc_out=. proto/url_shortener.proto

go mod tidy
```
This will generate two files:
* `url_shortener.pb.go`: Contains Go structs and methods for the Protobuf messages
* `url_shortener_grpc.pb.go`: Contains gRPC service definitions

We can now use the generated code to create a gRPC server and client:
```golang
import (
    "context"
    pb "url-shortener/proto"
)
```
### Accessing generated types and methods
Once the `.proto` files are compiled, Go generates code that includes all the necessary types and service methods. These generated types make it easy to work with structured request and response data in the gRPC application:
* `pb.ShortenRequest`: Request struct for shortening a URL
* `pb.ShortenResponse`: Response struct for shortened URL
* `pb.GetRequest`: Request struct for retrieving a URL
* `pb.GetResponse`: Response struct for original URL
  
With the code now generated, we are ready to build the gRPC server and client. These generated types and methods will help us send and receive structured data between services. In the next steps, we will use them to handle real gRPC requests and responses.

### Creating a gRPC server
To build a working gRPC service, we need to create a server that implements the methods defined in the `.proto` file. [The following code sets up `URLShortenerServer`, which can shorten URLs and retrieve the original URLs using an in-memory map. It also starts a gRPC server that listens for requests on port `50051`.](./0601-url-shortener/server/server.go)

First, import the required Go libraries, including standard packages for networking and logging, and the generated Protobuf code (`pb`):

```bash
go run server/server.go
```
With this setup, we now have a basic but functional gRPC server that can shorten URLs and retrieve the original ones using in-memory storage. While this example keeps things simple to focus on the core ideas, we can continue building on it by adding persistence, validation, or a frontend interface to turn it into a more complete application.

### Creating a gRPC client
[The following code shows how to build a gRPC client that connects to the URL shortener server. It first establishes a connection, then sends a request to shorten a URL, and prints the shortened version. After that, it sends another request to retrieve the original URL using the shortened one and prints the result.](./0601-url-shortener/client/client.go)

This example shows how a basic gRPC client works and helps to understand how to send requests, receive responses, and interact with a remote service.
```bash
go run client/client.go
```
Now that we’ve seen both the server and the client in action, we should have a good foundation to build more advanced features or integrate gRPC into other projects.

### Best practices for defining Protobuf files
Here are some best practices to follow when defining `.proto` files:
 * **Use meaningful names**: Choose clear and descriptive names for services and messages
 * **Add comments**: Add comments in the `.proto` file to document the functionality
 * **Group related services**: If you have multiple services, group them logically under a common package

In this section, we learned how to define service contracts using Protobuf in a gRPC-based URL shortener. We defined services and messages in a `.proto` file, generated Go code from those definitions, and used it to implement a gRPC server and client. Clear, well-structured Protobuf contracts ensure efficient and reliable communication between microservices.

## Implementing streaming with gRPC
In this section, we will explore how to implement gRPC streaming in a URL shortener application. Streaming allows continuous data exchange between the client and server, enabling real-time processing and high-performance communication.

gRPC supports three types of streaming:
* **Server streaming**: The server sends multiple responses to a single client request
* **Client streaming**: The client sends multiple requests and receives a single response
* **Bidirectional streaming**: Both the client and the server send a stream of messages to each other
 
In this section, we will do the following:
* Implement server-side streaming to send a list of shortened URLs
* Implement client-side streaming to batch shorten multiple URLs
* Implement bidirectional streaming to send and receive real-time URL statistics

### Overview of gRPC streaming
Streaming is useful when you need to send or receive a continuous flow of data, instead of waiting for everything at once. This can be great for real-time updates, processing batches of data, or transferring large files more efficiently.

Here’s how it works, step by step:
1. **Connection established**: The client and server open a connection and keep it open for the stream.
2. **Data sent in chunks**: Instead of sending all the data at once, they send pieces (chunks) as needed.
3. **Connection closed**: Once everything is sent and received, the connection is closed.

You might use streaming if you need live updates, for example, tracking how many times a short URL is clicked. It also helps when sending large or multiple items at once without overwhelming the server or client.

###  Defining gRPC streaming methods
We will add streaming methods to our [`url_shortener.proto`](./0602-url-shortener/proto/url_shortener.proto) Protobuf file. 

Each method in the service has a specific purpose and uses a different type of gRPC communication. Before we implement them, let’s take a quick look at what each method does. This will help to understand how and when to use them:
* `ShortenURL`: A simple request-response to shorten a single URL
* `GetOriginalURL`: Fetches the original URL from a shortened one
* `ListShortenedURLs`: Server-side streaming to return a list of all shortened URLs
* `BatchShortenURLs`: Client-side streaming to shorten multiple URLs in one request
* `MonitorURLStats`: Bidirectional streaming to monitor real-time stats for each URL

Run the following command to regenerate the Go code:
```bash
protoc --go_out=. --go-grpc_out=. proto/url_shortener.proto
```
This will update the `url_shortener.pb.go` and `url_shortener_grpc.pb.go` files.

With these updates, the `.proto` file now supports different types of gRPC streaming methods. You are now ready to build more powerful and flexible communication patterns between your services. In the next steps, we will implement these methods on the server and client sides.

### Implementing server-side streaming
To begin, open the [server/server.go](./0602-url-shortener/server/server.go) file and update it to include the server-side streaming method.

### Implementing the client for server streaming
To use the server-side stream, update the [client/client.go](./0602-url-shortener/client/client.go) file to receive and process the streamed responses from the server.

### Implementing client-side streaming
Let’s update the [server/server.go](./0602-url-shortener/server/server.go) file to handle client-side streaming by implementing the `BatchShortenURLs` method.

### Implementing the client for client streaming
Update [client/client.go](./0602-url-shortener/client/client.go)  to send a stream of URLs.

### Implementing bidirectional streaming
Now, let’s update the [server/server.go](./0602-url-shortener/server/server.go) file to support bidirectional streaming by implementing the `MonitorURLStats` method.

### Implementing the client for bidirectional streaming
Let’s now implement the client-side logic to support bidirectional streaming. In this code, the client will send a list of shortened URLs to the server and receive real-time statistics for each of them. This allows the client to send and receive data at the same time over a single connection.

This implementation shows how bidirectional streaming can help build interactive and responsive applications. As you can see, the client sends multiple requests using a goroutine, while another goroutine listens for server responses. This concurrent setup enables the client to receive real-time updates without waiting for the entire input to be sent. It is especially useful in scenarios where both the client and server need to continuously exchange information, such as live analytics or chat systems.

### Running the gRPC server and client
To run the gRPC server and client successfully, follow these detailed steps.

#### Running the gRPC server
The gRPC server is responsible for handling incoming client requests, processing them, and sending responses or streaming data back to the client:
```bash
go run server/server.go
```
When the server starts, it should display a message similar to the following:
```bash
Starting gRPC server on port :50051
```
This means that the server is listening for incoming gRPC requests on port `50051`, and it is ready to process the methods defined in the `url_shortener.proto` file.

##### Running the gRPC client
The client is responsible for sending requests to the server and handling the responses:
```bash
go run client/client.go
```
##### What to expect for each streaming type
Let’s see what happens when you run the client and use each of the gRPC streaming methods. These examples will help you understand how data flows between the client and server in different streaming scenarios.

###### Server-side streaming: ListShortenedURLs
In server-side streaming, the client sends a single request, and the server responds by sending back a stream of results, one at a time. This is useful when the server needs to return a list of items:
* The client sends a request to list all shortened URLs
* The server streams each shortened URL back to the client
You should see output similar to the following:
```bash
Shortened URL: short.ly/1
Shortened URL: short.ly/2
Shortened URL: short.ly/3
```
This output confirms that the server is streaming the list of shortened URLs to the client one by one.

###### Client-side streaming: BatchShortenURLs
In client-side streaming, the client sends a stream of data to the server, and when it’s done, the server sends back one response. This is helpful when sending a batch of inputs to be processed together:
* The client sends multiple URLs to be shortened as a stream
* Once all URLs are sent, the server processes them and returns the number of URLs shortened
  
Expected output is similar to the following:
```bash
Batch shortening completed. 3 URLs shortened.
```

This confirms that the client successfully sent a stream of URLs and received a summary from the server.

###### Bidirectional Streaming: MonitorURLStats
In bidirectional streaming, both the client and server send data streams at the same time. This is great for real-time use cases such as live monitoring or chat apps:

* The client sends URLs to monitor for statistics
* The server continuously streams real-time statistics for each URL back to the client
* 
Expected output is similar to the following:
```bash
Stats for short.ly/1: 100 clicks
Stats for short.ly/2: 120 clicks
```

This confirms that both the client and server are sending and receiving streams of data simultaneously.
## Best practices for gRPC in microservices
When working with gRPC in a microservices architecture, it’s essential to follow best practices to ensure high performance, security, and maintainability.

### Efficient connection management and load balancing
In a microservices environment, multiple gRPC clients may connect to the same server. To avoid opening and closing connections repeatedly, it’s best to use connection pooling and implement load balancing when working with multiple instances.

#### Connection pooling
[Instead of creating a new connection for every request, reuse a single connection for multiple requests.](./0603-url-shortener/client/client.go)

#### Implementing load balancing
Use load balancing if multiple instances of the gRPC server are running. You can use a load balancer such as Envoy or a gRPC load balancer.

Here is an example of a client with load balancing:
```golang
grpc.Dial("dns:///my-grpc-service.default.svc.cluster.local:50051", grpc.WithInsecure())
```
Using load balancing ensures that traffic is spread across multiple server instances, which helps improve reliability, performance, and scalability. It allows the gRPC client to automatically distribute requests without needing to manage server addresses manually. This setup is useful in production environments where services are running in Kubernetes or other distributed systems.
### Proper error handling and retries
Errors can occur due to network failures, server timeouts, or invalid inputs. Implementing error handling ensures that the application can gracefully handle these scenarios.
#### Basic error handling in the client
When building a gRPC client, it is important to handle errors properly to make the application more reliable and easier to debug. [This example shows how the client can safely call the server to fetch the original URL from a shortened one, while also checking for any issues that might happen during the request.](./0604-url-shortener/client/client.go)
#### Implementing retries with backoff
When connecting to a gRPC server, network issues or server unavailability can sometimes cause temporary failures. [Implementing retries with backoff can help the client automatically recover from these failures by attempting to reconnect after short delays, improving the reliability of the application.](./0604-url-shortener/client/client.go)

### Enabling security with TLS
gRPC supports secure communication using Transport Layer Security (TLS), ensuring that the data exchanged between clients and servers is encrypted.

#### Enabling TLS on the server
To protect data exchanged between the client and server, it is important to enable TLS. TLS encrypts the communication channel and ensures that the server is authenticated, which helps prevent eavesdropping and tampering. [This section shows how to configure a gRPC server to use TLS by loading a certificate and private key and using them when starting the server.](./0605-url-shortener/server/server.go)

Now, the gRPC server is configured to use TLS, which means all communication with clients will be encrypted and secure. This setup ensures that sensitive information, such as URLs and statistics, is protected during transmission.

Enabling TLS is a critical step toward building reliable and production-ready services that users can trust.

#### Configuring TLS on the client
To securely connect a gRPC client to a server that uses TLS, the client must be configured to trust the server’s certificate. This is done by loading the **Certificate Authority (CA)** certificate and using it to establish a secure connection. The CA certificate helps the client verify that it is communicating with the right server and that the connection is encrypted:
```golang
func createSecureClient() pb.URLShortenerClient {
    creds, err := credentials.NewClientTLSFromFile("ca.crt", "")
    if err != nil {
        log.Fatalf("Failed to load CA cert: %v", err)
    }
    conn, err := grpc.Dial(
        "localhost:50051",
        grpc.WithTransportCredentials(creds)
    )
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    return pb.NewURLShortenerClient(conn)
}
```
Once the secure connection is established, the client can safely send and receive data over the network. This setup improves trust and security in communication and is recommended for production environments.
### Using deadlines and timeouts
Deadlines and timeouts prevent long-running requests from consuming system resources indefinitely.

#### Versioning gRPC services
As gRPC services evolve, you need to introduce changes without breaking existing clients. Versioning allows smooth transitions between different API versions:
```proto
syntax = "proto3";
package urlshortener.v1;
service URLShortenerV1 {
    rpc ShortenURL(URLRequest) returns (URLResponse);
}
package urlshortener.v2;
service URLShortenerV2 {
    rpc ShortenURL(URLRequestV2) returns (URLResponseV2);
}
```
This setup allows supporting both old and new clients at the same time. It gives the flexibility to roll out updates gradually, while keeping the service stable and reliable for users already using earlier versions.

### Monitoring and logging gRPC requests
Monitoring and logging help track performance issues and debug errors in production environments. One way to achieve this is by using a unary interceptor, which lets us add cross-cutting logic around every gRPC request. We can define a logging interceptor like this:
```golang
import (
    "context"
    "log"
    "net"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    pb "url-shortener/proto"
)
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    log.Printf("Received request for method: %s", info.FullMethod)
    resp, err := handler(ctx, req)
    if err != nil {
        log.Printf("Error handling request: %v", err)
    }
    return resp, err
}
```
Then, we can attach it to the gRPC server during setup:
```golang
func main() {
    certFile := "server.crt"
    keyFile := "server.key"
    creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
    if err != nil {
        log.Fatalf("Failed to load TLS keys: %v", err)
    }
    server := grpc.NewServer(
        grpc.Creds(creds),
        grpc.UnaryInterceptor(loggingInterceptor),
    )
    pb.RegisterURLShortenerServer(server, &URLShortenerServer{
        urlMap: make(map[string]string),
    })
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    log.Println(
        "gRPC server with TLS and logging interceptor listening on
        port 50051"
    )
    server.Serve(lis)
}
```
If Prometheus is used for observability, we can integrate metrics using the `go-grpc-middleware` ecosystem. It provides out-of-the-box middleware for logging, monitoring, authentication, and more. For example, we can add a Prometheus interceptor like this:
```golang
import (
    grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)
server := grpc.NewServer(
    grpc.Creds(creds),
    grpc.ChainUnaryInterceptor(
        loggingInterceptor,
        grpc_prometheus.UnaryServerInterceptor,
    ),
)
```
Adding logging like this gives better visibility into how the gRPC server is being used. It helps to catch problems early and understand what’s happening during each request, making it easier to maintain and improve the service.