package eventhandling

import (
	"encoding/json"
	"fmt"
	events "github.com/keptn-sandbox/echo-service/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"

	_ "github.com/keptn/go-utils/pkg/lib"
	keptn "github.com/keptn/go-utils/pkg/lib"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
)

/**
 * loads a cloud event from the passed test json file and initializes a keptn object with it
 */
func initializeTestObjects(eventFileName string) (*keptn.Keptn, *cloudevents.Event, error) {
	// load sample event
	eventFile, err := ioutil.ReadFile(eventFileName)
	if err != nil {
		return nil, nil, fmt.Errorf("Cant load %s: %s", eventFileName, err.Error())
	}

	incomingEvent := &cloudevents.Event{}
	err = json.Unmarshal(eventFile, incomingEvent)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing: %s", err.Error())
	}

	var keptnOptions = keptn.KeptnOpts{}
	keptnOptions.UseLocalFileSystem = true
	myKeptn, err := keptn.NewKeptn(incomingEvent, keptnOptions)

	return myKeptn, incomingEvent, err
}

func TestHandleEchoEvent(t *testing.T) {
	testLogger := testWriter{}
	keptn, incomingEvent, err := initializeTestObjects("../../test-events/echo-event.json")
	require.Nil(t, err)

	data := &events.EchoEventData{}
	err = incomingEvent.DataAs(data)
	require.Nil(t, err)

	err = HandleEchoEvent(keptn, *incomingEvent, data, &testLogger)
	require.Nil(t, err)

	assert.Equal(t, "GOT A MESSAGE: hello\n", testLogger.Last())

}

type testWriter struct {
	messages []string
}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	tw.messages = append(tw.messages, msg)
	return len(p), nil
}

func (tw *testWriter) Last() string {
	return tw.messages[len(tw.messages)-1]
}
