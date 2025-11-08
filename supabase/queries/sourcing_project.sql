
-- name: ProjectInputList :many
SELECT
	sourcing_project_id value,
	sourcing_project_name label
FROM
	sourcing_schema.sourcing_project where tenant_id = IS_NULL_REPLACE(@tenant_id , tenant_id);
