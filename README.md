# yml2env

![Build Status](http://ci.engineerbetter.com/api/v1/teams/main/pipelines/yml2env/jobs/test/badge)

Either executes a command with environment variables taken from a YAML file, or prints a load of `export`s you can `eval`

```sh
yml2env <path-to-yaml-file> [<command> | --eval]

# Run command with env vars from YAML file
$ yml2env vars.yml tests.sh

# Set env vars in current shell
$ eval "$(yml2env var.yml --eval)"
```

## Why?

It's quite handy for using Concourse `--load-vars-from` files when running local tasks, like tests. The `--eval` feature is useful when you need to get lots of stuff from the output of a Concourse Terraform resource as env vars.

## Example

Given a YAML file stored in `ci/vars/local.yml`:

```
---
cf_username: admin
cf_password: whevsmate
```

...running `yml2env ci/vars/local.yml fly execute ci/tasks/system-tests.yml` is equivalent to running

```
CF_USERNAME=admin CF_PASSWORD=whevsmate fly execute ci/tasks/system-tests.yml
```

## Installation

### Go developers

```
go get github.com/EngineerBetter/yml2env
```

### Everyone else

1. Download [a release](https://github.com/EngineerBetter/yml2env/releases)
1. Move to `$PATH` and rename to `yml2env`
1. `chmod +x yml2env`

## Testing

```
go test ./...
```
