
# Start from the latest golang base image
FROM golang:1.19-alpine3.16 as builder

# Add Maintainer Info
LABEL maintainer="Shourya Pratap Singh <shourya@finbox.in>"

# Set the Current Working Directory inside the container
WORKDIR /app

ARG GH_PAT
ARG PROMETHEUS_PORT=9090

ENV GO111MODULE=on
RUN apk update && apk add --no-cache git openssh-client ca-certificates tzdata && update-ca-certificates
RUN apk add build-base
RUN apk add curl
RUN apk add libxml2-dev libxslt-dev xz-dev zlib-dev

# Make ssh dir
# RUN mkdir /root/.ssh/

# Copy over private key, and set permissions
# RUN echo "$SSH_PRIVATE_KEY" > /root/.ssh/id_rsa && \
#     chmod 600 /root/.ssh/id_rsa

# # Create known_hosts and add GitHub to prevent unknown host error
# RUN touch /root/.ssh/known_hosts && \
#     ssh-keyscan github.com >> /root/.ssh/known_hosts


RUN git config --global url."https://${GH_PAT}@github.com/".insteadOf "https://github.com/"

RUN export GOPRIVATE=github.com/finbox-in/*

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -o go-events .
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o go-events .

######## Start a new stage from scratch #######
FROM alpine:3.16

RUN apk --no-cache add ca-certificates

RUN apk add --no-cache tzdata
RUN apk add curl
RUN apk add build-base
RUN apk add libxml2-dev libxslt-dev xz-dev zlib-dev

RUN addgroup -S lendinguser && adduser -S lendinguser -G lendinguser

# creating folder for dynamic config usage (certificates uploads, etc)
RUN mkdir /etc/lendinguser/
RUN mkdir /app/
RUN chown lendinguser /etc/lendinguser
RUN chown lendinguser /app

USER lendinguser
WORKDIR /etc/lendinguser

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/go-events .

# Copy migration scripts
COPY migrations/scripts scripts

# Copy casbin model
COPY rbac/config config
# Copy partner push config
COPY conf/partnerpush.yml config
COPY conf/temporalsubmitconfig.yml config

# Copy keys yaml for ssm
COPY conf/keys.yaml config/keys.yaml

# Expose port 3335 to the outside world	
EXPOSE 3335

EXPOSE ${PROMETHEUS_PORT}

# Build Args
ARG LOG_DIR=/app/logs

# Create Log Directory
RUN mkdir -p ${LOG_DIR}

# Environment Variables

ENV LOG_DIR=${LOG_DIR}/app.log
ENV PROMETHEUS_PORT=${PROMETHEUS_PORT}

# Declare volumes to mount
# VOLUME [${LOG_DIR}]

# check if we want to exit on health check fail
#HEALTHCHECK --start-period=300s --interval=30s --timeout=2s --retries=5 CMD curl -f http://localhost:3335/ || exit 1 

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/go-events .

COPY --from=builder /app/start.sh .
RUN sed -i 's/\r$//' ./start.sh && chmod +x ./start.sh


# Command to run the executable
ENTRYPOINT ["sh", "-c", "./go-events"]
