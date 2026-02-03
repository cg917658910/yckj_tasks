<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">积分历史</div>
        <div class="panel-sub">用户积分变更记录</div>
      </div>
    </div>

    <div class="filter-bar">
      <input v-model="keyword" placeholder="搜索用户名" />
      <input v-model="type" placeholder="类型（如 task_reward）" />
    </div>

    <div class="table">
      <div class="table-head">
        <div>用户</div>
        <div>变更积分</div>
        <div>类型</div>
        <div>备注</div>
        <div>时间</div>
      </div>
      <div class="table-row" v-for="item in list" :key="item.id">
        <div>{{ item.username }}</div>
        <div>{{ item.change_points }}</div>
        <div>{{ item.type }}</div>
        <div>{{ item.remark || '-' }}</div>
        <div class="time">{{ item.created_at }}</div>
      </div>
    </div>

    <Pagination :total="total" :page="page" :page-size="pageSize" @update:page="onPageChange" />
  </section>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { fetchPointsLogs } from '../api/admin'
import Pagination from '../components/Pagination.vue'

const list = ref([])
const total = ref(0)
const keyword = ref('')
const type = ref('')
const page = ref(1)
const pageSize = ref(10)

const load = async () => {
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (type.value.trim()) params.type = type.value.trim()
    if (keyword.value.trim()) params.keyword = keyword.value.trim()
    const res = await fetchPointsLogs(params)
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

watch([keyword, type], () => {
  page.value = 1
  load()
})

onMounted(load)
</script>
