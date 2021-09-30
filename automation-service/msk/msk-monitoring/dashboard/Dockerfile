FROM node:12.11.1-alpine as builder

RUN apk add --no-cache git curl gcc libc-dev

RUN git clone https://github.com/GeneralMills/BurrowUI.git

WORKDIR /BurrowUI

RUN npm install
RUN npm install -g @angular/cli@6.1.1

RUN ng build --prod

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add --update nodejs npm

RUN npm config set unsafe-perm true
RUN npm update -g npm

RUN mkdir -p /BurrowUI/server /BurrowUI/dist
WORKDIR /BurrowUI

COPY ./package.json /BurrowUI/
COPY --from=builder /BurrowUI/server.js /BurrowUI/
COPY --from=builder /BurrowUI/server /BurrowUI/server
COPY --from=builder /BurrowUI/dist /BurrowUI/dist

RUN npm install --production
RUN rm -rf /usr/lib/node_modules/npm/

# Create a group and user
RUN addgroup -S wwgroup && adduser -S wwuser -G wwgroup
RUN chown -R wwuser /BurrowUI
# Tell docker that all future commands should run as the appuser wwuser
USER wwuser

EXPOSE 3000

CMD [ "node", "server" ]