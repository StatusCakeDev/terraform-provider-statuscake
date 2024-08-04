# How To Contribute

Contributions are welcome as we strive to make this application as useful as
possible for everyone. However time is not always on our side, and changes may
not be reviewed or merged in a timely manner.

If this application is found to be missing in functionality, please open an
issue describing the proposed change - discussing changes ahead of time reduces
friction within pull requests.

## Prerequisites

You will need the following things properly installed on your computer:

- [Git](https://git-scm.com/)
- [Go](https://go.dev/) (1.21+)
- [Terraform](https://www.terraform.io/)
- [Make](https://en.wikipedia.org/wiki/Make_(software)) (optional)

## Installation

If you wish to work on the provider, you'll first need
[Go](http://www.golang.org) installed on your machine (version 1.19+ is
*required*). You'll also need to correctly setup a
[GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding
`$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider that
should then be moved to the `$GOPATH/bin` directory.

```sh
make build
mv terraform-provider-statuscake $GOPATH/bin/terraform-provider-statuscake
```

To use the compiled binary the following must be included in the
`~/.terraformrc` file, having replaced `FULL_PATH_TO_GO_BIN` with the full
directory path to the `$GOPATH/bin` directory. This informs the Terraform CLI
tool to lookup the binary in the `$GOPATH` instead of the regular location.

```terraformrc
provider_installation {
  dev_overrides {
    "statuscakedev/statuscake" = "FULL_PATH_TO_GO_BIN"
  }
  direct {}
}
```

## Running tests

- `make testacc`

## Making Changes

For additional contributing guidelines visit
[devhandbook.org](https://devhandbook.org/contributing)
