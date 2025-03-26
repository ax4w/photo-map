import React, { useEffect, useState } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import L from 'leaflet';
import '../styles/Map.css';
import { fetchLocations, Locations } from '../services/api';
import ImageGallery from './ImageGallery';

interface MapProps {
    onSelectImage: (imageUrl: string) => void;
}

const Map: React.FC<MapProps> = ({ onSelectImage }) => {
    const [locations, setLocations] = useState<Locations>({});
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const loadLocations = async () => {
            try {
                const data = await fetchLocations();
                setLocations(data);
                setLoading(false);
            } catch (err) {
                console.error('Failed to load locations:', err);
                setError('Failed to load map data');
                setLoading(false);
            }
        };

        loadLocations();
    }, []);

    const pinIcon = L.divIcon({
        className: 'emoji-pin',
        html: '📍', 
        iconSize: [30, 30],
        iconAnchor: [15, 30]
    });

    if (loading) {
        return <div style={{ color: 'white', textAlign: 'center', paddingTop: '50px' }}>Loading map...</div>;
    }

    if (error) {
        return <div style={{ color: 'white', textAlign: 'center', paddingTop: '50px' }}>Error: {error}</div>;
    }

    return (
        <MapContainer 
            center={[30, 0]} 
            zoom={2} 
            minZoom={2} 
            maxZoom={6}
            id="map"
        >
            <TileLayer
                url="https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png"
                attribution="&copy; <a href='https://carto.com/'>Carto</a>"
                noWrap={true}
            />
            {Object.entries(locations).map(([region, { Lat, Long }]) => (
                <Marker 
                    key={region} 
                    position={[Lat, Long]} 
                    icon={pinIcon}
                >
                    <Popup maxWidth={800}>
                        <ImageGallery 
                            region={region} 
                            onSelectImage={onSelectImage} 
                        />
                    </Popup>
                </Marker>
            ))}
        </MapContainer>
    );
};

export default Map; 