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
	"net/url"
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

	if event.Type() == events.EchoEventTriggeredType {
		log.Println("Processing Echo Triggered Event")

		// 1. send started event
		if err := sendStartEvent(shkeptncontext, event); err != nil {
			return err
		}

		// 2. process event
		eventData := &events.EchoTriggeredEventData{}
		err := event.DataAs(eventData)
		if err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}
		if err := eventhandling.HandleEchoEvent(eventData, log.Writer(), eventhandling.NewConfigurableSleeper(5*time.Second, time.Sleep)); err != nil {
			return err
		}

		// 3. send finish event
		if err := sendFinishEvent(shkeptncontext, event, keptnv2.ResultPass); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func sendStartEvent(shkeptnctx string, incomingEvent cloudevents.Event) error {

	echoStartedEventData := events.EchoStartedEventData{}
	echoTriggeredEventData := events.EchoTriggeredEventData{}

	if err := incomingEvent.DataAs(&echoTriggeredEventData); err != nil {
		return err
	}

	echoStartedEventData.Status = keptnv2.StatusSucceeded
	echoStartedEventData.EventData = echoTriggeredEventData.EventData
	event := cloudevents.NewEvent()
	event.SetType(events.EchoStartedEventType)
	event.SetSource(events.ServiceName)
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", shkeptnctx)
	event.SetExtension("triggeredid", incomingEvent.ID())
	event.SetData(cloudevents.ApplicationJSON, echoStartedEventData)

	return sendEvent(event)

}

func sendFinishEvent(shkeptnctx string, incomingEvent cloudevents.Event, result keptnv2.ResultType) error {
	echoFinishedEventData := events.EchoFinishedEventData{}
	echoTriggeredEventData := events.EchoTriggeredEventData{}

	if err := incomingEvent.DataAs(&echoTriggeredEventData); err != nil {
		return err
	}

	echoFinishedEventData.Result = result
	echoFinishedEventData.Status = keptnv2.StatusSucceeded
	echoFinishedEventData.EventData = echoTriggeredEventData.EventData
	event := cloudevents.NewEvent()
	event.SetType(events.EchoFinishedEventType)
	event.SetSource(events.ServiceName)
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", shkeptnctx)
	event.SetExtension("triggeredid", incomingEvent.ID())
	event.SetData(cloudevents.ApplicationJSON, echoFinishedEventData)

	return sendEvent(event)

}

func sendEvent(event cloudevents.Event) error {

	endPoint, err := getServiceEndpoint("EVENTBROKER")
	if err != nil {
		return errors.New("Failed to retrieve endpoint of eventbroker. %s" + err.Error())
	}

	if endPoint.Host == "" {
		return errors.New("Host of eventbroker not set")
	}
	keptnHandler, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{
		EventBrokerURL: endPoint.String(),
	})

	if err != nil {
		return errors.New("Failed to initialize Keptn handler: " + err.Error())
	}

	return keptnHandler.SendCloudEvent(event)
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
