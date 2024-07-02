# syntax=docker/dockerfile:1

FROM golang:1.22.4

WORKDIR /app


COPY ./src .
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

RUN go mod download

RUN go build -o /weather-forecast-api

EXPOSE 8080

CMD [ "/weather-forecast-api" ]