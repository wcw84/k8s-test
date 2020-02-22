package main
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)
const (
	url = "http://logic:8081"
)
func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Printf("receive request at %v\n", time.Now())
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("error")
			fmt.Fprintf(writer, "something is error")
		} else {
			data, err := ioutil.ReadAll(resp.Body)
			s := string(data)
			if err != nil {
				fmt.Println("body error")
			} else {
				fmt.Printf("response %v\n", s)
				fmt.Fprintf(writer, "%v", s)

			}
		}
	})
	log.Fatal(http.ListenAndServe(":8080",nil))
}