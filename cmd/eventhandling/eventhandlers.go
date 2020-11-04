package eventhandling

import (
	"fmt"
	events "github.com/keptn-sandbox/echo-service/pkg"
	"io"
)

func HandleEchoEvent(data *events.EchoTriggeredEventData, writer io.Writer, sleeper Sleeper) error {
	_, err := fmt.Fprintf(writer, "GOT ECHO TRIGGERED EVENT\n")
	sleeper.Sleep()
	return err
}
