# Task: Platform Features — Sharing, Access Control, Multi-Site, White Label

## Source
Plausible + Rybbit feature research

## Status: Not Started

## Objective
Add platform-level features for sharing, access control, multi-site management, and white-label support. These are lower priority and many are handled by the host application.

## Features

### Public Dashboards
- Shareable read-only dashboard URLs.
- No login required to view public dashboard.
- Configurable: host app can toggle public access per site.

### Private Link Sharing
- Password-protected dashboard links.
- Segment-limited: share a specific filter view (e.g., only "desktop traffic from Germany").

### Embedded Dashboards
- Embeddable dashboard widgets without statsstore branding.
- iframe-based or JS widget embed.

### Organizations / Multi-Site
- Multi-tenant site/team management.
- Site ID field on visitor records.
- Aggregate stats across multiple sites (consolidated view).
- Per-site configuration and access control.

### RBAC
- Role-based access control (admin, viewer, editor).
- Team management with role assignment.
- 2FA enforcement.
- **Note**: Likely handled by host application — statsstore is an embeddable library.

### SSO
- Single sign-on integration (SAML, OIDC).
- Enterprise feature.
- **Note**: Likely handled by host application.

### Saved Segments
- Save and share filter presets across the team.
- Stored per-user or per-organization.

### White Label
- Remove all statsstore branding from admin UI.
- Configurable theme/colors.
- Custom logo support.

### Script Proxying
- Serve tracking script as first-party from own domain.
- Avoids ad blocker detection.
- Reverse proxy configuration or built-in script serving.

### GA Import
- Import historical Google Analytics data.
- Requires GA API access and data mapping.

### Looker Studio Connector
- Data connector for Google Looker Studio.
- Expose data via the API layer (depends on `api-layer.md`).

### AI Assistant Traffic Tracking (Plausible)
- Detect and categorize traffic from AI tools (ChatGPT, Claude, Perplexity, etc.).
- Dedicated channel in traffic sources breakdown.
- User-agent pattern matching for known AI crawlers/assistants.

## Implementation Notes
- Many of these features (RBAC, SSO, auth) are typically handled by the host application embedding statsstore.
- Multi-site support requires a `site_id` field on visitor records — significant schema and query change.
- Public/private sharing requires an access control layer independent of the admin auth.
- White label is primarily a UI configuration task.

## Dependencies
- `api-layer.md` for Looker Studio connector and GA import.
- Schema changes for multi-site support.
- Access control layer for public/private sharing.

## Acceptance Criteria
- Public dashboards accessible via shareable URL without login.
- Private links password-protected with segment filtering.
- Multi-site support with per-site stats and consolidated view.
- White label configuration removes branding.
- Saved segments persist and are shareable.
