package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cocktailbot/tweet-filter/algolia"
	"github.com/cocktailbot/tweet-filter/twitter"
)

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}
	keywords := os.Args[1:]

	algoliaConfig := algolia.Config{
		APIKey: os.Getenv("COCKTAILBOT_ALGOLIA_API_KEY"),
		AppID:  os.Getenv("COCKTAILBOT_ALGOLIA_APP_ID"),
	}
	algoliaClient := algolia.Create(algoliaConfig)

	config := twitter.Config{
		ConsumerKey:    os.Getenv("COCKTAILBOT_TWITTER_API_KEY"),
		ConsumerSecret: os.Getenv("COCKTAILBOT_TWITTER_SECRET"),
		AccessToken:    os.Getenv("COCKTAILBOT_TWITTER_ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("COCKTAILBOT_TWITTER_ACCESS_TOKEN_SECRET"),
	}
	client := twitter.Create(config)
	twitter.Stream(client, keywords, func(tweet interface{}) {
		// Convert tweet to json
		jsn, err := json.Marshal(tweet)
		if err != nil {
			panic(err)
		}

		algolia.Index(algoliaClient, "tweets", jsn)
	})
}

func help() {
	fmt.Println("Access Twitter Streaming API and save tweets containing given keywords")
	fmt.Println("\nExample:")
	fmt.Println("COMMAND keyword1 keyword2")
}
