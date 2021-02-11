package apps

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dokku/dokku/plugins/common"
)

// checks if an app exists
func appExists(appName string) error {
	return common.VerifyAppName(appName)
}

// checks if an app is locked
func appIsLocked(appName string) bool {
	lockfilePath := fmt.Sprintf("%v/.deploy.lock", common.AppRoot(appName))
	_, err := os.Stat(lockfilePath)
	return !os.IsNotExist(err)
}

// verifies app name and creates an app
func createApp(appName string) error {
	if err := common.IsValidAppName(appName); err != nil {
		return err
	}

	if err := appExists(appName); err == nil {
		return errors.New("Name is already taken")
	}

	common.LogInfo1Quiet(fmt.Sprintf("Creating %s...", appName))
	os.MkdirAll(common.AppRoot(appName), 0755)

	if err := common.PlugnTrigger("post-create", []string{appName}...); err != nil {
		return err
	}

	return nil
}

// destroys an app
func destroyApp(appName string) error {
	if os.Getenv("DOKKU_APPS_FORCE_DELETE") != "1" {
		if err := common.AskForDestructiveConfirmation(appName, "app"); err != nil {
			return err
		}
	}

	common.LogInfo1(fmt.Sprintf("Destroying %s (including all add-ons)", appName))

	imageTag, _ := common.GetRunningImageTag(appName)
	if err := common.PlugnTrigger("pre-delete", []string{appName, imageTag}...); err != nil {
		return err
	}

	scheduler := common.GetAppScheduler(appName)
	removeContainers := "true"
	if err := common.PlugnTrigger("scheduler-stop", []string{scheduler, appName, removeContainers}...); err != nil {
		return err
	}
	if err := common.PlugnTrigger("scheduler-post-delete", []string{scheduler, appName, imageTag}...); err != nil {
		return err
	}
	if err := common.PlugnTrigger("post-delete", []string{appName, imageTag}...); err != nil {
		return err
	}

	forceCleanup := true
	common.DockerCleanup(appName, forceCleanup)

	// remove contents for apps that are symlinks to other folders
	if err := os.RemoveAll(fmt.Sprintf("%v/", common.AppRoot(appName))); err != nil {
		common.LogWarn(err.Error())
	}

	// then remove the folder and/or the symlink
	if err := os.RemoveAll(common.AppRoot(appName)); err != nil {
		common.LogWarn(err.Error())
	}

	return nil
}

func listImagesByAppLabel(appName string) ([]string, error) {
	command := []string{
		common.DockerBin(),
		"image",
		"list",
		"--quiet",
		"--filter",
		fmt.Sprintf("label=com.dokku.app-name=%v", appName),
	}

	var stderr bytes.Buffer
	listCmd := common.NewShellCmd(strings.Join(command, " "))
	listCmd.ShowOutput = false
	listCmd.Command.Stderr = &stderr
	b, err := listCmd.Output()

	if err != nil {
		return []string{}, errors.New(strings.TrimSpace(stderr.String()))
	}

	output := strings.Split(strings.TrimSpace(string(b[:])), "\n")
	return output, nil
}

func listImagesByImageRepo(imageRepo string) ([]string, error) {
	command := []string{
		common.DockerBin(),
		"image",
		"list",
		"--quiet",
		imageRepo,
	}

	var stderr bytes.Buffer
	listCmd := common.NewShellCmd(strings.Join(command, " "))
	listCmd.ShowOutput = false
	listCmd.Command.Stderr = &stderr
	b, err := listCmd.Output()

	if err != nil {
		return []string{}, errors.New(strings.TrimSpace(stderr.String()))
	}

	output := strings.Split(strings.TrimSpace(string(b[:])), "\n")
	return output, nil
}

// creates an app if allowed
func maybeCreateApp(appName string) error {
	if err := appExists(appName); err == nil {
		return nil
	}

	b, _ := common.PlugnTriggerOutput("config-get-global", []string{"DOKKU_DISABLE_APP_AUTOCREATION"}...)
	disableAutocreate := strings.TrimSpace(string(b[:]))
	if disableAutocreate == "true" {
		common.LogWarn("App auto-creation disabled.")
		return fmt.Errorf("Re-enable app auto-creation or create an app with 'dokku apps:create %s'", appName)
	}

	return common.SuppressOutput(func() error {
		return createApp(appName)
	})
}
