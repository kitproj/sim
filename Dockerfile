FROM golang:1.20-alpine as build

WORKDIR /go/src/github.com/kitproj/sim

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build .

FROM scratch as sim

COPY --from=build /go/src/github.com/kitproj/sim/sim .

VOLUME [ "/apis" ]

CMD [ "sim" ]