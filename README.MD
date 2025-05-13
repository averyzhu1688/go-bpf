# Go-BPF 

## Project introduction

Go-BPF is a web application skeleton developed using Go language and Gin framework, integrating commonly used functions such as user authentication, permission management, and database CRUD operations, providing a basic infrastructure for rapid development of enterprise level web applications.

## Technical components and frameworks
Go 1.23+
Gin Web framework
Gorm ORM
MySQL  Database
Redis Cache
JWT  Authentication
Zap  log
Viper config manage

## Project Structure
api: route
cmd: base script
docs: project doc
internal: project inner code
pkg: exteinal code
test: unit test code
config.yaml: system config file
main.go: System startup entrance

## Function characteristics
Cross domain support
Global exception handling
Automatic database migration
Cache system support 
Comprehensive logging system
Http Request limit

## Core Component Description

### Config manage
- vapper config manage,support yaml 
- server、database、cache、jwt

### Database
- Based on Gorm ORM MySQL database
- Database connection pool configuration

### Log 
- Using the Zap high-performance logging system
- Integrate with the gin framework

## Startup project

1. config.yaml
    - app server config
    - mysql config
    - redis config
    - jwt config
    - log config
2. start server
    bash 
    go run main.go
3.  access API
    default url: http://localhost:8080

## Contract
 - Email : averyzhu1688@gmail.com