<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import Home from "$lib/components/ui/pages/home.svelte";
  import Login from "$lib/components/ui/pages/login.svelte";
  import Register from "$lib/components/ui/pages/register.svelte";
  import UserProfile from "$lib/components/ui/pages/userProfile.svelte";
  import { isAuthenticated, checkAuthentication } from "$lib/../store";;
  import { fetchCurrentUser } from "$lib/components/ui/pages/userstore";
  import { onMount } from "svelte";

  onMount(() => {
    checkAuthentication();
    if (isAuthenticated) {
      fetchCurrentUser();
    }
  });
</script>

<Router>
  <div>
      <Route path="/:username" let:params>
      <UserProfile username={params.username} />      </Route>
    {#if $isAuthenticated}
    
      <Route path="/"><Home /></Route>

    {:else}
      <Route path="/"><Login /></Route>
      <Route path="/login"><Login /></Route>
      <Route path="/register"><Register /></Route>
    {/if}
  </div>
</Router>