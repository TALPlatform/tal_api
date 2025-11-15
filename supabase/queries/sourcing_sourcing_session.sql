-- name: SourcingSessionCreateUpdate :one
SELECT
   s.sourcing_session_id::INT sourcing_session_id,
   s.sourcing_session_name::VARCHAR(255) sourcing_session_name,
   s.min_experience::INT min_experience,
   s.max_experience::INT max_experience,
   s.required_contact_info::TEXT[] required_contact_info,
   s.timezone::TEXT timezone,
   s.locations::TEXT[] locations,
   s.job_titles::TEXT[] job_titles,
   s.job_seniority::TEXT job_seniority,
   s.job_functions::TEXT[] job_functions,
   s.companies::TEXT[] companies,
   s.company_headcount::TEXT company_headcount,
   s.company_funding::TEXT company_funding,
   s.industries::TEXT[] industries,
   s.keywords::TEXT[] keywords,
   s.skills::TEXT[] skills,
   s.education_levels::TEXT[] education_levels,
   s.languages::TEXT[] languages,
   s.filter_limit::INT filter_limit,
   s.sourcing_project_id::INT sourcing_project_id,
   s.created_by::INT created_by,
   s.created_at::TIMESTAMP created_at,
   s.updated_at::TIMESTAMP updated_at,
   s.deleted_at::TIMESTAMP deleted_at
FROM sourcing_schema.sourcing_session_create_update(
    in_sourcing_session_id := @sourcing_session_id::INT,
    in_sourcing_session_name := @sourcing_session_name::VARCHAR(255),
    in_min_experience := @min_experience::INT,
    in_max_experience := @max_experience::INT,
    in_required_contact_info := @required_contact_info::TEXT[],
    in_timezone := @timezone::TEXT,
    in_locations := @locations::TEXT[],
    in_job_titles := @job_titles::TEXT[],
    in_job_seniority := @job_seniority::TEXT,
    in_job_functions := @job_functions::TEXT[],
    in_companies := @companies::TEXT[],
    in_company_headcount := @company_headcount::TEXT,
    in_company_funding := @company_funding::TEXT,
    in_industries := @industries::TEXT[],
    in_keywords := @keywords::TEXT[],
    in_skills := @skills::TEXT[],
    in_education_levels := @education_levels::TEXT[],
    in_languages := @languages::TEXT[],
    in_filter_limit := @filter_limit::INT,
    in_sourcing_project_id := @sourcing_project_id::INT,
    in_created_by := @created_by::INT
) s;

-- name: SourcingSessionFind :one
SELECT 
  COALESCE(sourcing_session_id, 0)::INT AS sourcing_session_id,
  COALESCE(sourcing_session_name, '')::VARCHAR(255) AS sourcing_session_name,
  COALESCE(min_experience, 0)::INT AS min_experience,
  COALESCE(max_experience, 0)::INT AS max_experience,
  COALESCE(required_contact_info, '{}')::TEXT[] AS required_contact_info,
  COALESCE(timezone, '')::TEXT AS timezone,
  COALESCE(locations, '{}')::TEXT[] AS locations,
  COALESCE(job_titles, '{}')::TEXT[] AS job_titles,
  COALESCE(job_seniority, '')::TEXT AS job_seniority,
  COALESCE(job_functions, '{}')::TEXT[] AS job_functions,
  COALESCE(companies, '{}')::TEXT[] AS companies,
  COALESCE(company_headcount, '')::TEXT AS company_headcount,
  COALESCE(company_funding, '')::TEXT AS company_funding,
  COALESCE(industries, '{}')::TEXT[] AS industries,
  COALESCE(keywords, '{}')::TEXT[] AS keywords,
  COALESCE(skills, '{}')::TEXT[] AS skills,
  COALESCE(education_levels, '{}')::TEXT[] AS education_levels,
  COALESCE(languages, '{}')::TEXT[] AS languages,
  COALESCE(filter_limit, 0)::INT AS filter_limit,
  
  COALESCE(session_created_at, '1970-01-01 00:00:00')::TIMESTAMP AS session_created_at,
  COALESCE(session_updated_at, '1970-01-01 00:00:00')::TIMESTAMP AS session_updated_at,
  COALESCE(session_deleted_at, '1970-01-01 00:00:00')::TIMESTAMP AS session_deleted_at,
  
  COALESCE(sourcing_project_id, 0)::INT AS sourcing_project_id,
  COALESCE(sourcing_project_name, '')::VARCHAR(255) AS sourcing_project_name,
  COALESCE(sourcing_project_breif, '')::TEXT AS sourcing_project_breif,
  COALESCE(tenant_id, 0)::INT AS tenant_id,
  COALESCE(project_created_at, '1970-01-01 00:00:00')::TIMESTAMP AS project_created_at,
  COALESCE(project_updated_at, '1970-01-01 00:00:00')::TIMESTAMP AS project_updated_at,
  COALESCE(project_deleted_at, '1970-01-01 00:00:00')::TIMESTAMP AS project_deleted_at,
  
  COALESCE(created_by, 0)::INT AS created_by,
  COALESCE(creator_name, '')::VARCHAR(255) AS creator_name,
  COALESCE(creator_email, '')::VARCHAR(255) AS creator_email,
  COALESCE(s.profiles::JSONB, '[]')::JSONB AS profiles
FROM sourcing_schema.sourcing_session_find(
in_sourcing_session_id := @sourcing_session_id::INT
) s;


