# Kubedock

Kubedock is an minimal implementation of the docker api that will orchestrate containers on a kubernetes cluster, rather than running containers locally. The main driver for this project is to run tests that require docker-containers inside a container, without the requirement of running docker-in-docker within resource heavy containers. Containers that are orchestrated by kubedock are considered short-lived and emphemeral and not intended to run production services. An example use case is running [testcontainers-java](https://www.testcontainers.org) enabled unit-tests in a tekton pipeline. In this use case, running kubedock in a sidecar can help orchestrating containers inside the kubernetes cluster instead of within the task container itself.

## Quick start

Running this locally with a testcontainers enabled unit-test requires to run kubedock (`kubedock server`). After that start the unit tests in another terminal with the below environment variables set, for example:

```bash
export TESTCONTAINERS_RYUK_DISABLED=true
export DOCKER_HOST=tcp://127.0.0.1:8999
mvn test
```

The default configuration for kubedock is to orchestrate in the namespace that has been set in the current context. This can be overruled with -n argument (or via the `NAMESPACE` environment variable). The service requires permissions to create Deployments, Services and Configmaps in the namespace.

To see a complete list of available options and additional examples: `kubedock --help`.

# Implementation

When kubedock is started with `kubedock server` it will start an API server on port :8080, which can be used as a drop-in replacement for the default docker api server. Additionally, kubedock can also start listening to an unix-socket (`docker.sock`).

## Containers

Container API calls are translated towards kubernetes Deployment resources. When a container is started, it will create port-forwards for the ports that should be exposed (only tcp is supported). Starting a container is a blocking call that will wait until the Deployment results in a running Pod. By default it will wait for maximum 1 minute, but this is configurable with the `--timeout` argument. The logs API calls will always return the complete history of logs, and doesn't differentiate between stdout/stderr. All log output is send as stdout. Executions in the containers are supported.

## Volumes

Volumes are implemented by copying over the source content towards the container by means of an init-container that is started before the actual container is started. By default the kubedock image with the same version as the running kubedock is used as the init container. However, this can be any image that has tar available and can be configured with the `--initimage` argument.

Volumes are one-way copies and emphemeral. This typically means, any data that is written into the volume is not available locally. This also means that mounts to devices, or sockets are not supported (e.g. mounting a docker-socket).

Copying data from a running container back towards the client is not supported either. Also be aware that copying data towards a container will implicitly start the container. This is different compared to a real docker api, where a container can be in an unstarted state. To 'workaround' this, use a volume instead.

## Networking

Kubedock flattens all networking, which basicly means that everything will run in the same namespace. This should be sufficient for most use-cases. Network aliases are supported. When a network alias is present, it will create a service exposing all ports that have been exposed by the container. If no ports are configured, kubedock is able to fetch ports that are exposed in the container image. To do this, kubedock should be started with the `--inspector` argument.

## Images

Kubedock implements the images API by tracking which images are requested. It is not able to actually build images. If kubedock is started with `--inspector`, kubedock will fetch configuration information about the image by calling external container registries. This configuration includes ports that are exposed by the container image itself, and increases network aliases support. The registries should be configured by the client (for example by doing a `skopeo login`).

## Namespace locking

If multiple kubedocks are using the namespace, it might be possible there will be collisions in network aliases. Since networks are flattend (see Networking), all network aliases will result in a Service with the name of the given network alias. To ensure tests don't fail because of these name collisions, kubedock can lock the namespace while it's running. When enabling this with the `--lock` argument, kubedock will create a Configmap called `kubedock-lock` in the namespace in which it tracks the current ownership.

## Resources cleanup

Kubedock will dynamically create deployments and services in the configured namespace. If kubedock is requested to delete a container, it will remove the deployment and related services. Kubedock will also delete all the resources (Services and Deployments) it created in the running instance before exiting (identified with the `kubedock.id` label).

### Automatic reaping

If a test fails and didn't clean up its started containers, these resources will remain in the namespace. To prevent unused deployments and services lingering around, kubedock will automatically delete deployments and services that are older than 15 minutes (default) if it's owned by the current process. If the deployment is not owned by the running process, it will delete it after 30 minutes if the deployment or service has the label `kubedock=true`. 

### Forced cleaning

The reaping of resources can also be enforced at startup. When kubedock is started with the `--prune-start` argument, it will delete all resources that have the `kubedock=true` before starting the API server. These resource includes resources created by other instances. 

# See also

* https://hub.docker.com/r/joyrex2001/kubedock
