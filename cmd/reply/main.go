package main

import (
	"os"

	"fmt"

	"time"

	"strconv"

	"github.com/bbalet/stopwords"
	"github.com/cocktailbot/tweet-filter/algolia"
	"github.com/cocktailbot/tweet-filter/twitter"
)

func main() {
	tweetsClient := algolia.Create(algolia.Config{
		APIKey: os.Getenv("COCKTAILBOT_ALGOLIA_API_KEY"),
		AppID:  os.Getenv("COCKTAILBOT_ALGOLIA_APP_ID"),
		Env:    os.Getenv("COCKTAILBOT_ENV"),
	})
	recipesClient := algolia.Create(algolia.Config{
		APIKey: os.Getenv("COCKTAILBOT_RECIPES_ALGOLIA_API_KEY"),
		AppID:  os.Getenv("COCKTAILBOT_RECIPES_ALGOLIA_APP_ID"),
		Env:    os.Getenv("COCKTAILBOT_ENV"),
	})
	twitterClient := twitter.Create(twitter.Config{
		ConsumerKey:    os.Getenv("COCKTAILBOT_TWITTER_API_KEY"),
		ConsumerSecret: os.Getenv("COCKTAILBOT_TWITTER_SECRET"),
		AccessToken:    os.Getenv("COCKTAILBOT_TWITTER_ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("COCKTAILBOT_TWITTER_ACCESS_TOKEN_SECRET"),
	})
	recent := tweetsClient.Search(
		"TWEETS",
		[]string{"objectID", "user.screen_name", "id_str", "text"},
		"",
		0,
		10)
	var replied []string

	for _, el := range recent {
		author := el["user"].(map[string]interface{})["screen_name"].(string)
		text := el["text"].(string)
		id, _ := strconv.ParseInt(el["id_str"].(string), 10, 64)
		objectID := el["objectID"].(string)
		text = stopwords.CleanString(text, "en", true)
		match := recipesClient.Search("local_recipes", []string{"url", "title", "search"}, text, 0, 1)

		if len(match) > 0 {
			// match["url"]
			ts := time.Now().Format(time.RFC850)
			message := fmt.Sprintf("Hey @%s, try %s %s %s", author, match[0]["title"].(string), match[0]["url"].(string), ts)
			fmt.Println(message + " in " + match[0]["search"].(string))
			twitter.Reply(twitterClient, message, id)
		}
		replied = append(replied, objectID)
	}

	if len(replied) > 0 {
		fmt.Printf("Deleting #%v\n", replied)
		//tweetsClient.DeleteByIds("TWEETS", replied)
	}
}
