FROM golang:latest

RUN apt-get install git

RUN mkdir -p /src/rebar
WORKDIR /src/rebar

# Install coveralls
RUN go get -u github.com/mattn/goveralls

COPY . .
