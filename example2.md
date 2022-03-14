<div align="center">
  <a href="https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-logo.png">
    <img src="https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-logo.png" alt="Logo" height="100">
  </a>
</div>
<div>
  <p align="center">
    </br>Catering Ecommerce Platform</br></br>
    <a href="https://app.swaggerhub.com/apis-docs/furqonzt99/snackbox/1">API Docs</a>
    ·
    <a href="https://whimsical.com/snackbox-UcYKhew5MBhFzJWaCXQbAb">Wireflow</a>
    ·
    <a href="https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-ucd.png">Use Case Diagram</a>
    ·
    <a href="https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-erd.png">Entity Relationship Diagram</a>
  </p>
</div>

# Snackbox

[![Go.Dev reference](https://img.shields.io/badge/echo-reference-blue?logo=go&logoColor=blue)](https://github.com/labstack/echo)
[![Go.Dev reference](https://img.shields.io/badge/gorm-reference-blue?logo=go&logoColor=blue)](https://pkg.go.dev/gorm.io/gorm?tab=doc)
[![Go.Dev reference](https://img.shields.io/badge/aws--s3-reference-orange)](https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/go/example_code/s3)
[![Go.Dev reference](https://img.shields.io/badge/maroto-reference-blue)](https://pkg.go.dev/github.com/johnfercher/maroto?tab=doc)

Organizations or offices need to set up a consumption section that is in charge of providing food, where they have to find restaurants and contact them all the time to track their orders to ensure the orders are processed properly.

Snackbox can make it easier for customers to order snacks and rice boxes and do tracking, providing security for customers because the connected partners have passed the verification process and this application will provide employment for partners.

## Features

- JWT Authentication
- Multi Role Middleware (Admin, Partner, User)
- Search Product By Nearest Partner Location
- Payment Gateway Integration - Invoice & Disbursment (Xendit)
- AWS S3 Integration
- PDF Export (Maroto)

## High Level Architecture

![High Level Architecture](https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-hla.png)

## Tech Stack

- [Github](https://github.com/) - Versioning Platform
- [Trello](https://trello.com/) - Collaboration Platform
- [Go](https://go.dev/) - Project Language
- [Echo](https://echo.labstack.com/) - Go Web Framework
- [MySql](https://www.mysql.com/) - SQL Database
- [Ngrok](https://ngrok.com/) - Expose local url to public url for test with third-party API
- [Xendit](https://docs.xendit.co/) - Payment gateway API
- [AWS S3](https://aws.amazon.com/s3/) - Object storage service
- [AWS EC2](https://aws.amazon.com/ec2/) - Virtual computer service
- [AWS RDS](https://aws.amazon.com/rds/) - Relational database service
- [Docker](https://www.docker.com/) - Container Registry
- [Kubernetes](https://kubernetes.io/) - Container Orchestration

## Structure

![Structure](https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-structure.png)

## Unit Test

![Unit Test](https://github.com/furqonzt99/snackbox/blob/documentation/documentation/test-result.png)

## Installation

- Clone this repo

```bash
git clone https://github.com/furqonzt99/snackbox.git snackbox
```

- Go to repository folder

```bash
cd snackbox
go mod tidy
```

- Create .env file and add the following environment (you can see the variables from .env.example)

- Run this app

```bash
go run .
```

## High Level Architecture

![High Level Architecture](https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-hla.png)

## Structure

![Structure](https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-structure.png)

## Unit Test

![Unit Test](https://github.com/furqonzt99/snackbox/blob/documentation/documentation/test-result.png)

## Tech Stack

- [Github](https://github.com/) - Versioning Platform
- [Trello](https://trello.com/) - Collaboration Platform
- [Go](https://go.dev/) - Project Language
- [Echo](https://echo.labstack.com/) - Go Web Framework
- [MySql](https://www.mysql.com/) - SQL Database
- [Ngrok](https://ngrok.com/) - Expose local url to public url for test with third-party API
- [Xendit](https://docs.xendit.co/) - Payment gateway API
- [AWS RDS](https://aws.amazon.com/rds/) - Relational database service
- [AWS S3](https://aws.amazon.com/s3/) - Object storage service
- [Okteto](https://www.okteto.com/) - Kubernetes development platform
- [Docker](https://www.docker.com/) - Container Registry
- [Kubernetes](https://kubernetes.io/) - Container Orchestration

## Authors

- [@furqonzt99](https://github.com/furqonzt99) - Product Owner & Developer
- [@yogawahyudi7](https://github.com/yogawahyudi7) - Developer