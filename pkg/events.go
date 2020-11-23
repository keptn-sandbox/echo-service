package events

import "github.com/keptn/go-utils/pkg/lib/v0_2_0"

// EchoEventTriggeredType is the name of an echo triggered event
const EchoEventTriggeredType = "sh.keptn.event.echo.triggered"

// EchoStartedEventType is the name of an echo started event
const EchoStartedEventType = "sh.keptn.event.echo.started"

// EchoFinishedEventType is the name of an echo finished event
const EchoFinishedEventType = "sh.keptn.event.echo.finished"

// ServiceName is the name of this service
const ServiceName = "echo-service"

// EchoTriggeredEventData is the data of an echo triggered event
type EchoTriggeredEventData struct {
	v0_2_0.EventData
}

// EchoStartedEventData is the data of an echo started event
type EchoStartedEventData struct {
	v0_2_0.EventData
}

// EchoFinishedEventData is the data of an echo finished event
type EchoFinishedEventData struct {
	v0_2_0.EventData
}
