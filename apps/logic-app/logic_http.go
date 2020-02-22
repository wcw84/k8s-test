package main
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func getEnv() string {
	return os.Getenv("MODULE_NAME")
}
func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		moduleName := getEnv()
		if moduleName == "" {
			moduleName = "unknown module"
		}
		hostName, _ := os.Hostname()
		time.Now().Date()
		fmt.Printf("[%v][%v] Hello at %v", moduleName, hostName, time.Now())
		fmt.Fprintf(writer, "[%v][%v] Hello at %v", moduleName, hostName, time.Now())
	})
	log.Fatal(http.ListenAndServe(":8080",nil))
}