FROM golang:1.24 AS build-stage

WORKDIR /app

COPY . ./

RUN make


FROM debian:11-slim AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/build/front-service /front-service

EXPOSE 7070

# USER nonroot:nonroot

ENTRYPOINT ["/front-service"]
