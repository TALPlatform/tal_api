
LANG=en_US.UTF-8
SHELL=/bin/bash
.SHELLFLAGS=--norc --noprofile -e -u -o pipefail -c
# Include the main .env file
include config/state.env
# Construct the variable name based on STATE
CURRENT_STATE_FILE = config/$(STATE).env
# Include the appropriate .env file (e.g., dev.env or prod.env)
include $(CURRENT_STATE_FILE)

# Include the additional .env file
include config/shared.env


raw_profile_load:
	PGPASSWORD=postgres psql \
  -h 127.0.0.1 \
  -p 54322 \
  -U postgres \
  -d postgres \
  < raw_profile_data.sql

raw_profile_bu:
	PGPASSWORD=postgres pg_dump \
  -h 127.0.0.1 \
  -p 54322 \
  -U postgres \
  -d postgres \
  -t people_schema.raw_profile \
  --data-only \
  --column-inserts \
  > raw_profile_data.sql

mign : 
	supabase migration new $(name)
testdb:
	go test ./db/... -v

testdate:
	go test ./pkg/dateutils/... -v

testapi:
	go test ./api/... -v --race


seed_storage:
	devkit seed storage -f seeds/assets -i seeds/icons 

deploy:
	docker compose build && make dtag v=$(v) && make dpush v=$(v)

seed_accounts:
	devkit seed accounts_schema --file-path seeds/schemas/accounts.xlsx -e


# seed_raw_profiles:
# 	devkit seed table   --file-path ./seeds/profiles/raw_profiles_01_with_embedding.json   --function people_schema.raw_profiles_bulk_create_update   -e
seed_sourcing:
	devkit seed sourcing_schema --file-path ./seeds/schemas/sourcing.xlsx -e
seed_public:
	devkit seed public --file-path seeds/schemas/public.xlsx -e


seed_properties:
	devkit seed properties_schema --file-path seeds/schemas/properties.xlsx -e
	
seed_tenants:
	devkit seed tenants_schema --file-path seeds/schemas/tenant.xlsx -e

seed_tenants_accounts:
	devkit seed accounts_schema --file-path seeds/schemas/tenant_accounts.xlsx -e



seed_super_user:
	devkit seed super-user -e admin@devkit.com -n "super admin user"

supabase_reset:
	supabase db reset 
			
rdbr:
	supabase db reset --linked

rdbrr:
	make rdbr seed_super_user seed_accounts  seed_storage seed_public seed_tenants seed_tenants_accounts 


SCHEMA_FILE    := weaviate_schema.json
CLASS_NAME     := $(shell jq -r '.class' $(SCHEMA_FILE))
.PHONY: init_weaviate_schema

init_weaviate_schema:
	  curl -s -X POST http://localhost:8080/v1/schema \
	    -H "Content-Type: application/json" \
	    --data-binary @$(SCHEMA_FILE) 

refresh_vector_db:
	curl -X DELETE http://localhost:8080/v1/schema/CommandPallete && make init_weaviate_schema

seed_all:
	make seed_super_user seed_accounts seed_storage seed_public seed_tenants seed_tenants_accounts seed_sourcing raw_profile_load
rdb:
	make supabase_reset refresh_vector_db seed_all
run:
	go run main.go
buf_push:
	cd proto && buf push

dtag:
	docker tag devkit_api exploremelon/tal_api:${v}

dpush:
	docker push  exploremelon/tal_api:${v}

buf:
	rm -rf proto_gen/tal/v1/*.pb.go && cd proto && buf lint && buf generate 
sqlc:
	rm -rf db/*.sql.go && sqlc generate	
gen:
	buf generate && sqlc generate

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/TALPlatform/tal_api/db Store
test:
	make mock && go test ./... -v --cover


