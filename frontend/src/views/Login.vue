<template>
  <div class="login-page">
    <div class="login-panel">
      <div class="login-copy">
        <div class="brand-row">
          <div class="brand-mark">AI</div>
          <span>AIPanel</span>
        </div>
        <h1>服务器面板</h1>
        <p>登录后管理系统状态和 Docker 容器。</p>
      </div>

      <el-form ref="formRef" :model="form" :rules="rules" class="login-card" @submit.prevent="handleLogin">
        <h2>登录</h2>
        <el-form-item prop="username">
          <el-input v-model="form.username" size="large" placeholder="用户名" :prefix-icon="User" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            size="large"
            type="password"
            show-password
            placeholder="密码"
            :prefix-icon="Lock"
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-button type="primary" size="large" native-type="submit" :loading="loading" class="login-button">
          登录
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { Lock, User } from '@element-plus/icons-vue'

import { login } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  username: 'admin',
  password: 'admin123',
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const { data } = await login(form)
    auth.setToken(data.token)
    router.push('/dashboard')
  } finally {
    loading.value = false
  }
}
</script>
