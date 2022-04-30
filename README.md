# kabuka

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go test](https://github.com/shionit/kabuka/actions/workflows/go.yml/badge.svg)](https://github.com/shionit/kabuka/actions/workflows/go.yml)
[![reviewdog](https://github.com/shionit/kabuka/workflows/reviewdog/badge.svg)](https://github.com/shionit/kabuka/actions?query=workflow%3Areviewdog)

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

kabuka supports json format.
```shell
$ kabuka 3994.T -f=json
{"symbol":"3994.T","current_price":"8420"}

## same as
$ kabuka 3994.T --format=json
```
and also supports csv format.
```shell
$ kabuka 3994.T -f=csv
symbol,current_price
3994.T,5850

## same as
$ kabuka 3994.T --format=csv
```

## Dependency

kabuka command depends on Yahoo! Japan Finance website.
It fetch stock price info from the website.
https://finance.yahoo.co.jp/
