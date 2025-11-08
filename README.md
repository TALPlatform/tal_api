# devkit-api

**Boost your Go development with devkit-api: A robust and scalable foundation for building backend applications with Clean Architecture.**

devkit-api is a Go project that provides a solid foundation for building backend applications. It leverages Clean Architecture principles to promote maintainability, testability, and scalability. With built-in features like Supabase integration, authentication, email sending, and more, devkit-api allows you to focus on building your core application logic instead of writing repetitive boilerplate code.

## Key Features

* **Supabase Integration:** Seamlessly integrates with Supabase, a powerful open-source Firebase alternative, for database and authentication functionalities.
* **Authentication and Authorization:** Provides built-in user authentication and authorization mechanisms, including role-based access control.
* **Email Sending:**  Integrated email sending functionality using Resend, making it easy to send transactional emails and notifications.
* **Internationalization:** Supports multiple languages, allowing you to build applications that cater to a global audience.
* **Storage Management:**  Provides tools for managing files and icons in Supabase storage.
* **Dynamic Table Imports:** Enables dynamic importing of database tables, simplifying database interactions.
* **Role-Based Access Control (RBAC):**  Implements a robust RBAC system for fine-grained control over user permissions.
* **Dynamic Navigation Bars:**  Supports dynamic generation of navigation bars based on user roles and permissions.

## Dependencies

devkit-api relies on several external libraries and services to provide its functionality. These dependencies are managed using Go modules.

**Key Dependencies:**

* **Supabase:**  A powerful open-source Firebase alternative that provides database, authentication, and storage services.
* **SQLc:**  A tool that generates type-safe Go code from SQL queries, improving code maintainability and reducing errors.
* **Buf:**  A modern Protobuf toolkit that simplifies API design and management.
* **Connect-Go:** A high-performance gRPC framework for Go.
* **Resend:**  An email API that makes it easy to send transactional emails and notifications.
* **Paseto:**  A secure and stateless token format used for authentication.
* **Redis:**  An in-memory data store used for caching and session management.
* **Various Go Libraries:**  Other Go libraries are used for tasks such as string manipulation, data validation, and error handling.

**Go Modules:**

The `go.mod` file lists all the project's dependencies and their versions. You can use `go get` to install or update these dependencies.

## Architecture

devkit-api follows Clean Architecture principles, promoting a separation of concerns and ensuring that your application is maintainable, testable, and scalable. The project is structured into the following layers:

* **API Layer:**  Handles incoming API requests and responses.
* **Domain Layer:**  Contains the core business logic of the application, organized into domains (e.g., Accounts, Products).
    * **Use Case Layer:**  Defines and implements use cases, which represent specific actions or operations within a domain.
    * **Repository Layer:**  Provides an abstraction for data access, hiding the specifics of the underlying database.
    * **Adapter Layer:**  Handles the transformation of data between the API layer and the repository layer.
* **SQL Layer:**  Interacts with the database using SQL queries.

<img src="https://miro.medium.com/v2/resize:fit:554/format:webp/1*EwIhtYgXBLMIfw7LjzW_iA.png" alt="Clean Architecture Diagram" width="500">

## Database Schema

**Public Schema:**

```sql
CREATE TABLE tags(
    tag varchar NOT NULL UNIQUE
);
 
CREATE TABLE input_types(
    input_type_id serial PRIMARY KEY,
    input_type_name varchar(20) NOT NULL UNIQUE
);

CREATE TABLE setting(
    setting_id serial PRIMARY KEY,
    input_type_id int NOT NULL,
    FOREIGN KEY (input_type_id) REFERENCES input_types(input_type_id),
    setting_key varchar(100) NOT NULL UNIQUE,
    setting_value text NOT NULL,
    updated_at timestamp

);
CREATE TABLE translations(
    translation_key varchar(200) NOT NULL UNIQUE,
    arabic_value varchar(200) NOT NULL ,
    english_value varchar(200) NOT NULL ,
    primary key(translation_key)
); 

CREATE TABLE icons(
    icon_id serial PRIMARY KEY,
    icon_name varchar(200) NOT NULL UNIQUE,
    icon_content text  NOT NULL
); 
```

