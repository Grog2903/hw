FROM golang:1.22 as build

ENV BIN_FILE /opt/calendar/calendar-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ../../build ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar/main.go

FROM alpine:3.9

ENV BIN_FILE /opt/calendar/calendar-app
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/calendar_config.yaml

COPY ../configs/calendar_config.yaml ${CONFIG_FILE}
CMD ${BIN_FILE} -config ${CONFIG_FILE}