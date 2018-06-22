# go-octopusdeploy
A Go wrapper for the Octopus Deploy REST API.

[![Build status](https://ci.appveyor.com/api/projects/status/5t5gbqjyl8hpou52?svg=true)](https://ci.appveyor.com/project/MattHodge/go-octopusdeploy)

> :warning: The Octopus Deploy REST Client is in heavy development.

# Go Dependencies
* Dependencies are managed using [dep](https://golang.github.io/dep/docs/new-project.html)

# Using the main.go Example

```bash
export OCTOPUS_URL=http://localhost:8081/
export OCTOPUS_APIKEY=API-FAKEAPIKEYFAKEAPIKEY

go run main.go # creates a project
```
