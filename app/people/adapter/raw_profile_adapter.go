package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

// FilterCondition defines the structure of a single filter condition for Crustdata
type FilterCondition struct {
	Column string `json:"column"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}

// FilterGroup defines logical groups of filters for Crustdata ("and" / "or")
type FilterGroup struct {
	Op         string        `json:"op"`
	Conditions []interface{} `json:"conditions"`
}

func (a *PeopleAdapter) RawProfileListEnrichAndMarshal(
	ctx context.Context,
	session_id int32,
	profiles *crustdata.PeopleSearchResponse,
) (*db.RawProfilesBulkCreateUpdateParams, error) {
	embeddingSourceTexts := make([]string, len(profiles.Profiles))
	for i, p := range profiles.Profiles {
		sourceParts := []string{
			p.Name,
			p.Headline,
			p.Region, // New key added
			p.Summary,
		}
		skillsStr := strings.Join(p.Skills, " ")
		sourceParts = append(sourceParts, skillsStr)
		embeddingSourceTexts[i] = strings.Join(sourceParts, " ")
	}
	enrichedProfiles := make([]map[string]interface{}, len(profiles.Profiles))
	// embeddings, err := a.embedBulk(ctx, embeddingSourceTexts, &genaiclient.EmbedOptions{Dimensions: 1536})
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	// }
	//
	for i, p := range profiles.Profiles {
		profileJSON, _ := json.Marshal(p)
		var profileMap map[string]interface{}
		if err := json.Unmarshal(profileJSON, &profileMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal profile to map: %w", err)
		}
		profileMap["embedding_source_text"] = embeddingSourceTexts[i]
		// profileMap["full_profile_embedding"] = pgvector.NewVector(embeddings[i])
		enrichedProfiles[i] = profileMap
	}

	jsonBytes, err := json.Marshal(enrichedProfiles)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal enriched profiles: %w", err)
	}

	return &db.RawProfilesBulkCreateUpdateParams{
		SessionID: session_id,
		Profiles:  jsonBytes,
	}, nil
}

// RawProfileListCrustDataFromGrpc converts RawProfileListRequest → Crustdata PeopleSearchRequest
func (a *PeopleAdapter) RawProfileListCrustDataFromGrpc(req *talv1.RawProfileListRequest) *crustdata.PeopleSearchRequest {
	var conditions []interface{}

	// 1. Job title / Query-related filters (if query provided)
	if len(req.Query) > 0 {
		jobConditions := []interface{}{}
		titles := []string{
			"backend developer", "frontend developer", "full stack developer",
			req.Query, // include raw query term too
		}

		for _, title := range titles {
			jobConditions = append(jobConditions,
				FilterCondition{Column: "current_employers.title", Type: "(.)", Value: title},
				FilterCondition{Column: "past_employers.title", Type: "(.)", Value: title},
			)
		}

		conditions = append(conditions, FilterGroup{
			Op:         "or",
			Conditions: jobConditions,
		})
	}

	// 2. Companies filter (current and past)
	if len(req.Companies) > 0 {
		companyConditions := []interface{}{}
		for _, c := range req.Companies {
			companyConditions = append(companyConditions,
				FilterCondition{Column: "current_employers.name", Type: "(.)", Value: c},
				FilterCondition{Column: "past_employers.name", Type: "(.)", Value: c},
			)
		}

		conditions = append(conditions, FilterGroup{
			Op:         "or",
			Conditions: companyConditions,
		})
	}

	// 3. Locations filter (region)
	if len(req.Locations) > 0 {
		locationConditions := []interface{}{}
		for _, l := range req.Locations {
			locationConditions = append(locationConditions,
				FilterCondition{Column: "region", Type: "(.)", Value: l},
			)
		}

		conditions = append(conditions, FilterGroup{
			Op:         "or",
			Conditions: locationConditions,
		})
	}

	// 4. Skills filter (AND all)
	if len(req.Skills) > 0 {
		skillConditions := []interface{}{}
		for _, s := range req.Skills {
			skillConditions = append(skillConditions,
				FilterCondition{Column: "skills", Type: "(.)", Value: s},
			)
		}

		conditions = append(conditions, FilterGroup{
			Op:         "and",
			Conditions: skillConditions,
		})
	}

	// 5. Industries filter (OR)
	if len(req.Industries) > 0 {
		industryConditions := []interface{}{}
		for _, i := range req.Industries {
			industryConditions = append(industryConditions,
				FilterCondition{Column: "industry", Type: "(.)", Value: i},
			)
		}

		conditions = append(conditions, FilterGroup{
			Op:         "or",
			Conditions: industryConditions,
		})
	}

	if req.YearsOfExperienceFrom > 0 || req.YearsOfExperienceTo > 0 {
		expConditions := []interface{}{}

		// Check for minimum experience (YearsOfExperienceFrom)
		if req.YearsOfExperienceFrom > 0 {
			expConditions = append(expConditions,
				FilterCondition{
					Column: "years_of_experience_raw", // Replace with your actual column name
					Type:   ">=",
					Value:  string(req.YearsOfExperienceFrom),
				},
			)
		}

		// Check for maximum experience (YearsOfExperienceTo)
		// Only apply if it's greater than 0 AND greater than 'From' to avoid conflict,
		// or if 'From' is 0, just apply the 'To' limit.
		if req.YearsOfExperienceTo > 0 {
			// This check ensures 'To' is only used as an upper limit
			// and isn't a simple 'YearsOfExperienceFrom' value that defaulted 'To' to 0.
			if req.YearsOfExperienceTo >= req.YearsOfExperienceFrom {
				expConditions = append(expConditions,
					FilterCondition{
						Column: "years_of_experience_raw", // Replace with your actual column name
						Type:   "<=",
						Value:  string(req.YearsOfExperienceTo),
					},
				)
			}
		}

		// Experience filters must ALL apply (AND)
		if len(expConditions) > 0 {
			conditions = append(conditions, FilterGroup{
				Op:         "and",
				Conditions: expConditions,
			})
		}
	}
	// 6. Health filters (ensure profiles are complete)
	healthFilters := []interface{}{
		// summary must not be empty
		FilterCondition{Column: "summary", Type: "!=", Value: ""},
		// skills must not be empty
		FilterCondition{Column: "skills", Type: "!=", Value: ""},
		// projects must not be empty
		// experience fields should have non-zero data
		FilterCondition{Column: "current_employers.title", Type: "!=", Value: ""},
		FilterCondition{Column: "past_employers.title", Type: "!=", Value: ""},
	}

	conditions = append(conditions, FilterGroup{
		Op:         "and",
		Conditions: healthFilters,
	})

	// 7. Build final Crustdata filter group
	root := FilterGroup{
		Op:         "and",
		Conditions: conditions,
	}

	// 8. Return Crustdata PeopleSearchRequest
	return &crustdata.PeopleSearchRequest{
		Filters: root,
		Limit:   int(req.Limit),
	}
}

// Converts a gRPC request to SQLC query parameters.
func (a *PeopleAdapter) RawProfileListSqlFromGrpc(req *talv1.RawProfileListRequest) *db.RawProfileListParams {
	if req == nil {
		return &db.RawProfileListParams{}
	}

	if req.Limit == 0 {
		req.Limit = 10
	}
	// queryEmbedding, _ := a.embed(context.Background(), req.Query, &genaiclient.EmbedOptions{Dimensions: 1536})
	return &db.RawProfileListParams{
		Query: req.Query,
		// Embedding:  pgvector.NewVector(queryEmbedding),
		Industries: req.Industries,
		Locations:  req.Locations,
		Skills:     req.Skills,
		Companies:  req.Companies,
		Projects:   req.Projects,
		Limit:      req.Limit,
	}
}
func (a *PeopleAdapter) RawProfileListGrpcFromSql(r *db.RawProfileListRow) *talv1.RawProfileListResponse {
	if r == nil {
		return nil
	}
	return &talv1.RawProfileListResponse{
		PersonId:       r.PersonID,
		Name:           r.Name,
		Headline:       r.Headline,
		Location:       r.Region,
		CurrentTitle:   r.CurrentTitle,
		CurrentCompany: r.CurrentCompany,
		Industry:       r.Industry,
		Summary:        r.Summary,
		Skills:         r.Skills,
		SemanticScore:  float32(r.SemanticScore),
		TextRank:       float32(r.TextRank),
		HybridScore:    float32(r.HybridScore),
	}
}

// Converts SQLC rows to gRPC response.
//
//	func (a *PeopleAdapter) RawProfileListGrpcFromSql(rows *[]*db.RawProfileListRow, crustDataReq ...*[]*db.RawProfileListRow) *talv1.RawProfileListResponse {
//		resp := &talv1.RawProfileListResponse{
//			Records: make([]*talv1.RawProfileListRow, len(*rows)),
//		}
//
//		for index, r := range *rows {
//			row := a.rawProfileRowToGrpc(r)
//			resp.Records[index] = row
//		}
//		if crustDataReq != nil {
//			for _, row := range crustDataReq {
//				for _, r := range *row {
//					row := a.rawProfileRowToGrpc(r)
//					resp.Records = append(resp.Records, row)
//				}
//			}
//
//		}
//
//		return resp
//	}

func (a *PeopleAdapter) RawProfileFindGrpcFromSql(req *db.RawProfileFindRow) *talv1.RawProfileFindResponse {
	if req == nil {
		return nil
	}

	// Helper to decode JSON arrays or objects into structs
	decodeJSON := func(data []byte, v interface{}) {
		if len(data) == 0 {
			return
		}
		_ = json.Unmarshal(data, v)
	}

	// Helper to decode [][]byte JSON arrays (for arrays of objects)
	decodeJSONArray := func(arr [][]byte, v interface{}) {
		if len(arr) == 0 {
			return
		}

		sliceValue := reflect.ValueOf(v).Elem() // must be pointer to slice
		elemType := sliceValue.Type().Elem()

		for _, b := range arr {
			elemPtr := reflect.New(elemType)
			if err := json.Unmarshal(b, elemPtr.Interface()); err != nil {
				continue
			}
			sliceValue.Set(reflect.Append(sliceValue, elemPtr.Elem()))
		}
	}

	// Prepare destination slices
	var (
		currentEmployers    []*talv1.Employment
		pastEmployers       []*talv1.Employment
		educationBackground []*talv1.Education
		honors              []*talv1.Honor
		certifications      []*talv1.Certification
		openToCards         []bool
	)

	// Decode JSONB[] fields
	decodeJSONArray(req.CurrentEmployers, &currentEmployers)
	decodeJSONArray(req.PastEmployers, &pastEmployers)
	decodeJSONArray(req.EducationBackground, &educationBackground)

	// Decode single JSONB fields
	decodeJSON(req.Honors, &honors)
	decodeJSON(req.Certifications, &certifications)

	// Handle open_to_cards JSONB[] → []bool
	decodeJSONArray(req.OpenToCards, &openToCards)

	// Handle region_address_components TEXT or JSON TEXT[]
	var regionComponents []string
	if req.RegionAddressComponents.Valid && req.RegionAddressComponents.String != "" {
		_ = json.Unmarshal([]byte(req.RegionAddressComponents.String), &regionComponents)
	}

	return &talv1.RawProfileFindResponse{
		PersonId:                req.PersonID,
		Name:                    req.Name.String,
		FirstName:               req.FirstName.String,
		LastName:                req.LastName.String,
		Region:                  req.Region.String,
		RegionAddressComponents: regionComponents,
		Headline:                req.Headline.String,
		Summary:                 req.Summary.String,
		Skills:                  req.Skills,
		Languages:               req.Languages,
		ProfileLanguage:         req.ProfileLanguage.String,
		TwitterHandle:           req.TwitterHandle.String,
		OpenToCards:             openToCards,
		NumOfConnections:        req.NumOfConnections.Int32,
		RecentlyChangedJobs:     req.RecentlyChangedJobs.Bool,
		YearsOfExperience:       req.YearsOfExperience.String,
		YearsOfExperienceRaw:    req.YearsOfExperienceRaw.Int32,
		CurrentEmployers:        currentEmployers,
		PastEmployers:           pastEmployers,
		EducationBackground:     educationBackground,
		Honors:                  honors,
		Certifications:          certifications,
	}
}

// func (a *PeopleAdapter) 	RawProfileGrpcFromCrustdata(row *crustdata.PersonDBProfile) *talv1.RawProfileListResponse{
//
//		return &talv1.rawprofilefindresponse{
//			personid:                req.personid,
//			name:                    req.name.string,
//			firstname:               req.firstname.string,
//			lastname:                req.lastname.string,
//			region:                  req.region.string,
//			regionaddresscomponents: regioncomponents,
//			headline:                req.headline.string,
//			summary:                 req.summary.string,
//			skills:                  req.skills,
//			languages:               req.languages,
//			profilelanguage:         req.profilelanguage.string,
//			twitterhandle:           req.twitterhandle.string,
//			opentocards:             opentocards,
//			numofconnections:        req.numofconnections.int32,
//			recentlychangedjobs:     req.recentlychangedjobs.bool,
//			yearsofexperience:       req.yearsofexperience.string,
//			yearsofexperienceraw:    req.yearsofexperienceraw.int32,
//			currentemployers:        currentemployers,
//			pastemployers:           pastemployers,
//			educationbackground:     educationbackground,
//			honors:                  honors,
//			certifications:          certifications,
//		}
//	}
func (a *PeopleAdapter) RawProfileListGrpcFromCrustdata(row *crustdata.PersonDBProfile) *talv1.RawProfileListResponse {
	var currentTitle, currentCompany string
	if len(row.CurrentEmployers) > 0 {
		currentTitle = row.CurrentEmployers[0].Title
		currentCompany = row.CurrentEmployers[0].Name
	} else if len(row.PastEmployers) > 0 {
		currentTitle = row.PastEmployers[0].Title
		currentCompany = row.PastEmployers[0].Name
	}
	return &talv1.RawProfileListResponse{
		PersonId:       row.PersonID,
		Name:           row.Name,
		Headline:       row.Headline,
		Location:       row.Region,
		CurrentTitle:   currentTitle,
		CurrentCompany: currentCompany,
		Industry:       "", // Crustdata doesn’t provide a direct industry field; leave empty or derive later
		Summary:        row.Summary,
		Skills:         row.Skills,
		SemanticScore:  0, // To be calculated if needed
		TextRank:       0, // To be calculated if needed
		HybridScore:    0, // To be calculated if needed
	}

}
