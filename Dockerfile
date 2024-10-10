FROM golang:1.17-alpine as build

RUN apk update && apk upgrade
RUN apk add --no-cache wget

WORKDIR /src
COPY . /src

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o /build/invitesvc ./cmd/invitesvc

FROM scratch

COPY --from=build /build/invitesvc /bin/

ENTRYPOINT ["/bin/invitesvc"]