package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/bbalet/stopwords"
	"github.com/cocktailbot/tweet-filter/algolia"
	"github.com/cocktailbot/tweet-filter/twitter"
)

func main() {
	indexTweets := os.Getenv("COCKTAILBOT_ALGOLIA_INDEX_TWEETS")
	indexRecipes := os.Getenv("COCKTAILBOT_ALGOLIA_INDEX_RECIPES")
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
		indexTweets,
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
		fields := []string{"url", "title", "search"}
		matches := recipesClient.Search(indexRecipes, fields, text, 0, 5)

		if len(matches) > 0 {
			match := matches[rand.Intn(len(matches)-1)]
			title := match["title"].(string)
			url := match["url"].(string)
			message := fmt.Sprintf("Hey @%s, try %s %s", author, title, url)
			fmt.Println(message + " in " + match["search"].(string))
			twitter.Reply(twitterClient, message, id)
		} else {
			fmt.Println("No match: " + text)
			message := fmt.Sprintf(
				"Well @%s, we couldn't find a match 🧐. Try some other ingredients!",
				author)
			twitter.Reply(twitterClient, message, id)
		}
		replied = append(replied, objectID)
	}

	if len(replied) > 0 {
		fmt.Printf("Deleting #%v\n", replied)
		tweetsClient.DeleteByIds(indexTweets, replied)
	}
}
