<template>
  <v-chart class="chart" :option="chartOption" :theme="themeStore.theme" autoresize />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { SankeyChart } from 'echarts/charts';
import { TitleComponent, TooltipComponent } from 'echarts/components';
import VChart from 'vue-echarts';
import { useMetricsStore } from '../../stores/metrics';
import { useThemeStore } from '../../stores/theme';
import 'echarts/theme/dark.js';

// Initialize ECharts components
use([CanvasRenderer, SankeyChart, TitleComponent, TooltipComponent]);

const metricsStore = useMetricsStore();
const themeStore = useThemeStore();

const chartOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    triggerOn: 'mousemove',
  },
  series: [
    {
      type: 'sankey',
      layout: 'none',
      emphasis: {
        focus: 'adjacency',
      },
      data: metricsStore.getTraffic.nodes,
      links: metricsStore.getTraffic.links,
      lineStyle: {
        color: 'source',
        curveness: 0.5,
      },
      label: {
        color: themeStore.theme === 'dark' ? '#fff' : '#333',
      },
    },
  ],
}));
</script>

<style scoped>
.chart {
  height: 400px;
}
</style>