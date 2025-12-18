# Building executable binaries for different platforms
```bash
mkdir mycli
cd mycli
go mod init mycli
go build ./cmd/mycli
```
## Cross-compilation basics
To cross-compile a Go application, you need to set the `GOOS` and `GOARCH` environment variables to specify the target operating system and architecture.
For example, here’s how to compile a binary for Windows on a Linux machine:
```bash
$ GOOS=windows GOARCH=amd64 go build -o mycli.exe ./cmd/mycli
```
Similarly, here’s how to build for macOS:
```bash
$ GOOS=darwin GOARCH=amd64 go build -o mycli ./cmd/mycli
```
And the following code is for Linux:
```bash
$ GOOS=linux GOARCH=amd64 go build -o  ./cmd/mycli
````
Common GOOS values:
* darwin (macOS)
* linux (Linux)
* windows (Windows)

Common GOARCH values:
* amd64 (64-bit x86)
* 386 (32-bit x86)
* arm (32-bit ARM)
* arm64 (64-bit ARM)

## Handling platform-specific differences
### File paths
Use the `filepath` package to manage file paths in a platform-agnostic way. The `filepath.Join()` function automatically uses the correct path separator (`/` on Unix-like systems and `\` on Windows).
```golang
import "path/filepath"
func getConfigPath() string {
    return filepath.Join("config", "settings.json")
}
```
### Environment variables
The way environment variables are accessed and used can vary. Use the `os` package to interact with environment variables in a consistent manner:
```golang
import "os"
func getEnvVar(key string) string {
    return os.Getenv(key)
}
```
### Line endings
Text files might have different line endings (`\n` for Unix-like systems and `\r\n` for Windows). Use Go’s `bufio` package to handle line endings when reading or writing text files, ensuring compatibility across platforms.

* Example - reading files
```golang
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text() // 不含 \n 或 \r\n
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
```
* Example - writing files
```golang
package main

import (
	"bufio"
	"os"
)

func main() {
	file, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	lines := []string{
		"First line",
		"Second line",
		"Third line",
	}

	for _, line := range lines {
		writer.WriteString(line + "\n")
	}

	writer.Flush()
}
```
| Scenario(情境) | Sugesttion(建議)      |
| -------------- | --------------------- |
| 一行一行讀設定檔 / log | `bufio.Scanner` ⭐     |
| 大檔案 / 自訂分隔符    | `bufio.Reader`        |
| 寫文字檔           | `bufio.Writer + "\n"` |
| 跨 OS / Docker  | **永遠用 `\n`**          |

## Building tags for platform-specific code
### File-level build tags
We can create separate files for different platforms using naming conventions (`_windows.go`, `_linux.go`, and so on) or use a build tag comment at the top of a file:
* Here’s the Windows platform-specific code:
```golang
//go:build windows
// +build windows
package main
import "fmt"
func sayHello() {
    fmt.Println("Hello from Windows")
}
```
* This is the Linux platform-specific code:
```golang
//go:build linux
// +build linux
package main
import "fmt"
func sayHello() {
    fmt.Println("Hello from Linux")
}
```
The compiler will include the appropriate file depending on the target platform when you run go build.

### Inline build tags
Alternatively, we can use `//go:build` within a single file to conditionally compile code blocks based on the environment.

## External dependencies and tooling
Here are a few strategies to manage such dependencies:
* **Use containers for consistency**: If the CLI depends on tools such as Redis, Postgres, or non-Go CLIs, we can package these in Docker containers and manage their lifecycle (for example, spin them up during development or integration tests). This provides a uniform setup across platforms.
* **Fallback to prebuilt binaries**: If a dependency offers binaries for multiple platforms (for example, `jq` or `ffmpeg`), include logic in the CLI to detect the OS/architecture and download the correct binary if it’s missing.
* **Add dependency checks**: Proactively check whether external tools are installed using functions such as `exec.LookPath()` and display informative messages to the user or trigger automated installation when safe.
```golang
import (
    "fmt"
    "os/exec"
)
func checkDependency(tool string) error {
    _, err := exec.LookPath(tool)
    if err != nil {
        return fmt.Errorf("%s not found in PATH", tool)
    }
    return nil
}
```
Managing external dependencies carefully ensures the CLI tool behaves predictably and is easier to set up on any system.

## Optimizing binaries for distribution
* **Stripping debug information**: Use the -ldflags flag to remove debugging information, which reduces the size of the binary:
    ```bash
    $ go build -ldflags="-s -w" -o mycli
    ```
  * `-s`: Omit the symbol table, which includes information such as function names, variable names, and line number mappings. These symbols are useful for debugging but unnecessary in production builds.
  * `-w`: Omit the DWARF symbol table, which is used for debuggers and profilers to provide stack traces, variable inspection, and other runtime metadata.
