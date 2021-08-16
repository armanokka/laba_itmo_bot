package main

import (
    "database/sql"
    "github.com/armanokka/translobot/dashbot"
    "gorm.io/gorm"
    "time"
)

const (
    DashBotAPIKey = "cjVjdWDRijXDk5kl9yGi5TTS9XImME7HbZMOg09F"
    AdminID       = 579515224
    botToken string = "1737819626:AAEoc8WyCq_8rFQcY4q0vtkhqCKro8AudfI"
)

var (
    db  *gorm.DB
    bot *BotAPI
    analytics *dashbot.DashBot
    langs = [...]string{"en","ru","la","ja","ar","fr","de","af","uk","uz","es","ko","zh","hi","bn","pt","mr","te","ms","tr","vi","ta","ur","jv","it","fa","gu","ab","aa","ak","sq","am","an","hy","as","av","ae","ay","az","bm","ba","eu","be","bh","bi","bs","br","bg","my","ca","ch","ce","ny","cv","kw","co","cr","hr","cs","da","dv","nl","eo","et","ee","fo","fj","fi","ff","gl","ka","el","gn","ht","ha","he","hz","ho","hu","ia","id","ie","ga","ig","ik","io","is","iu","kl","kn","kr","ks","kk","km","ki","rw","ky","kv","kg","ku","kj","la","lb","lg","li","ln","lo","lt","lu","lv","gv","mk","mg","ml","mt","mi","mh","mn","na","nv","nb","nd","ne","ng","nn","no","ii","nr","oc","oj","cu","om","or","os","pa","pi","pl","ps","qu","rm","rn","ro","sa","sc","sd","se","sm","sg","sr","gd","sn","si","sk","sl","so","st","su","sw","ss","sv","tg","th","ti","bo","tk","tl","tn","to","ts","tt","tw","ty","ug","ve","vo","wa","cy","wo","fy","xh","yi","yo","za"}
)

// Users is table in DB
type Users struct {
    ID     int64 `gorm:"primaryKey;index;not null"`
    MyLang string `gorm:"default:en"`
    ToLang string `gorm:"default:fr"`
    Act sql.NullString `gorm:"default:null"`
    Mailing bool `gorm:"default:true"`
    Usings int64 `gorm:"default:0"`
    Lang string `gorm:"default:en"`
}

// Sponsorships is table in DB
type Sponsorships struct {
    ID int64 // ID of user that bought offer
    Text string // Advertise
    ToLangs string // String of languages separated by "," of users that must to receive advertise
    Start time.Time // When offer starts
    Finish time.Time // When offer finish
}

// SponsorshipsOffers is table in DB
type SponsorshipsOffers struct {
    ID int64 // ID of user that bought offer
    Text string // Advertise
    ToLangs string // String of languages separated by "," of users that must to receive advertise
    Days int // Days of sponsorship
}

type Localization struct {
    LanguageCode string
    Text string
}

// type Referrers struct {
//     ID int64 `gorm:"primaryKey;index;not null"`
//     Code string
//     Users int64
// }
