export namespace metricservice {
	
	export class Aggregate {
	    domain: string;
	    name: string;
	    sum: number;
	    avg: number;
	    count: number;
	    min: number;
	    max: number;
	
	    static createFrom(source: any = {}) {
	        return new Aggregate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.domain = source["domain"];
	        this.name = source["name"];
	        this.sum = source["sum"];
	        this.avg = source["avg"];
	        this.count = source["count"];
	        this.min = source["min"];
	        this.max = source["max"];
	    }
	}
	export class FilterOptions {
	    domain: string;
	    name: string;
	    // Go type: time
	    startTime: any;
	    // Go type: time
	    endTime: any;
	
	    static createFrom(source: any = {}) {
	        return new FilterOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.domain = source["domain"];
	        this.name = source["name"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.endTime = this.convertValues(source["endTime"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Metric {
	    id: number[];
	    domain: string;
	    name: string;
	    value: number;
	    // Go type: time
	    recordedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Metric(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.domain = source["domain"];
	        this.name = source["name"];
	        this.value = source["value"];
	        this.recordedAt = this.convertValues(source["recordedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace operationservice {
	
	export class Operation {
	    id: number[];
	    completed: boolean;
	    error?: string;
	    result?: any;
	
	    static createFrom(source: any = {}) {
	        return new Operation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.completed = source["completed"];
	        this.error = source["error"];
	        this.result = source["result"];
	    }
	}

}

export namespace pack {
	
	export enum PackageType {
	    UNKNOWN = 0,
	    ADDON = 1,
	    BUILD = 2,
	}
	export enum PackageState {
	    AVAILABLE = 0,
	    DOWNLOADING = 1,
	    CANCELLED = 2,
	    INSTALLED = 3,
	    ERROR = 4,
	}
	export class Package {
	    id?: number[];
	    type: PackageType;
	    state: PackageState;
	    reference?: string;
	    name?: string;
	    author?: string;
	    tag?: string;
	    path?: string;
	    installationPath?: string;
	    operations?: string[];
	    addons?: string[];
	    platform?: number;
	    source?: rocketpack.Source;
	    // Go type: semver
	    version?: any;
	    verified?: boolean;
	    // Go type: time
	    updatedAt?: any;
	
	    static createFrom(source: any = {}) {
	        return new Package(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.state = source["state"];
	        this.reference = source["reference"];
	        this.name = source["name"];
	        this.author = source["author"];
	        this.tag = source["tag"];
	        this.path = source["path"];
	        this.installationPath = source["installationPath"];
	        this.operations = source["operations"];
	        this.addons = source["addons"];
	        this.platform = source["platform"];
	        this.source = this.convertValues(source["source"], rocketpack.Source);
	        this.version = this.convertValues(source["version"], null);
	        this.verified = source["verified"];
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace packageservice {
	
	export class GetPackageResponse {
	    package?: pack.Package;
	
	    static createFrom(source: any = {}) {
	        return new GetPackageResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.package = this.convertValues(source["package"], pack.Package);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListPackagesResponse {
	    packages?: pack.Package[];
	
	    static createFrom(source: any = {}) {
	        return new ListPackagesResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.packages = this.convertValues(source["packages"], pack.Package);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace project {
	
	export class Project {
	    id: number[];
	    name?: string;
	    tags?: string[];
	    path?: string;
	    fileName?: string;
	    build?: string;
	    addons?: string[];
	    splashPath?: string;
	    thumbnailPath?: string;
	    version?: string;
	    // Go type: time
	    updatedAt?: any;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.tags = source["tags"];
	        this.path = source["path"];
	        this.fileName = source["fileName"];
	        this.build = source["build"];
	        this.addons = source["addons"];
	        this.splashPath = source["splashPath"];
	        this.thumbnailPath = source["thumbnailPath"];
	        this.version = source["version"];
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace projectservice {
	
	export class CreateProjectRequest {
	    name?: string;
	    tags?: string[];
	    path?: string;
	    fileName?: string;
	    build?: string;
	    addons?: string[];
	
	    static createFrom(source: any = {}) {
	        return new CreateProjectRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.tags = source["tags"];
	        this.path = source["path"];
	        this.fileName = source["fileName"];
	        this.build = source["build"];
	        this.addons = source["addons"];
	    }
	}
	export class GetProjectResponse {
	    project?: project.Project;
	
	    static createFrom(source: any = {}) {
	        return new GetProjectResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project = this.convertValues(source["project"], project.Project);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListProjectsResponse {
	    projects?: project.Project[];
	
	    static createFrom(source: any = {}) {
	        return new ListProjectsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projects = this.convertValues(source["projects"], project.Project);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateProjectRequest {
	    id: number[];
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProjectRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}

}

export namespace rocketpack {
	
	export class Source {
	    resource?: string;
	    // Go type: downloader
	    uri?: any;
	
	    static createFrom(source: any = {}) {
	        return new Source(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.resource = source["resource"];
	        this.uri = this.convertValues(source["uri"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

