# StatusCake Terraform Provider [![test](https://github.com/StatusCakeDev/terraform-provider-statuscake/actions/workflows/test.yaml/badge.svg?event=push)](https://github.com/StatusCakeDev/terraform-provider-statuscake/actions/workflows/test.yaml)

<!-- pyml disable md033 -->
<a href="https://terraform.io">
  <img src=".github/terraform_logo.svg" alt="Terraform logo" title="Terraform"
    align="right" height="50" />
</a>
<!-- pyml enable md033 -->

[![Support][support-badge]][support]

[support-badge]: https://img.shields.io/badge/support-terraform--statuscake-623CE4.svg?style=flat
[support]: https://developers.statuscake.com/guides/terraform/introduction

- Website: [terraform.io](https://terraform.io)
- Tutorials: [learn.hashicorp.com](https://learn.hashicorp.com/terraform?track=getting-started#getting-started)
- Support: [developers.statuscake.com](https://developers.statuscake.com/guides/terraform/introduction)
- Discussion: [Discuss](https://discuss.hashicorp.com/c/terraform-core/27)

## Prerequisites

You will need the following things properly installed on your computer:

- [Terraform](https://www.terraform.io/downloads.html) (1.0+)

## Installation

Grab the binary for your system from the [releases
page](https://github.com/StatusCakeDev/terraform-provider-statuscake/releases).

Alternatively, install with Go:

```bash
go install github.com/StatusCakeDev/terraform-provider-statuscake@latest
```

> [!NOTE]
> Installing from source using `go install` will build from the master branch
> and won't report version information via `terraform-provider-statuscake --version`.

### Documentation

Full documentation is available on the [Terraform
website](https://registry.terraform.io/providers/StatusCakeDev/statuscake/latest/docs)
and the [developer portal][support].

## Signing

The binaries on the [releases
page](https://github.com/StatusCakeDev/terraform-provider-statuscake/releases)
are signed with the signing key `DDEA0314FAE6420CA54CB80D9F8705905758331F`.

### Public Key

```text
-----BEGIN PGP PUBLIC KEY BLOCK-----

mDMEauKIPxYJKwYBBAHaRw8BAQdATGYPDtbEk25PRHeOojl6bn11f8wP9Gm9NT+H
BwBY6B+0RVN0YXR1c0Nha2UgKENvZGUgc2lnbmluZyBmb3Igc3RhdHVzY2FrZS5j
b20pIDxzaWduaW5nQHN0YXR1c2Nha2UuY29tPoiZBBMWCgBBFiEE3eoDFPrmQgyl
TLgNn4cFkFdYMx8FAmriiD8CGwMFCQX6nYAFCwkIBwICIgIGFQoJCAsCBBYCAwEC
HgcCF4AACgkQn4cFkFdYMx+8iQEAr5lnryZE09h3MHhK+O1jqk83ThL4AJLcgCSG
FETIl5YA/j2YAUtP7TfnY+xP0g3UVagT9LY2j83G2S/TlRDfhpMIuDgEauKIPxIK
KwYBBAGXVQEFAQEHQDKjCcdzlwsDYLlAAGszcondKuVtWfv3QufPl0hen6RTAwEI
B4h+BBgWCgAmFiEE3eoDFPrmQgylTLgNn4cFkFdYMx8FAmriiD8CGwwFCQX6nYAA
CgkQn4cFkFdYMx+tmgEAgC+E6x61nxHsmenUcxsXA+Zia07NWSsX0xkBab+oPN8A
/j8UfoJyozzbxq6gUNrTKi1sFDZtGeevB540wkakm9IL
=Knl0
-----END PGP PUBLIC KEY BLOCK-----
```

### Verifying Binaries

1. Import the signing key:

   ```bash
   gpg --keyserver keyserver.ubuntu.com --recv-keys DDEA0314FAE6420CA54CB80D9F8705905758331F
   ```

1. Verify that the key has been imported:

   ```bash
   gpg --list-keys signing@statuscake.com
   ```

   Check the output is:

   ```bash
   pub   ed25519 2026-10-28 [SC] [expires: 2030-01-01]
         DDEA0314FAE6420CA54CB80D9F8705905758331F
   ```

1. Download the binary and its corresponding shasums and signature from the
   [releases
   page](https://github.com/StatusCakeDev/terraform-provider-statuscake/releases).

1. Verify the signature of the shasums file:

   ```bash
   gpg --verify terraform-provider-statuscake_x.y.z_SHA256SUMS.sig terraform-provider-statuscake_x.y.z_SHA256SUMS
   ```

   This should give you a similar output to that below:

<!-- pyml disable md013 -->
   ```bash
   gpg: Signature made Sun Feb  8 15:30:24 2026 GMT
   gpg:                using EDDSA key DDEA0314FAE6420CA54CB80D9F8705905758331F
   gpg:                issuer "signing@statuscake.com"
   gpg: Good signature from "StatusCake (Code signing for statuscake.com) <signing@statuscake.com>" [ultimate]
   ```
<!-- pyml enable md013 -->

1. Verify the shasum of the binary matches the value in the shasums file:

   ```bash
   sha256sum -c --ignore-missing terraform-provider-statuscake_x.y.z_SHA256SUMS
   ```

   This should give you a similar output to that below:

   ```bash
   terraform-provider-statuscake_x.y.z_linux_amd64.zip: OK
   ```

## License

This project is licensed under the [Mozilla Public License](LICENSE).
