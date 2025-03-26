import React, { useEffect, useRef, useState } from 'react';
import { fetchImages, getThumbUrl, getImageUrl } from '../services/api';

interface ImageGalleryProps {
    region: string;
    onSelectImage: (imageUrl: string) => void;
}

interface GalleryState {
    offset: number;
    hasMore: boolean;
    loading: boolean;
    images: string[];
}

const ImageGallery: React.FC<ImageGalleryProps> = ({ region, onSelectImage }) => {
    const [state, setState] = useState<GalleryState>({
        offset: 0,
        hasMore: true,
        loading: false,
        images: []
    });

    const galleryRef = useRef<HTMLDivElement>(null);

    const loadImages = async () => {
        if (!state.hasMore || state.loading) return;

        setState(prev => ({ ...prev, loading: true }));

        try {
            const data = await fetchImages(region, state.offset);
            setState(prev => ({
                ...prev,
                images: [...prev.images, ...data.images],
                offset: prev.offset + data.images.length,
                hasMore: data.has_more,
                loading: false
            }));
        } catch (error) {
            console.error('Error loading images:', error);
            setState(prev => ({ ...prev, loading: false }));
        }
    };

    useEffect(() => {
        loadImages();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [region]); // Only load initial images when region changes

    const handleScroll = () => {
        if (!galleryRef.current) return;
        
        const { scrollTop, scrollHeight, clientHeight } = galleryRef.current;
        if (scrollHeight - (scrollTop + clientHeight) < 200 && !state.loading && state.hasMore) {
            loadImages();
        }
    };

    return (
        <div>
            <h3 style={{ margin: '0 0 15px 10px' }}>{region.toUpperCase()}</h3>
            <div 
                className="image-gallery" 
                ref={galleryRef}
                onScroll={handleScroll}
            >
                {state.images.map((img, index) => (
                    <img
                        key={index}
                        className="thumbnail"
                        src={getThumbUrl(region, img)}
                        alt={`Thumbnail ${index}`}
                        onClick={() => onSelectImage(getImageUrl(region, img))}
                    />
                ))}
            </div>
            {state.loading && <div className="loading" style={{ display: 'block' }}>Loading...</div>}
        </div>
    );
};

export default ImageGallery; 