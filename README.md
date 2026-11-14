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
are signed with the signing key `ED23BC86151A016D322936BCC8F1DBBD3BA98A96`.

### Public Key

```text
-----BEGIN PGP PUBLIC KEY BLOCK-----

mQINBGUf+F8BEADDuiMPSLplWsqhCbeZRNEKup16XfeXVlvOcVXw1l4yIxeR9p+1
o9Z9PDju2quwQI61GblwygyO6mWEcM06blwbEIMtMF+m99w6bUy9FJTMGKX720rH
fEaVbjU/vB/EFzFnkBpzbTKF7frnbVk9oI6I2VVKScDgCLIBv35qWIT1MVLx5qXh
hS3pyW944JmJLE6u7qh3MC27zZYYDeVtUPKe9TK+xbZLCqV+8QXC3kI1eelF2LxX
52np8xswgRugfpY4H5CoLnbPfwEc7frCgz3enQ+B7vt3gXQZzRE7J+TWr9yENaq5
IDKyM6ASoU0AFIQ5000TtcQpVnktV3V9w4PlQgcGxLMyeoC5xim7yJ8oYUG4HaFM
7D23V1KtHy9SO8Ou7w9+izt5BHYHZ4ELCVCUWBzRjeifbT8s1R2Wi4FPAE1Lph9j
Tiu2ENPZjXNc1taQsASIcCTnXldx21EeC2PhhJAIjmrkCXhT5AVl4ynviF38jnFY
xojYPLgVAbzoiqRpd6LWr5TGM4rCcnfu7Q2RUBGjxWO2vHQLAMp+4deTwwMkdVET
d2U0P05XVyfvohyR+P7pTByXG62myH0T4jxtJQJAxHmWNp+TumpiY2z4M420KQ8O
/XvnjhgZj3KHJF9S1Asa1szBHxXMvkwytQ3Qw/i4G/nE2jBm8PIF2XT4FQARAQAB
tCdTdGF0dXNDYWtlIDxlbmdpbmVlcmluZ0BzdGF0dXNjYWtlLmNvbT6JAlIEEwEI
ADwWIQTtI7yGFRoBbTIpNrzI8du9O6mKlgUCZR/4XwIbAwULCQgHAgMiAgEGFQoJ
CAsCBBYCAwECHgcCF4AACgkQyPHbvTupipYrNQ//SD5N1w67wCT7jdpbHhcNmcfs
TNRkH1+Gmt7ZO6WX3Pd4bKUw90t+1ZTLdH/1VKICm1HivjapVjlM7Q2U0DKZQd0m
6905lEWsDnkgkqi1x3Q121qh+uOq8rG190L812GanYaS+DwPWOmZvBqtiU3SFSlL
FoaYH/HIOUEPg5Z190KlmYJv3CtyF5TjpmXr7Z+I70Gwc4g37QbGA+2OA+kzQGFz
H5YzA9EZAx/DBcmz1u+JG7ZpD4OQyt+BRHA4MXCkCcYOF0ShW/pxoxNuAEUJKWEF
o7dDlkJ9TIIkktn6csi5JoXUGBtNnjVVeqa5X7aLbu9p8Xlz0/ts7hmWnW6CV2uo
Os9RJ2couPGD/KB5tohK9obPmhIlgDFTwnSzYMQA9fTKwldPc8ycyRC0SjjlIovI
5/M+Gz5JMjuXq0/HGkHs1R1oxZm4rkIvZbgSnu599cne3PGg2KXfhb68xcGHdz9h
XWegI2dvSCHfD/Z4UJf5xSlgmvhg0D6oVuGm8WykqUmPQiAPjScUUuztXiBIbOH/
rEPfaInCUDJIz9otp9Dg3MCpNtC+8d+OJ2Z+7JmndatnFRP29cQdWRe0/BOmBVWr
LdgvB6GEmkc3Paz+RpLj7kpPyR4cV7q+muJoIf73w4xGU4dQnTHlHzYSvoMwZ6Kt
ROuAgFRNVX/zqVBbQHS5Ag0EZR/4XwEQAMGaIoz+tVxWctW1v2C0KxWYoiwOc7d2
QeExxIuQI2jtGr3CitA44ZwpCOdGhfYZOvaHQPhIrtQiRUx5B9jn4ID29AxnMi62
GlH6yjWYRKRXnzbUaE/6x6Ac5luEXFRJqZP7/PQXMoOFrnyM1oWsEsKd0pR0l98K
2locy7lPn5HVVJ7RvZDaqb0i6iBSRSR8biehE4nXEvDcY+Tpay1CJTcvWCUzs0DC
6yrlcxrCgCm5jcQEW3v0E3RWU/Vj3K2QOfbbTjPjFEZA48qTv0LhJKfQqSF/uAkZ
v4C9s4aR3U9fQRKyDD/mTKj2zwidkyCOpIDsPZZ7MsZidpSwysQZPnUvelMu2W+y
UBQRmEiMw+RT8tQUlL0bmQo1CqbYl7m6ECkGHNKDV3uM8VeUaUo+ZyVsvI9LaY2x
nWx+7ntcDf/zeAXmjhkwzXrSahOdwqpoY8hzzASf4QalcG3TR5kQfLpMJf/81P77
ZGvx4yVspKKP5ovx4gx7ed7pogW+1tVzBzZ02NGAS5ZZW4IizoxgWeaFlPmXORk7
jze0PYQxT/Q3f2HSwYoxEZv54fCum3VwS6QgCla7XrTk4dfSzdCadrR5jzzOql0C
kB4x9OzZoet2UL+/ManKKAXixgcqI894OMjjRkTH1pyo6Klx5YZj63Z7clM3cpnw
abC7+vQDdkYJABEBAAGJAjYEGAEIACAWIQTtI7yGFRoBbTIpNrzI8du9O6mKlgUC
ZR/4XwIbDAAKCRDI8du9O6mKlmBhD/4u5nMSQXZke08DTbhAcREQlBRryX7AtJK+
rlpVlc1cxDHr8u2aJJzHfyV73JVxFXKir1YcAQLNaWoN7P19LKpNMWhjuNKO0EoB
vvPxXJ0scax5ZiGQ9zWffA08X0laG/qj/rB+XttrPgEQncZ4kKGcArbyE8/dvKk4
x4SEZOgN/NfKCuVIFf9cGj333G8mwNJx/VqLDZHkuZUJR+UBKtZw+PLCsa4+ZjVX
qcV75MYcXgOYlgdxwoRbBjZ1r1tFsDdzYWHXWsgv8tW2h5GB6Ov+vEvEUgHNSXfH
pqptB7RUGZWi5QZU6o3nUiZ7bLY6baoagRfZcZzRSx8qkMGbQWnSSV7DkHPzfi2K
n/j/TODV05S50+Gt9YhbayYdcfhx7kohSH1n8QD7IRmujSKdvim8HtBCdt+jY/EV
sqOZ+99xf0ivjTq5YMnUJsqWZfOTNGsRp41V/IhUEhaS2Yi8cdh5wezJdQHA7y+e
1DVtL66Y1jW4F4CDG3gGXBgReqGFCSxKhYe/iQqzw/Y+eRb/4uMh9qnJ9QvBdwlg
yP5DRZJfRpXi5jAuqSdTsqsmhFIcO51jgRwQ5/5uhZhgqRzq7z4y7Qcaq8ecIyWi
Bl4NpkTMuCoAmImotXb/JerRAqMdlTP/kwcMNcpX15UMl2yDtVTdZwDS8WVgkfYN
0luRfRPVyA==
=eDwk
-----END PGP PUBLIC KEY BLOCK-----
```

### Verifying Binaries

1. Import the signing key:

   ```bash
   gpg --keyserver keyserver.ubuntu.com --recv-keys ED23BC86151A016D322936BCC8F1DBBD3BA98A96
   ```

1. Verify that the key has been imported:

   ```bash
   gpg --list-keys engineering@statuscake.com
   ```

   Check the output is:

   ```bash
   pub   rsa4096 2023-10-06 [SC]
         ED23BC86151A016D322936BCC8F1DBBD3BA98A96
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
   gpg:                using RSA key ED23BC86151A016D322936BCC8F1DBBD3BA98A96
   gpg:                issuer "engineering@statuscake.com"
   gpg: Good signature from "StatusCake (Code signing for statuscake.com) <engineering@statuscake.com>" [ultimate]
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
