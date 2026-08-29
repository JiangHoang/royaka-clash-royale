begin;
select plan(14);

select has_table('public', 'profiles', 'profiles exists');
select has_table('public', 'troops', 'troops exists');
select has_table('public', 'towers', 'towers exists');
select has_table('public', 'legacy_sessions', 'legacy sessions exists');
select has_table('public', 'sessions', 'sessions exists');
select ok(row_security_active('public.profiles'), 'profiles RLS enabled');
select ok(row_security_active('public.troops'), 'troops RLS enabled');
select ok(row_security_active('public.towers'), 'towers RLS enabled');
select ok(row_security_active('public.legacy_sessions'), 'legacy sessions RLS enabled');
select ok(row_security_active('public.sessions'), 'sessions RLS enabled');
select has_index('public', 'profiles', 'profiles_username_lower_key', 'username index exists');
select has_index('public', 'legacy_sessions', 'legacy_sessions_expires_at_idx', 'expiry index exists');
select col_is_pk('public', 'troops', 'name', 'troop name is primary key');
select col_is_pk('public', 'towers', 'type', 'tower type is primary key');

select * from finish();
rollback;