-- name: SourcingSessionProfileCreateUpdate :one
SELECT 
  sourcing_session_id::INT,
  raw_profile_id::INT,
  score::INT,
  order_index::INT,
  note::TEXT,
  report_summary::TEXT,
  is_short_listed::BOOLEAN,
  summary_bullets::TEXT[],
  justification::TEXT
FROM sourcing_schema.sourcing_session_profile_create_update(
in_sourcing_session_id := @sourcing_session_id::INT,
in_raw_profile_id := @raw_profile_id::INT,
in_score := @score::INT,
in_order_index := @order_index::INT,
in_note := @note::TEXT,
in_report_summary := @report_summary::TEXT,
in_is_short_listed := @is_short_listed::BOOLEAN,
in_summary_bullets := @summary_bullets::TEXT[],
in_justification := @justification::TEXT
);

-- name: SourcingSessionList :many
select 
  s.sourcing_session_id sourcing_session_id,
  sp.sourcing_project_name sourcing_project_name,
  t.tenant_name tenant_name,
  u.user_name user_name,
  s.sourcing_session_name sourcing_session_name,
  s.min_experience min_experience,
  s.max_experience max_experience,
  s.required_contact_info required_contact_info,
  s.timezone timezone,
  s.locations locations,
  s.job_titles job_titles,
  s.job_seniority job_seniority,
  s.job_functions job_functions,
  s.companies companies,
  s.company_headcount company_headcount,
  s.company_funding company_funding,
  s.industries industries,
  s.keywords keywords,
  s.skills skills,
  s.education_levels education_levels,
  s.languages languages,
  s.filter_limit filter_limit,
  s.sourcing_project_id sourcing_project_id,
  s.created_by created_by,
  s.created_at created_at,
  s.updated_at updated_at,
  s.deleted_at deleted_at
from sourcing_schema.sourcing_session s
JOIN sourcing_schema.sourcing_project sp on s.sourcing_project_id = sp.sourcing_project_id
JOIN tenants_schema.tenant  t on t.tenant_id = sp.tenant_id
JOIN accounts_schema.user u on u.user_id = s.created_by
WHERE t.tenant_id = is_null_replace(@tenant_id::int , t.tenant_id)
;


-- name: SourcingSessionProfileFindForAI :one
SELECT 
  s.sourcing_session_name,
  s.min_experience::INT,
  s.max_experience::INT,
  s.required_contact_info::TEXT[],
  s.timezone::TEXT,
  s.locations::TEXT[],
  s.job_titles::TEXT[],
  s.job_seniority::TEXT,
  s.job_functions::TEXT[],
  s.companies::TEXT[],
  s.company_headcount::TEXT,
  s.company_funding::TEXT,
  s.industries::TEXT[],
  s.keywords::TEXT[],
  s.skills::TEXT[],
  s.education_levels::TEXT[],
  s.languages::TEXT[],
  s.filter_limit::INT,
  sp.sourcing_project_name,
  sp.sourcing_project_breif
FROM sourcing_schema.sourcing_session s
JOIN sourcing_schema.sourcing_project sp on s.sourcing_project_id = sp.sourcing_project_id
WHERE s.sourcing_session_id = @sourcing_session_id;


-- name: SourcingSessionCriteriaCreate :one
INSERT INTO sourcing_schema.sourcing_session_criteria(
  sourcing_session_id,
  steps
  ) VALUES (@sourcing_session_id , @steps) RETURNING *;

-- name: SourcingSessionCriteriaProfilesBulkInsert :copyfrom
INSERT INTO sourcing_schema.sourcing_session_criteria_profile (
  sourcing_session_criteria_id,
  raw_profile_id,
  score,
  order_index,
  criteria_matches
  ) VALUES ($1, $2, $3, $4, $5);

-- name: SourcingSessionApply :many
SELECT 
    COALESCE(person_id, 0)::BIGINT AS person_id,
    COALESCE(first_name, '')::TEXT AS first_name,
    COALESCE(last_name, '')::TEXT AS last_name,
    COALESCE(name, '')::TEXT AS name,
    COALESCE(headline, '')::TEXT AS headline,
    COALESCE(location, '')::TEXT AS location,
    COALESCE(current_title, '')::TEXT AS current_title,
    COALESCE(current_company, '')::TEXT AS current_company,
    COALESCE(industry, '')::TEXT AS industry,
    COALESCE(summary, '')::TEXT AS summary,
    COALESCE(linkedin_profile_url, '')::TEXT AS linkedin_profile_url,
    COALESCE(profile_picture_url, '')::TEXT AS profile_picture_url,
    COALESCE(years_of_experience, 0)::INT AS years_of_experience,
    COALESCE(num_of_connections, 0)::INT AS num_of_connections,
    COALESCE(skills, ARRAY[]::TEXT[])::TEXT[] AS skills,
    COALESCE(languages, ARRAY[]::TEXT[])::TEXT[] AS languages
FROM sourcing_schema.sourcing_session_apply(
    sqlc.arg(in_sourcing_session_id)::INT
);


-- name: SourcingSessionProfileSync :exec
SELECT sourcing_schema.sourcing_session_profile_sync(
    in_sourcing_session_id := @session_id::INT,
    in_crust_profiles := @crust_profiles::JSONB,
    in_db_profiles := @db_profiles::JSONB
);

-- name: SourcingSessionFiltersBuilder :one
-- select * from sourcing_schema.sourcing_session;

-- name: SourcingSessionFiltersInfered :one
-- select * from sourcing_schema.sourcing_session;
