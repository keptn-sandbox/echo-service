package eventhandling

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConfigurableSleeper(t *testing.T) {
	expectedSleepTime := 5 * time.Second

	spyTime := &SpyTime{}
	sleeper := ConfigurableSleeper{expectedSleepTime, spyTime.Sleep}
	sleeper.Sleep()

	assert.Equal(t, expectedSleepTime, spyTime.durationSlept)
}

type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
	s.durationSlept = duration
}
