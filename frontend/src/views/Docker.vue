<template>
  <section class="panel">
    <div class="panel-header">
      <div>
        <h2>容器管理</h2>
        <p>通过 Docker SDK 管理本机容器。</p>
      </div>
      <el-button :icon="Refresh" :loading="loading" @click="loadContainers">刷新</el-button>
    </div>

    <el-table v-loading="loading" :data="containers" class="dark-table" row-key="id">
      <el-table-column prop="id" label="Container ID" width="150" />
      <el-table-column prop="name" label="Name" min-width="160" />
      <el-table-column prop="image" label="Image" min-width="220" show-overflow-tooltip />
      <el-table-column label="Status" width="130">
        <template #default="{ row }">
          <el-tag :type="row.status === 'running' ? 'success' : 'info'" effect="dark">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="320" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="success" :icon="VideoPlay" @click="runAction('start', row.id)">
            启动
          </el-button>
          <el-button size="small" :icon="VideoPause" @click="runAction('stop', row.id)">停止</el-button>
          <el-button size="small" type="warning" :icon="RefreshRight" @click="runAction('restart', row.id)">
            重启
          </el-button>
          <el-popconfirm title="确认删除该容器？" confirm-button-text="删除" @confirm="runAction('remove', row.id)">
            <template #reference>
              <el-button size="small" type="danger" :icon="Delete">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Delete, Refresh, RefreshRight, VideoPause, VideoPlay } from '@element-plus/icons-vue'

import {
  getContainers,
  removeContainer,
  restartContainer,
  startContainer,
  stopContainer,
  type DockerContainer,
} from '@/api/docker'

type Action = 'start' | 'stop' | 'restart' | 'remove'

const containers = ref<DockerContainer[]>([])
const loading = ref(false)

async function loadContainers() {
  loading.value = true
  try {
    const { data } = await getContainers()
    containers.value = data
  } finally {
    loading.value = false
  }
}

async function runAction(action: Action, id: string) {
  const actionMap = {
    start: startContainer,
    stop: stopContainer,
    restart: restartContainer,
    remove: removeContainer,
  }

  await actionMap[action](id)
  ElMessage.success('操作成功')
  await loadContainers()
}

onMounted(loadContainers)
</script>
