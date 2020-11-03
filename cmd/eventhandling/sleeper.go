package eventhandling

import "time"

//Sleeper is responsible to pause execution
type Sleeper interface {
	//Sleep pauses execution
	Sleep()
}

//ConfigurableSleeper sleeps a configured amount of time
type ConfigurableSleeper struct {
	duration time.Duration
	sleep    func(time.Duration)
}

func (c ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
}

//NewConfigurableSleeper returns a new sleeper that will sleep for a specified duration
func NewConfigurableSleeper(duration time.Duration, sleepFunc func(time.Duration)) Sleeper {
	return &ConfigurableSleeper{
		duration: duration,
		sleep:    sleepFunc,
	}
}

//TestSleeper is an implementation of a Sleeper which pretends to sleep
type TestSleeper struct {
}

func (s *TestSleeper) Sleep() {
	//no-op
}
