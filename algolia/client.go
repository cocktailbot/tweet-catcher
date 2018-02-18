package algolia

import (
	"encoding/json"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

// Config needed to connect to Algolia
type Config struct {
	AppID  string
	APIKey string
}

// Create a client
func Create(c Config) algoliasearch.Client {
	return algoliasearch.NewClient(c.AppID, c.APIKey)
}

// Index json string into the search service
func Index(client algoliasearch.Client, index string, jsn []byte) {
	indx := client.InitIndex(index)

	var object algoliasearch.Object
	if err := json.Unmarshal(jsn, &object); err != nil {
		panic(err)
	}

	_, err := indx.AddObject(object)
	if err != nil {
		panic(err)
	}
}
