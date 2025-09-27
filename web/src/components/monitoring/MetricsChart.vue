<template>
  <v-chart class="chart" :option="chartOption" autoresize />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { LineChart } from 'echarts/charts';
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
} from 'echarts/components';
import VChart from 'vue-echarts';
import { useMetricsStore } from '../../stores/metrics';

// Initialize ECharts components
use([
  CanvasRenderer,
  LineChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
]);

const metricsStore = useMetricsStore();

const chartOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
  },
  legend: {
    data: metricsStore.getMetrics.datasets.map(d => d.label),
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true,
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: metricsStore.getMetrics.labels,
  },
  yAxis: {
    type: 'value',
  },
  series: metricsStore.getMetrics.datasets.map(dataset => ({
    name: dataset.label,
    type: 'line',
    stack: 'total',
    areaStyle: {},
    emphasis: {
      focus: 'series',
    },
    data: dataset.data,
    smooth: true,
  })),
}));
</script>

<style scoped>
.chart {
  height: 400px;
}
</style>