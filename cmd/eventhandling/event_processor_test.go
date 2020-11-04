package eventhandling

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	events "github.com/keptn-sandbox/echo-service/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type FakeEventSender struct {
	sentEvents []cloudevents.Event
}

func (fs *FakeEventSender) SendEvent(event cloudevents.Event) error {
	fs.sentEvents = append(fs.sentEvents, event)
	return nil
}

func TestEventsGetsSent(t *testing.T) {

	fakeEventSender := &FakeEventSender{}
	fakeSleeper := &TestSleeper{}
	eventProcessor := EchoCloudEventProcessor{
		EventSender: fakeEventSender,
		Sleeper:     fakeSleeper,
	}

	event := cloudevents.NewEvent()
	event.SetType(events.EchoEventTriggeredType)
	event.SetData(cloudevents.ApplicationJSON, events.EchoTriggeredEventData{SimulateWrongEventSeq: false})

	err := eventProcessor.Process(event)

	require.Nil(t, err)
	require.Equal(t, 2, len(fakeEventSender.sentEvents))
	assert.Equal(t, events.EchoStartedEventType, fakeEventSender.sentEvents[0].Type())
	assert.Equal(t, events.EchoFinishedEventType, fakeEventSender.sentEvents[1].Type())

}

func TestEventsGetsSentInWrongOrder(t *testing.T) {

	fakeEventSender := &FakeEventSender{}
	fakeSleeper := &TestSleeper{}
	eventProcessor := EchoCloudEventProcessor{
		EventSender: fakeEventSender,
		Sleeper:     fakeSleeper,
	}

	event := cloudevents.NewEvent()
	event.SetType(events.EchoEventTriggeredType)
	event.SetData(cloudevents.ApplicationJSON, events.EchoTriggeredEventData{SimulateWrongEventSeq: true})

	err := eventProcessor.Process(event)

	require.Nil(t, err)
	require.Equal(t, 2, len(fakeEventSender.sentEvents))
	assert.Equal(t, events.EchoFinishedEventType, fakeEventSender.sentEvents[0].Type())
	assert.Equal(t, events.EchoStartedEventType, fakeEventSender.sentEvents[1].Type())

}
