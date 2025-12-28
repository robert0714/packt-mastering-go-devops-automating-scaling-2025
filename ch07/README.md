# 7 Using Go to Build Custom Terraform Providers
## Setting up the project structure
* https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider
* https://github.com/hashicorp/terraform-provider-hashicups/blob/main/README.md
* SEP 26, 2014, Official documentation can be found at https://www.hashicorp.com/en/blog/writing-custom-terraform-providers ().
   * https://github.com/hashicorp/terraform/tree/v0.9.11/builtin/providers
* Terraform plugin protocol: https://developer.hashicorp.com/terraform/plugin/terraform-plugin-protocol
  
Here’s the basic file structure of the project:
```bash
/terraform-provider-urlshortener
├── /internal
    ├── provider.go
│   └── resource_shorturl.go
├── /main.go
└── /go.mod
```

If the provider supports multiple resources or data sources, it’s a good practice to place each in its own file (for example, `redirect.go`, `analytics.go`, and so on) within the `/internal` folder. This keeps the code modular and maintainable. The `provider.go` file typically registers all available resources and data sources with Terraform, serving as the entry point for the provider logic.

Now, let’s initialize the project:
```bash
# Create the project folder:
mkdir terraform-provider-urlshortener
cd terraform-provider-urlshortener
# Initialize the Go module:
go mod init github.com/example/terraform-provider-urlshortener
```
Once the project structure is in place and the Go module is initialized, we can begin writing the provider logic. Instead of installing dependencies manually (`go get <package_name>`), it’s recommended to import required packages directly in the `.go` files (for example, `provider.go`). Go modules will automatically track those imports.

After importing the required packages in Go files, run the following:
```bash
go mod tidy
```
This command automatically adds any missing dependencies to `go.mod` and downloads them. It also cleans up any unused dependencies, helping to keep the project tidy and up to date.

A typical workflow is to import what is needed in the `.go` files, then run `go mod tidy` to ensure that the module has exactly the dependencies it requires, nothing more, nothing less.

### Why this design
Provider-level configuration (API key) is defined in the provider schema and mapped to a typed Go model (`types.String`) with `tfsdk` tags.

We configure building the client and store it in `resp.ResourceData` or `resp.DataSourceData` so resources can retrieve it during configuration and CRUD operations.

Resources implement CRUD and use `req.Plan`, `req.State`, and `resp.State.Set` to read plan/state and write state, as the framework expects. Use `resp.State.RemoveResource()` from `Read` if the remote object no longer exists`; Delete()` need not call `RemoveResource()`, a successful `Delete()` operation implicitly removes the state.

## Defining the provider
The provider serves as the entry point for Terraform to interact with the custom resources. Providers define the structure and configuration logic that Terraform will use when initializing the provider.

First, create a new file called `provider.go` in the `./provider/` directory. In that file, define a type called `urlShortenerProvider` with a `string` field to hold the API key. Then, add a Configure() method to this type. This method will be called by Terraform when ivru&nting \`tfsdk:”api_key”\`

The `urlShortenerProvider` type holds the provider’s configuration, here, just an API key as a string:

The `Configure()` method reads the API key from the Terraform configuration when Terraform starts. It extracts the value and saves it into the `APIKey` field so the provider can use it later to authenticate requests. This method also reports any errors back to Terraform to help troubleshoot configuration issues. This setup is essential to make sure the provider has what it needs before it manages any resources.

Next, create the `main.go` file in the root directory to serve as the program’s entry point. This file uses the recommended pattern for launching Terraform providers with the latest `terraform-plugin-framework` package.
```bash
go run . 
```

The `main()` function is the entry point of the provider. It uses `providerserver.Serve()` from the Terraform plugin framework to start the provider. This function takes a context, a factory function (`provider.New(version)`) that returns an instance of the provider, and optional server configuration via `ServeOpts`. In this case, we optionally enable debug mode based on a command-line flag.
```bash
go run . -debug
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
            "name": "terraform-provider-urlshortener",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/ch07/07-01-terraform-provider-urlshortener", 
            "cwd": "${workspaceFolder}",
            "env": {},
            "args": ["-debug"]
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
This setup ensures that Terraform can correctly discover, initialize, and interact with the provider when users execute Terraform CLI commands.

## Implementing CRUD operations for resources

### Defining the resource schema

### Testing CRUD operations
Test steps
1. Compiling custom Terraform provider plugin:
   ```bash
   go build -o terraform-provider-url_shortener.exe
   ```
2. Configure main.tf
   ```tf
   terraform {
      required_providers {
        url_shortener = {
          source  = "dev/urlshortener"
          version = "0.1.0"
        }
      }
    }

    provider "url_shortener" {
      # add api_key ...etc
    }
   ```
3. Configure local provider path: `terraform.rc` or `.terraformrc` 
   * Windows
     * default path：`%APPDATA%\terraform.rc`
       example： `C:\Users\<USER_NAME>\AppData\Roaming\terraform.rc`
   * macOS / Linux
     * default path： `~/.terraformrc`
   * content:
     ```hcl
      provider_installation {
        filesystem_mirror {
          path = "terraform.d/plugins"
        }
        direct {}
      }
     ```     
     In folder `tf-test` , create directory  `terraform.d/plugins/registry.terraform.io/local/urlshortener/0.0.1/windows_amd64`，並將編譯好的 provider binary 複製到此處。範例檔名：
     - `terraform-provider-urlshortener_v0.0.1.exe`（Windows）
     - `terraform-provider-urlshortener_v0.0.1`（Linux）
      ```bash
      terraform.d/
      └── plugins/
          └── registry.terraform.io/
              └── local/
                  └── urlshortener/
                      └── 0.0.1/
                          └── windows_amd64/
                              └── terraform-provider-urlshortener_v0.0.1.exe
      ```
4. cmd:
    ```bash
    TF_LOG=TRACE terraform providers
    TF_LOG=TRACE terraform init
    TF_LOG=TRACE terraform init -reconfigure
    TF_LOG=TRACE terraform init -upgrade
	# .terraform.lock.hcl is generated by terraform init
    TF_LOG=TRACE terraform apply -auto-approve
    ```
    or
    ```poweshell
    $env:TF_LOG = "TRACE";$env:TF_CLI_CONFIG_FILE = "$PWD\cli.hcl"; terraform providers
    $env:TF_LOG = "TRACE";$env:TF_CLI_CONFIG_FILE = "$PWD\cli.hcl"; terraform init -input=false
	# .terraform.lock.hcl is generated by terraform init
    $env:TF_LOG = "TRACE";$env:TF_CLI_CONFIG_FILE = "$PWD\cli.hcl"; terraform apply -auto-approve
    ```  

## Managing resource state and handling the lifecycle   