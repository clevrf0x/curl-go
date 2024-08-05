package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func printColoredText(text string, colorCode string) {
	fmt.Printf("\033[%sm%s\033[0m", colorCode, text)
}

func main() {
	method := flag.String("method", "GET", "HTTP method (GET, POST, PUT, PATCH, DELETE, OPTIONS)")
	data := flag.String("data", "", "JSON data to send with request")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Please provide a URL as an argument")
		os.Exit(1)
	}

	url := flag.Arg(0)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	var req *http.Request
	var err error

	if *data != "" {
		if !json.Valid([]byte(*data)) {
			fmt.Println("Invalid JSON data provided")
			os.Exit(1)
		}
		req, err = http.NewRequest(strings.ToUpper(*method), url, bytes.NewBufferString(*data))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(strings.ToUpper(*method), url, nil)
	}

	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	printColoredText(fmt.Sprintf("HTTP/%d.%d %s\n", resp.ProtoMajor, resp.ProtoMinor, resp.Status), "1;34")

	for key, values := range resp.Header {
		printColoredText(fmt.Sprintf("%s: ", key), "1;36")
		fmt.Printf("%s\n", strings.Join(values, ", "))
	}
	fmt.Println()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	cmd := exec.Command("jq", ".")
	cmd.Stdin = strings.NewReader(string(body))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running jq:", err)
		fmt.Println("Raw response:")
		fmt.Println(string(body))
	}
}
