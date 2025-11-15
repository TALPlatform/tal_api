-- DROP and recreate clean schema (overrides current)
DROP SCHEMA IF EXISTS people_schema CASCADE;
CREATE SCHEMA people_schema;

-- Extensions (Supabase usually provides these; skip if already enabled)
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;  -- for trigram fuzzy search

-------------------------------
-- Lookup / normalized tables
-------------------------------
CREATE TABLE people_schema.location (
    location_id SERIAL PRIMARY KEY,
    region_name VARCHAR(255) NOT NULL UNIQUE,
    country VARCHAR(100),
    state VARCHAR(100),
    city VARCHAR(100)
);


CREATE TABLE people_schema.industry (
    industry_id SERIAL PRIMARY KEY,
    industry_name VARCHAR(255) NOT NULL UNIQUE
);
-------------------------------
-- Person (denormalized & vector)
-------------------------------
CREATE TABLE people_schema.person (
    person_id BIGSERIAL PRIMARY KEY,            
    fsd_profile_id TEXT UNIQUE,                       
    name TEXT NOT NULL,
    headline TEXT,
    summary TEXT,
    interests JSONB,
    projects JSONB,
    publications JSONB,
    volunteering JSONB,
    websites JSONB,
    title TEXT,
    company TEXT,
    location TEXT,
    industry TEXT,
    location_id INTEGER REFERENCES people_schema.location(location_id),
    embedding VECTOR(1536),
    embedding_text TEXT,
    embedding_model TEXT DEFAULT 'text-embedding-3-large',
    embedding_updated_at TIMESTAMPTZ,
    search_tsvector TSVECTOR,
    version TEXT,
    run_id BIGINT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE people_schema.skill (
    skill_id SERIAL PRIMARY KEY,
    skill_name VARCHAR(255) NOT NULL UNIQUE
);
CREATE TABLE people_schema.organization (
    organization_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    entity_type VARCHAR(50) NOT NULL DEFAULT 'company',
    linkedin_id VARCHAR(100),
    company_id BIGINT,
    website_domain VARCHAR(255),
    description TEXT,
    logo_url TEXT,
    business_emails JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT organization_name_entity_type_unique UNIQUE (name, entity_type)
);

CREATE TABLE people_schema.certification_type (
    name VARCHAR(255) NOT NULL,
    issuer_organization_id INTEGER REFERENCES people_schema.organization(organization_id),
    PRIMARY KEY (name, issuer_organization_id)
);

-------------------------------
-- Employment & education & others (normalized)
-------------------------------
CREATE TABLE people_schema.job_title (
    title VARCHAR(255) NOT NULL PRIMARY KEY
);
CREATE TABLE people_schema.honor(
    title VARCHAR(255) NOT NULL,
    description TEXT,
    issuer_organization_id INTEGER REFERENCES people_schema.organization(organization_id),
    PRIMARY KEY (title, issuer_organization_id)
);



CREATE TABLE people_schema.person_skill (
    person_id BIGINT REFERENCES people_schema.person(person_id) ON DELETE CASCADE,
    skill_id INTEGER REFERENCES people_schema.skill(skill_id) ON DELETE CASCADE,
    PRIMARY KEY (person_id, skill_id)
);
CREATE TABLE people_schema.organization_industry (
    organization_id INTEGER REFERENCES people_schema.organization(organization_id) ON DELETE CASCADE,
    industry_id INTEGER REFERENCES people_schema.industry(industry_id) ON DELETE CASCADE,
    PRIMARY KEY (organization_id, industry_id)
);
CREATE TABLE people_schema.person_experience (
    person_id BIGINT NOT NULL REFERENCES people_schema.person(person_id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL REFERENCES people_schema.job_title(title) ON DELETE CASCADE,
    organization_id INTEGER,
    start_date DATE,
    end_date DATE,
    is_current BOOLEAN,
    description TEXT,
    location VARCHAR(255),
    seniority_level VARCHAR(50),
    function_category VARCHAR(100),
    employer_is_default BOOLEAN,
    PRIMARY KEY (person_id, title, organization_id)
  
);;

CREATE TABLE people_schema.education_entry (
    person_id BIGINT REFERENCES people_schema.person(person_id) NOT NULL,
    organization_id INTEGER REFERENCES people_schema.organization(organization_id),
    degree_name VARCHAR(255),
    field_of_study VARCHAR(255),
    activities_and_societies TEXT,
    start_date DATE,
    end_date DATE,
    PRIMARY KEY (person_id, organization_id, degree_name)
);

CREATE TABLE people_schema.person_email (
    email_address VARCHAR(255) NOT NULL PRIMARY KEY,
    person_id BIGINT REFERENCES people_schema.person(person_id) NOT NULL,
    email_type VARCHAR(50),
    is_verified BOOLEAN
);

CREATE TABLE people_schema.person_certification (
    person_id BIGINT NOT NULL REFERENCES people_schema.person(person_id) ON DELETE CASCADE,
    certification_name VARCHAR(255) NOT NULL,
    certification_issuer_organization_id INTEGER NOT NULL,
    issued_date DATE,
    expiration_date DATE,
    PRIMARY KEY (person_id,certification_name,certification_issuer_organization_id),
    FOREIGN KEY (certification_name, certification_issuer_organization_id)
        REFERENCES people_schema.certification_type(name, issuer_organization_id)
        ON DELETE CASCADE
);

CREATE TABLE people_schema.person_honor (
    person_id BIGINT NOT NULL REFERENCES people_schema.person(person_id) ON DELETE CASCADE,
    honor_title VARCHAR(255) NOT NULL,
    honor_issuer_organization_id INTEGER NOT NULL,
    issued_date DATE,
    media_urls TEXT[],
    PRIMARY KEY (person_id, honor_title, honor_issuer_organization_id, issued_date),
    FOREIGN KEY (honor_title, honor_issuer_organization_id)
        REFERENCES people_schema.honor(title, issuer_organization_id)
        ON DELETE CASCADE
);
-- 2. Create the raw_profile table
-- 2. Create the raw_profile table (Keep 3072 dimensions)





CREATE TABLE people_schema.raw_profile (
    person_id BIGINT PRIMARY KEY,
    embedding_source_text TEXT,
    full_profile_embedding VECTOR(1536),
    embedding_model TEXT DEFAULT 'gemini-embedding-001',
    embedding_updated_at TIMESTAMP,
    search_terms TSVECTOR,
    first_name TEXT,
    last_name TEXT,
    flagship_profile_url TEXT,
    linkedin_profile_url TEXT,
    region TEXT,
    region_address_components TEXT,
    name TEXT,
    headline TEXT,
    location TEXT,
    current_title TEXT,
    current_company TEXT,
    industry TEXT,
    summary TEXT,
    experience JSONB,
    education JSONB,
    skills TEXT[],
    certifications JSONB,
    languages TEXT[],
    profile_language TEXT,
    profile_picture_url TEXT,
    profile_picture_permalink TEXT,
    twitter_handle TEXT,
    open_to_cards JSONB[],
    num_of_connections int,
    education_background JSONB[],
    current_employers JSONB[],
    past_employers JSONB[],
    last_updated TIMESTAMP,
    recently_changed_jobs BOOLEAN,
    years_of_experience TEXT,
    years_of_experience_raw INT,
    all_employers JSONB[],
    honors JSONB,
    interests JSONB,
    projects JSONB,
    publications JSONB,
    volunteering JSONB,
    websites JSONB,
    version TEXT,
    run_id BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
-- Index for Fuzzy/Full-Text Search on 'search_terms'
CREATE INDEX idx_raw_profile_search_terms ON people_schema.raw_profile USING gin (search_terms);

-- Index for Semantic Search on 'full_profile_embedding'
CREATE INDEX idx_raw_profile_embedding_hnsw
  ON people_schema.raw_profile USING hnsw (full_profile_embedding vector_cosine_ops)
  WITH (M = 16, ef_construction = 128); -- HNSW parameters remain optimal

CREATE INDEX idx_raw_profile_skills_gin ON people_schema.raw_profile USING gin (skills);

-- Create indexes to improve search performance
CREATE INDEX IF NOT EXISTS idx_raw_profile_experience 
    ON people_schema.raw_profile(years_of_experience_raw);

CREATE INDEX IF NOT EXISTS idx_raw_profile_location 
    ON people_schema.raw_profile USING gin(to_tsvector('english', COALESCE(location, '') || ' ' || COALESCE(region, '')));

CREATE INDEX IF NOT EXISTS idx_raw_profile_title 
    ON people_schema.raw_profile USING gin(to_tsvector('english', COALESCE(current_title, '') || ' ' || COALESCE(headline, '')));

CREATE INDEX IF NOT EXISTS idx_raw_profile_company 
    ON people_schema.raw_profile(current_company);

CREATE INDEX IF NOT EXISTS idx_raw_profile_industry 
    ON people_schema.raw_profile(industry);

CREATE INDEX IF NOT EXISTS idx_raw_profile_skills 
    ON people_schema.raw_profile USING gin(skills);

CREATE INDEX IF NOT EXISTS idx_raw_profile_languages 
    ON people_schema.raw_profile USING gin(languages);

CREATE INDEX IF NOT EXISTS idx_raw_profile_search_terms 
    ON people_schema.raw_profile USING gin(search_terms);

CREATE INDEX IF NOT EXISTS idx_raw_profile_all_employers 
    ON people_schema.raw_profile USING gin(all_employers);

CREATE INDEX IF NOT EXISTS idx_raw_profile_education 
    ON people_schema.raw_profile USING gin(education_background);
