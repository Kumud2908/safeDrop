package url

import "fmt"

func GenerateDownloadURL(id string, key string) string {
	return fmt.Sprintf("http://localhost:8080/v/%s/%s", id, key)
}
