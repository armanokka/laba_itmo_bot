package main

import "fmt"

func Localize(text, lang string, placeholders ...string) string {
    if lang == "" || lang == "en" {
        return fmt.Sprintf(text, placeholders)
    }
    
    var languages = map[string]map[string]string{
        "Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)": {
            "ru": "Извините, произошла ошибка.\n\nПожалуйста, не блокируйте бота, я исправлю ошибку в ближайшее время, администратор уже предупрежден об этой ошибке ;)",
            "es": "Lo sentimos, error causado.\n\nPor favor, no bloquees el bot, corregiré el error en un futuro cercano, el administrador ya ha sido advertido sobre este error ;)",
            "uk": "Вибачте, сталася помилка.\n\nБудь ласка, не блокуйте бота, я виправлю помилку найближчим часом, адміністратор вже попереджений про цю помилку ;)",
            "pt": "Desculpa, erro causado.\n\nPor favor, não bloqueie o bot, eu vou corrigir o bug em um futuro próximo, o administrador já foi avisado sobre este ",
            "id": "Maaf, kesalahan disebabkan.\n\nTolong, jangan halangi robot itu, aku akan memperbaiki bug dalam waktu dekat, administrator telah diperingatkan tentang kesalahan ini ;)",
            "it": "Scusa, errore causato.\n\nPer favore, non bloccare il bot, correggerò il bug nel prossimo futuro, l'amministratore è già stato avvertito di questo errore;)",
            "uz": "Uzr, xato sabab.\n\nIltimos, botni bloklamang, yaqin kelajakda xatoni tuzataman, administrator bu xato haqida allaqachon ogohlantirilgan ;)",
            "de": "Sorry, Fehler verursacht.\n\nBitte blockieren Sie den Bot nicht, ich werde den Fehler in naher Zukunft beheben, der Administrator wurde bereits vor diesem Fehler gewarnt ;)",
        },
        "Your language is - <b>%s</b>, and the language for translation is - <b>%s</b>.": {
            "ru": "Ваш язык - <b>%s</b>, а язык для перевода - <b>%s</b>.",
            "es": "Su idioma es - <b>%s</b>, y el idioma de la traducción es - <b>%s</b>.",
            "uk": "Ваша мова - <B>%s</b>, а мова для перекладу - <B>%s</b>.",
            "pt": "A sua língua é - <b>%s</b>, e a língua para a tradução é- <b>%s</b>.",
            "id": "Bahasa Anda - <b>%s</b>, dan bahasa untuk Terjemahan adalah - <b>%s</b>.",
            "it":"La tua lingua è - < b > %s< / b>, e la lingua per la traduzione è - <b>%s</b>.",
            "uz":"Sizning til - <b>%s</b>, va - <b tarjima qilish uchun til>%bo'ladi s</b>.",
            "de":"Ihre Sprache ist - <b>%s</b> und die Sprache für die Übersetzung ist - <b>%s< / b>.",
        },
        "💡Instruction": {
            "ru":"💡Инструкция",
            "es":"💡Instrucción",
            "uk":"💡Інструкція",
            "pt":"💡Instrucao",
            "id":"💡Instruksi",
            "it":"💡Istruzione",
            "uz":"💡Yo'riqnoma",
            "de":"💡Anweisung",
        },
        "My Language": {
            "ru":"Мой Язык",
            "es":"Mi Idioma",
            "uk":"Моя Мова",
            "pt":"A Minha Língua",
            "id":"Bahasa Saya",
            "it":"La mia lingua",
            "uz":"Tilimni",
            "de":"Meine Sprache",
        },
        "Translate Language": {
            "de": "Sprache übersetzen",
            "es": "Traducir idioma",
            "id": "Terjemahkan Bahasa",
            "it": "Traduci lingua",
            "pt": "Traduzir idioma",
            "ru": "Перевести язык",
            "uk": "Перекласти мову",
            "uz": "Tarjima qilish tili",
        },
        "To setup <b>your language</b>, do <b>one</b> of the following: 👇\n\nℹ️ Send <b>few words</b> in <b>your</b> language, for example: \"<code>Hi, how are you today?</code>\" - language will be English, or \"<code>L'amour ne fait pas d'erreurs</code>\" - language will be French, and so on.\nℹ️ Or send the <b>name</b> of your language <b>in English</b>, e.g. \"<code>Russian</code>\", or \"<code>Japanese</code>\", or  \"<code>Arabic</code>\", e.t.c.": {
            "uk": "Щоб налаштувати <b>свою мову</b>, виконайте <b>одне</b>з наступного: 👇\n\nℹ️ Надішліть <b>кілька слів</b><b>вашою</b>мовою, наприклад: \"<code>Привіт, як справи сьогодні?</code>\" - мова буде англійською, або \"< code> L'amour ne fait pas d'erreurs</code> \"- мовою буде французька тощо.\nℹ️ Або надішліть <b>назву</b>вашої мови <b>англійською мовою</b>, напр. \"<code>російська</code>\", або \"<code>японська</code>\", або \"<code>арабська</code>\", напр.",
            "pt": "Para configurar <b>seu idioma</b>, faça <b>um</b>dos seguintes: 👇\n\nℹ️ Envie <b>algumas palavras</b>em <b>seu</b>idioma, por exemplo: \"<code>Olá, como você está hoje?</code>\" - o idioma será inglês ou \"< code> L'amour ne fait pas d'erreurs</code> \"- o idioma será o francês e assim por diante.\nℹ️ Ou envie o <b>nome</b>do seu idioma <b>em inglês</b>, por exemplo \"<code>Russo</code>\", ou \"<code>Japonês</code>\" ou \"<code>Árabe</code>\", e.t.c.",
            "id": "Untuk menyiapkan <b>bahasa Anda</b>, lakukan <b>salah satu</b> berikut ini:\n\n️ Kirim <b>beberapa kata</b> dalam bahasa <b>Anda</b>, misalnya: \"<code>Hai, apa kabar hari ini?</code>\" - bahasa akan menjadi bahasa Inggris, atau \"< code>L'amour ne fait pas d'erreurs</code>\" - bahasa akan menjadi bahasa Prancis, dan seterusnya.\n️ Atau kirimkan <b>nama</b> bahasa Anda <b>dalam bahasa Inggris</b>, mis. \"<code>Rusia</code>\", atau \"<code>Jepang</code>\", atau \"<code>Arab</code>\", dll.",
            "it": "Per impostare la <b>lingua</b>, esegui <b>una</b> delle seguenti operazioni: 👇\n\nℹ️ Invia <b>poche parole</b> nella <b>tua</b> lingua, ad esempio: \"<code>Ciao, come stai oggi?</code>\" - la lingua sarà inglese, o \"< code>L'amour ne fait pas d'erreurs</code>\" - la lingua sarà il francese e così via.\nℹ️ Oppure invia il <b>nome</b> della tua lingua <b>in inglese</b>, ad es. \"<code>Russo</code>\", o \"<code>Giapponese</code>\", o \"<code>Arabo</code>\", ecc.",
            "uz": "O'zingizning tilingizni</b>sozlash uchun quyidagilardan <b>birini</b>bajaring: 👇\n\nℹ️ o'z tilingizda <b>bir nechta so'zlar</b>ni yuboring, masalan: \"<code>salom, bugun yaxshimisiz?</code>\" - bu til inglizcha bo'ladi yoki \"<\" code> L'amour ne fait pas d'erreurs</code> \"- til frantsuzcha bo'ladi va hokazo.\nℹ️ Yoki o'z tilingiz <b>ning ismini <b>ingliz tilida</b>yuboring, masalan. \"<code>Russian</code>\", yoki \"<code>Japanese</code>\", or \"<code>Arabic</code>\", e.t.c.",
            "de": "Um <b>Ihre Sprache</b> einzurichten, führen Sie <b>einen</b> der folgenden Schritte aus: 👇\n\nℹ️ Senden Sie <b>ein paar Wörter</b> in <b>Ihrer</b> Sprache, zum Beispiel: \"<code>Hallo, wie geht es Ihnen heute?</code>\" - Sprache ist Englisch oder \"< code>L'amour ne fait pas d'erreurs</code>\" - Sprache wird Französisch sein und so weiter.\nℹ️ Oder sende den <b>Namen</b> deiner Sprache <b>auf Englisch</b>, z.B. \"<code>Russisch</code>\", oder \"<code>Japanisch</code>\", oder \"<code>Arabisch</code>\", usw.",
            "ru": "Чтобы настроить <b>свой язык</b>, выполните <b>одно</b>из следующего: 👇\n\nℹ️ Отправьте <b>несколько слов</b>на <b>своем</b>языке, например: \"<code>Привет, как дела?</code>\" - язык будет английский, или \"<code>L'amour ne fait pas d'erreurs</code> \"- язык будет французский, и так далее.\nℹ️ Или отправьте <b>название</b>вашего языка <b>на английском</b>, например \"<code>русский</code>\", или \"<code>японский</code>\", или \"<code>арабский</code>\" и т. д.",
            "es": "Para configurar <b>su idioma</b>, haga <b>una</b>de las siguientes opciones: 👇\n\nℹ️ Envíe <b>algunas palabras</b>en <b>su</b>idioma, por ejemplo: \"<code>Hola, ¿cómo está hoy?</code>\" - el idioma será inglés o \"< code> L'amour ne fait pas d'erreurs</code> \"- el idioma será el francés, y así sucesivamente.\nℹ️ O envíe el <b>nombre</b>de su idioma <b>en inglés</b>, p. ej. \"<code>Ruso</code>\", o \"<code>Japonés</code>\", o \"<code>Árabe</code>\", e.t.c.",
        },
        "To setup <b>translate language</b>, do <b>one</b> of the following: 👇\n\nℹ️ Send <b>few words</b> in language <b>into you want to translate</b>, for example: \"<code>Hi, how are you today?</code>\" - language will be English, or \"<code>L'amour ne fait pas d'erreurs</code>\" - language will be French, and so on.\nℹ️ Or send the <b>name</b> of language <b>into you want to translate, in English</b>, e.g. \"<code>Russian</code>\", or \"<code>Japanese</code>\", or  \"<code>Arabic</code>\", e.t.c.": {
            "ru": "Чтобы настроить <b>язык перевода</b>, выполните <b>одно</b>из следующего: 👇\n\nℹ️ Отправьте <b>несколько слов</b>на языке <b>, на который вы хотите перевести</b>, например: «<code>Привет, как дела?</code>» - язык будет английский. , или \"<code>L'amour ne fait pas d'erreurs</code>\" - языком будет французский и т. д.\nℹ️ Или отправьте <b>название</b>языка <b>, на который вы хотите перевести, на английский</b>, например \"<code>русский</code>\", или \"<code>японский</code>\", или \"<code>арабский</code>\" и т. д.",
            "es": "Para configurar el <b>idioma de traducción</b>, realice <b>una</b>de las siguientes acciones: 👇\n\nℹ️ Envíe <b>algunas palabras</b>en el idioma <b>al que desea traducir</b>, por ejemplo: \"<code>Hola, ¿cómo estás hoy?</code>\" - el idioma será inglés , o \"<code>L'amour ne fait pas d'erreurs</code>\" - el idioma será el francés, y así sucesivamente.\nℹ️ O envíe el <b>nombre</b>del idioma <b>al que desea traducir, en inglés</b>, p. ej. \"<code>Ruso</code>\", o \"<code>Japonés</code>\", o \"<code>Árabe</code>\", e.t.c.",
            "uk": "Щоб встановити <b>переклад мови</b>, виконайте <b>одне</b>з наступного: 👇\n\nℹ️ Надішліть <b>кілька слів</b>мовою <b>на ту, яку хочете перекласти</b>, наприклад: «<code>Привіт, як справи сьогодні?</code>« - мова буде англійською , або \"<code>L'amour ne fait pas d'erreurs</code>\" - мовою буде французька тощо.\nℹ️ Або надішліть <b>назву</b>мови <b>на ту мову, яку потрібно перекласти, англійською мовою</b>, напр. \"<code>російська</code>\", або \"<code>японська</code>\", або \"<code>арабська</code>\", напр.",
            "pt": "Para configurar a <b>tradução do idioma</b>, faça <b>um</b>dos seguintes: 👇\n\nℹ️ Envie <b>algumas palavras</b>no idioma <b>para o qual deseja traduzir</b>, por exemplo: \"<code>Olá, como você está hoje?</code>\" - o idioma será o inglês ou \"<code>L'amour ne fait pas d'erreurs</code>\" - o idioma será o francês e assim por diante.\nℹ️ Ou envie o <b>nome</b>do idioma <b>para o qual deseja traduzir, em inglês</b>, por exemplo \"<code>Russo</code>\", ou \"<code>Japonês</code>\" ou \"<code>Árabe</code>\", e.t.c.",
            "id": "Untuk menyiapkan <b>terjemahkan bahasa</b>, lakukan <b>salah satu</b> berikut ini:\n\n️ Kirim <b>beberapa kata</b> dalam bahasa <b>ke dalam bahasa yang ingin Anda terjemahkan</b>, misalnya: \"<code>Hai, apa kabar hari ini?</code>\" - bahasa akan menjadi bahasa Inggris , atau \"<code>L'amour ne fait pas d'erreurs</code>\" - bahasa akan menjadi bahasa Prancis, dan seterusnya.\n️ Atau kirimkan <b>nama</b> bahasa <b>ke dalam bahasa yang ingin Anda terjemahkan, dalam bahasa Inggris</b>, mis. \"<code>Rusia</code>\", atau \"<code>Jepang</code>\", atau \"<code>Arab</code>\", dll.",
            "it": "Per impostare la <b>lingua di traduzione</b>, esegui <b>una</b> delle seguenti operazioni: 👇\n\nℹ️ Invia <b>poche parole</b> nella lingua <b>nella quale vuoi tradurre</b>, ad esempio: \"<code>Ciao, come stai oggi?</code>\" - la lingua sarà l'inglese , o \"<code>L'amour ne fait pas d'erreurs</code>\" - la lingua sarà il francese e così via.\nℹ️ Oppure invia il <b>nome</b> della lingua <b>in cui desideri tradurre, in inglese</b>, ad es. \"<code>Russo</code>\", o \"<code>Giapponese</code>\", o \"<code>Arabo</code>\", ecc.",
            "uz": "Tilni tarjima qilish</b>ni sozlash uchun quyidagilardan <b>birini</b>bajaring: 👇\n\nℹ️ tarjima qilmoqchi bo'lgan tilingizga <b>bir nechta so'zlarni <b>yuboring</b>, masalan: \"<kod> salom, bugun yaxshimisiz?</code>\" - til ingliz tilida bo'ladi , yoki \"<code>L'amour ne fait pas d'erreurs</code>\" - til frantsuzcha bo'ladi va hokazo.\nℹ️ Yoki tarjima qilmoqchi bo'lgan tilga <b>ismini</b>ingliz tilida</b>yuboring, masalan. \"<code>Russian</code>\", yoki \"<code>Japanese</code>\", or \"<code>Arabic</code>\", e.t.c.",
            "de": "Um die <b>Übersetzungssprache</b> einzurichten, führen Sie <b>einen</b> der folgenden Schritte aus:\n\nℹ️ Senden Sie <b>wenige Wörter</b> in der Sprache <b>in die Sie übersetzen möchten</b>, zum Beispiel: \"<code>Hallo, wie geht es Ihnen heute?</code>\" - die Sprache wird Englisch sein , oder \"<code>L'amour ne fait pas d'erreurs</code>\" - die Sprache ist Französisch und so weiter.\nℹ️ Oder senden Sie den <b>Namen</b> der Sprache <b>in die Sie übersetzen möchten, auf Englisch</b>, z.B. \"<code>Russisch</code>\", oder \"<code>Japanisch</code>\", oder \"<code>Arabisch</code>\", usw.",
        },
        "⬅Back":{},
        "<b>What can this bot do?</b>\n▫️ Translo allow you to translate messages into 182+ languages.\n<b>How to translate message?</b>\n▫️ Firstly, you have to setup your language, then setup translate language, next send text messages and bot will translate them quickly.\n<b>How to setup my language?</b>\n▫️ Click on the button below called \"My Language\"\n<b>How to setup language into I want to translate?</b>\n▫️ Click on the button below called \"Translate Language\"\n<b>Are there any other interesting things?\n</b>▫️ Yes, the bot support inline mode. Start typing the nickname @translobot in the message input field and then write there the text you want to translate.\n<b>I have a suggestion or I found bug!</b>\n▫️ 👉 Contact me pls - @armanokka": {
            "it": "<b>Cosa può fare questo bot?</b>\n▫️ Translo ti consente di tradurre i messaggi in oltre 182 lingue.\n<b>Come tradurre il messaggio?</b>\n▫️ In primo luogo, devi impostare la tua lingua, quindi impostare la lingua di traduzione, quindi inviare messaggi di testo e il bot li tradurrà rapidamente.\n<b>Come impostare la mia lingua?</b>\n▫️ Clicca sul pulsante in basso chiamato \"La mia lingua\"\n<b>Come impostare la lingua in cui voglio tradurre?</b>\n▫️ Clicca sul pulsante in basso chiamato \"Traduci lingua\"\n<b>Ci sono altre cose interessanti?\n</b>▫️ Sì, il bot supporta la modalità in linea. Inizia a digitare il nickname @translobot nel campo di inserimento del messaggio e poi scrivi lì il testo che vuoi tradurre.\n<b>Ho un suggerimento o ho trovato un bug!</b>\n▫️ 👉 Contattami per favore - @armanokka",
            "uz": "<b>Ushbu bot nima qila oladi?</b>\n▫️ Translo sizga xabarlarni 182+ tilga tarjima qilishga imkon beradi.\n<b>Xabarni qanday tarjima qilish kerak?</b>\n▫️ Birinchidan, siz o'z tilingizni o'rnatishingiz kerak, so'ngra tarjima tilini sozlashingiz kerak, keyin matnli xabarlarni yuboring va bot ularni tezda tarjima qiladi.\n<b>Mening tilimni qanday sozlashim kerak?</b>\n▫️ Quyidagi \"Mening tilim\" deb nomlangan tugmani bosing.\n<b>Tilni qanday tarjima qilishni xohlayman?</b>\n▫️ \"Tilni tarjima qilish\" deb nomlangan tugmani bosing.\n<b>Boshqa qiziqarli narsalar bormi?\n</b> ▫️ Ha, botni qo'llab-quvvatlash inline rejimida. Xabarlarni kiritish maydoniga @translobot taxallusini yozishni boshlang va keyin tarjima qilmoqchi bo'lgan matni yozing.\n<b>Menda bir taklif bor yoki men xato topdim!</b>\n▫️ 👉 Men bilan bog'laning pls - @armanokka",
            "de": "<b>Was kann dieser Bot tun?</b>\n▫️ Mit Translo können Sie Nachrichten in mehr als 182 Sprachen übersetzen.\n<b>Wie übersetzt man eine Nachricht?</b>\n▫️ Zuerst müssen Sie Ihre Sprache einrichten, dann die Übersetzungssprache einrichten, als nächstes Textnachrichten senden und der Bot wird sie schnell übersetzen.\n<b>Wie richte ich meine Sprache ein?</b>\n▫️ Klicken Sie unten auf die Schaltfläche \"Meine Sprache\"\n<b>Wie richte ich die Sprache ein, in die ich übersetzen möchte?</b>\n▫️ Klicken Sie unten auf die Schaltfläche \"Sprache übersetzen\"\n<b>Gibt es noch andere interessante Dinge?\n</b>▫️ Ja, der Bot unterstützt den Inline-Modus. Geben Sie den Spitznamen @translobot in das Nachrichteneingabefeld ein und schreiben Sie dort den Text, den Sie übersetzen möchten.\n<b>Ich habe einen Vorschlag oder ich habe einen Fehler gefunden!</b>\n▫️ 👉 Kontaktieren Sie mich bitte - @armanokka",
            "ru": "<b>Что умеет этот бот?</b>\n▫️ Translo позволяет переводить сообщения на 182+ языков.\n<b>Как перевести сообщение?</b>\n▫️ Во-первых, вам нужно настроить свой язык, затем настроить язык перевода, затем отправить текстовые сообщения, и бот быстро их переведет.\n<b>Как настроить мой язык?</b>\n▫️ Нажмите кнопку ниже под названием «Мой язык».\n<b>Как установить язык, на котором я хочу переводить?</b>\n▫️ Нажмите кнопку ниже под названием «Перевести язык».\n<b>Есть еще что-нибудь интересное?\n</b> ▫️ Да, бот поддерживает встроенный режим. Начните вводить псевдоним @translobot в поле ввода сообщения, а затем впишите туда текст, который хотите перевести.\n<b>У меня есть предложение или я обнаружил ошибку!</b>\n▫️ 👉 Свяжитесь со мной, пожалуйста - @armanokka",
            "es": "<b>¿Qué puede hacer este bot?</b>\n▫️ Translo te permite traducir mensajes a más de 182 idiomas.\n<b>¿Cómo traducir un mensaje?</b>\n▫️ En primer lugar, debe configurar su idioma, luego configurar el idioma de traducción, luego enviar mensajes de texto y el bot los traducirá rápidamente.\n<b>¿Cómo configurar mi idioma?</b>\n▫️ Haga clic en el botón de abajo llamado \"Mi idioma\"\n<b>¿Cómo configurar el idioma al que quiero traducir?</b>\n▫️ Haga clic en el botón de abajo llamado \"Traducir idioma\"\n<b>¿Hay otras cosas interesantes?\n</b> ▫️ Sí, el bot admite el modo en línea. Comience a escribir el apodo @translobot en el campo de entrada del mensaje y luego escriba allí el texto que desea traducir.\n<b>¡Tengo una sugerencia o encontré un error!</b>\n▫️ 👉 Contáctame por favor - @armanokka",
            "uk": "<b>Що може зробити цей бот?</b>\n▫️ Translo дозволяє перекладати повідомлення 182+ мовами.\n<b>Як перекласти повідомлення?</b>\n▫️ По-перше, вам потрібно налаштувати свою мову, потім налаштувати переклад мови, наступне надіслати текстові повідомлення і бот швидко їх перекладе.\n<b>Як налаштувати мову?</b>\n▫️ Клацніть на кнопку під назвою \"Моя мова\"\n<b>Як налаштувати мову на мову, яку я хочу перекласти?</b>\nКлацніть на кнопку нижче під назвою \"Перекласти мову\"\n<b>Чи є якісь інші цікаві речі?\n</b> ▫️ Так, бот підтримує вбудований режим. Почніть вводити псевдонім @translobot у поле введення повідомлення, а потім напишіть туди текст, який потрібно перекласти.\n<b>У мене є пропозиція або я знайшов помилку!</b>\n👉 Зв’яжіться зі мною pls - @armanokka",
            "pt": "<b>O que este bot pode fazer?</b>\n▫️ Translo permite que você traduza mensagens em mais de 182 idiomas.\n<b>Como traduzir mensagem?</b>\n▫️ Em primeiro lugar, você deve configurar seu idioma, depois configurar traduzir o idioma, em seguida enviar mensagens de texto e o bot irá traduzi-las rapidamente.\n<b>Como configurar meu idioma?</b>\n▫️ Clique no botão abaixo chamado \"Meu Idioma\"\n<b>Como configurar o idioma para o qual desejo traduzir?</b>\n▫️ Clique no botão abaixo chamado \"Traduzir Idioma\"\n<b>Existem outras coisas interessantes?\n</b> ▫️ Sim, o bot suporta o modo inline. Comece digitando o apelido @translobot no campo de entrada da mensagem e então escreva lá o texto que deseja traduzir.\n<b>Tenho uma sugestão ou encontrei um bug!</b>\n▫️ 👉 Contate-me, pls - @armanokka",
            "id": "<b>Apa yang bisa dilakukan bot ini?</b>\n️ Translo memungkinkan Anda menerjemahkan pesan ke dalam 182+ bahasa.\n<b>Bagaimana cara menerjemahkan pesan?</b>\n️ Pertama, Anda harus mengatur bahasa Anda, lalu mengatur bahasa terjemahan, selanjutnya mengirim pesan teks dan bot akan menerjemahkannya dengan cepat.\n<b>Bagaimana cara mengatur bahasa saya?</b>\n️ Klik tombol di bawah yang disebut \"Bahasa Saya\"\n<b>Bagaimana cara mengatur bahasa yang ingin saya terjemahkan?</b>\n️ Klik tombol di bawah yang disebut \"Terjemahkan Bahasa\"\n<b>Apakah ada hal menarik lainnya?\n</b>▫️ Ya, bot mendukung mode sebaris. Mulai ketikkan nickname @translobot di kolom input pesan lalu tulis di sana teks yang ingin Anda terjemahkan.\n<b>Saya punya saran atau saya menemukan bug!</b>\n️ Hubungi saya pls - @armanokka",
        },
        "Failed to detect the language. Please enter something else": {
            "it": "Impossibile rilevare la lingua. Per favore inserisci qualcos'altro",
            "uz": "Til aniqlanmadi. Iltimos, yana bir narsani kiriting",
            "de": "Die Sprache konnte nicht erkannt werden. Bitte geben Sie etwas anderes ein",
            "ru": "Не удалось определить язык. Пожалуйста, введите что-нибудь еще",
            "es": "No se pudo detectar el idioma. Por favor ingrese algo más",
            "uk": "Не вдалося виявити мову. Будь ласка, введіть щось інше",
            "pt": "Falha ao detectar o idioma. Por favor, insira outra coisa",
            "id": "Gagal mendeteksi bahasa. Silakan masukkan sesuatu yang lain",
        },
        "Now your language is %s\n\nPress \"⬅Back\" to exit to menu": {
            "es": "Ahora tu idioma es %s\n\nPresione \"⬅Back\" para salir al menú",
            "uk": "Зараз ваша мова% %s\n\nНатисніть \"⬅Back\", щоб вийти в меню",
            "pt": "Agora seu idioma é %s\n\nPressione \"⬅Back\" para sair para o menu",
            "id": "Sekarang bahasa Anda adalah %s\n\nTekan \"⬅Back\" untuk keluar ke menu",
            "it": "Ora la tua lingua è %s\n\nPremere \"⬅Back\" per uscire dal menu",
            "uz": "Endi sizning tilingiz %s\n\nMenyuga chiqish uchun \"⬅Back\" tugmasini bosing",
            "de": "Ihre Sprache ist jetzt %s\n\nDrücken Sie \"⬅Back\", um das Menü zu verlassen",
            "ru": "Теперь ваш язык %s\n\nНажмите «⬅Back» для выхода в меню.",
        },
        "Now translate language is %s\n\nPress \"⬅Back\" to exit to menu": {
            "it": "Ora la lingua di traduzione è %s\n\nPremere \"⬅Back\" per uscire dal menu",
            "uz": "Endi tarjima qilish tili %s\n\nMenyuga chiqish uchun \"⬅Back\" tugmasini bosing",
            "de": "Die Übersetzungssprache ist jetzt %s\n\nDrücken Sie \"⬅Back\", um das Menü zu verlassen",
            "ru": "Теперь язык перевода %s\n\nНажмите «⬅Back» для выхода в меню.",
            "es": "Ahora el idioma de traducción es %s\n\nPresione \"⬅Back\" para salir al menú",
            "uk": "Тепер мова перекладу - %s\n\nНатисніть \"⬅Back\", щоб вийти в меню",
            "pt": "Agora o idioma de tradução é %s\n\nPressione \"⬅Back\" para sair para o menu",
            "id": "Sekarang bahasa terjemahan adalah %s\n\nTekan \"⬅Back\" untuk keluar ke menu",
        },
        "⏳ Translating...": {
            "uz": "⏳ Tarjima qilish...",
            "de": "⏳ Übersetzen...",
            "ru": "⏳ Перевожу...",
            "es": "⏳ Traducir...",
            "uk": "⏳ Перекласти...",
            "pt": "⏳ Traduzir...",
            "id": "⏳ Terjemahkan...",
            "it": "⏳ Traduci...",
        },
        "Please, send text message": {
            "uz": "Iltimos, SMS yuboring",
            "de": "Bitte SMS senden",
            "ru": "Пожалуйста, отправьте текстовое сообщение",
            "es": "Por favor, envíe un mensaje de texto",
            "uk": "Будь ласка, надішліть текстове повідомлення",
            "pt": "Por favor, envie mensagem de texto",
            "id": "Tolong, kirim pesan teks",
            "it": "Per favore, invia un messaggio di testo",
        },
        "Too big text": {
            "es": "Texto demasiado grande",
            "uk": "Занадто великий текст",
            "pt": "Texto muito grande",
            "id": "Teks terlalu besar",
            "it": "Testo troppo grande",
            "uz": "Juda katta matn",
            "de": "Zu großer Text",
            "ru": "Слишком большой текст",
        },
        "Unsupported language or internal error": {
            "ru": "Неподдерживаемый язык или внутренняя ошибка",
            "es": "Idioma no admitido o error interno",
            "uk": "Непідтримувана мова або внутрішня помилка",
            "pt": "Idioma não suportado ou erro interno",
            "id": "Bahasa yang tidak didukung atau kesalahan internal",
            "it": "Lingua non supportata o errore interno",
            "uz": "Qo'llab-quvvatlanmaydigan til yoki ichki xato",
            "de": "Nicht unterstützte Sprache oder interner Fehler",
        },
        "Empty result": {
            "ru": "Пустой результат",
            "es": "Resultado vacío",
            "uk": "Порожній результат",
            "pt": "Resultado vazio",
            "id": "Hasil kosong",
            "it": "Risultato vuoto",
            "uz": "Bo'sh natija",
            "de": "Leeres Ergebnis",
        },
        "Error, sorry.": {
            "uk": "Помилка, вибачте.",
            "pt": "Erro, desculpe.",
            "id": "Kesalahan, maaf.",
            "it": "Errore, scusa.",
            "uz": "Xato, uzr.",
            "de": "Fehler, tut mir leid.",
            "ru": "Ошибка, извините.",
            "es": "Error, lo siento.",
        },
    }
    if df, ok := languages[text]; ok {
        if v, ok := df[lang]; ok {
            if len(placeholders) == 0 {
                return v
            } else {
                return fmt.Sprintf(v, placeholders)
            }
            
        } else {
            if len(placeholders) == 0 {
                return text
            } else {
                return fmt.Sprintf(text, placeholders)
            }
        }
    } else {
        if len(placeholders) == 0 {
            return text
        } else {
            return fmt.Sprintf(text, placeholders)
        }
    }
    
}