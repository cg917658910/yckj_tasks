<template>
  <div :class="['app', { 'app--authless': hideLayout }]">
    <Sidebar v-if="!hideLayout" />
    <main class="main">
      <RouterView />
    </main>
    <ToastHost />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import Sidebar from './components/Sidebar.vue'
import ToastHost from './components/ToastHost.vue'

const route = useRoute()
const hideLayout = computed(() => route.meta.hideLayout)
</script>
