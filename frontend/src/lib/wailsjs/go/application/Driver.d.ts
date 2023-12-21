// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {projectservice} from '../models';
import {uuid} from '../models';
import {config} from '../models';
import {packageservice} from '../models';
import {runtime} from '../models';

export function CreateProject(arg1:projectservice.CreateProjectRequest):Promise<void>;

export function DeleteProject(arg1:uuid.UUID):Promise<void>;

export function ExploreProject(arg1:uuid.UUID):Promise<void>;

export function GetApplicationConfig():Promise<config.Config>;

export function GetPackage(arg1:uuid.UUID):Promise<packageservice.GetPackageResponse>;

export function GetPlatform():Promise<runtime.Platform>;

export function GetProject(arg1:uuid.UUID):Promise<projectservice.GetProjectResponse>;

export function GetRocketBlendConfig():Promise<config.Config>;

export function ListPackages(arg1:string,arg2:string,arg3:boolean):Promise<packageservice.ListPackagesResponse>;

export function ListProjects(arg1:string):Promise<projectservice.ListProjectsResponse>;

export function Quit():Promise<void>;

export function RenderProject(arg1:uuid.UUID):Promise<void>;

export function RunProject(arg1:uuid.UUID):Promise<void>;

export function UpdateProject(arg1:projectservice.UpdateProjectRequest):Promise<void>;
