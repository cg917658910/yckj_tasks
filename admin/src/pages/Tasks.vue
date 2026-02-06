<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">任务管理</div>
        <div class="panel-sub">发布、编辑、下架任务</div>
      </div>
      <button class="primary-btn" @click="openCreate">发布任务</button>
    </div>

    <div class="filter-bar">
      <input v-model="keyword" placeholder="搜索任务标题" />
      <select v-model.number="statusFilter">
        <option :value="-1">全部状态</option>
        <option :value="1">上架</option>
        <option :value="2">下架</option>
        <option :value="0">草稿</option>
      </select>
    </div>

    <div class="table">
      <div class="table-head">
        <div>任务名称</div>
        <div>奖励积分</div>
        <div>领取状态</div>
        <div>状态</div>
        <div>创建时间</div>
        <div>操作</div>
      </div>
      <div class="table-row" v-for="item in list" :key="item.id">
        <div>
          <div class="task-title">{{ item.title }}</div>
          <div class="task-desc">{{ item.summary }}</div>
        </div>
        <div class="points">{{ item.reward_points }} 积分</div>
        <div>
          <span class="chip" :class="item.claimed ? 'claimed' : 'online'">
            {{ item.claimed ? '已领取' : '未领取' }}
          </span>
        </div>
        <div><span class="chip" :class="statusClass(item.status)">{{ statusText(item.status) }}</span></div>
        <div class="time">{{ item.created_at }}</div>
        <div>
          <button class="link" @click="edit(item)">编辑</button>
          <button v-if="item.status !== 1" class="link success" @click="on(item)">上架</button>
          <button v-else class="link danger" @click="off(item)">下架</button>
        </div>
      </div>
    </div>

    <Pagination :total="total" :page="page" :page-size="pageSize" @update:page="onPageChange" />

    <div v-if="showForm" class="modal">
      <div class="modal-card">
        <div class="modal-title">{{ form.id ? '编辑任务' : '发布任务' }}</div>
        <div class="form-grid">
          <label>
            标题
            <input v-model="form.title" />
          </label>
          <label>
            简介
            <input v-model="form.summary" />
          </label>
          <label>
            奖励积分
            <input v-model.number="form.reward_points" type="number" />
          </label>
          <label class="full">
            详情
            <textarea v-model="form.detail" rows="4" />
          </label>
          <label class="full">
            文档附件
            <input type="file" @change="handleUpload" />
            <div class="upload-tip" v-if="uploading">上传中...</div>
            <div class="upload-tip" v-else-if="form.doc_url">已上传：{{ form.doc_url }}</div>
          </label>
          <label class="full">
            文档链接
            <input v-model="form.doc_url" placeholder="可粘贴文件链接" />
          </label>
        </div>
        <div class="modal-actions">
          <button class="ghost-btn" @click="close">取消</button>
          <button class="primary-btn" @click="save">保存</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { createTask, fetchTasks, offTask, onTask, updateTask, uploadImage } from '../api/admin'
import Pagination from '../components/Pagination.vue'
import { notify } from '../store/notify'

const list = ref([])
const total = ref(0)
const keyword = ref('')
const statusFilter = ref(-1)
const page = ref(1)
const pageSize = ref(10)

const showForm = ref(false)
const uploading = ref(false)
const form = ref({
  id: null,
  title: '',
  summary: '',
  detail: '',
  doc_url: '',
  reward_points: 0,
})

const load = async () => {
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (statusFilter.value >= 0) params.status = statusFilter.value
    if (keyword.value.trim()) params.keyword = keyword.value.trim()

    const res = await fetchTasks(params)
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

const openCreate = () => {
  form.value = { id: null, title: '', summary: '', detail: '', doc_url: '', reward_points: 0 }
  showForm.value = true
}

const edit = (item) => {
  form.value = { ...item }
  showForm.value = true
}

const close = () => {
  showForm.value = false
}

const save = async () => {
  try {
    if (form.value.id) {
      await updateTask(form.value.id, form.value)
      notify('任务已更新', 'success')
    } else {
      await createTask(form.value)
      notify('任务已发布', 'success')
    }
    await load()
    showForm.value = false
  } catch (err) {
    notify(err.message || '保存失败', 'error')
  }
}

const off = async (item) => {
  try {
    await offTask(item.id)
    await load()
    notify('任务已下架', 'success')
  } catch (err) {
    notify(err.message || '下架失败', 'error')
  }
}

const on = async (item) => {
  try {
    await onTask(item.id)
    await load()
    notify('任务已上架', 'success')
  } catch (err) {
    notify(err.message || '上架失败', 'error')
  }
}

const handleUpload = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return
  const formData = new FormData()
  formData.append('file', file)
  try {
    uploading.value = true
    const res = await uploadImage(formData)
    form.value.doc_url = res.data?.url || ''
    notify('上传成功', 'success')
  } catch (err) {
    notify(err.message || '上传失败', 'error')
  } finally {
    uploading.value = false
  }
}

const statusText = (status) => {
  if (status === 1) return '上架'
  if (status === 2) return '下架'
  return '草稿'
}

const statusClass = (status) => {
  if (status === 1) return 'online'
  if (status === 2) return 'claimed'
  return 'pending'
}

onMounted(() => {
  load()
  // ?create=true means open create form
  if (window.location.search.includes('create=true')) {
    openCreate()
  }
})
</script>
