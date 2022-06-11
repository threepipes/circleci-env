# ccienv

Manage CircleCI Project's Environment Variables.
(Currently, only available on github)

# Installation

```
go install github.com/threepipes/ccienv@latest
```

uninstall
```
rm $(which ccienv)
```

# Setup local environment

Set these environment variables to run the command.
- CIRCLECI_API_TOKEN
    - A personal api token of CircleCI
- CIRCLECI_ORGANIZATION_NAME
    - Organization name or github username of your repository

# Run

```
ccienv --repo <your_repo_name> <cmd> [<args>]
```

example.
```
ccienv --repo circleci-cli ls
```

## Help

```
ccienv -h
```
