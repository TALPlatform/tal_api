CREATE OR REPLACE FUNCTION get_first_word (input_text text)
	RETURNS text
	AS $$
BEGIN
	-- extract the first part before the first underscore
	RETURN split_part(input_text, '_', 1);
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION snake_to_spaced (input_text text)
	RETURNS text
	AS $$
BEGIN
	-- Replace underscores with spaces
	RETURN REGEXP_REPLACE(input_text, '_', ' ', 'g');
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION iif (condition boolean, true_result ANYELEMENT, false_result ANYELEMENT)
	RETURNS ANYELEMENT
	LANGUAGE plpgsql
	AS $$
BEGIN
	IF condition THEN
		RETURN true_result;
	ELSE
		RETURN false_result;
	END IF;
END
$$;

CREATE OR REPLACE FUNCTION snake_to_camel (input_text text)
	RETURNS text
	AS $$
DECLARE
	result text := '';
	word text;
BEGIN
	-- Split the input string by underscores
	FOR word IN
	SELECT
		unnest(string_to_array(input_text, '_'))
		LOOP
			-- Capitalize the first letter of each word
			result := result || initcap(word);
		END LOOP;
	RETURN result;
	END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION is_null (in_value ANYELEMENT)
	RETURNS boolean
	LANGUAGE plpgsql
	AS $$
DECLARE
	value_type varchar(30);
BEGIN
	IF in_value IS NULL THEN
		RETURN TRUE;
	END IF;
	SELECT
		pg_typeof(in_value) INTO value_type;
	IF value_type = 'character varying' OR value_type = 'text' THEN
		IF in_value = '' OR in_value IS NULL THEN
			RETURN TRUE;
		ELSE
			RETURN FALSE;
		END IF;
	ELSIF value_type = 'integer'
			OR value_type = 'real' THEN
			IF in_value = 0 OR in_value IS NULL THEN
				RETURN TRUE;
			ELSE
				RETURN FALSE;
			END IF;
		ELSEIF value_type LIKE '%[]' THEN
		IF array_length(in_value, 1) IS NULL THEN
			RETURN TRUE;
		ELSE
			RETURN FALSE;
		END IF;
	ELSIF value_type = 'time with time zone'
			OR value_type = 'time without time zone'
			OR value_type = 'timestamp with time zone'
			OR value_type = 'timestamp without time zone'
			OR value_type = 'boolean'
			OR value_type = 'boolean' THEN
			IF in_value IS NULL THEN
				RETURN TRUE;
			ELSE
				RETURN FALSE;
			END IF;
	ELSE
		-- Handle other data types if needed
		RAISE EXCEPTION 'Unsupported data type: %', pg_typeof(in_value);
	END IF;
END
$$;

CREATE OR REPLACE FUNCTION is_null_replace (in_value ANYELEMENT, in_target_value ANYELEMENT)
	RETURNS ANYELEMENT
	LANGUAGE plpgsql
	AS $$
DECLARE
	value_type varchar(30);
BEGIN
	IF is_null(in_value) THEN
		RETURN in_target_value;
	ELSE
		RETURN in_value;
	END IF;
END
$$;

CREATE OR REPLACE FUNCTION setting_update (keys text[], vals text[])
	RETURNS void
	LANGUAGE plpgsql
	AS $$
BEGIN
	-- Create a temporary table to hold the new values
	CREATE TEMP TABLE temp_settings AS
	SELECT
		unnest($1) AS setting_key,
		unnest($2) AS setting_value;
	-- Update the main table based on the temporary table
	UPDATE
		setting AS s
	SET
		setting_value = t.setting_value
	FROM
		temp_settings AS t
	WHERE
		s.setting_key = t.setting_key;
	-- Drop the temporary table
	DROP TABLE temp_settings;
END
$$;

CREATE OR REPLACE FUNCTION translation_create_update_bulk (in_keys text[], in_english_values text[], in_arabic_values text[])
	RETURNS SETOF translation
	LANGUAGE plpgsql
	AS $$
