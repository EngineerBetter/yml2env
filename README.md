# yml2env

Executes a command with environment variables taken from a YAML file.

```
yml2env <path-to-yaml-file> <command>
```

## Why?

It's quite handy for using Concourse `--load-vars-from` files when running local tasks, like tests.

## Example

Given a YAML file:

```
---
cf_username: admin
cf_password: whevsmate
```

...running `yml2env ci/vars/local.yml fly execute ci/tasks/system-tests.yml` is equivalent to running

```
CF_USERNAME=admin CF_PASSWORD=admin fly execute ci/tasks/system-tests.yml
```

## Installation

### Go developers

```
go get github.com/EngineerBetter/yml2env
```

### Everyone else

1. Download [a release](https://github.com/EngineerBetter/yml2env/releases)
1. Move to `$PATH` and rename to `yml2env`

## Testing

```
go test ./...
```