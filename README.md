# xcopy

When copied / backup files from/to server file system to/from mounted container, VMs, NFS and or ridiculous CIFS, you've ever experienced certain IO issue.
xcopy is a Go (golang) CLI with implementations of `ioprogress` and `bwio` that copying files and folders with

* IO Bandwidth control
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

Display progress and limts bw to 10 MB/s:
![Progress](https://github.com/thimoonxy/xcopy/blob/master/1.gif)

Display progress in human readable format, task stats shows IO is shaped to 10 MB/s as expected :
![Progress](https://github.com/thimoonxy/xcopy/blob/master/2.gif)

## Environment

* This CLI is designed to work on both `Linux` & `Windows` environments.
* Golang1.7 or newer.
