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
	TmHelloCmd      = "hello"                        // Команда приветствия
	TmNSFWVideo     = "nsfw_mp4"                     // Команда получения nsfw видео
	TmRealGirls     = "real_girls"                   // Команда получения item по тегу real girls
	TmDebugStartCmd = "debugStart"                   // Команда включения режима отладки
	TmDebugEndCmd   = "debugStop"                    // Команда выключения режима отладки
	TmReports       = "reports"                      // Команда получения отчетов
	TmTopViewersCmd = "\xF0\x9F\x93\x8A Top Viewers" // Команда получения отчета топ-зрителей
)
