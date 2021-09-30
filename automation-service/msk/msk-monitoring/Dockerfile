FROM golang:alpine as builder

ENV BURROW_SRC /usr/src/

RUN apk add --no-cache git curl gcc libc-dev

WORKDIR $BURROW_SRC
RUN git clone https://github.com/linkedin/Burrow.git
WORKDIR ./Burrow

RUN go mod tidy && go build -o /tmp/burrow .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /tmp/burrow /app/
ADD ./burrow.toml /etc/burrow/
RUN mkdir -p /var/private/ssl

# Create a group and user
RUN addgroup -S wwgroup && adduser -S wwuser -G wwgroup
RUN chown -R wwuser /app
RUN chown -R wwuser /etc/burrow
RUN chown -R wwuser /var/private/ssl
# Tell docker that all future commands should run as the appuser wwuser
USER wwuser

EXPOSE 8000

CMD ["/app/burrow", "--config-dir", "/etc/burrow"]