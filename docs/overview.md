# statsstore Overview

## Purpose
The `statsstore` module provides persistence and querying for website visitor analytics. It offers a unified interface for collecting, storing, and retrieving visitor sessions so downstream applications can power dashboards, reporting, and automation based on real traffic data.

## Core Concepts
- **Store**: Encapsulates database access for visitor records. It manages table creation, CRUD operations, counting, and query composition.
- **Visitor**: Data object representing a single visit. Fluent getters/setters keep records consistent while tracking dirty fields for efficient updates.
- **VisitorQueryOptions**: Filtering and pagination options that drive the SQL builder, enabling flexible list and count operations.

## Key Features
- Auto-migration to create the visitor table with standardized columns (path, IP, device, geo, timestamps, soft-delete).
- Goqu-based query builder for parameterized SQL across multiple drivers.
- Soft-delete support via `deleted_at` with option to include deleted entries.
- Distinct counting and paging helpers for analytics use cases.
- Request integration through `VisitorRegister`, which infers path, IP, and user agent directly from `*http.Request`.

## Usage Flow
1. **Initialize store** via `NewStore` with database handle, table name, and optional flags (auto-migrate, debug).
2. **Record visits** using `VisitorRegister` or by constructing a `Visitor` manually and calling `VisitorCreate`.
3. **Query analytics** with `VisitorList` and `VisitorCount`, passing `VisitorQueryOptions` to filter by ID, date range, country, and more.
4. **Maintain records** using `VisitorUpdate`, `VisitorDelete`, or `VisitorSoftDelete` depending on retention policies.

## Extensibility
- Columns and behaviors leverage the shared `sb` schema builder, making schema extensions straightforward.
- The store accepts any SQL driver supported by `database/sql`; driver name is inferred when omitted.
- Helper methods expose raw `*sql.DB` access for advanced use cases while keeping high-level operations consistent.

## Related Packages
- `github.com/dracory/dataobject` for visitor field management.
- `github.com/dracory/database` for executing parameterized SQL.
- `github.com/dromara/carbon/v2` for timestamp creation and parsing.
- `github.com/dracory/req` to extract request metadata during registration.
