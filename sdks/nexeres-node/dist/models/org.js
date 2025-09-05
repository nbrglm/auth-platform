export class Org {
    id;
    name;
    slug;
    description;
    avatarURL;
    settings;
    createdAt;
    updatedAt;
    deletedAt;
    constructor(data) {
        this.id = data.id;
        this.name = data.name;
        this.slug = data.slug;
        this.description = data.description;
        this.avatarURL = data.avatarURL;
        this.settings = data.settings;
        this.createdAt = data.createdAt instanceof Date ? data.createdAt : new Date(data.createdAt);
        this.updatedAt = data.updatedAt instanceof Date ? data.updatedAt : new Date(data.updatedAt);
        this.deletedAt = data.deletedAt ? (data.deletedAt instanceof Date ? data.deletedAt : new Date(data.deletedAt)) : undefined;
    }
}
//# sourceMappingURL=org.js.map