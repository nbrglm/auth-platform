-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    password_hash,
    first_name,
    last_name,
    avatar_url
  )
VALUES (
    $1,
    $2,
    sqlc.narg('password_hash'),
    sqlc.narg('first_name'),
    sqlc.narg('last_name'),
    sqlc.narg('avatar_url')
  )
RETURNING id,
  email,
  email_verified,
  first_name,
  last_name,
  avatar_url,
  created_at,
  updated_at;

-- name: UpdateUser :one
UPDATE users
SET first_name = coalesce(sqlc.narg('first_name'), first_name),
  last_name = coalesce(sqlc.narg('last_name'), last_name),
  avatar_url = coalesce(sqlc.narg('avatar_url'), avatar_url),
  updated_at = NOW()
WHERE email = sqlc.arg('email')
RETURNING id,
  email,
  email_verified,
  first_name,
  last_name,
  avatar_url,
  created_at,
  updated_at;

-- name: MarkUserEmailVerified :exec
UPDATE users
SET email_verified = TRUE,
  updated_at = NOW()
WHERE id = sqlc.arg('id');

-- name: SetUserBackupCodes :exec
UPDATE users
SET backup_codes = sqlc.arg('backup_codes'),
  updated_at = NOW()
WHERE email = sqlc.arg('email');

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = sqlc.arg('password_hash'),
  updated_at = NOW()
WHERE email = sqlc.arg('email');

-- name: GetLoginInfoForUser :one
SELECT *
FROM users
WHERE email = sqlc.arg('email');

-- name: GetInfoForSessionRefresh :one
SELECT u.first_name AS user_fname,
  u.last_name AS user_lname,
  u.email AS user_email,
  u.email_verified AS user_email_verified,
  u.avatar_url AS user_avatar_url,
  o.name AS org_name,
  o.slug AS org_slug,
  uo.role AS user_org_role
FROM users u
  INNER JOIN user_orgs uo ON u.id = uo.user_id
  INNER JOIN orgs o ON o.id = uo.org_id
WHERE u.id = sqlc.arg('user_id')
  AND o.id = sqlc.arg('org_id')
  AND uo.status != 'banned';

-- name: GetUserByEmail :one
SELECT id,
  email,
  email_verified,
  first_name,
  last_name,
  avatar_url,
  created_at,
  updated_at
FROM users
WHERE email = sqlc.arg('email');

-- name: GetUserByID :one
SELECT id,
  email,
  email_verified,
  first_name,
  last_name,
  avatar_url,
  created_at,
  updated_at
FROM users
WHERE id = sqlc.arg('id');

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW(),
  updated_at = NOW()
WHERE email = sqlc.arg('email');

-- name: CreateOrg :one
INSERT INTO orgs (
    id,
    slug,
    name,
    description,
    avatar_url,
    settings
  )
VALUES (
    $1,
    $2,
    $3,
    sqlc.narg('description'),
    sqlc.narg('avatar_url'),
    sqlc.narg('settings')
  ) ON CONFLICT (slug) DO NOTHING
RETURNING *;

-- name: UpdateOrg :one
UPDATE orgs
SET name = coalesce(sqlc.narg('name'), name),
  description = coalesce(sqlc.narg('description'), description),
  avatar_url = coalesce(sqlc.narg('avatar_url'), avatar_url),
  settings = coalesce(sqlc.narg('settings'), settings),
  updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateOrgWhereSlug :one
UPDATE orgs
SET name = coalesce(sqlc.narg('name'), name),
  description = coalesce(sqlc.narg('description'), description),
  avatar_url = coalesce(sqlc.narg('avatar_url'), avatar_url),
  settings = coalesce(sqlc.narg('settings'), settings),
  updated_at = NOW()
WHERE slug = sqlc.arg('slug')
RETURNING *;

-- name: GetOrgByID :one
SELECT *
FROM orgs
WHERE id = sqlc.arg('id');

-- name: GetOrgBySlug :one
SELECT *
FROM orgs
WHERE slug = sqlc.arg('slug');

-- name: SoftDeleteOrg :exec
UPDATE orgs
SET deleted_at = NOW(),
  updated_at = NOW()
WHERE id = sqlc.arg('id')
  AND slug != 'default';

-- name: GetOrgByDomain :one
SELECT o.*
FROM orgs o
  INNER JOIN org_domains od ON o.id = od.org_id
WHERE od.domain = sqlc.arg('domain');

-- name: GetOrgForDomainIfAutoJoin :one
SELECT o.*
FROM orgs o
  INNER JOIN org_domains od ON o.id = od.org_id
WHERE od.domain = sqlc.arg('domain')
  AND od.auto_join_enabled = TRUE
  AND od.verified = TRUE;

-- name: AddDomainToOrg :one
INSERT INTO org_domains (org_id, domain)
VALUES (
    sqlc.arg('org_id'),
    sqlc.arg('domain')
  ) ON CONFLICT (org_id, domain) DO NOTHING
RETURNING *;

-- name: RemoveDomainFromOrg :exec
DELETE FROM org_domains
WHERE org_id = sqlc.arg('org_id')
  AND domain = sqlc.arg('domain');

-- name: LinkUserToOrg :exec
INSERT INTO user_orgs (user_id, org_id, role)
VALUES (
    sqlc.arg('user_id'),
    sqlc.arg('org_id'),
    coalesce(sqlc.narg('role'), 'member')
  );

