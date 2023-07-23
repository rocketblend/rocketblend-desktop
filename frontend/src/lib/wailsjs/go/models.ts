export namespace blendconfig {
	
	export class BlendConfig {
	    projectPath: string;
	    blendFileName: string;
	    // Go type: rocketfile
	    rocketFile?: any;
	
	    static createFrom(source: any = {}) {
	        return new BlendConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectPath = source["projectPath"];
	        this.blendFileName = source["blendFileName"];
	        this.rocketFile = this.convertValues(source["rocketFile"], null);
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
	    key: string;
	    blendFile?: blendconfig.BlendConfig;
	    settings?: projectsettings.ProjectSettings;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.blendFile = this.convertValues(source["blendFile"], blendconfig.BlendConfig);
	        this.settings = this.convertValues(source["settings"], projectsettings.ProjectSettings);
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

export namespace projectsettings {
	
	export class ThumbnailSettings {
	    width?: number;
	    height?: number;
	    startFrame?: number;
	    endFrame?: number;
	    renderType?: string;
	
	    static createFrom(source: any = {}) {
	        return new ThumbnailSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	        this.startFrame = source["startFrame"];
	        this.endFrame = source["endFrame"];
	        this.renderType = source["renderType"];
	    }
	}
	export class ProjectSettings {
	    name?: string;
	    tags?: string[];
	    thumbnailSettings?: ThumbnailSettings;
	    thumbnailFilePath?: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.tags = source["tags"];
	        this.thumbnailSettings = this.convertValues(source["thumbnailSettings"], ThumbnailSettings);
	        this.thumbnailFilePath = source["thumbnailFilePath"];
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

