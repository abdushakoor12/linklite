# LinkLite

A lightweight URL shortener built with Go, using templ for templating and NES.css for retro-style UI.

## Features
- Shorten long URLs to memorable short links
- Clean, retro-style user interface
- Simple JSON-based storage
- Redirect service for shortened URLs

## Tech Stack
- Go 1.x
- [templ](https://templ.guide/) for HTML templating
- [NES.css](https://nostalgic-css.github.io/NES.css/) for retro-style UI
- JSON file for persistent storage

## Getting Started

### Prerequisites
- Go 1.x or higher
- templ CLI tool

### Installation
```bash
# Install templ
go install github.com/a-h/templ/cmd/templ@latest

# Install dependencies
go mod tidy
```

### Development
```bash
# Generate templates
templ generate

# Run the server
go run main.go
```

### Usage
1. Visit http://localhost:8080
2. Enter a long URL in the form
3. Get a shortened URL
4. Use the shortened URL to be redirected to the original URL
