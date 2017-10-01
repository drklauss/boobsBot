package config

const (
	LogFile = "./bot.log"
	// Telegram Bot
	TmApiUrl           = "https://api.telegram.org/bot"
	TmToken            = "425337521:AAGcOjS44c86oAStJdn5xqWOfGcPIBeMiw4"
	TmFullBotName      = "@DornBot"
	TmUpdateTime       = 1        // Отправка запроса для апдейта
	TmSkipMessagesTime = 60       // Вычитывает апдейты за указанный промежуток времени
	TmDevUserId        = 90310429 // Мой userId
	// Комманды бота
	TmHelloCmd      = "hello"        // Команда приветствия
	TmHotCmd        = RdtHotCategory // Команда получения порновидосика
	TmTopViewersCmd = "topViewers"   // Команда получения отчета топ-зрителей
	TmDebugStartCmd = "debugStart"   // Команда включения режима отладки
	TmDebugEndCmd   = "debugEnd"     // Команда выключения режима отладки
)
