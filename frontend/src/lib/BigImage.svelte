<script lang="ts">
    let { isOpen, imageUrl, onClose } = $props();

    let loading = $state(false);
    let imageLoaded = $state(false);
    
    $effect(() => {
        // run everytime when isOpen or imageUrl changes
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
        top: 20px;
        right: 20px;
        width: 40px;
        height: 40px;
        background: transparent;
        border: none;
        color: white;
        font-size: 32px;
        cursor: pointer;
        display: flex;
        justify-content: center;
        align-items: center;
        transition: all 0.2s ease;
        outline: none;
        padding: 0;
        text-align: center;
        line-height: 1;
    }

    .close:hover {
        transform: scale(1.1);
    }
</style>
    