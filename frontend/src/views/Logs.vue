<template>
  <section class="panel">
    <div class="panel-header">
      <div>
        <h2>系统日志</h2>
        <p>读取 systemd 和 kernel 日志。</p>
      </div>
      <div class="toolbar">
        <el-select v-model="type" style="width: 140px" @change="loadLogs">
          <el-option label="System" value="system" />
          <el-option label="Kernel" value="kernel" />
        </el-select>
        <el-input v-model="keyword" placeholder="搜索日志" clearable style="width: 220px" />
        <el-button :loading="loading" :icon="Refresh" @click="loadLogs">刷新</el-button>
        <el-button :icon="MagicStick" disabled>AI分析</el-button>
      </div>
    </div>

    <div v-loading="loading" class="log-viewer">
      <div v-for="(line, index) in filteredLogs" :key="index" class="log-line">
        {{ line }}
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { MagicStick, Refresh } from '@element-plus/icons-vue'

import { getSystemLogs } from '@/api/logs'

const type = ref('system')
const keyword = ref('')
const loading = ref(false)
const logs = ref<string[]>([])

const filteredLogs = computed(() => {
  if (!keyword.value) return logs.value
  return logs.value.filter((line) => line.toLowerCase().includes(keyword.value.toLowerCase()))
})

async function loadLogs() {
  loading.value = true
  try {
    const { data } = await getSystemLogs(type.value, 100)
    logs.value = data.logs
  } finally {
    loading.value = false
  }
}

onMounted(loadLogs)
</script>
