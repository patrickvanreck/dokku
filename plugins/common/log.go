package common

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// LogFail is the failure log formatter
// prints text to stderr and exits with status 1
func LogFail(text string) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(" !     %s", text))
	os.Exit(1)
}

// LogFailQuiet is the failure log formatter (with quiet option)
// prints text to stderr and exits with status 1
func LogFailQuiet(text string) {
	if os.Getenv("DOKKU_QUIET_OUTPUT") == "" {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(" !     %s", text))
	}
	os.Exit(1)
}

// Log is the log formatter
func Log(text string) {
	fmt.Println(text)
}

// LogQuiet is the log formatter (with quiet option)
func LogQuiet(text string) {
	if os.Getenv("DOKKU_QUIET_OUTPUT") == "" {
		fmt.Println(text)
	}
}

// LogInfo1 is the info1 header formatter
func LogInfo1(text string) {
	fmt.Println(fmt.Sprintf("-----> %s", text))
}

// LogInfo1Quiet is the info1 header formatter (with quiet option)
func LogInfo1Quiet(text string) {
	if os.Getenv("DOKKU_QUIET_OUTPUT") == "" {
		LogInfo1(text)
	}
}

// LogInfo2 is the info2 header formatter
func LogInfo2(text string) {
	fmt.Println(fmt.Sprintf("=====> %s", text))
}

// LogInfo2Quiet is the info2 header formatter (with quiet option)
func LogInfo2Quiet(text string) {
	if os.Getenv("DOKKU_QUIET_OUTPUT") == "" {
		LogInfo2(text)
	}
}

// LogVerbose is the verbose log formatter
// prints indented text to stdout
func LogVerbose(text string) {
	fmt.Println(fmt.Sprintf("       %s", text))
}

// LogVerboseQuiet is the verbose log formatter
// prints indented text to stdout (with quiet option)
func LogVerboseQuiet(text string) {
	if os.Getenv("DOKKU_QUIET_OUTPUT") == "" {
		LogVerbose(text)
	}
}

// LogVerboseQuietContainerLogs is the verbose log formatter for container logs
func LogVerboseQuietContainerLogs(containerID string) {
	sc := NewShellCmdWithArgs(DockerBin(), "container", "logs", containerID)
	sc.ShowOutput = false
	b, err := sc.CombinedOutput()
	if err != nil {
		LogExclaim(fmt.Sprintf("Failed to fetch container logs: %s", containerID))
	}

	output := strings.TrimSpace(string(b))
	if len(output) == 0 {
		return
	}

	for _, line := range strings.Split(output, "\n") {
		if line != "" {
			LogVerboseQuiet(line)
		}
	}
}

// LogVerboseQuietContainerLogsTail is the verbose log formatter for container logs with tail mode enabled
func LogVerboseQuietContainerLogsTail(containerID string, lines int, tail bool) {
	args := []string{"container", "logs", containerID}
	if lines > 0 {
		args = append(args, "--tail", strconv.Itoa(lines))
	}
	if tail {
		args = append(args, "--follow")
	}
	sc := NewShellCmdWithArgs(DockerBin(), args...)
	stdout, _ := sc.Command.StdoutPipe()
	sc.Command.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		LogVerboseQuiet(m)
	}
	sc.Command.Wait()
}

// LogWarn is the warning log formatter
func LogWarn(text string) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(" !     %s", text))
}

// LogExclaim is the log exclaim formatter
func LogExclaim(text string) {
	fmt.Println(fmt.Sprintf(" !     %s", text))
}

// LogStderr is the stderr log formatter
func LogStderr(text string) {
	fmt.Fprintln(os.Stderr, text)
}

// LogDebug is the debug log formatter
func LogDebug(text string) {
	if os.Getenv("DOKKU_TRACE") == "1" {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(" ?     %s", strings.TrimPrefix(text, " ?     ")))
	}
}
