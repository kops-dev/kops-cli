package deploy

const (
	Golang = `FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates
COPY main ./main
RUN chmod +x /main
EXPOSE {{ . }}
CMD ["/main"]`

	Js = `FROM node:10-alpine
RUN mkdir -p /home/node/app/node_modules && chown -R node:node /home/node/app
WORKDIR /home/node/app

COPY package*.json ./
USER node
RUN npm install
COPY --chown=node:node . .

EXPOSE {{ . }}
CMD [ "node", "app.js" ]`

	Java = `FROM openjdk

COPY app.jar app.jar

# copy project dependencies
ADD lib/ lib/

ENTRYPOINT ["java","-jar","app.jar"]`
)
