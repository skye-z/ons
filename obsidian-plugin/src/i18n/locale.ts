// come from https://github.com/joethei/obsidian-rss/blob/master/src/l10n/locale.ts
import en from "./en";
import zh from "./zh";

const locale = (window.moment) ? window.moment.locale() : "en";

const localeMap: { [k: string]: Partial<typeof en> } = {
    en,
    "zh-cn": zh,
};

const userLocale = localeMap[locale];

export default function t(str: keyof typeof en, ...inserts: string[]): string {
    let localeStr = (userLocale && userLocale[str]) ?? en[str];

    for (let i = 0; i < inserts.length; i++) {
        localeStr = localeStr.replace(`%${i + 1}`, inserts[i]);
    }

    return localeStr;
}