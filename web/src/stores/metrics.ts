import { defineStore } from 'pinia';

interface MetricsData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    backgroundColor: string;
    borderColor: string;
    borderWidth: number;
  }[];
}

interface TrafficNode {
  id: string;
  name: string;
}

interface TrafficLink {
  source: string;
  target: string;
  value: number;
}

interface TrafficData {
  nodes: TrafficNode[];
  links: TrafficLink[];
}

export const useMetricsStore = defineStore('metrics', {
  state: () => ({
    // Mock data for metrics chart
    metricsData: {
      labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
      datasets: [
        {
          label: 'CPU Usage',
          data: [65, 59, 80, 81, 56, 55, 40],
          backgroundColor: 'rgba(75, 192, 192, 0.2)',
          borderColor: 'rgba(75, 192, 192, 1)',
          borderWidth: 1,
        },
        {
          label: 'Memory Usage',
          data: [28, 48, 40, 19, 86, 27, 90],
          backgroundColor: 'rgba(153, 102, 255, 0.2)',
          borderColor: 'rgba(153, 102, 255, 1)',
          borderWidth: 1,
        },
      ],
    } as MetricsData,
    // Mock data for traffic graph
    trafficData: {
      nodes: [
        { id: '1', name: 'Agent' },
        { id: '2', name: 'Service A' },
        { id: '3', name: 'Service B' },
        { id: '4', name: 'Database' },
      ],
      links: [
        { source: '1', target: '2', value: 5 },
        { source: '1', target: '3', value: 8 },
        { source: '2', target: '4', value: 3 },
        { source: '3', target: '4', value: 6 },
      ],
    } as TrafficData,
  }),

  getters: {
    getMetrics(state): MetricsData {
      return state.metricsData;
    },
    getTraffic(state): TrafficData {
      return state.trafficData;
    },
  },

  actions: {
    updateMetrics(newMetrics: { cpu: number[]; memory: number[] }) {
      // This action will be called by the WebSocket client later
      this.metricsData.datasets[0].data = newMetrics.cpu;
      this.metricsData.datasets[1].data = newMetrics.memory;
    },
    updateTraffic(newTraffic: TrafficData) {
      // This action will be called by the WebSocket client later
      this.trafficData = newTraffic;
    },
  },
});