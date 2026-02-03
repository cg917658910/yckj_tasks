<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">积分历史</div>
        <div class="panel-sub">积分变更记录</div>
      </div>
    </div>

    <div class="history-item" v-for="item in list" :key="item.id">
      <div>
        <div class="task-title">{{ item.type }}</div>
        <div class="task-desc">{{ item.remark || '-' }}</div>
      </div>
      <div class="points">{{ item.change_points }}</div>
    </div>

    <EmptyState v-if="list.length === 0" title="暂无积分记录" desc="完成任务或提现后会显示记录" />

    <Pagination :total="total" :page="page" :page-size="pageSize" @update:page="onPageChange" />
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { fetchPointsLogs } from '../api/user'
import EmptyState from '../components/EmptyState.vue'
import Pagination from '../components/Pagination.vue'

const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const load = async () => {
  try {
    const res = await fetchPointsLogs({ page: page.value, page_size: pageSize.value })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (err) {
    console.error(err)
  }
}

const onPageChange = (next) => {
  page.value = next
  load()
}

onMounted(load)
</script>
