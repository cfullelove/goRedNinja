# goRedNinja

## Description

goRedNinja reads in lines from a file (usually a serial port, or stdin) and publishes these lines onto a
MQTT topic

## Usage
```
$ ./goRedNinja --host hostname --base-topic topic filename
```
Where:
```
Usage:
--host hostname      :   Hostname of MQTT Broker [default: localhost]
--base-topic topic   :   Base MQTT topic [default: ""]
filename             :   Filename to read in messages [e.g. /dev/ttyO1]["-" for stdin]
```

## How it works
goRedNinja reads lines from file and will publish each line on the `base-topic` with `/read` appended e.g. `base-topic/read`

goRedNinja will also subscribe to the topic `base-topic` with `write` appended (e.g. `base-topic/write`), however at this time messages received are not written to file but rather logged to the console

