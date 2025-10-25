<script lang="js">
  import RunfileModal from './RunfileModal.svelte';

  // Props using $props in runes mode
  const { runfile } = $props();

  // State
  let showModal = $state(false);

  $effect(() => {
    if (runfile && !runfile.identifier) {
      runfile.identifier = runfile.name || runfile.filename;
    }
  });

  // Handle card click to open modal
  function openModal() {
    showModal = true;
  }

  function closeModal() {
    showModal = false;
  }

  // Image fallback handler
  function handleImageError(event) {
    event.target.src = 'https://placehold.co/600x400/667788/ffffff?text=No+Image';
  }
</script>

<div class="card-container">
  <div 
    class="card" 
    onclick={openModal} 
    onkeydown={(e) => e.key === 'Enter' && openModal()} 
    tabindex="0" 
    role="button" 
    aria-label={`View details for ${runfile.name}`}
  >
    <div class="card-front" style="background-image: url('{runfile.background_url || ''}')">
      <div class="card-overlay">
        <p class="click-hint">Click for details</p>
      </div>
    </div>
  </div>
</div>

{#if showModal}
  <RunfileModal {runfile} onClose={closeModal} />
{/if}

<style>
  .card-container {
    width: 460px;
    height: 215px;
    perspective: 1000px;
    margin: 1rem;
  }
  
  .card {
    position: relative;
    width: 100%;
    height: 100%;
    cursor: pointer;
    box-shadow: var(--shadow-light, 0 4px 8px rgba(0,0,0,0.1));
    border-radius: 8px;
    overflow: hidden;
    transition: all 0.3s ease;
  }
  
  .card:hover {
    box-shadow: var(--shadow-medium, 0 8px 16px rgba(0,0,0,0.15));
    transform: translateY(-4px);
  }
  
  .card-front {
    width: 100%;
    height: 100%;
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
    position: relative;
  }
  
  .card-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.3));
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    transition: background 0.3s ease;
  }

  .card:hover .card-overlay {
    background: linear-gradient(135deg, rgba(0, 0, 0, 0.6), rgba(0, 0, 0, 0.4));
  }

  .click-hint {
    color: rgba(255, 255, 255, 0.8);
    font-size: 0.9rem;
    margin: 0;
    text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.7);
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .card:hover .click-hint {
    opacity: 1;
  }
  
  @media (max-width: 600px) {
    .card-container {
      width: 100%;
      max-width: 460px;
      margin: 0.5rem auto;
    }
  }
</style>