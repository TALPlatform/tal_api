DROP SCHEMA IF EXISTS sourcing_schema CASCADE;
CREATE SCHEMA sourcing_schema;

-------------------------------
-- Lookup / normalized tables
-------------------------------
CREATE TABLE sourcing_schema.sourcing_project (
  sourcing_project_id SERIAL PRIMARY KEY,
  sourcing_project_name VARCHAR(255) NOT NULL,
  sourcing_project_breif TEXT,
	tenant_id int,
	FOREIGN KEY (tenant_id) REFERENCES tenants_schema.tenant (tenant_id),
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp,
	deleted_at timestamp
);

CREATE TABLE sourcing_schema.sourcing_session (
  sourcing_session_id SERIAL PRIMARY KEY,
  sourcing_session_name VARCHAR(255) NOT NULL,
  initial_filters JSONB,
	sourcing_project_id int NOT NULL,
	FOREIGN KEY (sourcing_project_id) REFERENCES sourcing_schema.sourcing_project(sourcing_project_id),
	created_by int NOT NULL ,
	FOREIGN KEY (created_by) REFERENCES accounts_schema.user(user_id),
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp,
	deleted_at timestamp
);
-- all search profiles will stored here
CREATE TABLE sourcing_schema.sourcing_session_profile (
  sourcing_session_id INT ,
	FOREIGN KEY (sourcing_session_id) REFERENCES sourcing_schema.sourcing_session(sourcing_session_id),
  raw_profile_id INT ,
	FOREIGN KEY (raw_profile_id) REFERENCES people_schema.raw_profile (person_id),
  score int,
  order_index INT NOT NULL,
  note TEXT,
  report_summary TEXT,
  is_short_listed BOOLEAN,
  summary_bullets TEXT[],
  justification TEXT,
  PRIMARY KEY (sourcing_session_id , raw_profile_id)
);
