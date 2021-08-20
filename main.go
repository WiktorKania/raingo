package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	SpartathlonID int = 865167211944345600
	session       *discordgo.Session
	dbClient      *mongo.Client
)

const (
	boomerJoke = "Przychodzi facet do jasnowidzki.\n- Dzień dobry, Kamilu.\n- Ale ja nie jestem Kamil.\n- Wiem."
)

func tellJoke(session *discordgo.Session, msg *discordgo.MessageCreate, joke string) {
	session.ChannelMessageSend(msg.ChannelID, joke)
}

func replyToChannel(channelID string, msg string) {
	session.ChannelMessageSend(channelID, msg)
}

func findLatestMessage(guildID string) {
	channels, err := session.GuildChannels(guildID)
	if err != nil {
		log.Println("Couldn't find guilds channels", err)
	}

	for _, c := range channels {
		fmt.Println(c.Name, ": ", len(c.Messages), " ", c.LastMessageID, " ", c.Messages)
	}
}

func handleMessage(session *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == session.State.User.ID {
		return
	}
	fmt.Println("Got a message, ", msg.Content)
	findLatestMessage(msg.GuildID)
	message := strings.ToLower(msg.Content)

	if strings.HasPrefix(message, "go ") {
		command := strings.Split(message[3:], " ")
		fmt.Println("Got command, ", command)
		switch command[0] {
		case "joke":
			fetchAndRespond := func(jokeType string) {
				joke, err := fetchJoke(jokeType)
				if err != nil {
					replyToChannel(msg.ChannelID, "Joke failed")
					log.Println(err)
					return
				}
				tellJoke(session, msg, joke)
			}
			if len(command) == 2 {
				if command[1] == "boomer" {
					tellJoke(session, msg, boomerJoke)
				} else {
					fetchAndRespond(command[1])
				}
				break
			}
			fetchAndRespond("Any")
		case "help":
			session.ChannelMessageSend(msg.ChannelID, "I'll look for therapy places for you in my free time")
		case "comic":
			sendComic(session, msg)
		case "meme":
			var subreddit string
			if len(command) > 1 {
				subreddit = command[1]
			}
			sendMeme(subreddit, session, msg)
		}
	}

	// if time.Since(findLatestMessage()).Hours()/24 > 5 {
	// 	userDM, err := session.UserChannelCreate(msg.Author.ID)
	// 	if err != nil {
	// 		log.Println("Couldn't get into user's DMs", err)
	// 	}
	// 	session.ChannelMessageSend(userDM.ID, "I think they all moved to Raino. If you want, you can give me a message for them. Just type `go hermes *your-message*` and I will give it to them!")

	// }
}

func CacheMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("got ya")
	decoder := json.NewDecoder(r.Body)
	newMsg := struct {
		Msg       string    `json:"msg"`
		User      string    `json:"userName"`
		Timestamp time.Time `json:"timestamp,omitempty"`
	}{}

	err := decoder.Decode(&newMsg)
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()
	newMsg.Timestamp = time.Now()
	col := dbClient.Database("raingo-cache").Collection("messageQueue")
	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "Timestamp", Value: -1}})
	findOptions.SetLimit(1)
	cur, err := col.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Fatalln(err)
	}
	for cur.Next(ctx) {
		oldMsg := struct {
			_id       string
			Msg       string    `json:"msg"`
			User      string    `json:"userName"`
			Timestamp time.Time `json:"timestamp,omitempty"`
		}{}
		if err = cur.Decode(&oldMsg); err != nil {
			log.Fatal(err)
		}
		fmt.Println(oldMsg)
		id, err := primitive.ObjectIDFromHex(oldMsg._id)
		if err != nil {
			log.Fatal(err)
		}
		deleteResult, _ := col.DeleteOne(context.TODO(), bson.M{"_id": id})
		if deleteResult.DeletedCount == 0 {
			log.Fatal("Error on deleting one Hero", err)
		}
		res, err := col.InsertOne(ctx, newMsg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("res: ", res)
	}
}

func createHttpServer() {
	router := httprouter.New()
	router.POST("/api/raino", listenToRaindrops)
	router.GET("/api/wake", WakeMeUpBeforeYouGoGo)
	router.POST("/api/wake", CacheMessage)
	port, present := os.LookupEnv("PORT")
	if !present {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func connectToMongo() *mongo.Client {
	mongoUser, present := os.LookupEnv("MONGODB_USER")
	if !present {
		panic("No bot token found!")
	}
	mongoPass, present := os.LookupEnv("MONGODB_PASS")
	if !present {
		panic("No bot token found!")
	}
	mongoURL := fmt.Sprintf("mongodb+srv://%s:%s@raingo-cache.fuqhm.mongodb.net/raingo-cache?retryWrites=true&w=majority", mongoUser, mongoPass)
	clientOptions := options.Client().
		ApplyURI(mongoURL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	godotenv.Load(".env")
	botToken, present := os.LookupEnv("BOT_TOKEN")
	if !present {
		panic("No bot token found!")
	}

	dbClient = connectToMongo()

	bot, err := discordgo.New("Bot " + botToken)
	fmt.Println("API version:", discordgo.APIVersion)
	if err != nil {
		fmt.Println("Error creating bot session!")
		panic(err)
	}
	bot.AddHandler(handleMessage)
	bot.Open()
	session = bot
	createHttpServer()
	bot.Close()
}
