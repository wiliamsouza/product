# build stage
FROM golang:alpine AS build
RUN apk update \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates
RUN apk --no-cache add build-base git
ADD . /src
RUN cd /src && make deps && make build

# final stage
FROM alpine
WORKDIR /usr/sbin
COPY --from=build /src/dist/productd_linux_amd64/productd /usr/sbin/
COPY --from=build /src/dist/productctl/productctl /usr/sbin/
ENTRYPOINT ["/usr/sbin/productd"]
CMD ["help"]
