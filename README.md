# xcopy

xcopy is a Go (golang) CLI with implementations of `ioprogress` and `bwio` that copying files and folders with

* Bandwidth control
* Display copying progress with human readable units.

## Usage

Here is an example of outputting a basic usage of this CLI

```
$ ./xcp.exe --help
Usage of ~\github.com\thimoonxy\xcopy\xcp.exe:
  -d string
        dst path (default "~")
  -h    Print human readable output .
  -l int
        Bandwitdh limts, e.g. 50 means 50MB/s
  -p    Print progress status .
  -s string
        src path (default "~")
  -v    Print Verbose.

```

## Example

![Progress](https://github.com/thimoonxy/xcopy/blob/master/1.gif)

![Progress](https://github.com/thimoonxy/xcopy/blob/master/2.gif)
## Requirements

This CLI is designed to work on both Linux & Windows environments.
