# Sourcing & People Management API

## Overview

This document covers the new sourcing and people management endpoints that enable AI-powered talent search and candidate pipeline management. These endpoints allow you to create sourcing sessions, search profiles using both text and vector-based AI search, and manage candidate evaluations within sourcing projects.

## Sourcing Session Management

### Create or Update a Sourcing Session

Use this endpoint to create new sourcing sessions or update existing ones. Sourcing sessions organize your candidate searches within specific projects.

**Endpoint**: `SourcingSessionCreateUpdate`

**Example: Creating a New Session**

```json
{
  "sourcing_session_name": "Q4 2024 Senior Backend Engineers",
  "initial_filters": "{\"query\": \"senior backend engineer\", \"industries\": [\"Technology\", \"SaaS\"], \"locations\": [\"San Francisco\", \"New York\", \"Remote\"], \"skills\": [\"Go\", \"PostgreSQL\", \"gRPC\"], \"companies\": [\"Google\", \"Meta\", \"Amazon\"], \"projects\": [\"microservices\", \"distributed systems\"]}",
  "sourcing_project_id": 1,
  "created_by": 42
}
```

**Response**:

```json
{
  "sourcing_session_id": 123,
  "sourcing_session_name": "Q4 2024 Senior Backend Engineers",
  "initial_filters": "{\"query\": \"senior backend engineer\", \"industries\": [\"Technology\", \"SaaS\"], \"locations\": [\"San Francisco\", \"New York\", \"Remote\"], \"skills\": [\"Go\", \"PostgreSQL\", \"gRPC\"], \"companies\": [\"Google\", \"Meta\", \"Amazon\"], \"projects\": [\"microservices\", \"distributed systems\"]}",
  "sourcing_project_id": 1,
  "created_by": 42,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "deleted_at": ""
}
```

**Example: Updating an Existing Session**

```json
{
  "sourcing_session_id": 123,
  "sourcing_session_name": "Q4 2024 Senior Backend Engineers - Updated",
  "initial_filters": "{\"query\": \"senior backend engineer\", \"industries\": [\"Technology\", \"SaaS\", \"FinTech\"], \"locations\": [\"San Francisco\", \"New York\", \"Remote\", \"Europe\"], \"skills\": [\"Go\", \"PostgreSQL\", \"gRPC\", \"Kubernetes\"], \"companies\": [\"Google\", \"Meta\", \"Amazon\", \"Stripe\"]}",
  "sourcing_project_id": 1
}
```

### Get Sourcing Session Details

Retrieve complete session information including all associated candidate profiles with their evaluation scores and notes.

**Endpoint**: `SourcingSessionFind`

**Example Request**:

```json
{
  "sourcing_session_id": 2
}
```

**Example Response**:

```json
{
  "sourcing_session_id": 2,
  "sourcing_session_name": "AI/ML Engineering Talent Pool",
  "initial_filters": "{\"skills\": [\"machine learning\", \"python\", \"tensorflow\"], \"industries\": [\"Technology\", \"AI\"]}",
  "session_created_at": "2024-01-10T14:30:00Z",
  "session_updated_at": "2024-01-12T09:15:00Z",
  "sourcing_project_id": 5,
  "sourcing_project_name": "AI Product Development",
  "sourcing_project_breif": "Build AI-powered features for enterprise customers",
  "tenant_id": 1,
  "created_by": 42,
  "creator_name": "Sarah Johnson",
  "creator_email": "sarah@company.com",
  "profiles": [
    {
      "sourcing_session_id": 2,
      "raw_profile_id": 1586,
      "score": 85,
      "order_index": 1,
      "note": "Strong candidate with extensive experience in AI and ML",
      "report_summary": "Abhishek has 12+ years in tech with focus on AI adoption. Built and scaled products in education, finance, and e-commerce. Currently working on AI commercial adoption.",
      "is_short_listed": true,
      "summary_bullets": [
        "12+ years of tech industry experience",
        "Built and scaled products across multiple sectors",
        "Strong focus on AI/ML and commercial adoption",
        "Experience with education, finance, and e-commerce",
        "Leadership experience in product development"
      ],
      "justification": "Excellent match for the role due to extensive AI/ML experience and proven track record in building scalable products. Strong technical skills in Java, C++, Python, and machine learning. Leadership experience makes them suitable for senior positions.",
      "profile": {
        "person_id": 1586,
        "name": "Abhishek Kumar",
        "headline": "Senior AI Engineer | Machine Learning Specialist",
        "location": "San Francisco, California",
        "current_title": "Lead AI Engineer",
        "current_company": "Tech Innovations Inc",
        "industry": "Artificial Intelligence",
        "summary": "Experienced AI engineer with 12+ years building machine learning solutions across multiple industries. Specialized in computer vision and natural language processing.",
        "skills": [
          "python",
          "tensorflow",
          "machine learning",
          "deep learning",
          "java",
          "c++"
        ],
        "years_of_experience": "12+ years",
        "num_of_connections": 850,
        "profile_picture_url": "https://example.com/profiles/abhishek.jpg",
        "linkedin_profile_url": "https://linkedin.com/in/abhishekkumar"
      }
    }
  ]
}
```