**Accounts Schema:**

```sql
CREATE SCHEMA accounts_schema;
create table accounts_schema.user_types(
	user_type_id serial primary key,
	user_type_name varchar(200) not null unique
);

create table accounts_schema.users(
	user_id serial primary key,
	user_name varchar(200) not null,
	user_security_level int NOT NULL,
	user_type_id int NOT NULL,
	FOREIGN KEY (user_type_id) REFERENCES accounts_schema.user_types  (user_type_id),
	user_phone varchar(200) unique,
	user_email varchar(200) not null unique,
	user_password varchar(200),
	created_at timestamp not null default now(),
	updated_at timestamp,
	deleted_at timestamp 
);
```

## gRPC Service Definition

```protobuf
syntax = "proto3";

package tal.v1;
service TalService {
  // INJECT METHODS

//////////////////////////////////////////////////////////////////////////////////////////////
// public 
//////////////////////////////////////////////////////////////////////////////////////////////
    // settings
    rpc SettingsUpdate(SettingsUpdateRequest) returns (SettingsUpdateResponse) {}
    rpc SettingsFindForUpdate(SettingsFindForUpdateRequest) returns (SettingsFindForUpdateResponse)  {
        option idempotency_level = NO_SIDE_EFFECTS;
    }

    // icons
    rpc IconsInputList(google.protobuf.Empty) returns (IconsListResponse)  {
        option idempotency_level = NO_SIDE_EFFECTS;
    }
    rpc IconsCreateUpdateBulk(IconsCreateUpdateBulkRequest) returns (IconsListResponse)  {
        option idempotency_level = NO_SIDE_EFFECTS;
    }    // storage
    rpc FilesList(FilesListRequest) returns (FilesListResponse) {
        option idempotency_level = NO_SIDE_EFFECTS;
    }
    rpc BucketsList(google.protobuf.Empty) returns (BucketsListResponse)  {
        option idempotency_level = NO_SIDE_EFFECTS;
    }
    rpc UploadFile(UploadFileRequest) returns (UploadFileResponse) {}
    rpc UploadFiles(UploadFilesRequest) returns (UploadFileResponse) {}
    rpc FilesDelete(FilesDeleteRequest) returns (FilesDeleteResponse) {}
    rpc ImportTable(ImportTableRequest) returns (ImportTableResponse) {}
    // emails
    rpc SendEmail(SendEmailRequest) returns (SendEmailResponse) {}
    // translations
    rpc TranslationsCreateUpdateBulk(TranslationsCreateUpdateBulkRequest) returns (TranslationsListResponse) {}
    rpc TranslationsDelete(TranslationsDeleteRequest) returns (TranslationsListResponse) {}
    rpc TranslationsList(google.protobuf.Empty) returns (TranslationsListResponse)  {
        option idempotency_level = NO_SIDE_EFFECTS;
    }
//////////////////////////////////////////////////////////////////////////////////////////////
// accounts
//////////////////////////////////////////////////////////////////////////////////////////////
    // roles
//*******************************************************************************************//
    rpc RolesList(google.protobuf.Empty) returns (RolesListResponse) {
        option idempotency_level = NO_SIDE_EFFECTS;
    }
    rpc RoleCreateUpdate(RoleCreateUpdateRequest) returns (RoleCreateUpdateResponse) {}
    rpc RolesDeleteRestore(DeleteRestoreRequest) returns (google.protobuf.Empty) {}
//*******************************************************************************************//
    // users
//*******************************************************************************************//
    rpc UsersList(google.protobuf.Empty) returns (UsersListResponse)  {
        option idempotency_level = NO_SIDE_EFFECTS;
    }
    rpc UserCreateUpdate(UserCreateUpdateRequest) returns (UserCreateUpdateResponse) {}
    rpc UserDelete(google.protobuf.Empty) returns (AccountsSchemaUser) {}
    rpc UsersDeleteRestore(DeleteRestoreRequest) returns (google.protobuf.Empty) {}
//*******************************************************************************************//
    // auth
//*******************************************************************************************//
    rpc UserLoginProviderCallback(UserLoginProviderCallbackRequest) returns (UserLoginResponse) {
        option idempotency_level = NO_SIDE_EFFECTS;
    }
    rpc UserResetPassword(UserResetPasswordRequest) returns (UserLoginResponse) {}
    rpc UserResetPasswordEmail(UserResetPasswordEmailRequest) returns (UserResetPasswordEmailResponse) {}
    rpc UserLoginProvider(UserLoginProviderRequest) returns (UserLoginProviderResponse) {}
    rpc UserInvite(UserInviteRequest) returns (UserInviteResponse) {}
    rpc UserAuthorize(google.protobuf.Empty) returns (UserLoginResponse) {}
    rpc UserLogin(UserLoginRequest) returns (UserLoginResponse) {}
}
```

