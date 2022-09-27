package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/joho/godotenv"
)

type Result struct {
	URL      string `json:"url"`
	IPFS_URL string `json:"ipfs_url"`
	CID      string `json:"cid"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Get Infura secrets and IPFS config
	apiEndpoint := getEnv("IPFS_API_ENDPOINT", "https://ipfs.infura.io:5001")
	gateway := getEnv("IPFS_GATEWAY", "https://ipfs.io")
	projectId := os.Getenv("INFURA_PROJECT_ID")
	projectSecret := os.Getenv("INFURA_PROJECT_SECRET")

	//Infura secrets cannot be empty
	if len(projectId) == 0 || len(projectSecret) == 0 {
		log.Fatal("Provide Infura project information as INFURA_PROJECT_ID and INFURA_PROJECT_SECRET")
	}

	//Must provide a path as arg
	if len(os.Args) < 2 {
		log.Fatal("Provide a filepath to be uploaded")
	}

	//Connect to IPFS API, use transport wrapper to be able to add basic auth (based on Infura docs)
	shell := shell.NewShellWithClient(apiEndpoint, &http.Client{
		Transport: authTransport{
			RoundTripper:  http.DefaultTransport,
			ProjectId:     projectId,
			ProjectSecret: projectSecret,
		},
	})

	path := os.Args[1]

	//Get path info
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	var cid string

	//Use AddDir if the path is a directory
	if fileInfo.IsDir() {
		cid, err = shell.AddDir(path)
		if err != nil {
			log.Fatal(err)
		}
	} else { //Upload single file otherwise
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
		cid, err = shell.Add(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	//Pin the file to prevent garbage collection
	err = shell.Pin(cid)
	if err != nil {
		log.Fatal(err)
	}

	//Make sure the gateway has correct protocol set
	if !strings.HasPrefix(gateway, "https://") {
		gateway = fmt.Sprintf("https://%s", gateway)
	}

	//Generate serializable output
	result := Result{
		URL:      fmt.Sprintf("%s/ipfs/%s", gateway, cid),
		IPFS_URL: fmt.Sprintf("ipfs://%s", cid),
		CID:      cid,
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Print(cid)
		log.Fatal(err)
	}

	fmt.Println(string(output))
}

// authTransport decorates each request with a basic auth header.
type authTransport struct {
	http.RoundTripper
	ProjectId     string
	ProjectSecret string
}

func (t authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(t.ProjectId, t.ProjectSecret)
	return t.RoundTripper.RoundTrip(r)
}

func getEnv(key string, _default string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return _default
}
