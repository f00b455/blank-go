FROM gcr.io/distroless/static-debian12:nonroot

COPY bin/api /app/api

USER nonroot:nonroot
EXPOSE 3002

ENTRYPOINT ["/app/api"]