### Manage Session Candidates

Add candidates to sourcing sessions or update their evaluation scores and notes.

**Endpoint**: `SourcingSessionProfileCreateUpdate`

**Example Request**:

```json
{
  "sourcing_session_id": 2,
  "raw_profile_id": 1586,
  "score": 85,
  "order_index": 1,
  "note": "Strong candidate with extensive experience in AI and ML",
  "report_summary": "Abhishek has 12+ years in tech with focus on AI adoption. Built and scaled products in education, finance, and e-commerce. Currently working on AI commercial adoption.",
  "is_short_listed": true,
  "summary_bullets": [
    "12+ years of tech industry experience",
    "Built and scaled products across multiple sectors",
    "Strong focus on AI/ML and commercial adoption",
    "Experience with education, finance, and e-commerce",
    "Leadership experience in product development"
  ],
  "justification": "Excellent match for the role due to extensive AI/ML experience and proven track record in building scalable products. Strong technical skills in Java, C++, Python, and machine learning. Leadership experience makes them suitable for senior positions."
}
```

**Response**:

```json
{
  "sourcing_session_id": 2,
  "raw_profile_id": 1586,
  "score": 85,
  "order_index": 1,
  "note": "Strong candidate with extensive experience in AI and ML",
  "report_summary": "Abhishek has 12+ years in tech with focus on AI adoption. Built and scaled products in education, finance, and e-commerce. Currently working on AI commercial adoption.",
  "is_short_listed": true,
  "summary_bullets": [
    "12+ years of tech industry experience",
    "Built and scaled products across multiple sectors",
    "Strong focus on AI/ML and commercial adoption",
    "Experience with education, finance, and e-commerce",
    "Leadership experience in product development"
  ],
  "justification": "Excellent match for the role due to extensive AI/ML experience and proven track record in building scalable products. Strong technical skills in Java, C++, Python, and machine learning. Leadership experience makes them suitable for senior positions."
}
```

## Project Reference Data

### Get Available Projects

Retrieve the list of sourcing projects for dropdown selection when creating new sessions.

**Endpoint**: `ProjectInputList`

**Example Request**:

```json
{}
```

**Example Response**:

```json
{
  "options": [
    {
      "value": "1",
      "label": "Backend Team Expansion 2024",
      "note": "Active - 3 sessions, 45 candidates"
    },
    {
      "value": "2",
      "label": "Frontend Development Q4",
      "note": "Planning - 0 sessions"
    },
    {
      "value": "5",
      "label": "AI Product Development",
      "note": "Active - 2 sessions, 28 candidates"
    },
    {
      "value": "8",
      "label": "Data Science Hiring",
      "note": "Completed - 5 hires"
    }
  ]
}
```

## AI-Powered Profile Search

### Search Candidate Profiles

This powerful endpoint combines traditional text search with AI-powered vector similarity search to find the most relevant candidates.

**Endpoint**: `RawProfileList`

**Example: Text-Based Search**

```json
{
  "query": "human resources",
  "skills": ["human resources"],
  "sourcing_session_id": 2,
  "limit": 5
}
```

**Example: Comprehensive Search with Filters**

```json
{
  "query": "senior backend engineer distributed systems",
  "industries": ["Technology", "SaaS", "FinTech"],
  "locations": ["San Francisco", "New York", "Remote", "Europe"],
  "skills": ["Go", "PostgreSQL", "gRPC", "Kubernetes", "Docker"],
  "companies": ["Google", "Stripe", "Amazon", "Startups"],
  "projects": ["microservices", "distributed systems", "payment processing"],
  "sourcing_session_id": 2,
  "limit": 20
}
```

**Example Response**:

```json
{
  "records": [
    {
      "person_id": 1586,
      "name": "Abhishek Kumar",
      "headline": "Senior AI Engineer | Machine Learning Specialist",
      "location": "San Francisco, California",
      "current_title": "Lead AI Engineer",
      "current_company": "Tech Innovations Inc",
      "industry": "Artificial Intelligence",
      "summary": "Experienced AI engineer with 12+ years building machine learning solutions across multiple industries. Specialized in computer vision and natural language processing.",
      "years_of_experience": "12+ years",
      "num_of_connections": 850,
      "profile_picture_url": "https://example.com/profiles/abhishek.jpg",
      "linkedin_profile_url": "https://linkedin.com/in/abhishekkumar",
      "skills": [
        "python",
        "tensorflow",
        "machine learning",
        "deep learning",
        "java",
        "c++"
      ],
      "semantic_score": 0.89,
      "text_rank": 0.76,
      "hybrid_score": 0.85
    },
    {
      "person_id": 2045,
      "name": "Maria Rodriguez",
      "headline": "HR Director | Talent Acquisition Leader",
      "location": "New York, New York",
      "current_title": "Director of Human Resources",
      "current_company": "Global Tech Corp",
      "industry": "Human Resources",
      "summary": "Strategic HR leader with 15+ years experience in talent acquisition and employee development.",
      "years_of_experience": "15+ years",
      "num_of_connections": 1200,
      "profile_picture_url": "https://example.com/profiles/maria.jpg",
      "linkedin_profile_url": "https://linkedin.com/in/mariarodriguez",
      "skills": [
        "human resources",
        "talent acquisition",
        "employee engagement",
        "recruitment",
        "HR strategy"
      ],
      "semantic_score": 0.82,
      "text_rank": 0.88,
      "hybrid_score": 0.84
    }
  ]
}
```

