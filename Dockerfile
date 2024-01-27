# FROM golang:1.21.6-alpine

# # Working directory
# WORKDIR /app

# # Copy the files in the work dir
# COPY go.mod go.sum ./

# # Download dependecies
# RUN go mod download

# # Copy source soce in the work dir
# COPY . .

# # Build the Go app
# RUN go build -o main

# # Expose port
# EXPOSE 5000

# # Set the entry point of the container to the executable
# CMD ["./main"]


# the above docker file is about 600 MB

# the below is 27 MB

FROM golang:alpine AS builder

# Working directory
WORKDIR /app

# Copy the files in the work dir
COPY go.mod go.sum ./

# Download dependecies
RUN go mod download

# Copy source soce in the work dir
COPY . .

# Build the Go app
RUN go build -o main .

FROM alpine

WORKDIR /app

# Copies the compiled executable (main) from the builder stage to the /app directory in the current stage (the final image).
# --from=builder specifies that the source of the copy is the previous build stage named builder.
COPY --from=builder /app/main /app/main

# Set environment variables
ENV HTTP_LISTEN_ADDRESS=:5000
ENV JWT_SECRET=aSh2QuIevrHrbegqjgZx5OFxpgebqiSaTgoIXgWccWc2
ENV MONGO_DB_NAME=go-hotel
ENV MONGO_DB_URL=mongodb://localhost:27017/

# Expose port
EXPOSE 5000

# Set the entry point of the container to the executable
CMD ["./main"]