async function h() {
    const {session: e} = await l(["session"]);
    !e || (Date.now() - e.timestamp) / 6e4 > 30 ? await p() : await chrome.storage.sync.set({
        session: {
            ...e,
            timestamp: Date.now()
        }
    })
}
function translate(from, to, text) {
    await h();
    const t = await l(["selectedTargetLanguage", "session", "installationId"]), n = new Date, r = {
        sourceLang: "auto",
        targetLangCode: t.selectedTargetLanguage,
        sessionId: t.session.id,
        isHtml: !0,
        sourceTexts: e.payload.texts
    }, s = await Zt(r), o = {isFirstTime: e.payload.isFirstTime, processingTimeInMs: new Date - n};
    try {
        Mt({
            sessionId: t.session.id,
            installationId: t.installationId,
            trigger: 1,
            websiteLanguage: e.payload.sourceLang,
            targetLanguage: t.selectedTargetLanguage,
            characterCount: e.payload.texts.reduce(((e, t) => e + t.length), 0),
            translationPerformance: o
        })
    } catch (e) {
        console.error(e)
    }
    return s.texts
};
translate("en", "ru", "hey").then((e => n(e))).catch((e => {
    n(Rt(e))
}))