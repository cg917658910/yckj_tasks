<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">用户管理</div>
        <div class="panel-sub">查看用户积分与状态</div>
      </div>
    </div>

    <div class="filter-bar">
      <input v-model="keyword" placeholder="搜索用户名" />
      <select v-model.number="statusFilter">
        <option :value="-1">全部状态</option>
        <option :value="1">正常</option>
        <option :value="0">禁用</option>
      </select>
    </div>

    <div class="table">
      <div class="table-head">
        <div>用户名</div>
        <div>可用积分</div>
        <div>冻结积分</div>
        <div>累计提现</div>
        <div>状态</div>
        <div>操作</div>
      </div>
      <div class="table-row" v-for="item in list" :key="item.id">
        <div>{{ item.username }}</div>
        <div>{{ item.available_points }}</div>
        <div>{{ item.frozen_points }}</div>
        <div>{{ item.total_withdrawn || 0 }} 元</div>
        <div><span class="chip" :class="item.status === 1 ? 'online' : 'claimed'">{{ item.status === 1 ? '正常' : '禁用' }}</span></div>
        <div>
          <button class="link" @click="toggleStatus(item)">切换状态</button>
          <button class="link" @click="openAdjust(item)">调整积分</button>
        </div>
      </div>
    </div>

    <Pagination :total="total" :page="page" :page-size="pageSize" @update:page="onPageChange" />

    <div v-if="showAdjust" class="modal">
      <div class="modal-card">
        <div class="modal-title">调整用户积分</div>
        <div class="form-grid">
          <label>
            调整积分（正数/负数）
            <input v-model.number="adjustPoints" type="number" />
          </label>
          <label>
            备注
            <input v-model="adjustRemark" />
          </label>
        </div>
        <div class="modal-actions">
          <button class="ghost-btn" @click="closeAdjust">取消</button>
          <button class="primary-btn" @click="confirmAdjust">提交</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { fetchUsers, updateUserStatus, adjustUserPoints } from '../api/admin'
import { notify } from '../store/notify'
import Pagination from '../components/Pagination.vue'

const list = ref([])
const total = ref(0)
const keyword = ref('')
const statusFilter = ref(-1)
const page = ref(1)
const pageSize = ref(10)

const showAdjust = ref(false)
const current = ref(null)
const adjustPoints = ref(0)
const adjustRemark = ref('')

const load = async () => {
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (statusFilter.value >= 0) params.status = statusFilter.value
    if (keyword.value.trim()) params.keyword = keyword.value.trim()

    const res = await fetchUsers(params)
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

const toggleStatus = async (item) => {
  try {
    const nextStatus = item.status === 1 ? 0 : 1
    await updateUserStatus(item.id, { status: nextStatus })
    await load()
    notify('状态已更新', 'success')
  } catch (err) {
    notify(err.message || '更新失败', 'error')
  }
}

const openAdjust = (item) => {
  current.value = item
  adjustPoints.value = 0
  adjustRemark.value = ''
  showAdjust.value = true
}

const closeAdjust = () => {
  showAdjust.value = false
}

const confirmAdjust = async () => {
  try {
    await adjustUserPoints(current.value.id, {
      points: adjustPoints.value,
      remark: adjustRemark.value,
    })
    showAdjust.value = false
    await load()
    notify('积分已调整', 'success')
  } catch (err) {
    notify(err.message || '调整失败', 'error')
  }
}

onMounted(load)
</script>
