FROM golang:1.17-alpine as build-env

# environment variables
ENV APP_NAME Ben-Bot-Go
ENV CMD_PATH main.go

# create source code directory
COPY . $GOPATH/src/$APP_NAME
WORKDIR $GOPATH/src/$APP_NAME

# install dependencies

RUN go mod download

# build 
RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

# Run Stage
FROM alpine:3.14
 
# Set environment variable
ENV APP_NAME Ben-Bot-Go
 
# copy data into image
COPY --from=build-env /$APP_NAME .
 
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

# Start app
CMD ./$APP_NAME