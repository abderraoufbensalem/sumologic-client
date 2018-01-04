FROM alpine
RUN apk add --no-cache ca-certificates
COPY main /
ENTRYPOINT ["./main"]
EXPOSE 8001