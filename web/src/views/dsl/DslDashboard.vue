<template>
  <div class="p-4 sm:p-6 lg:p-8">
    <h1 class="text-3xl font-bold text-base-content">DSL Bypass Management</h1>
    <p class="mt-2 text-base-content/70">Monitor and control your Keenetic DSL optimization.</p>

    <div class="mt-8">
      <DSLBypassPanel />
    </div>
  </div>
</template>

<script setup lang="ts">
import DSLBypassPanel from '../../components/dsl/DSLBypassPanel.vue';

// This view acts as a container for the DSL panel component.
// Additional logic or components related to the DSL page can be added here in the future.
</script>