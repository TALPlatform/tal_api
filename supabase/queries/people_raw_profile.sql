
-- name: RawProfileList :many
SELECT
  person_id::BIGINT,
  name::TEXT,
  headline::TEXT,
  location::TEXT,
  current_title::TEXT,
  current_company::TEXT,
  industry::TEXT,
  summary::TEXT,
  years_of_experience::TEXT,
  num_of_connections::INT,
  profile_picture_url::TEXT,
  linkedin_profile_url::TEXT,
  skills::TEXT[],
  semantic_score::FLOAT,
  text_rank::FLOAT,
  hybrid_score::FLOAT hybrid_score
FROM people_schema.raw_profile_search(
    in_query := sqlc.arg('query')::TEXT,
    in_embedding := sqlc.arg('embedding')::VECTOR(1536),
    in_industries := sqlc.arg('industries')::TEXT[],
    in_locations := sqlc.arg('locations')::TEXT[],
    in_skills := sqlc.arg('skills')::TEXT[],
    in_companies := sqlc.arg('companies')::TEXT[],
    in_projects := sqlc.arg('projects')::TEXT[],
    in_limit := sqlc.arg('limit')::INT
);

-- name: RawProfilesBulkCreateUpdate :exec
SELECT people_schema.raw_profiles_bulk_create_update(
  @session_id::int,
  @profiles::JSONB
);
