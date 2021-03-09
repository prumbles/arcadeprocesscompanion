package proc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
)

type CheckProcessesResult struct {
	KeepAliveProcessFound  bool
	PriorityProcessMatcher string
	PriorityProcessFound   bool
	PID int
}

func GetProcCmdLine(pid int) (string, error) {
	pidStr := strconv.Itoa(pid)

	content, err := ioutil.ReadFile("/proc/" + pidStr + "/cmdline")

	if err != nil {
		return "", errors.New("Couldn't find process with pid: " + pidStr)
	}

	return strings.ReplaceAll(string(content), "\000", " "), nil
}

func CheckProcesses(processPriorityMatchings []string, keepAliveProcessMatcher string) CheckProcessesResult {
	processes, err := ps.Processes()
	processKeepAlive := keepAliveProcessMatcher != "*"

	if err != nil {
		fmt.Printf("WARNING: Failure to retrieve list of processes.\n")
		return CheckProcessesResult{
			false, "", false, 0,
		}
	} else {
		foundKeepAlive := !processKeepAlive
		foundPriorityProcess := false
		priorityProcessMatcher := ""
		foundPID := 0

		for priorityMatcherIndex := range processPriorityMatchings {
			priorityMatcher := processPriorityMatchings[priorityMatcherIndex]

			for i := range processes {
				//skip the currently running process
				if processes[i].Pid() == os.Getpid() {
					continue
				}

				pname, procError := GetProcCmdLine(processes[i].Pid())

				if procError == nil {
					if strings.HasPrefix(pname, "go run") {
						continue
					}
				} else {
					continue
				}

				if foundPriorityProcess && foundKeepAlive {
					return CheckProcessesResult{
						foundKeepAlive, priorityProcessMatcher, foundPriorityProcess, foundPID,
					}
				}

				if !foundPriorityProcess && strings.Contains(pname, priorityMatcher) {
					foundPriorityProcess = true
					priorityProcessMatcher = priorityMatcher
					foundPID = processes[i].Pid()
				}

				if !foundKeepAlive && strings.Contains(pname, keepAliveProcessMatcher) {
					foundKeepAlive = true
				}
			}
		}

		return CheckProcessesResult{
			foundKeepAlive, priorityProcessMatcher, foundPriorityProcess, foundPID, 
		}
	}
}
