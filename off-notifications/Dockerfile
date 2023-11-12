ARG GO_VERSION=1.19.0

FROM golang:${GO_VERSION}-alpine AS dev
RUN go env -w GOPROXY=direct
RUN apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR /app
# install reflex for hot reloading
RUN go install github.com/cespare/reflex@latest
# copy files to container
COPY . .
# download dependencies
RUN go mod tidy

FROM dev as build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \ 
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo -o notifications .

FROM gcr.io/distroless/static AS prod
USER nonroot:nonroot
COPY --from=build --chown=nonroot:nonroot /app/notifications .
ENTRYPOINT ["./notifications"]