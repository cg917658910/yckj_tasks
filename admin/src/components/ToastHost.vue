<template>
  <div v-if="state.show" class="toast" :class="state.type">
    {{ state.message }}
  </div>
</template>

<script setup>
import { notifyState as state } from '../store/notify'
</script>
