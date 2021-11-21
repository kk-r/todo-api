FROM golang:1.17-alpine AS build_base

ARG PACT_BROKER_URL
ENV BROKER_URL $PACT_BROKER_URL
ARG PACT_BROKER_TOKEN
ENV PACT_BROKER_TOKEN $PACT_BROKER_TOKEN


RUN apk add --no-cache git curl ca-certificates wget bash

# Set the Current Working Directory inside the container
WORKDIR /tmp/go-todo-app


#setup binary for Pact provider test
RUN cd /opt && \ 
    curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash 

ENV PATH "$PATH:/opt/pact/bin"

RUN echo "export PATH=/opt/pact/bin:${PATH}" >> /root/.bashrc

RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
&& wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.29-r0/glibc-2.29-r0.apk \
&& apk add glibc-2.29-r0.apk

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/go-todo-app .

# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build_base /tmp/go-todo-app/out/go-todo-app /app/go-todo-app

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/go-todo-app"]