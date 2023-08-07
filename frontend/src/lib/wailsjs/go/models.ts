export namespace projectservice {
	
	export class Project {
	    id?: string;
	    name?: string;
	    tags?: string[];
	    path?: string;
	    fileName?: string;
	    build?: string;
	    addons?: string[];
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

