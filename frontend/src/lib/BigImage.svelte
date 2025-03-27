<script lang="ts">
    let { isOpen, imageUrl, onClose } = $props();

    let loading = $state(false);
    let imageLoaded = $state(false);
    
    $effect(() => {
        if (isOpen && imageUrl) {
            loading = true;
            imageLoaded = false;
            
            const img = new Image();
            img.src = imageUrl;
            img.onload = () => {
                loading = false;
                imageLoaded = true;
            }   
            img.onerror = () => {
                loading = false;
                console.error('Failed to load image:', imageUrl);
                onClose();
            }
        }
    })  
</script>

<div>
    {#if isOpen && loading}
        <div class="loading-modal">Loading...</div>
    {/if}
    {#if isOpen && imageLoaded}
        <div class="modal">
            <button class="close" onclick={onClose}>&times;</button>
            <img 
                class="modal-img" 
                src={imageUrl} 
                alt="Enlarged view" 
            />
        </div>
    {/if}
</div>

<style>
    .modal {
        display: flex;
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0,0,0,0.9);
        justify-content: center;
        align-items: center;
        z-index: 1000;
    }

    .modal-img {
        max-width: 90%;
        max-height: 90%;
        object-fit: contain;
        transition: opacity 0.2s;
    }

    .loading-modal {
        display: flex;
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0,0,0,0.8);
        color: white;
        font-size: 1.5em;
        text-shadow: 0 0 5px rgba(0,0,0,0.5);
        justify-content: center;
        align-items: center;
        z-index: 1000;
    }

    .close {
        position: absolute;
        top: 25px;
        right: 35px;
        color: white;
        font-size: 45px;
        cursor: pointer;
        text-shadow: 0 0 5px rgba(0,0,0,0.5);
    }
</style>
    