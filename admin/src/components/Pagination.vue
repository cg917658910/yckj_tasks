<template>
  <div class="pagination" v-if="total > 0">
    <button class="ghost-btn" :disabled="page === 1" @click="change(page - 1)">上一页</button>
    <span class="page-info">第 {{ page }} / {{ totalPages }} 页</span>
    <button class="ghost-btn" :disabled="page === totalPages" @click="change(page + 1)">下一页</button>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  total: { type: Number, default: 0 },
  page: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
})

const emit = defineEmits(['update:page'])

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))

const change = (next) => {
  emit('update:page', next)
}
</script>
