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
        "Now your language is %s": {
            {
                LanguageCode: "de",
                Text:         "Ihre Sprache ist jetzt %s",
            },
            {
                LanguageCode: "es",
                Text:         "Ahora tu idioma es %s",
            },
            {
                LanguageCode: "id",
                Text:         "Sekarang bahasa Anda adalah %s",
            },
            {
                LanguageCode: "it",
                Text:         "Ora la tua lingua è %s",
            },
            {
                LanguageCode: "pt",
                Text:         "Agora seu idioma é %s",
            },
            {
                LanguageCode: "ru",
                Text:         "Теперь ваш язык %s",
            },
            {
                LanguageCode: "uk",
                Text:         "Тепер ваша мова - %s",
            },
            {
                LanguageCode: "uz",
                Text:         "Endi sizning tilingiz - %s",
            },
            {
                LanguageCode: "en",
                Text:         "Now your language is %s",
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
        "Now translate language is %s": {
            {
                LanguageCode: "de",
                Text:         "Die Übersetzungssprache ist jetzt %s",
            },
            {
                LanguageCode: "es",
                Text:         "Ahora el idioma de traducción es %s",
            },
            {
                LanguageCode: "id",
                Text:         "Sekarang bahasa terjemahan adalah %s",
            },
            {
                LanguageCode: "it",
                Text:         "Ora la lingua di traduzione è %s",
            },
            {
                LanguageCode: "pt",
                Text:         "Agora o idioma de tradução é %s",
            },
            {
                LanguageCode: "ru",
                Text:         "Теперь язык перевода %s",
            },
            {
                LanguageCode: "uk",
                Text:         "Тепер мовою перекладу є %s",
            },
            {
                LanguageCode: "uz",
                Text:         "Endi tarjima tili - %s",
            },
            {
                LanguageCode: "en",
                Text:         "Now translate language is %s",
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
                Text:         "Werbung in @TransloBot\n\nAnzeigen werden in jeder vom Bot übersetzten Nachricht platziert\n\n• Sie zahlen nur für den Vermögenswert\n• Viel billiger als der Versand\n• Inline-Modus ist in den Bedingungen enthalten\n• Wenn ein Benutzer den Inline-Modus in einer Gruppe verwendet, wird Ihre Anzeige von allen Chat-Teilnehmern gesehen, d. h. von Zehntausenden\n\nGeben Sie Ihren Nachrichtentext ein,\n- 130 Zeichen inklusive Leerzeichen\n- ohne Fotos und Gifs\n- keine Wörter verwenden: abonnieren, reinkommen, gehen usw.\n- das Markup funktioniert nicht, d.h. fetter, kursiver Text wird wie gewohnt angezeigt",
            },
            {
                LanguageCode: "es",
                Text:         "Publicidad en @TransloBot\n\nLos anuncios se colocan en cada mensaje traducido por el bot.\n\n• Paga solo por el activo\n• Mucho más barato que enviar por correo\n• El modo en línea está incluido en las condiciones\n• Si un usuario usa el modo en línea en un grupo, todos los participantes del chat ven su anuncio, es decir, decenas de miles.\n\nIngrese el texto de su mensaje,\n- 130 caracteres incluyendo espacios\n- sin fotos y gifs\n- no uses palabras: suscríbete, entra, ve, etc.\n- el marcado no funciona, es decir, el texto en negrita y cursiva será normal",
            },
            {
                LanguageCode: "id",
                Text:         "Beriklan di @TransloBot\n\nIklan ditempatkan di setiap pesan yang diterjemahkan oleh bot\n\n• Anda hanya membayar untuk aset\n• Jauh lebih murah daripada mengirim surat\n• Mode sebaris termasuk dalam ketentuan\n• Jika pengguna menggunakan mode sebaris dalam grup, maka iklan Anda dilihat oleh semua peserta obrolan, yaitu puluhan ribu\n\nMasukkan teks pesan Anda,\n- 130 karakter termasuk spasi\n- tanpa foto dan gif\n- jangan menggunakan kata-kata: berlangganan, masuk, pergi, dll.\n- markup tidak berfungsi, yaitu huruf tebal, teks miring akan seperti biasa",
            },
            {
                LanguageCode: "it",
                Text:         "Pubblicità in @TransloBot\n\nGli annunci vengono inseriti in ogni messaggio tradotto dal bot\n\n• Paghi solo per il bene\n• Molto più economico della spedizione\n• La modalità in linea è inclusa nelle condizioni\n• Se un utente utilizza la modalità in linea in un gruppo, il tuo annuncio viene visualizzato da tutti i partecipanti alla chat, ovvero decine di migliaia\n\nInserisci il testo del tuo messaggio,\n- 130 caratteri spazi inclusi\n- senza foto e gif\n- non usare parole: iscriviti, entra, vai, ecc.\n- il markup non funziona, ovvero il testo in grassetto e corsivo sarà normale",
            },
            {
                LanguageCode: "pt",
                Text:         "Publicidade em @TransloBot\n\nOs anúncios são colocados em cada mensagem traduzida pelo bot\n\n• Você paga apenas pelo ativo\n• Muito mais barato do que enviar pelo correio\n• O modo inline está incluído nas condições\n• Se um usuário usa o modo inline em um grupo, seu anúncio é visto por todos os participantes do bate-papo, ou seja, dezenas de milhares\n\nDigite o texto da sua mensagem,\n- 130 caracteres incluindo espaços\n- sem fotos e gifs\n- não use palavras: inscreva-se, entre, vá, etc.\n- a marcação não funciona, ou seja, o texto em negrito e itálico ficará normal",
            },
            {
                LanguageCode: "ru",
                Text:         "Реклама в @TransloBot\n\nРеклама размещается в каждом переведенном ботом сообщении \n\n• Вы платите только за актив\n• Значительно дешевле рассылки\n• Инлайн-режим включен в условия\n• Если пользователь пользуется инлайн-режимом в группе, то вашу рекламу видят все участники чата, то есть десятки тысяч\n\nВведите текст вашего сообщения, \n- 130 символов с учетом пробелов\n- без фото и гиф\n- не использовать слова : подписывайся, заходи, переходи и т.д.\n- разметка не работает, т.е жирный, курсивный текст будет как обычный",
            },
            {
                LanguageCode: "uk",
                Text:         "Реклама в @TransloBot\n\nРеклама розміщується в кожному перекладеному ботом повідомленні\n\n• Ви платите тільки за актив\n• Значно дешевше розсилки\n• Інлайн-режим включений в умови\n• Якщо користувач користується інлайн-режимом в групі, то вашу рекламу бачать всі учасники чату, тобто десятки тисяч\n\nВведіть текст вашого повідомлення,\n- 130 символів, враховуючи пробіли\n- без фото і гіф\n- не використовувати слова: підписуйся, заходь, переходь і т.д.\n- розмітка не працює, тобто жирний, курсивний текст буде як звичайний",
            },
            {
                LanguageCode: "uz",
                Text:         "@TransloBot -da reklama\n\nBot tomonidan tarjima qilingan har bir xabarda reklama joylashtiriladi\n\n• Siz faqat aktiv uchun to'laysiz\n• Pochta yuborishdan ancha arzon\n• Inline rejimi shartlarga kiritilgan\n• Agar foydalanuvchi guruhda inline rejimidan foydalansa, u holda sizning reklamangizni chatning barcha ishtirokchilari ko'radi, ya'ni o'n minglab\n\nXabar matnini kiriting,\n- bo'shliqlarni o'z ichiga olgan 130 ta belgi\n- fotosuratlar va sovg'alarsiz\n- so'zlarni ishlatmang: obuna bo'ling, kiring, kiring va hokazo.\n- belgilash ishlamaydi, ya'ni qalin, kursiv matn odatdagidek bo'ladi",
            },
            {
                LanguageCode: "en",
                Text:         "Advertising in @TransloBot\n\nAds are placed in every message translated by the bot\n\n• You pay only for the asset\n• Much cheaper than mailing\n• Inline mode is included in the conditions\n• If a user uses inline mode in a group, then your ad is seen by all chat participants, that is, tens of thousands\n\nEnter your message text,\n- 130 characters including spaces\n- without photos and gifs\n- do not use words: subscribe, come in, go, etc.\n- the markup does not work, i.e. bold, italic text will be as normal",
            },
        },
        "sponsorship_set_days": {
            {
                LanguageCode: "de",
                Text:         "Wählen Sie nun die Anzahl der Tage aus, die Sie inserieren möchten (die Tage beginnen mit dem Zahlungsdatum + Verrechnung mit den bereits belegten Tagen).",
            },
            {
                LanguageCode: "es",
                Text:         "Ahora seleccione la cantidad de días que desea anunciar (los días comenzarán a partir de la fecha de pago + compensación de los días ya ocupados).",
            },
            {
                LanguageCode: "id",
                Text:         "Sekarang pilih jumlah hari yang ingin Anda iklankan (hari akan dimulai dari tanggal pembayaran + offset dari hari yang sudah ditempati).",
            },
            {
                LanguageCode: "it",
                Text:         "Ora seleziona il numero di giorni che vuoi pubblicizzare (i giorni inizieranno dalla data di pagamento + offset dai giorni già occupati).",
            },
            {
                LanguageCode: "pt",
                Text:         "Agora selecione o número de dias que deseja anunciar (os dias começarão a partir da data do pagamento + compensação dos dias já ocupados).",
            },
            {
                LanguageCode: "ru",
                Text:         "Теперь выберите количество дней, которые хотите рекламироваться (отсчёт дней начнётся со дня оплаты + смещение от уже занятых дней).",
            },
            {
                LanguageCode: "uk",
                Text:         "Тепер виберіть кількість днів, які хочете рекламуватися (відлік днів почнеться з дня оплати + зміщення від уже зайнятих днів).",
            },
            {
                LanguageCode: "uz",
                Text:         "Endi siz reklama qilmoqchi bo'lgan kunlar sonini tanlang (kunlar to'langan kundan boshlanadi + ishg'ol qilingan kunlardan ofset).",
            },
            {
                LanguageCode: "en",
                Text:         "Now select the number of days you want to advertise (the days will start from the date of payment + offset from the days already occupied).",
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
        "Выберите язык, на котором хотите переводить текст": {
            {
                LanguageCode: "de",
                Text:         "Wählen Sie die Sprache aus, in die Sie den Text übersetzen möchten",
            },
            {
                LanguageCode: "es",
                Text:         "Seleccione el idioma en el que desea traducir el texto.",
            },
            {
                LanguageCode: "id",
                Text:         "Pilih bahasa di mana Anda ingin menerjemahkan teks",
            },
            {
                LanguageCode: "it",
                Text:         "Seleziona la lingua in cui vuoi tradurre il testo",
            },
            {
                LanguageCode: "pt",
                Text:         "Selecione o idioma no qual deseja traduzir o texto",
            },
            {
                LanguageCode: "ru",
                Text:         "Выберите язык, на котором хотите переводить текст",
            },
            {
                LanguageCode: "uk",
                Text:         "Виберіть мову, на якому хочете перекладати текст",
            },
            {
                LanguageCode: "uz",
                Text:         "Matnni tarjima qilmoqchi bo'lgan tilni tanlang",
            },
            {
                LanguageCode: "en",
                Text:         "Select the language in which you want to translate the text",
            },
        },
        "Выберите ваш родной язык": {
           {
                LanguageCode: "de",
                Text:         "Wählen Sie Ihre Muttersprache",
            },
           {
                LanguageCode: "es",
                Text:         "Elija su lengua materna",
            },
           {
                LanguageCode: "id",
                Text:         "Pilih bahasa ibu Anda",
            },
           {
                LanguageCode: "it",
                Text:         "Scegli la tua lingua madre",
            },
           {
                LanguageCode: "pt",
                Text:         "Escolha sua língua nativa",
            },
           {
                LanguageCode: "ru",
                Text:         "Выберите ваш родной язык",
            },
           {
                LanguageCode: "uk",
                Text:         "Виберіть ваш рідну мову",
            },
           {
                LanguageCode: "uz",
                Text:         "Ona tilingizni tanlang",
            },
           {
                LanguageCode: "en",
                Text:         "Choose your native language",
            },
        },
        "bot_advertise": {
            {
                LanguageCode: "de",
                Text:         "Wir sind Ihnen dankbar, wenn Sie Ihren Freunden von uns erzählen. Leiten Sie diese Nachricht einfach an sie weiter.\n\n😎 Übersetzen Sie Nachrichten schnell und einfach, ohne Telegram zu verlassen.\n@translobot",
            },
            {
                LanguageCode: "es",
                Text:         "Estaremos agradecidos si les cuenta a sus amigos sobre nosotros. Simplemente reenvíeles este mensaje.\n\n😎 Traduce mensajes de forma rápida y sencilla sin salir de Telegram.\n@translobot",
            },
            {
                LanguageCode: "id",
                Text:         "Kami akan berterima kasih jika Anda memberi tahu teman Anda tentang kami. Teruskan saja pesan ini kepada mereka.\n\n😎 Terjemahkan pesan dengan cepat dan mudah tanpa meninggalkan Telegram.\n@translobot",
            },
            {
                LanguageCode: "it",
                Text:         "Ti saremo grati se parlerai di noi ai tuoi amici. Basta inoltrare loro questo messaggio.\n\n😎 Traduci i messaggi in modo rapido e semplice senza uscire da Telegram.\n@translobot",
            },
            {
                LanguageCode: "pt",
                Text:         "Ficaremos gratos se você contar a seus amigos sobre nós. Apenas encaminhe esta mensagem para eles.\n\n😎 Traduza mensagens de forma rápida e fácil sem sair do Telegram.\n@translobot",
            },
            {
                LanguageCode: "ru",
                Text:         "Мы будем благодарны, если вы расскажете о нас друзьям. Просто перешлите им это сообщение.\n\n😎 Переводите сообщения быстро и легко, не выходя из Telegram.\n@translobot",
            },
            {
                LanguageCode: "uk",
                Text:         "Ми будемо вдячні, якщо ви розповісте про нас друзям. Просто перешліть їм це повідомлення.\n\n😎 Переводите повідомлення швидко і легко, не виходячи з Telegram.\n@translobot",
            },
            {
                LanguageCode: "uz",
                Text:         "Agar do'stlaringizga biz haqimizda aytib bersangiz, minnatdor bo'lamiz. Bu xabarni ularga yuboring.\n\n😎 Telegramdan chiqmasdan xabarlarni tez va oson tarjima qiling.\n@translobot",
            },
            {
                LanguageCode: "en",
                Text:         "We will be grateful if you tell your friends about us. Just forward this message to them.\n\n😎 Translate messages quickly and easily without leaving Telegram.\n@translobot",
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