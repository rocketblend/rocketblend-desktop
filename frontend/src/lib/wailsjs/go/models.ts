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
	
	export class Package {
	    id?: number[];
	    type?: number;
	    reference?: string;
	    name?: string;
	    author?: string;
	    tag?: string;
	    path?: string;
	    installationPath?: string;
	    // Go type: semver
	    version?: any;
	    addons?: string[];
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
	        this.reference = source["reference"];
	        this.name = source["name"];
	        this.author = source["author"];
	        this.tag = source["tag"];
	        this.path = source["path"];
	        this.installationPath = source["installationPath"];
	        this.version = this.convertValues(source["version"], null);
	        this.addons = source["addons"];
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
	    id?: number[];
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
	    id?: number[];
	    name?: string;
	    tags?: string[];
	    path?: string;
	    fileName?: string;
	    build?: string;
	    addons?: string[];
	
	    static createFrom(source: any = {}) {
	        return new UpdateProjectRequest(source);
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
	    }
	}

}

