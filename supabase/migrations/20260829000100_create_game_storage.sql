create extension if not exists pgcrypto;

create table public.profiles (
  auth_id uuid primary key default gen_random_uuid(),
  legacy_id text not null unique,
  username text not null check (char_length(btrim(username)) between 1 and 32),
  password_hash text not null,
  created_at timestamptz not null default now(),
  last_login timestamptz not null default now(),
  is_active boolean not null default true,
  exp integer not null default 0 check (exp >= 0),
  level integer not null default 1 check (level >= 1),
  games_played integer not null default 0 check (games_played >= 0),
  games_won integer not null default 0 check (games_won >= 0 and games_won <= games_played),
  avatar text not null default '1',
  gold integer not null default 0 check (gold >= 0)
);

create unique index profiles_username_lower_key on public.profiles (lower(btrim(username)));

create table public.troops (
  name text primary key,
  max_hp double precision not null,
  hp double precision not null,
  dmg double precision not null default 0,
  atk double precision not null,
  def double precision not null,
  mana integer not null check (mana >= 0),
  crit integer not null check (crit between 0 and 100),
  exp integer not null default 0,
  speed double precision not null,
  range double precision not null,
  type text not null,
  image text not null,
  description text not null default '',
  attack_speed double precision not null,
  aggro_priority text not null,
  rarity text not null
);

create table public.towers (
  type text primary key,
  max_hp double precision not null,
  hp double precision not null,
  atk double precision not null,
  def double precision not null,
  crit double precision not null,
  exp integer not null default 0,
  range double precision not null,
  attack_speed double precision not null
);

create table public.legacy_sessions (
  session_id text primary key,
  profile_id uuid not null references public.profiles(auth_id) on delete cascade,
  authenticated boolean not null default true,
  expires_at timestamptz not null,
  created_at timestamptz not null default now()
);

create table public.sessions (
  session_id text primary key,
  profile_id uuid not null references public.profiles(auth_id) on delete cascade,
  expires_at timestamptz not null,
  created_at timestamptz not null default now()
);

create index legacy_sessions_profile_id_idx on public.legacy_sessions(profile_id);
create index legacy_sessions_expires_at_idx on public.legacy_sessions(expires_at);

alter table public.profiles enable row level security;
alter table public.troops enable row level security;
alter table public.towers enable row level security;
alter table public.legacy_sessions enable row level security;
alter table public.sessions enable row level security;

revoke all on table public.profiles from anon, authenticated;
revoke all on table public.troops from anon, authenticated;
revoke all on table public.towers from anon, authenticated;
revoke all on table public.legacy_sessions from anon, authenticated;
revoke all on table public.sessions from anon, authenticated;

grant select, insert, update, delete on table public.profiles to service_role;
grant select, insert, update, delete on table public.troops to service_role;
grant select, insert, update, delete on table public.towers to service_role;
grant select, insert, update, delete on table public.legacy_sessions to service_role;
grant select, insert, update, delete on table public.sessions to service_role;
