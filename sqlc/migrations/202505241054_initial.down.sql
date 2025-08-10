-- NBRGLM Auth Platform (NAP) - Migration Down
-- This file reverses the NAP schema by dropping all indexes, tables, and data.
-- Execute with caution - this will permanently delete all NAP data.
-- Drop composite indexes first
DROP INDEX IF EXISTS idx_audit_logs_org_created;

DROP INDEX IF EXISTS idx_invitations_org_email;

DROP INDEX IF EXISTS idx_user_orgs_user_status;

DROP INDEX IF EXISTS idx_sessions_user_org;

-- Drop audit logs indexes
DROP INDEX IF EXISTS idx_audit_logs_resource;

DROP INDEX IF EXISTS idx_audit_logs_action;

DROP INDEX IF EXISTS idx_audit_logs_created_at;

DROP INDEX IF EXISTS idx_audit_logs_user_id;

DROP INDEX IF EXISTS idx_audit_logs_org_id;

-- Drop invitations indexes
DROP INDEX IF EXISTS idx_invitations_status;

DROP INDEX IF EXISTS idx_invitations_expires_at;

DROP INDEX IF EXISTS idx_invitations_token;

DROP INDEX IF EXISTS idx_invitations_email;

DROP INDEX IF EXISTS idx_invitations_org_id;

-- Drop OIDC refresh tokens indexes
DROP INDEX IF EXISTS idx_oidc_refresh_tokens_access_token_id;

DROP INDEX IF EXISTS idx_oidc_refresh_tokens_expires_at;

DROP INDEX IF EXISTS idx_oidc_refresh_tokens_token_hash;

-- Drop OIDC access tokens indexes
DROP INDEX IF EXISTS idx_oidc_access_tokens_client_id;

DROP INDEX IF EXISTS idx_oidc_access_tokens_user_id;

DROP INDEX IF EXISTS idx_oidc_access_tokens_expires_at;

DROP INDEX IF EXISTS idx_oidc_access_tokens_token_hash;

-- Drop OIDC auth codes indexes
DROP INDEX IF EXISTS idx_oidc_auth_codes_user_id;

DROP INDEX IF EXISTS idx_oidc_auth_codes_client_id;

DROP INDEX IF EXISTS idx_oidc_auth_codes_expires_at;

DROP INDEX IF EXISTS idx_oidc_auth_codes_code;

-- Drop OIDC clients indexes
DROP INDEX IF EXISTS idx_oidc_clients_client_id;

DROP INDEX IF EXISTS idx_oidc_clients_org_id;

-- Drop scopes indexes
DROP INDEX IF EXISTS idx_scopes_default;

DROP INDEX IF EXISTS idx_scopes_service;

DROP INDEX IF EXISTS idx_scopes_name;

-- Drop sessions indexes
DROP INDEX IF EXISTS idx_sessions_expires_at;

DROP INDEX IF EXISTS idx_sessions_token_hash;

DROP INDEX IF EXISTS idx_sessions_org_id;

DROP INDEX IF EXISTS idx_sessions_user_id;

-- Drop MFA factors indexes
DROP INDEX IF EXISTS idx_mfa_factors_verified;

DROP INDEX IF EXISTS idx_mfa_factors_type;

DROP INDEX IF EXISTS idx_mfa_factors_user_id;

-- Drop user OAuth identities indexes
DROP INDEX IF EXISTS idx_user_oauth_identities_email;

DROP INDEX IF EXISTS idx_user_oauth_identities_provider;

DROP INDEX IF EXISTS idx_user_oauth_identities_user_id;

-- Drop OAuth providers indexes
DROP INDEX IF EXISTS idx_oauth_providers_enabled;

DROP INDEX IF EXISTS idx_oauth_providers_org_id;

-- Drop user orgs indexes
DROP INDEX IF EXISTS idx_user_orgs_role;

DROP INDEX IF EXISTS idx_user_orgs_status;

DROP INDEX IF EXISTS idx_user_orgs_org_id;

DROP INDEX IF EXISTS idx_user_orgs_user_id;

-- Drop users indexes
DROP INDEX IF EXISTS idx_users_deleted_at;

DROP INDEX IF EXISTS idx_users_email;

-- Drop orgs indexes
DROP INDEX IF EXISTS idx_orgs_deleted_at;

DROP INDEX IF EXISTS idx_orgs_domain;

DROP INDEX IF EXISTS idx_orgs_slug;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS audit_logs;

DROP TABLE IF EXISTS invitations;

DROP TABLE IF EXISTS oidc_refresh_tokens;

DROP TABLE IF EXISTS oidc_access_tokens;

DROP TABLE IF EXISTS oidc_auth_codes;

DROP TABLE IF EXISTS oidc_clients;

DROP TABLE IF EXISTS scopes;

DROP TABLE IF EXISTS sessions;

DROP TABLE IF EXISTS mfa_factors;

DROP TABLE IF EXISTS user_oauth_identities;

DROP TABLE IF EXISTS oauth_providers;

DROP TABLE IF EXISTS user_orgs;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS org_domains;

DROP TABLE IF EXISTS orgs;