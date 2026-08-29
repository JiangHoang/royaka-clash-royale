insert into public.towers (type, max_hp, hp, atk, def, crit, exp, range, attack_speed) values
  ('King Tower', 2904, 2904, 72, 215, 10, 200, 5, 1),
  ('Guard Tower', 1600, 1600, 55, 120, 10, 100, 5, 1)
on conflict (type) do update set
  max_hp = excluded.max_hp, hp = excluded.hp, atk = excluded.atk,
  def = excluded.def, crit = excluded.crit, exp = excluded.exp,
  range = excluded.range, attack_speed = excluded.attack_speed;

insert into public.troops
  (name, max_hp, hp, dmg, atk, def, mana, crit, exp, speed, range, type,
   image, description, attack_speed, aggro_priority, rarity)
select
  'Local Troop ' || n, 1000 + n * 25, 1000 + n * 25, 100, 150 + n * 5,
  75, 3, 10, 5, 0.6, 1.2, 'fighter', 'Knight',
  'Local Supabase development troop', 1.2, 'troop', 'common'
from generate_series(1, 8) as n
on conflict (name) do nothing;
