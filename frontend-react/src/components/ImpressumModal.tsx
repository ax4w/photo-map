import React from 'react';
import '../styles/Modal.css';

interface ImpressumModalProps {
    isOpen: boolean;
    onClose: () => void;
}

const ImpressumModal: React.FC<ImpressumModalProps> = ({ isOpen, onClose }) => {
    if (!isOpen) return null;

    return (
        <div 
            className="impressum-modal" 
            style={{ display: 'block' }}
            onClick={(e) => {
                if (e.target === e.currentTarget) onClose();
            }}
        >
            <span className="impressum-close" onClick={onClose}>&times;</span>
            <h2>Impressum</h2>
            <p>This is a photo map application that displays images based on geographical locations.</p>
            <p>Created as a demo project for visualizing photos on a map.</p>
            <p>© 2023 - All rights reserved</p>
        </div>
    );
};

export default ImpressumModal; 