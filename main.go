package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn-sandbox/echo-service/cmd/eventhandling"
	events "github.com/keptn-sandbox/echo-service/pkg"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"log"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
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

func processKeptnCloudEvent(ctx context.Context, event cloudevents.Event) error {

	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	_, err := keptnv2.NewKeptn(&event, keptnOptions)
	log.Printf("gotEvent(%s): %s - %s", event.Type(), shkeptncontext, event.Context.GetID())

	if err != nil {
		log.Printf("failed to parse incoming cloudevent: %v", err)
		return err
	}

	if event.Type() == events.EchoEventType {
		log.Printf("Processing Echo Event")

		eventData := &events.EchoEventData{}
		err := event.DataAs(eventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}

		return eventhandling.HandleEchoEvent(eventData, log.Writer(), eventhandling.NewConfigurableSleeper(5*time.Second, time.Sleep))
	}
	// Unknown Event -> Throw Error!
	var errorMsg string
	errorMsg = fmt.Sprintf("Unhandled Keptn Cloud Event: %s", event.Type())

	log.Print(errorMsg)
	return errors.New(errorMsg)
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	os.Exit(_main(os.Args[1:], env))
}

/**
 * Opens up a listener on localhost:port/path and passes incoming requets to gotEvent
 */
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
	log.Fatal(c.StartReceiver(ctx, processKeptnCloudEvent))
	return 0
}
