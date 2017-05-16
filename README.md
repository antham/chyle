Chyle [![Build Status](https://travis-ci.org/antham/chyle.svg?branch=master)](https://travis-ci.org/antham/chyle) [![codecov](https://codecov.io/gh/antham/chyle/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/chyle) [![codebeat badge](https://codebeat.co/badges/c3867610-2741-4ae3-a195-d5e9711c7fcd)](https://codebeat.co/projects/github-com-antham-chyle-master) [![Go Report Card](https://goreportcard.com/badge/github.com/antham/chyle)](https://goreportcard.com/report/github.com/antham/chyle)
=====

Chyle produces a changelog from a git repository.

## Usage

```
Create a new changelog according to what is defined as config.

Changelog creation follows this process :

1 - fetch commits
2 - filter relevant commits
3 - extract informations from fetched datas
4 - contact third part services to retrieve additional informations from extracted datas
5 - send result to third part services

Usage:
  chyle create [flags]

Global Flags:
      --debug   enable debugging
```

## How it works ?

Chyle fetch a range of commits using given criterias from a git repository. From those commits you can extract relevant datas like commit message, author, and so on, and add it to original payload. We can afterwards, contact an external apis to enrich our payload with various useful datas. Finally, we can publish what we harvested to an external api for instance (currently only github release is added) and/or stdout.

## Setup

Download from release page according to your architecture chyle binary : https://github.com/antham/chyle/releases

You need afterwards to configure each module through environments variables : there are activated when you configure at least one environment variable they need to work.

### Summary

* [General config](#general-config)
* [Matchers](#matchers)
* [Extractors](#extractors)
* [Decorators](#decorators)
  * [Jira issue api](#jira-issue-api)
  * [Github issue api](#github-issue-api)
  * [Shell](#shell)
  * [Environment variable](#environment-variable)
* [Senders](#senders)
  * [Stdout](#stdout)
  * [Github release api](#github-release-api)

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

#### Jira issue api

Have a look to the [api documentation](https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getIssue) to know what you can fetch from this api.
First, you need to use an extractor to define a "jiraIssueId" key to extract jira issue id, let's consider our id is in commit message we would add as environment variable.

Name | Value
------------ | -------------
CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY | message
CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY | jiraIssueId
CHYLE_EXTRACTORS_JIRAISSUEID_REG | "(\w+-\d+)"

You need to define jira credentials and endpoint.

Name | Value
------------ | -------------
CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_URL | It's the endpoint of you remote jira access point
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
CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY | message
CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY | githubIssueId
CHYLE_EXTRACTORS_GITHUBISSUEID_REG | "(\w+-\d+)"

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
CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_COMMAND | Command to execute
CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_ORIGKEY | A field from which you want to use the content to pipe a command on
CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_DESTKEY | A name for the key which will receive the extracted value


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
CHYLE_SENDERS_STDOUT_FORMAT | "json" : output payload as JSON , "template" : output payload using golang template syntax
CHYLE_SENDERS_STDOUT_TEMPLATE | Linked to "template" stdout format, it must be set to defined a template following golang template syntax

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
CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE | It uses golang template syntax to produce a changelog from your commits (mandatory), eg: A metadata {{.Metadatas.title}} {{ range $key, $value := .Datas }}{{$value.authorName}} : {{$value.message}}{{ end }}
CHYLE_SENDERS_GITHUBRELEASE_RELEASE_UPDATE | Set to true if you want to update an existing changelog, typical usage would be when you produce a release through GUI github release system
CHYLE_SENDERS_GITHUBRELEASE_REPOSITORY_NAME | Github repository where we will publish the release (mandatory)
