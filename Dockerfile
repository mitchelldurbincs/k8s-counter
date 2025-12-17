# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /src

# Cache deps first
COPY go.mod ./
RUN go mod download

# Copy source + build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/app .

# Runtime stage
FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=build /out/app /app
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app"]
