# Nixiesearch Lambda Container Image
# Embedded search engine with S3 index storage

FROM public.ecr.aws/lambda/provided:al2023 AS builder

# Install build dependencies
RUN dnf install -y golang git tar gzip

# Set Go environment
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin

# Build the Lambda handler
WORKDIR /build
COPY backend/cmd/nixiesearch/ ./cmd/nixiesearch/
COPY backend/internal/ ./internal/
COPY backend/go.mod backend/go.sum ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/nixiesearch

# Final image
FROM public.ecr.aws/lambda/provided:al2023

# Copy the Lambda handler to the correct location
COPY --from=builder /build/bootstrap /var/runtime/bootstrap

# Set permissions
RUN chmod 755 /var/runtime/bootstrap

# Entrypoint for provided runtime
ENTRYPOINT ["/var/runtime/bootstrap"]
