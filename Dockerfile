# syntax=docker/dockerfile:1

FROM golang:1.17-alpine as builder

LABEL maintainer="cgil"

RUN apk update && apk add curl \
                          git \
                          bash \
                          make \
                          openssh-client && \
     rm -rf /var/cache/apk/*

WORKDIR /app

RUN git clone https://github.com/lao-tseu-is-alive/go-cloud-learning-01-http.git .

RUN make build

######## Start a new stage from scratch #######
FROM alpine:latest
#
WORKDIR /root/
#
## Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/todosServer .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./todosServer"]

