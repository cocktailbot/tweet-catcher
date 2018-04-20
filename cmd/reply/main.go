package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
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
		fields := []string{"url", "title", "search"}
		text := el["text"].(string)
		text = stopwords.CleanString(text, "en", true)
		author := el["user"].(map[string]interface{})["screen_name"].(string)
		id, _ := strconv.ParseInt(el["id_str"].(string), 10, 64)
		objectID := el["objectID"].(string)

		matches := recipesClient.Search(indexRecipes, fields, text, 0, 20)
		message := createMessage(author, matches)
		twitter.Reply(twitterClient, message, id)

		replied = append(replied, objectID)
	}

	if len(replied) > 0 {
		fmt.Printf("Deleting #%v\n", replied)
		tweetsClient.DeleteByIds(indexTweets, replied)
	}
}

func createMessage(author string, matches []algoliasearch.Map) string {
	message := randomNoMatchReply(author)
	if len(matches) > 0 {
		match := matches[rand.Intn(len(matches)-1)]
		title := match["title"].(string)
		url := match["url"].(string)
		message = randomMatchReply(author, title, url)
		fmt.Println(message + " in " + match["search"].(string))
	}
	return message
}

func randomMatchReply(username string, title string, url string) string {
	replies := []string{
		"Hey @%s, try %s %s",
		"Hi @%s, how about %s? %s",
		"Ok @%s, checkout %s %s",
	}

	reply := replies[rand.Intn(len(replies)-1)]

	return fmt.Sprintf(reply, username, title, url)
}

func randomNoMatchReply(username string) string {
	replies := []string{
		"Well @%s, we couldn't find a match 🧐. Try some other ingredients!",
		"Sorry @%s, there's no matches 🤭. Why not try some different ingredients?",
		"Hmm @%s, that doesn't seem to match 🤨. How about trying with different ingredients?",
	}
	reply := replies[rand.Intn(len(replies)-1)]

	return fmt.Sprintf(reply, username)
}
