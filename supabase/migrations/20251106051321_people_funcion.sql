DROP function IF EXISTS people_schema.people_ist;
create function people_schema.people_list(
    p_prompt text default null,
    p_location text default null,
    p_job_title text default null,
    p_years_of_experience_min numeric default null,
    p_years_of_experience_max numeric default null,
    p_industry text default null,
    p_skills text[] default null,
    p_company text default null,
    p_limit int default 20
)
RETURNS TABLE (
    person_id INT,
    name TEXT,
    headline TEXT,
    current_title TEXT,
    current_company_name TEXT,
    location TEXT,
    industry TEXT,
    score FLOAT
)
LANGUAGE plpgsql AS $$
DECLARE
    v_query_vector vector(1536);
BEGIN

    RETURN QUERY
    WITH base AS (
        SELECT
            p.person_id,
            p.name,
            p.headline,
            -- p.current_title,
            p.current_company_name,
            l.region_name AS location,
            p.industry_name AS industry,
            -- Hybrid score: combine vector similarity + FTS rank
            (
                COALESCE(1 - (p.embedding <=> v_query_vector), 0) * 0.7 +
                COALESCE(ts_rank_cd(p.search_tsvector, plainto_tsquery('simple', p_prompt)), 0) * 0.3
            ) AS score
        FROM people_schema.person p
        LEFT JOIN people_schema.location l ON p.current_location_id = l.location_id
        WHERE
            -- 2️⃣ Prompt relevance: only filter if we have a prompt
            (
                p_prompt IS NULL OR
                p.search_tsvector @@ plainto_tsquery('simple', p_prompt)
            )
            -- 3️⃣ Structured filters
            AND (p_location IS NULL OR l.region_name ILIKE '%' || p_location || '%')
            AND (p_job_title IS NULL OR p.current_title ILIKE '%' || p_job_title || '%')
            AND (p_industry IS NULL OR p.industry_name ILIKE '%' || p_industry || '%')
            AND (p_company IS NULL OR p.current_company_name ILIKE '%' || p_company || '%')
            AND (p_years_of_experience_min IS NULL OR p.years_of_experience_raw >= p_years_of_experience_min)
            AND (p_years_of_experience_max IS NULL OR p.years_of_experience_raw <= p_years_of_experience_max)
            AND (
                v_query_vector IS NULL 
                OR 1 - (p.embedding <=> v_query_vector) > 0.65
            )
               )
    SELECT *
    FROM base
    ORDER BY score DESC NULLS LAST
    LIMIT p_limit;
