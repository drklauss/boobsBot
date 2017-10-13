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
	TmCelebCmd     = "celeb"      // Команда получения item по тегу real celebnsfw
	TmHelloCmd     = "hello"      // Команда приветствия
	TmHelpCmd      = "help"       // Команда help
	TmRateCmd      = "rate"       // Команда получения ссылки для голосования
	// Админские
	TmAdmin         = "admin"                        // Команда получения админской клавиатуры
	TmDebugStartCmd = "\xF0\x9F\x94\xB4 Debug Start" // Команда включения режима отладки
	TmDebugStopCmd  = "\xF0\x9F\x94\xB5 Debug Stop"  // Команда выключения режима отладки
	TmTopViewersCmd = "\xF0\x9F\x93\x8A Top Viewers" // Команда получения отчета топ-зрителей
	TmUpdateCmd     = "\xF0\x9F\x94\x84 Update"      // Команда обновления
)