## Endpoint Creation Lifecycle

Here's an example of how an endpoint is created in devkit-api, using the "List Users" functionality:

1. **API Request:**  A client sends a request to the `/users` endpoint to retrieve a list of users.
2. **API Layer:** The API layer receives the request and passes it to the `ListUsers` use case in the Accounts domain.
3. **Use Case Layer:** The `ListUsers` use case retrieves the necessary data from the `User` repository.
4. **Repository Layer:** The `User` repository interacts with the database to fetch the list of users using the following SQL query:

    ```sql
    -- name: UsersList :many
    SELECT  
        user_id,
        user_name,
        user_security_level,
        user_type_id,
        user_phone,
        user_email,
        user_password,
        created_at,
        updated_at,
        deleted_at
    FROM accounts_schema.users;
    ```

5. **Adapter Layer:** The adapter layer transforms the database response into a format suitable for the API response, defined by the following Protobuf message:

    ```protobuf
    message UsersListResponse{
      repeated AccountsSchemaUser records= 1;
      repeated AccountsSchemaUser deleted_records = 2;
      ListDataOptions options = 3;
    }
    ```

6. **API Response:** The API layer sends the formatted response back to the client.

**Code Snippets:**

* **app/accounts/adapter/users_adapter.go:**

    ```go
    func (a *AccountsAdapter) UserEntityGrpcFromSql(resp *db.AccountsSchemaUser) *talv1.AccountsSchemaUser {
        // ... (code to transform database response to gRPC response)
    }

    func (a *AccountsAdapter) UsersListGrpcFromSql(resp []db.AccountsSchemaUser) *talv1.UsersListResponse {
        // ... (code to transform a list of database responses to gRPC response)
    }
    ```

* **app/accounts/usecase/users_usecase.go:**

    ```go
    func (u *AccountsUsecase) UsersList(ctx context.Context) (*apiv1.UsersListResponse, error) {
        // ... (code to handle the "List Users" use case)
    }
    ```

* **api/api.go:**

    ```go
    func (api *Api) UsersList(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[apiv1.UsersListResponse], error) {
        // ... (code to handle the API request and call the use case)
    }
    ```

## Getting Started

1. **Prerequisites:**
    * Go 1.16 or higher
    * Supabase CLI
    * SQLc CLI
    * Buf CLI
    * Resend account
    * Google Cloud Project (for Google authentication)

2. **Installation:**
    ```bash
    git clone [invalid URL removed]
    cd devkit-api
    ```

3. **Configuration**

devkit-api utilizes a robust and flexible configuration system that allows you to easily manage settings for different environments. Here's how it works:

**Environment Management:**

* **`state.env`:** This file specifies the current environment (e.g., `dev`, `prod`, `staging`). 
* **Environment-specific files:**  Based on the value in `state.env`, the application loads the corresponding environment-specific file (e.g., `dev.env`, `prod.env`, `staging.env`). This allows you to define different settings for each environment.
* **`shared.env`:**  This file can be used to store settings that are shared across all environments.

