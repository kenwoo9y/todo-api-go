# todo-api-go

This is a ToDo Web API implemented with Go, designed for simplicity and extensibility.

## Tech Stack

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MySQL](https://img.shields.io/badge/mysql-4479A1.svg?style=for-the-badge&logo=mysql&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/github%20actions-%232671E5.svg?style=for-the-badge&logo=githubactions&logoColor=white)

### Programming Languages
- [Go](https://go.dev/) v1.23.2 - Primary development language

### Backend
- Standard library `net/http` - Go's built-in HTTP server
- Standard library `database/sql` - Go's built-in database interface

### Database
- [MySQL](https://www.mysql.com/) v8.0 - Primary relational database
- [PostgreSQL](https://www.postgresql.org/) v16 - Alternative relational database

### Development Environment
- [Docker](https://www.docker.com/) with Compose - Containerization platform for building and managing applications

### Testing & Quality Assurance
- [go testing](https://pkg.go.dev/testing) - Go's built-in testing framework
- [go vet](https://pkg.go.dev/cmd/vet) - Go's built-in static analysis tool
- [go fmt](https://pkg.go.dev/cmd/gofmt) - Go's built-in code formatter

### CI/CD
- GitHub Actions - Continuous Integration and Deployment

## Setup
### Initial Setup
1. Clone this repository:
    ```
    $ git clone https://github.com/kenwoo9y/todo-api-go.git
    $ cd todo-api-go
    ```

2. Create environment file:
    ```
    $ cp .env.example .env
    ```
    Edit `.env` file to match your environment if needed.

3. Build the required Docker images:
    ```
    $ make build-local
    ```

4. Start the containers:
    ```
    $ make up
    ```

5. Apply database migrations:
    ```
    $ make migrate-mysql
    $ make migrate-psql
    ```

## Usage
### Container Management
- Check container status:
    ```
    $ make ps
    ```
- View container logs:
    ```
    $ make logs
    ```
- Stop containers:
    ```
    $ make down
    ```

## Development
### Running Tests
- Run tests:
    ```
    $ make test
    ```
- Run tests with coverage:
    ```
    $ make test-coverage
    ```

### Code Quality Checks
- Lint check:
    ```
    $ make lint
    ```
- Apply code formatting:
    ```
    $ make format

## Database
### Switching Database
1. Edit `.env` file:

For MySQL:
```
DB_TYPE=mysql
DB_HOST=mysql-db
DB_PORT=3306
DB_NAME=todo
DB_USER=<your_username>
DB_PASSWORD=<your_password>
```

For PostgreSQL:
```
DB_TYPE=psql
DB_HOST=postgresql-db
DB_PORT=5432
DB_NAME=todo
DB_USER=<your_username>
DB_PASSWORD=<your_password>
```

2. Rebuild and restart the application:
```
$ make build-local
$ make up
$ make migrate-mysql
$ make migrate-psql
```

### Database Access
- Access MySQL database:
    ```
    $ make mysql
    ```
- Access PostgreSQL database:
    ```
    $ make psql
    ```
