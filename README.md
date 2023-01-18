# cinefinder-api

Find movies based on the vaguest memory. Aspiring to stand on the shoulder of machine learning giants.

## Description

Cinefinder is a web application that allows users to search for movies and TV shows from their vaguest memory. It is built with Go and uses Elasticsearch as the search engine.
The goal is to learn how to build a database of movies and TV shows, and then implement NLP methods on top of it to perform semantic search.
This is only a prototype of what we want it to look like. In the meantime, we'll be outsourcing those technologies from existing search engine giants and magicks.

## Installation
- Git clone this repo
- Setup env variables
- Run `go mod tidy`
- Run `go run main.go`

Or run from docker

- Build docker image
```sh
docker build -t cinefinder-api .
```

- Run container
```
docker run -it -p 8009:8009 cinefinder-api
```
- Visit `http://localhost:8009`

## Resources
- [Elasticsearch](https://www.elastic.co/docs)
- [Playwright](https://github.com/playwright-community/playwright-go)