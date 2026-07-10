<template>
  <section class="panel">
    <div class="panel-header">
      <div>
        <h2>操作审计</h2>
        <p>记录 Docker 和 Terminal 等关键操作。</p>
      </div>
      <el-button :icon="Refresh" :loading="loading" @click="loadLogs">刷新</el-button>
    </div>

    <el-table v-loading="loading" :data="logs" class="dark-table" row-key="id">
      <el-table-column prop="username" label="用户" width="140" />
      <el-table-column prop="action" label="操作" width="180" />
      <el-table-column prop="target" label="目标" min-width="220" show-overflow-tooltip />
      <el-table-column prop="ip" label="IP" width="160" />
      <el-table-column label="时间" width="220">
        <template #default="{ row }">
          {{ new Date(row.created_at).toLocaleString() }}
        </template>
      </el-table-column>
    </el-table>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Refresh } from '@element-plus/icons-vue'

import { getAuditLogs, type AuditLog } from '@/api/audit'

const loading = ref(false)
const logs = ref<AuditLog[]>([])

async function loadLogs() {
  loading.value = true
  try {
    const { data } = await getAuditLogs()
    logs.value = data
  } finally {
    loading.value = false
  }
}

onMounted(loadLogs)
</script>
