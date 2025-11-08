DROP FUNCTION IF EXISTS sourcing_schema.sourcing_session_create_update;
CREATE OR REPLACE FUNCTION sourcing_schema.sourcing_session_create_update(
    in_sourcing_session_id INT DEFAULT NULL,
    in_sourcing_session_name VARCHAR(255) DEFAULT NULL,
    in_initial_filters JSONB DEFAULT NULL,
    in_sourcing_project_id INT DEFAULT NULL,
    in_created_by INT DEFAULT NULL
)
RETURNS TABLE (
    sourcing_session_id INT,
    sourcing_session_name VARCHAR(255),
    initial_filters JSONB,
    sourcing_project_id INT,
    created_by INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    -- If no session_id provided, do a direct INSERT
    IF is_null(in_sourcing_session_id) THEN
        RETURN QUERY
        INSERT INTO sourcing_schema.sourcing_session (
            sourcing_session_name,
            initial_filters,
            sourcing_project_id,
            created_by
        )
        VALUES (
            in_sourcing_session_name,
            in_initial_filters,
            in_sourcing_project_id,
            in_created_by
        )
        RETURNING
            sourcing_schema.sourcing_session.sourcing_session_id,
            sourcing_schema.sourcing_session.sourcing_session_name,
            sourcing_schema.sourcing_session.initial_filters,
            sourcing_schema.sourcing_session.sourcing_project_id,
            sourcing_schema.sourcing_session.created_by,
            sourcing_schema.sourcing_session.created_at,
            sourcing_schema.sourcing_session.updated_at,
            sourcing_schema.sourcing_session.deleted_at;
    ELSE
        -- If session_id is provided, do UPDATE
        RETURN QUERY
        UPDATE sourcing_schema.sourcing_session
        SET
            sourcing_session_name = is_null_replace(in_sourcing_session_name, sourcing_session.sourcing_session_name),
            initial_filters = is_null_replace(in_initial_filters, sourcing_session.initial_filters),
            sourcing_project_id = is_null_replace(in_sourcing_project_id, sourcing_session.sourcing_project_id),
            updated_at = NOW()
        WHERE sourcing_session.sourcing_session_id = in_sourcing_session_id
        RETURNING
            sourcing_session.sourcing_session_id,
            sourcing_session.sourcing_session_name,
            sourcing_session.initial_filters,
            sourcing_session.sourcing_project_id,
            sourcing_session.created_by,
            sourcing_session.created_at,
            sourcing_session.updated_at,
            sourcing_session.deleted_at;
    END IF;
END;
$$;
DROP FUNCTION IF EXISTS sourcing_schema.sourcing_session_profile_create_update;
CREATE OR REPLACE FUNCTION sourcing_schema.sourcing_session_profile_create_update(
    in_sourcing_session_id INT DEFAULT NULL,
    in_raw_profile_id INT DEFAULT NULL,
    in_score INT DEFAULT NULL,
    in_order_index INT DEFAULT NULL,
    in_note TEXT DEFAULT NULL,
    in_report_summary TEXT DEFAULT NULL,
    in_is_short_listed BOOLEAN DEFAULT NULL,
    in_summary_bullets TEXT[] DEFAULT NULL,
    in_justification TEXT DEFAULT NULL
)
RETURNS TABLE (
    sourcing_session_id INT,
    raw_profile_id INT,
    score INT,
    order_index INT,
    note TEXT,
    report_summary TEXT,
    is_short_listed BOOLEAN,
    summary_bullets TEXT[],
    justification TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    INSERT INTO sourcing_schema.sourcing_session_profile (
        sourcing_session_id,
        raw_profile_id,
        score,
        order_index,
        note,
        report_summary,
        is_short_listed,
        summary_bullets,
        justification
    )
    VALUES (
        in_sourcing_session_id,
        in_raw_profile_id,
        in_score,
        in_order_index,
        in_note,
        in_report_summary,
        in_is_short_listed,
        in_summary_bullets,
        in_justification
    )
    ON CONFLICT (sourcing_session_id, raw_profile_id)
    DO UPDATE SET
        score = is_null_replace(EXCLUDED.score, sourcing_schema.sourcing_session_profile.score),
        order_index = is_null_replace(EXCLUDED.order_index, sourcing_schema.sourcing_session_profile.order_index),
        note = is_null_replace(EXCLUDED.note, sourcing_schema.sourcing_session_profile.note),
        report_summary = is_null_replace(EXCLUDED.report_summary, sourcing_schema.sourcing_session_profile.report_summary),
        is_short_listed = is_null_replace(EXCLUDED.is_short_listed, sourcing_schema.sourcing_session_profile.is_short_listed),
        summary_bullets = is_null_replace(EXCLUDED.summary_bullets, sourcing_schema.sourcing_session_profile.summary_bullets),
        justification = is_null_replace(EXCLUDED.justification, sourcing_schema.sourcing_session_profile.justification)
    RETURNING
        sourcing_schema.sourcing_session_profile.sourcing_session_id,
        sourcing_schema.sourcing_session_profile.raw_profile_id,
        sourcing_schema.sourcing_session_profile.score,
        sourcing_schema.sourcing_session_profile.order_index,
        sourcing_schema.sourcing_session_profile.note,
        sourcing_schema.sourcing_session_profile.report_summary,
        sourcing_schema.sourcing_session_profile.is_short_listed,
        sourcing_schema.sourcing_session_profile.summary_bullets,
        sourcing_schema.sourcing_session_profile.justification;
END;
$$;
