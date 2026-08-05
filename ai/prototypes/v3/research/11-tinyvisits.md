# Tinyvisits — Research Summary

## Overview
Tinyvisits is a minimalistic, privacy-centric website analytics tool. Cookie-free, consent-less, compatible with GDPR and ePrivacy. Aligned with CNIL exemption for analytics. Built from the ground up to fully anonymize personal and device data on first touch.

## Core Philosophy
- **Truly privacy-centric**: Anonymizes personal and device data on first touch
- **No cookies, no fingerprinting**: Aggregated, purely stateless, irreversibly anonymous pageview counts
- **Consent-less**: Compatible with GDPR and ePrivacy without requiring user consent
- **Opt-out available**: Visitors can opt out of anonymous traffic count collection

## Tracking Approach Comparison
| Tool Type | Tracking Approach | GDPR & ePrivacy Compatibility |
|---|---|---|
| tinyvisits | No cookies, no fingerprinting. Aggregated, stateless, irreversibly anonymous counts. Opt-out available. | Compatible with GDPR and ePrivacy. Aligns with CNIL exemption |
| Other "privacy-friendly" tools | No cookies, but fingerprint hashing (IP + browser data + URL + salt, daily rotation) | Not explicitly approved. EDPB/DSK state fingerprinting hashes require prior user consent |

## Data Anonymization Approach
| Data Type | Processing Approach |
|---|---|
| Single pageview event | Split into separate statistical hits by URL, referring domain, major browser, device, major OS |
| IP address | Only processed to lookup geo country name, then discarded |
| Timestamp | Truncated to day |
| User agent | Truncated to 3 independent labels: major browser version, major OS, device type |
| Website URL | Standard URL parameters (ad campaign IDs, UTM tags) automatically pruned |
| Referring URL | Truncated to domain only |

## Key Differentiator: Zero Fingerprinting
Unlike other cookieless tools that use hash fingerprinting (IP + browser data + URL + salt), tinyvisits:
- Does NOT use fingerprinting based on IP address and device metadata
- Automatically splits pageviews into anonymous statistical counts
- Processed data is impossible to reverse
- Avoids any personal profile and/or device identification
- Statistical counts are rounded to nearest 10

## Two Modes
1. **Minimalistic analytics mode** (default): Consent-less aggregate statistical measurement. Immediate anonymization in RAM. Only anonymous aggregate statistics stored.
2. **Full analytics mode**: Uses cookies and persistent identifiers for detailed visitor behavior analysis. Requires valid prior end-user consent via CMP.

## Features (Minimalistic Mode)
- Minimalistic statistical pageview counts
- Data export
- DPIA pre-assessment template docs
- LIA template docs
- Dashboard graphs and tables with anonymous statistical counts
- No filters to single out specific users or segments

## Privacy Safeguards
- Analytics limited to single website
- No combining/matching/merging data from unrelated services
- No tracking of visitors outside single website
- UTM tags, ad campaign IDs, user identifiers stripped from URLs
- Referring URLs stripped to just domains
- No external data import possible
- IP address processed temporarily only for geo country name
- No measurement of reach possible
- Statistical counts rounded to nearest 10

## Integration
- Simple, lightweight JS snippet in page HTML template
- Tracks pageview count events automatically
- Opt-out via `tinyvisits_optout=1` in local storage (no unique identifier, just preference)

## Design Principles
- **Minimalistic**: Only the essential pageview counts, nothing more
- **Privacy-first by architecture**: Anonymization at the point of collection, not after
- **Consent-less**: No cookie banners, no consent flow needed
- **Stateless**: No persistent identifiers, no cross-session tracking
- **Compliant by design**: Engineered to meet CNIL analytics exemption requirements

## Sources
- https://tinyvisits.com/
- https://tinyvisits.com/privacy-policy
- https://tinyvisits.com/index.html
- https://tinyvisits.com/terms
