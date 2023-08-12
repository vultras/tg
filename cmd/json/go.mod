module jsoned

go 1.20

require (
	github.com/mojosa-software/goscript v0.0.0-20230626091305-86a004b7769c
	github.com/mojosa-software/got v0.0.0-20230812125405-bbe076f29abe
)

require github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 // indirect

replace github.com/mojosa-software/got => ./../..

