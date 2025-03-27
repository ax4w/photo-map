<script lang="ts">
    import { fetchImages, getThumbUrl } from "../api";

    let { region } = $props();
    
    let galleryElement: HTMLDivElement;

    let galleryState = $state({
        offset: 0,
        hasMore: true,
        loading: false,
        images: [] as string[]
    });

    async function loadImages() {
        if (!galleryState.hasMore || galleryState.loading) return;

        galleryState.loading = true;

        try {
            const data = await fetchImages(region, galleryState.offset);
            galleryState.images = [...galleryState.images, ...data.images];
            galleryState.offset += data.images.length;
            galleryState.hasMore = data.has_more;
            galleryState.loading = false;
        } catch (error) {
            console.error('Error loading images:', error);
            galleryState.loading = false;
        }
    }

    function handleScroll() {
        if (!galleryElement) return;

        const { scrollTop, scrollHeight, clientHeight } = galleryElement;
        if (scrollHeight - (scrollTop + clientHeight) < 200 && !galleryState.loading && galleryState.hasMore) {
            loadImages();
        }
    }

    function handleImageClick(img: string) {
        // TODO: Implement image selection logic
        console.log('Image clicked:', img);
    }

    $effect(() => {
        loadImages();
    });
    
</script>

<div>
    <h3 style="margin: 0 0 15px 10px">{region.toUpperCase()}</h3>
    <div 
        class="image-gallery" 
        bind:this={galleryElement}
        onscroll={handleScroll}
    >
        {#each galleryState.images as img, index}
            <img
                class="thumbnail"
                src={getThumbUrl(region, img)}
                alt={`Thumbnail ${index}`}
                onclick={() => handleImageClick(img)}
            />
        {/each}
        
        {#if galleryState.loading}
            <div class="loading" style="display: block">Loading...</div>
        {/if}
    </div>
</div>

<style>
    .image-gallery {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
    max-height: 30vh;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 10px;
    box-sizing: border-box;
    width: 100%;
    max-width: 780px;
}

.thumbnail {
    width: 100%;
    height: 120px;
    object-fit: cover;
    cursor: pointer;
    border-radius: 8px;
    transition: transform 0.2s;
    min-width: 120px;
    aspect-ratio: 1/1;
}

.thumbnail:hover {
    transform: scale(1.05);
}


.loading {
    display: none;
    position: absolute;
    bottom: 10px;
    left: 50%;
    transform: translateX(-50%);
    padding: 8px 20px;
    background: rgba(0,0,0,0.8);
    color: white;
    border-radius: 20px;
}
</style>