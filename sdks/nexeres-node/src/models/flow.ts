import { Org, OrgJSON } from "./org.js";

export type FlowType = "login" | "change-password" | "sso";

export interface FlowJSON {
  id: string;
  type: FlowType;
  userId: string;
  email: string;
  orgs: Array<OrgJSON>;
  mfaRequired: boolean;
  mfaVerified: boolean;
  ssoProvider?: string | undefined;
  ssoUserId?: string | undefined;
  returnTo?: string | undefined;
  createdAt: string | Date;
  expiresAt: string | Date;
}

export class Flow {
  id: string;
  type: FlowType;
  userId: string;
  email: string;
  orgs: Array<Org>;
  mfaRequired: boolean;
  mfaVerified: boolean;
  ssoProvider?: string | undefined;
  ssoUserId?: string | undefined;
  returnTo?: string | undefined;
  createdAt: Date;
  expiresAt: Date;

  constructor(data: FlowJSON) {
    this.id = data.id;
    this.type = data.type;
    this.userId = data.userId;
    this.email = data.email;
    this.orgs = data.orgs.map(o => new Org(o));
    this.mfaRequired = data.mfaRequired;
    this.mfaVerified = data.mfaVerified;
    this.ssoProvider = data.ssoProvider;
    this.ssoUserId = data.ssoUserId;
    this.returnTo = data.returnTo;
    this.createdAt = data.createdAt instanceof Date ? data.createdAt : new Date(data.createdAt);
    this.expiresAt = data.expiresAt instanceof Date ? data.expiresAt : new Date(data.expiresAt);
  }
}