export type Locale =
	| 'en'
	| 'ru'
	| 'de'
	| 'fr'
	| 'es'
	| 'zh'
	| 'ja'
	| 'ko'
	| 'th';

export const locales: { code: Locale; label: string; region: string }[] = [
	{ code: 'en', label: 'English', region: 'EU' },
	{ code: 'ru', label: 'Русский', region: 'EU' },
	{ code: 'de', label: 'Deutsch', region: 'EU' },
	{ code: 'fr', label: 'Français', region: 'EU' },
	{ code: 'es', label: 'Español', region: 'EU' },
	{ code: 'zh', label: '中文', region: 'Asia' },
	{ code: 'ja', label: '日本語', region: 'Asia' },
	{ code: 'ko', label: '한국어', region: 'Asia' },
	{ code: 'th', label: 'ไทย', region: 'Asia' }
];

type Messages = Record<string, string>;

const en: Messages = {
	'app.name': 'EquiDrug',
	'app.tagline': 'Your stuff, anywhere',
	'nav.lookup': 'Lookup',
	'nav.locker': 'Locker',
	'nav.planner': 'Planner',
	'nav.diet': 'Diet',
	'nav.wiki': 'Wiki',
	'nav.avoid': 'Avoid',
	'lookup.title': 'Find local equivalent',
	'lookup.hint': 'Enter product name, URL, or snap a photo',
	'lookup.query': 'Product name or ingredient',
	'lookup.country': 'Destination country',
	'lookup.city': 'City (optional)',
	'lookup.search': 'Find equivalent',
	'lookup.example': 'Try: Allegra, fexofenadine, Zyrtec',
	'locker.title': 'My locker',
	'locker.hint': 'Photo your cabinet or add items manually',
	'locker.scan': 'Scan shelf photo',
	'locker.add': 'Add item',
	'planner.title': 'Trip planner',
	'planner.hint': 'Plan quantities for your destination',
	'planner.convert': 'Convert for trip',
	'wiki.title': 'My equivalents wiki',
	'wiki.empty': 'Lookups and trip conversions build your personal dictionary',
	'avoid.title': 'Avoid & diet traps',
	'avoid.hint': 'Allergens and naming differences across countries',
	'diet.title': 'Travel-safe diet',
	'diet.hint': 'Hit your macros abroad — snap a plate, menu line, or type a dish',
	'diet.targets': 'Daily targets',
	'diet.protein': 'Protein',
	'diet.carbs': 'Carbs',
	'diet.fat': 'Fat',
	'diet.remaining': 'left today',
	'diet.analyze': 'Analyze food',
	'diet.add': 'Add to today',
	'diet.photo': 'Photo of plate',
	'diet.menu': 'Menu line or description',
	'diet.today': "Today's log",
	'disclaimer':
		'Not medical advice. Always consult a healthcare professional before switching products.',
	'category.consume': 'Consume',
	'category.avoid': 'Avoid'
};

const translations: Partial<Record<Locale, Messages>> = {
	ru: {
		'app.name': 'EquiDrug',
		'app.tagline': 'Ваши препараты в любой стране',
		'nav.lookup': 'Поиск',
		'nav.locker': 'Аптечка',
		'nav.planner': 'Поездка',
		'nav.wiki': 'Моя база',
		'nav.avoid': 'Избегать',
		'lookup.title': 'Найти аналог',
		'lookup.search': 'Найти',
		'disclaimer': 'Не является медицинской консультацией.'
	},
	de: {
		'app.name': 'EquiDrug',
		'app.tagline': 'Deine Produkte überall',
		'nav.lookup': 'Suche',
		'nav.locker': 'Schrank',
		'nav.planner': 'Reise',
		'nav.wiki': 'Mein Wiki',
		'nav.avoid': 'Meiden',
		'lookup.title': 'Lokales Äquivalent finden',
		'lookup.search': 'Suchen',
		'disclaimer': 'Keine medizinische Beratung.'
	},
	ja: {
		'app.name': 'EquiDrug',
		'app.tagline': '旅先でも同じ成分を',
		'nav.lookup': '検索',
		'nav.locker': 'ロッカー',
		'nav.planner': '旅行',
		'nav.wiki': 'マイWiki',
		'nav.avoid': '回避',
		'lookup.title': '現地の同等品を探す',
		'lookup.search': '検索',
		'disclaimer': '医療アドバイスではありません。'
	}
};

export function t(locale: Locale, key: string): string {
	return translations[locale]?.[key] ?? en[key] ?? key;
}

export { en };
