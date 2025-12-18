# Spring Boot Multitenancy

## Stack
* Grafana OSS

## Usage

You can use Docker Compose to run the necessary backing services for observability, authentication, and AI.

From the project root folder, run Docker Compose.

```bash
docker-compose up -d
```
* prometheus: http://localhost:9090/
* grafana: http://localhost:3000/
* chroma: http://localhost:8000/
* tempo: http://localhost:3110
* Zipkin: http://localhost:9411
