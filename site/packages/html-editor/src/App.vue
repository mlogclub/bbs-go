<script setup lang="ts">
import MEditor from './components/MEditor.vue'
import { ref } from 'vue'

const content = ref('')

// 主题切换
const toggleTheme = () => {
  const htmlElement = document.documentElement
  if (htmlElement.classList.contains('dark')) {
    htmlElement.classList.remove('dark')
  } else {
    htmlElement.classList.add('dark')
  }
}

</script>

<template>
  <div class="app">
    <div class="app-header">
      <button class="theme-toggle-btn" @click="toggleTheme">
        切换主题
      </button>
    </div>
    
    <MEditor 
      v-model="content" 
    />
    
    <div class="content-preview">
      <h4>HTML</h4>
      <div v-text="content"></div>
    </div>
  </div>
</template>

<style>
.app {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem;
  position: relative;
  color: var(--editor-text);
}

.app-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 1rem;
}

.theme-toggle-btn {
  padding: 0.5rem 1rem;
  background: var(--editor-toolbar-bg);
  border: 1px solid var(--editor-border);
  border-radius: 4px;
  color: var(--editor-text);
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
}

.theme-toggle-btn:hover {
  background: var(--editor-hover);
}

@media (prefers-color-scheme: dark) {
  body {
    background-color: #1a1a1a;
  }
}

html.dark body {
  background-color: #1a1a1a;
}

.test-instructions {
  background: var(--editor-bg);
  border: 1px solid var(--editor-border);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.test-instructions h3 {
  margin: 0 0 0.5rem 0;
  color: var(--editor-text);
}

.test-instructions ul {
  margin: 0;
  padding-left: 1.5rem;
}

.test-instructions li {
  margin-bottom: 0.5rem;
  color: var(--editor-text);
}

.content-preview {
  margin-top: 2rem;
  padding: 1rem;
  background: var(--editor-bg);
  border: 1px solid var(--editor-border);
  border-radius: 8px;
}

.content-preview h4 {
  margin: 0 0 1rem 0;
  color: var(--editor-text);
}
</style>
