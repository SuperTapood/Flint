package util

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
)

var printLock sync.Mutex

type DeployPrint struct {
	CompletionLength int
	ObjectTypeLength int
}

func CreateDeployPrint(stackName string, prettyNames []string, stackMetadata map[string]any) DeployPrint {
	var maxTypeLength int

	for _, name := range prettyNames {
		lastInd := strings.LastIndex(name, "::")
		objType := name[:lastInd]
		if len(objType) > maxTypeLength {
			maxTypeLength = len(objType)
		}
	}

	return DeployPrint{
		CompletionLength: 1 + (2 * int(math.Log10(float64(len(prettyNames))))),
		ObjectTypeLength: maxTypeLength + 1,
	}
}

func padRight(length int, str string) string {
	if len(str) >= length {
		return str
	}
	// Calculate how many spaces are needed
	padding := length - len(str)
	return str + fmt.Sprintf("%*s", padding, "")
}

func (self *DeployPrint) PrettyPrint(stackName string, current int, total int, status string, resourceName string) {
	printLock.Lock()
	defer printLock.Unlock()

	const maxStatusLength = 10

	lastInd := strings.LastIndex(resourceName, "::")
	objType := resourceName[:lastInd]
	objName := resourceName[lastInd+2:]

	fmt.Printf("%v | %v/%v | %v | %v | %v\n", stackName, strconv.Itoa(current), strconv.Itoa(total), padRight(maxStatusLength, status), padRight(self.ObjectTypeLength, objType), objName)
}
