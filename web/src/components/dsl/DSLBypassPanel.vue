<script setup lang="ts">
import { ref, onMounted } from 'vue';

// Mock data for DSL status - in a real implementation, this would come from an API call
const dslStatus = ref({
  status: 'loading...',
  speed_mbps: 0,
  snr_db: 0,
  description: 'Fetching DSL status...',
});

const isLoading = ref(true);
const errorMessage = ref('');

// Function to fetch DSL status from the new backend endpoint
async function fetchDslStatus() {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    // This is a placeholder for the actual API call
    // In a real app, you would use a service like axios or fetch
    // For now, we simulate a successful API call with mock data
    await new Promise(resolve => setTimeout(resolve, 1000)); // simulate network delay
    const response = {
      status: 'excellent',
      speed_mbps: 100,
      snr_db: 55,
      description: 'DSL connection is optimal (placeholder).',
    };
    dslStatus.value = response;
  } catch (error) {
    errorMessage.value = 'Failed to fetch DSL status. Please try again.';
    console.error(error);
  } finally {
    isLoading.value = false;
  }
}

// Function to trigger the DSL boost
async function triggerDslBoost() {
  alert('DSL Boost Initiated! (Placeholder)');
  // In a real app, this would make a POST request to the /api/dsl/boost endpoint
}

// Fetch status when the component is mounted
onMounted(() => {
  fetchDslStatus();
});
</script>

<template>
  <div class="dsl-bypass-panel p-4 bg-gray-100 dark:bg-gray-800 rounded-lg">
    <h3 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">Keenetic DSL Optimization</h3>

    <div v-if="isLoading" class="text-center">
      <p class="text-lg text-gray-600 dark:text-gray-300">Loading DSL Status...</p>
    </div>

    <div v-else-if="errorMessage" class="text-center text-red-500">
      <p>{{ errorMessage }}</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <div class="p-4 bg-white dark:bg-gray-700 rounded-md shadow">
        <h4 class="font-semibold text-gray-500 dark:text-gray-400">Status</h4>
        <p class="text-2xl font-bold text-green-500 capitalize">{{ dslStatus.status }}</p>
      </div>
      <div class="p-4 bg-white dark:bg-gray-700 rounded-md shadow">
        <h4 class="font-semibold text-gray-500 dark:text-gray-400">Connection Speed</h4>
        <p class="text-2xl font-bold text-blue-500">{{ dslStatus.speed_mbps }} <span class="text-lg">Mbps</span></p>
      </div>
      <div class="p-4 bg-white dark:bg-gray-700 rounded-md shadow">
        <h4 class="font-semibold text-gray-500 dark:text-gray-400">Signal-to-Noise Ratio</h4>
        <p class="text-2xl font-bold text-purple-500">{{ dslStatus.snr_db }} <span class="text-lg">dB</span></p>
      </div>
    </div>

    <div class="text-center">
      <button
        @click="triggerDslBoost"
        class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg transition-colors duration-300"
      >
        Boost Connection
      </button>
    </div>

    <p class="text-sm text-gray-500 dark:text-gray-400 mt-4 text-center">
      {{ dslStatus.description }}
    </p>
  </div>
</template>

<style scoped>
/* Scoped styles can be added here if needed */
</style>