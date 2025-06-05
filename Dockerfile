FROM golang AS build
WORKDIR /src
COPY . .
RUN go build -o /out/container-log-sanitizer main.go

FROM alpine
COPY --from=build /out/container-log-sanitizer /bin/container-log-sanitizer
ENTRYPOINT ["/bin/container-log-sanitizer"]
