"use client"

import * as React from "react"
import { I18nextProvider, useTranslation } from "react-i18next"

import { createI18nextInstance, type Locale, type TFunction } from "."

const I18nContext = React.createContext<{
  locale: Locale
  t: TFunction
  setLocale?: (locale: Locale) => void
}>({
  locale: "en-US",
  t: (key) => key,
})

export function I18nProvider({
  children,
  locale = "en-US",
  setLocale,
}: {
  children: React.ReactNode
  locale?: Locale
  setLocale?: (locale: Locale) => void
}) {
  const [i18n] = React.useState(() => createI18nextInstance(locale))

  return (
    <I18nextProvider i18n={i18n}>
      <I18nContextBridge locale={locale} setLocale={setLocale}>
        {children}
      </I18nContextBridge>
    </I18nextProvider>
  )
}

function I18nContextBridge({
  children,
  locale,
  setLocale,
}: {
  children: React.ReactNode
  locale: Locale
  setLocale?: (locale: Locale) => void
}) {
  const { i18n, t: translate } = useTranslation()

  React.useEffect(() => {
    if (i18n.resolvedLanguage !== locale) {
      void i18n.changeLanguage(locale)
    }
  }, [i18n, locale])

  const t = React.useCallback<TFunction>(
    (key, params) => {
      const value = translate(key, params)
      return typeof value === "string" && value ? value : key
    },
    [translate]
  )

  return (
    <I18nContext.Provider value={{ locale, t, setLocale }}>
      {children}
    </I18nContext.Provider>
  )
}

export function useI18n() {
  return React.useContext(I18nContext)
}
