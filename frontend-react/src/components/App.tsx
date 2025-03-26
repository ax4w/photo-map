import React, { useState } from 'react';
import Map from './Map';
import ImageModal from './ImageModal';
import ImpressumModal from './ImpressumModal';
import '../styles/Map.css';
import '../styles/Modal.css';

const App: React.FC = () => {
    const [selectedImage, setSelectedImage] = useState<string>('');
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isImpressumOpen, setIsImpressumOpen] = useState(false);

    const handleSelectImage = (imageUrl: string) => {
        setSelectedImage(imageUrl);
        setIsModalOpen(true);
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
    };

    const handleToggleImpressum = () => {
        setIsImpressumOpen(!isImpressumOpen);
    };

    return (
        <div>
            <Map onSelectImage={handleSelectImage} />
            <ImageModal 
                isOpen={isModalOpen} 
                imageUrl={selectedImage} 
                onClose={handleCloseModal} 
            />
            <div className="impressum-link" onClick={handleToggleImpressum}>
                Impressum
            </div>
            <ImpressumModal 
                isOpen={isImpressumOpen} 
                onClose={() => setIsImpressumOpen(false)} 
            />
        </div>
    );
};

export default App; 