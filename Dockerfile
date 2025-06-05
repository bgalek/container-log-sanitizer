FROM golang AS build
WORKDIR /src
COPY . .
RUN go build -o /out/container-log-sanitizer main.go

FROM scratch
COPY --from=build /out/container-log-sanitizer /container-log-sanitizer
ENTRYPOINT ["/container-log-sanitizer"]
