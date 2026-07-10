<template>
  <section class="panel settings-panel">
    <div class="panel-header">
      <div>
        <h2>AI Provider</h2>
        <p>OpenAI Compatible API 基础配置。</p>
      </div>
    </div>

    <el-form label-position="top" :model="form" class="settings-form">
      <el-form-item label="Provider">
        <el-select v-model="form.provider">
          <el-option label="OpenAI Chat Completions / Compatible" value="openai_chat" />
          <el-option label="OpenAI Responses API" value="openai_responses" />
          <el-option label="Gemini GenerateContent" value="gemini" />
          <el-option label="Anthropic Messages API" value="anthropic" />
        </el-select>
      </el-form-item>
      <el-form-item label="Base URL">
        <el-input v-model="form.base_url" placeholder="https://api.openai.com/v1" />
      </el-form-item>
      <el-form-item label="API Key">
        <el-input v-model="form.api_key" type="password" show-password placeholder="sk-..." />
      </el-form-item>
      <el-form-item label="Model">
        <el-input v-model="form.model" placeholder="gpt-5-mini" />
      </el-form-item>
      <div class="form-actions">
        <el-button :loading="testing" @click="handleTest">测试连接</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </div>
    </el-form>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'

import { getAIConfig, saveAIConfig, testAIConfig, type AIConfig } from '@/api/ai'

const saving = ref(false)
const testing = ref(false)
const form = reactive<AIConfig>({
  name: 'default',
  provider: 'openai_chat',
  base_url: 'https://api.openai.com/v1',
  api_key: '',
  model: 'gpt-5-mini',
})

async function loadConfig() {
  const { data } = await getAIConfig()
  Object.assign(form, data)
}

async function handleSave() {
  saving.value = true
  try {
    await saveAIConfig(form)
    ElMessage.success('保存成功')
  } finally {
    saving.value = false
  }
}

async function handleTest() {
  testing.value = true
  try {
    await saveAIConfig(form)
    const { data } = await testAIConfig()
    if (data.success) {
      ElMessage.success(`连接成功：${data.model}`)
    } else {
      ElMessage.error(data.error || '连接失败')
    }
  } finally {
    testing.value = false
  }
}

onMounted(loadConfig)
</script>
