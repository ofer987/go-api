package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

const DEBUG_MESSAGE = "DEBUG"

type HelloWorld struct {
	Message []string `json:"message"`
}

type FullerWorld struct {
	GivenNames []string   `json:"given_names"`
	HelloWorld HelloWorld `json:"hello_world"`
	SurName    string     `json:"sur_name"`
}

var helloRegexp = regexp.MustCompile("Hello")

func findPattern(filepath string) []byte {
	fmt.Printf("filepath: %v\n", filepath)

	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result := helloRegexp.Find(b)
	if result == nil {
		fmt.Printf("Failed to find %s\n", helloRegexp)

		// Empty result
		return []byte{}
	}

	return append([]byte{}, result...)
}

type responseRecorder struct {
	http.ResponseWriter
	body   []byte
	status int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)

	return 0, nil
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func printHelloDanBeforeAndOferAfter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(writer, r)

		var helloWorld HelloWorld
		jsonDecoder := json.NewDecoder(bytes.NewReader(writer.body))
		if err := jsonDecoder.Decode(&helloWorld); err != nil {
			fmt.Printf("Failed to decode bytes into 'HelloWorld' type\n")
			fmt.Println(helloWorld)

			return
		}

		fullerWorld := FullerWorld{
			GivenNames: []string{"Dan"},
			HelloWorld: helloWorld,
			SurName:    "Ofer",
		}

		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.Encode(fullerWorld)
	})
}

func main() {
	result := findPattern("test.txt")
	resultLength := len(result)
	message := fmt.Sprintf("<div>Hello, World! %s with length %d</div>", result, resultLength)
	messages := append([]string{}, message)
	helloWorldResult := HelloWorld{
		Message: messages,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.Encode(helloWorldResult)
	})

	fmt.Printf("[%s] Found '%s'\n", DEBUG_MESSAGE, result)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", printHelloDanBeforeAndOferAfter(http.DefaultServeMux))
}
