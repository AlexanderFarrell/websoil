package media

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"github.com/AlexanderFarrell/websoil/core"
)

const (
	defaultUploadFolder = "./uploads"
	defaultGlobalMedia  = "false"
)

var (
	uploadFolder string
	globalMedia  bool
)

func InitEnvVariables(shareMedia bool) {

	uploadFolder = web.GetEnvVar("UPLOAD_FOLDER", defaultUploadFolder)
	globalMedia = shareMedia

	if !CheckIfExists(uploadFolder) && uploadFolder != defaultUploadFolder {
		panic("Set upload folder: " + uploadFolder + " does not exist")
	}
}

func DeleteFolder(username string, path string) error {
	path = SanitizedForFileIO(path)
	base := baseFolder(username)
	p := filepath.Join(base, path)
	err := os.RemoveAll(p)
	return err
}

func SanitizedForFileIO(s string) string {
	s = strings.Replace(s, "..", "", -1)
	//s = strings.Replace(s, "/", "", -1)
	//s = strings.Replace(s, "\\", "", -1)
	return s
}

func UserPath(username string, rel string) (string, error) {
	rel = strings.TrimSpace(rel)

	if strings.Contains(rel, "\x00") {
		return "", fmt.Errorf("invalid path")
	}

	cleanRel := filepath.Clean("/" + rel)
	cleanRel = strings.TrimPrefix(cleanRel, string(os.PathSeparator))

	base := baseFolder(username)
	full := filepath.Clean(filepath.Join(base, cleanRel))

	if full != base && !strings.HasPrefix(full, base+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid path")
	}

	return full, nil
}

func CheckIfExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return false
}

func baseFolder(username string) string {
	if globalMedia {
		return filepath.Clean(uploadFolder)
	}
	username = SanitizedForFileIO(username)
	return filepath.Clean(filepath.Join(uploadFolder, username))
}

func CalculateEtag(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetAt() {
	
}
