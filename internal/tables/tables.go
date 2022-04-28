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
		"start_info": map[string]string{
			"id": "Sekarang saya akan terjemahkan dari %s ke %s dan kembali lagi. Jika anda ingin mengubah bahasa, menulis /memulai\n\n\nMengirim saya pesan, dan saya akan menerjemahkannya",
			"ar": "الآن سوف تترجم من %s %s والعودة مرة أخرى. إذا كنت ترغب في تغيير لغة الكتابة /بدء\n\n\nأرسل لي رسالة وأنا سوف تترجم لهم",
			"ru": "Теперь я буду переводить с %s на %s и обратно. Если захочешь сменить языки, напишешь /start\n\n\nОтправляй мне сообщения, а я буду их переводить",
			"uz": "Endi %s dan %s va orqaga tarjima qilaman. Agar tillarni o'zgartirishni xohlasangiz, yozing / boshlang\n\n\nMenga xabar yuboring va men ularni tarjima qilaman",
			"it": "Ora ho intenzione di tradurre da %s a %s e indietro. Se vuoi cambiare le lingue, scrivi /start\n\n\nMi mandare messaggi, e io sarò loro di tradurre",
			"en": "Now I will translate from %s to %s and back again. If you want to change languages, write /start\n\n\nSend me messages, and I'll translate them",
			"pt": "Agora vou traduzir de %s para %s e volta. Se quiseres mudar de línguas, irmão /start\n\n\nОтправляй-me mensagens, e eu serei o seu traduzir",
			"uk": "Тепер я буду перекладати з %s на %s і назад. Якщо захочеш змінити мови, напишеш /start\n\n\nВідправляй мені повідомлення, а я буду їх переказувати",
			"es": "Ahora voy a traducir de %s a %s, y viceversa. Si quieres cambiar el idioma, escribas /start\n\n\nEnvía me un mensaje, y me voy a traducirlas",
			"de": "Jetzt werde ich das übersetzen von %s auf %s und wieder zurück. Wenn Sie möchten ändern die Sprachen, schreiben /start\n\n\nSenden Sie mir eine Nachricht, ich werden Sie zu übersetzen",
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
		"Как переводить еще удобнее": map[string]string{
			"es": "Cómo traducir aún mejor",
			"uk": "Як перекладати ще зручніше",
			"id": "Bagaimana menerjemahkan bahkan lebih nyaman",
			"it": "Come tradurre ancora più conveniente",
			"pt": "Como traduzir ainda mais",
			"ar": "كيفية ترجمة أكثر سهولة",
			"de": "Wie übersetzt man noch bequemer",
			"ru": "Как переводить еще удобнее",
			"uz": "Tarjima qilish yanada qulayroq",
			"en": "How to translate even more conveniently",
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
		"Вас устраивает этот перевод?": map[string]string{
			"uz": "Ushbu tarjimadan qoniqasizmi?",
			"ru": "Вас устраивает этот перевод?",
			"it": "Siete soddisfatti di questa traduzione?",
			"en": "Are you satisfied with this translation?",
			"uk": "Вас влаштовує цей переклад?",
			"ar": "هل أنت راض عن هذه الترجمة ؟ ",
			"de": "Sie ordnet diese übersetzung?",
			"pt": "Você está satisfeito com esta tradução?",
			"es": "Usted satisfecho con esta traducción?",
			"id": "Apakah anda puas dengan terjemahan ini?",
		},
		"Мы учли ваш голос\nВаши отзывы помогут нам улучшить сервис": map[string]string{
			"ru": "Мы учли ваш голос\nВаши отзывы помогут нам улучшить сервис",
			"it": "Abbiamo preso in considerazione la tua voce\nI vostri commenti ci aiuteranno a migliorare il servizio",
			"de": "Wir berücksichtigen Ihre Stimme\nIhre Rückmeldungen helfen uns, den Service zu verbessern",
			"ar": "أخذنا التصويت بعين الاعتبار\nملاحظاتك سوف تساعدنا في تحسين الخدمة",
			"uk": "Ми врахували ваш голос\nВаші відгуки допоможуть нам поліпшити сервіс",
			"uz": "Sizning ovozingizni hisobga oldik\nSizning fikringiz xizmatni yaxshilashga yordam beradi",
			"en": "We took your vote into account\nYour feedback will help us improve the service",
			"es": "Hemos tenido en cuenta su voz\nSus comentarios nos ayudarán a mejorar el servicio",
			"id": "Kami mengambil suara anda ke rekening\nTanggapan anda akan membantu kami meningkatkan layanan",
			"pt": "Nós levamos em consideração a sua voz\nOs seus comentários ajudam-nos a melhorar o serviço",
		},
		"Спасибо!": map[string]string{
			"id": "Terima kasih!!!",
			"en": "Thanks!",
			"es": "Gracias!",
			"pt": "Obrigado!",
			"it": "Grazie!",
			"ar": "وذلك بفضل!",
			"de": "Danke!",
			"uz": "Rahmat!",
			"ru": "Спасибо!",
			"uk": "Спасибі!",
		},
		"Время ожидания жалобы на перевод истекло. Переведите сообщение еще раз": map[string]string{
			"id": "Waktu tunggu untuk transfer keluhan telah berakhir. Menerjemahkan pesan lagi",
			"de": "Die Wartezeit auf Beschwerde übersetzungen abgelaufen ist. Übersetzen Sie die Nachricht erneut",
			"it": "Il tempo di attesa reclami sulla traduzione è scaduto. Portare il messaggio ancora una volta",
			"ar": "وقت الانتظار لنقل الشكوى قد انتهت. ترجمة الرسالة مرة أخرى",
			"uz": "Transfer haqidagi shikoyatni kutish muddati tugadi. Xabarni qayta tarjima qiling",
			"pt": "O tempo de espera reclamações de tradução expirou. Coloque novamente a mensagem",
			"ru": "Время ожидания жалобы на перевод истекло. Переведите сообщение еще раз",
			"en": "The waiting time for a transfer complaint has expired. Translate the message again",
			"uk": "Час очікування скарги на переклад минув. Перекажіть повідомлення ще раз",
			"es": "El tiempo de espera de la queja en la traducción agotado. Ponga de nuevo el mensaje",
		},
		"format_translation_tip": map[string]string{
			"es": "💡Sugerencia:\nSólo Translo puede traducir estos textos con formato, como: \"𝑴𝒓 𝒑𝒂𝒖𝒍𝒔𝒐𝒏 𝑷𝒊𝒆𝒕𝒆𝒓\"\nNi google, ni yahoo! ni ningún otro traductor porque no pueden. Pruebe y compruebe en este",
			"id": "💡Petunjuk:\nHanya Translo dapat menerjemahkan diformat teks-teks seperti: \"𝑴𝒓 𝒑𝒂𝒖𝒍𝒔𝒐𝒏 𝑷𝒊𝒆𝒕𝒆𝒓\"\nGoogle atau Yandex, atau translator tidak bisa melakukan itu. Mencobanya dan melihat untuk diri sendiri",
			"en": "💡Hint:\nOnly Translo can translate formatted texts such as: \"𝑴𝒓 𝒑𝒂𝒖𝒍𝒔𝒐𝒏 𝑷𝒊𝒆𝒕𝒆𝒓\"\nAny Google or Yandex, or any other translator can't do that. Try it out and see for yourself",
			"de": "💡Tipp:\nNur Translo übersetzen kann diese formatierte Texte, wie: \"𝑴𝒓 𝒑𝒂𝒖𝒍𝒔𝒐𝒏 𝑷𝒊𝒆𝒕𝒆𝒓\"\nWeder Google, Yandex, noch irgendein anderer übersetzer nicht so können. Probieren Sie es aus und überzeugen Sie sich",
			"uk": "💡Підказка:\nТільки Translo може переводити такі форматовані тексти, як: \"𝑴𝒓 𝒑𝒂𝒖𝒍𝒔𝒐𝒏 𝑷𝒊𝒆𝒕𝒆𝒓\"\nНі гугл, ні яндекс, ні будь-який інший перекладач так не можуть. Спробуйте і переконаєтеся в цьому",
			"pt": "💡Dica:\nSó Translo pode traduzir tais textos formatados, como: \"𝑴𝒓 𝒑𝒂𝒖𝒍𝒔𝒐𝒏 𝑷𝒊𝒆𝒕𝒆𝒓\"\nNem o google, nem ela, nem qualquer outro tradutor não podem. Tente e certifique-se neste",
			"it": "💡Suggerimento:\nSolo Translo in grado di tradurre i testi formattati, come: \"𝑴𝒓 𝒑𝒂𝒖𝒍𝒔𝒐𝒏 𝑷𝒊𝒆𝒕𝒆𝒓\"\nNé google, né yandex, né qualsiasi altro traduttore non possono. Provare e vedere di",
			"uz": "Подсказ uchi:\nFaqat Translo quyidagi formatlangan matnlarni tarjima qilishi mumkin: \"𝑴𝒓 𝒑𝒂𝒖𝒍𝒔𝒐𝒏 𝑷𝒊𝒆𝒕𝒆𝒓\"\nNa Google, na Yandex, na boshqa tarjimon ham shunday qila olmaydi. Uni sinab ko'ring va ishonch hosil qiling",
			"ru": "💡Подсказка:\nТолько Translo может переводить такие форматированные тексты, как: \"𝑴𝒓 𝒑𝒂𝒖𝒍𝒔𝒐𝒏 𝑷𝒊𝒆𝒕𝒆𝒓\"\nНи гугл, ни яндекс, ни любой другой переводчик так не могут. Попробуйте и убедитесь в этом",
			"ar": "💡تلميح:\nفقط Translo يمكن أن تترجم تنسيق النصوص مثل: \"𝑴𝒓 𝒑𝒂𝒖𝒍𝒔𝒐𝒏 𝑷𝒊𝒆𝒕𝒆𝒓\"\nأي جوجل أو ياندكس ، أو أي مترجم لا تستطيع أن تفعل ذلك. محاولة وانظر لنفسك",
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
