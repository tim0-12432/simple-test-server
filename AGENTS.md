<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# Project Rules (Simple Test Server)

## Description

Simple Test Server is an application made for software testers and developers.
It should make it really easy to spin up a test server providing a specified protocol quickly.
This is done by starting a Docker container providing that specified connectivity via a web application.
There the configuration of that container to start should be set up, the status observed and communication checked.

Examples:

1. If you want to test some MQTT connectivity, you should be able to start a small MQTT broker in the webapp, which you can connect to easily. Then you should be able to observe published topics and messages. In the end you could stop the running container.
2. To test some REST requests, you could start a Webserver, upload some resources, that you can request then.
3. For testing mail sending, you should be able to start a small mail server. The application to test could send mails over it, that should be listed then in the webapp.

## Tech Stack

### Backend

- Language: Go v1.25
- API Framework: Gin v1.10
- Configuration parsing: Viper v1.20

### Frontend

- is served by backend
- Framework: React v19
- Language: Typescript v5
- Package manager: Bun v1.2
- Build tool: Vite v7
- Components: Shadcn UI
- Styling: TailwindCSS v4

### Database

- Pocketbase NoSQL
- started from inside Go backend
- reachable over path `/pb/`

## Folder Structure

### Overall

- .github/: GitHub Automations/Workflows
- .vscode/: Configuration files for IDE + REST test files
- ... Backend folders
- custom_images/: Dockerfiles and configuration for custom Docker images
- frontend/: Frontend stuff
- migrations/: Database migration files
- pb_data/: Database data
- app.env: Environment variables
- docker-compose.yml: Docker compose for deployment
- Dockerfile: Definition for Docker image 

### Backend

- config/: Configuration parsing from environment and default handling
- controllers/: REST API definitions calling the different services
- db/
  - dtos/: Database entities
  - services/: Services managing the different database collections of the DTOs
  - collection_bootstrap.go: Sets up the collections if database started freshly
  - pocketbase.go: Manages the connection to the database
  - utils.go: Some helper functions for database interaction
- docker/
  - servers/: Definitions for the different test/development servers
  - builder.go: Builds custom images if they don't exist yet
  - manager.go: Manages whole interaction with Docker
  - runner.go: Starts new Docker containers
- progress/hub.go: Registers eventsources for enabling progress tracking
- protocols/: Section for test server type specific logic
  - routes.go: Registers test server type specific APIs
  - mqtt/ or web/ etc.: Controller and services for this specific test server type
- go.mod: Go project definition and dependencies
- main.go: Application entrypoint starting the server

### Frontend

- frontend/
  - dist/: Build frontend files shipped by Go server
  - src/
    - assets/: All images, fonts whatever to be included in the webapp
    - components/: Components collected from different shadcn libraries e.g. KiboUi, but some customized or self implemented
    - lib/: Components and hooks used in different places over the webapp
    - tabs/: Factory for the different tab pages for each test server type
    - types/: Typescript types for DTOs
    - App.tsx
    - index.css
    - main.tsx
  - components.json: Shadcn configuration
  - eslint.config.js
  - index.html
  - package.json
  - tsconfig.json
  - vite.config.ts

## Testing

### Frontend

1. `cd frontend`
2. if dependencies needed: `bun install`
3. `bun run test`

### Backend

`go test ./...` executes all tests
