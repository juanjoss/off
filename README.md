# off

This is toy project which uses the Open Food Facts API to simulate a system in where users (actually, a generator) have monitoring devices for their products (or SSDs), and the SSDs (again, the generator) communicate with a REST backend service by sending a set of randomly generated and scheduled events (by the generator) in order to mantain the stock of products beign monitor by a SSD.

_Why a generator and not an UI?... Because I don't like frontend stuff and I just wanted to learn about microservices and how to schedule jobs._

_Why create a weird word for an hipothetical IoT device which monitors the products we ask it to?... Because if someone implements this system in real life I don't need to go ever again to the supermarket._

_Doesn't something like this already exist?... No. I'm talking about a system which should handle supermarket shopping in the same time or less it would took you. Probably logistically imposible today, but who knows..._

## What's inside?

The components of the system are:

### off-etl

This service extracts data from the OFF API and loads it into the database.

### off-generator

This service schedules and sends events to the order and user services. It acts as the "real users" and the, not at all weird, "SSDs" of the system.

The possible events are `user-registration` and `product-order`.

### off-users

This service handles user requests and communicates with other services over NATS.

HTTP endpoints:
- `POST /api/register`

NATS communication:
- `Sub orders.random` (receive request from generator's product-order event)
- `Sub orders.completed` (receive order and add the product to the user's ssd)
- `Pub orders.new` (reply to generator's product-order request)

### off-orders

This service handles order requests and communicates with other services over NATS.

HTTP endpoints:
- `POST /api/orders`

NATS communication:
- `Sub orders.random` (receive request from generator's product-order event)
- `Pub orders.new` (reply to generator's product-order request)
- `Pub orders.pending` (publish a new order)
- `Sub orders.shipped` (receive and update order status)
- `Sub orders.completed` (receive and update order status)

### off-notifications

This service is used to mock the behind the scenes product-order workflow. It acts as a communication medium with third party external services (supplier and delivery), and relies on `NATS` to read and update an order's state as follows:

1. Subscribes to `orders.pending` and `orders.shipped`.

2. When a new order is created, it will read the order from `orders.pending`, sleep for _t_ minutes (supplier processing time), update its status and publish the order to `orders.shipped`.

2. When the order is read from `orders.shipped`, it will sleep for _t_ minutes (delivery processing time), updates its status and publish it to `orders.completed`.

The sleeping time _t_ (in minutes) will be randomly generated when an order arrives via `orders.pending` and `orders.shipped`.

## Database Schema (ERD)

![db_schema](https://drive.google.com/uc?export=view&id=1ajSHXDxV_ZJ_CnMePGv4mCcYYwZMKAHg)

## Architectural Diagrams

_green = done, red = !green_.

### High level diagram of the project components

![arch_basic](https://drive.google.com/uc?export=view&id=1kRnklQk-EVtD-bonvvCYNwBA7MnEfZW6)

### Architecture of the service layer
![arch_service](https://drive.google.com/uc?export=view&id=1hOVoc3tpvdBLaKKwlg45KcMx2MgVtkyO)

### Workflow for the product-order event 
![orders_workflow](https://drive.google.com/uc?export=view&id=1c6LWkgnjMpJpM263uh90kqXWavSWRR7N)

## Generate swagger docs

The user and order services both have Swagger documentation available. You need to generate and serve it:

```bash
swagger generate spec -o docs.json
swagger serve docs.json
```

## Run it with Docker Compose

The `.env` file is used to for configuration.

```bash
docker compose up -d
```