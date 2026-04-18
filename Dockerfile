FROM alpine:3.14

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root

# Copy binary đã build sẵn và assets
COPY api .
COPY assets ./assets

EXPOSE 8080

CMD ["./api", "start"]