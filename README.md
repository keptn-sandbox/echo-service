# Keptn Echo Service

Simple Keptn Service that does nothing but echoing a string whenever it receives a ```sh.keptn.event.echo.triggered```
event.

## Concept of Execution

1. ```sh.keptn.event.echo.triggered``` event gets received
2. ```sh.keptn.event.echo.started``` event gets sent
3. The echo service sleeps for a specific amount of time. (default: ```2000ms```)
4. ```sh.keptn.event.echo.finished``` event gets sent

## Adaptation Data
- ```SLEEP_TIME_MS```:  amount of milliseconds to sleep between sending the ```sh.keptn.event.echo.started``` and the
```sh.keptn.event.echo.finished``` event. Format: ```int``` Default: ```1000```
- ```SIMULATE_EVENTS_OUT_OF_ORDER```: Flag indicating whether the echo service shall send ```sh.keptn.event.echo.started```
and the ```sh.keptn.event.echo.finished``` event in correct order, or not. Useful for testing. Format: ```bool``` Default: ```false```
- ```EVENTBROKER```: address of the keptn event broker