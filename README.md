# gocovdedup - Go Cover Deduplicate

Small utility application to merge and deduplicate go test cover files.

This program deduplicates over lapping profiles by unioning the line coverage.

TThe purpose of this tool is to get accurate line coverage.  Statement coverage uses the maximum of  overlapping lines value from each profile, and as such may not be strictly accurate.

## Install

```sh
go install github.com/nehemming/gocovdedup@latest
```

## Usage

The tool can either be used to pipe profiles to via stdin or be supplied paths to each coverage file to include.  The deduplicated and merged output is always sent to stdout.

### Piped

```sh
cat files | gocovdedup - > cover.out
```

### Args

```sh
gocovdedup package_one.out package_tow.out commontests.out > cover.out
```

