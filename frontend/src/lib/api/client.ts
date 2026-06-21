const API_URL = import.meta.env.PUBLIC_API_URL ?? 'http://localhost:16125';

export type Category = 'consume' | 'avoid';

export interface LookupRequest {
	query: string;
	country_code: string;
	city?: string;
	category?: Category;
	image_url?: string;
	product_url?: string;
}

export interface Product {
	id: string;
	type: string;
	category: Category;
	brand_name: string;
	display_name: string;
	country_code: string;
	image_url?: string;
	retailer_hint?: string;
	price_hint?: string;
}

export interface LookupMatch {
	product: Product;
	equivalence: {
		confidence: number;
		notes?: string;
		source: string;
	};
}

export interface LookupResponse {
	query: string;
	country_code: string;
	matches: LookupMatch[];
	disclaimer: string;
}

export interface LockerItem {
	id?: string;
	custom_name: string;
	dose_per_day: number;
	dose_unit: string;
	frequency: string;
	category: Category;
	product_id?: string;
}

export interface TripPlan {
	title: string;
	origin_country: string;
	dest_country: string;
	dest_city?: string;
	start_date: string;
	end_date: string;
	spare_percent?: number;
	preferred_brands?: string[];
}

async function api<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`${API_URL}${path}`, {
		...init,
		headers: {
			'Content-Type': 'application/json',
			...init?.headers
		}
	});
	if (!res.ok) {
		const err = await res.json().catch(() => ({ error: res.statusText }));
		throw new Error(err.error ?? 'Request failed');
	}
	return res.json();
}

export function lookup(req: LookupRequest) {
	return api<LookupResponse>('/api/v1/lookup', {
		method: 'POST',
		body: JSON.stringify(req)
	});
}

export function listLocker() {
	return api<LockerItem[]>('/api/v1/locker');
}

export function addLockerItem(item: LockerItem) {
	return api<LockerItem>('/api/v1/locker', {
		method: 'POST',
		body: JSON.stringify(item)
	});
}

export function scanLocker(imageUrl: string) {
	return api<{ status: string; message: string }>('/api/v1/locker/scan', {
		method: 'POST',
		body: JSON.stringify({ image_url: imageUrl })
	});
}

export function createTrip(trip: TripPlan) {
	return api<TripPlan & { id: string }>('/api/v1/trips', {
		method: 'POST',
		body: JSON.stringify(trip)
	});
}

export function convertTrip(id: string) {
	return api<unknown>(`/api/v1/trips/${id}/convert`, { method: 'POST', body: '{}' });
}

export function listWiki() {
	return api<LookupMatch[]>('/api/v1/wiki');
}

export const countries = [
	{ code: 'US', name: 'United States' },
	{ code: 'DE', name: 'Germany' },
	{ code: 'EE', name: 'Estonia' },
	{ code: 'FR', name: 'France' },
	{ code: 'ES', name: 'Spain' },
	{ code: 'JP', name: 'Japan' },
	{ code: 'TH', name: 'Thailand' },
	{ code: 'KR', name: 'South Korea' },
	{ code: 'CN', name: 'China' }
];
