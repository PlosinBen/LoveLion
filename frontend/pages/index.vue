<template>
  <div class="container">
    <header>
      <h1>🦁 LoveLion</h1>
      <p>Personal Bookkeeping & Travel Assistant</p>
    </header>

    <main>
      <div class="status-card">
        <h2>System Status</h2>
        <div class="status-item">
          <span>Backend API:</span>
          <span :class="backendStatus ? 'status-ok' : 'status-error'">
            {{ backendStatus ? '✓ Connected' : '✗ Disconnected' }}
          </span>
        </div>
      </div>

      <div class="info-card">
        <h2>Quick Links</h2>
        <ul>
          <li>📊 <a href="/ledger">Personal Ledger</a></li>
          <li>✈️ <a href="/trips">My Trips</a></li>
        </ul>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
const config = useRuntimeConfig()
const backendStatus = ref(false)

onMounted(async () => {
  try {
    const response = await fetch(`${config.public.apiBase}/health`)
    backendStatus.value = response.ok
  } catch {
    backendStatus.value = false
  }
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  min-height: 100vh;
  color: #eee;
}

.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

header {
  text-align: center;
  margin-bottom: 3rem;
}

header h1 {
  font-size: 3rem;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

header p {
  color: #888;
  font-size: 1.2rem;
}

.status-card,
.info-card {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

h2 {
  font-size: 1.2rem;
  margin-bottom: 1rem;
  color: #f5576c;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-ok {
  color: #4ade80;
}

.status-error {
  color: #f87171;
}

ul {
  list-style: none;
}

li {
  padding: 0.5rem 0;
}

a {
  color: #60a5fa;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}
</style>
