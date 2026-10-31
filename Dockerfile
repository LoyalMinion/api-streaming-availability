FROM golang:1.27.3 AS build

WORKDIR /src
COPY . .
ARG REVISION=dev
RUN CGO_ENABLED=0 go build -trimpath -buildvcs=false \
    -ldflags "-X main.revision=${REVISION}" \
    -o /api-streaming-availability ./cmd/api-streaming-availability

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /api-streaming-availability /api-streaming-availability
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/api-streaming-availability"]