END;
$$;
DROP FUNCTION IF EXISTS people_schema.raw_profiles_bulk_create_update;
CREATE OR REPLACE FUNCTION people_schema.raw_profiles_bulk_create_update(in_session_id INT,in_data JSONB)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    -- 1️⃣ Create a temporary staging table
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

    FROM jsonb_array_elements(in_data) AS p(json_data);

    -- 2️⃣ Upsert into main table
    INSERT INTO people_schema.raw_profile AS rp (
        person_id, embedding_model, embedding_updated_at, first_name, last_name,
        flagship_profile_url, linkedin_profile_url, region, region_address_components,
        name, headline, location, current_title, current_company, industry, summary,
        skills, languages, certifications, education, honors, interests, projects,
        publications, volunteering, websites, open_to_cards, profile_picture_url,
        profile_picture_permalink, twitter_handle, num_of_connections,
        recently_changed_jobs, years_of_experience, years_of_experience_raw,
        profile_language, version, run_id, last_updated,
        full_profile_embedding, embedding_source_text, search_terms,
        education_background, current_employers, past_employers, all_employers,
        updated_at
    )
    SELECT
        person_id, embedding_model, embedding_updated_at, first_name, last_name,
        flagship_profile_url, linkedin_profile_url, region, region_address_components,
        name, headline, location, current_title, current_company, industry, summary,
        skills, languages, certifications, education, honors, interests, projects,
        publications, volunteering, websites, open_to_cards, profile_picture_url,
        profile_picture_permalink, twitter_handle, num_of_connections,
        recently_changed_jobs, years_of_experience, years_of_experience_raw,
        profile_language, version, run_id, last_updated,
        full_profile_embedding, embedding_source_text, search_terms,
        education_background, current_employers, past_employers, all_employers,
        NOW()
    FROM tmp_raw_profiles
    ON CONFLICT (person_id)
    DO UPDATE SET
        embedding_model = EXCLUDED.embedding_model,
        embedding_updated_at = EXCLUDED.embedding_updated_at,
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        flagship_profile_url = EXCLUDED.flagship_profile_url,
        linkedin_profile_url = EXCLUDED.linkedin_profile_url,
        region = EXCLUDED.region,
        region_address_components = EXCLUDED.region_address_components,
        name = EXCLUDED.name,
        headline = EXCLUDED.headline,
        location = EXCLUDED.location,
        current_title = EXCLUDED.current_title,
        current_company = EXCLUDED.current_company,
        industry = EXCLUDED.industry,
        summary = EXCLUDED.summary,
        skills = EXCLUDED.skills,
        languages = EXCLUDED.languages,
        certifications = EXCLUDED.certifications,
        education = EXCLUDED.education,
        honors = EXCLUDED.honors,
        interests = EXCLUDED.interests,
        projects = EXCLUDED.projects,
        publications = EXCLUDED.publications,
        volunteering = EXCLUDED.volunteering,
        websites = EXCLUDED.websites,
        open_to_cards = EXCLUDED.open_to_cards,
        education_background = EXCLUDED.education_background,
        current_employers = EXCLUDED.current_employers,
        past_employers = EXCLUDED.past_employers,
        all_employers = EXCLUDED.all_employers,
        profile_picture_url = EXCLUDED.profile_picture_url,
        profile_picture_permalink = EXCLUDED.profile_picture_permalink,
        twitter_handle = EXCLUDED.twitter_handle,
        num_of_connections = EXCLUDED.num_of_connections,
        recently_changed_jobs = EXCLUDED.recently_changed_jobs,
        years_of_experience = EXCLUDED.years_of_experience,
        years_of_experience_raw = EXCLUDED.years_of_experience_raw,
        profile_language = EXCLUDED.profile_language,
        version = EXCLUDED.version,
        run_id = EXCLUDED.run_id,
        last_updated = EXCLUDED.last_updated,
        full_profile_embedding = EXCLUDED.full_profile_embedding,
        embedding_source_text = EXCLUDED.embedding_source_text,
        search_terms = EXCLUDED.search_terms,
        updated_at = NOW();
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
    ) SELECT in_session_id , p.person_id , 0,1,'','',FALSE,array[]::TEXT[],'' FROM tmp_raw_profiles p 
ON CONFLICT ON CONSTRAINT sourcing_session_profile_pkey DO NOTHING;

    -- 3️⃣ Drop temp table
    DROP TABLE IF EXISTS tmp_raw_profiles;
END;
$$;



DROP FUNCTION  IF EXISTS people_schema.raw_profile_search;
CREATE OR REPLACE FUNCTION people_schema.raw_profile_search(
    in_query TEXT DEFAULT NULL,
    in_embedding VECTOR(1536) DEFAULT NULL,
    in_industries TEXT[] DEFAULT NULL,
    in_locations TEXT[] DEFAULT NULL,
    in_skills TEXT[] DEFAULT NULL,
    in_companies TEXT[] DEFAULT NULL,
    in_projects TEXT[] DEFAULT NULL,
    in_limit INT DEFAULT 20
)
RETURNS TABLE (
    person_id BIGINT,
    name TEXT,
    headline TEXT,
    region TEXT,
    current_title TEXT,
    current_company TEXT,
    industry TEXT,
    summary TEXT,
    years_of_experience TEXT,
    num_of_connections INT,
    profile_picture_url TEXT,
    linkedin_profile_url TEXT,
    skills TEXT[],
    semantic_score FLOAT,
    text_rank FLOAT,
    hybrid_score FLOAT
)
LANGUAGE plpgsql
AS $$
DECLARE
    ts_query TSQUERY;
