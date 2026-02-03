<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">我的任务</div>
        <div class="panel-sub">查看当前任务与历史提交</div>
      </div>
    </div>

    <div v-if="current" class="task-current">
      <div>
        <div class="task-title">{{ current.task_title }}</div>
        <div class="task-desc">{{ current.summary }}</div>
      </div>
      <button class="primary-btn" @click="showSubmit = true">提交成果</button>
    </div>
    <EmptyState v-else title="暂无进行中的任务" desc="去任务列表领取一个任务开始吧" />

    <div class="history">
      <div class="panel-sub">历史记录</div>
      <div class="history-item" v-for="item in history" :key="item.id">
        <div>
          <div class="task-title">{{ item.task_title }}</div>
          <div class="task-desc">状态：{{ statusText(item.status) }}</div>
          <div v-if="item.status === 4 && item.reject_reason" class="task-desc">
            驳回原因：{{ item.reject_reason }}
          </div>
        </div>
        <div class="points">{{ item.reward_points }} 积分</div>
      </div>
      <EmptyState v-if="history.length === 0" title="暂无历史任务" desc="完成任务后会在这里显示记录" />
    </div>

    <div v-if="showSubmit" class="modal">
      <div class="modal-card">
        <div class="modal-title">提交任务成果</div>
        <label>
          截图（可多图）
          <input type="file" multiple @change="handleUpload" />
        </label>
        <div class="upload-tip" v-if="uploading">上传中...</div>
        <div v-else class="upload-list">
          <img v-for="(img, idx) in images" :key="idx" :src="img" alt="截图预览" />
        </div>
        <label>
          备注说明
          <textarea v-model="remark" rows="4" />
        </label>
        <div class="modal-actions">
          <button class="ghost-btn" @click="showSubmit = false">取消</button>
          <button class="primary-btn" @click="submit">提交</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { fetchCurrentClaim, fetchClaimHistory, submitClaim, uploadImage } from '../api/user'
import EmptyState from '../components/EmptyState.vue'
import { notify } from '../store/notify'

const current = ref(null)
const history = ref([])
const showSubmit = ref(false)
const images = ref([])
const remark = ref('')
const uploading = ref(false)

const load = async () => {
  try {
    const resCurrent = await fetchCurrentClaim()
    current.value = resCurrent.data?.current || null

    const resHistory = await fetchClaimHistory()
    history.value = resHistory.data?.list || []
  } catch (err) {
    console.error(err)
  }
}

const handleUpload = async (event) => {
  const files = Array.from(event.target.files || [])
  if (files.length === 0) {
    notify('请选择图片文件', 'error')
    return
  }
  uploading.value = true
  try {
    const results = []
    for (const file of files) {
      const formData = new FormData()
      formData.append('file', file)
      const res = await uploadImage(formData)
      if (res.data?.url) results.push(res.data.url)
    }
    images.value = results
    if (results.length) notify('上传成功', 'success')
  } catch (err) {
    notify(err.message || '上传失败', 'error')
  } finally {
    uploading.value = false
  }
}

const submit = async () => {
  if (!current.value) return
  if (images.value.length === 0) {
    notify('请先上传截图', 'error')
    return
  }
  try {
    await submitClaim(current.value.id, { images: images.value, remark: remark.value })
    showSubmit.value = false
    images.value = []
    remark.value = ''
    await load()
    notify('提交成功，等待审核', 'success')
  } catch (err) {
    notify(err.message || '提交失败', 'error')
  }
}

const statusText = (status) => {
  if (status === 1) return '进行中'
  if (status === 2) return '待审核'
  if (status === 3) return '已完成'
  if (status === 4) return '驳回'
  return '-'
}

onMounted(load)
</script>
