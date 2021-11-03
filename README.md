# kabuka

[![Go](https://github.com/shionit/kabuka/actions/workflows/go.yml/badge.svg)](https://github.com/shionit/kabuka/actions/workflows/go.yml)

kabuka is command line tool for display stock price.

## Usage

Show Japanese market stock price.
```shell
$ kabuka 3994.T
8090	3994.T

$ kabuka 4373
2782	4373.T
```

Also show US market stock price.
```shell
$ kabuka AAPL
150.02	AAPL
```

## Dependency

kabuka command depends on Yahoo! Japan Finance website.
It fetch stock price info from the website.
https://finance.yahoo.co.jp/
