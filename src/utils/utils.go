package utils

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
)

func CheckOSType() string {
	return runtime.GOOS
}

func CheckCPUCount() int {
	return runtime.NumCPU()
}

func CheckWSL2() bool {
	// WSL2 typically sets a WSL2 or WSL_DISTRO_NAME environment variable
	_, wsl2Exists := os.LookupEnv("WSL2")
	_, wslDistroExists := os.LookupEnv("WSL_DISTRO_NAME")

	return wsl2Exists || wslDistroExists
}

// extractVersion uses regular expression to find the version number in the given string
func ExtractVersion(str string) (string, error) {
	// Define a regular expression pattern for the version number
	// This pattern looks for sequences like 1.0.3, 10.23.456 etc.
	re := regexp.MustCompile(`\d+\.\d+\.\d+`)

	// Find the first match in the string
	matches := re.FindStringSubmatch(str)
	if len(matches) > 0 {
		return matches[0], nil
	}

	return "", fmt.Errorf("no version number found")
}

// ReplaceValuesInFile reads a file line by line and replaces the first occurrence of specified key's values.
func ReplaceValuesInFile(filePath string, replacements map[string]string) error {
	// Open the file for reading.
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Track which keys have been replaced.
	replaced := make(map[string]bool)
	for key := range replacements {
		replaced[key] = false
	}

	// Create a temporary buffer to store the modified content.
	var buffer strings.Builder

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains a key that needs to be replaced.
		for key, newValue := range replacements {
			//do not preserve indentation
			lineNoIndent := strings.TrimSpace(line)
			// Enable line if it was commented out
			if !replaced[key] && (strings.HasPrefix(lineNoIndent, key+" =") || strings.HasPrefix(lineNoIndent, "# "+key+" =") || strings.HasPrefix(lineNoIndent, "#"+key+" =")) {
				// Replace the line with the new value.
				line = key + " = \"" + newValue + "\""
				replaced[key] = true
				break
			}
		}

		buffer.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Write the modified content back to the file.
	return os.WriteFile(filePath, []byte(buffer.String()), 0644)
}

// DownloadFile downloads a file from the specified URL and saves it to the specified local path.
// It overwrites the file if it exists and returns the SHA256 checksum of the downloaded file.
func DownloadFile(url, filePath string) (string, error) {
	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file, overwriting it if it already exists
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a SHA256 hasher
	hasher := sha256.New()

	// Create a multi-writer to write to both file and hasher
	multiWriter := io.MultiWriter(file, hasher)

	// Write the response's body to the file and hasher
	if _, err = io.Copy(multiWriter, resp.Body); err != nil {
		return "", err
	}

	// Compute and return the SHA256 checksum in hexadecimal format
	checksum := hex.EncodeToString(hasher.Sum(nil))
	return checksum, nil
}

// IsValidURL tests a string to determine if it is a well-structured URL or not.
func IsValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func WriteError(err string) {
	fmt.Fprintln(os.Stderr, err)
}

// IPResponse represents the structure of the response from httpbin.org
type IPResponse struct {
	Origin string `json:"origin"`
}

// GetExternalIP fetches the external IP address of the current machine
func GetExternalIP() (string, error) {
	response, err := http.Get("https://httpbin.org/ip")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var ipResp IPResponse
	err = json.Unmarshal(body, &ipResp)
	if err != nil {
		return "", err
	}

	return ipResp.Origin, nil
}
