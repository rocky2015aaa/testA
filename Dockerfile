FROM golang:1.22 AS build

RUN mkdir /app
ADD . /app
WORKDIR /app

# Define build-time arguments
ARG VERSION=1.0.0
ARG BUILD=production
ARG DATE=1970-01-01_00:00:00

# Build the Go application with build-time arguments
RUN CGO_ENABLED=0 go build -ldflags "-X github.com/shmoulana/hashocr/internal/config.Version=${VERSION} -X github.com/shmoulana/hashocr/internal/config.Build=${BUILD} -X github.com/shmoulana/hashocr/internal/config.Date=${DATE}" -o hashocrsvr cmd/hashocrsvr/main.go

FROM alpine:latest

# Add a community repository if available
RUN apk add --no-cache \
    poppler-utils \
    && apk add --no-cache --virtual .build-deps gcc musl-dev \
    && apk del .build-deps

RUN mkdir /application
WORKDIR /application

COPY --from=build /app/hashocrsvr .
COPY --from=build /app/.env .

EXPOSE 9081

CMD ["./hashocrsvr"]
