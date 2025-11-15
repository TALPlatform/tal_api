create extension http with schema extensions;
-- Allow authenticated users full access to the abchotels bucket
CREATE POLICY "Allow unauth INSERTs" ON storage.buckets
	FOR INSERT
		WITH CHECK (
TRUE);

CREATE POLICY "Allow unauth INSERTs" ON storage.objects
	FOR INSERT
		WITH CHECK (
TRUE);

CREATE POLICY "Non-Authenticated users can manage abchotels bucket" ON storage.objects
	FOR ALL
		USING (bucket_id = 'abchotels')
		WITH CHECK (bucket_id = 'abchotels');

SELECT
	*
FROM
	accounts_schema.permissions_populate(execluded_tables => ARRAY['navigation_bar_item', 'log', 'setting', 'input_type'], added_tables => ARRAY['storage.bucket', 'storage.file', 'auth.session']);

INSERT INTO accounts_schema.permission (
  permission_function,
  permission_name,
  permission_description,
  permission_group
  ) values (
  'TenantDashboard',
  'tenant dashboard',
  'dashboard to view data related to tenats',
  'tenant'
  );

INSERT INTO accounts_schema.role (
  role_name,
  role_description,
  role_security_level)
VALUES (
  'super admin',
  'this is the most privlidged role that can do all the permissions on the sysstem',
  100);

INSERT INTO accounts_schema.role_permission (
  role_id,
  permission_id)
  SELECT
    1,
    permission_id
  FROM
  accounts_schema.permission;

INSERT INTO input_type (
  input_type_name)
VALUES (
  'text'),
  (
  'toggle'),
  (
  'richtext'),
  (
  'date');

INSERT INTO accounts_schema.user_type (
  user_type_name)
VALUES (
  'admin'),
  (
  'tenant');

INSERT INTO setting (
  input_type_id,
  setting_key,
  setting_value)
VALUES (
  1,
  'site_name',
  'Abc Hotels');


update 
accounts_schema.permission set permission_group = 'sourcing_session_criteria' where lower(permission_function) like '%sourcingsessioncriteria%';


update 
accounts_schema.permission set permission_group = 'sourcing_session' where lower(permission_function) like '%sourcingsession%';
update 
accounts_schema.permission set permission_group = 'sourcing_project' where lower(permission_function) like '%sourcingpro%';
