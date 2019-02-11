# lic
Lic is an easily extensible & flexible report generator to statically analyse your local sources and create a report on the fly or upload said report to a server.

## Usage
```
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
```
Usage:
  lic report [command]

Aliases:
  report, r

Available Commands:
  golang      Generates a report of current working directory

Flags:
  -h, --help   help for report

Global Flags:
  -v, --verbose   verbose output
```

## Build the Software
Make sure to run `go get` to install missing dependencies and then just run `make` to get an executable into the /bin/ folder.

## Contribute
This tool should support multiple languages. Currently I'm working on golang. Feel free to chip in there or contribute a new language set (I'd be happy to see Java, JavaScript/Typescript).
The goal of this tool for me to get a reliable list of used software components on any sources I throw at it, to be able to check licenses against whatever licenses I want to use/don't want to use.

The implementation will most likely change over time and is as of now a moving spec. The report might need to have more information to it than I currently foresee.

Also I want to connect this to some sort of CVE search, to be able to tell if I have any vulnerabilities in my code. This doesn't have to be on API level and I can see this as a stretch goal to build a server-side component or re-use an existing tool and format the output the right way, so that result visualization over time, storing of scan results and updates on findings as they become known are a possibility.
