import i18next, { createInstance, type i18n as I18nextInstance } from "i18next"
import { initReactI18next } from "react-i18next"

import enUS from "./messages/en-US"
import zhCN from "./messages/zh-CN"

export type Locale = "en-US" | "zh-CN"
export type TFunction = (key: string, params?: Record<string, string | number>) => string

export const messages = {
  "en-US": enUS,
  "zh-CN": zhCN,
} as const

export const supportedLocales = Object.keys(messages) as Locale[]

const resources = {
  "en-US": { translation: enUS },
  "zh-CN": { translation: zhCN },
}

const i18nextOptions = {
  fallbackLng: "en-US",
  supportedLngs: supportedLocales,
  defaultNS: "translation",
  interpolation: {
    escapeValue: false,
    prefix: "{",
    suffix: "}",
  },
  react: {
    useSuspense: false,
  },
  resources,
  returnNull: false,
  initAsync: false,
} as const

export function normalizeLocale(value: string | undefined | null): Locale {
  return value === "zh-CN" ? "zh-CN" : "en-US"
}

export function createI18nextInstance(locale: Locale = "en-US"): I18nextInstance {
  const instance = createInstance()
  void instance.use(initReactI18next).init({
    ...i18nextOptions,
    lng: locale,
  })
  return instance
}

export function createT(locale: Locale): TFunction {
  const fixedT = createI18nextInstance(locale).getFixedT(locale)

  return (key, params) => {
    const value = fixedT(key, params)
    if (typeof value !== "string" || value === "") {
      return key
    }

    return value
  }
}

export const i18n = i18next
void i18n.use(initReactI18next).init({
  ...i18nextOptions,
  lng: "en-US",
})
