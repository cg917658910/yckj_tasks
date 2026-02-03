<template>
  <div class="empty-state">
    <div class="empty-title">{{ title }}</div>
    <div class="empty-desc">{{ desc }}</div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: '暂无数据' },
  desc: { type: String, default: '当前没有可展示内容' },
})
</script>
