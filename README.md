# yml2env

Executes a command with environment variables taken from a YAML file.

```
yml2env <path-to-yaml-file> <command> 
```

## Installation

### Go developers

```
go get github.com/EngineerBetter/yml2env
```

## Example

Given a YAML file:

```
---
cf_username: admin
cf_password: admin
```

...running `yml2env ci/vars/local.yml fly execute ci/tasks/system-tests.yml` is equivalent to running

```
CF_USERNAME=admin CF_PASSWORD=admin fly execute ci/tasks/system-tests.yml
```

## Testing

```
go test ./...
```