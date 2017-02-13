Chyle [![Build Status](https://travis-ci.org/antham/chyle.svg?branch=master)](https://travis-ci.org/antham/chyle) [![codecov](https://codecov.io/gh/antham/chyle/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/chyle) [![codebeat badge](https://codebeat.co/badges/1fd5d776-6145-4a3f-9705-731d14e7283e)](https://codebeat.co/projects/github-com-antham-chyle)
=====

Chyle produces a changelog from a git repository.


## How it works ?

Chyle fetch a range of commits using given criterias from a git repository. From those commits you can extract relevant datas from commit message, author, etc... , those datas are added to original datas fetched from commits. We can afterwards contact a external apis to enrich our payload with various useful datas (currently only jira is added). Finally, we can publish what we harvested to an external api for instance (currently only github release is added) or stdout.

## Setup

Download from release page according to your architecture chyle binary : https://github.com/antham/chyle/releases

You need afterwards to configure each module through environments variables, their are activated when you configure environments variables they need to work.

### Matchers

Matchers filters commits according to criterias.

Name | Value
------------ | -------------
CHYLE_MATCHERS_NUMPARENTS | A number of parents a commit have, 1 will be a regular commit, 2 will be a merge cmmit
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

### Expanders

Expanders request remote api to enrich your commit payload with datas.

#### Jira ticket api

You need to define everytime both a "DESTKEY" key and a "FIELD" key, replace * with a name convenient to you, you can get as many value as you want.

Name | Value
------------ | -------------
CHYLE_EXPANDERS_JIRA_CREDENTIALS_URL | It's the endpoint of you remote jira access point
CHYLE_EXPANDERS_JIRA_CREDENTIALS_USERNAME | Jira username
CHYLE_EXPANDERS_JIRA_CREDENTIALS_PASSWORD | Jira password
CHYLE_EXPANDERS_JIRA_KEYS_*_DESTKEY | A name for the key which will receive the extracted value
CHYLE_EXPANDERS_JIRA_KEYS_*_FIELD | The field to extract from jira api response payload, use dot notation to extract a deep value (eg: "fields.summary")

### Senders

Senders are called when all operations are done on payload to render final result.

#### Stdout

Dump result to stdout

Name | Value
------------ | -------------
CHYLE_SENDERS_STDOUT_FORMAT | Only json is supported at the moment

#### Github release api

It creates a new release in github with a template from datas you harvested, if tag exists already it will fail.

Name | Value
------------ | -------------
CHYLE_SENDERS_GITHUB_TEMPLATE | It uses golang template syntax to produce a changelog from you commits
CHYLE_SENDERS_GITHUB_TAG | Release tag to create
CHYLE_SENDERS_GITHUB_CREDENTIALS_OWNER | Github owner
CHYLE_SENDERS_GITHUB_REPOSITORY_NAME | Github repository where we will publish the release
CHYLE_SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN | Github oauth token used to publish a release
