package eventhandling

import (
	"context"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/protocol"
	httpprotocol "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"log"
	"net/url"
	"time"
)

const maxSendRetries = 3

// EventSender defines a cloud event sender
type EventSender interface {

	// SendEvent sends a cloud event
	SendEvent(event cloudevents.Event) error
}

// HTTPEventSender is a HTTP based implementation of a cloud event sender
type HTTPEventSender struct {
	svcEndpoint url.URL
	client      client.Client
}

// NewHTTPEventSender creates a new instance of a HTTP based cloud event sender
func NewHTTPEventSender(svcEndpoint url.URL) EventSender {

	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	return &HTTPEventSender{
		svcEndpoint: svcEndpoint,
		client:      c,
	}
}

// SendEvent sends a cloud event
func (es HTTPEventSender) SendEvent(event cloudevents.Event) error {

	ctx := cloudevents.ContextWithTarget(context.Background(), es.svcEndpoint.String())
	ctx = cloudevents.WithEncodingStructured(ctx)

	var result protocol.Result
	for i := 0; i <= maxSendRetries; i++ {
		result = es.client.Send(ctx, event)
		httpResult, ok := result.(*httpprotocol.Result)
		if ok {
			if httpResult.StatusCode >= 200 && httpResult.StatusCode < 300 {
				return nil
			}
			<-time.After(keptn.GetExpBackoffTime(i + 1))
		} else if cloudevents.IsUndelivered(result) {
			<-time.After(keptn.GetExpBackoffTime(i + 1))
		} else {
			return nil
		}
	}
	return errors.New("Failed to send cloudevent: " + result.Error())
}
