import { Org } from "./org.js";
export class Flow {
    id;
    type;
    userId;
    email;
    orgs;
    mfaRequired;
    mfaVerified;
    ssoProvider;
    ssoUserId;
    returnTo;
    createdAt;
    expiresAt;
    constructor(data) {
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
//# sourceMappingURL=flow.js.map