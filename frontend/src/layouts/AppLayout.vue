<template>
  <el-container class="shell">
    <el-aside width="232px" class="sidebar">
      <div class="brand">
        <div class="brand-mark">AI</div>
        <div>
          <div class="brand-name">AIPanel</div>
          <div class="brand-subtitle">Server Panel</div>
        </div>
      </div>

      <el-menu :default-active="activePath" router class="nav-menu">
        <el-menu-item index="/dashboard">
          <el-icon><Monitor /></el-icon>
          <span>Dashboard</span>
        </el-menu-item>
        <el-menu-item index="/docker">
          <el-icon><Box /></el-icon>
          <span>Docker</span>
        </el-menu-item>
        <el-menu-item index="/terminal">
          <el-icon><Operation /></el-icon>
          <span>Terminal</span>
        </el-menu-item>
        <el-menu-item index="/logs">
          <el-icon><Document /></el-icon>
          <span>Logs</span>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>Settings</span>
        </el-menu-item>
        <el-menu-item index="/audit">
          <el-icon><Tickets /></el-icon>
          <span>Audit</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div>
          <h1>{{ pageTitle }}</h1>
          <p>Lightweight AI-first Server Panel</p>
        </div>
        <el-button :icon="SwitchButton" plain @click="handleLogout">退出</el-button>
      </el-header>

      <el-main class="content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Box, Document, Monitor, Operation, Setting, SwitchButton, Tickets } from '@element-plus/icons-vue'

import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const activePath = computed(() => route.path)
const pageTitle = computed(() => {
  if (route.path.startsWith('/docker')) return 'Docker'
  if (route.path.startsWith('/terminal')) return 'Terminal'
  if (route.path.startsWith('/logs')) return 'Logs'
  if (route.path.startsWith('/settings')) return 'Settings'
  if (route.path.startsWith('/audit')) return 'Audit'
  return 'Dashboard'
})

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>
