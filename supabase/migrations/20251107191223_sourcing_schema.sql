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

    -- Filter Columns (Replaces initial_filters JSONB)
    min_experience INT DEFAULT NULL,
    max_experience INT DEFAULT NULL,
    required_contact_info TEXT[] DEFAULT NULL,
    timezone TEXT DEFAULT 'any',
    locations TEXT[] DEFAULT NULL,
    job_titles TEXT[] DEFAULT NULL,
    job_seniority TEXT DEFAULT NULL,
    job_functions TEXT[] DEFAULT NULL,
    companies TEXT[] DEFAULT NULL,
    company_headcount TEXT DEFAULT NULL,
    company_funding TEXT DEFAULT NULL,
    industries TEXT[] DEFAULT NULL,
    keywords TEXT[] DEFAULT NULL,
    skills TEXT[] DEFAULT NULL,
    education_levels TEXT[] DEFAULT NULL,
    languages TEXT[] DEFAULT NULL,
    filter_limit INT DEFAULT 50, -- Renamed from in_limit to avoid potential SQL keyword conflict

    -- Relationships and Audit Columns
    sourcing_project_id INT NOT NULL,
    FOREIGN KEY (sourcing_project_id) REFERENCES sourcing_schema.sourcing_project(sourcing_project_id),
    created_by INT NOT NULL,
    FOREIGN KEY (created_by) REFERENCES accounts_schema.user(user_id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
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

CREATE TABLE sourcing_schema.sourcing_session_criteria(
  sourcing_session_criteria_id SERIAL PRIMARY KEY,
  sourcing_session_id INT NOT NULL,
	FOREIGN KEY (sourcing_session_id) REFERENCES sourcing_schema.sourcing_session(sourcing_session_id),
  steps TEXT[] NOT NULL
);


CREATE TABLE sourcing_schema.sourcing_session_criteria_profile(
  sourcing_session_criteria_id INT ,
  FOREIGN KEY (sourcing_session_criteria_id) REFERENCES sourcing_schema.sourcing_session_criteria(sourcing_session_criteria_id),
  raw_profile_id INT ,
	FOREIGN KEY (raw_profile_id) REFERENCES people_schema.raw_profile (person_id),
  score int,
  order_index INT NOT NULL,
  criteria_matches JSONB[],
  PRIMARY KEY (sourcing_session_criteria_id , raw_profile_id)
);