* **UPX compression**: Another way to reduce binary size is to compress the executable using `upx` (Ultimate Packer for eXecutables). Install `upx` and then compress your binary:
  ```bash 
  $ upx mycli
  ```
  While `upx` can significantly reduce the size of binaries, it may have a slight impact on startup time due to decompression.

Another useful flag to consider is `-trimpath`, which removes file system paths from the compiled binary. This is particularly helpful for improving build reproducibility and protecting sensitive or system-specific paths in distributed binaries:
```bash 
$ go build -trimpath -ldflags="-s -w" -o mycli  ./cmd/mycli
```
Here, `-trimpath` removes all file system paths from the compiled binary, replacing them with module or import paths.

Using `-trimpath` in combination with `-s` and `-w` is a common practice in CI pipelines or public distributions, as it helps create cleaner, more portable builds.

> [!NOTE]
> While Go offers other advanced flags for tuning inlining or garbage collection, they are generally recommended only in highly specialized performance scenarios and are not typically necessary for standard CLI tooling.

## Automating the build process
### Using shell scripts
Create a shell script to automate the build process for multiple platforms:
```bash
#!/bin/bash
platforms=("windows/amd64" "darwin/amd64" "linux/amd64")
output="mycli"
for platform in "${platforms[@]}"
do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    output_name=$output'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name
done
```
### Using Makefiles
Makefiles offer a structured approach to defining build rules:
```
BINARY_NAME=mycli
all: windows macos linux
windows:
    GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe
macos:
    GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)
linux:
    GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)
```
Run the make command to build for all platforms:
```bash
$ make
```
### Using Goreleaser
```bash
 choco install -y goreleaser
```
`goreleaser` simplifies cross-platform builds, packaging, and releasing binaries. Install `goreleaser` and create a configuration file (`.goreleaser.yml`):
```yaml
version: 2
builds:
  - binary: mycli
    goos:
      - windows
      - darwin
      - linux
    goarch:
      - amd64
```
By default, this configuration requires environment variables (such as `GITHUB_TOKEN`) to publish a release. To disable this publishing behavior, add the following section to the `.goreleaser.yml` file:
```yaml
release:
  disable: true
```
Run goreleaser to automate the entire release process:
```bash
$ goreleaser release
```
Use the `--snapshot` flag to build locally without triggering the full release process:
```bash
$ goreleaser release --snapshot
```
This creates versioned and platform-specific binaries in the `dist/` directory without needing authentication tokens or remote release permissions.

By automating the build process, we can save time, ensure consistent results across different platforms, and focus more on improving the application instead of managing repetitive build tasks.

# Dockerizing CLI applications
## Creating a Dockerfile for a Go CLI application
### Step 1: Start with a lightweight base image that includes Go
For most CLI applications, the `golang:alpine` image is a popular choice because it is small and efficient.
```Dockerfile
# Use the official Golang image as a base
FROM golang:alpine
```
### Step 2: Set the working directory inside the container where the application’s source code will be copied and compiled.
```Dockerfile
# Set the working directory
WORKDIR /app
```
### Step 3: Copy the application’s source code from the host machine to the container.
```Dockerfile
# Copy the source code to the working directory
COPY . .
```
### Step 4: Run the go build command to compile the application inside the container.
```Dockerfile
# Build the Go application
RUN go build -o mycli
```
### Step 5: Specify the command that should run when the container starts. In this case, it’s the compiled CLI application.
```Dockerfile
# Define the entry point for the container
ENTRYPOINT ["./mycli"]
```
Here’s a full Dockerfile example:
```Dockerfile
FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o mycli
ENTRYPOINT ["./mycli"]
```
## Optimizing the Docker image
### Multi-stage builds
Multi-stage builds allow you to use multiple FROM statements in a Dockerfile, where each stage can have its own base image. This approach helps in building the application in one stage and copying only the necessary files to a minimal base image in the final stage.
```Dockerfile
# First stage: Build the application
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o mycli
# Second stage: Create a smaller image for the final executable
FROM alpine
WORKDIR /root/
COPY --from=builder /app/mycli .
ENTRYPOINT ["./mycli"]
```
### Using a scratch image
For minimal Go applications, you can use the scratch base image, which is an empty image with nothing pre-installed. This drastically reduces the image size.
```Dockerfile
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o mycli
FROM scratch
COPY --from=builder /app/mycli .
ENTRYPOINT ["./mycli"]
```
### Removing unnecessary files
Ensure that only the necessary files and binaries are copied into the final image. Use `.dockerignore` to exclude files that are not needed, similar to `.gitignore` in Git.
Sample `.dockerignore` file:
```.dockerignore
# Exclude the vendor directory
/vendor
# Ignore local binaries
/bin
```
By applying these techniques, you can significantly reduce the image size and improve performance.

## Using Docker Compose for local development
```bash
mkdir cli-container
go mod init cli-container
go mod tidy
go get -u github.com/lib/pq
```
## Managing environment variables and configuration

