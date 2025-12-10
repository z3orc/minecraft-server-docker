FROM golang:1.25-alpine AS build
WORKDIR /app

# Download Go modules
COPY go.mod ./
RUN go mod download

# Copy source code and build
COPY ./ ./
RUN apk add --no-cache make
RUN make build

# Create container
FROM eclipse-temurin:21-jre

WORKDIR /app
COPY --from=build /app/out/runner /app
COPY --from=build /app/entrypoint.sh /app

# RUN chmod +x runner
RUN chmod +x entrypoint.sh

ARG UNAME=minecraft
RUN groupadd -g 1000 -o $UNAME
RUN useradd -m -u 1000 -g 1000 -o -s /bin/bash $UNAME
USER $UNAME


ENV SERVER_JAR=server.jar
ENV TIMEOUT=60
ENV USE_SIGKILL=false

WORKDIR /data

ENTRYPOINT [ "/app/entrypoint.sh" ]

