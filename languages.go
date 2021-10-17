package main

import (
    "fmt"
    "strings"
)

var commandTranslations = map[string][]string{
    "my language": {
        "meine Sprache",
        "mi idioma",
        "моя мова",
        "bahasaku",
        "my language",
        "mio linguaggio",
        "mening tilim",
        "мой язык",
        "minha língua",
    },
    "translate language": {
        "lengua de llegada",
        "translate language",
        "язык перевода",
        "Zielsprache",
        "мова перекладу",
        "idioma de chegada",
        "bahasa sasaran",
        "lingua di destinazione",
        "maqsadli til",
    },
}

var translations = map[string]map[string]string{
    "Сейчас бот переводит на %s. Выберите язык для перевода": map[string]string{
        "de": "Der Bot übersetzt derzeit in %s. Wählen Sie die Sprache für die Übersetzung aus",
        "es": "El bot se está traduciendo actualmente a %s. Seleccione el idioma para traducir",
        "ru": "Сейчас бот переводит на %s. Выберите язык для перевода",
        "uk": "Зараз бот переводить на %s. Виберіть мову для перекладу",
        "uz": "Bot hozirda %s tiliga tarjima qilinmoqda. Tarjima qilish uchun tilni tanlang",
        "en": "The bot is currently translating to %s. Select the language for translation",
        "id": "Bot saat ini menerjemahkan ke %s. Pilih bahasa untuk terjemahan",
        "it": "Il bot sta attualmente traducendo in %s. Seleziona la lingua per la traduzione",
        "pt": "O bot está traduzindo para %s. Selecione o idioma para tradução",
    },
    "Ваш язык %s. Выберите Ваш язык.": map[string]string{
        "it": "La tua lingua è %s. Scegli la tua LINGUA.",
        "uz": "Sizning tilingiz - %s. Tilingizni tanlang.",
        "en": "Your language is %s. Choose your language.",
        "de": "Ihre Sprache ist %s. Wähle deine Sprache.",
        "es": "Tu idioma es %s. Elige tu idioma.",
        "id": "Bahasa Anda adalah %s. Pilih bahasamu.",
        "pt": "Seu idioma é %s. Escolha seu idioma.",
        "ru": "Ваш язык %s. Выберите Ваш язык.",
        "uk": "Мову %s. Виберіть мову.",
    },
    "Bot language": {
        "it": "Linguaggio Bot",
        "de": "Bot-Sprache",
        "es": "Lenguaje bot",
        "pt": "Linguagem de bot",
        "id": "Bahasa bot",
        "uk": "Бот-мова",
        "uz": "Bot tili",
        "ru": "Язык бота",
    },
    "Too big text": {
        "id": "Teks terlalu besar",
        "it": "Testo troppo grande",
        "uz": "Juda katta matn",
        "de": "Zu großer Text",
        "ru": "Слишком большой текст",
        "es": "Texto demasiado grande",
        "uk": "Занадто великий текст",
        "pt": "Texto muito grande",
        "en": "Too big text",
    },
    "Available only for idioms, nouns, verbs and adjectives": {
        "es": "Disponible solo para modismos, sustantivos, verbos y adjetivos",
        "pt": "Disponível apenas para expressões idiomáticas, substantivos, verbos e adjetivos",
        "it": "Disponibile solo per modi di dire, nomi, verbi e aggettivi",
        "uz": "Faqat iboralar, otlar, fe'llar va sifatlar uchun mavjud",
        "de": "Nur für Redewendungen, Substantive, Verben und Adjektive verfügbar",
        "ru": "Доступно только для идиом, существительных, глаголов и прилагательных.",
        "en": "Available only for idioms, nouns, verbs and adjectives",
        "uk": "Доступно лише для ідіом, іменників, дієслів та прикметників",
        "id": "Hanya tersedia untuk idiom, kata benda, kata kerja dan kata sifat",
    },
    "Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)": {
        "es": "Lo sentimos, error causado.\n\nPor favor, no bloquees el bot, corregiré el error en un futuro cercano, el administrador ya ha sido advertido sobre este error ;)",
        "uk": "Вибачте, сталася помилка.\n\nБудь ласка, не блокуйте бота, я виправлю помилку найближчим часом, адміністратор вже попереджений про цю помилку ;)",
        "pt": "Desculpa, erro causado.\n\nPor favor, não bloqueie o bot, eu vou corrigir o bug em um futuro próximo, o administrador já foi avisado sobre este ",
        "id": "Maaf, kesalahan disebabkan.\n\nTolong, jangan halangi robot itu, aku akan memperbaiki bug dalam waktu dekat, administrator telah diperingatkan tentang kesalahan ini ;)",
        "it": "Scusa, errore causato.\n\nPer favore, non bloccare il bot, correggerò il bug nel prossimo futuro, l'amministratore è già stato avvertito di questo errore;)",
        "uz": "Uzr, xato sabab.\n\nIltimos, botni bloklamang, yaqin kelajakda xatoni tuzataman, administrator bu xato haqida allaqachon ogohlantirilgan ;)",
        "de": "Sorry, Fehler verursacht.\n\nBitte blockieren Sie den Bot nicht, ich werde den Fehler in naher Zukunft beheben, der Administrator wurde bereits vor diesem Fehler gewarnt ;)",
        "ru": "Извините, произошла ошибка.\n\nПожалуйста, не блокируйте бота, я исправлю ошибку в ближайшее время, администратор уже предупрежден об этой ошибке ;)",
    },
    "Delete": {
        "id": "Menghapus",
        "it": "Elimina",
        "ru": "Удалить",
        "uz": "O'chirish",
        "pt": "Excluir",
        "es": "Borrar",
        "en": "Delete",
        "de": "Löschen",
        "uk": "Видалити",
    },
    "Now your language is %s": {
        "es": "Ahora tu idioma es %s",
        "it": "Ora la tua lingua è %s",
        "uk": "Тепер ваша мова - %s",
        "ru": "Теперь ваш язык %s",
        "uz": "Endi sizning tilingiz - %s",
        "en": "Now your language is %s",
        "de": "Ihre Sprache ist jetzt %s",
        "id": "Sekarang bahasa Anda adalah %s",
        "pt": "Agora seu idioma é %s",
    },
    "Другие языки": {
        "uz": "Boshqa tillar",
        "en": "Other languages",
        "ru": "Другие языки",
        "es": "Otros idiomas",
        "id": "Bahasa lainnya",
        "it": "Altre lingue",
        "pt": "Outras línguas",
        "uk": "Інші мови",
        "de": "Andere Sprachen",
    },
    "⏳ Translating...": {
        "ru": "⏳ Перевожу...",
        "es": "⏳ Traducir...",
        "uk": "⏳ Перекласти...",
        "pt": "⏳ Traduzir...",
        "id": "⏳ Terjemahkan...",
        "it": "⏳ Traduci...",
        "uz": "⏳ Tarjima qilish...",
        "de": "⏳ Übersetzen...",
    },
    "Variants": {
        "it": "varianti",
        "pt": "Variantes",
        "ru": "Варианты",
        "uk": "Варіанти",
        "uz": "Variantlar",
        "de": "Varianten",
        "es": "Variantes",
        "id": "Varian",
    },
    "/start": {
        "ru": "Ваш язык - %s, а язык перевода - %s.",
        "uk": "Ваша мова - %s, а мова перекладу - %s.",
        "it": "La tua lingua è - %se la lingua per la traduzione è - %s.",
        "es": "Su idioma es - %s, y el idioma de traducción es - %s.",
        "pt": "Seu idioma é - %s, e o idioma para tradução é - %s.",
        "id": "Bahasa Anda adalah - %s, dan bahasa terjemahannya adalah - %s.",
        "en": "Your language is - %s, and the language for translation is - %s.",
        "uz": "Sizning tilingiz - %s, tarjima qilish uchun - %s.",
        "de": "Ihre Sprache ist - %s und die Sprache für die Übersetzung ist - %s.",
    },
    "Unsupported language or internal error": {
        "pt": "Idioma não suportado ou erro interno",
        "id": "Bahasa yang tidak didukung atau kesalahan internal",
        "it": "Lingua non supportata o errore interno",
        "uz": "Qo'llab-quvvatlanmaydigan til yoki ichki xato",
        "de": "Nicht unterstützte Sprache oder interner Fehler",
        "ru": "Неподдерживаемый язык или внутренняя ошибка",
        "es": "Idioma no admitido o error interno",
        "uk": "Непідтримувана мова або внутрішня помилка",
    },
    "Please, send text message": {
        "it": "Per favore, invia un messaggio di testo",
        "uz": "Iltimos, SMS yuboring",
        "de": "Bitte SMS senden",
        "ru": "Пожалуйста, отправьте текстовое сообщение",
        "es": "Por favor, envíe un mensaje de texto",
        "uk": "Будь ласка, надішліть текстове повідомлення",
        "pt": "Por favor, envie mensagem de texto",
        "id": "Tolong, kirim pesan teks",
    },
    "bot_advertise":{
        "es": "Estaremos agradecidos si les cuenta a sus amigos sobre nosotros. Reenvíeles este mensaje. <b>Translo Bot</b> Traductor en línea rápido y conveniente <i>Admite 182</i> <b>(!)</b> <i>Idiomas</i> @translobot",
        "uk": "We will be grateful if you tell your friends about us. Forward them this message. <b>Translo Bot</b> Fast and convenient inline translator <i>Supports 182</i> <b>(!)</b> <i>Languages</i> @translobot",
        "pt": "Ficaremos gratos se você contar a seus amigos sobre nós. Encaminhe esta mensagem para eles. <b>Translo Bot Conversor</b> embutido rápido e conveniente <i>Suporta 182</i> <b>(!)</b> <i>Idiomas</i> @translobot",
        "en": "We will be grateful if you tell your friends about us. Forward them this message. <b>Translo Bot</b> Fast and convenient inline translator <i>Supports 182</i> <b>(!)</b> <i>Languages</i> @translobot",
        "it": "Ti saremo grati se parlerai di noi ai tuoi amici. Inoltra loro questo messaggio. <b>Translo Bot</b> Traduttore in linea veloce e conveniente <i>Supporta 182</i> <b>(!)</b> <i>Lingue</i> @translobot",
        "uz": "Agar do&#39;stlaringizga biz haqimizda aytib bersangiz, minnatdor bo&#39;lamiz. Bu xabarni ularga yuboring. <b>Translo Bot</b> Tez va qulay inline tarjimon <i>182</i> <b>(!)</b> <i>Tilni</i> <i>qo&#39;llab -quvvatlaydi</i> @translobot",
        "de": "Wir sind Ihnen dankbar, wenn Sie Ihren Freunden von uns erzählen. Leiten Sie diese Nachricht weiter. <b>Translo Bot</b> Schneller und bequemer Inline-Übersetzer <i>Unterstützt 182</i> <b>(!)</b> <i>Sprachen</i> @translobot",
        "ru": "We will be grateful if you tell your friends about us. Forward them this message.\n\n<b>Translo Bot</b> Fast and convenient inline translator\n<i>Supports 182</i> <b>(!)</b> <i>languages</i>\n\n@translobot",
        "id": "Kami akan berterima kasih jika Anda memberi tahu teman Anda tentang kami. Teruskan pesan ini kepada mereka. <b>Translo Bot</b> Penerjemah sebaris yang cepat dan nyaman <i>Mendukung 182</i> <b>(!)</b> <i>Bahasa</i> @translobot",
    },
    "bot_advertise_without_html_tags": {
        "en": "Translo Bot Fast and convenient inline translator Supports 182 (!) Languages \n\n@translobot",
        "uk": "Translo Bot Швидкий і зручний інлайн перекладач Підтримує 182 (!) Мови @translobot",
        "id": "Translo Bot Penerjemah sebaris yang cepat dan nyaman Mendukung 182 (!) Bahasa @translobot",
        "pt": "Translo Bot Conversor embutido rápido e conveniente Suporta 182 (!) Idiomas @translobot",
        "it": "Translo Bot Traduttore in linea veloce e conveniente Supporta 182 (!) Lingue @translobot",
        "uz": "Translo Bot Tez va qulay inline tarjimon 182 (!) Tilni qo&#39;llab -quvvatlaydi @translobot",
        "de": "Translo Bot Schneller und bequemer Inline-Übersetzer Unterstützt 182 (!) Sprachen @translobot",
        "ru": "Translo Bot\nБыстрый и удобный инлайн переводчик \nПоддерживает 182(!) языка\n\n@translobot",
        "es": "Translo Bot Traductor en línea rápido y conveniente Admite 182 (!) Idiomas @translobot",
    },
    "Share": {
        "de": "Teilen",
        "ru": "Поделиться",
        "es": "Cuota",
        "id": "Membagikan",
        "en": "Share",
        "it": "Condividere",
        "uz": "Ulashing",
        "uk": "Поділитися",
        "pt": "Compartilhado",
    },
    "My Language": {
        "it": "La mia lingua",
        "uz": "Tilimni",
        "de": "Meine Sprache",
        "ru": "Мой Язык",
        "es": "Mi Idioma",
        "uk": "Моя Мова",
        "pt": "A Minha Língua",
        "en": "My Language",
        "id": "Bahasa Saya",
    },
    "Now translate language is %s": {
        "en": "Now translate language is %s",
        "de": "Die Übersetzungssprache ist jetzt %s",
        "id": "Sekarang bahasa terjemahan adalah %s",
        "it": "Ora la lingua di traduzione è %s",
        "pt": "Agora o idioma de tradução é %s",
        "ru": "Теперь язык перевода %s",
        "uz": "Endi tarjima tili - %s",
        "es": "Ahora el idioma de traducción es %s",
        "uk": "Тепер мовою перекладу є %s",
    },
    "Please, select bot language": {
        "it": "Per favore, seleziona la lingua del bot",
        "ru": "Пожалуйста, выберите язык бота",
        "de": "Bitte Bot-Sprache auswählen",
        "en": "Please, select bot language",
        "uz": "Iltimos, bot tilini tanlang",
        "uk": "Виберіть мову бота",
        "es": "Por favor, seleccione el idioma del bot",
        "pt": "Por favor, selecione o idioma do bot",
        "id": "Silakan, pilih bahasa bot",
    },
    "Error, sorry.": {
        "pt": "Erro, desculpe.",
        "id": "Kesalahan, maaf.",
        "it": "Errore, scusa.",
        "uz": "Xato, uzr.",
        "de": "Fehler, tut mir leid.",
        "ru": "Ошибка, извините.",
        "es": "Error, lo siento.",
        "uk": "Помилка, вибачте.",
    },
    "Ok, close": {
        "uz": "Yaxshi, yaqin",
        "de": "Okay, nah",
        "ru": "Хорошо, закрыть",
        "es": "Ok cerrar",
        "pt": "Ok fechar",
        "id": "Oke, tutup",
        "uk": "Гаразд, близько",
        "it": "Ok, chiudi",
    },
    "Выберите страны пользователей, которые получат вашу рассылку. Когда закончите, нажмите кнопку Далее": {
        "id": "Pilih negara pengguna yang akan menerima buletin Anda. Setelah selesai, klik Berikutnya",
        "it": "Seleziona i paesi degli utenti che riceveranno la tua newsletter. Al termine, fai clic su Avanti",
        "pt": "Selecione os países dos usuários que receberão sua newsletter. Quando terminar, clique em Avançar",
        "ru": "Выберите страны пользователей, которые получат вашу рассылку. Когда закончите, нажмите кнопку Далее",
        "de": "Wählen Sie die Länder der Benutzer aus, die Ihren Newsletter erhalten. Wenn Sie fertig sind, klicken Sie auf Weiter",
        "es": "Seleccione los países de los usuarios que recibirán su newsletter. Cuando termine, haga clic en Siguiente",
        "uk": "Виберіть країни користувачів, які отримають вашу розсилку. Коли закінчите, натисніть кнопку Далі",
        "uz": "Sizning axborot byulleteningizni oladigan foydalanuvchilarning mamlakatlarini tanlang. Ish tugagach, Keyingiga bosing",
        "en": "Select the countries of the users who will receive your newsletter. When done, click Next",
    },
    "Empty result": {
        "pt": "Resultado vazio",
        "id": "Hasil kosong",
        "it": "Risultato vuoto",
        "uz": "Bo'sh natija",
        "de": "Leeres Ergebnis",
        "ru": "Пустой результат",
        "es": "Resultado vacío",
        "uk": "Порожній результат",
    },
    "sponsorship_set_days": {
        "uk": "Тепер виберіть кількість днів, які хочете рекламуватися (відлік днів почнеться з дня оплати + зміщення від уже зайнятих днів).",
        "uz": "Endi siz reklama qilmoqchi bo'lgan kunlar sonini tanlang (kunlar to'langan kundan boshlanadi + ishg'ol qilingan kunlardan ofset).",
        "en": "Now select the number of days you want to advertise (the days will start from the date of payment + offset from the days already occupied).",
        "de": "Wählen Sie nun die Anzahl der Tage aus, die Sie inserieren möchten (die Tage beginnen mit dem Zahlungsdatum + Verrechnung mit den bereits belegten Tagen).",
        "it": "Ora seleziona il numero di giorni che vuoi pubblicizzare (i giorni inizieranno dalla data di pagamento + offset dai giorni già occupati).",
        "ru": "Теперь выберите количество дней, которые хотите рекламироваться (отсчёт дней начнётся со дня оплаты + смещение от уже занятых дней).",
        "es": "Ahora seleccione la cantidad de días que desea anunciar (los días comenzarán a partir de la fecha de pago + compensación de los días ya ocupados).",
        "id": "Sekarang pilih jumlah hari yang ingin Anda iklankan (hari akan dimulai dari tanggal pembayaran + offset dari hari yang sudah ditempati).",
        "pt": "Agora selecione o número de dias que deseja anunciar (os dias começarão a partir da data do pagamento + compensação dos dias já ocupados).",
    },
    "Now press /start 👈": {
        "ru": "Теперь нажмите /start 👈",
        "es": "Ahora presione /start 👈",
        "de": "Drücken Sie nun /start 👈",
        "pt": "Agora pressione /start 👈",
        "uz": "Endi /start tugmachasini bosing 👈",
        "en": "Now press /start 👈",
        "uk": "Тепер натисніть /start 👈",
        "it": "Ora premi /start 👈",
        "id": "Sekarang tekan /start 👈",
    },
    "/help": map[string]string{
        "es": "Solo envíame un mensaje de texto y lo traduciré.\n\n Prueba también nuestro modo integrado 👇",
        "uk": "Просто надішліть мені текст, і я його перекладу\n\n Також спробуйте наш вбудований режим 👇",
        "uz": "Menga faqat matn yuboring, men uni tarjima qilaman\n\n Bizning o'rnatilgan rejimimizni ham ko'ring 👇",
        "de": "Schicken Sie mir einfach einen Text und ich übersetze ihn\n\n Probiere auch unseren eingebauten Modus 👇",
        "ru": "Просто пришлите мне текст и я его переведу\n\n Также попробуйте наш инлайн режим 👇",
        "pt": "Basta me enviar um texto e eu irei traduzi-lo\n\n Experimente também o nosso modo integrado 👇",
        "id": "Kirimkan saja saya teks dan saya akan menerjemahkannya\n\n Coba juga mode bawaan kami 👇",
        "en": "Just send me a text and I will translate it\n\nAlso try our built-in mode 👇",
        "it": "Mandami un messaggio e lo traduco\n\n Prova anche la nostra modalità integrata 👇",
    },
    "Just send me a text and I will translate it": {
        "en": "Just send me a text and I will translate it",
        "ru": "Просто пришлите мне текст, и я его переведу",
        "pt": "Basta me enviar um texto e eu irei traduzi-lo",
        "it": "Mandami un messaggio e lo traduco",
        "uz": "Menga faqat matn yuboring, men uni tarjima qilaman",
        "de": "Schicken Sie mir einfach einen Text und ich übersetze ihn",
        "es": "Solo envíame un mensaje de texto y lo traduciré.",
        "uk": "Просто надішліть мені текст, і я його перекладу",
        "id": "Kirimkan saja saya teks dan saya akan menerjemahkannya",
    },
    "пишите сюда": map[string]string{
        "it": "scrivere qui",
        "uk": "пишіть сюди",
        "de": "hier schreiben",
        "es": "escriba aqui",
        "ru": "пишите сюда",
        "uz": "shu yerga yozing",
        "en": "write here",
        "id": "tulis disini",
        "pt": "escreva aqui",
    },
    "Далее": {
        "de": "Weiter ▶",
        "es": "Siguiente ▶",
        "it": "Avanti ▶",
        "uk": "Далі ▶",
        "uz": "Keyingi ▶",
        "id": "Selanjutnya",
        "pt": "Próximo ▶",
        "ru": "Далее ▶",
        "en": "Next ▶",
    },
    "Failed to detect the language. Please enter something else": {
        "uz": "Til aniqlanmadi. Iltimos, yana bir narsani kiriting",
        "de": "Die Sprache konnte nicht erkannt werden. Bitte geben Sie etwas anderes ein",
        "ru": "Не удалось определить язык. Пожалуйста, введите что-нибудь еще",
        "es": "No se pudo detectar el idioma. Por favor ingrese algo más",
        "uk": "Не вдалося виявити мову. Будь ласка, введіть щось інше",
        "pt": "Falha ao detectar o idioma. Por favor, insira outra coisa",
        "id": "Gagal mendeteksi bahasa. Silakan masukkan sesuatu yang lain",
        "it": "Impossibile rilevare la lingua. Per favore inserisci qualcos'altro",
    },
    "Translate Language": {
        "uk": "Мова перекладу",
        "uz": "Tarjima qilish uchun til",
        "de": "Sprache zum Übersetzen",
        "en": "Translate Language",
        "es": "Idioma para traducir",
        "id": "Bahasa untuk menerjemahkan",
        "it": "Lingua per tradurre",
        "pt": "Língua para tradução",
        "ru": "Язык перевода",
    },
    "To voice": {
        "uk": "озвучити",
        "it": "esprimere",
        "ru": "Озвучить",
        "es": "a voz",
        "pt": "dar voz",
        "id": "untuk menyuarakan",
        "uz": "ovoz berish",
        "de": "aussprechen",
    },
    "Выберите ваш родной язык": {
        "de": "Wählen Sie Ihre Muttersprache",
        "es": "Elija su lengua materna",
        "pt": "Escolha sua língua nativa",
        "ru": "Выберите ваш родной язык",
        "en": "Choose your native language",
        "id": "Pilih bahasa ibu Anda",
        "it": "Scegli la tua lingua madre",
        "uk": "Виберіть ваш рідну мову",
        "uz": "Ona tilingizni tanlang",
    },
    "Введите целое число без лишних символов": {
        "uz": "Qo'shimcha belgilarsiz butun sonni kiriting",
        "id": "Masukkan bilangan bulat tanpa karakter tambahan",
        "pt": "Insira um número inteiro sem caracteres extras",
        "ru": "Введите целое число без лишних символов",
        "uk": "Введіть ціле число без зайвих символів",
        "de": "Geben Sie eine ganze Zahl ohne zusätzliche Zeichen ein",
        "es": "Ingrese un número entero sin caracteres adicionales",
        "it": "Inserisci un numero intero senza caratteri extra",
        "en": "Enter an integer without extra characters",
    },
    "%s language is not supported": {
        "en": "%s language is not supported",
        "uz": "%s tili qo'llab -quvvatlanmaydi",
        "de": "%s Sprache wird nicht unterstützt",
        "pt": "O idioma %s não é compatível",
        "id": "%s bahasa tidak didukung",
        "it": "%s lingua non è supportata",
        "ru": "%s язык не поддерживается",
        "es": "%s idioma no es compatible",
        "uk": "Мова %s не підтримується",
    },
    "v.": { // verb
        "de": "Verb",
        "es": "verbo",
        "pt": "verbo",
        "id": "kata kerja",
        "en": "verb",
        "it": "verbo",
        "uz": "fe'l",
        "ru": "глагол",
        "uk": "дієслово",
    },
    "infl.": { // verb
        "de": "Verb",
        "es": "verbo",
        "pt": "verbo",
        "id": "kata kerja",
        "en": "verb",
        "it": "verbo",
        "uz": "fe'l",
        "ru": "глагол",
        "uk": "дієслово",
    },
    "n.": { // noun
        "en": "noun",
        "de": "Substantiv",
        "ru": "сущ.",
        "uk": "іменник",
        "id": "kata benda",
        "it": "sostantivo",
        "uz": "ism",
        "es": "sustantivo",
        "pt": "substantivo",
    },
    "nf.": { // noun
        "en": "noun",
        "de": "Substantiv",
        "ru": "сущ.",
        "uk": "іменник",
        "id": "kata benda",
        "it": "sostantivo",
        "uz": "ism",
        "es": "sustantivo",
        "pt": "substantivo",
    },
    "nn.": { // noun
        "en": "noun",
        "de": "Substantiv",
        "ru": "сущ.",
        "uk": "іменник",
        "id": "kata benda",
        "it": "sostantivo",
        "uz": "ism",
        "es": "sustantivo",
        "pt": "substantivo",
    },
    "nm.": { // noun
        "en": "noun",
        "de": "Substantiv",
        "ru": "сущ.",
        "uk": "іменник",
        "id": "kata benda",
        "it": "sostantivo",
        "uz": "ism",
        "es": "sustantivo",
        "pt": "substantivo",
    },
    "adv.": {
        "es": "adverbio",
        "pt": "advérbio",
        "id": "kata keterangan",
        "en": "adverb",
        "it": "avverbio",
        "uz": "olmosh",
        "de": "Adverb",
        "ru": "нар.",
        "uk": "прислівник",
    },
    "conj.": {
        "en": "conjunction",
        "de": "Verbindung",
        "es": "conjunción",
        "id": "konjungsi",
        "it": "congiunzione",
        "uz": "birikma",
        "ru": "союз",
        "uk": "сполучення",
        "pt": "conjunção",
    },
    "pron.": {
        "it": "pronome",
        "uk": "займенник",
        "pt": "pronome",
        "id": "kata ganti",
        "en": "pronoun",
        "uz": "olmosh",
        "de": "Pronomen",
        "ru": "местоим.",
        "es": "pronombre",
    },
    "adj.": {
        "en": "adjective",
        "de": "Adjektiv",
        "es": "adjetivo",
        "id": "kata sifat",
        "it": "aggettivo",
        "uz": "sifat",
        "ru": "прил.",
        "uk": "прикметник",
        "pt": "adjetivo",
    },
    "Dictionary": {
        "en": "Dictionary",
        "it": "Dizionario",
        "uk": "Словник",
        "id": "Kamus",
        "uz": "Lug'at",
        "de": "Wörterbuch",
        "ru": "Словарь",
        "es": "Diccionario",
        "pt": "Dicionário",
    },
    "To voice source text": {
        "en": "To voice source",
        "it": "Voce del testo originale",
        "pt": "Soar o texto original",
        "uz": "Manba matnini tarjima qiling",
        "de": "Originaltext vortragen",
        "ru": "Озвучить исходник",
        "es": "Expresar el texto original",
        "uk": "Озвучити вихідний текст",
        "id": "Voice the source text",
    },
    "Select the source language of your text if it was not defined correctly": {
        "id": "Pilih bahasa asli teks Anda jika tidak didefinisikan dengan benar",
        "uz": "Agar matn to'g'ri aniqlanmagan bo'lsa, uning asl tilini tanlang",
        "de": "Wählen Sie die Originalsprache Ihres Textes aus, wenn diese nicht richtig definiert wurde",
        "ru": "Выберите исходный язык вашего текста, если он был определен неправильно",
        "es": "Seleccione el idioma original de su texto si no se definió correctamente",
        "en": "Select the source language of your text if it was not defined correctly",
        "it": "Seleziona la lingua originale del tuo testo se non è stata definita correttamente",
        "uk": "Виберіть вихідний мову вашого тексту, якщо він був визначений неправильно",
        "pt": "Selecione o idioma original do seu texto se não foi definido corretamente",
    },
    "Welcome. We are glad that you are with us again. ✋": {
        "en": "Welcome. We are glad that you are with us again. ✋",
        "it": "Accoglienza. Siamo lieti che tu sia di nuovo con noi. ✋",
        "uz": "Xush kelibsiz. Siz yana biz bilan ekanligingizdan xursandmiz. ✋",
        "ru": "Welcome. We are glad that you are with us again. ✋",
        "uk": "Welcome. We are glad that you are with us again. ✋",
        "pt": "Receber. Estamos felizes por você estar conosco novamente. ✋",
        "de": "Willkommen. Wir freuen uns, dass Sie wieder bei uns sind. ✋",
        "es": "Bienvenido. Nos alegra que vuelva a estar con nosotros. ✋",
        "id": "Selamat datang. Kami senang Anda bersama kami lagi. ✋",
    },
    "The original text language and the target language are the same, please set different": {
        "en": "The original text language and the target language are the same, please set different",
        "uz": "Asl matn va maqsadli til bir xil, iltimos, boshqasini o'rnating",
        "pt": "O idioma do texto original e o idioma de chegada são os mesmos, por favor, defina diferentes",
        "it": "La lingua del testo originale e la lingua di destinazione sono le stesse, impostane una diversa",
        "de": "Die Originaltextsprache und die Zielsprache sind gleich, bitte unterschiedlich einstellen",
        "ru": "Исходный язык текста и язык перевода одинаковы, установите разные",
        "es": "El idioma del texto original y el idioma de destino son los mismos, establezca diferentes",
        "uk": "Вхідна мова тексту і мова перекладу однакові, встановіть різні",
        "id": "Bahasa teks asli dan bahasa target sama, harap atur berbeda",
    },
    "Examples": {
        "en": "Examples",
        "it": "Esempi",
        "de": "Beispiele",
        "pt": "Exemplos",
        "uz": "Misollar",
        "ru": "Примеры",
        "es": "Ejemplos de",
        "uk": "Приклади",
        "id": "Contoh",
    },
    "Meaning": {
        "en": "Meaning",
        "uz": "Ma'nosi",
        "de": "Bedeutung",
        "es": "Sentido",
        "uk": "Значення",
        "id": "Arti",
        "it": "Significato",
        "ru": "Значение",
        "pt": "Significado",
    },
}


func localize(text, code string, placeholders ...interface{}) string {
    if trigger, ok := translations[text]; ok { // Текст подходит под варианты
        if result, ok := trigger[code]; ok {
            return fmt.Sprintf(result, placeholders...)
        }
    }
   return fmt.Sprintf(text, placeholders...)
}

func command(intent string) []string {
    intent = strings.ToLower(intent)
    v, ok := commandTranslations[intent]
    if !ok {
        return []string{}
    }
    return v
}