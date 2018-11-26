<p align="center">
  <img alt="logo" height="120" src="https://github.com/subfinder/research/blob/master/logo.png"/>
<p>

# SubZero
> ❄️  Research project for the development of [SubFinder](https://github.com/subfinder/subfinder) core API `v2`

# Installing SubZero
If you're interested in testing out the alpha version of SubFinder's next command-line application experience:
> **Note**: this is alpha software. It's not perfect, and will be eventually merged into SubFinder's main project. This is purely for those interested in early testing to help out with development.
```console
$ go get -u github.com/subfinder/research/subzero
```

## Command-Line Usage
You can `enumerate` a domain, like `google.com`, and optionally use the `--verbose` flag for more information.
```console
$ subzero enumerate google.com --verbose
```

Get help for any command to learn about more options.
```console
$ subzero help enumerate
```

Generic (top-level) help menu.
```console
$ subzero help
```

### Help Menu

```console
$ subzero --help
Usage:
  subzero [command]

Available Commands:
  enumerate   Enumerate subdomains for the given domains.
  help        Help about any command

Flags:
  -h, --help   help for subzero

Use "subzero [command] --help" for more information about a command.
```

```console
$ subzero enumerate --help
Enumerate subdomains for the given domains

Usage:
  subzero enumerate [domains to enumerate] [flags]

Flags:
  -h, --help          help for enumerate
      --insecure      include potentially insecure sources using http
      --labels        show source of the domain in output
      --limit int     limit the reported results to the given number
      --no-timeout    do not timeout
      --recursive     use results to find more results
      --timeout int   number of seconds until timeout (default 30)
      --uniq          filter uniq results
      --verbose       show errors and other available diagnostic information
```

#### Run Tests
```console
$ cd /path/to/research
$ GOCACHE=off go test -v -run=. ./... -race
# just core
$ GOCACHE=off go test -v -run=. ./core -race
```

# Why?
That's a great question, and I hope these reasons will make sense:

* Reduce boilerplate code needed to add or develop new sudomain `Sources`
* Utilize the [`interface`](https://gobyexample.com/interfaces) type in golang to define our internal API semantics.
* Provide tests and examples of all its functionality.
* Ensure a solid foundation for future tools to be built upon (CLI, HTTP/REST API)

# How can I help?
Make an issue or pull request to this repository to improve it!

* Try running the tests.
* Fix the tests.
* Write some tests.
* Ask questions about parts that confuse you.
* Add documentation.
* Provide solutions for better word choices in documentation.
* Spellchecking is awesome.
* Be creative!
