package main

import (
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
	indexTweets := os.Getenv("COCKTAILBOT_ALGOLIA_INDEX_TWEETS")
	algoliaConfig := algolia.Config{
		APIKey: os.Getenv("COCKTAILBOT_ALGOLIA_API_KEY"),
		AppID:  os.Getenv("COCKTAILBOT_ALGOLIA_APP_ID"),
		Env:    os.Getenv("COCKTAILBOT_ENV"),
	}
	algoliaClient := algolia.Create(algoliaConfig)

	config := twitter.Config{
		ConsumerKey:    os.Getenv("COCKTAILBOT_TWITTER_API_KEY"),
		ConsumerSecret: os.Getenv("COCKTAILBOT_TWITTER_SECRET"),
		AccessToken:    os.Getenv("COCKTAILBOT_TWITTER_ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("COCKTAILBOT_TWITTER_ACCESS_TOKEN_SECRET"),
	}
	tc := twitter.Create(config)
	twitter.Stream(tc, keywords, func(tweet []byte) {
		fmt.Print(string(tweet[:]))
		algoliaClient.IndexJSON(indexTweets, tweet)
	})
}

func help() {
	fmt.Println("Access Twitter Streaming API and save tweets containing given keywords")
	fmt.Println("\nExample:")
	fmt.Println("COMMAND keyword1 keyword2")
}
