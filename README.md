# Keptn Echo Service

Simple Keptn Service that does nothing but echoing a string whenever it receives a ```sh.keptn.event.echo.triggered```
event.

## Concept of Execution

1. ```sh.keptn.event.echo.triggered``` event gets received
2. ```sh.keptn.event.echo.started``` event gets sent
3. The echo service sleeps for a specific amount of time. (default: ```1000ms```)
4. ```sh.keptn.event.echo.finished``` event gets sent

## Configuration
- ```SLEEP_TIME_MS```:  amount of milliseconds to sleep between sending the ```sh.keptn.event.echo.started``` and the
```sh.keptn.event.echo.finished``` event. Format: ```int``` Default: ```1000```
- ```SIMULATE_EVENTS_OUT_OF_ORDER```: Flag indicating whether the echo service shall send ```sh.keptn.event.echo.started```
and the ```sh.keptn.event.echo.finished``` event in correct order, or not. Useful for testing. Format: ```bool``` Default: ```false```
- ```EVENTBROKER```: address of the keptn event broker

## Compatibility Matrix
| Keptn Version    | [Echo-Service Service Image](https://hub.docker.com/r/keptnsandbox/echo-service/tags) |
|:----------------:|:----------------------------------------:|
|   0.8.0 (master) | keptnsandbox/echo-service:0.1.0          |

## Installation
As any Keptn Service the *echo-service* needs to be installed on the k8s cluster where you have installed Keptn!

### Deploy in your Kubernetes cluster

Please use the same namespace for the *echo-service* as you are using for Keptn, e.g: keptn.
    ```console
    kubectl apply -f deploy/service.yaml -n keptn
    ```
* This installs the *echo-service* into the `keptn` namespace, which you can verify using:
    ```console
    kubectl -n keptn get deployment echo-service -o wide
    kubectl -n keptn get pods -l run=echo-service
    ```
### Up- or Downgrading
Adapt and use the following command in case you want to up- or downgrade your installed version (specified by the `$VERSION` placeholder):
```console
kubectl -n keptn set image deployment/echo-service echo-service=keptncontrib/echo-service:$VERSION --record
```
### Uninstall
To delete a deployed *echo-service*, use the file `deploy/*.yaml` files from this repository and delete the Kubernetes resources:
```console
kubectl delete -f deploy/service.yaml -n keptn
```