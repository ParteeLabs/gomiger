# Mongomiger Contribution Guidelines

## Setup local environment for testing

To set up a local MongoDB instance for testing, you can use Docker Compose. A sample `docker-compose.yaml` file is provided in the repository. To start the MongoDB service, run the following command in the terminal:

```bash
docker-compose -f mongomiger/docker-compose.yaml up -d
```

This will start a MongoDB container that you can use for testing your migrations. Make sure Docker is installed and running on your machine before executing the command.

## Running Tests

To run the tests for Mongomiger, ensure that your local MongoDB instance is running (as described above). Then, execute the following command in the terminal from the root directory of the project:

```bash
go test ./mongomiger/... -coverprofile=coverage.out
```

This command will run all the tests in the `mongomiger` package and generate a coverage report. You can view the coverage report using:

```bash
go tool cover -html=coverage.out
```
