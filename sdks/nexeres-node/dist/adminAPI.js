export class AdminAPI {
    client;
    constructor(client) {
        this.client = client;
    }
    login(params) {
        return this.client.post("/api/admin/login", params);
    }
    verifyLogin(params) {
        // The empty string is a placeholder for the admin token, which is not needed for this endpoint
        return this.client.adminPost("/api/admin/login/verify", "", params);
    }
    getConfig(adminToken) {
        return this.client.adminGet("/api/admin/config", adminToken);
    }
}
//# sourceMappingURL=adminAPI.js.map