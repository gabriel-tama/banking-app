# Step 1: Build the binary
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container

COPY go.mod go.sum ./

COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Step 2: Use a minimal base image to run the application
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

ARG DB_NAME
ARG DB_PORT
ARG DB_HOST
ARG DB_USERNAME
ARG DB_PASSWORD
ARG DB_PARAMS
ARG BCRYPT_SALT
ARG JWT_SECRET
ARG S3_ID
ARG S3_SECRET_KEY
ARG S3_BUCKET_NAME
ARG S3_REGION
ARG PEM_FILE

ENV DB_NAME=$DB_NAME
ENV DB_PORT=$DB_PORT
ENV DB_HOST=$DB_HOST
ENV DB_USERNAME=$DB_USERNAME
ENV DB_PASSWORD=$DB_PASSWORD
ENV DB_PARAMS=$DB_PARAMS
ENV JWT_SECRET=$JWT_SECRET
ENV BCRYPT_SALT=$BCRYPT_SALT
ENV S3_ID=$S3_ID
ENV S3_SECRET_KEY=$S3_SECRET_KEY
ENV S3_BUCKET_NAME=$S3_BUCKET_NAME
ENV S3_REGION=$S3_REGION

COPY $PEM_FILE ap-southeast-1-bundle.pem
RUN chmod 600 ap-southeast-1-bundle.pem


# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