BEGIN
    -- 1️⃣ Convert text input into a tsquery (safe for NULL)
    ts_query := CASE 
        WHEN in_query IS NOT NULL AND in_query <> '' THEN plainto_tsquery('simple', in_query)
        ELSE NULL 
    END;

    -- 2️⃣ Hybrid query: semantic + text + structured filters (all nullable)
    RETURN QUERY
    WITH filtered AS (
        SELECT
            rp.person_id,
            rp.name,
            rp.headline,
            rp.region,
            rp.current_title,
            rp.current_company,
            rp.industry,
            rp.summary,
            rp.skills,
            rp.years_of_experience,
            rp.num_of_connections,
            rp.profile_picture_url,
            rp.linkedin_profile_url,
            CAST(
                CASE
                    WHEN in_embedding IS NOT NULL THEN 1 - (rp.full_profile_embedding <=> in_embedding)
                    ELSE 0
                END AS FLOAT
            ) AS semantic_score,
            CAST(
                CASE
                    WHEN ts_query IS NOT NULL THEN ts_rank_cd(rp.search_terms, ts_query)
                    ELSE 0
                END AS FLOAT
            ) AS text_rank
        FROM people_schema.raw_profile rp
        WHERE
            -- 🏭 INDUSTRY FILTER (nullable)
            (in_industries IS NULL OR cardinality(in_industries) = 0 OR rp.industry ILIKE ANY (
                ARRAY(SELECT '%' || i || '%' FROM unnest(in_industries) AS i)
            ))

            -- 🌍 LOCATION FILTER (nullable)
            AND (in_locations IS NULL OR cardinality(in_locations) = 0 OR rp.region ILIKE ANY (
                ARRAY(SELECT '%' || l || '%' FROM unnest(in_locations) AS l)
            ))

            -- 🧠 SKILLS FILTER (nullable)
            AND (in_skills IS NULL OR cardinality(in_skills) = 0 OR EXISTS (
                SELECT 1
                FROM unnest(rp.skills) s
                WHERE s ILIKE ANY (
                    ARRAY(SELECT '%' || sk || '%' FROM unnest(in_skills) AS sk)
                )
            ))

            -- 🏢 COMPANY FILTER (nullable)
            AND (in_companies IS NULL OR cardinality(in_companies) = 0 OR (
                rp.current_company ILIKE ANY (
                    ARRAY(SELECT '%' || c || '%' FROM unnest(in_companies) AS c)
                )
                OR EXISTS (
                    SELECT 1
                    FROM jsonb_array_elements(rp.experience) AS e
                    WHERE e ->> 'company_name' ILIKE ANY (
                        ARRAY(SELECT '%' || c || '%' FROM unnest(in_companies) AS c)
                    )
                )
            ))

            -- 🚀 PROJECTS FILTER (nullable)
            AND (in_projects IS NULL OR cardinality(in_projects) = 0 OR EXISTS (
                SELECT 1
                FROM jsonb_array_elements(rp.projects) AS p
                WHERE (p ->> 'name') ILIKE ANY (
                    ARRAY(SELECT '%' || prj || '%' FROM unnest(in_projects) AS prj)
                )
                OR (p ->> 'description') ILIKE ANY (
                    ARRAY(SELECT '%' || prj || '%' FROM unnest(in_projects) AS prj)
                )
            ))

            -- 🔍 FUZZY / FULL-TEXT FILTER
            -- AND (
            --     ts_query IS NULL OR rp.search_terms @@ ts_query
            -- )
    )
    SELECT
        f.person_id,
        coalesce(f.name, '') AS name,
        coalesce(f.headline, '') AS headline,
        coalesce(f.region, '') AS region,
        coalesce(f.current_title, '') AS current_title,
        coalesce(f.current_company, '') AS current_company,
        coalesce(f.industry, '') AS industry,
        coalesce(f.summary, '') AS summary,
        coalesce(f.years_of_experience, '') AS years_of_experience,
        coalesce(f.num_of_connections, 0) AS num_of_connections,
        coalesce(f.profile_picture_url, '') AS profile_picture_url,
        coalesce(f.linkedin_profile_url, '') AS linkedin_profile_url,
        coalesce(f.skills, ARRAY[]::text[]) AS skills,
        coalesce(f.semantic_score, 0) AS semantic_score,
        coalesce(f.text_rank, 0) AS text_rank,
        coalesce(ROUND((0.7 * f.semantic_score + 0.3 * f.text_rank)::numeric, 6) , 0)::FLOAT  AS hybrid_score
    FROM filtered f
    ORDER BY hybrid_score DESC
    LIMIT in_limit;
END;
$$;



DROP FUNCTION IF EXISTS people_schema.raw_profile_list;

