package algolia

import (
	"encoding/json"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

// Config of connection details
type Config struct {
	AppID  string
	APIKey string
	Env    string
}

// Client for connecting
type Client struct {
	Connection algoliasearch.Client
}

// Create a client
func Create(c Config) Client {
	return Client{Connection: algoliasearch.NewClient(c.AppID, c.APIKey)}
}

// IndexJSON json string into the search service
func (c Client) IndexJSON(index string, jsn []byte) {
	indx := c.Connection.InitIndex(index)

	var object algoliasearch.Object
	if err := json.Unmarshal(jsn, &object); err != nil {
		panic(err)
	}

	_, err := indx.AddObject(object)
	if err != nil {
		panic(err)
	}
}

// Search a given index for items matching query
func (c Client) Search(index string, sentence string, page int, perPage int) []algoliasearch.Map {
	idx := c.Connection.InitIndex(index)
	params := algoliasearch.Map{
		"hitsPerPage": perPage,
		"page":        page,
		"removeWordsIfNoResults": "allOptional",
	}
	res, err := idx.Search(sentence, params)

	if err != nil {
		panic(err)
	}

	return res.Hits
}

// DeleteByIds remove objects matching given ids
func (c Client) DeleteByIds(index string, ids []string) {
	idx := c.Connection.InitIndex(index)
	idx.DeleteObjects(ids)
}
