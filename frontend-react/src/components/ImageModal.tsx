import React, { useEffect, useState } from 'react';
import '../styles/Modal.css';

interface ImageModalProps {
    isOpen: boolean;
    imageUrl: string;
    onClose: () => void;
}

const ImageModal: React.FC<ImageModalProps> = ({ isOpen, imageUrl, onClose }) => {
    const [loading, setLoading] = useState(false);
    const [imageLoaded, setImageLoaded] = useState(false);
    
    useEffect(() => {
        if (isOpen && imageUrl) {
            setLoading(true);
            setImageLoaded(false);
            
            const img = new Image();
            img.src = imageUrl;
            img.onload = () => {
                setLoading(false);
                setImageLoaded(true);
            };
            img.onerror = () => {
                setLoading(false);
                console.error('Failed to load image:', imageUrl);
                onClose();
            };
        }
    }, [imageUrl, isOpen, onClose]);

    if (!isOpen) return null;

    return (
        <div 
            className="modal" 
            style={{ display: 'flex' }}
            onClick={(e) => {
                if (e.target === e.currentTarget) onClose();
            }}
        >
            <span className="close" onClick={onClose}>&times;</span>
            {loading && <div className="loading-modal">Loading...</div>}
            <img 
                className="modal-img" 
                src={imageUrl} 
                alt="Enlarged view" 
                style={{ opacity: imageLoaded ? 1 : 0 }}
                onLoad={() => setImageLoaded(true)}
            />
        </div>
    );
};

export default ImageModal; 