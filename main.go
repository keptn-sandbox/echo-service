package main

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/keptn-sandbox/echo-service/cmd/eventhandling"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"log"
	"net/url"
	"os"
	"time"
)

var keptnOptions = keptncommon.KeptnOpts{}

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int `envconfig:"RCV_PORT" default:"8080"`
	// Path to which cloudevents are sent
	Path string `envconfig:"RCV_PATH" default:"/"`
	// Whether we are running locally (e.g., for testing) or on production
	Env string `envconfig:"ENV" default:"local"`
	// URL of the Keptn configuration service (this is where we can fetch files from the config repo)
	ConfigurationServiceUrl string `envconfig:"CONFIGURATION_SERVICE" default:""`
	// URL of the Keptn event broker (this is where this service sends cloudevents to)
	EventBrokerUrl string `envconfig:"EVENTBROKER" default:""`
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

	processor := eventhandling.EchoCloudEventProcessor{
		EventSender: eventSender,
		Sleeper:     eventhandling.NewConfigurableSleeper(5*time.Second, time.Sleep),
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