BEGIN
	DROP TABLE IF EXISTS temp_prepared_input;
	CREATE temp TABLE temp_prepared_input AS
	SELECT
		unnest(in_keys) translation_key,
		unnest(in_english_values) english_value,
		unnest(in_arabic_values) arabic_value;
	DROP TABLE IF EXISTS temp_to_be_updated;
	CREATE temp TABLE temp_to_be_updated AS
	SELECT
		t.translation_key,
		i.english_value,
		i.arabic_value
	FROM
		translation t
		JOIN temp_prepared_input i ON t.translation_key = i.translation_key;
	DROP TABLE IF EXISTS temp_to_be_created;
	CREATE temp TABLE temp_to_be_created AS
	SELECT
		i.translation_key,
		i.english_value,
		i.arabic_value
	FROM
		temp_prepared_input i
	LEFT JOIN translation t ON t.translation_key = i.translation_key
WHERE
	t.translation_key IS NULL;
	INSERT INTO translation (
		translation_key,
		english_value,
		arabic_value)
	SELECT
		translation_key,
		english_value,
		arabic_value
	FROM
		temp_to_be_created;
	UPDATE
		translation
	SET
		english_value = i.english_value,
		arabic_value = i.arabic_value
	FROM
		temp_to_be_updated i
	WHERE
		translation.translation_key = i.translation_key;
	RETURN query
	SELECT
		t.translation_key,
		t.english_value,
		t.arabic_value
	FROM
		translation t
		JOIN temp_prepared_input i ON t.translation_key = i.translation_key;
END
$$;

CREATE OR REPLACE FUNCTION icon_create_update_bulk (in_icon_name text[], in_icon_content text[])
	RETURNS SETOF icon
	LANGUAGE plpgsql
	AS $$
BEGIN
	DROP TABLE IF EXISTS temp_prepared_input;
	CREATE temp TABLE temp_prepared_input AS
	SELECT
		unnest(in_icon_name) icon_name,
		unnest(in_icon_content) icon_content;
	DROP TABLE IF EXISTS temp_to_be_updated;
	CREATE temp TABLE temp_to_be_updated AS
	SELECT
		ic.icon_id,
		i.icon_name,
		i.icon_content
	FROM
		icon ic
		JOIN temp_prepared_input i ON ic.icon_name = i.icon_name;
	DROP TABLE IF EXISTS temp_to_be_created;
	CREATE temp TABLE temp_to_be_created AS
	SELECT
		i.icon_name,
		i.icon_content
	FROM
		temp_prepared_input i
	LEFT JOIN icon ic ON ic.icon_name = i.icon_name
WHERE
	ic.icon_name IS NULL;
	INSERT INTO icon (
		icon_name,
		icon_content)
	SELECT
		icon_name,
		icon_content
	FROM
		temp_to_be_created;
	UPDATE
		icon
	SET
		icon_content = i.icon_content
	FROM
		temp_to_be_updated i
	WHERE
		icon.icon_name = i.icon_name;
	RETURN query
	SELECT
		ic.icon_id,
		ic.icon_name,
		ic.icon_content
	FROM
		icon ic
		JOIN temp_prepared_input i ON i.icon_name = ic.icon_name;
END
$$;

CREATE OR REPLACE FUNCTION nullable_foreign (in_foreign_id int)
	RETURNS int
	AS $$
BEGIN
	RETURN iif(in_foreign_id = 0, NULL, in_foreign_id);
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION nullable_value (value ANYELEMENT)
	RETURNS ANYELEMENT
	AS $$
BEGIN
	RETURN iif(is_null(value), NULL, value);
END;
$$
LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION "public"."describe"("p_schema_name" character varying, "p_table_name" character varying) RETURNS TABLE("column_name" character varying, "data_type" character varying, "is_nullable" character varying, "column_default" "text")
    LANGUAGE "sql"
    AS $$
SELECT
    c.column_name::VARCHAR,
    c.data_type::VARCHAR,
    c.is_nullable::VARCHAR,
    c.column_default::TEXT
FROM
    information_schema.columns c
WHERE
    c.table_schema = p_schema_name
    AND c.table_name = p_table_name
ORDER BY
    c.ordinal_position;
$$;
