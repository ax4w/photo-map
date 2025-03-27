<script>
	import { Map, TileLayer, Marker, Popup, DivIcon } from 'sveaflet';
    import { fetchLocations } from '../api';
  import Gallery from './Gallery.svelte';
    
    let locationsPromise = $state(fetchLocations());
</script>

<div id="body">
	<Map
		style="height: 100vh; width: 100%; background-color: #262626;"
		options={{
			center: [30, 0],
			zoom: 5,
			minZoom: 4,
			maxZoom: 6,
		}}>
		<TileLayer 
            url={"https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png"} 
            options = {{
                 attribution: "&copy; <a href='https://carto.com/'>Carto</a>",
                 noWrap : true
            }}/>
        {#await locationsPromise}
            <div>Loading...</div>
        {:then locations}
            {#each Object.entries(locations).map(([region, location]) => ({
                region,
                lat: location.Lat,
                lng: location.Long
            })) as location}
                <Marker latLng={[location.lat, location.lng]}>
                    <DivIcon options={{
                        className: 'emoji-pin',
                        html: '📍', 
                        iconSize: [30, 30],
                        iconAnchor: [15, 30]
                    }} />
                    <Popup options={{
                        maxWidth: 800
                    }}>
                        <Gallery region={location.region} />
                    </Popup>
                    
                </Marker>
            {/each}
        {/await}    
	</Map>
</div>

<style>
    #body {
        background-color: #262626;
        margin: 0;
        width: 100vw;
        height: 100vh;
        position: absolute;
        top: 0;
        left: 0;
        padding: 0;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
            'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
            sans-serif;
    }

    :global(.emoji-pin) {
        font-size: 24px; 
        line-height: 30px;
    } 


    :global(.leaflet-popup-content-wrapper) {
        max-width: 100% !important;
        width: auto !important;
        min-width: 400px !important;
        border-radius: 12px !important;
    }

    :global(.leaflet-popup-content) {
        overflow: auto;
        margin: 13px 19px;
        line-height: 1.4;
        overflow-x: hidden;
        width: auto !important;
    }
</style>