<template>
  <div class="p-4 sm:p-6 lg:p-8">
    <h1 class="text-2xl font-bold text-gray-900 mb-6">Agent Monitoring</h1>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">Metrics Overview</h2>
        <MetricsChart />
      </div>
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">System Traffic</h2>
        <TrafficGraph />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import MetricsChart from '../components/monitoring/MetricsChart.vue';
import TrafficGraph from '../components/monitoring/TrafficGraph.vue';
</script>