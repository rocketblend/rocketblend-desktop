export type RadioOption = {
    value: number;
    key: string;
};

export type ImageSourceFormat = {
    src: string;
    w: number;
    h?: number;
};
  
export type ImageSources = {
    [key: string]: ImageSourceFormat[];
};
  
export type ImageSrc = {
    img: ImageSourceFormat;
    sources?: ImageSources;
};
  
export type ImgType = {
    src: ImageSrc;
    alt?: string;
};
