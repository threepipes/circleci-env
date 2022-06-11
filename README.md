# ccienv

A management tool for CircleCI Project's Environment Variables.  
(Currently, only available on GitHub)

## Installation

```
go install github.com/threepipes/circleci-env/cmd/ccienv@latest
```

uninstall
```
rm $(which ccienv)
```

## Setup local environment

Set these environment variables to run the command.
- CIRCLECI_API_TOKEN
    - A personal api token of CircleCI
- CIRCLECI_ORGANIZATION_NAME
    - Organization name or github username of your repository

## Run

```
ccienv --repo <your_repo_name> <cmd> [<args>]
```

example.
```
ccienv --repo circleci-env ls
```

## Help

```
ccienv -h
```
