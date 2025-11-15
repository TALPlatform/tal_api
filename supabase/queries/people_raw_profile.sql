
-- name: RawProfileList :many
select
  person_id::bigint,
  name::text,
  headline::text,
  region::text,
  current_title::text,
  current_company::text,
  industry::text,
  summary::text,
  years_of_experience::text,
  num_of_connections::int,
  profile_picture_url::text,
  linkedin_profile_url::text,
  skills::text[],
  semantic_score::float,
  text_rank::float,
  hybrid_score::float hybrid_score
from people_schema.raw_profile_search(
    in_query := sqlc.arg('query')::text,
    in_embedding := sqlc.arg('embedding')::vector(1536),
    in_industries := sqlc.arg('industries')::text[],
    in_locations := sqlc.arg('locations')::text[],
    in_skills := sqlc.arg('skills')::text[],
    in_companies := sqlc.arg('companies')::text[],
    in_projects := sqlc.arg('projects')::text[],
    in_limit := sqlc.arg('limit')::int
);

-- name: RawProfilesBulkCreateUpdate :exec
SELECT people_schema.raw_profiles_bulk_create_update(
  @session_id::int,
  @profiles::JSONB
);

-- name: RawProfileFind :one
SELECT
    person_id,
    coalesce(embedding_source_text , '') embedding_source_text,
    embedding_model,
    embedding_updated_at,
    first_name,
    last_name,
    flagship_profile_url,
    linkedin_profile_url,
    region,
    region_address_components,
    name,
    headline,
    location,
    current_title,
    current_company,
    industry,
    summary,
    experience,
    education,
    skills,
    certifications,
    languages,
    profile_language,
    profile_picture_url,
    profile_picture_permalink,
    twitter_handle,
    open_to_cards,
    num_of_connections,
    education_background,
    current_employers,
    past_employers,
    last_updated,
    recently_changed_jobs,
    years_of_experience,
    years_of_experience_raw,
    all_employers,
    honors,
    interests,
    projects,
    publications,
    volunteering,
    websites,
    version,
    run_id,
    created_at,
    updated_at
FROM
    people_schema.raw_profile
WHERE
    person_id = @person_id;

-- name: RawProfileListRequestBuild :many
-- select * from people_schema.raw_profile;
