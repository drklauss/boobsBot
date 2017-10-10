package config

const (
	LogFile = "./bot.log"
	// Telegram Bot
	TmApiUrl           = "https://api.telegram.org/bot"
	TmToken            = "425337521:AAGcOjS44c86oAStJdn5xqWOfGcPIBeMiw4"
	TmFullBotName      = "@DornBot"
	TmUpdateTime       = 2        // Отправка запроса для апдейта
	TmSkipMessagesTime = 60       // Вычитывает апдейты за указанный промежуток времени
	TmDevUserId        = 90310429 // Мой userId
	// Комманды бота
	TmNSFWCmd      = "nsfw"       // Команда получения nsfw
	TmRealGirlsCmd = "real_girls" // Команда получения item по тегу real girls
	TmCelebCmd     = "celeb_nsfw" // Команда получения item по тегу real celebnsfw
	TmHelloCmd     = "hello"      // Команда приветствия
	// Админские
	TmAdmin         = "admin"                        // Команда получения админской клавиатуры
	TmDebugStartCmd = "debugStart"                   // Команда включения режима отладки
	TmDebugStopCmd  = "debugStop"                    // Команда выключения режима отладки
	TmTopViewersCmd = "\xF0\x9F\x93\x8A Top Viewers" // Команда получения отчета топ-зрителей
)
