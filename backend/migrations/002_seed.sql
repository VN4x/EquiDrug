-- Seed: demo user + Allegra / fexofenadine cross-country example
INSERT INTO users (id, email, display_name, locale)
VALUES ('00000000-0000-0000-0000-000000000001', 'demo@equidrug.local', 'Traveler Demo', 'en')
ON CONFLICT (id) DO NOTHING;

INSERT INTO active_ingredients (id, inn, name, strength, unit) VALUES
  ('a1000000-0000-0000-0000-000000000001', 'fexofenadine', 'Fexofenadine hydrochloride', 180, 'mg'),
  ('a1000000-0000-0000-0000-000000000002', 'bilastine', 'Bilastine', 20, 'mg'),
  ('a1000000-0000-0000-0000-000000000003', 'cetirizine', 'Cetirizine', 10, 'mg')
ON CONFLICT (id) DO NOTHING;

INSERT INTO products (id, type, category, brand_name, display_name, country_code, retailer_hint, price_hint) VALUES
  ('p1000000-0000-0000-0000-000000000001', 'drug', 'consume', 'Allegra', 'Allegra 180mg (fexofenadine)', 'US', 'Walmart, CVS, Walgreens', '~$15-25 / 30ct'),
  ('p1000000-0000-0000-0000-000000000002', 'drug', 'consume', 'Telfast', 'Telfast 180mg (fexofenadine)', 'DE', 'Apotheke, dm-drogerie', '~€8-15 / 20ct'),
  ('p1000000-0000-0000-0000-000000000003', 'drug', 'consume', 'Bilaxten', 'Bilaxten 20mg (bilastine)', 'DE', 'Apotheke', '~€10-18 / 20ct'),
  ('p1000000-0000-0000-0000-000000000004', 'drug', 'consume', 'Fexofast', 'Fexofast 180mg', 'EE', 'Benu, Apotheka', '~€6-12 / 10ct'),
  ('p1000000-0000-0000-0000-000000000005', 'drug', 'consume', 'Zyrtec', 'Zyrtec 10mg (cetirizine)', 'US', 'Costco, Target', '~$12-20 / 30ct'),
  ('p1000000-0000-0000-0000-000000000006', 'drug', 'consume', 'Cetirizin-ratiopharm', 'Cetirizin 10mg', 'DE', 'Apotheke', '~€3-8 / 20ct')
ON CONFLICT (id) DO NOTHING;

INSERT INTO product_ingredients (product_id, ingredient_id) VALUES
  ('p1000000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001'),
  ('p1000000-0000-0000-0000-000000000002', 'a1000000-0000-0000-0000-000000000001'),
  ('p1000000-0000-0000-0000-000000000003', 'a1000000-0000-0000-0000-000000000002'),
  ('p1000000-0000-0000-0000-000000000004', 'a1000000-0000-0000-0000-000000000001'),
  ('p1000000-0000-0000-0000-000000000005', 'a1000000-0000-0000-0000-000000000003'),
  ('p1000000-0000-0000-0000-000000000006', 'a1000000-0000-0000-0000-000000000003')
ON CONFLICT DO NOTHING;

INSERT INTO equivalences (id, origin_product_id, foreign_product_id, confidence, notes, source) VALUES
  ('e1000000-0000-0000-0000-000000000001', 'p1000000-0000-0000-0000-000000000001', 'p1000000-0000-0000-0000-000000000002', 0.95, 'Same active ingredient (fexofenadine 180mg); brand differs US→DE', 'curated'),
  ('e1000000-0000-0000-0000-000000000002', 'p1000000-0000-0000-0000-000000000001', 'p1000000-0000-0000-0000-000000000004', 0.92, 'Fexofenadine available in EE; Allegra brand not sold', 'curated'),
  ('e1000000-0000-0000-0000-000000000003', 'p1000000-0000-0000-0000-000000000001', 'p1000000-0000-0000-0000-000000000003', 0.75, 'Bilastine is a related 2nd-gen antihistamine; not identical — consult pharmacist', 'curated'),
  ('e1000000-0000-0000-0000-000000000004', 'p1000000-0000-0000-0000-000000000005', 'p1000000-0000-0000-0000-000000000006', 0.93, 'Cetirizine 10mg equivalent', 'curated')
ON CONFLICT (origin_product_id, foreign_product_id) DO NOTHING;

INSERT INTO avoid_rules (locale_from, term_from, locale_to, term_to, safe_alternate, severity, notes) VALUES
  ('en', 'corn syrup', 'de', 'Glukosesirup', 'Check label for maize-derived sweeteners', 'info', 'Naming differs; not always corn'),
  ('en', 'peanut', 'th', 'ถั่วลิสง', 'Ask for ปลอดถั่วลิสง', 'critical', 'Thai labeling'),
  ('en', 'gluten-free', 'fr', 'sans gluten', NULL, 'info', 'EU regulated claim'),
  ('ru', 'арахис', 'en', 'peanut', 'Peanut allergy', 'critical', 'Cyrillic ↔ Latin trap')
ON CONFLICT DO NOTHING;
