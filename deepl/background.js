var background = function (e) {
    "use strict";
    var t, n = new Uint8Array(16);

    function r() {
        if (!t && !(t = "undefined" != typeof crypto && crypto.getRandomValues && crypto.getRandomValues.bind(crypto) || "undefined" != typeof msCrypto && "function" == typeof msCrypto.getRandomValues && msCrypto.getRandomValues.bind(msCrypto))) throw new Error("crypto.getRandomValues() not supported. See https://github.com/uuidjs/uuid#getrandomvalues-not-supported");
        return t(n)
    }

    var s = /^(?:[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}|00000000-0000-0000-0000-000000000000)$/i;

    function o(e) {
        return "string" == typeof e && s.test(e)
    }

    for (var i = [], a = 0; a < 256; ++a) i.push((a + 256).toString(16).substr(1));

    function c(e, t, n) {
        var s = (e = e || {}).random || (e.rng || r)();
        if (s[6] = 15 & s[6] | 64, s[8] = 63 & s[8] | 128, t) {
            n = n || 0;
            for (var a = 0; a < 16; ++a) t[n + a] = s[a];
            return t
        }
        return function (e) {
            var t = arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : 0,
                n = (i[e[t + 0]] + i[e[t + 1]] + i[e[t + 2]] + i[e[t + 3]] + "-" + i[e[t + 4]] + i[e[t + 5]] + "-" + i[e[t + 6]] + i[e[t + 7]] + "-" + i[e[t + 8]] + i[e[t + 9]] + "-" + i[e[t + 10]] + i[e[t + 11]] + i[e[t + 12]] + i[e[t + 13]] + i[e[t + 14]] + i[e[t + 15]]).toLowerCase();
            if (!o(n)) throw TypeError("Stringified UUID is invalid");
            return n
        }(s)
    }

    function u() {
        const e = new Uint32Array(28);
        return crypto.getRandomValues(e), Array.from(e, (e => ("0" + e.toString(16)).slice(-2))).join("")
    }

    function d(e) {
        try {
            return JSON.parse(atob(e.split(".")[1]))
        } catch (e) {
            return null
        }
    }

    async function l(e) {
        if (e.includes("session")) {
            const e = await chrome.storage.sync.get(["session"]);
            e && !e.session && await p()
        }
        return new Promise(((t, n) => {
            chrome.storage.sync.get(e, (r => {
                r ? t(r) : n(`Can not get ${e}`)
            }))
        }))
    }

    function p() {
        return chrome.storage.sync.set({session: {id: c(), timestamp: Date.now()}})
    }

    async function h() {
        const {session: e} = await l(["session"]);
        !e || (Date.now() - e.timestamp) / 6e4 > 30 ? await p() : await chrome.storage.sync.set({
            session: {
                ...e,
                timestamp: Date.now()
            }
        })
    }

    async function g(e, t, n) {
        const r = d(e);
        await chrome.storage.local.set({
            a_t: e,
            r_t: t,
            i_t: n,
            isLoggedIn: !(!e || !t),
            isProUser: !(!r || "Pro" !== r.sessionType)
        })
    }

    async function f() {
        await chrome.storage.local.remove(["a_t", "r_t", "i_t"]), await chrome.storage.local.set({
            isLoggedIn: !1,
            isProUser: !1
        })
    }

    async function _() {
        const e = await chrome.storage.local.get(["a_t", "r_t", "i_t"]);
        return {access_token: e ? e.a_t : void 0, refresh_token: e ? e.r_t : void 0, id_token: e ? e.i_t : void 0}
    }

    async function m() {
        const e = await chrome.storage.local.get(["isLoggedIn"]);
        return !(!e || !e.isLoggedIn)
    }

    async function y() {
        const e = await chrome.storage.local.get(["isProUser"]);
        return !(!e || !e.isProUser)
    }

    const v = [{k: "DE", v: "German"}, {k: "EN", v: "English"}, {k: "FR", v: "French"}, {
            k: "ES",
            v: "Spanish"
        }, {k: "IT", v: "Italian"}, {k: "PL", v: "Polish"}, {k: "NL", v: "Dutch"}, {k: "PT", v: "Portuguese"}, {
            k: "RU",
            v: "Russian"
        }, {k: "ZH", v: "Chinese"}, {k: "JA", v: "Japanese"}, {k: "BG", v: "Bulgarian"}, {k: "CS", v: "Czech"}, {
            k: "DA",
            v: "Danish"
        }, {k: "ET", v: "Estonian"}, {k: "FI", v: "Finnish"}, {k: "EL", v: "Greek"}, {k: "HU", v: "Hungarian"}, {
            k: "LV",
            v: "Latvian"
        }, {k: "LT", v: "Lithuanian"}, {k: "RO", v: "Romanian"}, {k: "SK", v: "Slovak"}, {
            k: "SL",
            v: "Slovenian"
        }, {k: "SV", v: "Swedish"}], S = 0, T = "inline", w = "automatic", E = "manual", b = "input", x = 1, I = 2, R = 1,
        L = 2, k = 3, O = 1, C = 2, N = 1, P = 2, A = 3, D = "HINT_KEYBOARD_SHORTCUT",
        j = "INLINE_TRANSLATION_TRIGGER_ICON", M = "INPUT_TRANSLATION_TRIGGER_ICON", U = "trigger-translation",
        q = "CLICK_ON_SHORTCUT_NOTIFICATION", F = () => {
            const e = navigator.languages ? navigator.languages[0].substr(0, 2).toUpperCase() : "EN";
            return v.map((e => e.k)).includes(e) ? e : "EN"
        }, V = () => {
            const e = {};
            return v.forEach((t => {
                e[t.k] = 0
            })), e
        }, G = "production", H = "development";

    function $() {
        return G === H
    }

    /*! *****************************************************************************
  Copyright (c) Microsoft Corporation.

  Permission to use, copy, modify, and/or distribute this software for any
  purpose with or without fee is hereby granted.

  THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
  REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
  AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
  INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
  LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
  OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
  PERFORMANCE OF THIS SOFTWARE.
  ***************************************************************************** */
    var Y = function (e, t) {
        return Y = Object.setPrototypeOf || {__proto__: []} instanceof Array && function (e, t) {
            e.__proto__ = t
        } || function (e, t) {
            for (var n in t) t.hasOwnProperty(n) && (e[n] = t[n])
        }, Y(e, t)
    };
    var J, W, B, K, z = function () {
        return z = Object.assign || function (e) {
            for (var t, n = 1, r = arguments.length; n < r; n++) for (var s in t = arguments[n]) Object.prototype.hasOwnProperty.call(t, s) && (e[s] = t[s]);
            return e
        }, z.apply(this, arguments)
    };

    function Q(e) {
        var t = "function" == typeof Symbol && Symbol.iterator, n = t && e[t], r = 0;
        if (n) return n.call(e);
        if (e && "number" == typeof e.length) return {
            next: function () {
                return e && r >= e.length && (e = void 0), {value: e && e[r++], done: !e}
            }
        };
        throw new TypeError(t ? "Object is not iterable." : "Symbol.iterator is not defined.")
    }

    function X(e, t) {
        var n = "function" == typeof Symbol && e[Symbol.iterator];
        if (!n) return e;
        var r, s, o = n.call(e), i = [];
        try {
            for (; (void 0 === t || t-- > 0) && !(r = o.next()).done;) i.push(r.value)
        } catch (e) {
            s = {error: e}
        } finally {
            try {
                r && !r.done && (n = o.return) && n.call(o)
            } finally {
                if (s) throw s.error
            }
        }
        return i
    }

    function Z() {
        for (var e = [], t = 0; t < arguments.length; t++) e = e.concat(X(arguments[t]));
        return e
    }

    function ee() {
        return "[object process]" === Object.prototype.toString.call("undefined" != typeof process ? process : 0)
    }

    function te(e, t) {
        return e.require(t)
    }

    !function (e) {
        e.Ok = "ok", e.Exited = "exited", e.Crashed = "crashed", e.Abnormal = "abnormal"
    }(J || (J = {})), function (e) {
        e.Ok = "ok", e.Errored = "errored", e.Crashed = "crashed"
    }(W || (W = {})), function (e) {
        e.Explicit = "explicitly_set", e.Sampler = "client_sampler", e.Rate = "client_rate", e.Inheritance = "inheritance"
    }(B || (B = {})), function (e) {
        e.BeforeSend = "before_send", e.EventProcessor = "event_processor", e.NetworkError = "network_error", e.QueueOverflow = "queue_overflow", e.RateLimitBackoff = "ratelimit_backoff", e.SampleRate = "sample_rate"
    }(K || (K = {}));
    var ne = {};

    function re() {
        return ee() ? global : "undefined" != typeof window ? window : "undefined" != typeof self ? self : ne
    }

    function se(e) {
        return "[object String]" === Object.prototype.toString.call(e)
    }

    function oe(e) {
        return "[object Object]" === Object.prototype.toString.call(e)
    }

    function ie(e) {
        return Boolean(e && e.then && "function" == typeof e.then)
    }

    function ae(e, t) {
        try {
            return e instanceof t
        } catch (e) {
            return !1
        }
    }

    var ce = re(), ue = "Sentry Logger ";

    function de(e) {
        var t = re();
        if (!("console" in t)) return e();
        var n = t.console, r = {};
        ["debug", "info", "warn", "error", "log", "assert"].forEach((function (e) {
            e in t.console && n[e].__sentry_original__ && (r[e] = n[e], n[e] = n[e].__sentry_original__)
        }));
        var s = e();
        return Object.keys(r).forEach((function (e) {
            n[e] = r[e]
        })), s
    }

    var le = function () {
        function e() {
            this._enabled = !1
        }

        return e.prototype.disable = function () {
            this._enabled = !1
        }, e.prototype.enable = function () {
            this._enabled = !0
        }, e.prototype.log = function () {
            for (var e = [], t = 0; t < arguments.length; t++) e[t] = arguments[t];
            this._enabled && de((function () {
                ce.console.log(ue + "[Log]: " + e.join(" "))
            }))
        }, e.prototype.warn = function () {
            for (var e = [], t = 0; t < arguments.length; t++) e[t] = arguments[t];
            this._enabled && de((function () {
                ce.console.warn(ue + "[Warn]: " + e.join(" "))
            }))
        }, e.prototype.error = function () {
            for (var e = [], t = 0; t < arguments.length; t++) e[t] = arguments[t];
            this._enabled && de((function () {
                ce.console.error(ue + "[Error]: " + e.join(" "))
            }))
        }, e
    }();
    ce.__SENTRY__ = ce.__SENTRY__ || {};
    var pe = ce.__SENTRY__.logger || (ce.__SENTRY__.logger = new le), he = "<anonymous>";

    function ge(e) {
        try {
            return e && "function" == typeof e && e.name || he
        } catch (e) {
            return he
        }
    }

    function fe(e, t, n) {
        if (t in e) {
            var r = e[t], s = n(r);
            if ("function" == typeof s) try {
                s.prototype = s.prototype || {}, Object.defineProperties(s, {
                    __sentry_original__: {
                        enumerable: !1,
                        value: r
                    }
                })
            } catch (e) {
            }
            e[t] = s
        }
    }

    function _e(e) {
        var t, n;
        if (oe(e)) {
            var r = e, s = {};
            try {
                for (var o = Q(Object.keys(r)), i = o.next(); !i.done; i = o.next()) {
                    var a = i.value;
                    void 0 !== r[a] && (s[a] = _e(r[a]))
                }
            } catch (e) {
                t = {error: e}
            } finally {
                try {
                    i && !i.done && (n = o.return) && n.call(o)
                } finally {
                    if (t) throw t.error
                }
            }
            return s
        }
        return Array.isArray(e) ? e.map(_e) : e
    }

    function me(e) {
        return e && /^function fetch\(\)\s+\{\s+\[native code\]\s+\}$/.test(e.toString())
    }

    function ye() {
        if (!function () {
            if (!("fetch" in re())) return !1;
            try {
                return new Headers, new Request(""), new Response, !0
            } catch (e) {
                return !1
            }
        }()) return !1;
        var e = re();
        if (me(e.fetch)) return !0;
        var t = !1, n = e.document;
        if (n && "function" == typeof n.createElement) try {
            var r = n.createElement("iframe");
            r.hidden = !0, n.head.appendChild(r), r.contentWindow && r.contentWindow.fetch && (t = me(r.contentWindow.fetch)), n.head.removeChild(r)
        } catch (e) {
            pe.warn("Could not create sandbox iframe for pure fetch check, bailing to window.fetch: ", e)
        }
        return t
    }

    var ve, Se = re(), Te = {}, we = {};

    function Ee(e) {
        if (!we[e]) switch (we[e] = !0, e) {
            case"console":
                !function () {
                    if (!("console" in Se)) return;
                    ["debug", "info", "warn", "error", "log", "assert"].forEach((function (e) {
                        e in Se.console && fe(Se.console, e, (function (t) {
                            return function () {
                                for (var n = [], r = 0; r < arguments.length; r++) n[r] = arguments[r];
                                xe("console", {args: n, level: e}), t && Function.prototype.apply.call(t, Se.console, n)
                            }
                        }))
                    }))
                }();
                break;
            case"dom":
                !function () {
                    if (!("document" in Se)) return;
                    var e = xe.bind(null, "dom"), t = Oe(e, !0);
                    Se.document.addEventListener("click", t, !1), Se.document.addEventListener("keypress", t, !1), ["EventTarget", "Node"].forEach((function (t) {
                        var n = Se[t] && Se[t].prototype;
                        n && n.hasOwnProperty && n.hasOwnProperty("addEventListener") && (fe(n, "addEventListener", (function (t) {
                            return function (n, r, s) {
                                if ("click" === n || "keypress" == n) try {
                                    var o = this,
                                        i = o.__sentry_instrumentation_handlers__ = o.__sentry_instrumentation_handlers__ || {},
                                        a = i[n] = i[n] || {refCount: 0};
                                    if (!a.handler) {
                                        var c = Oe(e);
                                        a.handler = c, t.call(this, n, c, s)
                                    }
                                    a.refCount += 1
                                } catch (e) {
                                }
                                return t.call(this, n, r, s)
                            }
                        })), fe(n, "removeEventListener", (function (e) {
                            return function (t, n, r) {
                                if ("click" === t || "keypress" == t) try {
                                    var s = this, o = s.__sentry_instrumentation_handlers__ || {}, i = o[t];
                                    i && (i.refCount -= 1, i.refCount <= 0 && (e.call(this, t, i.handler, r), i.handler = void 0, delete o[t]), 0 === Object.keys(o).length && delete s.__sentry_instrumentation_handlers__)
                                } catch (e) {
                                }
                                return e.call(this, t, n, r)
                            }
                        })))
                    }))
                }();
                break;
            case"xhr":
                !function () {
                    if (!("XMLHttpRequest" in Se)) return;
                    var e = [], t = [], n = XMLHttpRequest.prototype;
                    fe(n, "open", (function (n) {
                        return function () {
                            for (var r = [], s = 0; s < arguments.length; s++) r[s] = arguments[s];
                            var o = this, i = r[1];
                            o.__sentry_xhr__ = {
                                method: se(r[0]) ? r[0].toUpperCase() : r[0],
                                url: r[1]
                            }, se(i) && "POST" === o.__sentry_xhr__.method && i.match(/sentry_key/) && (o.__sentry_own_request__ = !0);
                            var a = function () {
                                if (4 === o.readyState) {
                                    try {
                                        o.__sentry_xhr__ && (o.__sentry_xhr__.status_code = o.status)
                                    } catch (e) {
                                    }
                                    try {
                                        var n = e.indexOf(o);
                                        if (-1 !== n) {
                                            e.splice(n);
                                            var s = t.splice(n)[0];
                                            o.__sentry_xhr__ && void 0 !== s[0] && (o.__sentry_xhr__.body = s[0])
                                        }
                                    } catch (e) {
                                    }
                                    xe("xhr", {args: r, endTimestamp: Date.now(), startTimestamp: Date.now(), xhr: o})
                                }
                            };
                            return "onreadystatechange" in o && "function" == typeof o.onreadystatechange ? fe(o, "onreadystatechange", (function (e) {
                                return function () {
                                    for (var t = [], n = 0; n < arguments.length; n++) t[n] = arguments[n];
                                    return a(), e.apply(o, t)
                                }
                            })) : o.addEventListener("readystatechange", a), n.apply(o, r)
                        }
                    })), fe(n, "send", (function (n) {
                        return function () {
                            for (var r = [], s = 0; s < arguments.length; s++) r[s] = arguments[s];
                            return e.push(this), t.push(r), xe("xhr", {
                                args: r,
                                startTimestamp: Date.now(),
                                xhr: this
                            }), n.apply(this, r)
                        }
                    }))
                }();
                break;
            case"fetch":
                !function () {
                    if (!ye()) return;
                    fe(Se, "fetch", (function (e) {
                        return function () {
                            for (var t = [], n = 0; n < arguments.length; n++) t[n] = arguments[n];
                            var r = {args: t, fetchData: {method: Ie(t), url: Re(t)}, startTimestamp: Date.now()};
                            return xe("fetch", z({}, r)), e.apply(Se, t).then((function (e) {
                                return xe("fetch", z(z({}, r), {endTimestamp: Date.now(), response: e})), e
                            }), (function (e) {
                                throw xe("fetch", z(z({}, r), {endTimestamp: Date.now(), error: e})), e
                            }))
                        }
                    }))
                }();
                break;
            case"history":
                !function () {
                    if (!function () {
                        var e = re(), t = e.chrome, n = t && t.app && t.app.runtime,
                            r = "history" in e && !!e.history.pushState && !!e.history.replaceState;
                        return !n && r
                    }()) return;
                    var e = Se.onpopstate;

                    function t(e) {
                        return function () {
                            for (var t = [], n = 0; n < arguments.length; n++) t[n] = arguments[n];
                            var r = t.length > 2 ? t[2] : void 0;
                            if (r) {
                                var s = ve, o = String(r);
                                ve = o, xe("history", {from: s, to: o})
                            }
                            return e.apply(this, t)
                        }
                    }

                    Se.onpopstate = function () {
                        for (var t = [], n = 0; n < arguments.length; n++) t[n] = arguments[n];
                        var r = Se.location.href, s = ve;
                        if (ve = r, xe("history", {from: s, to: r}), e) try {
                            return e.apply(this, t)
                        } catch (e) {
                        }
                    }, fe(Se.history, "pushState", t), fe(Se.history, "replaceState", t)
                }();
                break;
            case"error":
                Ce = Se.onerror, Se.onerror = function (e, t, n, r, s) {
                    return xe("error", {
                        column: r,
                        error: s,
                        line: n,
                        msg: e,
                        url: t
                    }), !!Ce && Ce.apply(this, arguments)
                };
                break;
            case"unhandledrejection":
                Ne = Se.onunhandledrejection, Se.onunhandledrejection = function (e) {
                    return xe("unhandledrejection", e), !Ne || Ne.apply(this, arguments)
                };
                break;
            default:
                pe.warn("unknown instrumentation type:", e)
        }
    }

    function be(e) {
        e && "string" == typeof e.type && "function" == typeof e.callback && (Te[e.type] = Te[e.type] || [], Te[e.type].push(e.callback), Ee(e.type))
    }

    function xe(e, t) {
        var n, r;
        if (e && Te[e]) try {
            for (var s = Q(Te[e] || []), o = s.next(); !o.done; o = s.next()) {
                var i = o.value;
                try {
                    i(t)
                } catch (t) {
                    pe.error("Error while triggering instrumentation handler.\nType: " + e + "\nName: " + ge(i) + "\nError: " + t)
                }
            }
        } catch (e) {
            n = {error: e}
        } finally {
            try {
                o && !o.done && (r = s.return) && r.call(s)
            } finally {
                if (n) throw n.error
            }
        }
    }

    function Ie(e) {
        return void 0 === e && (e = []), "Request" in Se && ae(e[0], Request) && e[0].method ? String(e[0].method).toUpperCase() : e[1] && e[1].method ? String(e[1].method).toUpperCase() : "GET"
    }

    function Re(e) {
        return void 0 === e && (e = []), "string" == typeof e[0] ? e[0] : "Request" in Se && ae(e[0], Request) ? e[0].url : String(e[0])
    }

    var Le, ke;

    function Oe(e, t) {
        return void 0 === t && (t = !1), function (n) {
            if (n && ke !== n && !function (e) {
                if ("keypress" !== e.type) return !1;
                try {
                    var t = e.target;
                    if (!t || !t.tagName) return !0;
                    if ("INPUT" === t.tagName || "TEXTAREA" === t.tagName || t.isContentEditable) return !1
                } catch (e) {
                }
                return !0
            }(n)) {
                var r = "keypress" === n.type ? "input" : n.type;
                (void 0 === Le || function (e, t) {
                    if (!e) return !0;
                    if (e.type !== t.type) return !0;
                    try {
                        if (e.target !== t.target) return !0
                    } catch (e) {
                    }
                    return !1
                }(ke, n)) && (e({
                    event: n,
                    name: r,
                    global: t
                }), ke = n), clearTimeout(Le), Le = Se.setTimeout((function () {
                    Le = void 0
                }), 1e3)
            }
        }
    }

    var Ce = null;
    var Ne = null;

    function Pe() {
        var e = re(), t = e.crypto || e.msCrypto;
        if (void 0 !== t && t.getRandomValues) {
            var n = new Uint16Array(8);
            t.getRandomValues(n), n[3] = 4095 & n[3] | 16384, n[4] = 16383 & n[4] | 32768;
            var r = function (e) {
                for (var t = e.toString(16); t.length < 4;) t = "0" + t;
                return t
            };
            return r(n[0]) + r(n[1]) + r(n[2]) + r(n[3]) + r(n[4]) + r(n[5]) + r(n[6]) + r(n[7])
        }
        return "xxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx".replace(/[xy]/g, (function (e) {
            var t = 16 * Math.random() | 0;
            return ("x" === e ? t : 3 & t | 8).toString(16)
        }))
    }

    var Ae = function () {
        function e(e) {
            var t = this;
            this._state = "PENDING", this._handlers = [], this._resolve = function (e) {
                t._setResult("RESOLVED", e)
            }, this._reject = function (e) {
                t._setResult("REJECTED", e)
            }, this._setResult = function (e, n) {
                "PENDING" === t._state && (ie(n) ? n.then(t._resolve, t._reject) : (t._state = e, t._value = n, t._executeHandlers()))
            }, this._attachHandler = function (e) {
                t._handlers = t._handlers.concat(e), t._executeHandlers()
            }, this._executeHandlers = function () {
                if ("PENDING" !== t._state) {
                    var e = t._handlers.slice();
                    t._handlers = [], e.forEach((function (e) {
                        e.done || ("RESOLVED" === t._state && e.onfulfilled && e.onfulfilled(t._value), "REJECTED" === t._state && e.onrejected && e.onrejected(t._value), e.done = !0)
                    }))
                }
            };
            try {
                e(this._resolve, this._reject)
            } catch (e) {
                this._reject(e)
            }
        }

        return e.resolve = function (t) {
            return new e((function (e) {
                e(t)
            }))
        }, e.reject = function (t) {
            return new e((function (e, n) {
                n(t)
            }))
        }, e.all = function (t) {
            return new e((function (n, r) {
                if (Array.isArray(t)) if (0 !== t.length) {
                    var s = t.length, o = [];
                    t.forEach((function (t, i) {
                        e.resolve(t).then((function (e) {
                            o[i] = e, 0 === (s -= 1) && n(o)
                        })).then(null, r)
                    }))
                } else n([]); else r(new TypeError("Promise.all requires an array as input."))
            }))
        }, e.prototype.then = function (t, n) {
            var r = this;
            return new e((function (e, s) {
                r._attachHandler({
                    done: !1, onfulfilled: function (n) {
                        if (t) try {
                            return void e(t(n))
                        } catch (e) {
                            return void s(e)
                        } else e(n)
                    }, onrejected: function (t) {
                        if (n) try {
                            return void e(n(t))
                        } catch (e) {
                            return void s(e)
                        } else s(t)
                    }
                })
            }))
        }, e.prototype.catch = function (e) {
            return this.then((function (e) {
                return e
            }), e)
        }, e.prototype.finally = function (t) {
            var n = this;
            return new e((function (e, r) {
                var s, o;
                return n.then((function (e) {
                    o = !1, s = e, t && t()
                }), (function (e) {
                    o = !0, s = e, t && t()
                })).then((function () {
                    o ? r(s) : e(s)
                }))
            }))
        }, e.prototype.toString = function () {
            return "[object SyncPromise]"
        }, e
    }(), De = {
        nowSeconds: function () {
            return Date.now() / 1e3
        }
    };
    var je = ee() ? function () {
        try {
            return te(module, "perf_hooks").performance
        } catch (e) {
            return
        }
    }() : function () {
        var e = re().performance;
        if (e && e.now) return {
            now: function () {
                return e.now()
            }, timeOrigin: Date.now() - e.now()
        }
    }(), Me = void 0 === je ? De : {
        nowSeconds: function () {
            return (je.timeOrigin + je.now()) / 1e3
        }
    }, Ue = De.nowSeconds.bind(De), qe = Me.nowSeconds.bind(Me), Fe = qe;
    !function () {
        var e = re().performance;
        if (e && e.now) {
            var t = 36e5, n = e.now(), r = Date.now(), s = e.timeOrigin ? Math.abs(e.timeOrigin + n - r) : t, o = s < t,
                i = e.timing && e.timing.navigationStart, a = "number" == typeof i ? Math.abs(i + n - r) : t;
            (o || a < t) && (s <= a && e.timeOrigin)
        }
    }();
    var Ve = function () {
        function e() {
            this._notifyingListeners = !1, this._scopeListeners = [], this._eventProcessors = [], this._breadcrumbs = [], this._user = {}, this._tags = {}, this._extra = {}, this._contexts = {}
        }

        return e.clone = function (t) {
            var n = new e;
            return t && (n._breadcrumbs = Z(t._breadcrumbs), n._tags = z({}, t._tags), n._extra = z({}, t._extra), n._contexts = z({}, t._contexts), n._user = t._user, n._level = t._level, n._span = t._span, n._session = t._session, n._transactionName = t._transactionName, n._fingerprint = t._fingerprint, n._eventProcessors = Z(t._eventProcessors), n._requestSession = t._requestSession), n
        }, e.prototype.addScopeListener = function (e) {
            this._scopeListeners.push(e)
        }, e.prototype.addEventProcessor = function (e) {
            return this._eventProcessors.push(e), this
        }, e.prototype.setUser = function (e) {
            return this._user = e || {}, this._session && this._session.update({user: e}), this._notifyScopeListeners(), this
        }, e.prototype.getUser = function () {
            return this._user
        }, e.prototype.getRequestSession = function () {
            return this._requestSession
        }, e.prototype.setRequestSession = function (e) {
            return this._requestSession = e, this
        }, e.prototype.setTags = function (e) {
            return this._tags = z(z({}, this._tags), e), this._notifyScopeListeners(), this
        }, e.prototype.setTag = function (e, t) {
            var n;
            return this._tags = z(z({}, this._tags), ((n = {})[e] = t, n)), this._notifyScopeListeners(), this
        }, e.prototype.setExtras = function (e) {
            return this._extra = z(z({}, this._extra), e), this._notifyScopeListeners(), this
        }, e.prototype.setExtra = function (e, t) {
            var n;
            return this._extra = z(z({}, this._extra), ((n = {})[e] = t, n)), this._notifyScopeListeners(), this
        }, e.prototype.setFingerprint = function (e) {
            return this._fingerprint = e, this._notifyScopeListeners(), this
        }, e.prototype.setLevel = function (e) {
            return this._level = e, this._notifyScopeListeners(), this
        }, e.prototype.setTransactionName = function (e) {
            return this._transactionName = e, this._notifyScopeListeners(), this
        }, e.prototype.setTransaction = function (e) {
            return this.setTransactionName(e)
        }, e.prototype.setContext = function (e, t) {
            var n;
            return null === t ? delete this._contexts[e] : this._contexts = z(z({}, this._contexts), ((n = {})[e] = t, n)), this._notifyScopeListeners(), this
        }, e.prototype.setSpan = function (e) {
            return this._span = e, this._notifyScopeListeners(), this
        }, e.prototype.getSpan = function () {
            return this._span
        }, e.prototype.getTransaction = function () {
            var e, t, n, r, s = this.getSpan();
            return (null === (e = s) || void 0 === e ? void 0 : e.transaction) ? null === (t = s) || void 0 === t ? void 0 : t.transaction : (null === (r = null === (n = s) || void 0 === n ? void 0 : n.spanRecorder) || void 0 === r ? void 0 : r.spans[0]) ? s.spanRecorder.spans[0] : void 0
        }, e.prototype.setSession = function (e) {
            return e ? this._session = e : delete this._session, this._notifyScopeListeners(), this
        }, e.prototype.getSession = function () {
            return this._session
        }, e.prototype.update = function (t) {
            if (!t) return this;
            if ("function" == typeof t) {
                var n = t(this);
                return n instanceof e ? n : this
            }
            return t instanceof e ? (this._tags = z(z({}, this._tags), t._tags), this._extra = z(z({}, this._extra), t._extra), this._contexts = z(z({}, this._contexts), t._contexts), t._user && Object.keys(t._user).length && (this._user = t._user), t._level && (this._level = t._level), t._fingerprint && (this._fingerprint = t._fingerprint), t._requestSession && (this._requestSession = t._requestSession)) : oe(t) && (t = t, this._tags = z(z({}, this._tags), t.tags), this._extra = z(z({}, this._extra), t.extra), this._contexts = z(z({}, this._contexts), t.contexts), t.user && (this._user = t.user), t.level && (this._level = t.level), t.fingerprint && (this._fingerprint = t.fingerprint), t.requestSession && (this._requestSession = t.requestSession)), this
        }, e.prototype.clear = function () {
            return this._breadcrumbs = [], this._tags = {}, this._extra = {}, this._user = {}, this._contexts = {}, this._level = void 0, this._transactionName = void 0, this._fingerprint = void 0, this._requestSession = void 0, this._span = void 0, this._session = void 0, this._notifyScopeListeners(), this
        }, e.prototype.addBreadcrumb = function (e, t) {
            var n = "number" == typeof t ? Math.min(t, 100) : 100;
            if (n <= 0) return this;
            var r = z({timestamp: Ue()}, e);
            return this._breadcrumbs = Z(this._breadcrumbs, [r]).slice(-n), this._notifyScopeListeners(), this
        }, e.prototype.clearBreadcrumbs = function () {
            return this._breadcrumbs = [], this._notifyScopeListeners(), this
        }, e.prototype.applyToEvent = function (e, t) {
            var n;
            if (this._extra && Object.keys(this._extra).length && (e.extra = z(z({}, this._extra), e.extra)), this._tags && Object.keys(this._tags).length && (e.tags = z(z({}, this._tags), e.tags)), this._user && Object.keys(this._user).length && (e.user = z(z({}, this._user), e.user)), this._contexts && Object.keys(this._contexts).length && (e.contexts = z(z({}, this._contexts), e.contexts)), this._level && (e.level = this._level), this._transactionName && (e.transaction = this._transactionName), this._span) {
                e.contexts = z({trace: this._span.getTraceContext()}, e.contexts);
                var r = null === (n = this._span.transaction) || void 0 === n ? void 0 : n.name;
                r && (e.tags = z({transaction: r}, e.tags))
            }
            return this._applyFingerprint(e), e.breadcrumbs = Z(e.breadcrumbs || [], this._breadcrumbs), e.breadcrumbs = e.breadcrumbs.length > 0 ? e.breadcrumbs : void 0, this._notifyEventProcessors(Z(function () {
                var e = re();
                return e.__SENTRY__ = e.__SENTRY__ || {}, e.__SENTRY__.globalEventProcessors = e.__SENTRY__.globalEventProcessors || [], e.__SENTRY__.globalEventProcessors
            }(), this._eventProcessors), e, t)
        }, e.prototype._notifyEventProcessors = function (e, t, n, r) {
            var s = this;
            return void 0 === r && (r = 0), new Ae((function (o, i) {
                var a = e[r];
                if (null === t || "function" != typeof a) o(t); else {
                    var c = a(z({}, t), n);
                    ie(c) ? c.then((function (t) {
                        return s._notifyEventProcessors(e, t, n, r + 1).then(o)
                    })).then(null, i) : s._notifyEventProcessors(e, c, n, r + 1).then(o).then(null, i)
                }
            }))
        }, e.prototype._notifyScopeListeners = function () {
            var e = this;
            this._notifyingListeners || (this._notifyingListeners = !0, this._scopeListeners.forEach((function (t) {
                t(e)
            })), this._notifyingListeners = !1)
        }, e.prototype._applyFingerprint = function (e) {
            e.fingerprint = e.fingerprint ? Array.isArray(e.fingerprint) ? e.fingerprint : [e.fingerprint] : [], this._fingerprint && (e.fingerprint = e.fingerprint.concat(this._fingerprint)), e.fingerprint && !e.fingerprint.length && delete e.fingerprint
        }, e
    }();
    var Ge, He = function () {
        function e(e) {
            this.errors = 0, this.sid = Pe(), this.duration = 0, this.status = J.Ok, this.init = !0, this.ignoreDuration = !1;
            var t = qe();
            this.timestamp = t, this.started = t, e && this.update(e)
        }

        return e.prototype.update = function (e) {
            if (void 0 === e && (e = {}), e.user && (!this.ipAddress && e.user.ip_address && (this.ipAddress = e.user.ip_address), this.did || e.did || (this.did = e.user.id || e.user.email || e.user.username)), this.timestamp = e.timestamp || qe(), e.ignoreDuration && (this.ignoreDuration = e.ignoreDuration), e.sid && (this.sid = 32 === e.sid.length ? e.sid : Pe()), void 0 !== e.init && (this.init = e.init), !this.did && e.did && (this.did = "" + e.did), "number" == typeof e.started && (this.started = e.started), this.ignoreDuration) this.duration = void 0; else if ("number" == typeof e.duration) this.duration = e.duration; else {
                var t = this.timestamp - this.started;
                this.duration = t >= 0 ? t : 0
            }
            e.release && (this.release = e.release), e.environment && (this.environment = e.environment), !this.ipAddress && e.ipAddress && (this.ipAddress = e.ipAddress), !this.userAgent && e.userAgent && (this.userAgent = e.userAgent), "number" == typeof e.errors && (this.errors = e.errors), e.status && (this.status = e.status)
        }, e.prototype.close = function (e) {
            e ? this.update({status: e}) : this.status === J.Ok ? this.update({status: J.Exited}) : this.update()
        }, e.prototype.toJSON = function () {
            return _e({
                sid: "" + this.sid,
                init: this.init,
                started: new Date(1e3 * this.started).toISOString(),
                timestamp: new Date(1e3 * this.timestamp).toISOString(),
                status: this.status,
                errors: this.errors,
                did: "number" == typeof this.did || "string" == typeof this.did ? "" + this.did : void 0,
                duration: this.duration,
                attrs: _e({
                    release: this.release,
                    environment: this.environment,
                    ip_address: this.ipAddress,
                    user_agent: this.userAgent
                })
            })
        }, e
    }(), $e = function () {
        function e(e, t, n) {
            void 0 === t && (t = new Ve), void 0 === n && (n = 4), this._version = n, this._stack = [{}], this.getStackTop().scope = t, e && this.bindClient(e)
        }

        return e.prototype.isOlderThan = function (e) {
            return this._version < e
        }, e.prototype.bindClient = function (e) {
            this.getStackTop().client = e, e && e.setupIntegrations && e.setupIntegrations()
        }, e.prototype.pushScope = function () {
            var e = Ve.clone(this.getScope());
            return this.getStack().push({client: this.getClient(), scope: e}), e
        }, e.prototype.popScope = function () {
            return !(this.getStack().length <= 1) && !!this.getStack().pop()
        }, e.prototype.withScope = function (e) {
            var t = this.pushScope();
            try {
                e(t)
            } finally {
                this.popScope()
            }
        }, e.prototype.getClient = function () {
            return this.getStackTop().client
        }, e.prototype.getScope = function () {
            return this.getStackTop().scope
        }, e.prototype.getStack = function () {
            return this._stack
        }, e.prototype.getStackTop = function () {
            return this._stack[this._stack.length - 1]
        }, e.prototype.captureException = function (e, t) {
            var n = this._lastEventId = Pe(), r = t;
            if (!t) {
                var s = void 0;
                try {
                    throw new Error("Sentry syntheticException")
                } catch (e) {
                    s = e
                }
                r = {originalException: e, syntheticException: s}
            }
            return this._invokeClient("captureException", e, z(z({}, r), {event_id: n})), n
        }, e.prototype.captureMessage = function (e, t, n) {
            var r = this._lastEventId = Pe(), s = n;
            if (!n) {
                var o = void 0;
                try {
                    throw new Error(e)
                } catch (e) {
                    o = e
                }
                s = {originalException: e, syntheticException: o}
            }
            return this._invokeClient("captureMessage", e, t, z(z({}, s), {event_id: r})), r
        }, e.prototype.captureEvent = function (e, t) {
            var n = Pe();
            return "transaction" !== e.type && (this._lastEventId = n), this._invokeClient("captureEvent", e, z(z({}, t), {event_id: n})), n
        }, e.prototype.lastEventId = function () {
            return this._lastEventId
        }, e.prototype.addBreadcrumb = function (e, t) {
            var n = this.getStackTop(), r = n.scope, s = n.client;
            if (r && s) {
                var o = s.getOptions && s.getOptions() || {}, i = o.beforeBreadcrumb, a = void 0 === i ? null : i,
                    c = o.maxBreadcrumbs, u = void 0 === c ? 100 : c;
                if (!(u <= 0)) {
                    var d = Ue(), l = z({timestamp: d}, e), p = a ? de((function () {
                        return a(l, t)
                    })) : l;
                    null !== p && r.addBreadcrumb(p, u)
                }
            }
        }, e.prototype.setUser = function (e) {
            var t = this.getScope();
            t && t.setUser(e)
        }, e.prototype.setTags = function (e) {
            var t = this.getScope();
            t && t.setTags(e)
        }, e.prototype.setExtras = function (e) {
            var t = this.getScope();
            t && t.setExtras(e)
        }, e.prototype.setTag = function (e, t) {
            var n = this.getScope();
            n && n.setTag(e, t)
        }, e.prototype.setExtra = function (e, t) {
            var n = this.getScope();
            n && n.setExtra(e, t)
        }, e.prototype.setContext = function (e, t) {
            var n = this.getScope();
            n && n.setContext(e, t)
        }, e.prototype.configureScope = function (e) {
            var t = this.getStackTop(), n = t.scope, r = t.client;
            n && r && e(n)
        }, e.prototype.run = function (e) {
            var t = Je(this);
            try {
                e(this)
            } finally {
                Je(t)
            }
        }, e.prototype.getIntegration = function (e) {
            var t = this.getClient();
            if (!t) return null;
            try {
                return t.getIntegration(e)
            } catch (t) {
                return pe.warn("Cannot retrieve integration " + e.id + " from the current Hub"), null
            }
        }, e.prototype.startSpan = function (e) {
            return this._callExtensionMethod("startSpan", e)
        }, e.prototype.startTransaction = function (e, t) {
            return this._callExtensionMethod("startTransaction", e, t)
        }, e.prototype.traceHeaders = function () {
            return this._callExtensionMethod("traceHeaders")
        }, e.prototype.captureSession = function (e) {
            if (void 0 === e && (e = !1), e) return this.endSession();
            this._sendSessionUpdate()
        }, e.prototype.endSession = function () {
            var e, t, n, r, s;
            null === (n = null === (t = null === (e = this.getStackTop()) || void 0 === e ? void 0 : e.scope) || void 0 === t ? void 0 : t.getSession()) || void 0 === n || n.close(), this._sendSessionUpdate(), null === (s = null === (r = this.getStackTop()) || void 0 === r ? void 0 : r.scope) || void 0 === s || s.setSession()
        }, e.prototype.startSession = function (e) {
            var t = this.getStackTop(), n = t.scope, r = t.client, s = r && r.getOptions() || {}, o = s.release,
                i = s.environment, a = (re().navigator || {}).userAgent,
                c = new He(z(z(z({release: o, environment: i}, n && {user: n.getUser()}), a && {userAgent: a}), e));
            if (n) {
                var u = n.getSession && n.getSession();
                u && u.status === J.Ok && u.update({status: J.Exited}), this.endSession(), n.setSession(c)
            }
            return c
        }, e.prototype._sendSessionUpdate = function () {
            var e = this.getStackTop(), t = e.scope, n = e.client;
            if (t) {
                var r = t.getSession && t.getSession();
                r && n && n.captureSession && n.captureSession(r)
            }
        }, e.prototype._invokeClient = function (e) {
            for (var t, n = [], r = 1; r < arguments.length; r++) n[r - 1] = arguments[r];
            var s = this.getStackTop(), o = s.scope, i = s.client;
            i && i[e] && (t = i)[e].apply(t, Z(n, [o]))
        }, e.prototype._callExtensionMethod = function (e) {
            for (var t = [], n = 1; n < arguments.length; n++) t[n - 1] = arguments[n];
            var r = Ye(), s = r.__SENTRY__;
            if (s && s.extensions && "function" == typeof s.extensions[e]) return s.extensions[e].apply(this, t);
            pe.warn("Extension method " + e + " couldn't be found, doing nothing.")
        }, e
    }();

    function Ye() {
        var e = re();
        return e.__SENTRY__ = e.__SENTRY__ || {extensions: {}, hub: void 0}, e
    }

    function Je(e) {
        var t = Ye(), n = Ke(t);
        return ze(t, e), n
    }

    function We() {
        var e = Ye();
        return Be(e) && !Ke(e).isOlderThan(4) || ze(e, new $e), ee() ? function (e) {
            var t, n, r;
            try {
                var s = null === (r = null === (n = null === (t = Ye().__SENTRY__) || void 0 === t ? void 0 : t.extensions) || void 0 === n ? void 0 : n.domain) || void 0 === r ? void 0 : r.active;
                if (!s) return Ke(e);
                if (!Be(s) || Ke(s).isOlderThan(4)) {
                    var o = Ke(e).getStackTop();
                    ze(s, new $e(o.client, Ve.clone(o.scope)))
                }
                return Ke(s)
            } catch (t) {
                return Ke(e)
            }
        }(e) : Ke(e)
    }

    function Be(e) {
        return !!(e && e.__SENTRY__ && e.__SENTRY__.hub)
    }

    function Ke(e) {
        return e && e.__SENTRY__ && e.__SENTRY__.hub || (e.__SENTRY__ = e.__SENTRY__ || {}, e.__SENTRY__.hub = new $e), e.__SENTRY__.hub
    }

    function ze(e, t) {
        return !!e && (e.__SENTRY__ = e.__SENTRY__ || {}, e.__SENTRY__.hub = t, !0)
    }

    function Qe() {
        var e, t, n,
            r = (void 0 === e && (e = We()), null === (n = null === (t = e) || void 0 === t ? void 0 : t.getScope()) || void 0 === n ? void 0 : n.getTransaction());
        r && (pe.log("[Tracing] Transaction: " + Ge.InternalError + " -> Global error occured"), r.setStatus(Ge.InternalError))
    }

    !function (e) {
        e.Ok = "ok", e.DeadlineExceeded = "deadline_exceeded", e.Unauthenticated = "unauthenticated", e.PermissionDenied = "permission_denied", e.NotFound = "not_found", e.ResourceExhausted = "resource_exhausted", e.InvalidArgument = "invalid_argument", e.Unimplemented = "unimplemented", e.Unavailable = "unavailable", e.InternalError = "internal_error", e.UnknownError = "unknown_error", e.Cancelled = "cancelled", e.AlreadyExists = "already_exists", e.FailedPrecondition = "failed_precondition", e.Aborted = "aborted", e.OutOfRange = "out_of_range", e.DataLoss = "data_loss"
    }(Ge || (Ge = {})), function (e) {
        e.fromHttpCode = function (t) {
            if (t < 400 && t >= 100) return e.Ok;
            if (t >= 400 && t < 500) switch (t) {
                case 401:
                    return e.Unauthenticated;
                case 403:
                    return e.PermissionDenied;
                case 404:
                    return e.NotFound;
                case 409:
                    return e.AlreadyExists;
                case 413:
                    return e.FailedPrecondition;
                case 429:
                    return e.ResourceExhausted;
                default:
                    return e.InvalidArgument
            }
            if (t >= 500 && t < 600) switch (t) {
                case 501:
                    return e.Unimplemented;
                case 503:
                    return e.Unavailable;
                case 504:
                    return e.DeadlineExceeded;
                default:
                    return e.InternalError
            }
            return e.UnknownError
        }
    }(Ge || (Ge = {}));
    var Xe, Ze = function () {
        function e(e) {
            void 0 === e && (e = 1e3), this.spans = [], this._maxlen = e
        }

        return e.prototype.add = function (e) {
            this.spans.length > this._maxlen ? e.spanRecorder = void 0 : this.spans.push(e)
        }, e
    }(), et = function (e) {
        function t(t, n) {
            var r = e.call(this, t) || this;
            return r._measurements = {}, r._hub = We(), ae(n, $e) && (r._hub = n), r.name = t.name || "", r.metadata = t.metadata || {}, r._trimEnd = t.trimEnd, r.transaction = r, r
        }

        return function (e, t) {
            function n() {
                this.constructor = e
            }

            Y(e, t), e.prototype = null === t ? Object.create(t) : (n.prototype = t.prototype, new n)
        }(t, e), t.prototype.setName = function (e) {
            this.name = e
        }, t.prototype.initSpanRecorder = function (e) {
            void 0 === e && (e = 1e3), this.spanRecorder || (this.spanRecorder = new Ze(e)), this.spanRecorder.add(this)
        }, t.prototype.setMeasurements = function (e) {
            this._measurements = z({}, e)
        }, t.prototype.setMetadata = function (e) {
            this.metadata = z(z({}, this.metadata), e)
        }, t.prototype.finish = function (t) {
            var n, r, s, o, i, a = this;
            if (void 0 === this.endTimestamp) {
                if (this.name || (pe.warn("Transaction has no name, falling back to `<unlabeled transaction>`."), this.name = "<unlabeled transaction>"), e.prototype.finish.call(this, t), !0 !== this.sampled) return pe.log("[Tracing] Discarding transaction because its trace was not chosen to be sampled."), void (null === (i = null === (s = null === (n = this._hub.getClient()) || void 0 === n ? void 0 : (r = n).getTransport) || void 0 === s ? void 0 : (o = s.call(r)).recordLostEvent) || void 0 === i || i.call(o, K.SampleRate, "transaction"));
                var c = this.spanRecorder ? this.spanRecorder.spans.filter((function (e) {
                    return e !== a && e.endTimestamp
                })) : [];
                this._trimEnd && c.length > 0 && (this.endTimestamp = c.reduce((function (e, t) {
                    return e.endTimestamp && t.endTimestamp ? e.endTimestamp > t.endTimestamp ? e : t : e
                })).endTimestamp);
                var u = {
                    contexts: {trace: this.getTraceContext()},
                    spans: c,
                    start_timestamp: this.startTimestamp,
                    tags: this.tags,
                    timestamp: this.endTimestamp,
                    transaction: this.name,
                    type: "transaction",
                    debug_meta: this.metadata
                };
                return Object.keys(this._measurements).length > 0 && (pe.log("[Measurements] Adding measurements to transaction", JSON.stringify(this._measurements, void 0, 2)), u.measurements = this._measurements), pe.log("[Tracing] Finishing " + this.op + " transaction: " + this.name + "."), this._hub.captureEvent(u)
            }
        }, t.prototype.toContext = function () {
            var t = e.prototype.toContext.call(this);
            return _e(z(z({}, t), {name: this.name, trimEnd: this._trimEnd}))
        }, t.prototype.updateWithContext = function (t) {
            var n;
            return e.prototype.updateWithContext.call(this, t), this.name = null != (n = t.name) ? n : "", this._trimEnd = t.trimEnd, this
        }, t
    }(function () {
        function e(e) {
            if (this.traceId = Pe(), this.spanId = Pe().substring(16), this.startTimestamp = Fe(), this.tags = {}, this.data = {}, !e) return this;
            e.traceId && (this.traceId = e.traceId), e.spanId && (this.spanId = e.spanId), e.parentSpanId && (this.parentSpanId = e.parentSpanId), "sampled" in e && (this.sampled = e.sampled), e.op && (this.op = e.op), e.description && (this.description = e.description), e.data && (this.data = e.data), e.tags && (this.tags = e.tags), e.status && (this.status = e.status), e.startTimestamp && (this.startTimestamp = e.startTimestamp), e.endTimestamp && (this.endTimestamp = e.endTimestamp)
        }

        return e.prototype.child = function (e) {
            return this.startChild(e)
        }, e.prototype.startChild = function (t) {
            var n = new e(z(z({}, t), {parentSpanId: this.spanId, sampled: this.sampled, traceId: this.traceId}));
            return n.spanRecorder = this.spanRecorder, n.spanRecorder && n.spanRecorder.add(n), n.transaction = this.transaction, n
        }, e.prototype.setTag = function (e, t) {
            var n;
            return this.tags = z(z({}, this.tags), ((n = {})[e] = t, n)), this
        }, e.prototype.setData = function (e, t) {
            var n;
            return this.data = z(z({}, this.data), ((n = {})[e] = t, n)), this
        }, e.prototype.setStatus = function (e) {
            return this.status = e, this
        }, e.prototype.setHttpStatus = function (e) {
            this.setTag("http.status_code", String(e));
            var t = Ge.fromHttpCode(e);
            return t !== Ge.UnknownError && this.setStatus(t), this
        }, e.prototype.isSuccess = function () {
            return this.status === Ge.Ok
        }, e.prototype.finish = function (e) {
            this.endTimestamp = "number" == typeof e ? e : Fe()
        }, e.prototype.toTraceparent = function () {
            var e = "";
            return void 0 !== this.sampled && (e = this.sampled ? "-1" : "-0"), this.traceId + "-" + this.spanId + e
        }, e.prototype.toContext = function () {
            return _e({
                data: this.data,
                description: this.description,
                endTimestamp: this.endTimestamp,
                op: this.op,
                parentSpanId: this.parentSpanId,
                sampled: this.sampled,
                spanId: this.spanId,
                startTimestamp: this.startTimestamp,
                status: this.status,
                tags: this.tags,
                traceId: this.traceId
            })
        }, e.prototype.updateWithContext = function (e) {
            var t, n, r, s, o;
            return this.data = null != (t = e.data) ? t : {}, this.description = e.description, this.endTimestamp = e.endTimestamp, this.op = e.op, this.parentSpanId = e.parentSpanId, this.sampled = e.sampled, this.spanId = null != (n = e.spanId) ? n : this.spanId, this.startTimestamp = null != (r = e.startTimestamp) ? r : this.startTimestamp, this.status = e.status, this.tags = null != (s = e.tags) ? s : {}, this.traceId = null != (o = e.traceId) ? o : this.traceId, this
        }, e.prototype.getTraceContext = function () {
            return _e({
                data: Object.keys(this.data).length > 0 ? this.data : void 0,
                description: this.description,
                op: this.op,
                parent_span_id: this.parentSpanId,
                span_id: this.spanId,
                status: this.status,
                tags: Object.keys(this.tags).length > 0 ? this.tags : void 0,
                trace_id: this.traceId
            })
        }, e.prototype.toJSON = function () {
            return _e({
                data: Object.keys(this.data).length > 0 ? this.data : void 0,
                description: this.description,
                op: this.op,
                parent_span_id: this.parentSpanId,
                span_id: this.spanId,
                start_timestamp: this.startTimestamp,
                status: this.status,
                tags: Object.keys(this.tags).length > 0 ? this.tags : void 0,
                timestamp: this.endTimestamp,
                trace_id: this.traceId
            })
        }, e
    }());

    function tt() {
        var e = this.getScope();
        if (e) {
            var t = e.getSpan();
            if (t) return {"sentry-trace": t.toTraceparent()}
        }
        return {}
    }

    function nt(e, t, n) {
        return function (e) {
            var t;
            return void 0 === e && (e = null === (t = We().getClient()) || void 0 === t ? void 0 : t.getOptions()), !!e && ("tracesSampleRate" in e || "tracesSampler" in e)
        }(t) ? void 0 !== e.sampled ? (e.setMetadata({transactionSampling: {method: B.Explicit}}), e) : ("function" == typeof t.tracesSampler ? (r = t.tracesSampler(n), e.setMetadata({
            transactionSampling: {
                method: B.Sampler,
                rate: Number(r)
            }
        })) : void 0 !== n.parentSampled ? (r = n.parentSampled, e.setMetadata({transactionSampling: {method: B.Inheritance}})) : (r = t.tracesSampleRate, e.setMetadata({
            transactionSampling: {
                method: B.Rate,
                rate: Number(r)
            }
        })), function (e) {
            if (isNaN(e) || "number" != typeof e && "boolean" != typeof e) return pe.warn("[Tracing] Given sample rate is invalid. Sample rate must be a boolean or a number between 0 and 1. Got " + JSON.stringify(e) + " of type " + JSON.stringify(typeof e) + "."), !1;
            if (e < 0 || e > 1) return pe.warn("[Tracing] Given sample rate is invalid. Sample rate must be between 0 and 1. Got " + e + "."), !1;
            return !0
        }(r) ? r ? (e.sampled = Math.random() < r, e.sampled ? (pe.log("[Tracing] starting " + e.op + " transaction - " + e.name), e) : (pe.log("[Tracing] Discarding transaction because it's not included in the random sample (sampling rate = " + Number(r) + ")"), e)) : (pe.log("[Tracing] Discarding transaction because " + ("function" == typeof t.tracesSampler ? "tracesSampler returned 0 or false" : "a negative sampling decision was inherited or tracesSampleRate is set to 0")), e.sampled = !1, e) : (pe.warn("[Tracing] Discarding transaction because of invalid sample rate."), e.sampled = !1, e)) : (e.sampled = !1, e);
        var r
    }

    function rt(e, t) {
        var n, r, s = (null === (n = this.getClient()) || void 0 === n ? void 0 : n.getOptions()) || {},
            o = new et(e, this);
        return (o = nt(o, s, z({
            parentSampled: e.parentSampled,
            transactionContext: e
        }, t))).sampled && o.initSpanRecorder(null === (r = s._experiments) || void 0 === r ? void 0 : r.maxSpans), o
    }

    function st() {
        var e = Ye();
        if (e.__SENTRY__) {
            var t = {
                mongodb: function () {
                    return new (te(module, "./integrations/node/mongo").Mongo)
                }, mongoose: function () {
                    return new (te(module, "./integrations/node/mongo").Mongo)({mongoose: !0})
                }, mysql: function () {
                    return new (te(module, "./integrations/node/mysql").Mysql)
                }, pg: function () {
                    return new (te(module, "./integrations/node/postgres").Postgres)
                }
            }, n = Object.keys(t).filter((function (e) {
                return !!function (e) {
                    var t;
                    try {
                        t = te(module, e)
                    } catch (e) {
                    }
                    try {
                        var n = te(module, "process").cwd;
                        t = te(module, n() + "/node_modules/" + e)
                    } catch (e) {
                    }
                    return t
                }(e)
            })).map((function (e) {
                try {
                    return t[e]()
                } catch (e) {
                    return
                }
            })).filter((function (e) {
                return e
            }));
            n.length > 0 && (e.__SENTRY__.integrations = Z(e.__SENTRY__.integrations || [], n))
        }
    }

    (Xe = Ye()).__SENTRY__ && (Xe.__SENTRY__.extensions = Xe.__SENTRY__.extensions || {}, Xe.__SENTRY__.extensions.startTransaction || (Xe.__SENTRY__.extensions.startTransaction = rt), Xe.__SENTRY__.extensions.traceHeaders || (Xe.__SENTRY__.extensions.traceHeaders = tt)), ee() && st(), be({
        callback: Qe,
        type: "error"
    }), be({callback: Qe, type: "unhandledrejection"});
    const ot = "ACCESS_DENIED", it = "INVALID_REQUEST", at = "INVALID_GRANT", ct = "GET_TOKEN_ERROR",
        ut = "INVALID_NONCE";

    class dt extends Error {
        constructor(e, t) {
            super(e), this.name = ot, this.message = e, this.httpStatusCode = t
        }
    }

    class lt extends Error {
        constructor(e, t) {
            super(e), this.name = it, this.message = e, this.httpStatusCode = t
        }
    }

    class pt extends Error {
        constructor(e, t) {
            super(e), this.name = at, this.message = e, this.httpStatusCode = t
        }
    }

    class ht extends Error {
        constructor(e, t) {
            super(e), this.name = ct, this.message = e, this.httpStatusCode = t
        }
    }

    class gt extends Error {
        constructor(e) {
            super(e), this.name = ut, this.message = e
        }
    }

    function ft(e) {
        chrome.runtime.sendMessage({action: "dlTrackError", errorMessage: e})
    }

    const _t = "PROTOCOL_ERROR", mt = "NETWORK_ERROR", yt = "FULL_PAGE_SERVER_ERROR", vt = "INSTALLATION_ERROR",
        St = "TRANSLATED_INPUT_NOT_SET", Tt = "HTML_NODE_CHAR_LENGTH_TOO_LONG";

    class wt extends Error {
        constructor(e) {
            super(e), this.name = _t, this.message = e
        }
    }

    class Et extends Error {
        constructor(e) {
            super(e), this.name = mt, this.message = e
        }
    }

    class bt extends Error {
        constructor(e) {
            super(e), this.name = yt, this.message = e
        }
    }

    class xt extends Error {
        constructor(e, t) {
            super(t), this.name = e, this.message = t
        }
    }

    function It(e, t) {
        switch (console.error(e, t), e) {
            case _t:
            case mt:
                return "error_message_generic_error";
            case"429":
                return "error_message_too_many_requests";
            case yt:
                return "full_page_translation_error_service_not_available";
            case Tt:
                return "full_page_translation_error_exceeds_char_limit";
            case ot:
            case it:
            case at:
                return "login_error_default_message";
            case"-32003":
                return "error_message_generic_error";
            case ct:
            case ut:
                return "login_error_default_message";
            default:
                return "error_message_generic_error"
        }
    }

    function Rt(e) {
        try {
            ft(`${e.name} ${e.message}`)
        } catch (e) {
            console.error(e)
        }
        return {error: It(e.name, e.message)}
    }

    function Lt() {
        return {url: "https://s.deepl.com/chrome/statistics"}
    }

    function kt(e) {
        switch (e) {
            case E:
                return 1;
            case w:
                return 2;
            case T:
                return 3;
            case b:
                return 4;
            default:
                return 0
        }
    }

    function Ot(e) {
        switch (e) {
            case R:
                return O;
            case L:
                return C;
            default:
                return 0
        }
    }

    function Ct(e) {
        switch (e) {
            case R:
                return N;
            case L:
                return P;
            case k:
                return A;
            default:
                return 0
        }
    }

    function Nt(e) {
        return e === D ? 1 : 0
    }

    function Pt(e) {
        switch (e) {
            case M:
                return 1;
            case j:
                return 2;
            default:
                return 0
        }
    }

    function At(e) {
        return e === U ? 1 : 0
    }

    function Dt(e) {
        return e === q ? 60011 : 0
    }

    function jt() {
        return "0.15.0".split(".").length < 4
    }

    async function Mt(e) {
        let t = null;
        const n = await m();
        return jt() && (t = e.domainName), fetch("https://s.deepl.com/chrome/statistics", {
            method: "POST",
            headers: {Accept: "application/json", "Content-Type": "application/json"},
            body: JSON.stringify({
                eventId: 60001,
                sessionId: e.sessionId,
                extensionVersion: "0.15.0",
                userInfos: {userType: n ? I : x},
                translationRequestData: {
                    extensionVersion: "0.15.0", ...t && {domainName: t},
                    userId: e.installationId,
                    translationTrigger: kt(e.trigger),
                    translationData: {
                        sourceLang: e.websiteLanguage,
                        targetLang: e.targetLanguage,
                        sourceLength: e.characterCount
                    }, ...e.translationPerformance && {translationPerformance: e.translationPerformance}
                }
            })
        }).then((e => {
            e.ok || console.log(`DAP 60001 request failed: ${e}`)
        }))
    }

    async function Ut(e) {
        const t = await m();
        return fetch("https://s.deepl.com/chrome/statistics", {
            method: "POST",
            headers: {Accept: "application/json", "Content-Type": "application/json"},
            body: JSON.stringify({
                eventId: 60002,
                sessionId: e.sessionId,
                extensionVersion: "0.15.0",
                userInfos: {userType: t ? I : x},
                errorData: {extensionVersion: "0.15.0", errorCode: e.errorCode, errorMessage: e.errorMessage}
            })
        }).then((e => {
            e.ok || console.log(`DAP 60002 request failed: ${e}`)
        }))
    }

    async function qt(e) {
        const t = await m();
        return fetch("https://s.deepl.com/chrome/statistics", {
            method: "POST",
            headers: {Accept: "application/json", "Content-Type": "application/json"},
            body: JSON.stringify({
                eventId: 60003,
                sessionId: e.sessionId,
                extensionVersion: "0.15.0",
                userInfos: {userType: t ? I : x},
                translationErrorData: {
                    errorData: {
                        extensionVersion: "0.15.0",
                        errorCode: e.errorCode,
                        errorMessage: e.errorMessage
                    }, targetLang: e.targetLang
                }
            })
        }).then((e => {
            e.ok || console.log(`DAP 60003 request failed: ${e}`)
        }))
    }

    async function Ft(e) {
        const t = await m();
        return fetch("https://s.deepl.com/chrome/statistics", {
            method: "POST",
            headers: {Accept: "application/json", "Content-Type": "application/json"},
            body: JSON.stringify({
                eventId: 60004,
                sessionId: e.sessionId,
                extensionVersion: "0.15.0",
                userInfos: {userType: t ? I : x},
                contentScriptErrorData: {
                    errorData: {
                        extensionVersion: "0.15.0",
                        errorCode: e.errorCode,
                        errorMessage: e.errorMessage
                    }, domainName: e.domainName
                }
            })
        }).then((e => {
            e.ok || console.log(`DAP 60004 request failed: ${e}`)
        }))
    }

    async function Vt(e) {
        if (!jt()) return;
        const t = await m();
        return fetch("https://s.deepl.com/chrome/statistics", {
            method: "POST",
            headers: {Accept: "application/json", "Content-Type": "application/json"},
            body: JSON.stringify({
                eventId: 60005,
                sessionId: e.sessionId,
                extensionVersion: "0.15.0",
                userInfos: {userType: t ? I : x},
                inputInjectionData: {extensionVersion: "0.15.0", domainName: e.domainName}
            })
        }).then((e => {
            e.ok || console.log(`DAP 60005 request failed: ${e}`)
        }))
    }

    async function Gt(e) {
        if (!jt()) return;
        const t = await m();
        return fetch("https://s.deepl.com/chrome/statistics", {
            method: "POST",
            headers: {Accept: "application/json", "Content-Type": "application/json"},
            body: JSON.stringify({
                eventId: 60007,
                sessionId: e.sessionId,
                userInfos: {userType: t ? I : x},
                extensionVersion: "0.15.0",
                extensionOpenedData: {
                    extensionVersion: "0.15.0",
                    domainName: e.domainName,
                    extensionOpenedFrom: e.extensionOpenedFrom
                }
            })
        }).then((e => {
            e.ok || console.log(`DAP 60007 request failed: ${e}`)
        }))
    }

    async function Ht(e) {
        if (!jt()) return;
        const t = await m();
        return fetch("https://s.deepl.com/chrome/statistics", {
            method: "POST",
            headers: {Accept: "application/json", "Content-Type": "application/json"},
            body: JSON.stringify({
                eventId: 60008,
                sessionId: e.sessionId,
                userInfos: {userType: t ? I : x},
                extensionVersion: "0.15.0",
                logoClickedData: {extensionVersion: "0.15.0", logoClickedFrom: e.logoClickedFrom}
            })
        }).then((e => {
            e.ok || console.log(`DAP 60008 request failed: ${e}`)
        }))
    }

    async function $t(e) {
        const t = await m();
        return fetch("https://s.deepl.com/chrome/statistics", {
            method: "POST",
            headers: {Accept: "application/json", "Content-Type": "application/json"},
            body: JSON.stringify({
                eventId: 60009,
                sessionId: e.sessionId,
                extensionVersion: "0.15.0",
                userInfos: {userType: t ? I : x},
                notificationViewedData: {extensionVersion: "0.15.0", notificationType: Nt(e.notificationType)}
            })
        }).then((e => {
            e.ok || console.log(`DAP 60009 request failed: ${e}`)
        }))
    }

    async function Yt(e) {
        const t = await m();
        return fetch("https://s.deepl.com/chrome/statistics", {
            method: "POST",
            headers: {Accept: "application/json", "Content-Type": "application/json"},
            body: JSON.stringify({
                eventId: 60010,
                sessionId: e.sessionId,
                userInfos: {userType: t ? I : x},
                extensionVersion: "0.15.0",
                shortcutUsedData: {shortcutType: At(e.shortcutType), uiLocation: Pt(e.location), shortcut: e.shortcut}
            })
        }).then((e => {
            e.ok || console.log(`DAP 60010 request failed: ${e}`)
        }))
    }

    async function Jt(e) {
        const t = await m();
        return fetch("https://s.deepl.com/chrome/statistics", {
            method: "POST",
            headers: {Accept: "application/json", "Content-Type": "application/json"},
            body: JSON.stringify({
                eventId: Dt(e.uiAction),
                sessionId: e.sessionId,
                userInfos: {userType: t ? I : x},
                extensionVersion: "0.15.0"
            })
        }).then((e => {
            e.ok || console.log(`DAP 60011 request failed: ${e}`)
        }))
    }

    function Wt(e, t) {
        const n = e.split("?")[1];
        if (n) {
            var r = n.split("&");
            for (let e = 0; e < r.length; e++) {
                const n = r[e].split("=");
                if (decodeURIComponent(n[0]) === t) return decodeURIComponent(n[1])
            }
        }
    }

    async function Bt(e, t, n = "https://test.deepl.com/oidc-login") {
        const r = u(), s = u(), o = await async function (e) {
            const t = await crypto.subtle.digest("SHA-256", (new TextEncoder).encode(e));
            return btoa(String.fromCharCode(...new Uint8Array(t))).replace(/=/g, "").replace(/\+/g, "-").replace(/\//g, "_")
        }(s);
        return {
            url: n + "?" + new URLSearchParams([["client_id", e], ["redirect_uri", t], ["code_challenge", o], ["code_challenge_method", "S256"], ["nonce", r]]),
            nonce: r,
            codeVerifier: s
        }
    }

    function Kt(e, t) {
        if (200 !== e) {
            if ("invalid_request" === t.error) throw new lt(`dl-oidc (get tokens): invalid request - http status: ${e}`, e);
            if ("invalid_grant" === t.error) throw new pt(`dl-oidc (get tokens): invalid grant - http status: ${e}`, e);
            throw new ht(`dl-oidc (get tokens): http status: ${e}`, e)
        }
    }

    async function zt(e, t) {
        let n = t.access_token;

        async function r() {
            return await fetch(e, {...t, headers: {...t.headers, Authorization: n ? `Bearer ${n}` : "None"}})
        }

        let s = await r();
        return 401 !== s.status && 403 !== s.status || (n = await async function (e, t = "https://deepl.com/oidc/token") {
            const n = await fetch(t, {
                method: "POST",
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                body: new URLSearchParams({grant_type: "refresh_token", refresh_token: e})
            }), r = await n.json();
            return Kt(n.status, r), r.access_token
        }(t.refresh_token, t.token_url), s = await r()), {response: s, accessToken: n}
    }

    async function Qt(e, t, n, r) { //e is body, t is content-type, n=, r=
        const s = await _(), o = await y(),
            i = s && s.access_token && o ? "https://www.deepl.com/jsonrpc?client=chrome-extension,0.15.0" : "https://www2.deepl.com/jsonrpc?client=chrome-extension,0.15.0";
        try {
            const {response: o, accessToken: a} = await zt(i, {
                method: "POST",
                headers: {"Content-Type": t},
                body: e,
                access_token: s.access_token,
                refresh_token: s.refresh_token,
                token_url: "-"
            });
            await g(a, s.refresh_token, s.id_token);
            let c = function (e, t, n, r) {
                let s;
                try {
                    s = JSON.parse(e)
                } catch (e) {
                    throw 200 === t ? (qt({
                        sessionId: n,
                        errorCode: _t,
                        errorMessage: "JSON_RPC: Parsing of response object failed"
                    }), new wt("JSON_RPC: Parsing of response object failed")) : t >= 500 && r ? (qt({
                        sessionId: n,
                        errorCode: yt,
                        errorMessage: `Http status code ${t}`
                    }), new bt(`FPT failed: Http status code ${t}`)) : (qt({
                        sessionId: n,
                        errorCode: t.toString(),
                        errorMessage: `Http status code ${t}`
                    }), new Et(`Translation failed: Http status code ${t}`))
                }
                if ("2.0" !== s.jsonrpc || void 0 !== s.result && 200 !== t || void 0 === s.result && void 0 === s.error) throw qt({
                    sessionId: n,
                    errorCode: _t,
                    errorMessage: "JSON_RPC: missing field, wrong id, proper result without 200 status, or missing data"
                }), new wt("JSON_RPC: missing field, wrong id, proper result without 200 status, or missing data");
                if (s.error) {
                    const e = s.error;
                    if (void 0 === e.code) throw qt({
                        sessionId: n,
                        errorCode: _t,
                        errorMessage: 'JSON_RPC: Missing field "code" in error response'
                    }), new wt('JSON_RPC: Missing field "code" in error response');
                    throw qt({
                        sessionId: n,
                        errorCode: e.code.toString(),
                        errorMessage: e.message
                    }), new xt(e.code.toString(), e.message)
                }
                200 !== t && qt({sessionId: n, errorCode: t.toString(), errorMessage: `Http status code ${t}`});
                return s
            }(await o.text(), o.status, n, r);
            return {data: c, status: o.status}
        } catch (e) {
            e.name !== ot && e.name !== it && e.name !== at && e.name !== ct && e.name !== ut || await f();
            try {
                qt({sessionId: n, errorCode: e.name, errorMessage: e.message})
            } catch (e) {
            }
            throw e
        }
    }

    async function Xt(e, t, n = (e => e), r, s, o) {
        const i = function (e, t, n) {
            return {jsonrpc: "2.0", method: e, params: t, id: n}
        }(e, t, r), a = n(JSON.stringify(i), i);
        return (await Qt(a, "application/json; charset=utf-8", s, o)).data.result
    }

    async function Zt(e) {
        let t, n = {};
        for (const t in e.weights) e.weights[t] > 0 && (n[t] = e.weights[t]);
        t = e.isHtml ? e.sourceTexts : e.sourceTexts.map((e => e.slice(0, 4999)));
        const r = {
            texts: t.map((e => ({text: e}))), ...e.isHtml && {html: "enabled"}, ...!e.isHtml && {splitting: e.splitting ? e.splitting : "paragraphs"},
            lang: {
                target_lang: e.targetLangCode.toUpperCase(),
                source_lang_user_selected: "auto" === e.sourceLang ? "auto" : e.sourceLang.toUpperCase(),
                preference: {weight: n}
            }
        };
        let s = 1;
        const o = Date.now();
        t.forEach((e => s += (e.match(/[i]/g) || []).length)), r.timestamp = o - o % s + s;
        const i = await Xt("LMT_handle_texts", r, ((e, t) => e = (t.id + 3) % 13 == 0 || (t.id + 5) % 29 == 0 ? e.replace('"method":"', '"method" : "') : e.replace('"method":"', '"method": "')), e.id, e.sessionId, e.isHtml);
        return {
            texts: i.texts.map((e => e.text)),
            sourceLang: i.lang.toUpperCase(),
            detectedLanguages: i.detectedLanguages
        }
    }

    function en(e, t) {
        const n = 1e-5;
        return n * Math.floor((.97 * e + t) / n)
    }

    function tn() {
        let e = new Uint32Array(1);
        return crypto.getRandomValues(e), e[0]
    }

    function nn(e) {
        let t;
        try {
            t = new URL(e)
        } catch (e) {
            return !1
        }
        return "http:" === t.protocol || "https:" === t.protocol
    }

    const rn = {
        share_feedback_survey: "DOCS_GOOGLE_COM_FEEDBACK_FORM",
        try_DeepL_Pro: "DEEPL_COM_PRO",
        show_more_translation: "DEEPL_COM_TRANSLATOR",
        deepl_translator_home: "DEEPL_COM_HOMEPAGE",
        shortcut_settings: "SHORTCUT_SETTINGS"
    }, sn = {
        DOCS_GOOGLE_COM_FEEDBACK_FORM: "https://docs.google.com/forms/d/e/1FAIpQLSeDGytB1T8lGvMcOCn9_orOyMjILe-OQ2vAVWZlpid4avrw_w/viewform",
        DEEPL_COM_PRO: "https://www.deepl.com/pro",
        DEEPL_COM_TRANSLATOR: "https://www.deepl.com/translator",
        SHORTCUT_SETTINGS: "chrome://extensions/shortcuts"
    };

    function on(e) {
        const {destination: t, ...n} = e;
        if (void 0 !== t) switch (t) {
            case rn.share_feedback_survey:
                return sn.DOCS_GOOGLE_COM_FEEDBACK_FORM;
            case rn.try_DeepL_Pro:
                return n?.utmContent ? un(sn.DEEPL_COM_PRO, [["utm_source", "browser_extension"], ["utm_medium", "chrome"], ["utm_content", n?.utmContent]]) : sn.DEEPL_COM_PRO;
            case rn.show_more_translation:
                return n ? function ({sourceLang: e, targetLang: t, textToShare: n}) {
                    if (!e || !t || !n) return "";
                    const r = new URL(sn.DEEPL_COM_TRANSLATOR);
                    let s = `${e.toLowerCase()}/${t.toLowerCase()}/${encodeURIComponent(n)}`;
                    return r.hash = s, r.href
                }(n) : void console.error("Invalid action.");
            case rn.deepl_translator_home:
                return n?.utmContent ? un(sn.DEEPL_COM_TRANSLATOR, [["utm_source", "browser_extension"], ["utm_medium", "chrome"], ["utm_content", n?.utmContent]]) : sn.DEEPL_COM_TRANSLATOR;
            case rn.shortcut_settings:
                return sn.SHORTCUT_SETTINGS;
            default:
                return
        }
    }

    function an(e) {
        const t = on(e);
        t && chrome.tabs.create({url: cn(t), active: !0})
    }

    function cn(e) {
        return new URL(e).href
    }

    function un(e, t) {
        if (0 === t.length || !e) return;
        const n = new URL(e);
        return t.forEach((e => {
            2 === e.length && n.searchParams.append(e[0], e[1])
        })), n.href
    }

    const dn = {developer_tools: "DEVELOPER_TOOLS"}, ln = "dev-page.html";

    function pn({payload: e}) {
        const {destination: t} = e;
        if (void 0 === t) return;
        let n = "";
        if (t === dn.developer_tools) r = ln, n = chrome.runtime.getURL(r);
        var r;
        n && chrome.tabs.create({url: n, active: !0})
    }

    async function hn(e = !1) {
        const t = await l(["session"]), n = async () => {
            const n = await Bt("chromeExtension", chrome.identity.getRedirectURL(), "-");
            return new Promise(((r, s) => {
                chrome.identity.launchWebAuthFlow({
                    interactive: !0,
                    url: e ? `${n.url}#error_loginfailed` : n.url
                }, (async e => {
                    if (!e) {
                        const e = chrome.runtime.lastError.message;
                        return console.error("response url undefined", e), ft(`login - response url undefined. ${e}`), /The user did not approve access/gm.test(e) ? r() : s()
                    }
                    try {
                        const t = await async function (e, t, n, r, s, o = "https://deepl.com/oidc/token") {
                            const i = Wt(n, "code");
                            if (!i) {
                                const e = Wt(n, "error");
                                if ("access_denied" === e) throw new dt("dl-oidc (get tokens): access denied - undefined code");
                                if ("invalid_request" === e) throw new lt("dl-oidc (get tokens): invalid request - undefined code")
                            }
                            const a = await fetch(`${o}`, {
                                method: "POST",
                                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                                body: new URLSearchParams({
                                    client_id: e,
                                    code_verifier: s,
                                    code: i,
                                    grant_type: "authorization_code",
                                    redirect_uri: t
                                })
                            }), c = await a.json();
                            if (Kt(a.status, c), d(c.id_token).nonce !== r) throw new gt("dl-oidc (get tokens): invalid nonce");
                            return c
                        }("chromeExtension", chrome.identity.getRedirectURL(), e, n.nonce, n.codeVerifier, "-");
                        await g(t.access_token, t.refresh_token, t.id_token), r()
                    } catch (e) {
                        console.error(e);
                        try {
                            Ut({sessionId: t.session.id, errorCode: e.name, errorMessage: e.message})
                        } catch (e) {
                        }
                        s(e)
                    }
                }))
            }))
        };
        try {
            await n()
        } catch (e) {
            hn(!0)
        }
    }

    async function gn() {
        return await chrome.commands.getAll()
    }

    return chrome.runtime.onInstalled.addListener((e => {
        try {
            e.reason === chrome.runtime.OnInstalledReason.UPDATE && console.log(`Update from version ${e.previousVersion} to ${chrome.runtime.getManifest().version} was completed`), chrome.storage.sync.get(null, (function (e) {
                const t = {
                    settingSchemeVersion: 1.5,
                    blacklistDomains: [],
                    languagesForTranslation: v.reduce(((e, t) => ({...e, [t.k]: S})), {}),
                    selectedTargetLanguage: F(),
                    isTargetLanguageConfirmed: !1,
                    deepLProApiKey: {
                        IS_PROD: !0,
                        BASE_URLS: {
                            dapApi: "https://s.deepl.com/chrome/statistics",
                            toolingApiPro: "https://www.deepl.com/jsonrpc?client=chrome-extension,0.15.0",
                            toolingApi: "https://www2.deepl.com/jsonrpc?client=chrome-extension,0.15.0",
                            login: "-",
                            token: "-",
                            logout: "-",
                            feedbackSurvey: "https://docs.google.com/forms/d/e/1FAIpQLSeDGytB1T8lGvMcOCn9_orOyMjILe-OQ2vAVWZlpid4avrw_w/viewform",
                            deeplPro: "https://www.deepl.com/pro",
                            deeplTranslator: "https://www.deepl.com/translator"
                        },
                        APP_VERSION: "0.15.0",
                        FEATURE_FLAGS: {TRANSLATE_INPUT: !0, PAGE_TRANSLATION: !1, IMPROVE_WRITING: !1, PRO_LOGIN: !1}
                    }.DEEPL_PRO_API_KEY,
                    installationId: c(),
                    enableInputTranslation: !0,
                    languagePreferenceForWriting: V(),
                    languagePreferenceForReading: V()
                };
                if (!e || !e.settingSchemeVersion || e.settingSchemeVersion < t.settingSchemeVersion) {
                    let n = {};
                    for (const [r, s] of Object.entries(t)) n[r] = r in e && "settingSchemeVersion" !== r ? e[r] : s;
                    chrome.storage.sync.clear((() => {
                        chrome.storage.sync.set(n, (() => {
                            console.log("settings set: ", n)
                        }))
                    }))
                } else console.log(`settings with setting scheme version ${e.settingSchemeVersion} found`), chrome.storage.sync.get(null, (e => {
                    console.log("stored settings: ", e)
                }))
            })), chrome.contextMenus.removeAll((function () {
                const e = (() => {
                    const e = navigator.language || navigator.userLanguage;
                    switch (e) {
                        case"pt":
                        case"pt-PT":
                            return "pt_PT";
                        case"pt-BR":
                            return "pt_BR";
                        case"zh":
                        case"zh-CN":
                            return "zh_CN";
                        default:
                            return e.substr(0, 2)
                    }
                })();
                fetch(chrome.runtime.getURL(`/_locales/${e}/messages.json`)).then((e => e.json())).then((function (e) {
                    chrome.contextMenus.create({
                        title: e.context_menus_translate_section.message,
                        id: "dlTranslateInline",
                        contexts: ["selection"]
                    })
                })).catch((e => {
                    chrome.contextMenus.create({
                        title: "Translate this selection",
                        id: "dlTranslateInline",
                        contexts: ["selection"]
                    })
                })), $() && chrome.contextMenus.create({
                    title: "Translate DOM",
                    id: "dlTranslateDOMInline",
                    contexts: ["selection"]
                })
            })), chrome.tabs.query({}, (function (e) {
                for (var t = 0; t < e.length; ++t) nn(e[t].url) && chrome.scripting.executeScript({
                    target: {tabId: e[t].id},
                    files: ["build/content.js"]
                })
            }))
        } catch (e) {
            l(["session"]).then((t => {
                Ut({sessionId: t.session.id, errorCode: vt, errorMessage: e.message})
            })).catch((t => {
                Ut({sessionId: "undefined", errorCode: vt, errorMessage: e.message})
            }))
        }
    })), chrome.contextMenus.onClicked.addListener((function (e, t) {
        switch (e.menuItemId) {
            case"dlTranslateInline":
                chrome.tabs.query({active: !0, currentWindow: !0}, (function (e) {
                    chrome.storage.sync.get(null, (t => {
                        chrome.tabs.sendMessage(e[0].id, {
                            action: "dlTranslateInline",
                            targetLang: t.selectedTargetLanguage
                        })
                    }))
                }));
                break;
            case"dlTranslateDOMInline":
                chrome.tabs.query({active: !0, currentWindow: !0}, (function (e) {
                    chrome.storage.sync.get(null, (t => {
                        chrome.tabs.sendMessage(e[0].id, {
                            action: "dlTranslateDOMInline",
                            targetLang: t.selectedTargetLanguage
                        })
                    }))
                }))
        }
    })), chrome.tabs.onActivated.addListener((async function (e) {
        nn((await chrome.tabs.get(e.tabId)).url) && chrome.tabs.sendMessage(e.tabId, {action: "dlTabChanged"})
    })), chrome.runtime.onConnect.addListener((e => {
        $() && console.log(`port connected to tab with id: ${e.sender.tab.id}`)
    })), chrome.runtime.onMessage.addListener(((e, t, n) => {
        switch ($() && console.log(t.tab ? `message from a content script: ${t.tab.url} with tab id: ${t.tab.id}. action: ${e.action}` : `message from the extension with action: ${e.action}`), e.action) {
            case"dlRequestHtmlTranslation":
                (async function (e) {
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
                            trigger: e.payload.trigger,
                            websiteLanguage: e.payload.sourceLang,
                            targetLanguage: t.selectedTargetLanguage,
                            characterCount: e.payload.texts.reduce(((e, t) => e + t.length), 0),
                            translationPerformance: o
                        })
                    } catch (e) {
                        console.error(e)
                    }
                    return s.texts
                })(e).then((e => n(e))).catch((e => {
                    n(Rt(e))
                }));
                break;
            case"dlRequestInlineTranslation":
                (async function (e) {
                    await h();
                    const t = await l(["selectedTargetLanguage", "session", "installationId", "languagePreferenceForReading"]),
                        n = {
                            sourceLang: "auto",
                            targetLangCode: t.selectedTargetLanguage,
                            sessionId: t.session.id,
                            id: tn(),
                            sourceTexts: [e.payload.requests[0].text],
                            weights: t.languagePreferenceForReading,
                            splitting: "newlines"
                        }, r = await Zt(n);
                    try {
                        if (r.detectedLanguages) {
                            for (const e in r.detectedLanguages) (t.languagePreferenceForReading[e] || 0 === t.languagePreferenceForReading[e]) && (t.languagePreferenceForReading[e] = en(t.languagePreferenceForReading[e], r.detectedLanguages[e]));
                            chrome.storage.sync.set({languagePreferenceForReading: t.languagePreferenceForReading})
                        }
                    } catch (e) {
                        console.error(e)
                    }
                    try {
                        Mt({
                            sessionId: t.session.id,
                            installationId: t.installationId,
                            trigger: e.payload.trigger,
                            websiteLanguage: e.payload.sourceLang,
                            targetLanguage: t.selectedTargetLanguage,
                            characterCount: e.payload.requests.reduce(((e, t) => e + t.text.length), 0),
                            domainName: e.payload.domainName
                        })
                    } catch (e) {
                        console.error(e)
                    }
                    return [{text: r.texts[0], detected_source_language: r.sourceLang}]
                })(e).then((e => n(e))).catch((e => {
                    n(Rt(e))
                }));
                break;
            case"dlTriggerTranslateInline":
                !async function (e) {
                    const t = await l(["selectedTargetLanguage"]);
                    chrome.tabs.sendMessage(e, {action: "dlTranslateInline", targetLang: t.selectedTargetLanguage})
                }(t.tab.id);
                break;
            case"dlRequestInlineDOMTranslation":
                translateRequest(e).then((e => n(e))).catch((e => {
                    n(Rt(e))
                }));
                break;
            case"dlRequestInputTranslation":
                (async function (e) {
                    await h();
                    const t = await l(["targetLanguageUserInput", "session", "installationId", "languagePreferenceForWriting"]),
                        n = e.payload.requests, r = {
                            sourceLang: "auto",
                            targetLangCode: t.targetLanguageUserInput,
                            sessionId: t.session.id,
                            id: tn(),
                            sourceTexts: [n[0].text],
                            weights: t.languagePreferenceForWriting,
                            splitting: "paragraphs"
                        };
                    let s = await Zt(r);
                    try {
                        if (s.detectedLanguages) {
                            for (const e in s.detectedLanguages) (t.languagePreferenceForWriting[e] || 0 === t.languagePreferenceForWriting[e]) && (t.languagePreferenceForWriting[e] = en(t.languagePreferenceForWriting[e], s.detectedLanguages[e]));
                            chrome.storage.sync.set({languagePreferenceForWriting: t.languagePreferenceForWriting})
                        }
                    } catch (e) {
                        console.error(e)
                    }
                    try {
                        Mt({
                            sessionId: t.session.id,
                            installationId: t.installationId,
                            trigger: b,
                            websiteLanguage: s.sourceLang,
                            targetLanguage: t.targetLanguageUserInput,
                            characterCount: n.reduce(((e, t) => e + t.text.length), 0),
                            domainName: e.payload.domainName
                        })
                    } catch (e) {
                        console.error(e)
                    }
                    return [{detected_source_language: s.sourceLang, text: s.texts[0]}]
                })(e).then((e => n(e))).catch((e => {
                    n(Rt(e))
                }));
                break;
            case"dlRequestDetectLanguage":
                (async function (e) {
                    await h();
                    const t = await l(["selectedTargetLanguage", "session"]), n = {
                        sourceLang: "auto",
                        targetLangCode: t.selectedTargetLanguage,
                        sessionId: t.session.id,
                        id: tn(),
                        sourceTexts: [e.payload.sampleText]
                    };
                    return (await Zt(n)).sourceLang
                })(e).then((e => n(e))).catch((e => {
                    n(Rt(e)), ft(e.message)
                }));
                break;
            case"dlRequestLogin":
                hn().then((() => n(!0))).catch((e => {
                    n(Rt(e))
                }));
                break;
            case"dlRequestLogout":
                (async function () {
                    const e = await _(),
                        t = "-?" + new URLSearchParams([["id_token_hint", e.id_token], ["post_logout_redirect_uri", chrome.identity.getRedirectURL()]]);
                    chrome.identity.launchWebAuthFlow({url: t, interactive: !0}, (async () => {
                        await f()
                    }))
                })().then((() => n(!0))).catch((e => {
                    n(Rt(e))
                }));
                break;
            case"dlGetAppSettings":
                chrome.storage.sync.get(null, (function (e) {
                    n(e)
                }));
                break;
            case"dlUpdateAppSettings":
                chrome.storage.sync.set(e.payload, (function () {
                    n(e.payload)
                }));
                break;
            case"dlTranslatedInputTextNotSet":
                l(["session"]).then((t => {
                    Ft({
                        sessionId: t.session.id,
                        errorCode: St,
                        errorMessage: "The translated text could not be set",
                        domainName: e.payload
                    })
                })), n();
                break;
            case"dlGetUserState":
                (async function () {
                    return {isLoggedIn: await m(), isProUser: await y()}
                })().then((e => n(e))).catch((e => {
                    n(Rt(e))
                }));
                break;
            case"dlExternalURLRedirect":
                an(e.payload), n();
                break;
            case"dlOpenInternalPage":
                pn(e), n();
                break;
            case"dlTrackInputTranslationInjection":
                !async function (e) {
                    const t = await l(["session"]);
                    try {
                        Vt({sessionId: t.session.id, domainName: e})
                    } catch (e) {
                        console.error(e)
                    }
                }(e.payload), n();
                break;
            case"dlTrackExtensionOpenedEvent":
                !async function (e) {
                    const t = await l(["session"]);
                    try {
                        Gt({
                            sessionId: t.session.id,
                            domainName: e.domainName,
                            extensionOpenedFrom: Ot(e.extensionContext)
                        })
                    } catch (e) {
                        console.error(e)
                    }
                }(e.payload), n();
                break;
            case"dlTrackLogoClickedEvent":
                !async function (e) {
                    const t = await l(["session"]);
                    try {
                        Ht({sessionId: t.session.id, logoClickedFrom: Ct(e.extensionContext)})
                    } catch (e) {
                        console.error(e)
                    }
                }(e.payload), n();
                break;
            case"dlSendNotification":
                !async function (e) {
                    const t = await l(["notificationsViewed", "session"]);
                    if (e.type === D) {
                        const e = await chrome.storage.sync.get(["notificationsViewed"]);
                        if (e.notificationsViewed && e.notificationsViewed.hintKeyboardShortcut) return;
                        await chrome.storage.sync.set({
                            notificationsViewed: {
                                ...e.notificationsViewed,
                                hintKeyboardShortcut: !0
                            }
                        })
                    }
                    try {
                        $t({sessionId: t.session.id, notificationType: e.type})
                    } catch (e) {
                        console.error(e)
                    }
                    const n = await chrome.tabs.query({active: !0, currentWindow: !0});
                    chrome.tabs.sendMessage(n[0].id, {action: "dlShowNotification", payload: e})
                }(e.payload), n();
                break;
            case"dlCheckCommandShortcuts":
                gn().then((e => n(e)));
                break;
            case"dlTrackUseOfShortcut":
                !async function (e) {
                    const t = await l(["session"]), n = await gn();
                    try {
                        Yt({
                            sessionId: t.session.id,
                            shortcutType: e.shortcutType,
                            location: e.location,
                            shortcut: n.filter((t => t.name === e.shortcutType))[0].shortcut
                        })
                    } catch (e) {
                        console.error(e)
                    }
                }(e.payload), n();
                break;
            case"dlTrackUiAction":
                !async function (e) {
                    Jt({sessionId: (await l(["session"])).session.id, uiAction: e.uiAction})
                }(e.payload), n()
        }
        return !0
    })), chrome.commands.onCommand.addListener((e => {
        console.log(`Command "${e}" triggered`), e === U && chrome.tabs.query({
            active: !0,
            currentWindow: !0
        }, (function (e) {
            chrome.tabs.sendMessage(e[0].id, {action: "dlTriggerTranslationShortcut"})
        }))
    })), self.addEventListener("error", (function (e) {
        ft(e.message)
    })), e.INTERNAL_PAGE_KEYS = dn, e.URL_ACTIONS = rn, e.deleteAuthTokens = f, e.getAuthTokens = _, e.getDapApiConfig = Lt, e.getDapApiExtensionOpenedFromId = Ot, e.getDapApiLogoClickedFromId = Ct, e.getDapApiNotificationType = Nt, e.getDapApiShortcutTypes = At, e.getDapApiTranslationTriggerId = kt, e.getDapApiUiActionEventId = Dt, e.getDapApiUiLocation = Pt, e.getExternalUrl = on, e.getSettings = l, e.isProUser = y, e.isUserLoggedIn = m, e.jsonrpc_callFunction = Xt, e.openInternalPage = pn, e.redirectToExternalURL = an, e.sendDAPContentScriptErrorData = Ft, e.sendDAPErrorData = Ut, e.sendDAPExtensionOpened = Gt, e.sendDAPFTPBlockSizeExceeded = async function (e) {
        const t = await m();
        return fetch("https://s.deepl.com/chrome/statistics", {
            method: "POST",
            headers: {Accept: "application/json", "Content-Type": "application/json"},
            body: JSON.stringify({
                eventId: 60006,
                sessionId: e.sessionId,
                extensionVersion: "0.15.0",
                userInfos: {userType: t ? I : x},
                fptBlockSizeExceededData: {
                    translationData: {
                        sourceLang: e.websiteLanguage,
                        targetLang: e.targetLanguage,
                        sourceLength: e.characterCount
                    }, domainName: e.domainName
                }
            })
        }).then((e => {
            e.ok || console.log(`DAP 60006 request failed: ${e}`)
        }))
    }, e.sendDAPInputTranslationInjection = Vt, e.sendDAPLogoClicked = Ht, e.sendDAPNotificationViewed = $t, e.sendDAPShortcutUsed = Yt, e.sendDAPTranslationErrorData = qt, e.sendDAPTranslationRequest = Mt, e.sendDAPUiAction = Jt, e.setAuthTokens = g, e.translate = Zt, e.updateSession = h, Object.defineProperty(e, "__esModule", {value: !0}), e
}({});
//# sourceMappingURL=background.js.map
