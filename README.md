# Github Worflow Explorer

A graphical Terminal application that enables the listing, filtering and deletion of Github's Repos Workflow runs.

![Main Window](docs/imgs/01.png)

[Click for more sreenshots](docs/screenshots.md)

## Usage

```
Usage of shipper: gh-we [-c | --config <path>] [-a | --auth <string>]
        -c, --config  path to configuation json file
        -a, --auth    github Token used for authentication. Overwrites json configurations.
	-r, --repo    github repo to view. Overwrites json configurations.
        -h, --help    display this help message
```

## Configuration

By default the gh-we search for a `config.json` file in the same folder that invoke the application.
These json have the following structure:

```json
{
  "owner": "OWNER",
  "name": "REPO",
  "auth": "token",
  "githubApiVersion": "2022-11-28"
}
```

| Field            | Description                                    |
| ---------------- | ---------------------------------------------- |
| owner            | Repo's owner. Either is a user or organization |
| name             | Repo's name                                    |
| auth             | Token to use as the bearer authentication      |
| githubApiVersion | GitHub API Versions                            |
