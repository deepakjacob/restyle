# https://dev.to/plutov/docker-and-go-modules-3kkn
# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
RUN echo $(ls -1 /app)

# final stage
FROM scratch
COPY --from=builder /app/restyle /app/


EXPOSE 8000
ENTRYPOINT ["/app/restyle"]
