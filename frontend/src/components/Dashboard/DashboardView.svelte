<script>
    import QuickActionsCard from './cards/QuickActionsCard.svelte';
    import ConsoleCard from './cards/ConsoleCard.svelte';
    import SystemInfoCard from './cards/SystemInfoCard.svelte';
    import LogsCard from './cards/LogsCard.svelte';
    import QuickActionsCardAlternative from './cards/QuickActionsCardAlternative.svelte';
</script>
  
  <div class="dashboard-container">
    <div class="dashboard-grid">
      <!-- <QuickActionsCard /> -->
      <QuickActionsCardAlternative />
      <SystemInfoCard />
      <ConsoleCard />
      <LogsCard />

    </div>
  </div>
  
  <style>
    .dashboard-container {
        min-height: 100%;
        width: 100%;
        background-image: url('https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/544550/ss_e81cd4d2f058ccc070295edae95beef36a25822e.1920x1080.jpg?t=1752831512');
        background-size: cover;
        background-repeat: no-repeat;
        background-position: center;
        border-radius: 8px;
    }


    .dashboard-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        grid-auto-rows: minmax(0, 400px);
        gap: 1.5rem;
        padding: 1rem;
    }
  
    @media (max-width: 1024px) {
        .dashboard-grid {
            grid-template-columns: 1fr;
        }
    }
  
    @media (min-width: 1025px) and (max-width: 1400px) {
        .dashboard-grid {
            grid-template-columns: repeat(2, 1fr);
        }
    }
  </style>