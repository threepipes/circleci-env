# ccienv

A management tool for CircleCI Project's environment variables.  
(Currently, only available on GitHub)

## Installation

```
go install github.com/threepipes/circleci-env/cmd/ccienv@latest
```

### Uninstallation

```
rm $(which ccienv)
```

## Setup local environment

These environment variables are required.

- CIRCLECI_API_TOKEN
    - A personal API token of CircleCI
- CIRCLECI_ORGANIZATION_NAME
    - GitHub organization name or GitHub username of your repository

## Run

```
ccienv -r <your_repo_name> <cmd> [<args>]
```

### Example

```
ccienv -r circleci-env ls
```

## Help

```
ccienv -h
```
