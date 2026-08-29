drop trigger if exists on_auth_user_created on auth.users;
drop function if exists public.handle_new_auth_user();

alter table public.profiles add column if not exists password_hash text;

update public.profiles p
set password_hash = u.encrypted_password
from auth.users u
where u.id = p.auth_id and p.password_hash is null;

alter table public.profiles drop constraint if exists profiles_auth_id_fkey;
alter table public.profiles alter column auth_id set default gen_random_uuid();

create table if not exists public.sessions (
  session_id text primary key,
  profile_id uuid not null references public.profiles(auth_id) on delete cascade,
  expires_at timestamptz not null,
  created_at timestamptz not null default now()
);

alter table public.sessions enable row level security;
revoke all on table public.sessions from anon, authenticated;
grant select, insert, update, delete on table public.sessions to service_role;

-- If imported users were not present in auth.users, rerun the JSON importer
-- before enforcing this constraint.
do $$
begin
  if not exists (select 1 from public.profiles where password_hash is null) then
    alter table public.profiles alter column password_hash set not null;
  end if;
end $$;
