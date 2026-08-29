alter table public.profiles
  drop constraint if exists profiles_username_check;

alter table public.profiles
  add constraint profiles_username_check
  check (char_length(btrim(username)) between 1 and 32);
