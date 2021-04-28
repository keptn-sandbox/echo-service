package main

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/keptn-sandbox/echo-service/cmd/eventhandling"
	"log"
	"net/url"
	"os"
	"time"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int `envconfig:"RCV_PORT" default:"8080"`
	// Path to which cloudevents are sent
	Path string `envconfig:"RCV_PATH" default:"/"`
	// URL of the Keptn event broker (this is where this service sends cloudevents to)
	EventBrokerURL string `envconfig:"EVENTBROKER" default:""`
	// Duration in milliseconds the echo service will sleep between sending a start and finished event
	SleepTimeMillis int `envconfig:"SLEEP_TIME_MS" default:"1000"`
	// Flag indicating whether the service shall send the .started and .finished events in correct or reversed order
	SimulateEventOutOfOrder bool `envconfig:"SIMULATE_EVENTS_OUT_OF_ORDER" default:"false"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {

	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)

	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	svcEndpoint, _ := getServiceEndpoint("EVENTBROKER")
	eventSender := eventhandling.NewHTTPEventSender(svcEndpoint)

	var processor eventhandling.CloudEventProcessor
	if env.SimulateEventOutOfOrder {
		log.Println("Using Broken cloud event processor that will send the .started and .finished events out of order")
		processor = eventhandling.BrokenEchoCloudEventProcessor{
			EventSender: eventSender,
			Sleeper:     eventhandling.NewConfigurableSleeper(time.Duration(env.SleepTimeMillis)*time.Millisecond, time.Sleep),
		}
	} else {
		processor = eventhandling.EchoCloudEventProcessor{
			EventSender: eventSender,
			Sleeper:     eventhandling.NewConfigurableSleeper(time.Duration(env.SleepTimeMillis)*time.Millisecond, time.Sleep),
		}
	}

	log.Fatal(c.StartReceiver(ctx, processor.Process))

	return 0
}

func getServiceEndpoint(service string) (url.URL, error) {
	url, err := url.Parse(os.Getenv(service))
	if err != nil {
		return *url, fmt.Errorf("Failed to retrieve value from ENVIRONMENT_VARIABLE: %s", service)
	}

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	return *url, nil
}
