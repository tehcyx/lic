# lic
Lic is an easily extensible & flexible report generator to statically analyse your local sources and create a report on the fly or upload said report to a server.

- [How to build lic](#how-to-build-lic)
- [How to use lic](#how-to-use-lic)
- [Roadmap](#roadmap)
- [Contribute](#contribute)

## How to build lic
Make sure to run `go get` to install missing dependencies and then just run `make` to get an executable into the /bin/ folder.
If you want to use the executable from your path, you can invoke `make install` to install the file in the specified `$GOBIN` location.

## How to use lic
To generate your first report invoke `lic report golang` in the project folder of your choice or specify an absolute path like `lic report golang --src $GOPATH/src/github.com/username/repository` to run the report in a different folder than the current working directory.

The executable supports various other commands that you can see via the help sub-command (`lic help`):

```shell
Usage:
  lic [command]

Available Commands:
  help        Help about any command
  report      Creates a report of sources
  version     Version of the lic CLI

Flags:
  -h, --help      help for lic
  -v, --verbose   verbose output
```

Various commands will have sub commands, for example the report command will differentiate with the supported languages (currently only golang).
```shell
Usage:
  lic report [command]

Aliases:
  report, r

Available Commands:
  golang      Generates a report of current working directory or specified path

Flags:
  -h, --help   help for report

Global Flags:
  -v, --verbose   verbose output
```

## Roadmap
- Extend language support
  - Java
  - JavaScript
  - Typescript
  - ...?
- Report generation
  - richer reports (HTML, JSON)
- Version detection
- Server-side component that receives reports, holds history

## Requirements

To not run into early rate-limiting of GitHub's API, the tool requires the `LIC_GITHUB_ACCESS_TOKEN` environment variable to be set with the value of a GitHub personal access token, which can be created here: [Get a personal access token](https://github.com/settings/tokens). The token does not need to have any checkboxes applied and runs with a standard no-permission token.

The `git` executable is assumed to be on the PATH, so that the tool can get the repositories tags via this command: `git describe --tags --always`. Future implementations should safeguard this by checking if the tool exists. Current implementation would probably result in an error.

## Contribute
This tool should support multiple languages. Currently I'm working on golang. Feel free to chip in or contribute a new language set (I'd be happy to see Java/JavaScript/Typescript). Tests are rare so far, so there's definitely more needed.

The goal of this tool for me to get a reliable list of used software components on any sources I throw at it, to be able to check licenses against whatever licenses I want to use/don't want to use.

The implementation will most likely change over time and is as of now a moving spec. The report might need to have more information to it than I currently foresee.

Also I want to connect this to some sort of CVE search, to be able to tell if I have any vulnerabilities in my code. This doesn't have to be on API level and I can see this as a stretch goal to build a server-side component or re-use an existing tool and format the output the right way, so that result visualization over time, storing of scan results and updates on findings as they become known are a possibility.
