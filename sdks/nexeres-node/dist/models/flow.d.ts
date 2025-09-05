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
export declare class Flow {
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
    constructor(data: FlowJSON);
}
//# sourceMappingURL=flow.d.ts.map