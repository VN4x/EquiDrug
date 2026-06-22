package service

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/equidrug/equidrug/internal/domain"
)

func RenderTripHTML(report domain.TripReport) string {
	var b strings.Builder
	b.WriteString(`<!DOCTYPE html><html lang="en"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>EquiDrug — `)
	b.WriteString(html.EscapeString(report.Trip.Title))
	b.WriteString(`</title>
<style>
*{box-sizing:border-box}body{font-family:system-ui,sans-serif;margin:0;padding:24px;background:#f4f7f6;color:#1a2b2b}
.wrap{max-width:800px;margin:0 auto}.hdr{background:#0d6e6e;color:#fff;padding:24px;border-radius:12px;margin-bottom:20px}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:12px;margin-bottom:20px}
.stat{background:#fff;padding:16px;border-radius:10px;text-align:center;box-shadow:0 2px 12px rgba(0,0,0,.06)}
.stat strong{display:block;font-size:1.5rem;color:#0d6e6e}
.row{background:#fff;border-radius:12px;padding:20px;margin-bottom:16px;box-shadow:0 2px 12px rgba(0,0,0,.06)}
.pair{display:grid;grid-template-columns:1fr auto 1fr;gap:16px;align-items:center}
.product{border:1px solid #d4e0de;border-radius:10px;padding:16px}
.product h3{margin:0 0 8px;font-size:1rem}.arrow{font-size:2rem;color:#0d6e6e;text-align:center}
.meta{font-size:.9rem;color:#5c6f6f;margin:4px 0}.badge{display:inline-block;background:#e8f2f1;color:#0d6e6e;padding:2px 8px;border-radius:999px;font-size:.75rem}
.shop{margin-top:12px;padding-top:12px;border-top:1px dashed #d4e0de}
.bought{text-decoration:line-through;opacity:.6}
.disclaimer{font-size:.8rem;color:#5c6f6f;border-left:3px solid #e8a838;padding-left:12px;margin-top:24px}
@media print{body{background:#fff}.hdr{-webkit-print-color-adjust:exact;print-color-adjust:exact}}
@media(max-width:600px){.pair{grid-template-columns:1fr}.arrow{transform:rotate(90deg)}.stats{grid-template-columns:repeat(2,1fr)}}
</style></head><body><div class="wrap">`)

	b.WriteString(`<div class="hdr"><h1 style="margin:0">`)
	b.WriteString(html.EscapeString(report.Trip.Title))
	b.WriteString(`</h1><p style="margin:8px 0 0;opacity:.9">`)
	b.WriteString(html.EscapeString(report.Trip.OriginCountry))
	b.WriteString(` → `)
	b.WriteString(html.EscapeString(report.Trip.DestCountry))
	if report.Trip.DestCity != "" {
		b.WriteString(` (`)
		b.WriteString(html.EscapeString(report.Trip.DestCity))
		b.WriteString(`)`)
	}
	b.WriteString(` · `)
	b.WriteString(report.Trip.StartDate.Format("Jan 2, 2006"))
	b.WriteString(` – `)
	b.WriteString(report.Trip.EndDate.Format("Jan 2, 2006"))
	b.WriteString(fmt.Sprintf(` · %d days + %.0f%% spare</p></div>`, report.Summary.TripDays, report.Trip.SparePercent))

	b.WriteString(`<div class="stats">`)
	for _, s := range []struct{ label, val string }{
		{"Items", fmt.Sprintf("%d", report.Summary.TotalItems)},
		{"Matched", fmt.Sprintf("%d", report.Summary.MatchedItems)},
		{"To find", fmt.Sprintf("%d", report.Summary.UnmatchedItems)},
		{"Bought", fmt.Sprintf("%d", report.Summary.BoughtItems)},
	} {
		b.WriteString(`<div class="stat"><strong>`)
		b.WriteString(s.val)
		b.WriteString(`</strong>`)
		b.WriteString(html.EscapeString(s.label))
		b.WriteString(`</div>`)
	}
	b.WriteString(`</div>`)

	for _, row := range report.Rows {
		cls := ""
		if row.LineItem.Bought {
			cls = ` bought`
		}
		b.WriteString(`<div class="row` + cls + `">`)
		b.WriteString(`<p><span class="badge">`)
		b.WriteString(html.EscapeString(row.LockerItem.CustomName))
		b.WriteString(`</span> · Need <strong>`)
		b.WriteString(fmt.Sprintf("%.1f %s", row.LineItem.QuantityNeeded, row.LineItem.QuantityUnit))
		b.WriteString(`</strong></p>`)

		b.WriteString(`<div class="pair">`)
		b.WriteString(productBox("Home", row.Origin, report.Trip.OriginCountry))
		b.WriteString(`<div class="arrow">→</div>`)
		b.WriteString(productBox("Buy locally", row.Foreign, report.Trip.DestCountry))
		b.WriteString(`</div>`)

		if row.Foreign != nil {
			b.WriteString(`<div class="shop">`)
			if row.LineItem.WhereToBuy != "" {
				b.WriteString(`<p class="meta">🛒 `)
				b.WriteString(html.EscapeString(row.LineItem.WhereToBuy))
				b.WriteString(`</p>`)
			}
			if row.LineItem.EstimatedPrice != "" {
				b.WriteString(`<p class="meta">💰 `)
				b.WriteString(html.EscapeString(row.LineItem.EstimatedPrice))
				b.WriteString(`</p>`)
			}
			if row.Confidence > 0 {
				b.WriteString(fmt.Sprintf(`<p class="meta">Confidence: %.0f%%`, row.Confidence*100))
				if row.MatchNotes != "" {
					b.WriteString(` — `)
					b.WriteString(html.EscapeString(row.MatchNotes))
				}
				b.WriteString(`</p>`)
			}
			b.WriteString(`</div>`)
		} else if row.MatchNotes != "" {
			b.WriteString(`<p class="meta">`)
			b.WriteString(html.EscapeString(row.MatchNotes))
			b.WriteString(`</p>`)
		}

		if row.LineItem.Notes != "" {
			b.WriteString(`<p class="meta">📝 `)
			b.WriteString(html.EscapeString(row.LineItem.Notes))
			b.WriteString(`</p>`)
		}
		if row.LineItem.Bought {
			b.WriteString(`<p class="meta">✅ Marked as bought</p>`)
		}
		b.WriteString(`</div>`)
	}

	b.WriteString(`<p class="disclaimer">`)
	b.WriteString(html.EscapeString(report.Disclaimer))
	b.WriteString(`</p>`)
	b.WriteString(`<p class="meta" style="text-align:center;margin-top:16px">Generated by EquiDrug · `)
	b.WriteString(time.Now().Format(time.RFC1123))
	b.WriteString(`</p></div></body></html>`)
	return b.String()
}

func productBox(label string, p *domain.Product, fallbackCountry string) string {
	var b strings.Builder
	b.WriteString(`<div class="product"><p class="meta" style="margin:0 0 8px">`)
	b.WriteString(html.EscapeString(label))
	b.WriteString(`</p>`)
	if p != nil {
		b.WriteString(`<h3>`)
		b.WriteString(html.EscapeString(p.DisplayName))
		b.WriteString(`</h3><p class="meta"><span class="badge">`)
		b.WriteString(html.EscapeString(p.BrandName))
		b.WriteString(`</span> · `)
		b.WriteString(html.EscapeString(p.CountryCode))
		b.WriteString(`</p>`)
	} else {
		b.WriteString(`<h3 style="color:#5c6f6f">—</h3><p class="meta">`)
		b.WriteString(html.EscapeString(fallbackCountry))
		b.WriteString(`</p>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}
