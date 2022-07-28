# ccienv

A management tool for CircleCI Project's environment variables.  
(Currently, only available on GitHub)

This repository will be supported until the official circleci-cli will support the project's variables management: [issue](https://github.com/CircleCI-Public/circleci-cli/issues/652)

## Installation

```
go install github.com/threepipes/circleci-env/cmd/ccienv@latest
```

### Uninstallation

```
rm $(which ccienv)
```

## Requirements

- git
    - Only if load repository name from .git

## Setup

```
ccienv config init
```

Set these variables.

- CircleCI API Token
    - A personal API token of CircleCI
- GitHub organization
    - GitHub organization name or GitHub username of your repository

Then, `$XDG_CONFIG_HOME/ccienv/config.yml` will be created.

## Run

```
ccienv -r <your_repo_name> <cmd> [<args>]
```

If `-r <your_repo_name>` is omitted, the current directory name is used as the target repository name.

### Example

```
ccienv -r circleci-env ls
```

## Help

```
ccienv -h
```
