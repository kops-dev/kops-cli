package templates

const (
	golang = `FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates
COPY main ./main
RUN chmod +x /main
EXPOSE {{ port }}
CMD ["/main"]`

	js = `FROM node:10-alpine
RUN mkdir -p /home/node/app/node_modules && chown -R node:node /home/node/app
WORKDIR /home/node/app

COPY package*.json ./
USER node
RUN npm install
COPY --chown=node:node . .

EXPOSE {{ port }}
CMD [ "node", "app.js" ]`

	java = `FROM openjdk

COPY app.jar app.jar

# copy project dependencies
ADD lib/ lib/

ENTRYPOINT ["java","-jar","app.jar"]`
)

var TmplMap = map[string]string{
	"go":   golang,
	"java": java,
	"js":   js,
}
