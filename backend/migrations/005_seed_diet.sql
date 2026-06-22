-- Demo macro profile + travel-friendly food reference
INSERT INTO macro_profiles (user_id, protein_g, carbs_g, fat_g, calories_kcal, notes)
VALUES (
  '00000000-0000-0000-0000-000000000001',
  130, 80, 25,
  130 * 4 + 80 * 4 + 25 * 9,
  'Default gym-travel macros'
)
ON CONFLICT (user_id) DO UPDATE SET
  protein_g = EXCLUDED.protein_g,
  carbs_g = EXCLUDED.carbs_g,
  fat_g = EXCLUDED.fat_g;

INSERT INTO food_reference (id, name, name_local, locale, country_code, serving_label, protein_g, carbs_g, fat_g, calories_kcal, tags) VALUES
  ('f1000000-0000-0000-0000-000000000001', 'Grilled chicken breast', 'Hähnchenbrust', 'en', 'DE', '150g cooked', 46, 0, 5, 230, ARRAY['protein', 'safe']),
  ('f1000000-0000-0000-0000-000000000002', 'Greek yogurt plain', 'Griechischer Joghurt', 'en', 'GR', '200g cup', 20, 8, 4, 150, ARRAY['protein', 'dairy']),
  ('f1000000-0000-0000-0000-000000000003', 'Pho bo (beef noodle soup)', 'Phở bò', 'vi', 'VN', '1 restaurant bowl', 25, 55, 8, 420, ARRAY['noodles', 'soup']),
  ('f1000000-0000-0000-0000-000000000004', 'Pad thai shrimp', 'ผัดไทยกุ้ง', 'th', 'TH', '1 plate', 22, 68, 18, 520, ARRAY['noodles', 'street']),
  ('f1000000-0000-0000-0000-000000000005', 'Onigiri salmon', '鮭おにぎり', 'ja', 'JP', '1 piece', 6, 38, 3, 210, ARRAY['rice', 'conbini']),
  ('f1000000-0000-0000-0000-000000000006', 'Protein shake whey', NULL, 'en', 'US', '1 scoop + water', 24, 3, 1, 120, ARRAY['supplement']),
  ('f1000000-0000-0000-0000-000000000007', 'Boiled eggs', 'Eier', 'en', NULL, '2 large', 12, 1, 10, 155, ARRAY['protein', 'portable']),
  ('f1000000-0000-0000-0000-000000000008', 'Cottage cheese', 'Hüttenkäse', 'en', 'DE', '200g', 24, 8, 4, 165, ARRAY['protein', 'dairy']),
  ('f1000000-0000-0000-0000-000000000009', 'Rice bowl tuna', 'マグロ丼', 'ja', 'JP', '1 donburi', 35, 72, 6, 480, ARRAY['rice', 'fish']),
  ('f1000000-0000-0000-0000-000000000010', 'Banh mi thit', 'Bánh mì thịt', 'vi', 'VN', '1 sandwich', 18, 42, 14, 380, ARRAY['bread', 'street']),
  ('f1000000-0000-0000-0000-000000000011', 'Som tam (papaya salad)', 'ส้มตำ', 'th', 'TH', '1 plate', 4, 22, 6, 160, ARRAY['salad', 'spicy']),
  ('f1000000-0000-0000-0000-000000000012', 'Oatmeal cooked', 'Haferflocken', 'en', NULL, '80g dry', 10, 54, 6, 300, ARRAY['breakfast', 'carbs'])
ON CONFLICT (id) DO NOTHING;
