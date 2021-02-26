package proc

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

func GetProcCmdLine(pid int) (string, error) {
	pidStr := strconv.Itoa(pid)

	content, err := ioutil.ReadFile("/proc/" + pidStr + "/cmdline")

	if err != nil {
		return "", errors.New("Couldn't find process with pid: " + pidStr)
	}

	return strings.ReplaceAll(string(content), "\000", " "), nil
}
