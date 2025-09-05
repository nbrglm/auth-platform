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

export class Org {
  id: string;
  name: string;
  slug: string;
  description?: string | undefined;
  avatarURL?: string | undefined;
  settings?: Record<string, any> | undefined;
  createdAt: Date;
  updatedAt: Date;
  deletedAt?: Date | undefined;

  constructor(data: OrgJSON) {
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