DROP FUNCTION IF EXISTS sourcing_schema.sourcing_session_create_update;
CREATE OR REPLACE FUNCTION sourcing_schema.sourcing_session_create_update(
    in_sourcing_session_id INT DEFAULT NULL,
    in_sourcing_session_name VARCHAR(255) DEFAULT NULL,
    in_min_experience INT DEFAULT NULL,
    in_max_experience INT DEFAULT NULL,
    in_required_contact_info TEXT[] DEFAULT NULL,
    in_timezone TEXT DEFAULT NULL,
    in_locations TEXT[] DEFAULT NULL,
    in_job_titles TEXT[] DEFAULT NULL,
    in_job_seniority TEXT DEFAULT NULL,
    in_job_functions TEXT[] DEFAULT NULL,
    in_companies TEXT[] DEFAULT NULL,
    in_company_headcount TEXT DEFAULT NULL,
    in_company_funding TEXT DEFAULT NULL,
    in_industries TEXT[] DEFAULT NULL,
    in_keywords TEXT[] DEFAULT NULL,
    in_skills TEXT[] DEFAULT NULL,
    in_education_levels TEXT[] DEFAULT NULL,
    in_languages TEXT[] DEFAULT NULL,
    in_filter_limit INT DEFAULT NULL,
    in_sourcing_project_id INT DEFAULT NULL,
    in_created_by INT DEFAULT NULL
)
RETURNS TABLE (
    sourcing_session_id INT,
    sourcing_session_name VARCHAR(255),
    min_experience INT,
    max_experience INT,
    required_contact_info TEXT[],
    timezone TEXT,
    locations TEXT[],
    job_titles TEXT[],
    job_seniority TEXT,
    job_functions TEXT[],
    companies TEXT[],
    company_headcount TEXT,
    company_funding TEXT,
    industries TEXT[],
    keywords TEXT[],
    skills TEXT[],
    education_levels TEXT[],
    languages TEXT[],
    filter_limit INT,
    sourcing_project_id INT,
    created_by INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    IF is_null(in_sourcing_session_id) THEN
        INSERT INTO sourcing_schema.sourcing_session (
            sourcing_session_name,
            -- New Filter Columns
            min_experience,
            max_experience,
            required_contact_info,
            timezone,
            locations,
            job_titles,
            job_seniority,
            job_functions,
            companies,
            company_headcount,
            company_funding,
            industries,
            keywords,
            skills,
            education_levels,
            languages,
            filter_limit,
            -- Existing Columns
            sourcing_project_id,
            created_by
        )
        VALUES (
            in_sourcing_session_name,
            -- New Filter Values
            in_min_experience,
            in_max_experience,
            in_required_contact_info,
            in_timezone,
            in_locations,
            in_job_titles,
            in_job_seniority,
            in_job_functions,
            in_companies,
            in_company_headcount,
            in_company_funding,
            in_industries,
            in_keywords,
            in_skills,
            in_education_levels,
            in_languages,
            in_filter_limit,
            -- Existing Values
            in_sourcing_project_id,
            in_created_by
        )
        RETURNING sourcing_session.sourcing_session_id INTO in_sourcing_session_id;
      ELSE
        UPDATE sourcing_schema.sourcing_session AS s
            SET
                sourcing_session_name = is_null_replace(in_sourcing_session_name, s.sourcing_session_name),
                sourcing_project_id = is_null_replace(in_sourcing_project_id, s.sourcing_project_id),
                min_experience = is_null_replace(in_min_experience, s.min_experience),
                max_experience = is_null_replace(in_max_experience, s.max_experience),
                required_contact_info = is_null_replace(in_required_contact_info, s.required_contact_info),
                timezone = is_null_replace(in_timezone, s.timezone),
                locations = is_null_replace(in_locations, s.locations),
                job_titles = is_null_replace(in_job_titles, s.job_titles),
                job_seniority = is_null_replace(in_job_seniority, s.job_seniority),
                job_functions = is_null_replace(in_job_functions, s.job_functions),
                companies = is_null_replace(in_companies, s.companies),
                company_headcount = is_null_replace(in_company_headcount, s.company_headcount),
                company_funding = is_null_replace(in_company_funding, s.company_funding),
                industries = is_null_replace(in_industries, s.industries),
                keywords = is_null_replace(in_keywords, s.keywords),
                skills = is_null_replace(in_skills, s.skills),
                education_levels = is_null_replace(in_education_levels, s.education_levels),
                languages = is_null_replace(in_languages, s.languages),
                filter_limit = is_null_replace(in_filter_limit, s.filter_limit),
                updated_at = NOW()
            WHERE s.sourcing_session_id = in_sourcing_session_id;
    END IF;

    RETURN query
        SELECT             
              s.sourcing_session_id,
              s.sourcing_session_name,
              s.min_experience,
              s.max_experience,
              s.required_contact_info,
              s.timezone,
              s.locations,
              s.job_titles,
              s.job_seniority,
              s.job_functions,
              s.companies,
              s.company_headcount,
              s.company_funding,
              s.industries,
              s.keywords,
              s.skills,
              s.education_levels,
              s.languages,
              s.filter_limit,
              s.sourcing_project_id,
              s.created_by,
              s.created_at,
              s.updated_at,
              s.deleted_at -- Return all columns from the updated row
  FROM sourcing_schema.sourcing_session s where s.sourcing_session_id = in_sourcing_session_id;
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

DROP FUNCTION IF EXISTS sourcing_schema.sourcing_session_find;
CREATE OR REPLACE FUNCTION sourcing_schema.sourcing_session_find(
    in_sourcing_session_id INT
)
RETURNS TABLE (
    -- Session fields
    sourcing_session_id INT,
    sourcing_session_name VARCHAR(255),
    
    -- Filter fields (NEW: Replacing initial_filters JSONB)
    min_experience INT,
    max_experience INT,
    required_contact_info TEXT[],
    timezone TEXT,
    locations TEXT[],
    job_titles TEXT[],
    job_seniority TEXT,
    job_functions TEXT[],
    companies TEXT[],
    company_headcount TEXT,
    company_funding TEXT,
    industries TEXT[],
    keywords TEXT[],
    skills TEXT[],
    education_levels TEXT[],
    languages TEXT[],
    filter_limit INT,

    session_created_at TIMESTAMP,
    session_updated_at TIMESTAMP,
    session_deleted_at TIMESTAMP,
    
    -- Project fields
    sourcing_project_id INT,
    sourcing_project_name VARCHAR(255),
    sourcing_project_breif TEXT,
    tenant_id INT,
    project_created_at TIMESTAMP,
    project_updated_at TIMESTAMP,
    project_deleted_at TIMESTAMP,
    
    -- Creator fields
    created_by INT,
    creator_name VARCHAR(255),
    creator_email VARCHAR(255),
    
    -- Profiles array
    profiles JSONB
)
LANGUAGE plpgsql
AS $$
BEGIN
    -- 1️⃣ Create temp table with profiles data (UNCHANGED)
    CREATE TEMP TABLE tmp_session_profiles AS
    SELECT
        ssp.sourcing_session_id,
        ssp.raw_profile_id,
        ssp.score,
        ssp.order_index,
        ssp.note,
        ssp.report_summary,
        ssp.is_short_listed,
        ssp.summary_bullets,
        ssp.justification,
        rp.person_id,
        rp.name,
        rp.headline,
        rp.location,
        rp.current_title,
        rp.current_company,
        rp.industry,
        rp.summary,
        rp.skills,
        rp.years_of_experience,
        rp.num_of_connections,
        rp.profile_picture_url,
        rp.linkedin_profile_url
    FROM sourcing_schema.sourcing_session_profile ssp
    LEFT JOIN people_schema.raw_profile rp ON ssp.raw_profile_id = rp.person_id
    WHERE ssp.sourcing_session_id = in_sourcing_session_id
    ORDER BY ssp.order_index;

    -- 2️⃣ Return main query with aggregated profiles (UPDATED)
    RETURN QUERY
    SELECT
        -- Session fields
        ss.sourcing_session_id,
        ss.sourcing_session_name,

        -- Filter fields (NEW)
        ss.min_experience,
        ss.max_experience,
        ss.required_contact_info,
        ss.timezone,
        ss.locations,
        ss.job_titles,
        ss.job_seniority,
        ss.job_functions,
        ss.companies,
        ss.company_headcount,
        ss.company_funding,
        ss.industries,
        ss.keywords,
        ss.skills,
        ss.education_levels,
        ss.languages,
        ss.filter_limit,

        ss.created_at AS session_created_at,
        ss.updated_at AS session_updated_at,
        ss.deleted_at AS session_deleted_at,
        
        -- Project fields
        sp.sourcing_project_id,
        sp.sourcing_project_name,
        sp.sourcing_project_breif,
        sp.tenant_id,
        sp.created_at AS project_created_at,
        sp.updated_at AS project_updated_at,
        sp.deleted_at AS project_deleted_at,
        
        -- Creator fields
        ss.created_by,
        u.user_name AS creator_name,
        u.user_email AS creator_email,
        
        -- Profiles array from temp table
        COALESCE(
            (
                SELECT jsonb_agg(
                    jsonb_build_object(
                        'sourcing_session_id', tmp.sourcing_session_id,
                        'raw_profile_id', tmp.raw_profile_id,
                        'score', tmp.score,
                        'order_index', tmp.order_index,
                        'note', tmp.note,
                        'report_summary', tmp.report_summary,
                        'is_short_listed', tmp.is_short_listed,
                        'summary_bullets', tmp.summary_bullets,
                        'justification', tmp.justification,
                        'profile', jsonb_build_object(
                            'person_id', tmp.person_id,
                            'name', tmp.name,
                            'headline', tmp.headline,
                            'location', tmp.location,
                            'current_title', tmp.current_title,
                            'current_company', tmp.current_company,
                            'industry', tmp.industry,
                            'summary', tmp.summary,
                            'skills', tmp.skills,
                            'years_of_experience', tmp.years_of_experience,
                            'num_of_connections', tmp.num_of_connections,
                            'profile_picture_url', tmp.profile_picture_url,
                            'linkedin_profile_url', tmp.linkedin_profile_url
                        )
                    )
                )
                FROM tmp_session_profiles tmp
            ),
            '[]'::JSONB
        ) AS profiles
        
    FROM sourcing_schema.sourcing_session ss
    JOIN sourcing_schema.sourcing_project sp 
        ON ss.sourcing_project_id = sp.sourcing_project_id
    LEFT JOIN accounts_schema.user u 
        ON ss.created_by = u.user_id
    WHERE ss.sourcing_session_id = in_sourcing_session_id;

    -- 3️⃣ Drop temp table (UNCHANGED)
    DROP TABLE IF EXISTS tmp_session_profiles;
END;
$$;
DROP FUNCTION IF EXISTS sourcing_schema.sourcing_session_apply;
CREATE OR REPLACE FUNCTION sourcing_schema.sourcing_session_apply(
    in_sourcing_session_id INT
)
RETURNS TABLE (
    person_id BIGINT,
    first_name TEXT,
    last_name TEXT,
    name TEXT,
    headline TEXT,
    location TEXT,
    current_title TEXT,
    current_company TEXT,
    industry TEXT,
    summary TEXT,
    linkedin_profile_url TEXT,
    profile_picture_url TEXT,
    years_of_experience INT,
    num_of_connections INT,
    skills TEXT[],
    languages TEXT[]
)
LANGUAGE plpgsql
AS $$
DECLARE
    session_record sourcing_schema.sourcing_session%ROWTYPE;
BEGIN
    -- Load sourcing session filters
    SELECT *
    INTO session_record
    FROM sourcing_schema.sourcing_session
    WHERE sourcing_session_id = in_sourcing_session_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Sourcing session with id % not found', in_sourcing_session_id;
    END IF;

    -- Apply filters to raw profiles
    RETURN QUERY
    SELECT 
        rp.person_id,
        rp.first_name,
        rp.last_name,
        rp.name,
        rp.headline,
        rp.location,
        rp.current_title,
        rp.current_company,
        rp.industry,
        rp.summary,
        rp.linkedin_profile_url,
        rp.profile_picture_url,
        rp.years_of_experience_raw AS years_of_experience,
        rp.num_of_connections,
        rp.skills,
        rp.languages
    FROM people_schema.raw_profile rp
    WHERE 
        -- Experience range filter
        (session_record.min_experience IS NULL OR rp.years_of_experience_raw >= session_record.min_experience)
        AND (session_record.max_experience IS NULL OR rp.years_of_experience_raw <= session_record.max_experience)

        -- Required contact info
        AND (
            session_record.required_contact_info IS NULL 
            OR (
                ('email' = ANY(session_record.required_contact_info) AND rp.flagship_profile_url IS NOT NULL)
                OR ('linkedin' = ANY(session_record.required_contact_info) AND rp.linkedin_profile_url IS NOT NULL)
                OR ('twitter' = ANY(session_record.required_contact_info) AND rp.twitter_handle IS NOT NULL)
            )
        )

        -- Location filter
        AND (
            session_record.locations IS NULL 
            OR EXISTS (
                SELECT 1 
                FROM unnest(session_record.locations) AS loc
                WHERE rp.location ILIKE '%' || loc || '%' 
                   OR rp.region ILIKE '%' || loc || '%'
            )
        )

        -- Job titles filter
        AND (
            session_record.job_titles IS NULL 
            OR EXISTS (
                SELECT 1 
                FROM unnest(session_record.job_titles) AS title
                WHERE rp.current_title ILIKE '%' || title || '%'
                   OR rp.headline ILIKE '%' || title || '%'
            )
        )

        -- Job seniority filter
        AND (
            session_record.job_seniority IS NULL
            OR (
                session_record.job_seniority = 'entry' AND (
                    rp.current_title ~* '\y(junior|entry|associate|assistant|intern|trainee)\y'
                    OR rp.headline ~* '\y(junior|entry|associate|assistant|intern|trainee)\y'
                )
            )
            OR (
                session_record.job_seniority = 'mid' AND (
                    rp.current_title ~* '\y(specialist|coordinator|analyst|engineer|developer)\y'
                    AND NOT rp.current_title ~* '\y(senior|lead|principal|chief|head|director|vp|vice president)\y'
                )
            )
            OR (
                session_record.job_seniority = 'senior' AND (
                    rp.current_title ~* '\y(senior|lead|principal)\y'
                    OR rp.headline ~* '\y(senior|lead|principal)\y'
                )
            )
            OR (
                session_record.job_seniority = 'executive' AND (
                    rp.current_title ~* '\y(chief|ceo|cto|cfo|coo|vp|vice president|director|head of)\y'
                    OR rp.headline ~* '\y(chief|ceo|cto|cfo|coo|vp|vice president|director|head of)\y'
                )
            )
        )

        -- Companies filter
        AND (
            session_record.companies IS NULL 
            OR EXISTS (
                SELECT 1 
                FROM unnest(session_record.companies) AS comp
                WHERE rp.current_company ILIKE '%' || comp || '%'
                   OR EXISTS (
                       SELECT 1 
                       FROM unnest(rp.all_employers) AS emp
                       WHERE emp->>'name' ILIKE '%' || comp || '%'
                   )
            )
        )

        -- Industries filter
        AND (
            session_record.industries IS NULL
            OR EXISTS (
                SELECT 1 FROM unnest(session_record.industries) AS ind
                WHERE
                    rp.industry ILIKE '%' || ind || '%'
                    OR EXISTS (
                        SELECT 1 FROM unnest(rp.all_employers) AS emp_json
                        WHERE emp_json->>'name' ILIKE '%' || ind || '%'
                           OR emp_json->>'industry' ILIKE '%' || ind || '%'
                    )
                    OR EXISTS (
                        SELECT 1 FROM jsonb_array_elements(rp.projects) AS proj
                        WHERE proj->>'title' ILIKE '%' || ind || '%'
                           OR proj->>'description' ILIKE '%' || ind || '%'
                    )
            )
        )
        -- Keywords filter
        AND (
            session_record.keywords IS NULL 
            OR EXISTS (
                SELECT 1 
                FROM unnest(session_record.keywords) AS kw
                WHERE rp.search_terms @@ plainto_tsquery('english', kw)
                   OR rp.headline ILIKE '%' || kw || '%'
                   OR rp.summary ILIKE '%' || kw || '%'
                   OR rp.current_title ILIKE '%' || kw || '%'
            )
        )

        -- Skills overlap filter
        AND (
            session_record.skills IS NULL 
            OR rp.skills && session_record.skills
        )

        -- Education level filter
        AND (
            session_record.education_levels IS NULL 
            OR EXISTS (
                SELECT 1 
                FROM unnest(rp.education_background) AS edu,
                     unnest(session_record.education_levels) AS level
                WHERE edu->>'degree' ILIKE '%' || level || '%'
            )
        )

        -- Languages filter
        AND (
            session_record.languages IS NULL 
            OR rp.languages && session_record.languages
        )
        
    ORDER BY 
        (CASE WHEN rp.linkedin_profile_url IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN rp.summary IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN rp.years_of_experience_raw IS NOT NULL THEN 1 ELSE 0 END) DESC,
        rp.num_of_connections DESC NULLS LAST,
        rp.updated_at DESC

    LIMIT COALESCE(session_record.filter_limit, 50);
END;
$$;
DROP FUNCTION IF EXISTS sourcing_schema.sourcing_session_profile_sync;
CREATE OR REPLACE FUNCTION sourcing_schema.sourcing_session_profile_sync(
    in_sourcing_session_id INT,
    in_crust_profiles JSONB,
    in_db_profiles JSONB
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    -- 1️⃣ Create temporary table for crust profiles
DROP TABLE IF EXISTS tmp_raw_profiles;
    CREATE TEMP TABLE tmp_raw_profiles AS
    SELECT
        (p.json_data ->> 'person_id')::BIGINT AS person_id,
        p.json_data ->> 'embedding_model' AS embedding_model,
        (p.json_data ->> 'embedding_updated_at')::TIMESTAMP AS embedding_updated_at,
        p.json_data ->> 'first_name' AS first_name,
        p.json_data ->> 'last_name' AS last_name,
        p.json_data ->> 'flagship_profile_url' AS flagship_profile_url,
        p.json_data ->> 'linkedin_profile_url' AS linkedin_profile_url,
        p.json_data ->> 'region' AS region,
        p.json_data ->> 'region_address_components' AS region_address_components,
        p.json_data ->> 'name' AS name,
        p.json_data ->> 'headline' AS headline,
        p.json_data ->> 'location' AS location,
        p.json_data ->> 'current_title' AS current_title,
        p.json_data ->> 'current_company' AS current_company,
        p.json_data ->> 'industry' AS industry,
        p.json_data ->> 'summary' AS summary,

        -- Array fields
CASE 
    WHEN jsonb_typeof(p.json_data -> 'skills') = 'array' 
    THEN ARRAY(SELECT jsonb_array_elements_text(p.json_data -> 'skills'))
    ELSE ARRAY[]::TEXT[]
END AS skills,
CASE 
    WHEN jsonb_typeof(p.json_data -> 'languages') = 'array' 
    THEN ARRAY(SELECT jsonb_array_elements_text(p.json_data -> 'languages'))
    ELSE ARRAY[]::TEXT[]
END AS languages,

        -- JSONB fields (scalar)
        p.json_data -> 'certifications' AS certifications,
        p.json_data -> 'education' AS education,
        p.json_data -> 'honors' AS honors,
        p.json_data -> 'interests' AS interests,
        p.json_data -> 'projects' AS projects,
        p.json_data -> 'publications' AS publications,
        p.json_data -> 'volunteering' AS volunteering,
        p.json_data -> 'websites' AS websites,

        -- JSONB[] fields
CASE 
    WHEN jsonb_typeof(p.json_data -> 'open_to_cards') = 'array' 
    THEN ARRAY(SELECT jsonb_array_elements(p.json_data -> 'open_to_cards'))
    ELSE ARRAY[]::JSONB[]
END AS open_to_cards,

CASE 
    WHEN jsonb_typeof(p.json_data -> 'education_background') = 'array' 
    THEN ARRAY(SELECT jsonb_array_elements(p.json_data -> 'education_background'))
    ELSE ARRAY[]::JSONB[]
END AS education_background,

CASE 
    WHEN jsonb_typeof(p.json_data -> 'current_employers') = 'array' 
    THEN ARRAY(SELECT jsonb_array_elements(p.json_data -> 'current_employers'))
    ELSE ARRAY[]::JSONB[]
END AS current_employers,

CASE 
    WHEN jsonb_typeof(p.json_data -> 'past_employers') = 'array' 
    THEN ARRAY(SELECT jsonb_array_elements(p.json_data -> 'past_employers'))
    ELSE ARRAY[]::JSONB[]
END AS past_employers,

CASE 
    WHEN jsonb_typeof(p.json_data -> 'all_employers') = 'array' 
    THEN ARRAY(SELECT jsonb_array_elements(p.json_data -> 'all_employers'))
    ELSE ARRAY[]::JSONB[]
END AS all_employers,
        -- Other fields
        p.json_data ->> 'profile_picture_url' AS profile_picture_url,
        p.json_data ->> 'profile_picture_permalink' AS profile_picture_permalink,
        p.json_data ->> 'twitter_handle' AS twitter_handle,
        (p.json_data ->> 'num_of_connections')::INT AS num_of_connections,
        (p.json_data ->> 'recently_changed_jobs')::BOOLEAN AS recently_changed_jobs,
        p.json_data ->> 'years_of_experience' AS years_of_experience,
        (p.json_data ->> 'years_of_experience_raw')::INT AS years_of_experience_raw,
        p.json_data ->> 'profile_language' AS profile_language,
        p.json_data ->> 'version' AS version,
        (p.json_data ->> 'run_id')::BIGINT AS run_id,
        (p.json_data ->> 'last_updated')::TIMESTAMP AS last_updated,

        -- Convert the embedding JSON/text array to a float vector
        string_to_array(trim(both '[]' from (p.json_data ->> 'full_profile_embedding')), ',')::FLOAT4[]::VECTOR(1536) AS full_profile_embedding,

        -- 🧮 Derived fields
        (
            coalesce(p.json_data ->> 'name', '') || ' ' ||
            coalesce(p.json_data ->> 'headline', '') || ' ' ||
            coalesce(p.json_data ->> 'current_title', '') || ' ' ||
            coalesce(p.json_data ->> 'current_company', '') || ' ' ||
            coalesce(p.json_data ->> 'location', '') || ' ' ||
            coalesce(p.json_data ->> 'industry', '') || ' ' ||
            coalesce(p.json_data ->> 'summary', '') || ' ' ||
            coalesce(
                CASE 
                    WHEN jsonb_typeof(p.json_data -> 'skills') = 'array' 
                    THEN array_to_string(ARRAY(SELECT jsonb_array_elements_text(p.json_data -> 'skills')), ' ')
                    ELSE ''
                END, 
                ''
            )
        ) AS embedding_source_text,

        to_tsvector('simple',
            coalesce(p.json_data ->> 'name', '') || ' ' ||
            coalesce(p.json_data ->> 'headline', '') || ' ' ||
            coalesce(p.json_data ->> 'summary', '') || ' ' ||
            coalesce(array_to_string(ARRAY(SELECT jsonb_array_elements_text(p.json_data -> 'skills')), ' '), '')
        ) AS search_terms

    FROM jsonb_array_elements(in_crust_profiles->'profiles') AS p(json_data);

    -- 2️⃣ Upsert into raw_profile
    INSERT INTO people_schema.raw_profile (
        person_id, first_name, last_name, name, headline, region,
        current_title, current_company, industry, summary, skills, languages,
        profile_language, current_employers, past_employers, education_background,
        certifications, honors, last_updated
    )
    SELECT
        person_id, first_name, last_name, name, headline, region,
        current_title, current_company, industry, summary, skills, languages,
        profile_language, current_employers, past_employers, education_background,
        certifications, honors, NOW()
    FROM tmp_raw_profiles
    ON CONFLICT (person_id)
    DO UPDATE SET
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        name = EXCLUDED.name,
        headline = EXCLUDED.headline,
        region = EXCLUDED.region,
        current_title = EXCLUDED.current_title,
        current_company = EXCLUDED.current_company,
        industry = EXCLUDED.industry,
        summary = EXCLUDED.summary,
        skills = EXCLUDED.skills,
        languages = EXCLUDED.languages,
        profile_language = EXCLUDED.profile_language,
        current_employers = EXCLUDED.current_employers,
        past_employers = EXCLUDED.past_employers,
        education_background = EXCLUDED.education_background,
        certifications = EXCLUDED.certifications,
        honors = EXCLUDED.honors,
        last_updated = NOW();

    -- 3️⃣ Insert DB profiles into sourcing_session_profile
    INSERT INTO sourcing_schema.sourcing_session_profile (
        sourcing_session_id, raw_profile_id, score, order_index, note,
        report_summary, is_short_listed, summary_bullets, justification
    )
    SELECT
        in_sourcing_session_id, (r ->> 'person_id')::INT, 0, 1, '', '', FALSE, '{}', ''
    FROM jsonb_array_elements(in_db_profiles) AS r
    ON CONFLICT ON CONSTRAINT sourcing_session_profile_pkey DO NOTHING;

    -- 4️⃣ Insert crust profiles into sourcing_session_profile
    INSERT INTO sourcing_schema.sourcing_session_profile (
        sourcing_session_id, raw_profile_id, score, order_index, note,
        report_summary, is_short_listed, summary_bullets, justification
    )
    SELECT
        in_sourcing_session_id, person_id::INT, 0, 1, '', '', FALSE, '{}', ''
    FROM tmp_raw_profiles
    ON CONFLICT ON CONSTRAINT sourcing_session_profile_pkey DO NOTHING;

    -- 5️⃣ Drop temp table
    DROP TABLE IF EXISTS tmp_raw_profiles;

END;
$$;
