package save

import (
	"fmt"
	"instaspy/pkg/config"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
Load image -> grab metrics -> go to username folder -> save image
*/
func Image(username, imagePath string) error {
	const op = "pkg.save.SaveImage"

	cfg := config.MustLoad()

	resp, err := http.Get(imagePath)
	if err != nil {
		return fmt.Errorf("Failed to download the image %s: %w\n", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to download the image at %s. Status code: %d\n", op, resp.StatusCode)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	if err := os.Chdir(cfg.DownloadPath); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}
	defer os.Chdir(currentDir)

	os.Mkdir(username, 0755)

	if err = os.Chdir(username); err != nil {
		return fmt.Errorf("failed to change directory to users")
	}

	fileName, err := returnName()
	if err != nil {
		return fmt.Errorf("Error getting number at %s: %s", op, err)
	}

	outputFile, err := os.Create(fileName + ".jpg")
	if err != nil {
		return fmt.Errorf("Failed to create the output file %s: %s\n", op, err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		fmt.Printf("Failed to save the image: %s\n", err)
	}

	return nil
}

func returnName() (string, error) {
	const op = "pkg.save.returnNumber"

	cmd := exec.Command("ls")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error executing command at %s: %s", op, err)
	}

	lines := strings.Split(string(output), "\n")

	fileCount := len(lines) - 1

	return strconv.Itoa(fileCount), nil
}
