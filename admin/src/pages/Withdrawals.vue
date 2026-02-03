<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">提现管理</div>
        <div class="panel-sub">审核提现申请并打款</div>
      </div>
    </div>

    <div class="filter-bar">
      <input v-model="keyword" placeholder="搜索用户名" />
      <select v-model.number="statusFilter">
        <option :value="-1">全部状态</option>
        <option :value="1">待审核</option>
        <option :value="2">已打款</option>
        <option :value="3">驳回</option>
      </select>
    </div>

    <div class="table">
      <div class="table-head">
        <div>用户</div>
        <div>申请金额</div>
        <div>消耗积分</div>
        <div>状态</div>
        <div>操作</div>
      </div>
      <div class="table-row" v-for="item in list" :key="item.id">
        <div>{{ item.username }}</div>
        <div>{{ item.amount }} 元</div>
        <div>{{ item.points_cost }}</div>
        <div><span class="chip" :class="statusClass(item.status)">{{ statusText(item.status) }}</span></div>
        <div>
          <button class="link success" @click="pay(item)">打款</button>
          <button class="link danger" @click="openReject(item)">驳回</button>
        </div>
      </div>
    </div>

    <Pagination :total="total" :page="page" :page-size="pageSize" @update:page="onPageChange" />

    <div v-if="showReject" class="modal">
      <div class="modal-card">
        <div class="modal-title">填写驳回原因</div>
        <textarea v-model="rejectReason" rows="4" />
        <div class="modal-actions">
          <button class="ghost-btn" @click="closeReject">取消</button>
          <button class="primary-btn" @click="confirmReject">提交</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { fetchWithdrawals, payWithdrawal, rejectWithdrawal } from '../api/admin'
import { notify } from '../store/notify'
import Pagination from '../components/Pagination.vue'

const list = ref([])
const total = ref(0)
const keyword = ref('')
const statusFilter = ref(-1)
const page = ref(1)
const pageSize = ref(10)

const showReject = ref(false)
const current = ref(null)
const rejectReason = ref('')

const load = async () => {
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (statusFilter.value >= 0) params.status = statusFilter.value
    if (keyword.value.trim()) params.keyword = keyword.value.trim()

    const res = await fetchWithdrawals(params)
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

watch([keyword, statusFilter], () => {
  page.value = 1
  load()
})

const pay = async (item) => {
  try {
    await payWithdrawal(item.id)
    await load()
    notify('已打款', 'success')
  } catch (err) {
    notify(err.message || '打款失败', 'error')
  }
}

const openReject = (item) => {
  current.value = item
  rejectReason.value = ''
  showReject.value = true
}

const confirmReject = async () => {
  try {
    await rejectWithdrawal(current.value.id, { reason: rejectReason.value })
    showReject.value = false
    await load()
    notify('已驳回', 'success')
  } catch (err) {
    notify(err.message || '驳回失败', 'error')
  }
}

const closeReject = () => {
  showReject.value = false
}

const statusText = (status) => {
  if (status === 2) return '已打款'
  if (status === 3) return '驳回'
  return '待审核'
}

const statusClass = (status) => {
  if (status === 2) return 'online'
  if (status === 3) return 'claimed'
  return 'pending'
}

onMounted(load)
</script>
