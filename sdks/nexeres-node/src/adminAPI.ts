import { NexeresAdminResponse, NexeresClient, NexeresResponse } from "./client.js";

export class AdminAPI {
  constructor(private client: NexeresClient) {
  }

  login(params: AdminLoginParams): NexeresResponse<AdminLoginResponse> {
    return this.client.post<AdminLoginResponse>("/api/admin/login", params);
  }

  verifyLogin(params: AdminLoginVerifyParams): NexeresAdminResponse<AdminLoginVerifyResponse> {
    // The empty string is a placeholder for the admin token, which is not needed for this endpoint
    return this.client.adminPost<AdminLoginVerifyResponse>("/api/admin/login/verify", "", params);
  }

  getConfig(adminToken: string): NexeresAdminResponse<NexeresAdminViewableConfig> {
    return this.client.adminGet<NexeresAdminViewableConfig>("/api/admin/config", adminToken);
  }
}

/** Parameters for the admin login */
export interface AdminLoginParams {
  /** The admin's email address */
  email: string;
}

/** Response for the admin login */
export interface AdminLoginResponse {
  /** Indicates if the login was successful */
  emailSent: string;

  /** The flow ID for the login verification step */
  flowId: string;
}

/** Parameters for the admin login verification */
export interface AdminLoginVerifyParams {
  /** The flow ID returned from the initial login request */
  flowId: string;
  /** The verification code sent to the admin's email */
  code: string;
}

/** Response for the admin login verification */
export interface AdminLoginVerifyResponse {
  /** Indicates if the verification was successful */
  success: boolean;

  /** The ephemeral token for the admin user, expires in 15 minutes of inactivity */
  token: string;
}

/** Viewable config for the admin */
export interface NexeresAdminViewableConfig {
  /** Indicates if backend is in debug mode */
  debug: boolean;

  /** The public configuration for the backend. Affects the URL generation in backend. */
  public: NexeresAdminPublicConfig;

  /** Indicates if the backend is in multitenancy mode */
  multitenancy: boolean;

  /**
   * The configuration for JWTs.
   */
  jwt: NexeresAdminJWTConfig;

  /**
   * Configuration for Notifications like Emails, SMS, etc.
   */
  notifications: NexeresAdminNotificationsConfig;

  /**
   * Configuration for Branding in the backend, for emails, sms, etc.
   */
  branding: NexeresAdminBrandingConfig;

  /**
   * Configuration for Security Settings, like API Keys, etc.
   */
  security: NexeresAdminSecurityConfig;
}

/** Public config for the admin */
export interface NexeresAdminPublicConfig {
  /** The scheme (http or https) used by the Nexeres instance */
  scheme: string;

  /** The domain used by the Nexeres UI instance, as configured in the backend. */
  domain: string;

  /** The subdomain used by the Nexeres UI instance, as configured in the backend. */
  subDomain: string;

  /** The debug base URL for the Nexeres instance, as configured in the backend. */
  debugBaseURL: string;
}

/** Configuration for JWTs */
export interface NexeresAdminJWTConfig {
  /** The expiration time for session tokens */
  sessionTokenExpiration: string;

  /** The expiration time for refresh tokens */
  refreshTokenExpiration: string;
}

/** Configuration for Notifications */
export interface NexeresAdminNotificationsConfig {
  /** Configuration for Email notifications */
  email: NexeresAdminEmailConfig;

  /** Configuration for SMS notifications */
  sms?: NexeresAdminSMSConfig;
}

/** Configuration for Email notifications */
export interface NexeresAdminEmailConfig {
  /** The email provider used for sending emails */
  provider: string;

  /**
   * The endpoints (full-URLs) used in emails for various actions like verification, password reset, etc.
   */
  endpoints: {
    verificationEmail: string;
    passwordReset: string;
  }
}

/** Configuration for SMS notifications */
export interface NexeresAdminSMSConfig {
  /** The SMS provider used for sending SMS */
  provider: string;
}

/** Configuration for Branding */
export interface NexeresAdminBrandingConfig {
  /** The name of the application, used in emails and other communications */
  appName: string;

  /** The full name of the Company */
  companyName: string;

  /** The short name of the Company */
  companyNameShort: string;

  /**
   * The SupportURL for the Company.
   *
   * Eg. "https://example.com/support" or "mailto:support@example.com"
   */
  supportURL: string;
}

/** Configuration for Security Settings */
export interface NexeresAdminSecurityConfig {
  /** Configuration for Audit Logs */
  auditLogs: {
    /** Indicates if audit logging is enabled */
    enable: boolean;
  }

  /** Configuration for API Keys */
  apiKeys: APIKeyAdminConfig[];

  /** Configuration for Rate Limiting */
  rateLimit: {
    /** The rate limit for API requests
     * 
     * Format: "R-U", where R is the number of requests and U is the time unit (e.g., "100-h" for 100 requests per hour)
     * 
     * Supported time units:
     * - s: Second
     * - m: Minute
     * - h: Hour
     * - d: Day
     */
    rate: string;
  }
}

/** Configuration for an API Key */
export interface APIKeyAdminConfig {
  /** The name of the API Key */
  name: string;

  /** The description of the API Key */
  description: string;
}