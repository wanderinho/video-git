[Русский](README.md) | [English](README.en.md)

---

# Video Git

**Video Git** is a backend project for version control and collaboration on video content, inspired by Git and GitHub.

The goal is to provide video editors with a structured way to manage projects, create branches, track video revisions, and collaborate without losing previous versions of their work.

> **Project status:** Work in progress. Core backend functionality is under active development.

## Features

### Implemented / In Progress

- User registration and authentication with JWT
- Project creation and management
- Project membership and access control
- Branch creation from existing video commits
- PostgreSQL transactions for multi-step operations
- Database migrations
- Docker-based development environment

### Planned

- Video uploads to S3-compatible object storage
- Video version history and parent-child commit relationships
- Asynchronous video processing with FFmpeg
- Video metadata extraction
- Commit processing statuses
- Background jobs and task queues
- Video previews and comparison between revisions
- Additional collaboration features

## Tech Stack

| Category | Technology |
|---|---|
| Language | Go |
| HTTP framework | Echo |
| Database | PostgreSQL |
| Database driver | pgx |
| Object storage | MinIO |
| Background jobs | Redis, Asynq |
| Video processing | FFmpeg |
| Authentication | JWT |
| Infrastructure | Docker, Docker Compose |
| Code quality | golangci-lint |

## Architecture

The application follows a layered architecture with a separation of responsibilities.

```text
HTTP Request
     |
     v
   Handler
     |
     v
   Service
     |
     v
  Repository
     |
     v
 PostgreSQL
```

**Handler layer** is responsible for HTTP request parsing, input validation, authentication context, and HTTP responses.

**Service layer** contains business rules and coordinates application operations.

**Repository layer** handles database queries, transactions, persistence, and database-specific error handling.

The project also includes a planned background worker for asynchronous video processing.

## Core Concepts

### Projects

A project represents a workspace for video editing and version management. Each project has an owner, members, and a default branch.

Project creation is performed within a database transaction to ensure that the project, its default branch, and the owner's membership are created consistently.

### Branches

Branches allow editors to work on different versions or directions of a video project.

A new branch is created from an existing, successfully processed commit. Its initial HEAD points to that commit without creating a duplicate commit.

Branch names must be unique within a project, using case-insensitive comparison.

### Commits

A commit represents a version of video content.

Commits are designed to reference parent commits, allowing the application to reconstruct version history and relationships between branches.

Video processing will be performed asynchronously, with commit status and metadata updated after processing.

## Data Consistency

The backend uses PostgreSQL transactions to maintain consistency across multi-step operations.

Examples include:

- Creating a project together with its default branch and membership record
- Validating project access and creating branches
- Maintaining relationships between commits and branch HEAD pointers

Concurrency control and transaction isolation are areas of ongoing development.

## Getting Started

The project is still under development, so installation and startup instructions will be finalized as the infrastructure evolves.

The intended local development environment uses Docker Compose for the API, PostgreSQL, Redis, MinIO, and background worker services.

## Development

The project uses `golangci-lint` for static analysis and formatting checks.

The development process focuses on:

- Clear separation of application layers
- Explicit error handling
- Transactional database operations
- Maintainable Go code
- Incremental implementation of backend features

## Roadmap

- [x] Initial project structure
- [x] Database migrations
- [x] Authentication foundation
- [x] Project creation
- [x] Branch creation logic
- [ ] Complete project and branch management
- [ ] Implement commit creation and history
- [ ] Integrate video storage with MinIO
- [ ] Implement asynchronous FFmpeg worker
- [ ] Add video metadata extraction
- [ ] Add automated tests
- [ ] Complete API documentation
- [ ] Prepare production deployment

## Motivation

Video editing workflows often rely on manually named files, duplicated project folders, and external communication to track revisions.

Video Git explores how concepts from software version control — commits, branches, history, and collaboration — can be adapted to video production.

The project is also an opportunity to practice backend engineering with Go, PostgreSQL, distributed components, asynchronous processing, and reliable data management.

## License

License to be determined.
