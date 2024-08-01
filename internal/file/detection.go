package file

import "os"

const (
	golang = "golang"
	java   = "java"
	js     = "js"
)

func Detect() string {
	switch {
	case checkFile("go.mod"):
		return golang
	case checkFile("package.json"):
		return js
	case checkFile("pom.xml") || checkFile("build.gradle"):
		return java
	}

	return ""
}

func checkFile(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		return false
	}

	return true
}
