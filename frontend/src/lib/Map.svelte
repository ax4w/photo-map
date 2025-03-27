<script>
	import { Map, TileLayer, Marker, Popup, DivIcon } from 'sveaflet';
    import { fetchLocations } from '../api';
    
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
                    }}>
                        <Popup options={{
                            maxWidth: 800
                        }}>
                            <h1>location.region</h1>
                        </Popup>
                </DivIcon>
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
</style>