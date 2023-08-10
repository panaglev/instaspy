package save

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"instaspy/src/config"
	"instaspy/src/logger"
	sqlite "instaspy/src/storage"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Image(username, imagePath string, db *sqlite.Storage) (sqlite.FileInfo, error) {
	const op = "src.save.SaveImage"

	// Load config and structure to return
	cfg := config.MustLoad()
	fileInfo := sqlite.FileInfo{}
	//Filling FileInfo
	fileInfo.Username = username

	// Download image
	resp, err := http.Get(imagePath)
	if err != nil {
		logger.HandleOpError(op, err)
		return sqlite.FileInfo{}, err
	}
	defer resp.Body.Close()

	// Create a buffer to cache response body contents
	var bodyBuffer bytes.Buffer
	_, err = io.Copy(&bodyBuffer, resp.Body)
	if err != nil {
		logger.HandleOpError(op, err)
		return sqlite.FileInfo{}, err
	}

	// Calculate hash to check if image already downloaded or not
	fileInfo.Hash, err = returnHash(ioutil.NopCloser(bytes.NewReader(bodyBuffer.Bytes())))
	if err != nil {
		logger.HandleOpError(op, err)
		return sqlite.FileInfo{}, err
	}

	// Check calculated hash
	status_code, err := db.CheckHash(fileInfo.Hash)
	if err != nil {
		logger.HandleOpError(op, err)
		return sqlite.FileInfo{}, err
	}

	// Unparse response from db
	if status_code == true {
		return sqlite.FileInfo{}, nil
	} else {
		if resp.StatusCode != http.StatusOK {
			logger.HandleOpError(op, err)
			return sqlite.FileInfo{}, err
		}

		// Get current working directory
		currentDir, err := os.Getwd()
		if err != nil {
			logger.HandleOpError(op, err)
			return sqlite.FileInfo{}, err
		}

		// Change working directory to config download path
		if err := os.Chdir(cfg.DownloadPath); err != nil {
			logger.HandleOpError(op, err)
			return sqlite.FileInfo{}, err
		}
		// After saving file return to current working directory
		defer os.Chdir(currentDir)

		// In current directory crearing user's folder
		os.Mkdir(username, 0755)

		// Change directory to user's
		if err = os.Chdir(username); err != nil {
			logger.HandleOpError(op, err)
			return sqlite.FileInfo{}, err
		}

		// Calculating current file name
		fileName, err := returnName()
		if err != nil {
			logger.HandleOpError(op, err)
			return sqlite.FileInfo{}, err
		}

		fileInfo.Picture_name = fileName

		// Creating empty file with required filename
		outputFile, err := os.Create(fileName + ".jpg")
		if err != nil {
			logger.HandleOpError(op, err)
			return sqlite.FileInfo{}, err
		}
		defer outputFile.Close()

		// Write resp.Body to file
		_, err = io.Copy(outputFile, &bodyBuffer)
		if err != nil {
			logger.HandleOpError(op, err)
			return sqlite.FileInfo{}, err
		}

		return fileInfo, nil
	}
}

func returnName() (string, error) {
	const op = "src.save.returnNumber"

	cmd := exec.Command("ls")
	output, err := cmd.Output()
	if err != nil {
		logger.HandleOpError(op, err)
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	fileCount := len(lines) - 1

	return strconv.Itoa(fileCount), nil
}

func returnHash(fileStream io.ReadCloser) (string, error) {
	const op = "src.save.returnHash"

	hash := sha256.New()

	_, err := io.Copy(hash, fileStream)
	if err != nil {
		logger.HandleOpError(op, err)
		return "", err
	}

	hashSum := hash.Sum(nil)
	hashSumHex := hex.EncodeToString(hashSum)
	return hashSumHex, nil
}
