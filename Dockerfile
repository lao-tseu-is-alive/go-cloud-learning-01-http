FROM golang:1-alpine3.15 as builder

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

######## Start a new stage  #######
FROM alpine:3.15
#
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10111 -S gouser && adduser -S -G gouser -H -u 10111 gouser
USER gouser

WORKDIR /goapp
#
## Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/todosServer .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./todosServer"]

