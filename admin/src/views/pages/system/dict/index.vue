<template>
  <div>
    <div class="dict-wrapper">
      <div class="type-container">
        <DictType />
      </div>
      <div class="dict-container">
        <Dict />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { defineOptions } from 'vue';
  import DictType from './components/DictType.vue';
  import Dict from './components/Dict.vue';

  defineOptions({
    name: 'Dict',
  });
</script>

<style scoped lang="less">
  .dict-wrapper {
    display: flex;
    .type-container {
      width: 300px;
    }
    .dict-container {
      flex: 1;
    }
  }
</style>
