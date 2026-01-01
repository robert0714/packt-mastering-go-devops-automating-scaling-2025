# 13 Serverless Computing Using AWS Lambda
* [Tutorial: Deploy a Hello World application with AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-getting-started-hello-world.html)
  * https://medium.com/@amirhosseinsoltani7/serverless-api-using-aws-lambda-api-gateway-and-golang-9b68f6012eb6
  * https://dave.dev/blog/2022/09/01-09-2022-aws-lambda-go-local/
    * https://github.com/davedotdev/aws-local-go-lambda-example
    * https://www.youtube.com/watch?v=1uVmHuMo4zU
  * https://github.com/aws/aws-sam-cli-app-templates
    * https://github.com/aws/aws-sam-cli-app-templates/tree/master/al2023/go/hello-img/%7B%7Bcookiecutter.project_name%7D%7D
* steps
```bash
sam init

You can preselect a particular runtime or package type when using the `sam init` experience.
Call `sam init --help` to learn more.

Which template source would you like to use?
        1 - AWS Quick Start Templates
        2 - Custom Template Location
Choice: 1

Choose an AWS Quick Start application template
        1 - Hello World Example
        2 - Data processing
        3 - Hello World Example with Powertools for AWS Lambda
        4 - Multi-step workflow
        5 - Scheduled task
        6 - Standalone function
        7 - Serverless API
        8 - Infrastructure event management
        9 - Lambda Response Streaming
        10 - GraphQLApi Hello World Example
        11 - Full Stack
        12 - Lambda EFS example
        13 - Serverless Connector Hello World Example
        14 - Multi-step workflow with Connectors
        15 - DynamoDB Example
        16 - Machine Learning
Template: 1

Use the most popular runtime and package type? (python3.13 and zip) [y/N]: n

Which runtime would you like to use?
        1 - dotnet8
        2 - dotnet6
        3 - go (provided.al2)
        4 - go (provided.al2023)
        5 - graalvm.java11 (provided.al2)
        6 - graalvm.java17 (provided.al2)
        7 - java21
        8 - java17
        9 - java11
        10 - java8.al2
        11 - nodejs22.x
        12 - nodejs20.x
        13 - python3.9
        14 - python3.13
        15 - python3.12
        16 - python3.11
        17 - python3.10
        18 - ruby3.4
        19 - ruby3.3
        20 - ruby3.2
        21 - rust (provided.al2)
        22 - rust (provided.al2023)
Runtime: 4

What package type would you like to use?
        1 - Zip
        2 - Image
Package type: 1

Based on your selections, the only dependency manager available is mod.
We will proceed copying the template using mod.

Would you like to enable X-Ray tracing on the function(s) in your application?  [y/N]: N

Would you like to enable monitoring using CloudWatch Application Insights?
For more info, please view https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/cloudwatch-application-insights.html [y/N]: N

Would you like to set Structured Logging in JSON format on your Lambda functions?  [y/N]: y
Structured Logging in JSON format might incur an additional cost. View https://docs.aws.amazon.com/lambda/latest/dg/monitoring-cloudwatchlogs.html#monitoring-cloudwatchlogs-pricing for more details

Project name [sam-app]: url-shortener-lambda

Cloning from https://github.com/aws/aws-sam-cli-app-templates (process may take a moment)

    -----------------------
    Generating application:
    -----------------------
    Name: url-shortener-lambda
    Runtime: go (provided.al2023)
    Architectures: x86_64
    Dependency Manager: mod
    Application Template: hello-world
    Output Directory: .
    Configuration file: url-shortener-lambda\samconfig.toml

    Next steps can be found in the README file at url-shortener-lambda\README.md


Commands you can use next
=========================
[*] Create pipeline: cd url-shortener-lambda && sam pipeline init --bootstrap
[*] Validate SAM template: cd url-shortener-lambda && sam validate
[*] Test Function in the Cloud: cd url-shortener-lambda && sam sync --stack-name {stack-name} --watch


SAM CLI update available (1.151.0); (1.145.0 installed)
To download: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html
```
## And run?
The Lambda created is an example of using the API Gateway functionality and can be executed right out of the box using the provided event. However, it didn’t work for me using the browser. The easiest way to test this is to provide a synthetic event which pretends to be an API call coming through the API gateway. The -e switch is for passing in the event file.
```bash
sam build
sam local invoke -e events/event.json

Mounting D:\Data\workspaces\golang\packt-mastering-go-devops-automating-scaling-2025\ch13\url-shortener-lambda\.aws-sam\build\UrlShortenerFunction as /var/task:ro,delegated,     
inside runtime container
SAM_CONTAINER_ID: 92171740e81fd0ac081231395922ec48dcd4c95da0c092f96fcddbcd25dd5433
START RequestId: 8587ce54-faf8-43b4-b560-f30fbff0dfec Version: $LATEST
END RequestId: fe13f34b-bdf0-47b1-a792-22592660d00a
REPORT RequestId: fe13f34b-bdf0-47b1-a792-22592660d00a  Init Duration: 0.06 ms  Duration: 39.43 ms      Billed Duration: 40 ms  Memory Size: 128 MB     Max Memory Used: 128 MB 
{"statusCode": 200, "headers": null, "multiValueHeaders": null, "body": "Hello, 127.0.0.1!\n"}
```
* Using the API Gateway
  * [Startup](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/using-sam-cli-configure.html)
    ```bash 
     sam local start-api --debug
    ```
  *  Send request 1
    ```bash 
     curl --location 'http://127.0.0.1:3000/hello' \
     --header 'Content-Type: application/json' 

     Hello, 127.0.0.1!
    ```
  *  Send request 2
    ```bash 
     curl --location 'http://127.0.0.1:3000/Prod/hello' \
     --header 'Content-Type: application/json' 

    {"message":"Missing Authentication Token"}
    ```  

