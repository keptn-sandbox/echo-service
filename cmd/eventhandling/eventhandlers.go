package eventhandling

import (
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	events "github.com/keptn-sandbox/echo-service/pkg"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"io"
)

func HandleEchoEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *events.EchoEventData, writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "GOT A MESSAGE: %s\n", data.Message)
	return err
}
