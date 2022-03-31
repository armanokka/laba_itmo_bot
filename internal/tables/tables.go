package tables

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Users is table in DB
type Users struct {
	ID           int64 `gorm:"primaryKey;index;not null"`
	MyLang       string
	ToLang       string
	Act          string
	Usings       int  `gorm:"default:0"`
	Blocked      bool `gorm:"default:false"`
	LastActivity time.Time
	Lang         string `gorm:"-"` // internal
}

func (u *Users) SetLang(lang string) {
	u.Lang = lang
}

func (u Users) Localize(key string, placeholders ...interface{}) string {
	key = strings.TrimSpace(key)
	localization := map[string]map[string]string{
		"Попробовать": map[string]string{
			"de": "Versuchen",
			"es": "tratar",
			"uk": "Спробувати",
			"uz": "harakat qilib ko'ring",
			"it": "Tentativo",
			"en": "try",
			"id": "mencoba",
			"pt": "experimentar",
			"ar": "يحاول",
			"ru": "Попробовать",
		},
		"Произошла ошибка": map[string]string{
			"en": "Sorry, an error has occurred\n\nPlease do not block the bot, we will fix everything soon.",
			"ru": "Извините, произошла ошибка\n\nПожалуйста, не блокируйте бота, скоро мы все исправим.",
			"es": "Lo sentimos, ha ocurrido un error\n\nPor favor, no bloquees el bot, arreglaremos todo pronto.",
			"pt": "Desculpe, ocorreu um erro\n\nPor favor, não bloqueie o bot, vamos corrigir tudo em breve.",
			"de": "Entschuldigung, ein Fehler ist aufgetreten\n\nBitte blockieren Sie den Bot nicht, wir werden alles bald beheben.",
			"uk": "Вибачте, сталася помилка\n\nБудь ласка, не блокуйте робота, скоро ми всі виправимо.",
			"uz": "Kechirasiz, xatolik yuz berdi\n\nIltimos, botni bloklamang, tez orada hammasini tuzatamiz.",
			"id": "Maaf, telah terjadi kesalahan\n\nTolong jangan blokir bot, kami akan segera memperbaiki semuanya.",
			"it": "Spiacenti, si è verificato un errore\n\nSi prega di non bloccare il bot, sistemeremo tutto al più presto.",
			"ar": "عذرا، حدث خطأ\n\nمن فضلك لا تحجب الروبوت ، سنصلح كل شيء قريبا.",
		},
		"Добро пожаловать. Мы рады, что вы снова с нами. ✋": map[string]string{
			"es": "Bienvenidos. Nos alegramos de que esté con nosotros de nuevo. ✋",
			"it": "Benvenuto. Siamo felici che tu sia di nuovo con noi. ✋",
			"ar": "أهلا بك. نحن سعداء لأنك معنا مرة أخرى. ✋",
			"pt": "Receber. Estamos felizes por você estar conosco novamente. ✋",
			"en": "Welcome. We are glad that you are with us again. ✋",
			"ru": "Добро пожаловать. Мы рады, что вы снова с нами. ✋",
			"de": "Willkommen zurück. Wir freuen uns, dass Sie wieder bei uns sind. ✋",
			"uk": "Ласкаво просимо. Ми раді, що ви знову з нами. ✋",
			"uz": "Xush kelibsiz. Yana biz bilan ekanligingizdan xursandmiz. ✋",
			"id": "Selamat datang. Kami senang Anda bersama kami lagi. ✋.",
		},
		"Порекомендовать бота": map[string]string{
			"de": "Empfehlen Sie einen Bot",
			"uk": "Рекомендувати робота",
			"it": "Consiglia un bot",
			"ar": "يوصي روبوت",
			"en": "Recommend a bot",
			"ru": "Порекомендовать бота",
			"es": "recomendar un bot",
			"uz": "Botni tavsiya eting",
			"id": "Rekomendasikan bot",
			"pt": "Recomende um bot",
		},
		"inline_ad": map[string]string{
			"de": "🔥 <a href=\"https://t.me/translobot\">Translo</a> 🌐 - beste in der Welt der bot-übersetzer in Telegram.",
			"uk": "🔥 <a href=\"https://t.me/translobot\">Translo</a> 🌐 - кращий в світі бот-перекладач в Telegram.",
			"id": "🔥 <a href=\"https://t.me/translobot\">Translo</a> 🌐 - terbaik di dunia bot penerjemah untuk Telegram.",
			"it": "🔥 <a href=\"https://t.me/translobot\">Translo</a> 🌐 - è la migliore al mondo bot-traduttore a Telegram.",
			"pt": "🔥 <a href=\"https://t.me/translobot\">Translo</a> 🌐 - o melhor do mundo bot-tradutor em Telegram.",
			"ar": "🔥 <a href=\"https://t.me/translobot\">Translo</a> 🌐 - العالم أفضل بوت مترجم برقية.",
			"en": "🔥 <a href=\"https://t.me/translobot\">Translo</a> 🌐 - the world's best bot translator for Telegram.",
			"ru": "🔥 <a href=\"https://t.me/translobot\">Translo</a>  🌐 - лучший в мире бот-переводчик в Telegram.",
			"es": "🔥 <a href=\"https://t.me/translobot\">Translo</a> 🌐 - mejor en el mundo de bots traductor en Telegram.",
			"uz": "🔥 <a href=\"https://t.me/translobot\">Translo</a> 🌐 - dunyodagi eng yaxshi telegram tarjimon boti",
		},
		"кликните, чтобы порекомендовать бота в чате": map[string]string{
			"es": "haga clic aquí para recomendar un bot de chat en vivo",
			"it": "cliccare per raccomandare un bot in chat",
			"en": "click to recommend the bot in the chat",
			"de": "klicken Sie, um zu empfehlen im Chat bot",
			"uz": "chatda botni tavsiya qilish uchun bosing",
			"id": "klik untuk merekomendasikan bot di chat",
			"pt": "clique para recomendar um bot no bate-papo",
			"ar": "انقر فوق أن يوصي بوت في الدردشة",
			"ru": "кликните, чтобы порекомендовать бота в чате",
			"uk": "натисніть, щоб порекомендувати бота в чаті",
		},
		"Озвучить": map[string]string{
			"uz": "Ovoz berish",
			"id": "Suara itu lebih dari",
			"it": "Esprimere",
			"pt": "Ouvir",
			"ar": "صوت أكثر",
			"en": "Voice it over",
			"ru": "Озвучить",
			"de": "Äußern",
			"es": "La voz de",
			"uk": "Озвучити",
		},
		"Примеры": map[string]string{
			"es": "Ejemplos",
			"uk": "Приклади",
			"it": "Esempi",
			"ar": "أمثلة",
			"ru": "Примеры",
			"de": "Beispiele",
			"id": "Contoh",
			"pt": "Exemplos",
			"en": "Examples",
			"uz": "Misollar",
		},
		"Переводы": map[string]string{
			"uk": "Переклади",
			"pt": "Transferências",
			"en": "Translations",
			"de": "Übersetzungen",
			"es": "Traducciones",
			"it": "Traduzioni",
			"ar": "نقل",
			"ru": "Переводы",
			"uz": "Tarjima",
			"id": "Transfer",
		},
		"Словарь": map[string]string{
			"uz": "Lug'at",
			"id": "Kamus",
			"pt": "Dicionário",
			"en": "Dictionary",
			"ru": "Словарь",
			"uk": "Словник",
			"ar": "القاموس",
			"de": "Wörterbuch",
			"es": "Diccionario",
			"it": "Dizionario",
		},
		"На какой язык перевести?": map[string]string{
			"en": "What language should I translate it into?",
			"ru": "На какой язык перевести?",
			"de": "Auf welche Sprache zu übersetzen?",
			"es": "En qué idioma a traducir?",
			"uz": "Qaysi tilni tarjima qilish kerak?",
			"uk": "На яку мову перевести?",
			"id": "Bahasa apa aku harus menerjemahkannya ke?",
			"it": "In quale lingua tradurre?",
			"pt": "Qualquer idioma para traduzir?",
			"ar": "ما هي اللغة التي يجب أن تترجم ذلك ؟ ",
		},
		"Я рекомендую @translobot": map[string]string{
			"en": "I recommend @translobot",
			"de": "Ich empfehle @translobot",
			"uk": "Я рекомендую @translobot",
			"id": "Saya merekomendasikan @translobot",
			"ar": "أوصي @translobot",
			"ru": "Я рекомендую @translobot",
			"es": "Recomiendo @translobot",
			"uz": "Kanallarning birortasi ham ushbu nashr bilan ulashmadi",
			"it": "Mi raccomando @translobot",
			"pt": "Eu recomendo @translobot",
		},
		"Рассказать про нас": map[string]string{
			"es": "Hablar sobre nosotros",
			"pt": "Falar sobre nós",
			"ar": "أخبرنا عنا",
			"en": "Tell us about us",
			"ru": "Рассказать про нас",
			"de": "Über uns erzählen",
			"uk": "Розповісти про нас",
			"uz": "Biz haqimizda aytib bering",
			"id": "Beritahu kami tentang kami",
			"it": "Parlare di noi",
		},
		"Понравился бот? 😎 Поделись с друзьями, нажав на кнопку": map[string]string{
			"ru": "Понравился бот? 😎 Поделись с друзьями, нажав на кнопку",
			"it": "Piaciuto il bot? 😎 Condividi con i tuoi amici, cliccando sul pulsante",
			"uz": "Bot kabi? 😎 Tugmasini bosib do'stlaringiz bilan baham ko'ring",
			"id": "Apakah anda suka bot? 😎 Berbagi dengan teman-teman anda dengan mengklik tombol",
			"pt": "Gostei do bot? 😎 Partilhe com os seus amigos, clicando no botão",
			"ar": "هل تحب الآلي ؟ 😎 مشاركتها مع أصدقائك عن طريق النقر على زر",
			"en": "Did you like the bot? 😎 Share with your friends by clicking on the button",
			"de": "Mochte ich den bot? 😎 Teilen Sie mit Ihren Freunden, indem Sie auf die Schaltfläche",
			"es": "Bien el bot? 😎 Comparte con tus amigos haciendo clic en el botón",
			"uk": "Сподобався бот? 😎 Поділіться з друзями, натиснувши на кнопку",
		},
		"Отправь текстовое сообщение, чтобы я его перевел": map[string]string{
			"uz": "Uni tarjima qilish uchun matnli xabar yuboring",
			"ar": "ترسل لي رسالة حتى استطيع ترجمته",
			"es": "Envía un mensaje de texto, para que lo tradujo",
			"uk": "Відправ текстове повідомлення, щоб я його переклав",
			"id": "Saya mengirim pesan teks sehingga saya dapat menerjemahkannya",
			"it": "Invia un messaggio di testo che ho tradotto",
			"pt": "Envie uma mensagem de texto para que eu o traduzi",
			"en": "Send me a text message so I can translate it",
			"ru": "Отправь текстовое сообщение, чтобы я его перевел",
			"de": "Sende eine Textnachricht, dass ich ihn übersetzt habe",
		},
		"Просто напиши мне текст, а я его переведу": map[string]string{
			"uk": "Просто напиши мені текст, а я його переведу",
			"uz": "Menga faqat matn yozing va uni tarjima qilaman",
			"pt": "Basta escrever-me um texto e eu traduzir",
			"ru": "Просто напиши мне текст, а я его переведу",
			"de": "Schreiben Sie mir einfach den Text, und ich werde",
			"es": "Simplemente me escribe el texto, y lo he transferir",
			"ar": "فقط اكتب لي رسالة وأنا سوف ترجمته",
			"en": "Just write me a text and I'll translate it",
			"id": "Hanya menulis saya sms dan saya akan menerjemahkannya",
			"it": "Mi scrivi il testo, e io lo traduco",
		},
		"%s не поддерживается": map[string]string{
			"de": "%s wird nicht unterstützt",
			"uk": "%s не підтримується",
			"id": "%s tidak didukung",
			"it": "%s non è supportato",
			"pt": "%s não é suportado",
			"ar": "%s غير معتمد",
			"en": "%s is not supported",
			"ru": "%s не поддерживается",
			"es": "%s no es compatible",
			"uz": "%s qo'llab-quvvatlanmaydi",
		},
		"Теперь я буду переводить с %s на %s и обратно. Если захочешь изменить, напишешь /start": map[string]string{
			"uz": "Endi %s dan %s va orqaga tarjima qilaman. Agar o'zgartirish bo'lsangiz, yozish / boshlash",
			"id": "Sekarang saya akan terjemahkan dari %s untuk %s dan kembali lagi. Jika anda ingin mengubah hal itu, menulis /memulai",
			"it": "Ora ho intenzione di tradurre %s su %s e indietro. Se vuoi cambiare, scrivi /start",
			"pt": "Agora eu vou traduzir %s por %s e volta. Se o quiseres alterar, escreva /start",
			"en": "Now I will translate from %s to %s and back again. If you want to change it, write /start",
			"de": "Jetzt werde ich das übersetzen von %s auf %s und wieder zurück. Wenn Sie wollen ändern, schreiben /start",
			"es": "Ahora voy a traducir de %s de %s y viceversa. Si quieres cambiar, escribas /start",
			"uk": "Тепер я буду перекладати з %s на %s і назад. Якщо захочеш змінити, напишеш /start",
			"ar": "الآن سوف تترجم من %s إلى %s والعودة مرة أخرى. إذا كنت ترغب في تغيير كتابة /بدء",
			"ru": "Теперь я буду переводить с %s на %s и обратно. Если захочешь изменить, напишешь /start",
		},
		"не получилось перевести": map[string]string{
			"de": "hat nicht funktioniert übersetzen",
			"uk": "не вийшло перевести",
			"it": "non è riuscito a tradurre",
			"ar": "لم أستطع ترجمته",
			"en": "couldn't translate it",
			"es": "no pude traducir",
			"uz": "tarjima qila olmadi",
			"id": "tidak bisa menerjemahkannya",
			"pt": "não conseguimos traduzir",
			"ru": "не получилось перевести",
		},
		"Report a bug or suggest a feature": map[string]string{
			"en": "Report a bug or suggest a feature",
			"es": "Informar de un error o sugerir una característica",
			"uz": "Xato haqida xabar bering yoki xususiyatni taklif qiling",
			"ru": "Сообщите об ошибке или предложите функцию",
			"de": "Report a bug oder ein feature vorschlagen",
			"uk": "Повідомте про помилку або запропонуйте функцію",
			"id": "Melaporkan bug atau menyarankan fitur",
			"it": "Segnalare un bug o suggerire una funzione",
			"pt": "Comunicar um erro ou \" sugerir um recurso",
			"ar": "تقرير الشوائب أو اقتراح ميزة",
		},
		"Завершить диалог": map[string]string{
			"ru": "Завершить диалог",
			"id": "Mengakhiri dialog",
			"pt": "Terminar a conversa",
			"en": "End the dialog",
			"de": "Dialog beenden",
			"es": "Completar un diálogo",
			"uk": "Завершити діалог",
			"uz": "Suhbatni yakunlang",
			"it": "Terminare il dialogo",
			"ar": "في نهاية الحوار",
		},
		"Диалог завершен": map[string]string{
			"it": "Conversazione",
			"pt": "A conversa terminou",
			"ar": "الحوار انتهى",
			"en": "The dialog is finished",
			"ru": "Диалог завершен",
			"de": "Dialog abgeschlossen",
			"es": "La conversación ha finalizado",
			"uk": "Діалог завершено",
			"uz": "Muloqot yakunlandi",
			"id": "Dialog selesai",
		},
		"Chat mode": map[string]string{
			"id": "Modus chatting",
			"it": "Modalità Chat",
			"en": "Chat mode",
			"uk": "Режим чату",
			"uz": "Chat tartibi",
			"pt": "Modo de bate-papo",
			"ar": "وضع دردشة",
			"ru": "Режим чата",
			"de": "Chat-Modus",
			"es": "El modo de Chat",
		},
		"пример текста": map[string]string{
			"es": "texto de ejemplo",
			"uk": "приклад тексту",
			"uz": "matn namunasi",
			"ru": "пример текста",
			"de": "Beispieltext",
			"id": "contoh teks",
			"it": "esempio di testo",
			"pt": "texto de exemplo",
			"ar": "النص عينة",
			"en": "sample text",
		},
		"Как переводить сообщения, не выходя из чата": map[string]string{
			"ru": "Как переводить сообщения, не выходя из чата",
			"de": "Wie übersetzt man Nachrichten, ohne aus dem Chat",
			"es": "Cómo traducir un mensaje sin salir de chat",
			"uk": "Як переводити повідомлення, не виходячи з чату",
			"it": "Come tradurre i messaggi, non lasciare la chat",
			"en": "How to translate messages without leaving the chat",
			"uz": "Suhbatdan chiqmasdan xabarlarni qanday tarjima qilish kerak",
			"id": "Bagaimana menerjemahkan pesan tanpa meninggalkan obrolan",
			"pt": "Como traduzir as mensagens sem sair do chat",
			"ar": "كيفية ترجمة الرسائل دون ترك الدردشة",
		},
	}

	if v, ok := localization[key]; ok {
		if v, ok := v[u.Lang]; ok {
			return fmt.Sprintf(v, placeholders...)
		}
	}
	return fmt.Sprintf(key, placeholders...)
}

type UsersLogs struct {
	ID      int64          // fk users.id
	Intent  sql.NullString // varchar(25)
	Text    string         // varchar(518)
	FromBot bool
	Date    time.Time
}

type Mailing struct {
	ID int64 // fk users.id
}

func (Mailing) TableName() string {
	return "mailing"
}
