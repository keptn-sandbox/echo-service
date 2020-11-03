package events

const EchoEventType = "sh.keptn.event.echo"

type EchoEventData struct {
	//Message is the message to echo
	Message string
}
