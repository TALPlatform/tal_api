# Build stage
FROM  golang:alpine3.20 AS builder
ENV GOTOOLCHAIN=auto
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY start.sh .
RUN chmod +x /app/start.sh
EXPOSE 9091
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
