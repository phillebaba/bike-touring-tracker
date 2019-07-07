FROM golang:alpine AS build

# Fetch dependencies
RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep
COPY Gopkg.lock Gopkg.toml /go/src/github.com/phillebaba/bike-tracker/
WORKDIR /go/src/github.com/phillebaba/bike-tracker/
RUN dep ensure -vendor-only

# Build go project
COPY . /go/src/github.com/phillebaba/bike-tracker/
RUN go build -o /bin/server cmd/server/main.go

# Only keep build artifact
FROM scratch
COPY --from=build /bin/server /bin/server
ENTRYPOINT ["/bin/server"]
