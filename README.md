# Scout Me

A sibling project to [zoubun](github.com/jacobsky/zoubun), this is a project where the primary motivation is to get comfortable with the go ecosystem.

In this case, the goal is to create a hypermedia driven microblog/portfolio application that sends me emails with exactly the information that I can about when being scouted.

For speed, the project's structure was created using the [go-blueprint](https://github.com/Melkeydev/go-blueprint) tool.

## Stack

- *Frontend*: Alpine.js (UI interactivity) + HTMX (backend interactivity)
- *Styling*: Tailwind + Pines (components) + picocss classless (base)
- *Backend*: Golang + Temple

## Goals
- [ ] Configure HTMX, TailwindCSS, Golang and Templ to create a reasonably presentable landing page.
- [ ] Using a flowing form, send me a well formatted email notifying me of a potential position.
- [ ] Small profile section to list out skills that may match your needs (basically a copy [jacobsky](github.com/jacobsky/jacobsky))

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
