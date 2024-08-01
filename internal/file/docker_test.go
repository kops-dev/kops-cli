package file

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/logging"
)

func Test_CreateDockerFile(t *testing.T) {
	var (
		goContent = `FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates
COPY main ./main
RUN chmod +x /main
EXPOSE 8080
CMD ["/main"]`
		jsContent = `FROM node:10-alpine
RUN mkdir -p /home/node/app/node_modules && chown -R node:node /home/node/app
WORKDIR /home/node/app

COPY package*.json ./
USER node
RUN npm install
COPY --chown=node:node . .

EXPOSE 8080
CMD [ "node", "app.js" ]`
		javaContent = `FROM openjdk

COPY app.jar app.jar

# copy project dependencies
ADD lib/ lib/

ENTRYPOINT ["java","-jar","app.jar"]`
	)

	defer os.Remove("Dockerfile")

	testCases := []struct {
		desc    string
		lang    string
		expErr  error
		content string
	}{
		{
			desc:   "language not supported yet",
			lang:   "pascal",
			expErr: errLanguageNotSupported,
		},
		{
			desc:    "golang DockerFile",
			lang:    golang,
			expErr:  nil,
			content: goContent,
		},
		{
			desc:    "javascript DockerFile",
			lang:    js,
			expErr:  nil,
			content: jsContent,
		},
		{
			desc:    "java DockerFile",
			lang:    java,
			expErr:  nil,
			content: javaContent,
		},
	}

	ctx := &gofr.Context{
		Container: &container.Container{Logger: logging.NewMockLogger(logging.INFO)},
	}

	for _, tc := range testCases {
		err := CreateDockerFile(ctx, tc.lang, "8080")

		assert.Equal(t, tc.expErr, err)

		if err == nil {
			testDockerFile(t, tc.content)
		}
	}
}

func testDockerFile(t *testing.T, expected string) {
	t.Helper()

	_, err := os.Stat("DockerFile")
	if err != nil {
		t.Errorf("Test failed, could not create DockerFile, error : %v", err)
	}

	file, err := os.Open("DockerFile")
	if err != nil {
		t.Errorf("Test failed, error opening DockerFile, error : %v", err)
	}

	b, err := io.ReadAll(file)
	if err != nil {
		t.Errorf("Test failed, error reading DockerFile, error : %v", err)
	}

	assert.Equal(t, expected, string(b))
}
