Chyle [![Build Status](https://travis-ci.org/antham/chyle.svg?branch=master)](https://travis-ci.org/antham/chyle) [![codecov](https://codecov.io/gh/antham/chyle/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/chyle) [![codebeat badge](https://codebeat.co/badges/1fd5d776-6145-4a3f-9705-731d14e7283e)](https://codebeat.co/projects/github-com-antham-chyle)
=====

Chyle produces a changelog from a git repository.

## Usage

```
Create a new changelog according to what is defined in config file.

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

Chyle fetch a range of commits using given criterias from a git repository. From those commits you can extract relevant datas like commit message, author, and so on, and add it to original payload. We can afterwards, contact an external apis to enrich our payload with various useful datas (currently only jira ticket api is added). Finally, we can publish what we harvested to an external api for instance (currently only github release is added) and/or stdout.

## Setup

Download from release page according to your architecture chyle binary : https://github.com/antham/chyle/releases

You need afterwards to configure each module through environments variables : their are activated when you configure at least one environment variable they need to work.

### Summary

* [General config](#general-config)
* [Matchers](#matchers)
* [Extractors](#extractors)
* [Decorators](#decorators)
  * [Jira ticket api](#jira-ticket-api)
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

Decorators request remote api to enrich your commit payload with datas.

#### Jira ticket api

First, you need to use an extractor to define a "jiraIssueId" key ot extract jira ticket id, let's consider our id is in commit message we would add as environment variable.

Name | Value
------------ | -------------
CHYLE_EXTRACTORS_JIRA_ORIGKEY | message
CHYLE_EXTRACTORS_JIRA_DESTKEY | jiraIssueId
CHYLE_EXTRACTORS_JIRA_REG | "(\w+-\d+)"

You need to define everytime both a "DESTKEY" key and a "FIELD" key, replace * with a name convenient to you, you can get as many value as you want.

Name | Value
------------ | -------------
CHYLE_DECORATORS_JIRA_CREDENTIALS_URL | It's the endpoint of you remote jira access point
CHYLE_DECORATORS_JIRA_CREDENTIALS_USERNAME | Jira username
CHYLE_DECORATORS_JIRA_CREDENTIALS_PASSWORD | Jira password
CHYLE_DECORATORS_JIRA_KEYS_*_DESTKEY | A name for the key which will receive the extracted value
CHYLE_DECORATORS_JIRA_KEYS_*_FIELD | The field to extract from jira api response payload, use dot notation to extract a deep value (eg: "fields.summary")

### Senders

Senders are called when all operations are done on payload to render final result.

#### Stdout

Dump result to stdout

Name | Value
------------ | -------------
CHYLE_SENDERS_STDOUT_FORMAT | Only json is supported at the moment

#### Github release api

It creates a new release in [github](https://developer.github.com/v3/repos/releases/#create-a-release) with a template from datas you harvested.

Name | Value
------------ | -------------
CHYLE_SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN | Github oauth token used to publish a release
CHYLE_SENDERS_GITHUB_CREDENTIALS_OWNER | Github owner
CHYLE_SENDERS_GITHUB_RELEASE_DRAFT | Create a draft (unpublished) release, boolean value, default is false
CHYLE_SENDERS_GITHUB_RELEASE_NAME | Release title (optional)
CHYLE_SENDERS_GITHUB_RELEASE_PRERELEASE | Create a prerelease release, boolean value, default is false
CHYLE_SENDERS_GITHUB_RELEASE_TAGNAME | Release tag to create, when you update a release it will be used to find out release tied to this tag
CHYLE_SENDERS_GITHUB_RELEASE_TARGETCOMMITISH | The commitish value that determines where the Git tag is created from (optional)
CHYLE_SENDERS_GITHUB_RELEASE_TEMPLATE | It uses golang template syntax to produce a changelog from you commits, eg: {{ range $key, $value := . }}{{$value.authorName}} : {{$value.message}}{{ end }}
CHYLE_SENDERS_GITHUB_RELEASE_UPDATE | Set to true if you want to update an existing changelog, typical usage would be when you produce a release through GUI github release system
CHYLE_SENDERS_GITHUB_REPOSITORY_NAME | Github repository where we will publish the release
