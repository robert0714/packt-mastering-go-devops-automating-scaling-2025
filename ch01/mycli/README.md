# Command parsing and argument handling
```bash
go mod init mycli

#Command parsing and argument handling
go run step01/main.go hello Anita

go run step02/main.go -name  robert -age 12


# Setting up cobra
go get -u github.com/spf13/cobra

# Creating the root command
go run step03/main.go

# Adding subcommands
go run step04/main.go hello
go run step04/main.go hello Ada

# we’re enhancing the hello subcommand by adding a custom flag. 
go run step05/main.go hello Mert --greeting Hi
```
> [!NOTE]
> `cobra` supports both persistent and local flags:
> * **Persistent flags**: Available to the root command and all subcommands
> * **Local flags**: Specific to a single command

## Best practices
When building CLI applications with `cobra`, it’s important to follow a few best practices to make the tool easier to use, more reliable, and easier to maintain:
* Organize commands clearly: Use subcommands to divide functionality, making it easier for users to explore features
* Documentation: Provide short and long descriptions for each command to make `cobra help` output informative
* Error handling: Use the `cmd.MarkFlagRequired()` function to enforce necessary flags and improve error messages

# Output formatting and user feedback

```bash
go run step06/main.go 

# Go’s ecosystem offers libraries to create even more dynamic and user-friendly output. 
go get -u github.com/olekukonko/tablewriter

go run step07/main.go
```
To check version , tablewriter v1.1.x 's API had been changed.
```bash
# check  
go list -m github.com/olekukonko/tablewriter
go doc github.com/olekukonko/tablewriter

```
Another key element of output formatting is color. Colors can help users quickly identify different types of information, such as warnings, errors, or success messages. 

```bash
go get -u github.com/fatih/color

go run step08/main.go
```
 A CLI application without progress indicators may leave users uncertain of whether it’s actively working or stalled. One effective solution is the progressbar package, which lets developers implement customizable progress bars.

The following example shows how to use the progressbar package to display progress during a task:
```bash
go get -u github.com/schollz/progressbar/v3

go run step08/main.go
```
# Networking best practices
Here’s an example that demonstrates concurrent HTTP requests using goroutines:
```bash
go run step10/main.go
```
## Error handling and retries
Go’s standard library doesn’t include a built-in retry mechanism, but retries can be implemented using a simple loop or a helper function. Here’s an example of retry logic that makes multiple attempts to fetch a URL with a delay between retries:
```bash
go run step11/main.go
```
## Securing connections
Go’s `http` package automatically enforces HTTPS connections when URLs start with `https://`, but additional measures may be needed for certificate validation or when working with self-signed certificates.
 ```bash
go run step12/main.go
```
Setting `InsecureSkipVerify` to `true` disables certificate verification, allowing connections to proceed. However, this approach introduces security risks and should only be used in controlled environments where certificate validation is impractical, such as development and testing.

## Handling timeouts
Go’s `http.Client` allows setting timeouts for both connection establishment and response retrieval, providing fine-grained control over network operations.
```bash
go run step13/main.go
```

## Rate limiting
```bash
go run step14/main.go
```
# Best practices for building robust CLIs
```bash
go run step15/main.go hello Mert --greeting Hi
```