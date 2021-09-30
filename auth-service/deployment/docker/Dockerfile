# USAGE: 
#   Build image - From /auth-serivce $ make docker
#   Run image - $ docker run  --env env=dev --env pass=somepassphrase --env cred=.  -it -p 4000:80 gftn/auth-service:latest

# used base node-alpine to get packages (image includes npm)
FROM node:10-alpine as builder
# FROM node as builder

RUN node -v
RUN npm -v

WORKDIR /app

# COPY package.json package-lock.json app.js ./
COPY authentication/build .

RUN pwd ; ls -a

# install npm for production
RUN npm i --prod

# delete out test kesy for twistlock scan
# see https://github.com/GFTN/Security-Gaps-and-Issues/issues/130
# RUN rm -Rf ./node_modules/agent-base/test
RUN find . -name '*.key' -delete

# use ww alpine images 
#   Note1: impage excludes npm 
#   Note2: last FROM stmt in docker file defines the final image to use for build
# FROM gftn/node:alpine-3.10
FROM ip-team-worldwire-docker-local.artifactory.swg-devops.com/gftn/node-alpine
# FROM node
# FROM node:10-alpine

WORKDIR /app

# copy over everything from build
COPY --from=builder /app .

# log output from running app
RUN mkdir -p /opt/logs/
RUN touch /opt/logs/logs.txt

# Create a group and user
RUN addgroup -S wwgroup && adduser -S wwuser -G wwgroup

# Change owner for log file
RUN chown wwuser /opt/logs/logs.txt
RUN chown wwuser /app

# Tell docker that all future commands should run as the appuser wwuser
USER wwuser

EXPOSE 8080

# start npm app
ENTRYPOINT [ "sh", "docker-entrypoint.sh" ]

# # To run image with access to cmd line:
# # $ sudo docker run -it -p 4000:80 [IMAGE_ID]
# ENTRYPOINT [ "sh" ]