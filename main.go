package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	projectId := os.Getenv("INFURA_PROJECT_ID")
	projectSecret := os.Getenv("INFURA_PROJECT_SECRET")
	gateway := os.Getenv("IPFS_GATEWAY")

	if len(projectId) == 0 || len(projectSecret) == 0 {
		log.Fatal("Provide Infura project information as INFURA_PROJECT_ID and INFURA_PROJECT_SECRET")
	}

	if len(os.Args) < 2 {
		log.Fatal("Provide a filepath to be uploaded")
	}

	path := os.Args[1]

	shell := shell.NewShellWithClient("https://ipfs.infura.io:5001", &http.Client{
		Transport: authTransport{
			RoundTripper:  http.DefaultTransport,
			ProjectId:     projectId,
			ProjectSecret: projectSecret,
		},
	})

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	var cid string
	if fileInfo.IsDir() {
		cid, err = shell.AddDir(path)
		if err != nil {
			log.Fatal(err)
		}
	} else {
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

	err = shell.Pin(cid)
	if err != nil {
		log.Fatal(err)
	}

	if len(gateway) == 0 {
		gateway = "https://ipfs.io"
	}

	if !strings.HasPrefix(gateway, "https://") {
		gateway = fmt.Sprintf("https://%s", gateway)
	}

	log.Print(fmt.Sprintf("Pinned: %s/ipfs/%s", gateway, cid))

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
