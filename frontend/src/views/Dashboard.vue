<template>
  <div class="page-stack">
    <section class="metric-grid">
      <div class="metric-card">
        <span>CPU</span>
        <strong>{{ status.cpu }}%</strong>
      </div>
      <div class="metric-card">
        <span>Memory</span>
        <strong>{{ status.memory }}%</strong>
      </div>
      <div class="metric-card">
        <span>Disk</span>
        <strong>{{ status.disk }}%</strong>
      </div>
      <div class="metric-card">
        <span>Hostname</span>
        <strong class="hostname">{{ status.hostname || '-' }}</strong>
      </div>
    </section>

    <section class="panel">
      <div class="panel-header">
        <div>
          <h2>CPU 趋势</h2>
          <p>每 5 秒自动刷新，显示最近 12 次采样。</p>
        </div>
        <el-tag effect="dark">Uptime {{ status.uptime || '-' }}</el-tag>
      </div>
      <div ref="chartRef" class="chart"></div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import * as echarts from 'echarts'

import { getSystemStatus, type SystemStatus } from '@/api/system'

const status = reactive<SystemStatus>({
  cpu: 0,
  memory: 0,
  disk: 0,
  hostname: '',
  uptime: '',
})

const chartRef = ref<HTMLDivElement>()
let chart: echarts.ECharts | null = null
let timer: number | undefined
const cpuHistory = ref<number[]>([])
const labels = ref<string[]>([])

async function refreshStatus() {
  const { data } = await getSystemStatus()
  Object.assign(status, data)

  const now = new Date().toLocaleTimeString()
  cpuHistory.value = [...cpuHistory.value.slice(-11), data.cpu]
  labels.value = [...labels.value.slice(-11), now]
  renderChart()
}

function renderChart() {
  if (!chart) return
  chart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 32, right: 24, top: 24, bottom: 32 },
    xAxis: {
      type: 'category',
      data: labels.value,
      axisLine: { lineStyle: { color: '#3b4354' } },
      axisLabel: { color: '#8f9bb3' },
    },
    yAxis: {
      type: 'value',
      min: 0,
      max: 100,
      axisLabel: { color: '#8f9bb3', formatter: '{value}%' },
      splitLine: { lineStyle: { color: '#202938' } },
    },
    series: [
      {
        name: 'CPU',
        type: 'line',
        smooth: true,
        symbolSize: 6,
        data: cpuHistory.value,
        areaStyle: { color: 'rgba(64, 158, 255, 0.14)' },
        lineStyle: { color: '#409eff', width: 3 },
        itemStyle: { color: '#67c23a' },
      },
    ],
  })
}

function resizeChart() {
  chart?.resize()
}

onMounted(async () => {
  if (chartRef.value) {
    chart = echarts.init(chartRef.value, 'dark')
  }
  await refreshStatus()
  timer = window.setInterval(refreshStatus, 5000)
  window.addEventListener('resize', resizeChart)
})

onBeforeUnmount(() => {
  if (timer) window.clearInterval(timer)
  window.removeEventListener('resize', resizeChart)
  chart?.dispose()
})
</script>
