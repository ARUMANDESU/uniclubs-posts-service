<a id="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://astanait.edu.kz/">
    <img src="https://static.tildacdn.pro/tild3764-6633-4663-b138-303730646233/aitu-logo__2.png" alt="Logo" height="80">
  </a>
  <h3 align="center">AITU UCMS Posts Service</h3>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#protofiles">Protofiles</a></li>
    <li><a href="#technologies-used">Technologies Used</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <ul>
      <li><a href="#prerequisites">Prerequisites</a></li>
      <li><a href="#installation">Installation</a></li>
      <li><a href="#configuration">Configuration</a></li>
    </ul>
    <li><a href="#running-the-service">Running the Service</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project
This service is part of the University Clubs Management System (UCMS) project. The service is responsible for managing posts( announcements, events, polls, etc.) for clubs. It provides a gRPC API for other services to interact with it.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- PROTOFILES -->
## Protofiles

* [Protofiles Repository][protofiles-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- TECHNOLOGIES USED -->
## Technologies Used

* [![Go][go-shield]][go-url]
* [![MongoDB][mongodb-shield]][mongodb-url]
* [![gRPC][grpc-shield]][go-url]
* [![RabbitMQ][rabbitmq-shield]][rabbitmq-url]
* [![Docker][docker-shield]][docker-url]
* [![Docker Compose][docker-compose-shield]][docker-compose-url]
* [![Taskfile][tasks-shield]][tasks-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Getting Started
### Prerequisites
- Go version 1.22.2
- Docker 26.1.4


### Installation
Clone the repository:
   ```bash
   git clone https://github.com/ARUMANDESU/uniclubs-posts-service.git
   cd uniclubs-posts-service
   go mod download
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


### Configuration
This Service requires a configuration file to specify various settings like database connections,
and service-specific parameters. 
Depending on your environment (development, test, or production), different configurations may be needed.

#### Setting Up Configuration
```dotenv
# Example configuration snippet
ENV=dev
# Database Configuration
MONGODB_URI=mongodb://username:password@host:port
MONGODB_PING_TIMEOUT=
MONGODB_DATABASE_NAME=
# Server Configuration
GRPC_PORT=
GRPC_TIMEOUT=
# RabbitMQ Configuration
RABBITMQ_USER=
RABBITMQ_PASSWORD=
RABBITMQ_HOST=
RABBITMQ_PORT=
# Other Services client configuration
USER_SERVICE_ADDRESS=
USER_SERVICE_TIMEOUT=
USER_SERVICE_RETRIES_COUNT=
CLUB_SERVICE_ADDRESS=
CLUB_SERVICE_TIMEOUT=
CLUB_SERVICE_RETRIES_COUNT=
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Running the Service
After setting up the database and configuring the service, you can run it as follows:
```bash
go run cmd/main.go
```

Or use the provided Taskfile to run the service:
```bash
task run:enviroment
 ```
or
```bash
task r:e
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[aitu-url]: https://astanait.edu.kz/
[aitu-ucms-url]: https://www.ucms.space/
[protofiles-url]: https://github.com/ARUMANDESU/uniclubs-protos

[go-url]: https://golang.org/
[docker-url]: https://www.docker.com/
[docker-compose-url]: https://docs.docker.com/compose/
[redis-url]: https://redis.io/
[mongodb-url]: https://www.mongodb.com/
[grpc-url]: https://grpc.io/
[rabbitmq-url]: https://www.rabbitmq.com/
[tasks-url]: https://taskfile.dev/

[go-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[docker-shield]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[docker-compose-shield]: https://img.shields.io/badge/Docker_Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white
[redis-shield]: https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white
[mongodb-shield]: https://img.shields.io/badge/MongoDB-47A248?style=for-the-badge&logo=mongodb&logoColor=white
[grpc-shield]: https://img.shields.io/badge/gRPC-008FC7?style=for-the-badge&logo=google&logoColor=white
[rabbitmq-shield]: https://img.shields.io/badge/RabbitMQ-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white
[tasks-shield]: https://img.shields.io/badge/Taskfile-00ADD8?style=for-the-badge&logo=go&logoColor=white
