<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">任务审核</div>
        <div class="panel-sub">待审核成果与驳回记录</div>
      </div>
    </div>

    <div class="filter-bar">
      <input v-model="keyword" placeholder="搜索任务或用户" />
      <select v-model.number="statusFilter">
        <option :value="-1">全部状态</option>
        <option :value="2">待审核</option>
        <option :value="3">通过</option>
        <option :value="4">驳回</option>
      </select>
    </div>

    <div class="table">
      <div class="table-head">
        <div>任务</div>
        <div>提交用户</div>
        <div>状态</div>
        <div>提交时间</div>
        <div>操作</div>
      </div>
      <div class="table-row" v-for="item in list" :key="item.id">
        <div>
          <div class="task-title">{{ item.task_title }}</div>
          <div class="task-desc">奖励 {{ item.reward_points }} 积分</div>
        </div>
        <div>{{ item.username }}</div>
        <div>
          <span class="chip" :class="statusClass(item.status)">{{ statusText(item.status) }}</span>
        </div>
        <div class="time">{{ item.submitted_at || '-' }}</div>
        <div>
          <button class="link success" @click="openApprove(item)" :disabled="item.status !== 2">通过</button>
          <button class="link danger" @click="reject(item)" :disabled="item.status !== 2">驳回</button>
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

    <div v-if="showApprove" class="modal">
      <div class="modal-card">
        <div class="modal-title">审核通过</div>
        <div class="form-grid">
          <label>
            奖励积分（可调整）
            <input type="number" v-model.number="approvePoints" />
          </label>
          <label class="full">
            通过备注
            <input v-model="approveRemark" placeholder="可选，记录审核说明" />
          </label>
        </div>
        <div class="modal-actions">
          <button class="ghost-btn" @click="closeApprove">取消</button>
          <button class="primary-btn" @click="confirmApprove">确定</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { fetchClaims, approveClaim, rejectClaim } from '../api/admin'
import Pagination from '../components/Pagination.vue'
import { notify } from '../store/notify'

const list = ref([])
const total = ref(0)
const keyword = ref('')
const statusFilter = ref(-1)
const page = ref(1)
const pageSize = ref(10)

const showReject = ref(false)
const showApprove = ref(false)
const current = ref(null)
const rejectReason = ref('')
const approvePoints = ref(0)
const approveRemark = ref('')

const load = async () => {
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (statusFilter.value >= 0) params.status = statusFilter.value
    if (keyword.value.trim()) params.keyword = keyword.value.trim()

    const res = await fetchClaims(params)
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

const openApprove = (item) => {
  if (item.status !== 2) return
  current.value = item
  approvePoints.value = item.reward_points || 0
  approveRemark.value = ''
  showApprove.value = true
}

const reject = (item) => {
  if (item.status !== 2) return
  current.value = item
  rejectReason.value = ''
  showReject.value = true
}

const confirmReject = async () => {
  try {
    await rejectClaim(current.value.id, { reason: rejectReason.value })
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

const confirmApprove = async () => {
  try {
    await approveClaim(current.value.id, { reward_points: approvePoints.value, remark: approveRemark.value })
    showApprove.value = false
    await load()
    notify('审核通过', 'success')
  } catch (err) {
    notify(err.message || '审核通过失败', 'error')
  }
}

const closeApprove = () => {
  showApprove.value = false
}

const statusText = (status) => {
  if (status === 2) return '待审核'
  if (status === 3) return '通过'
  if (status === 4) return '驳回'
  return '进行中'
}

const statusClass = (status) => {
  if (status === 2) return 'pending'
  if (status === 3) return 'online'
  if (status === 4) return 'claimed'
  return 'pending'
}

onMounted(load)
</script>
