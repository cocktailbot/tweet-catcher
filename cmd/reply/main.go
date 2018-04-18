package reply

import (
	"os"

	"strconv"

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
	recent := tweetsClient.Search("tweets", "", 1, 10)
	var replied []string

	for _, el := range recent {
		text := el["text"].(string)
		id := el["id"].(string)
		match := recipesClient.Search("local_recipes", text, 1, 1)

		if len(match) > 0 {
			// match["url"]
			message := "Found match " + match[0]["title"].(string)
			replyID, _ := strconv.ParseInt(id, 10, 64)
			twitter.Reply(twitterClient, message, replyID)
		}
		replied = append(replied, id)
	}

	if len(replied) > 0 {
		tweetsClient.DeleteByIds("tweets", replied)
	}
}
