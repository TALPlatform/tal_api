package crustdata

import (
	"context"
	"encoding/json"
	"fmt"
)

// PeopleSearchRequest represents the request for people discovery
type PeopleSearchRequest struct {
	Filters        interface{}     `json:"filters"` // Can be FilterCondition or FilterGroup
	Cursor         string          `json:"cursor,omitempty"`
	Limit          int             `json:"limit,omitempty"`
	PostProcessing *PostProcessing `json:"post_processing,omitempty"`
	Preview        bool            `json:"preview,omitempty"`
}

// FilterCondition represents a single filter condition
type FilterCondition struct {
	Column string      `json:"column"`
	Type   string      `json:"type"` // "=", "!=", "in", "not_in", ">", "<", "=>", "=<", "(.)", "[.]"
	Value  interface{} `json:"value"`
}

// FilterGroup represents a group of filters with logical operator
type FilterGroup struct {
	Op         string        `json:"op"`         // "and", "or"
	Conditions []interface{} `json:"conditions"` // Can be FilterCondition or FilterGroup
}

// PostProcessing represents additional filtering options
type PostProcessing struct {
	ExcludeProfiles []string `json:"exclude_profiles,omitempty"`
	ExcludeNames    []string `json:"exclude_names,omitempty"`
}

// PeopleSearchResponse represents the response from people discovery
type PeopleSearchResponse struct {
	Profiles   []*PersonDBProfile `json:"profiles"`
	NextCursor string             `json:"next_cursor,omitempty"`
	TotalCount int                `json:"total_count"`
}

// PersonDBProfile represents a person profile from the search API
type PersonDBProfile struct {
	PersonID                int64           `json:"person_id"`
	Name                    string          `json:"name"`
	FirstName               string          `json:"first_name"`
	LastName                string          `json:"last_name"`
	Region                  string          `json:"region"`
	RegionAddressComponents []string        `json:"region_address_components,omitempty"`
	Headline                string          `json:"headline"`
	Summary                 string          `json:"summary"`
	Skills                  []string        `json:"skills"`
	Languages               []string        `json:"languages"`
	ProfileLanguage         string          `json:"profile_language"`
	Emails                  []string        `json:"emails"`
	TwitterHandle           string          `json:"twitter_handle"`
	OpenToCards             []bool          `json:"open_to_cards"`
	NumOfConnections        int             `json:"num_of_connections"`
	RecentlyChangedJobs     bool            `json:"recently_changed_jobs"`
	YearsOfExperience       string          `json:"years_of_experience"`
	YearsOfExperienceRaw    float64         `json:"years_of_experience_raw"`
	CurrentEmployers        []Employment    `json:"current_employers"`
	PastEmployers           []Employment    `json:"past_employers"`
	EducationBackground     []Education     `json:"education_background"`
	Honors                  []Honor         `json:"honors"`
	Certifications          []Certification `json:"certifications"`
}

// Employment represents employment information
type Employment struct {
	Name                               string   `json:"name"`
	LinkedInID                         string   `json:"linkedin_id"`
	CompanyID                          int64    `json:"company_id"`
	CompanyLinkedInProfileURL          string   `json:"company_linkedin_profile_url"`
	CompanyWebsiteDomain               string   `json:"company_website_domain"`
	PositionID                         string   `json:"position_id"`
	Title                              string   `json:"title"`
	Description                        string   `json:"description"`
	Location                           string   `json:"location"`
	StartDate                          string   `json:"start_date"`
	EndDate                            string   `json:"end_date"`
	EmployerIsDefault                  bool     `json:"employer_is_default"`
	BusinessEmailVerified              bool     `json:"business_email_verified"`
	CompanyHeadquartersCountry         string   `json:"company_headquarters_country"`
	CompanyHQLocation                  string   `json:"company_hq_location"`
	CompanyHQLocationAddressComponents []string `json:"company_hq_location_address_components"`
	CompanyHeadcountRange              string   `json:"company_headcount_range"`
	CompanyHeadcountLatest             int      `json:"company_headcount_latest"`
	CompanyIndustries                  []string `json:"company_industries"`
	CompanyType                        string   `json:"company_type"`
	SeniorityLevel                     string   `json:"seniority_level"`
	FunctionCategory                   string   `json:"function_category"`
	YearsAtCompany                     string   `json:"years_at_company"`
	YearsAtCompanyRaw                  float64  `json:"years_at_company_raw"`
}

// Education represents education background
type Education struct {
	DegreeName             string `json:"degree_name"`
	InstituteName          string `json:"institute_name"`
	InstituteLinkedInID    string `json:"institute_linkedin_id"`
	InstituteLinkedInURL   string `json:"institute_linkedin_url"`
	InstituteLogoURL       string `json:"institute_logo_url"`
	FieldOfStudy           string `json:"field_of_study"`
	ActivitiesAndSocieties string `json:"activities_and_societies"`
	StartDate              string `json:"start_date"`
	EndDate                string `json:"end_date"`
}

// Honor represents honors and awards
type Honor struct {
	Title                            string   `json:"title"`
	IssuedDate                       string   `json:"issued_date"`
	Description                      string   `json:"description"`
	Issuer                           string   `json:"issuer"`
	MediaURLs                        []string `json:"media_urls"`
	AssociatedOrganizationLinkedInID string   `json:"associated_organization_linkedin_id"`
	AssociatedOrganization           string   `json:"associated_organization"`
}

// Certification represents professional certifications
type Certification struct {
	Name                         string `json:"name"`
	IssuedDate                   string `json:"issued_date"`
	ExpirationDate               string `json:"expiration_date"`
	URL                          string `json:"url"`
	IssuerOrganization           string `json:"issuer_organization"`
	IssuerOrganizationLinkedInID string `json:"issuer_organization_linkedin_id"`
	CertificationID              string `json:"certification_id"`
}

// SearchPeople implements the people discovery API
func (c *CrustdataService) PeopleSearch(ctx context.Context, req *PeopleSearchRequest) (*PeopleSearchResponse, error) {
	path := "/screener/persondb/search"

	body, err := c.Client.doRequest(ctx, "POST", path, req)
	if err != nil {
		return nil, err
	}

	var response PeopleSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
