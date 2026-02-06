<template>
  <section>
    <header class="topbar">
      <div class="title">任务列表</div>
      <div class="search">
        <input v-model="keyword" placeholder="搜索任务..." />
      </div>
      <select class="filter" v-model="filter">
        <option value="all">所有任务</option>
        <option value="available">可领取</option>
        <option value="claimed">已领取</option>
      </select>
    </header>

    <NoticeBar text="领取限制：同一时间只能领取 1 个任务，完成后可继续领取。" />

    <section class="cards">
      <article class="task-card" v-for="task in filtered" :key="task.id" @click="goDetail(task)">
        <div class="task-head">
          <div>
            <div class="task-title">{{ task.title }}</div>
            <div class="task-desc">{{ task.summary }}</div>
          </div>
          <span class="badge" v-if="task.claimed">已领取</span>
          <span class="badge" v-else>可领取</span>
        </div>
        <div class="task-foot">
          <div class="points">⭐ {{ task.reward_points }} 积分</div>
          <button
            class="primary-btn"
            :class="{ disabled: task.claimed }"
            @click.stop="handleClaim(task)"
            :disabled="task.claimed"
          >
            {{ task.claimed ? '已领取' : '领取任务' }}
          </button>
        </div>
      </article>
    </section>

    <EmptyState v-if="filtered.length === 0" title="暂无任务" desc="当前暂无可领取任务或搜索结果为空" />
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { claimTask, fetchTasks } from '../api/user'
import EmptyState from '../components/EmptyState.vue'
import NoticeBar from '../components/NoticeBar.vue'
import { notify } from '../store/notify'

const router = useRouter()
const keyword = ref('')
const filter = ref('all')
const tasks = ref([])

const load = async () => {
  try {
    const res = await fetchTasks()
    tasks.value = res.data?.list || []
  } catch (err) {
    console.error(err)
  }
}

const handleClaim = async (task) => {
  try {
   const res =  await claimTask(task.id)
   if (res.code !== 0) {
     throw new Error(res.message || '领取失败')
   }
    await load()
    notify('任务领取成功', 'success')
  } catch (err) {
    notify(err.message || '领取失败', 'error')
  }
}

const goDetail = (task) => {
  router.push(`/tasks/${task.id}`)
}

const filtered = computed(() => {
  const kw = keyword.value.trim()
  return tasks.value.filter((item) => {
    const matchKeyword = !kw || item.title?.includes(kw)
    const matchStatus =
      filter.value === 'all' ||
      (filter.value === 'available' && !item.claimed) ||
      (filter.value === 'claimed' && item.claimed)
    return matchKeyword && matchStatus
  })
})

onMounted(load)
</script>
