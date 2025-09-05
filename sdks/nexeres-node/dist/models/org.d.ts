export interface OrgJSON {
    id: string;
    name: string;
    slug: string;
    description?: string | undefined;
    avatarURL?: string | undefined;
    settings?: Record<string, any> | undefined;
    createdAt: string | Date;
    updatedAt: string | Date;
    deletedAt?: string | Date | undefined;
}
export declare class Org {
    id: string;
    name: string;
    slug: string;
    description?: string | undefined;
    avatarURL?: string | undefined;
    settings?: Record<string, any> | undefined;
    createdAt: Date;
    updatedAt: Date;
    deletedAt?: Date | undefined;
    constructor(data: OrgJSON);
}
//# sourceMappingURL=org.d.ts.map