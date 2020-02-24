# Build rgeoSrv in build-env
FROM golang:1.13.8-alpine3.11 AS build-env

# Add user to prevent running as root
RUN adduser -D -u 10000 rgeoSrv
RUN mkdir /rgeoSrv/ && chown rgeoSrv /rgeoSrv/
USER rgeoSrv

# Add files to container
WORKDIR /rgeoSrv/
ADD . /rgeoSrv/

# Get Go dependencies
RUN go get

# Build rgeoSrv
RUN go build -o rgeoSrv

# New container to run in
FROM alpine:3.11


# Add user to prevent running as root
RUN adduser -D -u 10000 rgeoSrv
USER rgeoSrv

# Copy files from build-env
WORKDIR /
COPY --from=build-env /rgeoSrv/rgeoSrv /

EXPOSE 8080

CMD ["/rgeoSrv", "-addr", ":8080"]
