import { createI18n } from 'vue-i18n'
// @ts-ignore
import en from '@languages/en.json'
// @ts-ignore
import ar from '@languages/ar.json'

const messages = {
  en,
  ar
}

export const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages
})

export const setLanguage = (locale: string) => {
  i18n.global.locale.value = locale as 'en' | 'ar'
  const dir = locale === 'ar' ? 'rtl' : 'ltr'
  document.documentElement.setAttribute('dir', dir)
  document.documentElement.setAttribute('lang', locale)
}
