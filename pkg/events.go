package events

import "github.com/keptn/go-utils/pkg/lib/v0_2_0"

const EchoEventTriggeredType = "sh.keptn.event.echo.triggered"

const EchoStartedEventType = "sh.keptn.event.echo.started"

const EchoFinishedEventType = "sh.keptn.event.echo.finished"

const ServiceName = "echo-service"

type EchoTriggeredEventData struct {
	v0_2_0.EventData
}

type EchoStartedEventData struct {
	v0_2_0.EventData
}

type EchoFinishedEventData struct {
	v0_2_0.EventData
}
