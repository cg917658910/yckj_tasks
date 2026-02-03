<template>
  <div class="notice" :class="type">
    <span class="notice-dot"></span>
    <span>{{ text }}</span>
  </div>
</template>

<script setup>
defineProps({
  text: { type: String, required: true },
  type: { type: String, default: 'info' },
})
</script>
