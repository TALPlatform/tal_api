## Core Changes:

- Renamed all package imports from `github.com/darwishdev/devkit-api` to `github.com/TALPlatform/tal_api`
- Updated proto package from `devkit.v1` to `tal.v1` and service from `DevkitService` to `TalService`
- Modified Makefile targets and build configurations

## New Feature Modules Added:

1. **Sourcing Module**:

   - New sourcing usecase with LLM integration
   - Sourcing schema seeding support
   - Added sourcing-related RPC endpoints

2. **People Module**:

   - New people usecase with Crustdata and LLM integration
   - Raw profile data management with backup/load functionality
   - Profile bulk operations support

3. **LLM Integration**:

   - Added new LLM client package for AI capabilities
   - Integration with both sourcing and people modules

4. **Crustdata Integration**:
   - New Crustdata API client for external data services
   - Used in people module for enhanced profile management

## Storage & Data Management:

- Added raw profile backup/load functionality via PostgreSQL
- Enhanced database seeding with sourcing data
- Updated database reset procedure to include profile data

## Removed Features:

- Removed entire property module (city/location management)
- Cleaned up unused property-related RPC endpoints and adapters

## Infrastructure:

- Updated Docker push references to new image names
- Modified API interceptors and middleware for new service name
- Enhanced Redis client integration for session management

## Testing:

- Updated all test references to new package structure
- Maintained existing test coverage with new naming

This rebranding and feature expansion positions the platform for talent acquisition and management use cases with enhanced AI capabilities and external data integrations.
