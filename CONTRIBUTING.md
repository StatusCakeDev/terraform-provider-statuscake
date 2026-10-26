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
- [Go](https://go.dev/) (1.27+)
- [Terraform](https://www.terraform.io/)
- [Make](https://www.gnu.org/software/make/) (4.0+)

## Installation

If you wish to work on the provider, you'll first need
[Go](http://www.golang.org) installed on your machine (the version named in
`go.mod`). You'll also need to correctly setup a
[GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding
`$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider that
should then be moved to the `$GOPATH/bin` directory.

```sh
git clone <repository-url>` this repository
cd terraform-provider-statuscake
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

You will need a StatusCake API token to try the provider against the real API.
Include it within terraform configuration, or set `STATUSCAKE_API_TOKEN`.

### Pointing the provider at another API

`custom_endpoint` within terraform configuration, or
`STATUSCAKE_CUSTOM_ENDPOINT` exists only in development builds, made with `make
build TAGS=dev`. A release build has no such setting and ignores it if set. Most
tests do not need the tag: they hand `adapters.Config` an `Endpoint`. Those that
drive the provider through its config file are built with `-tags dev`.

## Running Tests

- `make test` vets and tests both the release and the development build.

The tests run the provider through the Terraform CLI, as a user would, so they
create no real resources and need no API token. Where they need an API they talk
to a fake of it (`internal/provider/fakeapi_test.go`). Install
[Terraform](https://developer.hashicorp.com/terraform/install), or
[OpenTofu](https://opentofu.org) by pointing `TF_ACC_TERRAFORM_PATH` at it (and,
for OpenTofu, setting `TF_ACC_PROVIDER_HOST=registry.opentofu.org` and
`TF_ACC_PROVIDER_NAMESPACE=hashicorp`). Without one the tests that need the CLI
are skipped.

The tests that need the fake run only in a development build, as that is the one
that can be pointed at it (see below), so `make test` runs the tests twice: once
as a release and once with `-tags dev`. Run `go test -tags dev ./...` to run all
of them.

The tests drive the provider from the outside, and only the server that `main`
serves is in the way of it. If a behaviour cannot be exercised that way, prefer
a test that can over exporting something to reach it.

Two tests hold golden files, because users' configurations and states depend on
what they pin:

- `TestSchemaIsUnchanged` compares the schema that Terraform receives with
  `testdata/schema.golden.json`.
- `TestStateIsUnchanged` compares the state of a minimal resource of each kind
  with `testdata/state/`. It holds every default and every value read back.

If you change either deliberately, review the difference and refresh the files:

```sh
go test ./internal/provider -run TestSchemaIsUnchanged -update
go test -tags dev ./internal/provider -run TestStateIsUnchanged -update
```

A new resource should have a lifecycle test like those in
`internal/provider/resource_*_test.go`, a case in `TestStateIsUnchanged`, and
the fake API needs a response for its kind.

## Making Changes

For additional contributing guidelines visit
[devhandbook.org](https://devhandbook.org/contributing)
