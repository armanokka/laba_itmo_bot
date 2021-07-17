package main

import "fmt"

func Localize(text, lang string, placeholders ...interface{}) string {
    
    
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
        "/start": {
            "it": "La tua lingua è - %se la lingua per la traduzione è - %s.",
            "de": "Ihre Sprache ist - %s und die Sprache für die Übersetzung ist - %s.",
            "es": "Su idioma es - %s, y el idioma de traducción es - %s.",
            "id": "Bahasa Anda adalah - %s, dan bahasa terjemahannya adalah - %s.",
            "en": "Your language is - %s, and the language for translation is - %s.",
            "uz": "Sizning tilingiz - %s, tarjima qilish uchun - %s.",
            "ru": "Ваш язык - %s, а язык перевода - %s.",
            "uk": "Ваша мова - %s, а мова перекладу - %s.",
            "pt": "Seu idioma é - %s, e o idioma para tradução é - %s.",
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
            "de": "Sprache zum Übersetzen",
            "es": "Idioma para traducir",
            "id": "Bahasa untuk menerjemahkan",
            "it": "Lingua per tradurre",
            "pt": "Língua para tradução",
            "ru": "Язык перевода",
            "uk": "Мова перекладу",
            "uz": "Tarjima qilish uchun til",
        },
        "/my_lang": {
            "pt": "Para configurar  seu idioma , faça  uma  das seguintes opções: 👇\n\nℹ️ Envie  algumas palavras  em seu idioma, por exemplo: \"Hi, how are you today?\" - o idioma será o inglês, ou \"L'amour ne fait pas d'erreurs\" - o idioma será francês e assim por diante.\nℹ️ Ou envie o nome do seu idioma, por ex. \"Russian\", ou \"Japanese\", ou \"Arabic\", e.t.c.",
            "en": "To setup your language, do one of the following: 👇\n\nℹ️ Send few words in your language, for example: \"Hi, how are you today?\" - language will be English, or \"L'amour ne fait pas d'erreurs\" - language will be French, and so on.\nℹ️ Or send the name of your language, e.g. \"Russian\", or \"Japanese\", or  \"Arabic\", e.t.c.",
            "it": "Per impostare la tua lingua, esegui una delle seguenti operazioni: .\n\nℹ️ Invia poche parole nella tua lingua, ad esempio: \"Hi, how are you today?\" - la lingua sarà l'inglese, o \"L'amour ne fait pas d'erreurs\" - la lingua sarà francese, e così via.\nℹ️ Oppure invia il nome della tua lingua, ad es. \"Russo\", o \"Giapponese\", o \"Arabo\", ecc.",
            "de": "Führen Sie eine der folgenden Schritte aus, um Ihre Sprache einzurichten: 👇\n\nℹ️ Sende einige Wörter in deiner Sprache, zum Beispiel: \"Hi, how are you today?\" - Sprache wird Englisch sein, oder \"L'amour ne fait pas d'erreurs\" - Sprache wird französisch sein und so weiter.\nℹ️ Oder schicke den Namen deiner Sprache, z.B. \"Russisch\", oder \"Japanisch\", oder \"Arabisch\", usw.",
            "ru": "Чтобы настроить  свой язык , выполните  одно  из следующих действий: 👇\n\nℹ️ Отправьте  несколько слов  на своем языке, например: \"Hi, how are you today?\" - язык будет английский, или «« L'amour ne fait pas d'erreurs »- язык будет французским и так далее.\nℹ️ Или отправьте название своего языка, например «Русский», «Японский», «Арабский» и т. Д.",
            "es": "Para configurar  su idioma , haga  una  de las siguientes opciones: 👇\n\nℹ️ Envía  algunas palabras  en tu idioma, por ejemplo: \"Hi, how are you today?\" - el idioma será el inglés, o \"L'amour ne fait pas d'erreurs\" - idioma será francés, etc.\nℹ️ O envíe el nombre de su idioma, p. ej. \"Ruso ''\", o \"Japonés\", o \"Árabe\", e.t.c.",
            "uz": " Tilingizni  sozlash uchun quyidagilardan  birini bajaring: 👇\n\nℹ️ O'zingizning tilingizda  bir nechta so'zlarni  yuboring, masalan: \"Hi, how are you today? \"\" \"- til ingliz tilida bo'ladi yoki\" \"L'amour ne fait pas d'erreurs`\" - til frantsuzcha bo'ladi va hokazo.\nℹ️ Yoki o'z tilingiz nomini yuboring, masalan. \"\" \"Ruscha\" yoki \"\" Yaponcha \"yoki\" \"Arabcha\", e.t.c.",
            "uk": "Щоб налаштувати  свою мову , виконайте  одне  з наступного: 👇\n\nℹ️ Надішліть  кілька слів  своєю мовою, наприклад: \"Hi, how are you today?\" - мова буде англійською, або \"L'amour ne fait pas d'erreurs\" - мова буде французькою тощо.\nℹ️ Або надішліть назву вашої мови, напр. \"Російська\", або\"японська\", або  арабська, тощо",
            "id": "Untuk mengatur Bahasa anda, melakukan satu hal:\n\nℹ️ contoh: \"hai, bagaimana kabarmu hari ini?\"- bahasa akan bahasa Inggris, atau \"l'amour ne fait pas d'erreurs\" - bahasa akan Perancis, dan seterusnya.\nℹ️ atau kirim Nama bahasa Anda dalam bahasa Inggris, misal \"rusia\", atau \"jepang\", atau \"arab\", e.t.c.",
        },
        "/to_lang": {
            "uk": "Щоб налаштувати мову перекладу, виконайте одну з таких дій: 👇\n\nℹ️ Надішліть кілька слів мовою на мову, яку потрібно перекласти, наприклад: \"Hi, how are you?\" - мовою буде англійська, або \"L'amour ne fait pas d'erreurs\" - мовою буде французька тощо.\nℹ️ Або надішліть назву мови, напр. \"російська\", або \"японська\", або \"арабська\", напр.",
            "pt": "Para configurar o idioma de tradução, siga um destes procedimentos: 👇\n\nℹ️ Envie algumas palavras no idioma que deseja traduzir, por exemplo: \"Hi, how are you?\" - o idioma será inglês ou \"L'amour ne fait pas d'erreurs\" - o idioma será o francês e assim por diante.\nℹ️ Ou envie o nome do idioma, por ex. \"Russo\" ou \"Japonês\" ou \"Árabe\", e.t.c.",
            "id": "Untuk menerjemahkan bahasa, lakukan satu bahasa berikut:\n\nℹ mengirim beberapa kata dalam bahasa ke dalam Anda ingin menerjemahkan, misalnya: \"hai, apa kabar?\"- bahasa akan bahasa Inggris, atau \"l'amour ne fait pas d'erreurs\" - bahasa akan Perancis, dan seterusnya.\nℹ️atau kirim Nama bahasa, misalnya \"Rusia\", atau \"Jepang\", atau \"Arab\", misalnya. t.c.",
            "en": "To setup translate language, do one of the following: 👇\n\nℹ Send few words in language into you want to translate, for example: \"Hi, how are you?\" - language will be English, or \"L'amour ne fait pas d'erreurs\" - language will be French, and so on.\nℹ️ Or send the name of language, e.g. \"Russian\", or \"Japanese\", or  \"Arabic\", e.t.c.",
            "it": "Per impostare la lingua di traduzione, esegui una delle seguenti operazioni: 👇\n\nℹ️ Invia poche parole nella lingua che vuoi tradurre, ad esempio: \"Hi, how are you?\" - la lingua sarà l'inglese, o \"L'amour ne fait pas d'erreurs\" - la lingua sarà il francese, e così via.\nℹ️ Oppure inviare il nome della lingua, ad es. \"Russo\" o \"Giapponese\" o \"Arabo\", ecc.",
            "uz": "Tarjima tilini sozlash uchun quyidagilardan birini bajaring: 👇\n\nℹ️ O'zingizning tilingizga tarjima qilishni xohlagan bir nechta so'zlarni yuboring, masalan: \"Hi, how are you?\" - til inglizcha bo'ladi yoki \"L'amour ne fait pas d'erreurs\" - frantsuzcha bo'ladi va hokazo.\nℹ️ Yoki til nomini yuboring, masalan. \"Ruscha\", yoki \"Yaponcha\" yoki \"Arabcha\", e.t.c.",
            "de": "Um die Übersetzungssprache einzurichten, führen Sie einen der folgenden Schritte aus:\n\nℹ️ Schicke ein paar Wörter in die Sprache, die du übersetzen möchtest, zum Beispiel: \"Hi, how are you?\" - Sprache ist Englisch oder \"L'amour ne fait pas d'erreurs\" - Sprache ist Französisch und so weiter.\nℹ️ Oder sende den Namen der Sprache, z.B. \"Russisch\" oder \"Japanisch\" oder \"Arabisch\", usw.",
            "ru": "Чтобы настроить язык перевода, выполните одно из следующих действий: 👇\n\nℹ️ Отправьте несколько слов на языке, который вы хотите перевести, например: «Hi, how are you?» - язык будет английский, или \"L'amour ne fait pas d'erreurs\" - язык будет французским, и так далее.\nℹ️ Или отправьте название языка, например «русский», или «японский», или «арабский» и т. Д.",
            "es": "Para configurar el idioma de traducción, realice una de las siguientes acciones: 👇\n\nℹ️ Envíe algunas palabras en el idioma que desea traducir, por ejemplo: \"Hi, how are you?\" - el idioma será el inglés, o \"L'amour ne fait pas d'erreurs\" - el idioma será el francés, etc.\nℹ️ O envíe el nombre del idioma, p. ej. \"Ruso\", \"Japonés\" o \"Árabe\", etc.",
        },
        "💬 Change bot language":{
            "de": "💬 Bot-Sprache ändern",
            "ru": "💬 Изменить язык бота",
            "uk": "💬 Змінити мову бота",
            "pt": "💬 Alterar o idioma do bot",
            "en": "💬 Change bot language",
            "it": "💬 Cambia la lingua del bot",
            "uz": "💬 Bot Bot tilini o'zgartiring",
            "es": "💬 Cambiar el idioma del bot",
            "id": "💬 Ubah bahasa bot",
        },
        "Please, select bot language":{
            "es": "Por favor, seleccione el idioma del bot",
            "en": "Please, select bot language",
            "uk": "Виберіть мову бота",
            "pt": "Por favor, selecione o idioma do bot",
            "id": "Silakan, pilih bahasa bot",
            "it": "Per favore, seleziona la lingua del bot",
            "uz": "Iltimos, bot tilini tanlang",
            "de": "Bitte Bot-Sprache auswählen",
            "ru": "Пожалуйста, выберите язык бота",
        },
        "Now press /start 👈":{
            "it": "Ora premi /start 👈",
            "uz": "Endi /start tugmachasini bosing 👈",
            "id": "Sekarang tekan /start 👈",
            "en": "Now press /start 👈",
            "de": "Drücken Sie nun /start 👈",
            "ru": "Теперь нажмите /start 👈",
            "es": "Ahora presione /start 👈",
            "uk": "Тепер натисніть /start 👈",
            "pt": "Agora pressione /start 👈",
        },
        "⬅Back":{
            "de": "⬅Zurück",
            "es": "⬅Atrás",
            "id": "⬅Kembali",
            "it": "⬅Indietro",
            "pt": "⬅Back",
            "ru": "⬅Назад",
            "uk": "⬅Назад",
            "uz": "⬅Arka",
        },
        "/help": {
            "de": "Was kann dieser Bot tun?\n▫️ Mit Translo können Sie Nachrichten in mehr als 182 Sprachen übersetzen.\nWie übersetzt man eine Nachricht?\n▫️ Zuerst müssen Sie Ihre Sprache einrichten, dann die Übersetzungssprache einrichten, als nächstes Textnachrichten senden und der Bot wird sie schnell übersetzen.\nWie richte ich meine Sprache ein?\n▫️ Klicken Sie unten auf die Schaltfläche \"Meine Sprache\"\nWie richte ich die Sprache ein, in die ich übersetzen möchte?\n▫️ Klicken Sie unten auf die Schaltfläche \"Sprache übersetzen\"\nGibt es noch andere interessante Dinge?\n▫️ Ja, der Bot unterstützt den Inline-Modus. Geben Sie den Spitznamen @translobot in das Nachrichteneingabefeld ein und schreiben Sie dort den Text, den Sie übersetzen möchten.\nIch habe einen Vorschlag oder ich habe einen Fehler gefunden!\n▫️ 👉 Kontaktieren Sie mich bitte - @armanokka",
            "ru": " Что умеет этот бот? \n▫️ Translo позволяет переводить сообщения на 182+ языков.\n Как перевести сообщение? \n▫️ Во-первых, вам нужно настроить свой язык, затем настроить язык перевода, затем отправить текстовые сообщения, и бот быстро их переведет.\n Как настроить мой язык? \n▫️ Нажмите на кнопку под названием «Мой язык».\n Как установить язык, на котором я хочу переводить? \n▫️ Нажмите кнопку ниже под названием «Перевести язык».\n Есть еще что-нибудь интересное? \n▫️ Да, бот поддерживает встроенный режим. Начните вводить псевдоним @translobot в поле ввода сообщения, а затем впишите туда текст, который хотите перевести.\n У меня есть предложение или я обнаружил ошибку! \n▫️ 👉 Свяжитесь со мной, пожалуйста - @armanokka",
            "es": " ¿Qué puede hacer este bot? \n▫️ Translo te permite traducir mensajes a más de 182 idiomas.\n ¿Cómo traducir un mensaje? \n▫️ En primer lugar, debe configurar su idioma, luego configurar el idioma de traducción, luego enviar mensajes de texto y el bot los traducirá rápidamente.\n ¿Cómo configurar mi idioma? \n▫️ Haga clic en el botón de abajo llamado \"Mi idioma\"\n ¿Cómo configurar el idioma al que quiero traducir? \n▫️ Haga clic en el botón de abajo llamado \"Traducir idioma\"\n ¿Hay otras cosas interesantes? \n▫️ Sí, el bot admite el modo en línea. Comience a escribir el apodo @translobot en el campo de entrada del mensaje y luego escriba allí el texto que desea traducir.\n ¡Tengo una sugerencia o encontré un error! \n▫️ 👉 Contáctame por favor - @armanokka",
            "uk": " Що може зробити цей бот? \n▫️ Translo дозволяє перекладати повідомлення 182+ мовами.\n Як перекласти повідомлення? \n▫️ По-перше, вам потрібно налаштувати свою мову, потім налаштувати переклад мови, наступне надіслати текстові повідомлення і бот швидко їх перекладе.\n Як налаштувати мову? \n▫️ Клацніть на кнопку під назвою \"Моя мова\"\n Як налаштувати мову на мову, яку я хочу перекласти? \nКлацніть на кнопку нижче під назвою \"Перекласти мову\"\n Чи є якісь інші цікаві речі? \nТак, бот підтримує вбудований режим. Почніть вводити псевдонім @translobot у поле введення повідомлення, а потім напишіть туди текст, який потрібно перекласти.\n У мене є пропозиція або я знайшов помилку! \n👉 Зв’яжіться зі мною pls - @armanokka",
            "pt": " O que este bot pode fazer? \n▫️ Translo permite que você traduza mensagens em mais de 182 idiomas.\n Como traduzir a mensagem? \n▫️ Em primeiro lugar, você tem que configurar seu idioma, depois configurar traduzir o idioma, em seguida enviar mensagens de texto e o bot irá traduzi-las rapidamente.\n Como configurar meu idioma? \n▫️ Clique no botão abaixo chamado \"Meu Idioma\"\n Como configurar o idioma para o que desejo traduzir? \n▫️ Clique no botão abaixo chamado \"Traduzir Idioma\"\n Existem outras coisas interessantes? \n▫️ Sim, o bot suporta o modo inline. Comece digitando o apelido @translobot no campo de entrada da mensagem e então escreva lá o texto que deseja traduzir.\n Tenho uma sugestão ou encontrei um bug! \n▫️ 👉 Contate-me, pls - @armanokka",
            "id": "Apa yang bisa dilakukan robot ini?\n▫️ ️ooo memungkinkan Anda untuk menerjemahkan pesan ke 182+ bahasa.\nBagaimana menerjemahkan pesan?\n▫️  Anda harus menata bahasa Anda, kemudian menerjemahkan bahasa, teks kirim pesan teks dan keduanya akan menerjemahkannya dengan cepat.\n▫️ Bagaimana cara mengatur bahasaku?\n▫️ ️ Klik di tombol di bawah ini disebut \" Bahasa saya\"\nBagaimana mengatur bahasa ke dalam saya ingin menerjemahkan?\nKlik di tombol di bawah ini yang disebut \"menerjemahkan Bahasa\"\n▫️ ️ Apakah ada hal-hal yang menarik lainnya?\nYa, dukungan robot dalam mode inline. Mulai mengetik Nama panggilan @translobot dalam kolom masukan pesan dan kemudian menulis di sana teks yang ingin Anda terjemahkan.\n▫️ ️ Aku punya saran atau aku menemukan bug!\nHubungi aku pls - @armanokka",
            "en": "What can this bot do?\n▫️ Translo allow you to translate messages into 182+ languages.\nHow to translate message?\n▫️ Firstly, you have to setup your language, then setup translate language, next send text messages and bot will translate them quickly.\nHow to setup my language?\n▫️ Click on the button below called \"My Language\"\nHow to setup language into I want to translate?\n▫️ Click on the button below called \"Translate Language\"\nAre there any other interesting things?\n▫️ Yes, the bot support inline mode. Start typing the nickname @translobot in the message input field and then write there the text you want to translate.\nI have a suggestion or I found bug!\n▫️ 👉 Contact me pls - @armanokka",
            "it": "Cosa può fare questo bot?\n▫️ Translo ti consente di tradurre i messaggi in oltre 182 lingue.\nCome tradurre il messaggio?\n▫️ In primo luogo, devi impostare la tua lingua, quindi impostare la lingua di traduzione, quindi inviare messaggi di testo e il bot li tradurrà rapidamente.\nCome impostare la mia lingua?\n▫️ Clicca sul pulsante in basso chiamato \"La mia lingua\"\nCome impostare la lingua in cui voglio tradurre?\n▫️ Clicca sul pulsante in basso chiamato \"Traduci lingua\"\nCi sono altre cose interessanti?\n▫️ Sì, il bot supporta la modalità in linea. Inizia a digitare il nickname @translobot nel campo di inserimento del messaggio e poi scrivi lì il testo che vuoi tradurre.\nHo un suggerimento o ho trovato un bug!\n▫️ 👉 Contattami per favore - @armanokka",
            "uz": " Bu bot nima qila oladi? \n▫️ Translo sizga xabarlarni 182+ tilga tarjima qilishga imkon beradi.\n Xabarni qanday tarjima qilish kerak? \n▫️ Birinchidan, siz o'z tilingizni o'rnatishingiz kerak, so'ngra tarjima tilini sozlashingiz kerak, keyin matnli xabarlarni yuboring va bot ularni tezda tarjima qiladi.\n Mening tilimni qanday o'rnatish kerak? \n▫️ Quyidagi \"Mening tilim\" deb nomlangan tugmani bosing.\n Men tarjima qilishni xohlagan tilni qanday o'rnataman? \n▫️ \"Tilni tarjima qilish\" deb nomlangan tugmani bosing.\n Boshqa qiziqarli narsalar bormi? \n▫️ Ha, botni inline rejimida qo'llab-quvvatlash. Xabarlarni kiritish maydoniga @translobot taxallusini yozishni boshlang va keyin tarjima qilmoqchi bo'lgan matni yozing.\n Menda taklif bor yoki men xato topdim! \n▫️ 👉 Men bilan bog'laning pls - @armanokka",
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
        "To voice":{
            "it": "esprimere",
            "ru": "Озвучить",
            "es": "a voz",
            "pt": "dar voz",
            "id": "untuk menyuarakan",
            "uz": "ovoz berish",
            "de": "aussprechen",
            "uk": "озвучити",
        },
        "Ok, close": {
            "es": "Ok cerrar",
            "pt": "Ok fechar",
            "id": "Oke, tutup",
            "uk": "Гаразд, близько",
            "it": "Ok, chiudi",
            "uz": "Yaxshi, yaqin",
            "de": "Okay, nah",
            "ru": "Хорошо, закрыть",
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

    if df, ok := languages[text]; ok { // Текст подходит под варианты
        if v, ok := df[lang]; ok { // Есть соответствующий язык
            return fmt.Sprintf(v, placeholders...)
        }
    }
   return fmt.Sprintf(text, placeholders...)
}