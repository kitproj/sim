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

LABEL org.opencontainers.image.title="Sim" \
      org.opencontainers.image.description="Sim is straight-forward API simulation tool that's tiny, fast, secure and scalable." \
      org.opencontainers.image.url="https://github.com/kitproj/sim" \
      org.opencontainers.image.licenses="MIT"

CMD [ "/sim", "/apis" ]