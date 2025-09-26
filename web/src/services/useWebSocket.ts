import { ref, onMounted, onUnmounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useMetricsStore } from '../stores/metrics';
import { useSystemStore } from '../stores/system';

export function useWebSocket() {
  const ws = ref<WebSocket | null>(null);
  const isConnected = ref(false);
  const error = ref<string | null>(null);

  const authStore = useAuthStore();
  const metricsStore = useMetricsStore();
  const systemStore = useSystemStore();

  const connect = () => {
    if (!authStore.accessToken) {
      error.value = "Authentication token not found.";
      return;
    }

    // Determine the WebSocket protocol
    const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    const wsUrl = `${proto}//${host}/ws?token=${authStore.accessToken}`;

    ws.value = new WebSocket(wsUrl);

    ws.value.onopen = () => {
      isConnected.value = true;
      error.value = null;
      console.log("WebSocket connection established.");
    };

    ws.value.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);

        // Route message to the correct store based on its type
        switch (message.type) {
          case 'system_info':
            systemStore.updateSystemInfo(message.payload);
            break;
          case 'v2ray_status':
            systemStore.updateV2rayStatus(message.payload);
            break;
          case 'metrics_update':
            metricsStore.updateMetrics(message.payload);
            break;
          case 'traffic_update':
            metricsStore.updateTraffic(message.payload);
            break;
          default:
            console.warn("Unknown WebSocket message type:", message.type);
        }
      } catch (e) {
        console.error("Failed to parse WebSocket message:", e);
      }
    };

    ws.value.onerror = (event) => {
      console.error("WebSocket error:", event);
      error.value = "WebSocket connection error.";
    };

    ws.value.onclose = () => {
      isConnected.value = false;
      console.log("WebSocket connection closed.");
      // Optional: implement reconnection logic here
    };
  };

  const disconnect = () => {
    if (ws.value) {
      ws.value.close();
      ws.value = null;
    }
  };

  onMounted(() => {
    // Connect when the component using this composable is mounted
    connect();
  });

  onUnmounted(() => {
    // Disconnect when the component is unmounted
    disconnect();
  });

  return {
    isConnected,
    error,
    connect,
    disconnect,
  };
}