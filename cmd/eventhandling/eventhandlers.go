package eventhandling

import (
	"fmt"
	events "github.com/keptn-sandbox/echo-service/pkg"
	"io"
)

func HandleEchoEvent(data *events.EchoEventData, writer io.Writer, sleeper Sleeper) error {
	_, err := fmt.Fprintf(writer, "GOT A MESSAGE: %s\n", data.Message)
	sleeper.Sleep()
	return err
}
