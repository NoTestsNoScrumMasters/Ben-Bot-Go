FROM golang:1.19.1 as builder

# environment variables
ENV APP_NAME Ben-Bot-Go
ENV CMD_PATH main.go
RUN mkdir -p /go/src/github.com/NoTestsNoScrumMasters/Ben-Bot-Go
COPY . /go/src/github.com/NoTestsNoScrumMasters/Ben-Bot-Go/
WORKDIR /go/src/github.com/NoTestsNoScrumMasters/Ben-Bot-Go/
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/Ben-Bot-Go
RUN cp /go/src/github.com/NoTestsNoScrumMasters/Ben-Bot-Go/.env /app/

# Run Stage
FROM alpine:3.14
 
# Set environment variable
ENV APP_NAME Ben-Bot-Go

RUN apk add curl
RUN mkdir -p /app
# copy data into image
COPY --from=builder /app/Ben-Bot-Go .
COPY --from=builder /app/.env .
 
# listening on port  8080
EXPOSE 8080

# set argument vars in docker-run command
ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY
ARG AWS_DEFAULT_REGION

# ENV AWS agrument
ENV AWS_ACCESS_KEY_ID $AWS_ACCESS_KEY_ID
ENV AWS_SECRET_ACCESS_KEY $AWS_SECRET_ACCESS_KEY
ENV AWS_DEFAULT_REGION $AWS_DEFAULT_REGION

# set Discord APP token argument
ARG DISCORD_BOT_TOKEN

# ENV discord argument
ENV DISCORD_BOT_TOKEN $DISCORD_BOT_TOKEN

EXPOSE 3000

ENTRYPOINT [ "/app/Ben-Bot-Go", "-env", "docker" ]

# Start app
CMD ./"main.go"