package config

const (
	LogFile = "./bot.log"
	// Telegram Bot
	TmApiUrl      = "https://api.telegram.org/bot"
	TmToken       = "425337521:AAGcOjS44c86oAStJdn5xqWOfGcPIBeMiw4"
	TmFullBotName = "@DornBot"
	TmUpdateTime  = 2
	// Комманды бота
	HelloCom = "/hello" // Say Hello
	JokeCom  = "/joke"  // Tell me a joke
	DornCom  = "/dorn"  // Give me a corn
	// Reddit Bot
	RedditApiUrl = "https://oauth.reddit.com/api/v1/"
	RedditToken  = "8uEfARcCknLTAZ_G8dWvpE7ey9c"
	NSFWNew      = "https://www.reddit.com/r/NSFW_GIF/new/"
)
