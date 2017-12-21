package main

import (
	"os"
	"io"
	"bufio"
	"fmt"
	"strings"
	"flag"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

var baseTopic, inFile, mqttHost string
var reader io.Reader

func setup() bool {
	flag.StringVar( &baseTopic, "base-topic", "", "Base MQTT topic for pub/sub" )
	flag.StringVar( &mqttHost, "host", "localhost", "MQTT server host")

	flag.Parse()

	inFile = flag.Arg(0)

	if ( inFile == "" || inFile == "-" ) {
		reader = os.Stdin
	} else {
		var err error
		reader, err = os.Open( inFile )
		if err != nil {
			panic(err)
		}
	}

	return true

}

func main() {

	if ! setup() {
		os.Exit(1)
	}

	fmt.Printf( "baseTopic: %#v\n", baseTopic )
	fmt.Printf( "outHost: %#v\n", mqttHost )
	fmt.Printf( "file: %#v\n", inFile )

	inLine := make( chan string )

	//define a function for the default message handler
	var f MQTT.MessageHandler = func(client *MQTT.Client, msg MQTT.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}

	opts := MQTT.NewClientOptions().AddBroker("tcp://" + mqttHost + ":1883")
	opts.SetClientID("goRedNinja")
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe( baseTopic + "/write", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	go func ( line chan string, cli *MQTT.Client, baseTopic string ) {
		for {
			t := c.Publish( baseTopic + "/read", 0, false, <-line)
			t.Wait()
		}
	}( inLine, c, baseTopic )


	var line string
	stdin := bufio.NewReader( reader )

	for {
		line, _ = stdin.ReadString( '\n' )
		line = strings.Trim( line, "\n" )

		if len(line) == 0 { continue }
		killString := "\"G\":\"0\",\"V\":0,\"D\":2"

		if ( ! strings.Contains( line, killString ) ) {
			inLine <- line
		}
	}

}