**Example:**

If `state.env` contains `STATE=dev`, the application will load the settings from `dev.env`.

**Configuration Files:**

* **`.env` files:** These files contain key-value pairs representing environment variables.
* **`supabase/config.toml`:**  This file stores Supabase connection details.

**Example `dev.env` file:**

```
DB_USER=postgres
DB_PROJECT_REF=supaprojectref
DB_DRIVER=postgres
DB_PORT=54322
SUPABASE_SERVICE_ROLE_KEY=supaservicekey
SUPABASE_API_KEY=supaanonkey
DB_HOST=localhost
DB_NAME=postgres
DB_PASSWORD=postgres
GRPC_PORT=9090
RESEND_API_KEY=resendapikey

GRPC_HOST=0.0.0.0
DB_SOURCE=${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable
ACCESS_TOKEN_DURATION=730h
TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012
GRPC_SERVER_ADDRESS=${GRPC_HOST}:${GRPC_PORT}
CLIENT_BASE_URL=[http://z4.com:5173/](http://z4.com:5173/)

REDIS_PORT=6379
REDIS_HOST=localhost
REDIS_DATABASE=2
REDIS_PASSWORD=
```

**Programmatic Access:**

The `config` package provides functions for loading and accessing configuration values. You can modify the `Config` struct in `config/config.go` to add or remove configuration keys as needed.

**Example Usage:**

```go
package main

import (
    "fmt"

    "[invalid URL removed]"
)

func main() {
    // Load the state configuration
    stateConfig, err := config.LoadState("conf") 
    if err != nil {
        panic(err)
    }

    // Load the environment configuration
    envConfig, err := config.LoadConfig("conf", stateConfig.State) 
    if err != nil {
        panic(err)
    }

    fmt.Println(envConfig.DBSource) // Access the database source string
}
```


## devkit-cli Integration

