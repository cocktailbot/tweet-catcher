package algolia

import (
	"encoding/json"

	"fmt"

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
	Config     Config
}

// Create a client
func Create(c Config) Client {
	return Client{Connection: algoliasearch.NewClient(c.AppID, c.APIKey), Config: c}
}

// IndexJSON json string into the search service
func (c Client) IndexJSON(index string, jsn []byte) {
	idx := c.Connection.InitIndex(indexPrefix(c.Config.Env) + index)

	var object algoliasearch.Object
	if err := json.Unmarshal(jsn, &object); err != nil {
		panic(err)
	}

	_, err := idx.AddObject(object)
	if err != nil {
		panic(err)
	}
}

// Search a given index for items matching query
func (c Client) Search(index string, fields []string, sentence string, page int, perPage int) []algoliasearch.Map {
	idx := c.Connection.InitIndex(indexPrefix(c.Config.Env) + index)
	params := algoliasearch.Map{
		"attributesToRetrieve": fields,
		"hitsPerPage":          perPage,
		"page":                 page,
		"removeWordsIfNoResults": "allOptional",
		"typoTolerance":          false,
	}
	res, err := idx.Search(sentence, params)

	if err != nil {
		panic(err)
	}

	return res.Hits
}

// DeleteByIds remove objects matching given ids
func (c Client) DeleteByIds(index string, ids []string) {
	idx := c.Connection.InitIndex(indexPrefix(c.Config.Env) + index)
	res, err := idx.DeleteObjects(ids)

	if err != nil {
		panic(err)
	}

	err = idx.WaitTask(res.TaskID)

	if err != nil {
		panic(err)
	}
}

func indexPrefix(env string) string {
	return fmt.Sprintf("%s_", env)
}