CREATE OR REPLACE FUNCTION people_schema.raw_profile_list(
    in_min_experience INT DEFAULT NULL,
    in_max_experience INT DEFAULT NULL,
    in_required_contact_info TEXT[] DEFAULT NULL,
    in_timezone TEXT DEFAULT 'any',
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
    in_limit INT DEFAULT 50
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
BEGIN
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
        rp.years_of_experience_raw,
        rp.num_of_connections,
        rp.skills,
        rp.languages
    FROM people_schema.raw_profile rp
    WHERE 
        -- Experience range filter
        (in_min_experience IS NULL OR rp.years_of_experience_raw >= in_min_experience)
        AND (in_max_experience IS NULL OR rp.years_of_experience_raw <= in_max_experience)
        
        -- Required contact info (check for non-null fields)
        AND (
            in_required_contact_info IS NULL 
            OR (
                ('email' = ANY(in_required_contact_info) AND rp.flagship_profile_url IS NOT NULL)
                OR ('linkedin' = ANY(in_required_contact_info) AND rp.linkedin_profile_url IS NOT NULL)
                OR ('twitter' = ANY(in_required_contact_info) AND rp.twitter_handle IS NOT NULL)
            )
        )
        
        -- Location filter (case-insensitive partial match)
        AND (
            in_locations IS NULL 
            OR EXISTS (
                SELECT 1 
                FROM unnest(in_locations) AS loc
                WHERE rp.location ILIKE '%' || loc || '%' 
                   OR rp.region ILIKE '%' || loc || '%'
            )
        )
        
        -- Job titles filter (case-insensitive partial match on current_title and headline)
        AND (
            in_job_titles IS NULL 
            OR EXISTS (
                SELECT 1 
                FROM unnest(in_job_titles) AS title
                WHERE rp.current_title ILIKE '%' || title || '%'
                   OR rp.headline ILIKE '%' || title || '%'
            )
        )
        
        -- Job seniority (parsed from title/headline)
        AND (
            in_job_seniority IS NULL
            OR (
                in_job_seniority = 'entry' AND (
                    rp.current_title ~* '\y(junior|entry|associate|assistant|intern|trainee)\y'
                    OR rp.headline ~* '\y(junior|entry|associate|assistant|intern|trainee)\y'
                )
            )
            OR (
                in_job_seniority = 'mid' AND (
                    rp.current_title ~* '\y(specialist|coordinator|analyst|engineer|developer)\y'
                    AND NOT rp.current_title ~* '\y(senior|lead|principal|chief|head|director|vp|vice president)\y'
                )
            )
            OR (
                in_job_seniority = 'senior' AND (
                    rp.current_title ~* '\y(senior|lead|principal)\y'
                    OR rp.headline ~* '\y(senior|lead|principal)\y'
                )
            )
            OR (
                in_job_seniority = 'executive' AND (
                    rp.current_title ~* '\y(chief|ceo|cto|cfo|coo|vp|vice president|director|head of)\y'
                    OR rp.headline ~* '\y(chief|ceo|cto|cfo|coo|vp|vice president|director|head of)\y'
                )
            )
        )
        
        -- Companies filter (current and past)
        AND (
            in_companies IS NULL 
            OR EXISTS (
                SELECT 1 
                FROM unnest(in_companies) AS comp
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
            in_industries IS NULL 
            OR EXISTS (
                SELECT 1 
                FROM unnest(in_industries) AS ind
                WHERE rp.industry ILIKE '%' || ind || '%'
            )
        )
        
        -- Keywords search (full-text search on multiple fields)
        AND (
            in_keywords IS NULL 
            OR EXISTS (
                SELECT 1 
                FROM unnest(in_keywords) AS kw
                WHERE rp.search_terms @@ plainto_tsquery('english', kw)
                   OR rp.headline ILIKE '%' || kw || '%'
                   OR rp.summary ILIKE '%' || kw || '%'
                   OR rp.current_title ILIKE '%' || kw || '%'
            )
        )
        
        -- Skills filter (array overlap)
        AND (
            in_skills IS NULL 
            OR rp.skills && in_skills
        )
        
        -- Education levels (check education JSONB array)
        AND (
            in_education_levels IS NULL 
            OR EXISTS (
                SELECT 1 
                FROM unnest(rp.education_background) AS edu,
                     unnest(in_education_levels) AS level
                WHERE edu->>'degree' ILIKE '%' || level || '%'
            )
        )
        
        -- Languages filter
        AND (
            in_languages IS NULL 
            OR rp.languages && in_languages
        )
    
    ORDER BY 
        -- Prioritize profiles with more complete information
        (CASE WHEN rp.linkedin_profile_url IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN rp.summary IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN rp.years_of_experience_raw IS NOT NULL THEN 1 ELSE 0 END) DESC,
        rp.num_of_connections DESC NULLS LAST,
        rp.updated_at DESC
    
    LIMIT in_limit;
END;
$$;