devkit-api is designed to work seamlessly with [devkit-cli](https://github.com/darwishdev/devkit-cli), a command-line interface that streamlines the creation and management of Go backend applications. devkit-cli provides a suite of commands that automate common development tasks, such as:

* **`devkit new api`:** Bootstraps a new Go backend API application by forking the devkit-api repository, cloning it locally, and performing the necessary configurations.
* **`devkit new domain`:** Generates a new domain within your Go backend application, creating the necessary directory structure and boilerplate code.
* **`devkit new feature`:** Generates a new feature within a specified domain, including adapter, repository, use case, and API components.
* **`devkit new endpoint`:** Generates a new endpoint within a specified feature and domain, handling code generation for all relevant layers.
* **`devkit seed`:**  Automates the process of seeding your database tables with data from an Excel file.
* **`devkit seed storage`:** Seeds Supabase storage with files and icons from specified paths.

By using devkit-cli, you can significantly accelerate your development workflow and ensure consistency and best practices throughout your project.

**Example Usage:**

To create a new API application named "myapp" based on devkit-api:

```bash
devkit new api myapp
```

This will create a new repository by forking this repo, clone it locally, and set up the project with all the necessary configurations. You can then use the other devkit-cli commands to add domains, features, and endpoints to your application.

## Dynamic Table Imports

devkit-api leverages the [SQL Seeder](https://github.com/darwishdev/sqlseeder) Go package to enable dynamic importing of database tables. This feature simplifies database interactions by allowing you to seed your database with data from JSON or Excel files, while automatically handling relationships between tables.

**Key Features of SQL Seeder:**

* **Seed from JSON or Excel:** Generate SQL INSERT statements from structured data in either JSON or Excel format.
* **Relationship Support:** Handles one-to-many and many-to-many relationships between tables.
* **Customizable Delimiters:** Configure the delimiters used in your data for flexible parsing.
* **Templating:** Uses Go templates to generate the SQL statements, allowing for customization.

**How it Works in devkit-api:**

The `devkit seed` command in devkit-cli utilizes SQL Seeder to parse your data file and generate the corresponding SQL INSERT statements. It automatically detects and handles relationships between tables based on the column naming conventions defined by SQL Seeder.

**Runtime Table Importing:**

In addition to the `devkit seed` command, devkit-api provides a gRPC endpoint, `ImportTable`, that allows you to seed tables dynamically at runtime. This endpoint accepts an `ImportTableRequest` message, which includes the table name, schema name, sheet name (if applicable), and the data to be imported.

**`ImportTable` Endpoint:**

```protobuf
rpc ImportTable(ImportTableRequest) returns (ImportTableResponse) {}

message ImportTableRequest{
  string table_name = 1 [
  (buf.validate.field).string.min_len = 2
  ];
  string schema_name = 2;
  string sheet_name = 3;
  bytes reader = 4;
}

message ImportTableResponse {
  string message = 1;
}
```

This endpoint provides a flexible way to import data into your database on demand, without relying on the CLI.

**Benefits:**

* **Simplified Database Seeding:** Easily populate your database with data from external sources.
* **Automated Relationship Handling:**  No need to manually manage foreign key constraints and relationships.
* **Increased Efficiency:**  Reduces manual effort and potential errors associated with writing SQL INSERT statements.
* **Runtime Flexibility:** Seed tables dynamically at runtime using the `ImportTable` endpoint.

By incorporating SQL Seeder and providing the `ImportTable` endpoint, devkit-api provides a powerful and convenient way to manage your database and streamline your development workflow.

## Role-Based Access Control (RBAC)

devkit-api implements a robust Role-Based Access Control (RBAC) system to provide fine-grained control over user permissions. This system allows you to define roles, assign permissions to those roles, and then assign roles to users. This ensures that users can only access the resources and functionalities they are authorized for.

**Database Schema:**

The RBAC system is built upon the following database schema:

```sql
CREATE SCHEMA accounts_schema;

create table accounts_schema.permissions(
	permission_id serial PRIMARY KEY,
	permission_function varchar(200) NOT NULL UNIQUE,
	permission_name varchar(200) NOT NULL,
	permission_description varchar(200),
	permission_group varchar(200) NOT NULL
);

create table accounts_schema.roles(
	role_id serial primary key,
	role_name varchar(200) not null unique,
	role_description varchar(200),
	created_at timestamp not null default now(),
	updated_at timestamp,
	deleted_at timestamp 
);

CREATE TABLE accounts_schema.role_permission(
	role_id int NOT NULL,
	FOREIGN KEY (role_id) REFERENCES accounts_schema.roles(role_id),
	permission_id int NOT NULL,
	FOREIGN KEY (permission_id) REFERENCES accounts_schema.permissions(permission_id),
	PRIMARY KEY (role_id, permission_id)
);

-- ... (other tables)

CREATE TABLE accounts_schema.user_roles(
	user_id int NOT NULL,
	FOREIGN KEY (user_id) REFERENCES accounts_schema.users(user_id),
	role_id int NOT NULL,
	FOREIGN KEY (role_id) REFERENCES accounts_schema.roles(role_id),
	PRIMARY KEY (user_id, role_id)
);

CREATE TABLE accounts_schema.navigation_bars(
    navigation_bar_id serial PRIMARY KEY,
    menu_key varchar(200) UNIQUE NOT NULL,
    label varchar(200) NOT NULL,
    label_ar varchar(200),
    icon varchar(200),
    "route" varchar(200) UNIQUE,
    parent_id int,
    FOREIGN KEY (parent_id) REFERENCES accounts_schema.navigation_bars(navigation_bar_id),
    permission_id int
);
```

**Token-Based Authentication:**

devkit-api utilizes Paseto tokens and Redis to manage user authentication and authorization. Upon successful login, a token is generated and  the user's permissions are stored in Redis along. This token is then used to authenticate subsequent requests and determine the user's access level.

**Permission Caching:**

To optimize performance, user permissions are cached in Redis during the login process. This eliminates the need to query the database for permissions on every request.

**Dynamic Navigation Bars:**

The RBAC system also extends to dynamic navigation bars. The `accounts_schema.navigation_bars` table stores information about navigation menu items, including their associated permissions. This allows you to generate navigation bars that are tailored to each user's role and permissions.

**API Layer Integration:**

The API layer integrates with the RBAC system to enforce access control on endpoints. The `getAccessableActionsForGroup` function retrieves the user's permissions from Redis and determines which actions are allowed for a given resource group. This information is then included in the API response, allowing the frontend to dynamically display only the permitted actions.

**Example:**

```go
func (api *Api) RolesList(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[apiv1.RolesListResponse], error) {
    // ... (retrieve roles)

    options, err := api.getAccessableActionsForGroup(ctx, req.Header(), "roles")
    if err != nil {
        return nil, err
    }
    response.Options = options

    return connect.NewResponse(response), nil
}
```

**Benefits:**

* **Enhanced Security:** Ensures that users can only access authorized resources and functionalities.
* **Improved User Experience:** Provides users with a personalized experience by displaying only relevant actions and navigation options.
* **Simplified Development:**  Reduces the complexity of managing user permissions and access control.

By implementing a comprehensive RBAC system, devkit-api provides a secure and flexible foundation for building applications with varying levels of user access.

## Makefile Commands

devkit-api includes a Makefile that provides convenient commands for common development tasks. Here are some of the available commands:

* **`make mign name=<migration_name>`:** Creates a new Supabase migration with the specified name.
* **`make testdb`:** Runs tests for the database layer.
* **`make testapi`:** Runs tests for the API layer with race condition detection.
* **`make rdb`:** Resets the Supabase database.
* **`make run`:** Runs the devkit-api application.
* **`make buf`:** Generates Go code from Protobuf definitions using Buf.
* **`make sqlc`:** Generates Go code from SQL queries using SQLc.
* **`make gen`:** Generates both Protobuf and SQLc code.
* **`make mock`:** Generates mock implementations for database interfaces using Mockgen.
* **`make test`:** Runs all tests with code coverage reporting.

These commands simplify development tasks such as running tests, generating code, and managing the database. You can use `make help` to see a list of all available commands and their descriptions.

## Deployment

devkit-api is designed for easy deployment using Docker and Docker Compose. The project includes a Dockerfile and a docker-compose.yml file to streamline the process.

**Dockerfile:**

The Dockerfile defines the build process for creating a Docker image of the devkit-api application. It uses a multi-stage build to optimize the image size:

* **Build Stage:**  A Go image is used to build the application binary.
* **Run Stage:** A lightweight Alpine image is used to run the application, copying the binary from the build stage.

**Example Dockerfile:**

```dockerfile
# Build stage
FROM  golang:alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY start.sh .
COPY wait-for.sh . 
EXPOSE 9091
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
```

**docker-compose.yml:**

The docker-compose.yml file defines the services that make up the application, including the devkit-api service and any dependencies (e.g., database). It also handles volume mounting, port mapping, and the execution of startup scripts.

**Example docker-compose.yml:**

```yaml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: devkit_api
    volumes:
      - ./config:/app/config
    image: devkit_api
    ports:
      - 9091:9091
    entrypoint:
      [
        "/app/wait-for.sh",
        "postgres:5432",
        "--",
        "/app/start.sh"
      ]
    command: [ "/app/main" ]
```

**Deployment Steps:**

1. **Build the Docker Image:**

    ```bash
    docker-compose build
    ```

2. **Run the Application:**

    ```bash
    docker-compose up -d
    ```

This will start the devkit-api application and any dependencies in detached mode.

**Benefits of Using Docker:**

* **Consistency:** Ensures that the application runs the same way in all environments.
* **Portability:** Makes it easy to deploy the application to different platforms and cloud providers.
* **Isolation:** Isolates the application and its dependencies from the host system.
* **Scalability:** Simplifies scaling the application by running multiple containers.

By leveraging Docker and Docker Compose, devkit-api simplifies the deployment process and ensures a consistent and reliable runtime environment.
