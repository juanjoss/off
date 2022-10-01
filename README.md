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

The `go.work` file is used to work with multiple go modules.

The 4 components of the system are:

### off-etl

This service extracts data from the OFF API and loads it into the database.

### off-users-service

This service handles user and SSD requests. The endpoints are:

- `POST /api/users/register` (to handle an user-registration event)
- `POST /api/users/ssds/products` (to handle an add-product-to-ssd event)
- `GET /api/users/ssds/random` (to get a random SSD)

### off-orders-service

This service handles product and order requests. The endpoints are:

- `GET /api/products` (to get all products)
- `GET /api/products/randomProductFromUserSSD` (to get a random product from a user's SSD)
- `GET /api/products/random` (to get a random product, not necesarilly in a SSD)
- `POST /api/products/orders` (to handle a product-order event)

### off-generator

This service schedules and sends events to the order and user services. It acts as the "real users" and the not at all weird "SSDs" of the system.

The possible events are `user-registration`, `add-product-to-ssd` and `product-order`.

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

## Work in Progress

- Change most of the communication between services to `NATS`.
- Add [qrgen](https://github.com/juanjoss/qrgen) and [shorturl](https://github.com/juanjoss/shorturl) services.
- Create a CI/CD pipeline with `Tekton` and `ArgoCD` inside Kubernetes.
- Find a way to create Docker images inside Kubernetes (maybe [kaniko?](https://github.com/GoogleContainerTools/kaniko)).
- Create manifest files to deploy into Kubernetes, maybe [kompose?](https://kompose.io/) since docker compose files already exist.