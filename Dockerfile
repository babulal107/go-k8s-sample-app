# Build stage: Build the Go binary
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum first to download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# OR can copy on required app files
#COPY cmd /app
#COPY internal /app

# Build the Go application
RUN CGO_ENABLED=0 go build -v -o ./main ./cmd

# Run ls and keep the container open for debugging
RUN ls -la /app && sleep 5

# --------- Stage 2: Create Minimal Runtime Image using Distroless image as scratch ---------
# Fincal stage: size will be reduce around 4mb only instant of 890mb at above stage
FROM alpine:latest

# Used scrathc as distroless image for Golang, If it's Distroless, it may not support interactive shells at all.
# FROM scratch
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# FROM gcr.io/distroless/static:nonroot

# Copy the statically built binary and configs
WORKDIR /app
COPY --from=builder /app/main .

# Run the application as a non-root user (security best practice)
USER 1000:1000

# Expose port 8080
EXPOSE 8080

# Command to run the application
# CMD ["./main"]
ENTRYPOINT ["./main"]
