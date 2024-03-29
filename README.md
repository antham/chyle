<p align="left"><img src="logo/horizontalversion.png" alt="chyle" height="100px"></p>

# [![codecov](https://codecov.io/gh/antham/chyle/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/chyle) [![Go Report Card](https://goreportcard.com/badge/github.com/antham/chyle)](https://goreportcard.com/report/github.com/antham/chyle) [![GitHub tag](https://img.shields.io/github/tag/antham/chyle.svg)]()

Chyle produces a changelog from a git repository.

[![asciicast](https://asciinema.org/a/o2PDZ4ELfUP3F1eKWl1IqirzU.png)](https://asciinema.org/a/o2PDZ4ELfUP3F1eKWl1IqirzU)

---

- [Usage](#usage)
- [How it works ?](#how-it-works-)
- [Setup](#setup)
- [Documentation and examples](#documentation-and-examples)
- [Contribute](#contribute)

---

## Usage

```
Create a changelog from your commit history

Usage:
  chyle [command]

Available Commands:
  config      Configuration prompt
  create      Create a new changelog
  help        Help about any command

Flags:
      --debug   enable debugging
  -h, --help    help for chyle

Use "chyle [command] --help" for more information about a command.
```

### config

Run a serie of prompt to help generate quickly and easily a configuration.

### create

Generate changelog.

## How it works ?

Chyle fetch a range of commits using given criterias from a git repository. From those commits you can extract relevant datas from commit message, author, and so on, and add it to original payload. You can afterwards if needed, enrich your payload with various useful datas contacting external apps (shell command, apis, ....) and finally, you can publish what you harvested (to an external api, stdout, ....). You can mix all steps together, avoid some, combine some, it's up to you.

## Setup

Download from release page according to your architecture chyle binary : https://github.com/antham/chyle/releases

Look at the documentation and examples, run `chyle config` to launch the configuration prompt.

## Documentation and examples

Have a look to the [wiki of this project](https://github.com/antham/chyle/wiki).

## Contribute

If you want to add a new feature to chyle project, the best way is to open a ticket first to know exactly how to implement your changes in code.

### Setup

After cloning the repository you need to install vendors with `go mod vendor`
To test your changes locally you can run go tests with : `make test-all`
