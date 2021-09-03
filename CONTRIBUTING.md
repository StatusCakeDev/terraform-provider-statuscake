# How To Contribute

Contributions are welcome as we strive to make this application as useful as
possible for everyone. However time is not always on our side, and changes may
not be reviewed or merged in a timely manner.

If this application is found to be missing in functionality, please open an
issue describing the proposed change - discussing changes ahead of time reduces
friction within pull requests.

## Installation

If you wish to work on the provider, you'll first need
[Go](http://www.golang.org) installed on your machine (version 1.17+ is
*required*). You'll also need to correctly setup a
[GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding
`$GOPATH/bin` to your `$PATH`.

To compile the provider, run `go build`. This will build the provider and
should be moved to the `$GOPATH/bin` directory.

```sh
$ go build -o terraform-provider-statuscake
$ mv terraform-provider-statuscake $GOPATH/bin/terraform-provider-statuscake
...
```

To use the compiled binary the following must be included in the
`~/.terraformrc` file, having replaced `FULL_PATH_TO_GO_BIN` with the full
directory path to the `$GOPATH/bin` directory. This informs the Terraform CLI
tool to lookup the binary in the `$GOPATH` instead of the regular location.

```
provider_installation {
  dev_overrides {
    "statuscakedev/statuscake" = "FULL_PATH_TO_GO_BIN"
  }
  direct {}
}
```

## Linting

* `golint ./...`

## Running tests

* `make testacc`

## Making Changes

Begin by creating a new branch. It is appreciated if branch names are written
using kebab-case.

```bash
git checkout master
git pull --rebase
git checkout -b my-new-feature
```

Make the desired change, and ensure both the linter and test suite continue to
pass. Once this requirement is met push the change back to a fork of this
repository.

```bash
git push -u origin my-new-feature
```

Finally open a pull request through the GitHub UI. Upon doing this the CI suite
will be run to ensure changes do not break current functionality.

Changes are more likely to be approve if they:

- Include tests for new functionality,
- Are accompanied with a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html),
- Contain few commits (preferably a single commit),
- Do not contain merge commits,
- Maintain backward compatibility.
