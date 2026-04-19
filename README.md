# Power Consumption Metrics Tool

Queries HomeAssistant for a summary of your power usage over a period of time.
Home Assistant provides power usage data for each hour on its Energy dashboard, but does not have an API endpoint to query this data.
The websocket API does provide a way to do this, which is what the frontend uses.

I got fed up of trying to figure out how to get the same data that the Energy dashboard shows, so I wrote this tool to do it for me.

This tool queries the websocket API to get the power usage data for each hour over a period of days and outputs the data in various formats

## Installation

Go

```bash
go install github.com/poolski/powertracker@latest
```

[Releases](https://github.com/poolski/powertracker/releases)

## Configuration

This tool requires a configuration file to be present at `~/.config/powertracker/config.yaml`. If one does not exist, it will ask for input and create it for you.
The only things this tool needs are the URL of your Home Assistant instance and a long-lived access token.

You can generate a long-lived access token by going to your Home Assistant instance, clicking on your profile picture in the bottom left, then clicking on "Long-Lived Access Tokens" at the bottom of the list and creating a new one.

## Usage

```bash
$ powertracker --help

Usage:
  powertracker [flags]

Flags:
  -c, --config string     config file (default "$HOME_DIR/.config/powertracker/config.yaml")
  -f, --csv-file string   the path of the CSV file to write to (default "results.csv")
  -d, --days int          number of days to compute power stats for (default 7)
  -h, --help              help for powertracker
  -i  --insecure          skip TLS verification
  -o, --output string     output format (text, table, csv)

```

> **Personal note:** I changed the default `--days` value from 30 to 7, since I mostly want a weekly overview and found myself always passing `-d 7` anyway.

> **Personal note:** I use the `table` output format almost exclusively, so I set that as my default output format instead of `text`.

> **Personal note:** My Home Assistant instance uses a self-signed certificate, so I always need `--insecure`. Consider wrapping this in a shell alias: `alias pt='powertracker --insecure'`

> **Personal note:** For CSV exports I keep the output in `~/Documents/power-data/` rather than the default `results.csv` in the working directory. Useful to have a consistent location when running from cron: `pt -o csv -f ~/Documents/power-data/weekly.csv`

## Example output

```bash
$ powertracker -d 7 # 7 days' worth of data

+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+
|    0     |    1     |    2     |    3     |    4     |    5     |    6     |    7     |    8     |    9     |    10    |    11    |    12    |    13    |    14    |    15    |    16    |    17    |    18    |    19    |    20    |    21    |    22    |    23    |
+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+-------
```
