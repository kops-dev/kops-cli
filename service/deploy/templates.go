package deploy

const (
	Golang = `FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 go build -o main

FROM alpine:3.15

COPY --from=builder /app /

RUN chmod 777 /main

EXPOSE 9000

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
