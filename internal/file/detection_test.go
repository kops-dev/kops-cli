package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DetectGolang(t *testing.T) {
	f, err := os.Create("go.mod")
	if err != nil {
		t.Errorf("error creating go.mod file, err : %v", err)
	}

	defer os.Remove(f.Name())

	out := Detect()

	assert.Equal(t, golang, out)
}

func Test_DetectJavaScript(t *testing.T) {
	f, err := os.Create("package.json")
	if err != nil {
		t.Errorf("error creating go.mod file, err : %v", err)
	}

	defer os.Remove(f.Name())

	out := Detect()

	assert.Equal(t, js, out)
}

func Test_DetectJava(t *testing.T) {
	f, err := os.Create("pom.xml")
	if err != nil {
		t.Errorf("error creating go.mod file, err : %v", err)
	}

	defer os.Remove(f.Name())

	out := Detect()

	assert.Equal(t, java, out)
}

func Test_DetectFail(t *testing.T) {
	out := Detect()
	assert.Equal(t, "", out)
}
