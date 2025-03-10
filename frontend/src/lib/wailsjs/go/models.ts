export namespace application {
	
	export class AddPackageOpts {
	    reference: string;
	
	    static createFrom(source: any = {}) {
	        return new AddPackageOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.reference = source["reference"];
	    }
	}
	export class AddProjectPackageOpts {
	    id: number[];
	    reference: string;
	
	    static createFrom(source: any = {}) {
	        return new AddProjectPackageOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.reference = source["reference"];
	    }
	}
	export class AggregateMetricsOpts {
	    domain?: string;
	    name?: string;
	    // Go type: time
	    startTime?: any;
	    // Go type: time
	    endTime?: any;
	
	    static createFrom(source: any = {}) {
	        return new AggregateMetricsOpts(source);
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
		    if (a.slice && a.map) {
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
	export class AggregateMetricsResult {
	    Aggregate?: types.Aggregate;
	
	    static createFrom(source: any = {}) {
	        return new AggregateMetricsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Aggregate = this.convertValues(source["Aggregate"], types.Aggregate);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class CancelOperationOpts {
	    id: number[];
	
	    static createFrom(source: any = {}) {
	        return new CancelOperationOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class CreateProjectOpts {
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateProjectOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	    }
	}
	export class CreateProjectResult {
	    operationID: number[];
	
	    static createFrom(source: any = {}) {
	        return new CreateProjectResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.operationID = source["operationID"];
	    }
	}
	export class DeleteProjectOpts {
	    id: number[];
	
	    static createFrom(source: any = {}) {
	        return new DeleteProjectOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class Details {
	    version: string;
	    platform: string;
	    installationPath: string;
	    packagePath: string;
	    applicationConfigPath: string;
	    rocketblendConfigPath: string;
	
	    static createFrom(source: any = {}) {
	        return new Details(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.platform = source["platform"];
	        this.installationPath = source["installationPath"];
	        this.packagePath = source["packagePath"];
	        this.applicationConfigPath = source["applicationConfigPath"];
	        this.rocketblendConfigPath = source["rocketblendConfigPath"];
	    }
	}
	export class Feature {
	    addon: boolean;
	    developer: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Feature(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.addon = source["addon"];
	        this.developer = source["developer"];
	    }
	}
	export class FileFilter {
	    displayName: string;
	    pattern: string;
	
	    static createFrom(source: any = {}) {
	        return new FileFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.displayName = source["displayName"];
	        this.pattern = source["pattern"];
	    }
	}
	export class GetOperationOpts {
	    id: number[];
	
	    static createFrom(source: any = {}) {
	        return new GetOperationOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class GetOperationResult {
	    operation?: types.Operation;
	
	    static createFrom(source: any = {}) {
	        return new GetOperationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.operation = this.convertValues(source["operation"], types.Operation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class GetPackageOpts {
	    id: number[];
	
	    static createFrom(source: any = {}) {
	        return new GetPackageOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class GetPackageResult {
	    package?: types.Package;
	
	    static createFrom(source: any = {}) {
	        return new GetPackageResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.package = this.convertValues(source["package"], types.Package);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class GetProjectResult {
	    project?: types.Project;
	
	    static createFrom(source: any = {}) {
	        return new GetProjectResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project = this.convertValues(source["project"], types.Project);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class InstallPackageOpts {
	    id: number[];
	
	    static createFrom(source: any = {}) {
	        return new InstallPackageOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class InstallPackageResult {
	    operationID: number[];
	
	    static createFrom(source: any = {}) {
	        return new InstallPackageResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.operationID = source["operationID"];
	    }
	}
	export class ListMetricsOpts {
	    domain?: string;
	    name?: string;
	    // Go type: time
	    startTime?: any;
	    // Go type: time
	    endTime?: any;
	
	    static createFrom(source: any = {}) {
	        return new ListMetricsOpts(source);
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
		    if (a.slice && a.map) {
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
	export class ListMetricsResult {
	    Metrics: types.Metric[];
	
	    static createFrom(source: any = {}) {
	        return new ListMetricsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Metrics = this.convertValues(source["Metrics"], types.Metric);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class ListOperationsResult {
	    operations: types.Operation[];
	
	    static createFrom(source: any = {}) {
	        return new ListOperationsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.operations = this.convertValues(source["operations"], types.Operation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class ListPackagesOpts {
	    query: string;
	    references?: string[];
	    type: enums.PackageType;
	    state: enums.PackageState;
	
	    static createFrom(source: any = {}) {
	        return new ListPackagesOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.query = source["query"];
	        this.references = source["references"];
	        this.type = source["type"];
	        this.state = source["state"];
	    }
	}
	export class ListPackagesResult {
	    packages?: types.Package[];
	
	    static createFrom(source: any = {}) {
	        return new ListPackagesResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.packages = this.convertValues(source["packages"], types.Package);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class ListProjectsOpts {
	    query: string;
	
	    static createFrom(source: any = {}) {
	        return new ListProjectsOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.query = source["query"];
	    }
	}
	export class ListProjectsResult {
	    projects: types.Project[];
	
	    static createFrom(source: any = {}) {
	        return new ListProjectsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projects = this.convertValues(source["projects"], types.Project);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class OpenDialogOptions {
	    defaultDirectory?: string;
	    defaultFilename?: string;
	    title?: string;
	    filters?: FileFilter[];
	
	    static createFrom(source: any = {}) {
	        return new OpenDialogOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.defaultDirectory = source["defaultDirectory"];
	        this.defaultFilename = source["defaultFilename"];
	        this.title = source["title"];
	        this.filters = this.convertValues(source["filters"], FileFilter);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class OpenExplorerOptions {
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new OpenExplorerOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}
	export class Preferences {
	    watchPath: string;
	    feature: Feature;
	
	    static createFrom(source: any = {}) {
	        return new Preferences(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.watchPath = source["watchPath"];
	        this.feature = this.convertValues(source["feature"], Feature);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class RemoveProjectPackageOpts {
	    id: number[];
	    reference: string;
	
	    static createFrom(source: any = {}) {
	        return new RemoveProjectPackageOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.reference = source["reference"];
	    }
	}
	export class RenderProjectOpts {
	    id: number[];
	
	    static createFrom(source: any = {}) {
	        return new RenderProjectOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class RunProjectOpts {
	    id: number[];
	
	    static createFrom(source: any = {}) {
	        return new RunProjectOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class SaveDialogOptions {
	    defaultDirectory?: string;
	    defaultFilename?: string;
	    title?: string;
	    filters?: FileFilter[];
	
	    static createFrom(source: any = {}) {
	        return new SaveDialogOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.defaultDirectory = source["defaultDirectory"];
	        this.defaultFilename = source["defaultFilename"];
	        this.title = source["title"];
	        this.filters = this.convertValues(source["filters"], FileFilter);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class UninstallPackageOpts {
	    id: number[];
	
	    static createFrom(source: any = {}) {
	        return new UninstallPackageOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class UpdatePreferencesOpts {
	    watchPath: string;
	    feature: Feature;
	
	    static createFrom(source: any = {}) {
	        return new UpdatePreferencesOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.watchPath = source["watchPath"];
	        this.feature = this.convertValues(source["feature"], Feature);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class UpdateProjectOpts {
	    id: number[];
	    name?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProjectOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}

}

export namespace enums {
	
	export enum PackageState {
	    AVAILABLE = "available",
	    DOWNLOADING = "downloading",
	    INCOMPLETE = "incomplete",
	    INSTALLED = "installed",
	    ERROR = "error",
	}
	export enum PackageType {
	    BUILD = "build",
	    ADDON = "addon",
	}

}

export namespace semver {
	
	export class Version {
	    Major: number;
	    Minor: number;
	    Patch: number;
	
	    static createFrom(source: any = {}) {
	        return new Version(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Major = source["Major"];
	        this.Minor = source["Minor"];
	        this.Patch = source["Patch"];
	    }
	}

}

export namespace types {
	
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
	export class Dependency {
	    reference: string;
	    type: enums.PackageType;
	
	    static createFrom(source: any = {}) {
	        return new Dependency(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.reference = source["reference"];
	        this.type = source["type"];
	    }
	}
	export class Media {
	    filePath: string;
	    url: string;
	    splash: boolean;
	    thumbnail: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Media(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.url = source["url"];
	        this.splash = source["splash"];
	        this.thumbnail = source["thumbnail"];
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
		    if (a.slice && a.map) {
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
	export class Progress {
	    currentBytes: number;
	    totalBytes: number;
	    bytesPerSecond: number;
	
	    static createFrom(source: any = {}) {
	        return new Progress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentBytes = source["currentBytes"];
	        this.totalBytes = source["totalBytes"];
	        this.bytesPerSecond = source["bytesPerSecond"];
	    }
	}
	export class URI {
	    Scheme: string;
	    Opaque: string;
	    // Go type: url
	    User?: any;
	    Host: string;
	    Path: string;
	    RawPath: string;
	    OmitHost: boolean;
	    ForceQuery: boolean;
	    RawQuery: string;
	    Fragment: string;
	    RawFragment: string;
	
	    static createFrom(source: any = {}) {
	        return new URI(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Scheme = source["Scheme"];
	        this.Opaque = source["Opaque"];
	        this.User = this.convertValues(source["User"], null);
	        this.Host = source["Host"];
	        this.Path = source["Path"];
	        this.RawPath = source["RawPath"];
	        this.OmitHost = source["OmitHost"];
	        this.ForceQuery = source["ForceQuery"];
	        this.RawQuery = source["RawQuery"];
	        this.Fragment = source["Fragment"];
	        this.RawFragment = source["RawFragment"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class Package {
	    id: number[];
	    type: enums.PackageType;
	    state: enums.PackageState;
	    reference: string;
	    name: string;
	    author: string;
	    tag: string;
	    path: string;
	    verified: boolean;
	    installationPath: string;
	    operations: string[];
	    platform: string;
	    uri?: URI;
	    version?: semver.Version;
	    progress?: Progress;
	    // Go type: time
	    updatedAt: any;
	
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
	        this.verified = source["verified"];
	        this.installationPath = source["installationPath"];
	        this.operations = source["operations"];
	        this.platform = source["platform"];
	        this.uri = this.convertValues(source["uri"], URI);
	        this.version = this.convertValues(source["version"], semver.Version);
	        this.progress = this.convertValues(source["progress"], Progress);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	
	export class Project {
	    id: number[];
	    name: string;
	    tags: string[];
	    path: string;
	    mediaPath: string;
	    fileName: string;
	    dependencies: Dependency[];
	    media: Media[];
	    strict: boolean;
	    version: string;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.tags = source["tags"];
	        this.path = source["path"];
	        this.mediaPath = source["mediaPath"];
	        this.fileName = source["fileName"];
	        this.dependencies = this.convertValues(source["dependencies"], Dependency);
	        this.media = this.convertValues(source["media"], Media);
	        this.strict = source["strict"];
	        this.version = source["version"];
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

