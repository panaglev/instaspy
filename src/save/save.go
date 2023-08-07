package save

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"instaspy/src/config"
	sqlite "instaspy/src/storage"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// check if images folder exists
func Image(username, imagePath string, db *sqlite.Storage) (sqlite.FileInfo, error) {
	const op = "pkg.save.SaveImage"

	cfg := config.MustLoad()
	fileInfo := sqlite.FileInfo{}

	//Filling FileInfo
	fileInfo.Username = username

	// Download image
	resp, err := http.Get(imagePath)
	if err != nil {
		return sqlite.FileInfo{}, fmt.Errorf("Failed to download the image %s: %w\n", op, err)
	}
	defer resp.Body.Close()

	// Create a buffer to cache response body contents
	var bodyBuffer bytes.Buffer
	_, err = io.Copy(&bodyBuffer, resp.Body)
	if err != nil {
		return sqlite.FileInfo{}, fmt.Errorf("Failed to create image buffer at %s: %w\n", op, err)
	}

	fileInfo.Hash, err = returnHash(ioutil.NopCloser(bytes.NewReader(bodyBuffer.Bytes())))
	if err != nil {
		return sqlite.FileInfo{}, fmt.Errorf("Error having hash of resp.Body at %s: %w", op, err)
	}

	status_code := db.CheckHash(fileInfo.Hash)
	if status_code == 409 {
		return sqlite.FileInfo{Hash: "dont"}, nil
	} else if status_code == 200 {
		if resp.StatusCode != http.StatusOK {
			return sqlite.FileInfo{}, fmt.Errorf("Failed to download the image at %s. Status code: %d\n", op, resp.StatusCode)
		}

		// Get current working directory
		currentDir, err := os.Getwd()
		if err != nil {
			return sqlite.FileInfo{}, fmt.Errorf("failed to get current working directory: %w", err)
		}

		// Change working directory to config download path
		if err := os.Chdir(cfg.DownloadPath); err != nil {
			return sqlite.FileInfo{}, fmt.Errorf("failed to change directory: %w", err)
		}
		// After saving file return to current working directory
		defer os.Chdir(currentDir)

		// In current directory crearing user's folder
		os.Mkdir(username, 0755)

		// Change directory to user's
		if err = os.Chdir(username); err != nil {
			return sqlite.FileInfo{}, fmt.Errorf("failed to change directory to users")
		}

		// Calculating current file name
		fileName, err := returnName()
		if err != nil {
			return sqlite.FileInfo{}, fmt.Errorf("Error getting number at %s: %s", op, err)
		}

		fileInfo.Picture_name, err = strconv.Atoi(fileName)
		if err != nil {
			return sqlite.FileInfo{}, fmt.Errorf("Problem converting filename to fileInfo.Picture_name at %s: %w", op, err)
		}

		// Creating empty file with required filename
		outputFile, err := os.Create(fileName + ".jpg")
		if err != nil {
			return sqlite.FileInfo{}, fmt.Errorf("Failed to create the output file %s: %s\n", op, err)
		}
		defer outputFile.Close()

		// Write resp.Body to file
		_, err = io.Copy(outputFile, &bodyBuffer)
		if err != nil {
			fmt.Printf("Failed to save the image: %s\n", err)
		}

		return fileInfo, nil
	} else {
		return sqlite.FileInfo{}, fmt.Errorf("Error %s: %w", op, err)
	}
}

func returnName() (string, error) {
	const op = "pkg.save.returnNumber"

	cmd := exec.Command("ls")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error executing command at %s: %w", op, err)
	}

	lines := strings.Split(string(output), "\n")
	fileCount := len(lines) - 1

	return strconv.Itoa(fileCount), nil
}

func returnHash(fileStream io.ReadCloser) (string, error) {
	const op = "pkg.save.returnHash"

	hash := sha256.New()

	_, err := io.Copy(hash, fileStream)
	if err != nil {
		return "", err
	}

	hashSum := hash.Sum(nil)
	hashSumHex := hex.EncodeToString(hashSum)
	return hashSumHex, nil
}
