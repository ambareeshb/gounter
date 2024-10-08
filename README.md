# Gounter

**Gounter** is a simple counter application built in Go that provides RESTful APIs for managing counters. 

## Table of Contents

- [Gounter](#gounter)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Usage](#usage)
  - [API Documentation](#api-documentation)
  - [Makefile Commands](#makefile-commands)

## Installation

Make sure you have [Docker](https://www.docker.com/get-started) installed on your machine. Clone this repository and navigate to the project directory:

```bash
git clone https://github.com/ambareeshb/gounter
```

## Usage
```bash
cd gounter
make run
```

## API Documentation
The Swagger documentation for the APIs is available at:

http://localhost:8080

You can view the API endpoints and their details there.

## Makefile Commands
This project uses a Makefile for various tasks. Here are the available commands:

```bash
make build: Build the application.
make run: Run the application.
make test: Run unit tests.
make test-integration: Run integration tests (the application must be running).
make install-deps: Install richgo if it is not available.
make clean: Clean up build artifacts and stop the application.
make serve-swagger: Serve documentation on http://localhost:8080
make create-migration: Create a new database migration file.
```
