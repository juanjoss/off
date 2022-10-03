# off-project

This is toy project which uses the Open Food Facts API to simulate a system in where users (actually, a generator) have monitoring devices for their products (or SSDs), and the SSDs (again, the generator) communicate with a REST backend service by sending a set of randomly generated and scheduled events (by the generator) in order to mantain the stock of products beign monitor by a SSD.

_Why a generator and not an UI?... Because I don't like frontend stuff and I just wanted to learn about microservices and how to schedule jobs._

_Why create a weird word for an hipothetical IoT device which monitors the products we ask it to?... Because if someone implements this system in real life I don't need to go ever again to the supermarket._

## What's inside?

The files in the repo just provide the configuration to develop and run in docker all the components. You need to clone the repositories of all the services:

- https://github.com/juanjoss/off-etl
- https://github.com/juanjoss/off-generator
- https://github.com/juanjoss/off-orders-service
- https://github.com/juanjoss/off-users-service
- https://github.com/juanjoss/off-notifications-service

The `go.work` file is used to work with multiple go modules.

The components of the system are:

### off-etl

This service extracts data from the OFF API and loads it into the database.

### off-generator

This service schedules and sends events to the order and user services. It acts as the "real users" and the, not at all weird, "SSDs" of the system.

The possible events are `user-registration` and `product-order`.

### off-users-service

This service handles user requests and communicates with other services over NATS.

HTTP endpoints:
- `POST /api/register`

NATS communication:
- `Sub orders.random` (receive request from generator's product-order event)
- `Pub orders.new` (reply to generator's product-order request)

### off-orders-service

This service handles order requests and communicates with other services over NATS.

HTTP endpoints:
- `POST /api/orders`

NATS communication:
- `Sub orders.random` (receive request from generator's product-order event)
- `Pub orders.new` (reply to generator's product-order request)
- `Pub orders.pending` (publish a new order)
- `Sub orders.shipped` (receive and update order status)
- `Sub orders.completed` (receive and update order status)

### off-notifications-service

This service is used to mock the behind the scenes product-order workflow. It acts as a communication medium with third party external services (supplier and delivery), and relies on `NATS` to read and update an order's state as follows:

1. Subscribes to `orders.pending` and `orders.shipped`.

2. When a new order is created, it will read the order from `orders.pending`, sleep for _t_ minutes (supplier processing time), update its status and publish the order to `orders.shipped`.

2. When the order is read from `orders.shipped`, it will sleep for _t_ minutes (delivery processing time), updates its status and publish it to `orders.completed`.

The sleeping time _t_ (in minutes) will be randomly generated when an order arrives via `orders.pending` and `orders.shipped`.

## Architectural Diagrams

_green = done, red = !green_.

### High level diagram of the project components.
![arch_basic](https://drive.google.com/uc?export=view&id=1kRnklQk-EVtD-bonvvCYNwBA7MnEfZW6)

### Architecture of the service layer.
![arch_service](https://drive.google.com/uc?export=view&id=1xB-YAc2PKwYC5Pruw6V7xT4k-iTSLgfn)

### Workflow for the product-order event. 
![orders_workflow](https://drive.google.com/uc?export=view&id=14DvmCakoJZLWIhCLawNHIGPc4I7sKI01)

## Run it with docker compose:

Both environments use the `.env` file for configuration.

- Development environment:

```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

- Production environment:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

## Generate Swagger Docs

The user and order services both have Swagger documentation available. You need to generate and serve it:

```bash
swagger generate spec -o docs.json
swagger serve docs.json
```

## Tekton

### Install Tekton

k apply -f https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml

### Install `tkn` (Tekton CLI)

```bash
curl -LO https://github.com/tektoncd/cli/releases/download/v0.26.0/tektoncd-cli-0.26.0_Linux-64bit.deb
sudo dpkg -i <package-name>
```

### Run the Pipeline

```bash
k apply -f task.yaml
k apply -f pipeline.yaml
k create -f pipelinerun.yaml
tkn pipelinerun logs <name> -f
```

## Work in Progress

- Add tests to basicaly... Everything.
- Add authentication for the HTTP endpoints.
- Change most of the communication between services to [NATS](https://nats.io/).
- Add [qrgen](https://github.com/juanjoss/qrgen) and [shorturl](https://github.com/juanjoss/shorturl) services.
- Create a CI/CD pipeline with [Tekton](https://tekton.dev/) and [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) inside Kubernetes. Tekton will use [kaniko](https://github.com/GoogleContainerTools/kaniko) for building Docker images.
- Create manifest files to deploy into Kubernetes, maybe [kompose?](https://kompose.io/) since docker compose files already exist.
