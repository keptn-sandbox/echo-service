package eventhandling

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	events "github.com/keptn-sandbox/echo-service/pkg"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"log"
)

// CloudEventProcessor defines a processor on cloud events
type CloudEventProcessor interface {
	Process(event cloudevents.Event) error
}

// EchoCloudEventProcessor is the default implementation of a CloudEventProcessor that hhah
type EchoCloudEventProcessor struct {
	EventSender EventSender
	Sleeper     Sleeper
}

// Process processes a cloud event
func (ep EchoCloudEventProcessor) Process(event cloudevents.Event) error {

	if event.Type() == events.EchoEventTriggeredType {
		log.Println("Processing Echo Triggered Event")
		eventData := &events.EchoTriggeredEventData{}
		if err := event.DataAs(eventData); err != nil {
			log.Printf("Got Data Error: %s", err.Error())
			return err
		}

		//TODO: make this prettier
		if eventData.SimulateWrongEventSeq {
			if err := ep.EventSender.SendEvent(ep.createEchoFinishedEvent(event)); err != nil {
				log.Printf("Got Send Error: %s", err.Error())
				return err
			}

		} else {
			if err := ep.EventSender.SendEvent(ep.createEchoStartedEvent(event)); err != nil {
				log.Printf("Got Send Error: %s", err.Error())
				return err
			}
		}

		log.Printf("GOT EVENT :)\n")

		ep.Sleeper.Sleep()

		if eventData.SimulateWrongEventSeq {
			if err := ep.EventSender.SendEvent(ep.createEchoStartedEvent(event)); err != nil {
				log.Printf("Got Send Error: %s", err.Error())
				return err
			}

		} else {
			if err := ep.EventSender.SendEvent(ep.createEchoFinishedEvent(event)); err != nil {
				log.Printf("Got Send Error: %s", err.Error())
				return err
			}
		}
	}
	return nil
}

func (ep EchoCloudEventProcessor) createEchoStartedEvent(incomingEvent cloudevents.Event) event.Event {
	var shkeptnctx string
	incomingEvent.Context.ExtensionAs("shkeptncontext", &shkeptnctx)

	echoStartedEventData := events.EchoStartedEventData{}
	echoTriggeredEventData := events.EchoTriggeredEventData{}

	if err := incomingEvent.DataAs(&echoTriggeredEventData); err != nil {
		log.Println(err.Error())
		return event.Event{}
	}

	echoStartedEventData.Status = keptnv2.StatusSucceeded
	echoStartedEventData.EventData = echoTriggeredEventData.EventData
	outEvent := cloudevents.NewEvent()
	outEvent.SetType(events.EchoStartedEventType)
	outEvent.SetSource(events.ServiceName)
	outEvent.SetDataContentType(cloudevents.ApplicationJSON)
	outEvent.SetExtension("shkeptncontext", shkeptnctx)
	outEvent.SetExtension("triggeredid", incomingEvent.ID())
	outEvent.SetData(cloudevents.ApplicationJSON, echoStartedEventData)
	return outEvent
}

func (ep EchoCloudEventProcessor) createEchoFinishedEvent(incomingEvent cloudevents.Event) event.Event {

	var shkeptnctx string
	incomingEvent.Context.ExtensionAs("shkeptncontext", &shkeptnctx)

	echoFinishedEventData := events.EchoFinishedEventData{}
	echoTriggeredEventData := events.EchoTriggeredEventData{}

	if err := incomingEvent.DataAs(&echoTriggeredEventData); err != nil {
		log.Println(err.Error())
		return event.Event{}
	}

	echoFinishedEventData.Result = keptnv2.ResultPass
	echoFinishedEventData.Status = keptnv2.StatusSucceeded
	echoFinishedEventData.EventData = echoTriggeredEventData.EventData
	outEvent := cloudevents.NewEvent()
	outEvent.SetType(events.EchoFinishedEventType)
	outEvent.SetSource(events.ServiceName)
	outEvent.SetDataContentType(cloudevents.ApplicationJSON)
	outEvent.SetExtension("shkeptncontext", shkeptnctx)
	outEvent.SetExtension("triggeredid", incomingEvent.ID())
	outEvent.SetData(cloudevents.ApplicationJSON, echoFinishedEventData)

	return outEvent

}
