<template>
  <section class="panel terminal-panel">
    <div class="panel-header">
      <div>
        <h2>Web Terminal</h2>
        <p>登录态 WebSocket 终端。</p>
      </div>
      <el-button :type="connected ? 'danger' : 'primary'" @click="toggleConnection">
        {{ connected ? '断开' : '连接' }}
      </el-button>
    </div>
    <div ref="terminalRef" class="terminal-box"></div>
  </section>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'

import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const terminalRef = ref<HTMLDivElement>()
const connected = ref(false)
let term: Terminal | null = null
let fitAddon: FitAddon | null = null
let socket: WebSocket | null = null

function createTerminal() {
  term = new Terminal({
    cursorBlink: true,
    fontFamily: 'JetBrains Mono, Consolas, monospace',
    fontSize: 13,
    theme: {
      background: '#0b0f14',
      foreground: '#dce6f4',
      cursor: '#67c23a',
    },
  })
  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  if (terminalRef.value) {
    term.open(terminalRef.value)
    fitAddon.fit()
  }
}

function terminalURL() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${window.location.host}/api/terminal/ws?token=${encodeURIComponent(auth.token)}`
}

function connect() {
  if (!term || connected.value) return
  socket = new WebSocket(terminalURL())
  socket.binaryType = 'arraybuffer'
  socket.onopen = () => {
    connected.value = true
    term?.writeln('\r\nConnected.\r\n')
  }
  socket.onmessage = (event) => {
    if (event.data instanceof ArrayBuffer) {
      term?.write(new Uint8Array(event.data))
    } else {
      term?.write(event.data)
    }
  }
  socket.onclose = () => {
    connected.value = false
    term?.writeln('\r\nDisconnected.\r\n')
  }
  socket.onerror = () => {
    term?.writeln('\r\nTerminal connection error.\r\n')
  }
}

function disconnect() {
  socket?.close()
  socket = null
  connected.value = false
}

function toggleConnection() {
  if (connected.value) disconnect()
  else connect()
}

function resize() {
  fitAddon?.fit()
}

onMounted(async () => {
  createTerminal()
  term?.onData((data) => socket?.send(data))
  await nextTick()
  connect()
  window.addEventListener('resize', resize)
})

onBeforeUnmount(() => {
  disconnect()
  term?.dispose()
  window.removeEventListener('resize', resize)
})
</script>
