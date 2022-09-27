# go-ipfs-upload
Simple go experiment to try uploading content to IPFS

## Config

* `INFURA_PROJECT_ID` - Infura Project ID (mandatory)
* `INFURA_PROJECT_SECRET` - Infura Project/API Secret (mandatory)
* `IPFS_GATEWAY` - IPFS gateway to be used in output URL (optional, default: `https://ipfs.io`)
* `IPFS_API_ENDPOINT` - IPFS API endpoint to be used for upload (optional, default: `https://ipfs.infura.io:5001`) 

## Run

```
cp .env.example .env
#edit .env
go run main.go $PATH_TO_FILE_OR_DIR
```

Example output: 

```
{
  "url": "https://ipfs.io/ipfs/QmXU1Qd6aPpavzPMMNwLg4cxvqcPrZC1zjVCyvB3o3o62w",
  "ipfs_url": "ipfs://QmXU1Qd6aPpavzPMMNwLg4cxvqcPrZC1zjVCyvB3o3o62w",
  "cid": "QmXU1Qd6aPpavzPMMNwLg4cxvqcPrZC1zjVCyvB3o3o62w"
}

```
