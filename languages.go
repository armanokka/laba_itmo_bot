package main

import "fmt"


func Localize(text, lang string, placeholders ...interface{}) string {
    
    
    var languages = map[string][]Localization {
        "To voice": {
            {
            LanguageCode: "es",
                Text:         "a voz",
            },
            {
            LanguageCode: "pt",
                Text:         "dar voz",
            },
            {
            LanguageCode: "id",
                Text:         "untuk menyuarakan",
            },
            {
            LanguageCode: "uz",
                Text:         "ovoz berish",
            },
            {
            LanguageCode: "de",
                Text:         "aussprechen",
            },
            {
            LanguageCode: "uk",
                Text:         "озвучити",
            },
            {
            LanguageCode: "it",
                Text:         "esprimere",
            },
            {
            LanguageCode: "ru",
                Text:         "Озвучить",
            },
        },
        "My Language": {
            {
            LanguageCode: "id",
                Text:         "Bahasa Saya",
            },
            {
            LanguageCode: "it",
                Text:         "La mia lingua",
            },
            {
            LanguageCode: "uz",
                Text:         "Tilimni",
            },
            {
            LanguageCode: "de",
                Text:         "Meine Sprache",
            },
            {
            LanguageCode: "ru",
                Text:         "Мой Язык",
            },
            {
            LanguageCode: "es",
                Text:         "Mi Idioma",
            },
            {
            LanguageCode: "uk",
                Text:         "Моя Мова",
            },
            {
            LanguageCode: "pt",
                Text:         "A Minha Língua",
            },
        },
        "/help": {
            {
            LanguageCode: "uk",
                Text:         " Що може зробити цей бот? \n▫️ Translo дозволяє перекладати повідомлення 182+ мовами.\n Як перекласти повідомлення? \n▫️ По-перше, вам потрібно налаштувати свою мову, потім налаштувати переклад мови, наступне надіслати текстові повідомлення і бот швидко їх перекладе.\n Як налаштувати мову? \n▫️ Клацніть на кнопку під назвою \"Моя мова\"\n Як налаштувати мову на мову, яку я хочу перекласти? \nКлацніть на кнопку нижче під назвою \"Перекласти мову\"\n Чи є якісь інші цікаві речі? \nТак, бот підтримує вбудований режим. Почніть вводити псевдонім @translobot у поле введення повідомлення, а потім напишіть туди текст, який потрібно перекласти.\n У мене є пропозиція або я знайшов помилку! \n👉 Зв’яжіться зі мною pls - @armanokka",
            },
            {
            LanguageCode: "id",
                Text:         "Apa yang bisa dilakukan robot ini?\n▫️ ️ooo memungkinkan Anda untuk menerjemahkan pesan ke 182+ bahasa.\nBagaimana menerjemahkan pesan?\n▫️  Anda harus menata bahasa Anda, kemudian menerjemahkan bahasa, teks kirim pesan teks dan keduanya akan menerjemahkannya dengan cepat.\n▫️ Bagaimana cara mengatur bahasaku?\n▫️ ️ Klik di tombol di bawah ini disebut \" Bahasa saya\"\nBagaimana mengatur bahasa ke dalam saya ingin menerjemahkan?\nKlik di tombol di bawah ini yang disebut \"menerjemahkan Bahasa\"\n▫️ ️ Apakah ada hal-hal yang menarik lainnya?\nYa, dukungan robot dalam mode inline. Mulai mengetik Nama panggilan @translobot dalam kolom masukan pesan dan kemudian menulis di sana teks yang ingin Anda terjemahkan.\n▫️ ️ Aku punya saran atau aku menemukan bug!\nHubungi aku pls - @armanokka",
            },
            {
            LanguageCode: "en",
                Text:         "What can this bot do?\n▫️ Translo allow you to translate messages into 182+ languages.\nHow to translate message?\n▫️ Firstly, you have to setup your language, then setup translate language, next send text messages and bot will translate them quickly.\nHow to setup my language?\n▫️ Click on the button below called \"My Language\"\nHow to setup language into I want to translate?\n▫️ Click on the button below called \"Translate Language\"\nAre there any other interesting things?\n▫️ Yes, the bot support inline mode. Start typing the nickname @translobot in the message input field and then write there the text you want to translate.\nI have a suggestion or I found bug!\n▫️ 👉 Contact me pls - @armanokka",
            },
            {
            LanguageCode: "de",
                Text:         "Was kann dieser Bot tun?\n▫️ Mit Translo können Sie Nachrichten in mehr als 182 Sprachen übersetzen.\nWie übersetzt man eine Nachricht?\n▫️ Zuerst müssen Sie Ihre Sprache einrichten, dann die Übersetzungssprache einrichten, als nächstes Textnachrichten senden und der Bot wird sie schnell übersetzen.\nWie richte ich meine Sprache ein?\n▫️ Klicken Sie unten auf die Schaltfläche \"Meine Sprache\"\nWie richte ich die Sprache ein, in die ich übersetzen möchte?\n▫️ Klicken Sie unten auf die Schaltfläche \"Sprache übersetzen\"\nGibt es noch andere interessante Dinge?\n▫️ Ja, der Bot unterstützt den Inline-Modus. Geben Sie den Spitznamen @translobot in das Nachrichteneingabefeld ein und schreiben Sie dort den Text, den Sie übersetzen möchten.\nIch habe einen Vorschlag oder ich habe einen Fehler gefunden!\n▫️ 👉 Kontaktieren Sie mich bitte - @armanokka",
            },
            {
            LanguageCode: "ru",
                Text:         " Что умеет этот бот? \n▫️ Translo позволяет переводить сообщения на 182+ языков.\n Как перевести сообщение? \n▫️ Во-первых, вам нужно настроить свой язык, затем настроить язык перевода, затем отправить текстовые сообщения, и бот быстро их переведет.\n Как настроить мой язык? \n▫️ Нажмите на кнопку под названием «Мой язык».\n Как установить язык, на котором я хочу переводить? \n▫️ Нажмите кнопку ниже под названием «Перевести язык».\n Есть еще что-нибудь интересное? \n▫️ Да, бот поддерживает встроенный режим. Начните вводить псевдоним @translobot в поле ввода сообщения, а затем впишите туда текст, который хотите перевести.\n У меня есть предложение или я обнаружил ошибку! \n▫️ 👉 Свяжитесь со мной, пожалуйста - @armanokka",
            },
            {
            LanguageCode: "es",
                Text:         " ¿Qué puede hacer este bot? \n▫️ Translo te permite traducir mensajes a más de 182 idiomas.\n ¿Cómo traducir un mensaje? \n▫️ En primer lugar, debe configurar su idioma, luego configurar el idioma de traducción, luego enviar mensajes de texto y el bot los traducirá rápidamente.\n ¿Cómo configurar mi idioma? \n▫️ Haga clic en el botón de abajo llamado \"Mi idioma\"\n ¿Cómo configurar el idioma al que quiero traducir? \n▫️ Haga clic en el botón de abajo llamado \"Traducir idioma\"\n ¿Hay otras cosas interesantes? \n▫️ Sí, el bot admite el modo en línea. Comience a escribir el apodo @translobot en el campo de entrada del mensaje y luego escriba allí el texto que desea traducir.\n ¡Tengo una sugerencia o encontré un error! \n▫️ 👉 Contáctame por favor - @armanokka",
            },
            {
            LanguageCode: "pt",
                Text:         " O que este bot pode fazer? \n▫️ Translo permite que você traduza mensagens em mais de 182 idiomas.\n Como traduzir a mensagem? \n▫️ Em primeiro lugar, você tem que configurar seu idioma, depois configurar traduzir o idioma, em seguida enviar mensagens de texto e o bot irá traduzi-las rapidamente.\n Como configurar meu idioma? \n▫️ Clique no botão abaixo chamado \"Meu Idioma\"\n Como configurar o idioma para o que desejo traduzir? \n▫️ Clique no botão abaixo chamado \"Traduzir Idioma\"\n Existem outras coisas interessantes? \n▫️ Sim, o bot suporta o modo inline. Comece digitando o apelido @translobot no campo de entrada da mensagem e então escreva lá o texto que deseja traduzir.\n Tenho uma sugestão ou encontrei um bug! \n▫️ 👉 Contate-me, pls - @armanokka",
            },
            {
            LanguageCode: "it",
                Text:         "Cosa può fare questo bot?\n▫️ Translo ti consente di tradurre i messaggi in oltre 182 lingue.\nCome tradurre il messaggio?\n▫️ In primo luogo, devi impostare la tua lingua, quindi impostare la lingua di traduzione, quindi inviare messaggi di testo e il bot li tradurrà rapidamente.\nCome impostare la mia lingua?\n▫️ Clicca sul pulsante in basso chiamato \"La mia lingua\"\nCome impostare la lingua in cui voglio tradurre?\n▫️ Clicca sul pulsante in basso chiamato \"Traduci lingua\"\nCi sono altre cose interessanti?\n▫️ Sì, il bot supporta la modalità in linea. Inizia a digitare il nickname @translobot nel campo di inserimento del messaggio e poi scrivi lì il testo che vuoi tradurre.\nHo un suggerimento o ho trovato un bug!\n▫️ 👉 Contattami per favore - @armanokka",
            },
            {
            LanguageCode: "uz",
                Text:         " Bu bot nima qila oladi? \n▫️ Translo sizga xabarlarni 182+ tilga tarjima qilishga imkon beradi.\n Xabarni qanday tarjima qilish kerak? \n▫️ Birinchidan, siz o'z tilingizni o'rnatishingiz kerak, so'ngra tarjima tilini sozlashingiz kerak, keyin matnli xabarlarni yuboring va bot ularni tezda tarjima qiladi.\n Mening tilimni qanday o'rnatish kerak? \n▫️ Quyidagi \"Mening tilim\" deb nomlangan tugmani bosing.\n Men tarjima qilishni xohlagan tilni qanday o'rnataman? \n▫️ \"Tilni tarjima qilish\" deb nomlangan tugmani bosing.\n Boshqa qiziqarli narsalar bormi? \n▫️ Ha, botni inline rejimida qo'llab-quvvatlash. Xabarlarni kiritish maydoniga @translobot taxallusini yozishni boshlang va keyin tarjima qilmoqchi bo'lgan matni yozing.\n Menda taklif bor yoki men xato topdim! \n▫️ 👉 Men bilan bog'laning pls - @armanokka",
            },
        },
        "Now your language is %s\n\nPress \"⬅Back\" to exit to menu": {
            {
            LanguageCode: "ru",
                Text:         "Теперь ваш язык %s\n\nНажмите «⬅Back» для выхода в меню.",
            },
            {
            LanguageCode: "es",
                Text:         "Ahora tu idioma es %s\n\nPresione \"⬅Back\" para salir al menú",
            },
            {
            LanguageCode: "uk",
                Text:         "Зараз ваша мова% %s\n\nНатисніть \"⬅Back\", щоб вийти в меню",
            },
            {
            LanguageCode: "pt",
                Text:         "Agora seu idioma é %s\n\nPressione \"⬅Back\" para sair para o menu",
            },
            {
            LanguageCode: "id",
                Text:         "Sekarang bahasa Anda adalah %s\n\nTekan \"⬅Back\" untuk keluar ke menu",
            },
            {
            LanguageCode: "it",
                Text:         "Ora la tua lingua è %s\n\nPremere \"⬅Back\" per uscire dal menu",
            },
            {
            LanguageCode: "uz",
                Text:         "Endi sizning tilingiz %s\n\nMenyuga chiqish uchun \"⬅Back\" tugmasini bosing",
            },
            {
            LanguageCode: "de",
                Text:         "Ihre Sprache ist jetzt %s\n\nDrücken Sie \"⬅Back\", um das Menü zu verlassen",
            },
        },
        "Translate Language": {
            {
            LanguageCode: "id",
                Text:         "Bahasa untuk menerjemahkan",
            },
            {
            LanguageCode: "it",
                Text:         "Lingua per tradurre",
            },
            {
            LanguageCode: "pt",
                Text:         "Língua para tradução",
            },
            {
            LanguageCode: "ru",
                Text:         "Язык перевода",
            },
            {
            LanguageCode: "uk",
                Text:         "Мова перекладу",
            },
            {
            LanguageCode: "uz",
                Text:         "Tarjima qilish uchun til",
            },
            {
            LanguageCode: "de",
                Text:         "Sprache zum Übersetzen",
            },
            {
            LanguageCode: "es",
                Text:         "Idioma para traducir",
            },
        },
        "Please, send text message": {
            {
            LanguageCode: "it",
                Text:         "Per favore, invia un messaggio di testo",
            },
            {
            LanguageCode: "uz",
                Text:         "Iltimos, SMS yuboring",
            },
            {
            LanguageCode: "de",
                Text:         "Bitte SMS senden",
            },
            {
            LanguageCode: "ru",
                Text:         "Пожалуйста, отправьте текстовое сообщение",
            },
            {
            LanguageCode: "es",
                Text:         "Por favor, envíe un mensaje de texto",
            },
            {
            LanguageCode: "uk",
                Text:         "Будь ласка, надішліть текстове повідомлення",
            },
            {
            LanguageCode: "pt",
                Text:         "Por favor, envie mensagem de texto",
            },
            {
            LanguageCode: "id",
                Text:         "Tolong, kirim pesan teks",
            },
        },
        "Too big text": {
            {
            LanguageCode: "de",
                Text:         "Zu großer Text",
            },
            {
            LanguageCode: "ru",
                Text:         "Слишком большой текст",
            },
            {
            LanguageCode: "es",
                Text:         "Texto demasiado grande",
            },
            {
            LanguageCode: "uk",
                Text:         "Занадто великий текст",
            },
            {
            LanguageCode: "pt",
                Text:         "Texto muito grande",
            },
            {
            LanguageCode: "id",
                Text:         "Teks terlalu besar",
            },
            {
            LanguageCode: "it",
                Text:         "Testo troppo grande",
            },
            {
            LanguageCode: "uz",
                Text:         "Juda katta matn",
            },
        },
        "/my_lang": {
            {
            LanguageCode: "it",
                Text:         "Per impostare la tua lingua, esegui una delle seguenti operazioni: .\n\nℹ️ Invia poche parole nella tua lingua, ad esempio: \"Hi, how are you today?\" - la lingua sarà l'inglese, o \"L'amour ne fait pas d'erreurs\" - la lingua sarà francese, e così via.\nℹ️ Oppure invia il nome della tua lingua, ad es. \"Russo\", o \"Giapponese\", o \"Arabo\", ecc.",
            },
            {
            LanguageCode: "ru",
                Text:         "Чтобы настроить  свой язык , выполните  одно  из следующих действий: 👇\n\nℹ️ Отправьте  несколько слов  на своем языке, например: \"Hi, how are you today?\" - язык будет английский, или «« L'amour ne fait pas d'erreurs »- язык будет французским и так далее.\nℹ️ Или отправьте название своего языка, например «Русский», «Японский», «Арабский» и т. Д.",
            },
            {
            LanguageCode: "es",
                Text:         "Para configurar  su idioma , haga  una  de las siguientes opciones: 👇\n\nℹ️ Envía  algunas palabras  en tu idioma, por ejemplo: \"Hi, how are you today?\" - el idioma será el inglés, o \"L'amour ne fait pas d'erreurs\" - idioma será francés, etc.\nℹ️ O envíe el nombre de su idioma, p. ej. \"Ruso ''\", o \"Japonés\", o \"Árabe\", e.t.c.",
            },
            {
            LanguageCode: "uk",
                Text:         "Щоб налаштувати  свою мову , виконайте  одне  з наступного: 👇\n\nℹ️ Надішліть  кілька слів  своєю мовою, наприклад: \"Hi, how are you today?\" - мова буде англійською, або \"L'amour ne fait pas d'erreurs\" - мова буде французькою тощо.\nℹ️ Або надішліть назву вашої мови, напр. \"Російська\", або\"японська\", або  арабська, тощо",
            },
            {
            LanguageCode: "id",
                Text:         "Untuk mengatur Bahasa anda, melakukan satu hal:\n\nℹ️ contoh: \"hai, bagaimana kabarmu hari ini?\"- bahasa akan bahasa Inggris, atau \"l'amour ne fait pas d'erreurs\" - bahasa akan Perancis, dan seterusnya.\nℹ️ atau kirim Nama bahasa Anda dalam bahasa Inggris, misal \"rusia\", atau \"jepang\", atau \"arab\", e.t.c.",
            },
            {
            LanguageCode: "pt",
                Text:         "Para configurar  seu idioma , faça  uma  das seguintes opções: 👇\n\nℹ️ Envie  algumas palavras  em seu idioma, por exemplo: \"Hi, how are you today?\" - o idioma será o inglês, ou \"L'amour ne fait pas d'erreurs\" - o idioma será francês e assim por diante.\nℹ️ Ou envie o nome do seu idioma, por ex. \"Russian\", ou \"Japanese\", ou \"Arabic\", e.t.c.",
            },
            {
            LanguageCode: "en",
                Text:         "To setup your language, do one of the following: 👇\n\nℹ️ Send few words in your language, for example: \"Hi, how are you today?\" - language will be English, or \"L'amour ne fait pas d'erreurs\" - language will be French, and so on.\nℹ️ Or send the name of your language, e.g. \"Russian\", or \"Japanese\", or  \"Arabic\", e.t.c.",
            },
            {
            LanguageCode: "de",
                Text:         "Führen Sie eine der folgenden Schritte aus, um Ihre Sprache einzurichten: 👇\n\nℹ️ Sende einige Wörter in deiner Sprache, zum Beispiel: \"Hi, how are you today?\" - Sprache wird Englisch sein, oder \"L'amour ne fait pas d'erreurs\" - Sprache wird französisch sein und so weiter.\nℹ️ Oder schicke den Namen deiner Sprache, z.B. \"Russisch\", oder \"Japanisch\", oder \"Arabisch\", usw.",
            },
            {
            LanguageCode: "uz",
                Text:         " Tilingizni  sozlash uchun quyidagilardan  birini bajaring: 👇\n\nℹ️ O'zingizning tilingizda  bir nechta so'zlarni  yuboring, masalan: \"Hi, how are you today? \"\" \"- til ingliz tilida bo'ladi yoki\" \"L'amour ne fait pas d'erreurs`\" - til frantsuzcha bo'ladi va hokazo.\nℹ️ Yoki o'z tilingiz nomini yuboring, masalan. \"\" \"Ruscha\" yoki \"\" Yaponcha \"yoki\" \"Arabcha\", e.t.c.",
            },
        },
        "⏳ Translating...": {
            {
            LanguageCode: "de",
                Text:         "⏳ Übersetzen...",
            },
            {
            LanguageCode: "ru",
                Text:         "⏳ Перевожу...",
            },
            {
            LanguageCode: "es",
                Text:         "⏳ Traducir...",
            },
            {
            LanguageCode: "uk",
                Text:         "⏳ Перекласти...",
            },
            {
            LanguageCode: "pt",
                Text:         "⏳ Traduzir...",
            },
            {
            LanguageCode: "id",
                Text:         "⏳ Terjemahkan...",
            },
            {
            LanguageCode: "it",
                Text:         "⏳ Traduci...",
            },
            {
            LanguageCode: "uz",
                Text:         "⏳ Tarjima qilish...",
            },
        },
        "Variants": {
            {
            LanguageCode: "uk",
                Text:         "Варіанти",
            },
            {
            LanguageCode: "uz",
                Text:         "Variantlar",
            },
            {
            LanguageCode: "de",
                Text:         "Varianten",
            },
            {
            LanguageCode: "es",
                Text:         "Variantes",
            },
            {
            LanguageCode: "id",
                Text:         "Varian",
            },
            {
            LanguageCode: "it",
                Text:         "varianti",
            },
            {
            LanguageCode: "pt",
                Text:         "Variantes",
            },
            {
            LanguageCode: "ru",
                Text:         "Варианты",
            },
        },
        "Delete": {
            {
            LanguageCode: "es",
                Text:         "Borrar",
            },
            {
            LanguageCode: "id",
                Text:         "Menghapus",
            },
            {
            LanguageCode: "it",
                Text:         "Elimina",
            },
            {
            LanguageCode: "ru",
                Text:         "Удалить",
            },
            {
            LanguageCode: "en",
                Text:         "Delete",
            },
            {
            LanguageCode: "de",
                Text:         "Löschen",
            },
            {
            LanguageCode: "uk",
                Text:         "Видалити",
            },
            {
            LanguageCode: "uz",
                Text:         "O'chirish",
            },
            {
            LanguageCode: "pt",
                Text:         "Excluir",
            },
        },
        "Error, sorry.": {
            {
            LanguageCode: "pt",
                Text:         "Erro, desculpe.",
            },
            {
            LanguageCode: "id",
                Text:         "Kesalahan, maaf.",
            },
            {
            LanguageCode: "it",
                Text:         "Errore, scusa.",
            },
            {
            LanguageCode: "uz",
                Text:         "Xato, uzr.",
            },
            {
            LanguageCode: "de",
                Text:         "Fehler, tut mir leid.",
            },
            {
            LanguageCode: "ru",
                Text:         "Ошибка, извините.",
            },
            {
            LanguageCode: "es",
                Text:         "Error, lo siento.",
            },
            {
            LanguageCode: "uk",
                Text:         "Помилка, вибачте.",
            },
        },
        "Rate": {
            {
            LanguageCode: "uk",
                Text:         "Оцініть бота",
            },
            {
            LanguageCode: "id",
                Text:         "Tolong, beri peringkat bot",
            },
            {
            LanguageCode: "it",
                Text:         "Per favore, vota il bot",
            },
            {
            LanguageCode: "pt",
                Text:         "Por favor, avalie o bot",
            },
            {
            LanguageCode: "de",
                Text:         "Bitte bewerte den Bot",
            },
            {
            LanguageCode: "ru",
                Text:         "Пожалуйста, оцените бота",
            },
            {
            LanguageCode: "es",
                Text:         "Por favor, califica al bot",
            },
            {
            LanguageCode: "uz",
                Text:         "Iltimos, botga baho bering",
            },
            {
            LanguageCode: "en",
                Text:         "Please, rate the bot",
            },
        },
        "🙎‍♂️Profile": {
            {
            LanguageCode: "es",
                Text:         "🙍‍♂️Perfil",
            },
            {
            LanguageCode: "uz",
                Text:         "🙍‍♂️Profil",
            },
            {
            LanguageCode: "it",
                Text:         "🙍‍♂️Profilo",
            },
            {
            LanguageCode: "uk",
                Text:         "🙍‍♂️Профіль",
            },
            {
            LanguageCode: "id",
                Text:         "🙍‍♂️Profil",
            },
            {
            LanguageCode: "pt",
                Text:         "🙍‍♂️Perfil",
            },
            {
            LanguageCode: "ru",
                Text:         "🙍‍♂️Профиль",
            },
            {
            LanguageCode: "de",
                Text:         "🙍‍♂️Profil",
            },
        },
        "💬 Bot language": {
            {
            LanguageCode: "en",
                Text:         "💬 Bot language",
            },
            {
            LanguageCode: "de",
                Text:         "💬 Bot-Sprache",
            },
            {
            LanguageCode: "es",
                Text:         "💬 Lenguaje bot",
            },
            {
            LanguageCode: "pt",
                Text:         "💬 Linguagem de bot",
            },
            {
            LanguageCode: "uz",
                Text:         "💬 Bot tili",
            },
            {
            LanguageCode: "id",
                Text:         "💬 Bahasa bot",
            },
            {
            LanguageCode: "uk",
                Text:         "💬 Бот-мова",
            },
            {
            LanguageCode: "it",
                Text:         "💬 Linguaggio Bot",
            },
            {
            LanguageCode: "ru",
                Text:         "💬 Язык бота",
            },
        },
        "/start":{
            {
            LanguageCode: "uz",
                Text:         "Sizning tilingiz - %s, tarjima qilish uchun - %s.",
            },
            {
            LanguageCode: "pt",
                Text:         "Seu idioma é - %s, e o idioma para tradução é - %s.",
            },
            {
            LanguageCode: "de",
                Text:         "Ihre Sprache ist - %s und die Sprache für die Übersetzung ist - %s.",
            },
            {
            LanguageCode: "id",
                Text:         "Bahasa Anda adalah - %s, dan bahasa terjemahannya adalah - %s.",
            },
            {
            LanguageCode: "en",
                Text:         "Your language is - %s, and the language for translation is - %s.",
            },
            {
            LanguageCode: "ru",
                Text:         "Ваш язык - %s, а язык перевода - %s.",
            },
            {
            LanguageCode: "uk",
                Text:         "Ваша мова - %s, а мова перекладу - %s.",
            },
            {
            LanguageCode: "it",
                Text:         "La tua lingua è - %se la lingua per la traduzione è - %s.",
            },
            {
            LanguageCode: "es",
                Text:         "Su idioma es - %s, y el idioma de traducción es - %s.",
            },
        },
        "💡 Instruction": {
            {
            LanguageCode: "it",
                Text:         "💡 Istruzione",
            },
            {
            LanguageCode: "uz",
                Text:         "💡 Yo'riqnoma",
            },
            {
            LanguageCode: "de",
                Text:         "💡 Anweisung",
            },
            {
            LanguageCode: "ru",
                Text:         "💡 Инструкция",
            },
            {
            LanguageCode: "es",
                Text:         "💡 Instrucción",
            },
            {
            LanguageCode: "uk",
                Text:         "💡 Інструкція",
            },
            {
            LanguageCode: "pt",
                Text:         "💡 Instrucao",
            },
            {
            LanguageCode: "id",
                Text:         "💡 Instruksi",
            },
        },
        "Now translate language is %s\n\nPress \"⬅Back\" to exit to menu": {
            {
            LanguageCode: "uk",
                Text:         "Тепер мова перекладу - %s\n\nНатисніть \"⬅Back\", щоб вийти в меню",
            },
            {
            LanguageCode: "pt",
                Text:         "Agora o idioma de tradução é %s\n\nPressione \"⬅Back\" para sair para o menu",
            },
            {
            LanguageCode: "id",
                Text:         "Sekarang bahasa terjemahan adalah %s\n\nTekan \"⬅Back\" untuk keluar ke menu",
            },
            {
            LanguageCode: "it",
                Text:         "Ora la lingua di traduzione è %s\n\nPremere \"⬅Back\" per uscire dal menu",
            },
            {
            LanguageCode: "uz",
                Text:         "Endi tarjima qilish tili %s\n\nMenyuga chiqish uchun \"⬅Back\" tugmasini bosing",
            },
            {
            LanguageCode: "de",
                Text:         "Die Übersetzungssprache ist jetzt %s\n\nDrücken Sie \"⬅Back\", um das Menü zu verlassen",
            },
            {
            LanguageCode: "ru",
                Text:         "Теперь язык перевода %s\n\nНажмите «⬅Back» для выхода в меню.",
            },
            {
            LanguageCode: "es",
                Text:         "Ahora el idioma de traducción es %s\n\nPresione \"⬅Back\" para salir al menú",
            },
        },
        "Now press /start 👈": {
            {
            LanguageCode: "it",
                Text:         "Ora premi /start 👈",
            },
            {
            LanguageCode: "id",
                Text:         "Sekarang tekan /start 👈",
            },
            {
            LanguageCode: "de",
                Text:         "Drücken Sie nun /start 👈",
            },
            {
            LanguageCode: "pt",
                Text:         "Agora pressione /start 👈",
            },
            {
            LanguageCode: "uz",
                Text:         "Endi /start tugmachasini bosing 👈",
            },
            {
            LanguageCode: "en",
                Text:         "Now press /start 👈",
            },
            {
            LanguageCode: "ru",
                Text:         "Теперь нажмите /start 👈",
            },
            {
            LanguageCode: "es",
                Text:         "Ahora presione /start 👈",
            },
            {
            LanguageCode: "uk",
                Text:         "Тепер натисніть /start 👈",
            },
        },
        "/to_lang": {
            {
            LanguageCode: "uz",
                Text:         "Tarjima tilini sozlash uchun quyidagilardan birini bajaring: 👇\n\nℹ️ O'zingizning tilingizga tarjima qilishni xohlagan bir nechta so'zlarni yuboring, masalan: \"Hi, how are you?\" - til inglizcha bo'ladi yoki \"L'amour ne fait pas d'erreurs\" - frantsuzcha bo'ladi va hokazo.\nℹ️ Yoki til nomini yuboring, masalan. \"Ruscha\", yoki \"Yaponcha\" yoki \"Arabcha\", e.t.c.",
            },
            {
            LanguageCode: "de",
                Text:         "Um die Übersetzungssprache einzurichten, führen Sie einen der folgenden Schritte aus:\n\nℹ️ Schicke ein paar Wörter in die Sprache, die du übersetzen möchtest, zum Beispiel: \"Hi, how are you?\" - Sprache ist Englisch oder \"L'amour ne fait pas d'erreurs\" - Sprache ist Französisch und so weiter.\nℹ️ Oder sende den Namen der Sprache, z.B. \"Russisch\" oder \"Japanisch\" oder \"Arabisch\", usw.",
            },
            {
            LanguageCode: "ru",
                Text:         "Чтобы настроить язык перевода, выполните одно из следующих действий: 👇\n\nℹ️ Отправьте несколько слов на языке, который вы хотите перевести, например: «Hi, how are you?» - язык будет английский, или \"L'amour ne fait pas d'erreurs\" - язык будет французским, и так далее.\nℹ️ Или отправьте название языка, например «русский», или «японский», или «арабский» и т. Д.",
            },
            {
            LanguageCode: "es",
                Text:         "Para configurar el idioma de traducción, realice una de las siguientes acciones: 👇\n\nℹ️ Envíe algunas palabras en el idioma que desea traducir, por ejemplo: \"Hi, how are you?\" - el idioma será el inglés, o \"L'amour ne fait pas d'erreurs\" - el idioma será el francés, etc.\nℹ️ O envíe el nombre del idioma, p. ej. \"Ruso\", \"Japonés\" o \"Árabe\", etc.",
            },
            {
            LanguageCode: "pt",
                Text:         "Para configurar o idioma de tradução, siga um destes procedimentos: 👇\n\nℹ️ Envie algumas palavras no idioma que deseja traduzir, por exemplo: \"Hi, how are you?\" - o idioma será inglês ou \"L'amour ne fait pas d'erreurs\" - o idioma será o francês e assim por diante.\nℹ️ Ou envie o nome do idioma, por ex. \"Russo\" ou \"Japonês\" ou \"Árabe\", e.t.c.",
            },
            {
            LanguageCode: "en",
                Text:         "To setup translate language, do one of the following: 👇\n\nℹ Send few words in language into you want to translate, for example: \"Hi, how are you?\" - language will be English, or \"L'amour ne fait pas d'erreurs\" - language will be French, and so on.\nℹ️ Or send the name of language, e.g. \"Russian\", or \"Japanese\", or  \"Arabic\", e.t.c.",
            },
            {
            LanguageCode: "it",
                Text:         "Per impostare la lingua di traduzione, esegui una delle seguenti operazioni: 👇\n\nℹ️ Invia poche parole nella lingua che vuoi tradurre, ad esempio: \"Hi, how are you?\" - la lingua sarà l'inglese, o \"L'amour ne fait pas d'erreurs\" - la lingua sarà il francese, e così via.\nℹ️ Oppure inviare il nome della lingua, ad es. \"Russo\" o \"Giapponese\" o \"Arabo\", ecc.",
            },
            {
            LanguageCode: "uk",
                Text:         "Щоб налаштувати мову перекладу, виконайте одну з таких дій: 👇\n\nℹ️ Надішліть кілька слів мовою на мову, яку потрібно перекласти, наприклад: \"Hi, how are you?\" - мовою буде англійська, або \"L'amour ne fait pas d'erreurs\" - мовою буде французька тощо.\nℹ️ Або надішліть назву мови, напр. \"російська\", або \"японська\", або \"арабська\", напр.",
            },
            {
            LanguageCode: "id",
                Text:         "Untuk menerjemahkan bahasa, lakukan satu bahasa berikut:\n\nℹ mengirim beberapa kata dalam bahasa ke dalam Anda ingin menerjemahkan, misalnya: \"hai, apa kabar?\"- bahasa akan bahasa Inggris, atau \"l'amour ne fait pas d'erreurs\" - bahasa akan Perancis, dan seterusnya.\nℹ️atau kirim Nama bahasa, misalnya \"Rusia\", atau \"Jepang\", atau \"Arab\", misalnya. t.c.",
            },
        },
        "Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)": {
            {
            LanguageCode: "uz",
                Text:         "Uzr, xato sabab.\n\nIltimos, botni bloklamang, yaqin kelajakda xatoni tuzataman, administrator bu xato haqida allaqachon ogohlantirilgan ;)",
            },
            {
            LanguageCode: "de",
                Text:         "Sorry, Fehler verursacht.\n\nBitte blockieren Sie den Bot nicht, ich werde den Fehler in naher Zukunft beheben, der Administrator wurde bereits vor diesem Fehler gewarnt ;)",
            },
            {
            LanguageCode: "ru",
                Text:         "Извините, произошла ошибка.\n\nПожалуйста, не блокируйте бота, я исправлю ошибку в ближайшее время, администратор уже предупрежден об этой ошибке ;)",
            },
            {
            LanguageCode: "es",
                Text:         "Lo sentimos, error causado.\n\nPor favor, no bloquees el bot, corregiré el error en un futuro cercano, el administrador ya ha sido advertido sobre este error ;)",
            },
            {
            LanguageCode: "uk",
                Text:         "Вибачте, сталася помилка.\n\nБудь ласка, не блокуйте бота, я виправлю помилку найближчим часом, адміністратор вже попереджений про цю помилку ;)",
            },
            {
            LanguageCode: "pt",
                Text:         "Desculpa, erro causado.\n\nPor favor, não bloqueie o bot, eu vou corrigir o bug em um futuro próximo, o administrador já foi avisado sobre este ",
            },
            {
            LanguageCode: "id",
                Text:         "Maaf, kesalahan disebabkan.\n\nTolong, jangan halangi robot itu, aku akan memperbaiki bug dalam waktu dekat, administrator telah diperingatkan tentang kesalahan ini ;)",
            },
            {
            LanguageCode: "it",
                Text:         "Scusa, errore causato.\n\nPer favore, non bloccare il bot, correggerò il bug nel prossimo futuro, l'amministratore è già stato avvertito di questo errore;)",
            },
        },
        "Unsupported language or internal error": {
            {
            LanguageCode: "pt",
                Text:         "Idioma não suportado ou erro interno",
            },
            {
            LanguageCode: "id",
                Text:         "Bahasa yang tidak didukung atau kesalahan internal",
            },
            {
            LanguageCode: "it",
                Text:         "Lingua non supportata o errore interno",
            },
            {
            LanguageCode: "uz",
                Text:         "Qo'llab-quvvatlanmaydigan til yoki ichki xato",
            },
            {
            LanguageCode: "de",
                Text:         "Nicht unterstützte Sprache oder interner Fehler",
            },
            {
            LanguageCode: "ru",
                Text:         "Неподдерживаемый язык или внутренняя ошибка",
            },
            {
            LanguageCode: "es",
                Text:         "Idioma no admitido o error interno",
            },
            {
            LanguageCode: "uk",
                Text:         "Непідтримувана мова або внутрішня помилка",
            },
        },
        "⬅Back": {
            {
            LanguageCode: "it",
                Text:         "⬅Indietro",
            },
            {
            LanguageCode: "pt",
                Text:         "⬅Back",
            },
            {
            LanguageCode: "ru",
                Text:         "⬅Назад",
            },
            {
            LanguageCode: "uk",
                Text:         "⬅Назад",
            },
            {
            LanguageCode: "uz",
                Text:         "⬅Arka",
            },
            {
            LanguageCode: "de",
                Text:         "⬅Zurück",
            },
            {
            LanguageCode: "es",
                Text:         "⬅Atrás",
            },
            {
            LanguageCode: "id",
                Text:         "⬅Kembali",
            },
        },
        "Please, select bot language": {
            {
            LanguageCode: "de",
                Text:         "Bitte Bot-Sprache auswählen",
            },
            {
            LanguageCode: "es",
                Text:         "Por favor, seleccione el idioma del bot",
            },
            {
            LanguageCode: "en",
                Text:         "Please, select bot language",
            },
            {
            LanguageCode: "pt",
                Text:         "Por favor, selecione o idioma do bot",
            },
            {
            LanguageCode: "id",
                Text:         "Silakan, pilih bahasa bot",
            },
            {
            LanguageCode: "uz",
                Text:         "Iltimos, bot tilini tanlang",
            },
            {
            LanguageCode: "uk",
                Text:         "Виберіть мову бота",
            },
            {
            LanguageCode: "it",
                Text:         "Per favore, seleziona la lingua del bot",
            },
            {
            LanguageCode: "ru",
                Text:         "Пожалуйста, выберите язык бота",
            },
        },
        "Ok, close": {
            {
            LanguageCode: "pt",
                Text:         "Ok fechar",
            },
            {
            LanguageCode: "id",
                Text:         "Oke, tutup",
            },
            {
            LanguageCode: "uk",
                Text:         "Гаразд, близько",
            },
            {
            LanguageCode: "it",
                Text:         "Ok, chiudi",
            },
            {
            LanguageCode: "uz",
                Text:         "Yaxshi, yaqin",
            },
            {
            LanguageCode: "de",
                Text:         "Okay, nah",
            },
            {
            LanguageCode: "ru",
                Text:         "Хорошо, закрыть",
            },
            {
            LanguageCode: "es",
                Text:         "Ok cerrar",
            },
        },
        "Empty result": {
            {
            LanguageCode: "es",
                Text:         "Resultado vacío",
            },
            {
            LanguageCode: "uk",
                Text:         "Порожній результат",
            },
            {
            LanguageCode: "pt",
                Text:         "Resultado vazio",
            },
            {
            LanguageCode: "id",
                Text:         "Hasil kosong",
            },
            {
            LanguageCode: "it",
                Text:         "Risultato vuoto",
            },
            {
            LanguageCode: "uz",
                Text:         "Bo'sh natija",
            },
            {
            LanguageCode: "de",
                Text:         "Leeres Ergebnis",
            },
            {
            LanguageCode: "ru",
                Text:         "Пустой результат",
            },
        },
        "Failed to detect the language. Please enter something else": {
            {
            LanguageCode: "ru",
                Text:         "Не удалось определить язык. Пожалуйста, введите что-нибудь еще",
            },
            {
            LanguageCode: "es",
                Text:         "No se pudo detectar el idioma. Por favor ingrese algo más",
            },
            {
            LanguageCode: "uk",
                Text:         "Не вдалося виявити мову. Будь ласка, введіть щось інше",
            },
            {
            LanguageCode: "pt",
                Text:         "Falha ao detectar o idioma. Por favor, insira outra coisa",
            },
            {
            LanguageCode: "id",
                Text:         "Gagal mendeteksi bahasa. Silakan masukkan sesuatu yang lain",
            },
            {
            LanguageCode: "it",
                Text:         "Impossibile rilevare la lingua. Per favore inserisci qualcos'altro",
            },
            {
            LanguageCode: "uz",
                Text:         "Til aniqlanmadi. Iltimos, yana bir narsani kiriting",
            },
            {
            LanguageCode: "de",
                Text:         "Die Sprache konnte nicht erkannt werden. Bitte geben Sie etwas anderes ein",
            },
        },
        "Images": {
            {
            LanguageCode: "de",
                Text:         "Bilder",
            },
            {
            LanguageCode: "uk",
                Text:         "Зображення",
            },
            {
            LanguageCode: "es",
                Text:         "Imagenes",
            },
            {
            LanguageCode: "id",
                Text:         "Gambar-gambar",
            },
            {
            LanguageCode: "it",
                Text:         "immagini",
            },
            {
            LanguageCode: "pt",
                Text:         "Imagens",
            },
            {
            LanguageCode: "ru",
                Text:         "Картинки",
            },
            {
            LanguageCode: "uz",
                Text:         "Tasvirlar",
            },
        },
        "Powered by": {
            {
                LanguageCode: "de",
                Text:         "Gesponsert von",
            },
            {
                LanguageCode: "es",
                Text:         "Patrocinado por",
            },
            {
                LanguageCode: "id",
                Text:         "Disponsori oleh",
            },
            {
                LanguageCode: "it",
                Text:         "Sponsorizzato da",
            },
            {
                LanguageCode: "pt",
                Text:         "Distribuído por",
            },
            {
                LanguageCode: "ru",
                Text:         "При поддержке",
            },
            {
                LanguageCode: "uk",
                Text:         "Спонсорується",
            },
            {
                LanguageCode: "uz",
                Text:         "Tomonidan qo'llab-quvvatlanadi",
            },
            {
                LanguageCode: "en",
                Text:         "Powered by",
            },
        },
        "sponsorship": {
            {
                LanguageCode: "de",
                Text:         "Werbung in @TransloBot\n\nUnter <b> jeder </b> übersetzten Nachricht befindet sich die Aufschrift \"Unterstützt von {Ihrem Hyperlink}\"\nWarum so?\n\nVorteile\n• Sie zahlen nur für <b> Vermögenswerte </b>\n• Deutlich <b> günstiger als der Versand </b>\n• Inline-Modus ist in den Bedingungen enthalten\n• Wenn ein Nutzer den Inline-Modus in einer Gruppe verwendet, wird Ihre Anzeige von allen Chat-Teilnehmern gesehen, das heißt <b> Zehntausende </b>\n\nGeben Sie Ihren Nachrichtentext ein,\n- 130 Zeichen inklusive Leerzeichen\n- ohne Fotos und Gifs\n- keine Wörter verwenden: abonnieren, reinkommen, gehen usw.\n- das Markup funktioniert nicht, d.h. fetter, kursiver Text wird wie gewohnt angezeigt",
            },
            {
                LanguageCode: "es",
                Text:         "Publicidad en @TransloBot\n\nDebajo de <b> cada </b> mensaje traducido habrá una inscripción \"Compatible con {su hipervínculo}\".\n¿Por qué de esta manera?\n\nVentajas\n• Paga solo por <b> activo </b>\n• Significativamente <b> más barato que el envío por correo </b>\n• El modo en línea está incluido en las condiciones\n• Si un usuario usa el modo en línea en un grupo, todos los participantes del chat ven su anuncio, es decir, <b> decenas de miles </b>.\n\nIngrese el texto de su mensaje,\n- 130 caracteres incluyendo espacios\n- sin fotos y gifs\n- no uses palabras: suscríbete, entra, ve, etc.\n- el marcado no funciona, es decir, el texto en negrita y cursiva será normal",
            },
            {
                LanguageCode: "id",
                Text:         "Beriklan di @TransloBot\n\nDi bawah <b> setiap </b> pesan yang diterjemahkan akan ada tulisan \"Didukung oleh {your hyperlink}\"\nMengapa dengan cara ini?\n\nKeuntungan\n• Anda hanya membayar <b> aset </b>\n• Secara signifikan <b> lebih murah daripada mengirim surat </b>\n• Mode sebaris termasuk dalam ketentuan\n• Jika pengguna menggunakan mode sebaris dalam grup, maka iklan Anda dilihat oleh semua peserta obrolan, yaitu <b> puluhan ribu </b>\n\nMasukkan teks pesan Anda,\n- 130 karakter termasuk spasi\n- tanpa foto dan gif\n- jangan menggunakan kata-kata: berlangganan, masuk, pergi, dll.\n- markup tidak berfungsi, yaitu huruf tebal, teks miring akan seperti biasa",
            },
            {
                LanguageCode: "it",
                Text:         "Pubblicità in @TransloBot\n\nSotto <b> ogni </b> messaggio tradotto ci sarà la scritta \"Supportato da {your hyperlink}\"\nPerché in questo modo?\n\nVantaggi\n• Paghi solo per <b> asset </b>\n• Significativamente <b> più economico della spedizione postale </b>\n• La modalità in linea è inclusa nelle condizioni\n• Se un utente utilizza la modalità in linea in un gruppo, il tuo annuncio viene visualizzato da tutti i partecipanti alla chat, ovvero <b> decine di migliaia </b>\n\nInserisci il testo del tuo messaggio,\n- 130 caratteri spazi inclusi\n- senza foto e gif\n- non usare parole: iscriviti, entra, vai, ecc.\n- il markup non funziona, ovvero il testo in grassetto e corsivo sarà normale",
            },
            {
                LanguageCode: "pt",
                Text:         "Publicidade em @TransloBot\n\nSob <b> cada </b> mensagem traduzida, haverá uma inscrição \"Apoiado por {seu hiperlink}\"\nPor que assim?\n\nVantagens\n• Você paga apenas por <b> ativos </b>\n• Significativamente <b> mais barato do que enviar pelo correio </b>\n• O modo inline está incluído nas condições\n• Se um usuário usa o modo inline em um grupo, seu anúncio é visto por todos os participantes do bate-papo, ou seja, <b> dezenas de milhares </b>\n\nDigite o texto da sua mensagem,\n- 130 caracteres incluindo espaços\n- sem fotos e gifs\n- não use palavras: assine, entre, vá, etc.\n- a marcação não funciona, ou seja, o texto em negrito e itálico ficará normal",
            },
            {
                LanguageCode: "ru",
                Text:         "Реклама в @TransloBot\n\nПод <b>каждым</b> переведенным сообщением будет надпись \"При поддержке {ваша гиперссылка}\"\nПочему такой способ?\n\nПреимущества\n• Вы платите только за <b>актив</b>\n• Значительно <b>дешевле рассылки</b>\n• Инлайн-режим включен в условия\n• Если пользователь пользуется инлайн-режимом в группе, то вашу рекламу видят все участники чата, то есть <b>десятки тысяч</b>\n\nВведите текст вашего сообщения, \n- 130 символов с учетом пробелов\n- без фото и гиф\n- не использовать слова : подписывайся, заходи, переходи и т.д.\n- разметка не работает, т.е жирный, курсивный текст будет как обычный",
            },
            {
                LanguageCode: "uk",
                Text:         "Реклама в @TransloBot\n\nПід <b> кожним </b> переведеним повідомленням буде напис \"За підтримки {ваша гіперпосилання}\"\nЧому такий спосіб?\n\nпереваги\n• Ви платите тільки за <b> актив </b>\n• Значно <b> дешевше розсилки </b>\n• Інлайн-режим включений в умови\n• Якщо користувач користується інлайн-режимом в групі, то вашу рекламу бачать всі учасники чату, тобто <b> десятки тисяч </b>\n\nВведіть текст вашого повідомлення,\n- 130 символів, враховуючи пробіли\n- без фото і гіф\n- не використовувати слова: підписуйся, заходь, переходь і т.д.\n- розмітка не працює, тобто жирний, курсивний текст буде як звичайний",
            },
            {
                LanguageCode: "uz",
                Text:         "@TransloBot -da reklama\n\n<b> Har bir </b> tarjima qilingan xabar ostida \"{sizning giperhavola} tomonidan qo'llab -quvvatlanadi\" yozuvi bo'ladi.\nNega bunday?\n\nAfzalliklar\n• Siz faqat <b> aktiv </b> uchun to'laysiz\n• sezilarli darajada <b> pochta jo'natishdan ko'ra arzon </b>\n• Inline rejimi shartlarga kiritilgan\n• Agar foydalanuvchi guruhda inline rejimidan foydalansa, sizning reklamangizni chatning barcha ishtirokchilari ko'radi, ya'ni <b> o'n minglab </b>.\n\nXabar matnini kiriting,\n- bo'shliqlarni o'z ichiga olgan 130 ta belgi\n- fotosuratlar va sovg'alarsiz\n- so'zlarni ishlatmang: obuna bo'ling, kiring, kiring va hokazo.\n- belgilash ishlamaydi, ya'ni qalin, kursiv matn odatdagidek bo'ladi",
            },
            {
                LanguageCode: "en",
                Text:         "Advertising in @TransloBot\n\nUnder <b> each </b> translated message there will be an inscription \"Supported by {your hyperlink}\"\nWhy this way?\n\nAdvantages\n• You pay only for <b> asset </b>\n• Significantly <b> cheaper than mailing </b>\n• Inline mode is included in the conditions\n• If a user uses inline mode in a group, then your ad is seen by all chat participants, that is, <b> tens of thousands </b>\n\nEnter your message text,\n- 130 characters including spaces\n- without photos and gifs\n- do not use words: subscribe, come in, go, etc.\n- the markup does not work, i.e. bold, italic text will be as normal",
            },
        },
        "sponsorship_set_days": {
            {
                LanguageCode: "de",
                Text:         "Senden Sie nun die Anzahl der Tage, die Sie inserieren möchten (die Tage beginnen mit dem Zahlungsdatum). 1-30, ganze Zahl.",
            },
            {
                LanguageCode: "es",
                Text:         "Ahora envíe la cantidad de días que desea anunciar (los días comenzarán a partir de la fecha de pago). 1-30, número entero.",
            },
            {
                LanguageCode: "id",
                Text:         "Sekarang kirim jumlah hari Anda ingin beriklan (hari akan dimulai dari tanggal pembayaran). 1-30, bilangan bulat.",
            },
            {
                LanguageCode: "it",
                Text:         "Ora invia il numero di giorni che vuoi pubblicizzare (i giorni inizieranno dalla data di pagamento). 1-30, intero.",
            },
            {
                LanguageCode: "pt",
                Text:         "Envie agora a quantidade de dias que deseja anunciar (os dias contarão a partir da data do pagamento). 1-30, número inteiro.",
            },
            {
                LanguageCode: "ru",
                Text:         "Теперь отправьте количество дней, которые хотите рекламироваться (отсчёт дней начнётся со дня оплаты). 1-30, целое число.",
            },
            {
                LanguageCode: "uk",
                Text:         "Тепер відправте кількість днів, які хочете рекламуватися (відлік днів почнеться з дня оплати). 1-30, ціле число.",
            },
            {
                LanguageCode: "uz",
                Text:         "Endi reklama qilmoqchi bo'lgan kunlar sonini yuboring (kunlar to'lov sanasidan boshlanadi). 1-30, butun son.",
            },
            {
                LanguageCode: "en",
                Text:         "Now send the number of days you want to advertise (the days will start from the date of payment). 1-30, integer.",
            },
        },
        "Введите целое число без лишних символов": {
            {
                LanguageCode: "de",
                Text:         "Geben Sie eine ganze Zahl ohne zusätzliche Zeichen ein",
            },
            {
                LanguageCode: "es",
                Text:         "Ingrese un número entero sin caracteres adicionales",
            },
            {
                LanguageCode: "id",
                Text:         "Masukkan bilangan bulat tanpa karakter tambahan",
            },
            {
                LanguageCode: "it",
                Text:         "Inserisci un numero intero senza caratteri extra",
            },
            {
                LanguageCode: "pt",
                Text:         "Insira um número inteiro sem caracteres extras",
            },
            {
                LanguageCode: "ru",
                Text:         "Введите целое число без лишних символов",
            },
            {
                LanguageCode: "uk",
                Text:         "Введіть ціле число без зайвих символів",
            },
            {
                LanguageCode: "uz",
                Text:         "Qo'shimcha belgilarsiz butun sonni kiriting",
            },
            {
                LanguageCode: "en",
                Text:         "Enter an integer without extra characters",
            },
        },
        "Количество дней должно быть от 1 до 30 включительно": {
            {
                LanguageCode: "de",
                Text:         "Die Anzahl der Tage muss zwischen 1 und 30 liegen",
            },
            {
                LanguageCode: "es",
                Text:         "El número de días debe ser de 1 a 30 inclusive.",
            },
            {
                LanguageCode: "id",
                Text:         "Jumlah hari harus dari 1 hingga 30 inklusif",
            },
            {
                LanguageCode: "it",
                Text:         "Il numero di giorni deve essere compreso tra 1 e 30 inclusi",
            },
            {
                LanguageCode: "pt",
                Text:         "O número de dias deve ser de 1 a 30, inclusive",
            },
            {
                LanguageCode: "ru",
                Text:         "Количество дней должно быть от 1 до 30 включительно",
            },
            {
                LanguageCode: "uk",
                Text:         "Кількість днів має бути від 1 до 30 включно",
            },
            {
                LanguageCode: "uz",
                Text:         "Kunlar soni 1 dan 30 gacha bo'lishi kerak",
            },
            {
                LanguageCode: "en",
                Text:         "The number of days must be from 1 to 30 inclusive",
            },
        },
        "Выберите страны пользователей, которые получат вашу рассылку. Когда закончите, нажмите кнопку Далее": {
            {
                LanguageCode: "de",
                Text:         "Wählen Sie die Länder der Benutzer aus, die Ihren Newsletter erhalten. Wenn Sie fertig sind, klicken Sie auf Weiter",
            },
            {
                LanguageCode: "es",
                Text:         "Seleccione los países de los usuarios que recibirán su newsletter. Cuando termine, haga clic en Siguiente",
            },
            {
                LanguageCode: "id",
                Text:         "Pilih negara pengguna yang akan menerima buletin Anda. Setelah selesai, klik Berikutnya",
            },
            {
                LanguageCode: "it",
                Text:         "Seleziona i paesi degli utenti che riceveranno la tua newsletter. Al termine, fai clic su Avanti",
            },
            {
                LanguageCode: "pt",
                Text:         "Selecione os países dos usuários que receberão sua newsletter. Quando terminar, clique em Avançar",
            },
            {
                LanguageCode: "ru",
                Text:         "Выберите страны пользователей, которые получат вашу рассылку. Когда закончите, нажмите кнопку Далее",
            },
            {
                LanguageCode: "uk",
                Text:         "Виберіть країни користувачів, які отримають вашу розсилку. Коли закінчите, натисніть кнопку Далі",
            },
            {
                LanguageCode: "uz",
                Text:         "Sizning axborot byulleteningizni oladigan foydalanuvchilarning mamlakatlarini tanlang. Ish tugagach, Keyingiga bosing",
            },
            {
                LanguageCode: "en",
                Text:         "Select the countries of the users who will receive your newsletter. When done, click Next",
            },
        },
        "Далее": {
            {
                LanguageCode: "de",
                Text:         "Weiter ▶",
            },
            {
                LanguageCode: "es",
                Text:         "Siguiente ▶",
            },
            {
                LanguageCode: "id",
                Text:         "Selanjutnya",
            },
            {
                LanguageCode: "it",
                Text:         "Avanti ▶",
            },
            {
                LanguageCode: "pt",
                Text:         "Próximo ▶",
            },
            {
                LanguageCode: "ru",
                Text:         "Далее ▶",
            },
            {
                LanguageCode: "uk",
                Text:         "Далі ▶",
            },
            {
                LanguageCode: "uz",
                Text:         "Keyingi ▶",
            },
            {
                LanguageCode: "en",
                Text:         "Next ▶",
            },
        },
}
    
    if df, ok := languages[text]; ok { // Текст подходит под варианты
        if v, ok := CorrespLang(&df, lang); ok {
            return fmt.Sprintf(v, placeholders...)
        }
    }
   return fmt.Sprintf(text, placeholders...)
}