Chyle [![CircleCI](https://circleci.com/gh/antham/chyle.svg?style=svg)](https://circleci.com/gh/antham/chyle) [![codecov](https://codecov.io/gh/antham/chyle/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/chyle) [![codebeat badge](https://codebeat.co/badges/c3867610-2741-4ae3-a195-d5e9711c7fcd)](https://codebeat.co/projects/github-com-antham-chyle-master) [![Go Report Card](https://goreportcard.com/badge/github.com/antham/chyle)](https://goreportcard.com/report/github.com/antham/chyle) [![GoDoc](https://godoc.org/github.com/antham/chyle?status.svg)](http://godoc.org/github.com/antham/chyle)
=====

Chyle produces a changelog from a git repository.

---

* [Usage](#usage)
* [How it works ?](#how-it-works-)
* [Setup](#setup)
* [Examples](#examples)
* [Documentation](#documentation)
* [Contribute](#contribute)

---

## Usage

```
Create a changelog from your commit history

Usage:
  chyle [command]

Available Commands:
  config      Generate environments variables from a prompt session
  create      Create a new changelog
  help        Help about any command

Flags:
      --debug   enable debugging
  -h, --help    help for chyle

Use "chyle [command] --help" for more information about a command.
```

## How it works ?

Chyle fetch a range of commits using given criterias from a git repository. From those commits you can extract relevant datas from commit message, author, and so on, and add it to original payload. You can afterwards if needed, enrich your payload with various useful datas contacting external apps (shell command, apis, ....) and finally, you can publish what you harvested (to an external api, stdout, ....). You can mix all steps together, avoid some, combine some, it's up to you.

## Setup

Download from release page according to your architecture chyle binary : https://github.com/antham/chyle/releases

You need afterwards to configure each module through environments variables : there are activated when you configure at least one environment variable they need to work.

## Examples

We will use this repository : [https://github.com/antham/test-git](https://github.com/antham/test-git), created for chyle testing purpose only, you can try examples on it. Don't forget to clone repository and adapt some environment variables to your configuration.

---

* [Get a JSON ouput of all merge commits](#get-a-json-ouput-of-all-merge-commits)
* [Get a markdown ouput of merge and regular commits](#get-a-markdown-ouput-of-merge-and-regular-commits)
* [Get a JSON ouput of all merge commits and contact github issue api to enrich payload](#get-a-json-ouput-of-all-merge-commits-and-contact-github-issue-api-to-enrich-payload)
* [Populate a release in github from CircleCI](#populate-a-release-in-github-from-circleci)

---

### Get a JSON ouput of all merge commits

commands :

```bash
export CHYLE_GIT_REFERENCE_FROM=a00ee81c109c8787f0ea161a776d2c9795f816cd
export CHYLE_GIT_REFERENCE_TO=f617fb708dfa6fa290205615ea98c53a860e499d
export CHYLE_GIT_REPOSITORY_PATH=/your-local-path/test-git
export CHYLE_MATCHERS_TYPE=merge
export CHYLE_SENDERS_STDOUT_FORMAT="json"

chyle create
```

output :

```json
{
  "datas": [
    {
      "authorDate": "2017-05-10 22:24:40 +0200 +0200",
      "authorEmail": "antham@users.noreply.github.com",
      "authorName": "Anthony HAMON",
      "committerDate": "2017-05-10 22:24:40 +0200 +0200",
      "committerEmail": "noreply@github.com",
      "committerName": "GitHub",
      "id": "f617fb708dfa6fa290205615ea98c53a860e499d",
      "message": "Merge pull request #3 from antham/test2\n\nTest2",
      "type": "merge"
    },
    {
      "authorDate": "2017-05-10 22:22:03 +0200 +0200",
      "authorEmail": "antham@users.noreply.github.com",
      "authorName": "Anthony HAMON",
      "committerDate": "2017-05-10 22:22:03 +0200 +0200",
      "committerEmail": "noreply@github.com",
      "committerName": "GitHub",
      "id": "8fdfae00cbcc66936113a60f5146d110f2ba3c28",
      "message": "Merge pull request #1 from antham/test\n\nTest",
      "type": "merge"
    }
  ],
  "metadatas": {}
}
```

### Get a markdown ouput of merge and regular commits

commands :

```bash
export CHYLE_GIT_REFERENCE_FROM=a00ee81c109c8787f0ea161a776d2c9795f816cd
export CHYLE_GIT_REFERENCE_TO=f617fb708dfa6fa290205615ea98c53a860e499d
export CHYLE_GIT_REPOSITORY_PATH=/your-local-path/test-git
export CHYLE_SENDERS_STDOUT_FORMAT="template"
export CHYLE_SENDERS_STDOUT_TEMPLATE='{{ range $key, $value := .Datas }}
{{ $value.id }} => **{{ regexFind ".*?\n" $value.message | trim }}** *({{ $value.authorName }} - {{ $value.authorDate | date "2006-01-02 15:04:05" }})*
{{ end }}'

chyle create
```
output :

```markdown
f617fb708dfa6fa290205615ea98c53a860e499d => **Merge pull request #3 from antham/test2** *(Anthony HAMON - 2017-05-29 02:08:37)*

d8106fffee242f5b6394a103059b4064a83fcf3b => **Whatever** *(antham - 2017-05-29 02:08:37)*

e0a746c906fba7e2462f5717322b9eb55aca3943 => **Whatever** *(antham - 2017-05-29 02:08:37)*

118ad33a1d4ffc66bbeb74a1aba7524ef192ae62 => **Whatever** *(antham - 2017-05-29 02:08:37)*

78dcf412cc21d4054e06c534876200a89c04622e => **Whatever** *(antham - 2017-05-29 02:08:37)*

44fb3316ea67298df5a2b6fbb43795990575ec32 => **Whatever** *(antham - 2017-05-29 02:08:37)*

8fdfae00cbcc66936113a60f5146d110f2ba3c28 => **Merge pull request #1 from antham/test** *(Anthony HAMON - 2017-05-29 02:08:37)*
```

### Get a JSON ouput of all merge commits and contact github issue api to enrich payload

commands :

```bash
export CHYLE_GIT_REFERENCE_FROM=a00ee81c109c8787f0ea161a776d2c9795f816cd
export CHYLE_GIT_REFERENCE_TO=f617fb708dfa6fa290205615ea98c53a860e499d
export CHYLE_GIT_REPOSITORY_PATH=/your-local-path/test-git
export CHYLE_MATCHERS_TYPE=merge
export CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY=message
export CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY=githubIssueId
export CHYLE_EXTRACTORS_GITHUBISSUEID_REG="\#(\d+)"
export CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME=test-git
export CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN=token
export CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER=antham
export CHYLE_DECORATORS_GITHUBISSUE_KEYS_NUMBER_DESTKEY=ticketNumber
export CHYLE_DECORATORS_GITHUBISSUE_KEYS_NUMBER_FIELD=number
export CHYLE_DECORATORS_GITHUBISSUE_KEYS_COMMENTNUMBER_DESTKEY=commentNumber
export CHYLE_DECORATORS_GITHUBISSUE_KEYS_COMMENTNUMBER_FIELD=comments
export CHYLE_SENDERS_STDOUT_FORMAT="json"

chyle create
```
output :

```json
{
  "datas": [
    {
      "authorDate": "2017-05-10 22:24:40 +0200 +0200",
      "authorEmail": "antham@users.noreply.github.com",
      "authorName": "Anthony HAMON",
      "commentNumber": 0,
      "committerDate": "2017-05-10 22:24:40 +0200 +0200",
      "committerEmail": "noreply@github.com",
      "committerName": "GitHub",
      "githubIssueId": 3,
      "id": "f617fb708dfa6fa290205615ea98c53a860e499d",
      "message": "Merge pull request #3 from antham/test2\n\nTest2",
      "ticketNumber": 3,
      "type": "merge"
    },
    {
      "authorDate": "2017-05-10 22:22:03 +0200 +0200",
      "authorEmail": "antham@users.noreply.github.com",
      "authorName": "Anthony HAMON",
      "commentNumber": 0,
      "committerDate": "2017-05-10 22:22:03 +0200 +0200",
      "committerEmail": "noreply@github.com",
      "committerName": "GitHub",
      "githubIssueId": 1,
      "id": "8fdfae00cbcc66936113a60f5146d110f2ba3c28",
      "message": "Merge pull request #1 from antham/test\n\nTest",
      "ticketNumber": 1,
      "type": "merge"
    }
  ],
  "metadatas": {}
}
```

### Populate a release in github from CircleCI

Let's create a script release.sh

```bash
#!/usr/bin/env bash

set -e

### Functions

function setPreviousTag {
    if [ ! -n "$1" ]; then
        PREVIOUS_TAG="$(git describe --abbrev=0 --always --tags $2^)"
    fi

    if [[ $PREVIOUS_TAG =~ ^[0-9a-f]{40} ]];then
        PREVIOUS_TAG_SHA="$(git rev-list --max-parents=0 HEAD|head -n 1)"
        PREVIOUS_TAG="first commit"

        return
    fi

    PREVIOUS_TAG_SHA="$(git rev-parse $PREVIOUS_TAG^{commit})"
}

function setCurrentTag {
    local tag=""

    if [ -n "$1" ]; then
        tag="$1"
    fi

    if [ -n "$2" ]; then
        tag="$2"
    fi

    if [ ! -n "$tag" ]; then
        echo "You must declare CIRCLE_TAG or CURRENT_TAG variable"
        exit 1
    fi

    CURRENT_TAG_SHA="$(git rev-parse $tag^{commit})"
    CURRENT_TAG="$tag"
}

echo "-> Setup environment variables"

cd "$REPOSITORY_PATH"

setCurrentTag "$CURRENT_TAG" "$CIRCLE_TAG"
setPreviousTag "$PREVIOUS_TAG" "$CURRENT_TAG"

export CHYLE_GIT_REPOSITORY_PATH=$REPOSITORY_PATH
export CHYLE_GIT_REFERENCE_FROM=$PREVIOUS_TAG_SHA
export CHYLE_GIT_REFERENCE_TO=$CURRENT_TAG_SHA

# Setup matchers

## Pick only merge commits
export CHYLE_MATCHERS_TYPE=merge

# Setup decorators

## Setup github issue decorator

### Github issue id extractor
export CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY=message
export CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY=githubIssueId
export CHYLE_EXTRACTORS_GITHUBISSUEID_REG="\#(\d+)"

### Github credentials
export CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN=$GITHUB_TOKEN
export CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER=antham

### Git path
export CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME=$CIRCLE_PROJECT_REPONAME

### Extract title field
export CHYLE_DECORATORS_GITHUBISSUE_KEYS_TITLE_DESTKEY=title
export CHYLE_DECORATORS_GITHUBISSUE_KEYS_TITLE_FIELD=title

# Setup senders

## Setup github release

### Github credentials
export CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN=$GITHUB_TOKEN
export CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER=antham

### Github release config
export CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TAGNAME=$CIRCLE_TAG
export CHYLE_SENDERS_GITHUBRELEASE_RELEASE_UPDATE=true
export CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE='{{ range $key, $value := .Datas }}
{{ $value.id }} => **{{ $value.title | trim }}** *({{ $value.authorName }} - {{ $value.authorDate | date "2006-01-02 15:04:05" }})*
{{ end }}'

### Git path
export CHYLE_SENDERS_GITHUBRELEASE_REPOSITORY_NAME=$CIRCLE_PROJECT_REPONAME

# Generate changelog
echo "-> Generating changelog between $PREVIOUS_TAG ($PREVIOUS_TAG_SHA) and $CURRENT_TAG ($CURRENT_TAG_SHA)"
```

In CircleCI (v1) we add a deployment section to trigger a build when a tag is created (we need to declare credentials envs inside circle settings project) :

```yaml
machine:
  environment:
    REPOSITORY_PATH: $HOME/$CIRCLE_PROJECT_REPONAME

deployment:
  release:
    tag: /v[0-9]+(\.[0-9]+)*/
    owner: antham
    commands:
      - ./release.sh
```

Now we can create a release in github (we create a tag in the same time), it triggers a circle build and at the end, chyle generates a diff between two tags and populate the release, check [v1.0.0 release of test-git](https://github.com/antham/test-git/releases/tag/v1.0.0).

## Documentation

* [Default fields](#default-fields)
* [General config](#general-config)
* [Matchers](#matchers)
* [Extractors](#extractors)
* [Decorators](#decorators)
  * [Custom api decorator](#custom-api-decorator)
  * [Jira issue api](#jira-issue-api)
  * [Github issue api](#github-issue-api)
  * [Shell](#shell)
  * [Environment variable](#environment-variable)
* [Senders](#senders)
  * [Stdout](#stdout)
  * [Github release api](#github-release-api)
  * [Custom api sender](#custom-api-sender)
* [Help](#help)
  * [Template](#template)

---

### Default fields

Below you have an example, as a JSON, of a payload that shows fields extracted by default.

```json
{
  "datas": [
    {
      "authorDate": "2017-05-10 22:24:40 +0200 +0200",
      "authorEmail": "antham@users.noreply.github.com",
      "authorName": "Anthony HAMON",
      "committerDate": "2017-05-10 22:24:40 +0200 +0200",
      "committerEmail": "noreply@github.com",
      "committerName": "GitHub",
      "id": "f617fb708dfa6fa290205615ea98c53a860e499d",
      "message": "Merge pull request #3 from antham/test2\n\nTest2",
      "type": "merge"
    }
  ],
  "metadatas": {}
}
```

### General config

We need to define first, where the repository stand and which git range we want to target.

Name | Value
------------ | -------------
CHYLE_GIT_REPOSITORY_PATH | Path where your repository stand
CHYLE_GIT_REFERENCE_FROM | Git reference starting your range, could be an id, HEAD or a branch
CHYLE_GIT_REFERENCE_TO | Git reference ending your range, could be an id, HEAD or a branch

### Matchers
Matchers filters commits according to criterias.

Name | Value
------------ | -------------
CHYLE_MATCHERS_TYPE | Match commit by type, "merge" represents a merge commit, "regular" represents a usual commit (not merge)
CHYLE_MATCHERS_MESSAGE | A regexp that will be matched against a commit message
CHYLE_MATCHERS_COMMITTER | A regexp that will be matched against a committer field of a commit
CHYLE_MATCHERS_AUTHOR | A regexp that will be matched against an author field of a commit

### Extractors

Extractors defined from which commit field to extract datas, what to extract and under which name to store the extracted value, you can defined as many extractors you want.

You need to define those 3 values below in order to set an extractor, replace * with a name convenient to you.

Name | Value
------------ | -------------
CHYLE_EXTRACTORS_*_ORIGKEY | A commit field from which we want to extract datas (id, authorName, authorEmail, authorDate, committerName, committerEmail, committerMessage, type)
CHYLE_EXTRACTORS_*_DESTKEY | A name for the key which will receive the extracted value
CHYLE_EXTRACTORS_*_REG | A regexp used to extract a data

### Decorators

Decorators enrich your changelog with datas.

#### Custom api decorator

You can fetch a custom api to fetch any data.
First you need to define an id that will be added when calling your api.

Name | Value
------------ | -------------
CHYLE_EXTRACTORS_CUSTOMAPIID_ORIGKEY | Field from which you you want to extract the id
CHYLE_EXTRACTORS_CUSTOMAPIID_DESTKEY | customApiId
CHYLE_EXTRACTORS_CUSTOMAPIID_REG | A regexp to extract the id

You need to define a token header that will be given when calling your api and an url endpoint.

Name | Value
------------ | -------------
CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL | It's the endpoint of you remote api, use {{ID}} as a placeholder to interpolate the id you extracted before in URL
CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN | Token submitted as authorization header when calling your api

To extract data, you need to define everytime both a "DESTKEY" key and a "FIELD" key, replace * with a name convenient to you, you can get as many value as you want.

Name | Value
------------ | ------------
CHYLE_DECORATORS_CUSTOMAPI_KEYS_*_DESTKEY | A name for the key which will receive the extracted value
CHYLE_DECORATORS_CUSTOMAPI_KEYS_*_FIELD | The field to extract from your custom api response payload, use dot notation to extract a deep value (eg: "fields.summary")

#### Jira issue api

Have a look to the [api documentation](https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getIssue) to know what you can fetch from this api.
First, you need to use an extractor to define a "jiraIssueId" key to extract jira issue id, let's consider our id is in commit message we would add as environment variable.

Name | Value
------------ | -------------
CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY | Field from which you you want to extract the jira issue id
CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY | jiraIssueId
CHYLE_EXTRACTORS_JIRAISSUEID_REG | A regexp to extract jira issue id

You need to define jira credentials and endpoint.

Name | Value
------------ | -------------
CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL | It's the endpoint of you remote jira access point
CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME | Jira username
CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD | Jira password

To extract data, you need to define everytime both a "DESTKEY" key and a "FIELD" key, replace * with a name convenient to you, you can get as many value as you want.

Name | Value
------------ | ------------
CHYLE_DECORATORS_JIRAISSUE_KEYS_*_DESTKEY | A name for the key which will receive the extracted value
CHYLE_DECORATORS_JIRAISSUE_KEYS_*_FIELD | The field to extract from jira api response payload, use dot notation to extract a deep value (eg: "fields.summary")

#### Github issue api

You can get pull request datas or issue datas from this decorator as described in [api documentation](https://developer.github.com/v3/issues/#get-a-single-issue).
First, you need to use an extractor to define a "githubIssueId" key to extract github issue id, let's consider our id is in commit message we would add as environment variable.

Name | Value
------------ | -------------
CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY | Field from which you you want to extract the github issue id
CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY | githubIssueId
CHYLE_EXTRACTORS_GITHUBISSUEID_REG | A regexp to extract the github issue id

You need to define github credentials.

Name | Value
------------ | -------------
CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN | Github oauth token used to fetch issue datas
CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER | Github owner

To extract data, you need to define everytime both a "DESTKEY" key and a "FIELD" key, replace * with a name convenient to you, you can get as many value as you want.

Name | Value
------------ | ------------
CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_DESTKEY | A name for the key which will receive the extracted value
CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_FIELD | The field to extract from github issue api response payload, use dot notation to extract a deep value (eg: "fields.summary")

#### Shell

You can pipe any shell commands on every commit datas through this decorator.

To extract data, you need to define a threesome, a "COMMAND" key, an "ORIGKEY" key and a "DESTKEY" key, replace * with a name convenient to you, you can get as many value as you want.

Name | Value
------------ | ------------
CHYLE_DECORATORS_SHELL_*_COMMAND | Command to execute
CHYLE_DECORATORS_SHELL_*_ORIGKEY | A field from which you want to use the content to pipe a command on
CHYLE_DECORATORS_SHELL_*_DESTKEY | A name for the key which will receive the extracted value


#### Environment variable

This decorator dump an environment variable in metadatas changelog section.

You need to define everytime both an "VALUE" key and a "DESTKEY" key, replace * with a name convenient to you, you can get as many value as you want.

Name | Value
------------ | -------------
CHYLE_DECORATORS_ENV_*_VARNAME | Environment variable name to dump in metadatas
CHYLE_DECORATORS_ENV_*_DESTKEY | The name of the key where to store dumped value in metadatas

### Senders

Senders are called when all operations are done on payload to render final result.

#### Stdout

Dump result to stdout

Name | Value
------------ | -------------
CHYLE_SENDERS_STDOUT_FORMAT | "json" : output payload as JSON , "template" : output payload using golang template syntax look at the [help](#template) to get more informations
CHYLE_SENDERS_STDOUT_TEMPLATE | Linked to "template" stdout format, it must be set to defined a template following golang template syntax, look at the [help](#template) to get more informations

#### Github release api

It creates a new release in github with a template from datas you harvested, look at the [documentation](https://developer.github.com/v3/repos/releases/#create-a-release) to have more details about what you can do.

Name | Value
------------ | -------------
CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN | Github oauth token used to publish a release (mandatory)
CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER | Github owner (mandatory)
CHYLE_SENDERS_GITHUBRELEASE_RELEASE_DRAFT | Create a draft (unpublished) release, boolean value, default is false
CHYLE_SENDERS_GITHUBRELEASE_RELEASE_NAME | Release title
CHYLE_SENDERS_GITHUBRELEASE_RELEASE_PRERELEASE | Create a prerelease release, boolean value, default is false
CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TAGNAME | Release tag to create, when you update a release it will be used to find out release tied to this tag (mandatory)
CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TARGETCOMMITISH | The commitish value that determines where the Git tag is created from
CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE | It uses golang template syntax to produce a changelog from your commits (mandatory), look at the [help](#template) to get more informations
CHYLE_SENDERS_GITHUBRELEASE_RELEASE_UPDATE | Set to true if you want to update an existing changelog, typical usage would be when you produce a release through GUI github release system
CHYLE_SENDERS_GITHUBRELEASE_REPOSITORY_NAME | Github repository where we will publish the release (mandatory)

#### Custom api sender

Send release to a custom http endpoint through a POST request as a JSON payload.

Name | Value
------------ | -------------
CHYLE_SENDERS_CUSTOMAPI_CREDENTIALS_TOKEN | Access token given in request header "Authorization"
CHYLE_SENDERS_CUSTOMAPI_ENDPOINT_URL | The URL endpoint where the POST request will be made

### Help

#### Template

Chyle uses go template as template engine, documentation can be found in godoc [here](https://golang.org/pkg/text/template/#hdr-Text_and_spaces).

Let's have an example using the following release generated using JSON format :

```json
{
  "datas": [
    {
      "authorDate": "2017-05-10 22:24:40 +0200 +0200",
      "authorEmail": "antham@users.noreply.github.com",
      "authorName": "Anthony HAMON",
      "committerDate": "2017-05-10 22:24:40 +0200 +0200",
      "committerEmail": "noreply@github.com",
      "committerName": "GitHub",
      "date": "2017-05-10 22:24:40",
      "githubIssueId": 3,
      "githubTitle": "Test2",
      "id": "f617fb708dfa6fa290205615ea98c53a860e499d",
      "message": "Merge pull request #3 from antham/test2\n\nTest2",
      "type": "merge"
    },
    {
      "authorDate": "2017-05-10 22:22:03 +0200 +0200",
      "authorEmail": "antham@users.noreply.github.com",
      "authorName": "Anthony HAMON",
      "committerDate": "2017-05-10 22:22:03 +0200 +0200",
      "committerEmail": "noreply@github.com",
      "committerName": "GitHub",
      "date": "2017-05-10 22:22:03",
      "githubIssueId": 1,
      "githubTitle": "Whatever",
      "id": "8fdfae00cbcc66936113a60f5146d110f2ba3c28",
      "message": "Merge pull request #1 from antham/test\n\nTest",
      "type": "merge"
    }
  ],
  "metadatas": {
    "date": "jeu. mai 18 23:01:25 CEST 2017"
  }
}
```

If we want to display a markdown release with the list of pull request title and their authors we can do :

```go
### Release

{{ range $key, $value := .Datas }} ({{ $value.authorName }}) {{ end }}

Generated at {{ .Metadatas.date }}
```

We get :

```markdown
### Release

Test2 (Anthony HAMON)
Whatever (Anthony HAMON)
Generated at jeu. mai 18 23:01:25 CEST 2017%
```

To provide more functionalities to original golang template, [sprig](https://github.com/Masterminds/sprig) library is provided, it gives several useful additional helpers, documentation can be found [here](http://masterminds.github.io/sprig/).

For the sake of convenience, a custom global store is available as well, as templates cannot mutate defined variables : you can store a data using ```{{ set "key" "data"}}```, you can retrieve a data using ```{{ get "key" }}```, you can test if a key is set using ```{{ isset "key" }}```.


## Contribute

If you want to add a new feature to chyle project, the best way is to open a ticket first to know exactly how to implement your changes in code.

### Setup

After cloning the repository you need to install vendors with [dep](https://github.com/golang/dep).
To test your changes locally you can run go tests with : ```make run-quick-tests```, and you can run gometalinter check with : ```make gometalinter```, with those two commands you will fix lot of issues, other tests will be ran through travis so only open a pull request to see what break.
