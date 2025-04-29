<script>
    import { onMount } from 'svelte';
    import { initializeApiService } from './lib/services/api'; // Adjust path to your api.js file
  
    onMount(() => {
      initializeApiService();
    });
  </script>