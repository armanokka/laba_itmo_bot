module github.com/armanokka/translobot

// +heroku goVersion go1.15
go 1.15

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.0.0-rc1.0.20210608165722-fb1de2fb48dd => D:\Go-Workspace\FORKS\go-telegram-bot-api\telegram-bot-api\v5

require (
	github.com/PuerkitoBio/goquery v1.7.0
	github.com/armanokka/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/emvi/iso-639-1 v1.0.1
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.0.0-rc1.0.20210608165722-fb1de2fb48dd
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/m90/go-chatbase v2.0.0+incompatible
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/valyala/fasthttp v1.27.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.11
)
