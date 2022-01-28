FROM golang:1.17-alpine as build

WORKDIR /go/src/todoNote

COPY go.mod ./
COPY go.sum ./
COPY vendor ./
COPY ../.. .


#RUN go get -d -v ./...
#RUN go install -v ./...

RUN go build -o /todo-note .

FROM alpine:latest
WORKDIR /
COPY --from=build /todo-note /todo-note
CMD ["/todo-note"]

#ADR architecture deseision records
# vendoring go mod vendor
