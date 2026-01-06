package util

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
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

func (self *DeployPrint) SafePrint(format string, a ...any) {
	printLock.Lock()
	defer printLock.Unlock()
	fmt.Printf(format, a...)
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
	const maxStatusLength = 15
	const maxTimeLength = 8

	lastInd := strings.LastIndex(resourceName, "::")
	objType := resourceName[:lastInd]
	objName := resourceName[lastInd+2:]

	self.SafePrint("%v | %v/%v | %v | %v | %v | %v\n", stackName, strconv.Itoa(current), strconv.Itoa(total), padRight(maxTimeLength, time.Now().Format(time.TimeOnly)), padRight(maxStatusLength, status), padRight(self.ObjectTypeLength, objType), objName)
}
