-- Demo locker items linked to seed products for trip report demo
INSERT INTO locker_items (id, user_id, product_id, custom_name, dose_per_day, dose_unit, frequency, category) VALUES
  ('l1000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001',
   'p1000000-0000-0000-0000-000000000001', 'Allegra 180mg', 1, 'tablet', 'daily', 'consume'),
  ('l1000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001',
   'p1000000-0000-0000-0000-000000000005', 'Zyrtec 10mg', 1, 'tablet', 'daily', 'consume')
ON CONFLICT (id) DO NOTHING;

-- Placeholder product image paths (served by frontend static or CDN later)
UPDATE products SET image_url = '/products/allegra-us.svg' WHERE id = 'p1000000-0000-0000-0000-000000000001';
UPDATE products SET image_url = '/products/telfast-de.svg' WHERE id = 'p1000000-0000-0000-0000-000000000002';
UPDATE products SET image_url = '/products/fexofast-ee.svg' WHERE id = 'p1000000-0000-0000-0000-000000000004';
UPDATE products SET image_url = '/products/zyrtec-us.svg' WHERE id = 'p1000000-0000-0000-0000-000000000005';
UPDATE products SET image_url = '/products/cetirizin-de.svg' WHERE id = 'p1000000-0000-0000-0000-000000000006';

-- Germany trip equivalence for Zyrtec
INSERT INTO equivalences (id, origin_product_id, foreign_product_id, confidence, notes, source) VALUES
  ('e1000000-0000-0000-0000-000000000005', 'p1000000-0000-0000-0000-000000000005', 'p1000000-0000-0000-0000-000000000006', 0.93, 'Cetirizine 10mg — ask Apotheke', 'curated')
ON CONFLICT (origin_product_id, foreign_product_id) DO NOTHING;
