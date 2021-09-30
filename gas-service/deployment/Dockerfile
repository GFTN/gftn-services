
FROM gftn/node-alpine:latest

WORKDIR /usr/src/app

RUN pwd


COPY package*.json ./
COPY . .

RUN mkdir -p /opt/logs/
RUN touch /opt/logs/logs.txt

RUN npm install
RUN find . -name '*.key' -delete

# Create a group and user
RUN addgroup -S wwgroup && adduser -S wwuser -G wwgroup
# Change owner for log file
RUN chown wwuser /opt/logs/logs.txt
# Tell docker that all future commands should run as the appuser wwuser
USER wwuser

EXPOSE 8080
CMD [ "npm", "start" ]