## Managing triggers and event sources
### Setting up an S3 trigger
#### Method01: using `sam local invoke`
```bash
sam build
sam local invoke UrlShortenerFunction --event events/s3-event.json
```
#### Method02: using `sam local invoke`,`LocalStack` to stimulating AWS environment
* Step1 : build docker network   
  First, establish a shared Docker network to enable communication between the SAM Lambda containers and LocalStack:
  ```bash
  docker network create sam-local-network
  ```
* Step2 : Start LocalStack
  ```bash
  docker run -d --name localstack \
  --network sam-local-network \
  -p 4510-4559:4510-4559 \
  -e SERVICES=s3 \
  -e DEBUG=1  \
  localstack/localstack

  # inspect container
  docker ps --filter name=localstack
  ```
* Step3: Create an S3 Bucket in LocalStack
  ```bash
  aws --endpoint-url=http://localhost:4566 s3 mb s3://my-test-bucket
  
  # test cp
  aws --endpoint-url=http://localhost:4566 s3 cp test.json s3://my-test-bucket/
  ```
* Step4: Create an S3 event file for testing.
  ```json
  {
    "Records": [
      {
        "eventVersion": "2.1",
        "eventSource": "aws:s3",
        "eventName": "ObjectCreated:Put",
        "s3": {
          "bucket": {
            "name": "my-test-bucket"
          },
          "object": {
            "key": "test-file.json"
          }
        }
      }
    ]
  }
  ```
* Step5: Building a SAM application
  ```bash
  sam build
  ```
* Step6: Use SAM to locally call Lambda (connect to LocalStack)
  ```bash
  sam local invoke UrlShortenerFunction \
  --event events/s3-event.json \
  --docker-network sam-local-network \
  --env-vars env.json
  ```
* env.json 說明
  ```json
  {
    "UrlShortenerFunction": {
      "AWS_ENDPOINT_URL": "http://localstack:4566",
      "BUCKET_NAME": "my-test-bucket",
      "AWS_ACCESS_KEY_ID": "test",
      "AWS_SECRET_ACCESS_KEY": "test",
      "AWS_REGION": "us-east-1"
    }
  }
  ```
| 變數| 用途| 
|-----| ----| 
| AWS_ENDPOINT_URL| 指向 LocalStack (使用容器名稱 localstack)| 
| BUCKET_NAME| S3 bucket 名稱| 
| AWS_ACCESS_KEY_ID| LocalStack 接受任意值| 
| AWS_SECRET_ACCESS_KEY| LocalStack 接受任意值| 

* 如果 Lambda 需要讀取 S3 檔案, 更新 main.go 以連接 LocalStack：
```golang
import (
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

func handler(ctx context.Context, s3Event events.S3Event) error {
    // AWS SDK v2 會自動讀取 AWS_ENDPOINT_URL 環境變數
    cfg, _ := config.LoadDefaultConfig(ctx)
    client := s3.NewFromConfig(cfg)
    
    // 使用 client 讀取 S3...
}
```
* Clean resource
  ```bash 
  docker stop localstack && docker rm localstack
  docker network rm sam-local-network
  ```

### Setting up an API Gateway trigger
* Using the API Gateway
  * [Startup](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/using-sam-cli-configure.html)
    ```bash 
     sam local start-api --debug
    ```
  *  Send request 1
    ```bash 
     curl --location 'http://127.0.0.1:3000/hello' \
     --header 'Content-Type: application/json'  \
     --data  '{"message": "hello world"}'
     
     Hello, 127.0.0.1!
    ```
  *  Send request 2
    ```bash 
     curl --location 'http://127.0.0.1:3000/Prod/hello' \
     --header 'Content-Type: application/json' 

    {"message":"Missing Authentication Token"}
    ```  
