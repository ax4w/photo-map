<script lang="ts">
  import Map from './lib/Map.svelte';
  import BigImage from './lib/BigImage.svelte';
  let selectedImage = $state("");
  let isModalOpen = $state(false);

  function handleSelectImage(imageUrl: string) {
    selectedImage = imageUrl;
    isModalOpen = true;
  }

  function handleCloseModal() {
    isModalOpen = false;
  }
  
</script>

<main>
  <Map onSelectImage={handleSelectImage} />
  <BigImage 
    isOpen={isModalOpen} 
    imageUrl={selectedImage} 
    onClose={handleCloseModal} 
  />
</main>