## Key Features

### Hybrid Search Technology

- **Text Search**: Traditional keyword matching across profiles
- **Vector Search**: AI-powered semantic understanding using 1536-dimensional embeddings
- **Smart Ranking**: 70% semantic similarity + 30% text relevance = hybrid score
- **Structured Filtering**: Combine with industries, locations, skills, companies

### Session Integration

- Automatically save search results to sourcing sessions using `sourcing_session_id`
- Track candidate evaluations, scores, and notes
- Organize candidates across different projects and searches

### External Data Enrichment

- Integrates with external data sources for expanded candidate pool
- Automatic synchronization of external profiles
- Bulk profile creation and updates

## Usage Notes

- Use `sourcing_session_id` parameter to automatically save search results to sessions
- Combine text queries with structured filters for optimal results
- The hybrid scoring system (semantic_score + text_rank = hybrid_score) helps prioritize the most relevant candidates
- All sourcing operations are tied to projects for better organization and reporting

---

## AI-Powered Profile Request Builder

These endpoints help construct structured search requests for sourcing sessions using AI-assisted parsing of text prompts. They allow you to infer skills, locations, companies, industries, and projects from a natural language query.

### Build a Raw Profile List Request

**Endpoint**: `RawProfileListRequestBuild`

Use this endpoint to generate a **structured request** for the `RawProfileList` search. The AI analyzes the text prompt and fills the fields: `industries`, `locations`, `skills`, `companies`, and `projects`.

**Example Request**:

```json
{
  "text": "senior js engineer with 5 years and on NYc"
}
```

**Example Response**:

```json
{
  "structured_response": {
    "industries": [],
    "locations": ["New York City"],
    "skills": ["javascript"],
    "companies": [],
    "projects": []
  }
}
```

**Example Request with Fuzzy References**:

```json
{
  "text": "FANG frontend dev with Python experience, remote"
}
```

**Example Response**:

```json
{
  "structured_response": {
    "industries": [],
    "locations": ["Remote"],
    "skills": ["python"],
    "companies": ["Facebook", "Amazon", "Netflix", "Google"],
    "projects": []
  }
}
```

**Notes / Instructions for Users**:

1. The `text` field can include natural language prompts describing the candidate profile, skills, locations, companies, or projects.
2. The response **always contains only the five fields**: `industries`, `locations`, `skills`, `companies`, `projects`.
3. AI normalizes abbreviations and fuzzy terms:

   - `"js"` → `"javascript"`
   - `"NYC"` → `"New York City"`
   - `"FANG"` → `["Facebook", "Amazon", "Netflix", "Google"]`
   - `"FE"` → `"Frontend"`, `"BE"` → `"Backend"`

**Usage**:

- Typically used before calling `RawProfileList` to pre-fill structured filters from a natural-language prompt.
- Works best when combined with a `sourcing_session_id` to automatically save structured searches.
  **Reference**: [Current agent instructions](https://docs.google.com/document/d/18mT4n5yEF8wKSTYVA0egWkK7lYMn1bsk2bdyfZexRBI/edit?usp=sharing)

---

### Find a Raw Profile by Person ID

**Endpoint**: `RawProfileFind`

Retrieve detailed information for a specific candidate using their `person_id`.

**Example Request**:

```json
{
  "person_id": 1586
}
```

**Example Response**:

```json
{
  "person_id": 1586,
  "name": "Abhishek Kumar",
  "headline": "Senior AI Engineer | Machine Learning Specialist",
  "location": "San Francisco, California",
  "current_title": "Lead AI Engineer",
  "current_company": "Tech Innovations Inc",
  "industry": "Artificial Intelligence",
  "summary": "Experienced AI engineer with 12+ years building machine learning solutions across multiple industries. Specialized in computer vision and natural language processing.",
  "skills": [
    "python",
    "tensorflow",
    "machine learning",
    "deep learning",
    "java",
    "c++"
  ],
  "years_of_experience": "12+ years",
  "num_of_connections": 850,
  "profile_picture_url": "https://example.com/profiles/abhishek.jpg",
  "linkedin_profile_url": "https://linkedin.com/in/abhishekkumar"
}
```

**Notes**:

- This endpoint is useful to fetch a single candidate’s full profile after filtering candidates via `RawProfileList` or `RawProfileListRequestBuild`.

---

I can also **add a new "How to use RawProfileListRequestBuild + RawProfileList workflow" section** that shows the full workflow:

1. Build structured request from text prompt
2. Use it to call `RawProfileList`
3. Save results to a sourcing session

This will make the doc much clearer for users.

Do you want me to add that workflow section next?