-- name: BanUserFromOrg :exec
UPDATE user_orgs
SET STATUS = 'banned'
WHERE user_id = sqlc.arg('user_id')
  AND org_id = sqlc.arg('org_id');

-- name: UnlinkUserFromOrg :exec
DELETE FROM user_orgs
WHERE user_id = sqlc.arg('user_id')
  AND org_id = sqlc.arg('org_id');

-- name: GetUserOrgsByEmail :many
SELECT sqlc.embed(o),
  sqlc.embed(uo)
FROM orgs o
  INNER JOIN user_orgs uo ON o.id = uo.org_id
  INNER JOIN users u ON u.id = uo.user_id
WHERE u.email = sqlc.narg('email');

-- name: GetUserOrgsByID :many
SELECT o.id,
  o.slug,
  o.name,
  o.description,
  o.avatar_url
FROM orgs o
  INNER JOIN user_orgs uo ON o.id = uo.org_id
  INNER JOIN users u ON u.id = uo.user_id
WHERE u.id = sqlc.narg('id');

-- name: CreateSession :one
INSERT INTO sessions (
    id,
    user_id,
    org_id,
    token_hash,
    refresh_token_hash,
    mfa_verified,
    ip_address,
    user_agent,
    mfa_verified_at,
    expires_at
  )
VALUES (
    sqlc.arg('id'),
    sqlc.arg('user_id'),
    sqlc.arg('org_id'),
    sqlc.arg('token_hash'),
    sqlc.arg('refresh_token_hash'),
    sqlc.arg('mfa_verified'),
    sqlc.arg('ip_address'),
    sqlc.arg('user_agent'),
    sqlc.narg('mfa_verified_at'),
    sqlc.arg('expires_at')
  )
RETURNING *;

-- name: GetSessionByID :one
SELECT *
FROM sessions
WHERE id = sqlc.arg('id');

-- name: GetSessionByToken :one
SELECT *
FROM sessions
WHERE token_hash = sqlc.arg('token_hash');

-- name: GetSessionByRefreshToken :one
SELECT *
FROM sessions
WHERE refresh_token_hash = sqlc.arg('refresh_token_hash');

-- name: RefreshSession :one
UPDATE sessions
SET token_hash = coalesce(sqlc.narg('token_hash'), token_hash),
  refresh_token_hash = coalesce(
    sqlc.narg('refresh_token_hash'),
    refresh_token_hash
  ),
  expires_at = coalesce(sqlc.narg('expires_at'), expires_at)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateSessionMFA :one
UPDATE sessions
SET mfa_verified = sqlc.arg('mfa_verified'),
  mfa_verified_at = sqlc.narg('mfa_verified_at'),
  updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateUserSessionAgentAndIP :one
UPDATE sessions
SET user_agent = coalesce(sqlc.narg('user_agent'), user_agent),
  ip_address = coalesce(sqlc.narg('ip_address'), ip_address),
  updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = sqlc.arg('id');

-- name: DeleteSessionByToken :exec
DELETE FROM sessions
WHERE token_hash = sqlc.arg('token_hash');

-- name: DeleteSessionByRefreshToken :exec
DELETE FROM sessions
WHERE refresh_token_hash = sqlc.arg('refresh_token_hash');

-- name: GetSessionsByUserID :many
SELECT *
FROM sessions
WHERE user_id = sqlc.arg('user_id');

-- name: GetSessionsByOrgID :many
SELECT *
FROM sessions
WHERE org_id = sqlc.arg('org_id');

-- name: GetSessionsByUserIDAndOrgID :many
SELECT *
FROM sessions
WHERE user_id = sqlc.arg('user_id')
  AND org_id = sqlc.arg('org_id');

-- name: GetInvitationByID :one
SELECT *
FROM invitations
WHERE id = sqlc.arg('id')
  AND expires_at > NOW();

-- name: GetInvitationByIDUnsafe :one
SELECT *
FROM invitations
WHERE id = sqlc.arg('id');

-- name: CreateInvitation :one
INSERT INTO invitations (
    id,
    org_id,
    email,
    role,
    invited_by,
    token,
    expires_at
  )
VALUES (
    sqlc.arg('id'),
    sqlc.arg('org_id'),
    sqlc.arg('email'),
    sqlc.arg('role'),
    sqlc.arg('invited_by'),
    sqlc.arg('token'),
    sqlc.arg('expires_at')
  )
RETURNING *;

-- name: GetInvitationByToken :one
SELECT *
FROM invitations
WHERE token = sqlc.arg('token')
  AND expires_at > NOW();

-- name: GetInvitationByTokenUnsafe :one
SELECT *
FROM invitations
WHERE token = sqlc.arg('token');

-- name: RevokeInvitation :exec
DELETE FROM invitations
WHERE id = sqlc.arg('id');

-- name: RevokeInvitationByToken :exec
DELETE FROM invitations
WHERE token = sqlc.arg('token');

-- name: RevokeInvitationByEmail :exec
DELETE FROM invitations
WHERE email = sqlc.arg('email')
  AND org_id = sqlc.arg('org_id');

-- name: NewVerificationToken :one
INSERT INTO verification_tokens (
    id,
    user_id,
    TYPE,
    token_hash,
    expires_at
  )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetVerificationTokenByHash :one
SELECT *
FROM verification_tokens
WHERE token_hash = sqlc.arg('token_hash')
  AND expires_at > NOW();