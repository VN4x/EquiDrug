const API_URL = import.meta.env.PUBLIC_API_URL ?? 'http://localhost:16125';

export type Category = 'consume' | 'avoid';

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

export interface LockerItem {
	id?: string;
	custom_name: string;
	dose_per_day: number;
	dose_unit: string;
	frequency: string;
	category: Category;
	product_id?: string;
}

export interface TripLineItem {
	id: string;
	trip_id: string;
	locker_item_id: string;
	origin_product_id?: string;
	foreign_product_id?: string;
	quantity_needed: number;
	quantity_unit: string;
	estimated_price?: string;
	where_to_buy?: string;
	bought: boolean;
	notes?: string;
}

export interface TripReportRow {
	line_item: TripLineItem;
	locker_item: LockerItem;
	origin?: Product;
	foreign?: Product;
	confidence: number;
	match_notes?: string;
	match_status: 'matched' | 'unmatched';
}

export interface TripReportSummary {
	total_items: number;
	matched_items: number;
	unmatched_items: number;
	bought_items: number;
	trip_days: number;
}

export interface TripReport {
	trip: {
		id: string;
		title: string;
		origin_country: string;
		dest_country: string;
		dest_city?: string;
		start_date: string;
		end_date: string;
		spare_percent: number;
		preferred_brands?: string[];
		status: string;
	};
	summary: TripReportSummary;
	rows: TripReportRow[];
	disclaimer: string;
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

export interface LookupRequest {
	query: string;
	country_code: string;
	city?: string;
	category?: Category;
}

export interface LookupMatch {
	product: Product;
	equivalence: { confidence: number; notes?: string; source: string };
}

export interface LookupResponse {
	query: string;
	country_code: string;
	matches: LookupMatch[];
	disclaimer: string;
}

async function api<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`${API_URL}${path}`, {
		...init,
		headers: { 'Content-Type': 'application/json', ...init?.headers }
	});
	if (!res.ok) {
		const err = await res.json().catch(() => ({ error: res.statusText }));
		throw new Error(err.error ?? 'Request failed');
	}
	return res.json();
}

export function lookup(req: LookupRequest) {
	return api<LookupResponse>('/api/v1/lookup', { method: 'POST', body: JSON.stringify(req) });
}

export function listLocker() {
	return api<LockerItem[]>('/api/v1/locker');
}

export function addLockerItem(item: LockerItem) {
	return api<LockerItem>('/api/v1/locker', { method: 'POST', body: JSON.stringify(item) });
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

export function getTrip(id: string) {
	return api<TripReport>(`/api/v1/trips/${id}`);
}

export function convertTrip(id: string, preferredBrands?: string[]) {
	return api<TripReport>(`/api/v1/trips/${id}/convert`, {
		method: 'POST',
		body: JSON.stringify({ preferred_brands: preferredBrands })
	});
}

export function updateTripLineItem(
	tripId: string,
	itemId: string,
	patch: { bought?: boolean; notes?: string }
) {
	return api<TripReport>(`/api/v1/trips/${tripId}/items/${itemId}`, {
		method: 'PATCH',
		body: JSON.stringify(patch)
	});
}

export function tripExportUrl(tripId: string, format: 'html' | 'json' = 'html') {
	return `${API_URL}/api/v1/trips/${tripId}/export?format=${format}`;
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
	{ code: 'VN', name: 'Vietnam' },
	{ code: 'LA', name: 'Laos' },
	{ code: 'TH', name: 'Thailand' },
	{ code: 'KR', name: 'South Korea' },
	{ code: 'CN', name: 'China' }
];

export function countryName(code: string) {
	return countries.find((c) => c.code === code)?.name ?? code;
}

export function formatDate(iso: string) {
	return new Date(iso).toLocaleDateString(undefined, {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	});
}
