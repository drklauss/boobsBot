package config

const (
	LogFile = "./bot.log"
	// Telegram Bot
	TmApiUrl           = "https://api.telegram.org/bot"
	TmToken            = "425337521:AAGcOjS44c86oAStJdn5xqWOfGcPIBeMiw4"
	TmFullBotName      = "@DornBot"
	TmUpdateTime       = 1        // Отправка запроса для апдейта
	TmSkipMessagesTime = 60       // Вычитывает апдейты за указанный промежуток времени
	TmAdminUserId      = 90310429 // Мой chatId
	// Комманды бота
	TmHelloCmd      = "hello"        // Say TmHelloCmd
	TmHotCmd        = RdtHotCategory // Give me a corn
	TmTopViewersCmd = "topViewers"   // Top Viewers report
)
