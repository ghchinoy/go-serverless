# Building a Serverless Go App on Google Cloud

Inspired by Laurent Picard's [Building a serverless Python app in minutes with GCP](https://medium.com/google-cloud/building-a-serverless-python-app-in-minutes-with-gcp-5184d21a012f), this is a similar
set of steps, but uses a Go application.

## Prerequisite: Install the Google Cloud SDK cli

[Install the Google Cloud SDK command-line tool](https://cloud.google.com/sdk/downloads) - this'll be used throughout.


## Create a Google Cloud Project

Using the [`gcloud alpha create`](https://cloud.google.com/sdk/gcloud/reference/alpha/projects/create) command, create a new project. Note, this is an "alpha release" command, so it may change in the future. The non-alpha way is to use the [web-based cloud console](https://console.cloud.google.com).

```
$ gcloud alpha projects create go-serverless
Create in progress for [https://cloudresourcemanager.googleapis.com/v1/projects/go-serverless].
Waiting for [operations/pc.7509198304344852511] to finish...done.
```

## List the existing projects

If this isn't your first Google Cloud project, you may have multiple projects already, so here's a useful command for checking what projects you have

```
$ gcloud projects list
PROJECT_ID                 NAME                                PROJECT_NUMBER
bespokemirrorapi           bespokemirrorapi                    811093430365
gdgnoco-fortune            gdgnoco-fortune                     861018601285
go-serverless              go-serverless                       49448245715
...
```

## Set a default project

Subsequent `gcloud` commands will be associated with a project, so it's useful to set the default project to our newly created project. Creating a new project doesn't automatically set it to the default one.

List all the current project configs with `gcloud config list` and set a config item with `gcloud set ...`.

Set the `core/project` property to `go-serverless`:

```
$ gcloud config set core/project go-serverless
Updated property [core/project].
derby:go-serverless ghc$ gcloud config list
[compute]
region = us-central1
zone = us-central1-c
[core]
account = ghchinoy@gmail.com
disable_usage_reporting = False
project = go-serverless

Your active configuration is: [default]
```

## Create an app locally

* Create a project directory
* Create a deployment file `app.yaml`

Go projects, by convention, go into a folder under `$GOPATH\src`, typically associated with a repository.  Here, I create a project using my github as a path and change to the created directory:

```
$ mkdir -p $GOPATH/src/github.com/ghchinoy/go-serverless
$ cd $GOPATH/src/github.com/ghchinoy/go-serverless
```

Add a new file called `app.yaml`


```
runtime: go
env: flex
api_version: 1

skip_files:
- README.md
```


Add code


## Test locally

Run the app locally; view at http://localhost:8080

```
$ go run main.go template.go
2017/03/30 15:55:21 booklist
127.0.0.1 - - [30/Mar/2017:15:55:24 -0600] "GET / HTTP/1.1" 200 1584 "" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3056.0 Safari/537.36"
...
```