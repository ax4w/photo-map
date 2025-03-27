export interface Location {
    Lat: number;
    Long: number;
}

export interface Locations {
    [region: string]: Location;
}

export interface ImagesResponse {
    images: string[];
    has_more: boolean;
}

export const fetchLocations = async (): Promise<Locations> => {
    const response = await fetch('/api/regions/');
    if (!response.ok) {
        throw new Error('Failed to fetch locations');
    }
    return await response.json();
};

export const fetchImages = async (region: string, offset: number): Promise<ImagesResponse> => {
    const response = await fetch(`/api/images/${region}?offset=${offset}`);
    if (!response.ok) {
        throw new Error(`Failed to fetch images for region ${region}`);
    }
    return await response.json();
};

export const getThumbUrl = (region: string, image: string): string => {
    return `/thumbs/${region}/${image}`;
};

export const getImageUrl = (region: string, image: string): string => {
    return `/images/${region}/${image}`;
}; 