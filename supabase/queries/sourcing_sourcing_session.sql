-- name: SourcingSessionCreateUpdate :one
SELECT 
    sourcing_session_id::INT,
    sourcing_session_name::VARCHAR(255),
    initial_filters::JSONB,
    sourcing_project_id::INT,
    created_by::INT,
    created_at::TIMESTAMP,
    updated_at::TIMESTAMP,
    deleted_at::TIMESTAMP
FROM sourcing_schema.sourcing_session_create_update(
    in_sourcing_session_id := @sourcing_session_id::INT,
    in_sourcing_session_name := @sourcing_session_name::VARCHAR(255),
    in_initial_filters := @initial_filters::JSONB,
    in_sourcing_project_id := @sourcing_project_id::INT,
    in_created_by := @created_by::INT
);

-- name: SourcingSessionFind :one
SELECT 
    sourcing_session_id::INT,
    sourcing_session_name::VARCHAR(255),
    initial_filters::JSONB,
    session_created_at::TIMESTAMP,
    session_updated_at::TIMESTAMP,
    session_deleted_at::TIMESTAMP,
    sourcing_project_id::INT,
    sourcing_project_name::VARCHAR(255),
    sourcing_project_breif::TEXT,
    tenant_id::INT,
    project_created_at::TIMESTAMP,
    project_updated_at::TIMESTAMP,
    project_deleted_at::TIMESTAMP,
    created_by::INT,
    creator_name::VARCHAR(255),
    creator_email::VARCHAR(255),
    profiles::JSONB
FROM sourcing_schema.sourcing_session_find_by_id(
    in_sourcing_session_id := @sourcing_session_id::INT
);

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
