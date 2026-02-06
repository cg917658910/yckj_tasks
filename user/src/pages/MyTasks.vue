<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">我的任务</div>
        <div class="panel-sub">查看当前任务与历史提交</div>
      </div>
    </div>

    <div v-if="current && current?.status" class="panel">
      <div class="panel-head">
        <div>
          <div class="panel-title">{{ current.task_title }}</div>
          <div class="panel-sub">奖励 {{ current.reward_points }} 积分</div>
        </div>
        <button
          v-if="current.status === 1 || current.status === 4"
          class="primary-btn"
          @click="showSubmit = true"
        >
          提交成果
        </button>
        <button v-else class="primary-btn disabled" disabled>提交成果</button>
      </div>
      <div class="detail-block">
        <div class="detail-title">任务简介</div>
        <div class="detail-text">{{ current.summary }}</div>
      </div>

      <div class="detail-block">
        <div class="detail-title">任务详情</div>
        <div class="detail-text">{{ current.detail }}</div>
      </div>

      <div v-if="current.doc_url" class="detail-block">
        <div class="detail-title">文档资料</div>
        <a :href="current.doc_url" target="_blank">查看附件</a>
      </div>

      <div class="detail-block">
        <div class="detail-title">状态流转</div>
        <ol class="flow-list">
          <li>领取任务</li>
          <li>提交成果</li>
          <li>等待审核</li>
          <li>审核通过获得积分</li>
        </ol>
      </div>

      <div class="detail-block">
        <div class="detail-title">当前状态</div>
        <span class="badge" :class="statusBadgeClass(current.status)">{{ statusText(current.status) }}</span>
      </div>

      <div v-if="current.status === 2" class="detail-block">
        <div class="detail-title">审核结果</div>
        <div class="detail-text">审核中，请耐心等待</div>
      </div>
    </div>
    <div v-else class="panel empty-current">
      <EmptyState title="暂无进行中的任务" desc="去任务列表领取一个任务开始吧" />
      <div class="modal-actions">
        <button class="primary-btn" @click="goTasks">去领取任务</button>
      </div>
    </div>

    <div class="history">
      <div class="panel-sub">历史记录</div>
      <div class="history-item" v-for="item in history" :key="item.id">
        <div>
          <div class="task-title">{{ item.task_title }}</div>
          <div class="task-desc">状态：{{ statusText(item.status) }}</div>
          <div v-if="item.status === 4 && item.reject_reason" class="task-desc">
            驳回原因：{{ item.reject_reason }}
          </div>
          <div v-if="item.doc_url" class="task-desc">
            附件：
            <a :href="item.doc_url" target="_blank">查看附件</a>
          </div>
        </div>
        <div class="points">+{{ item.reward_points_final }} 积分</div>
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
import { useRouter } from 'vue-router'
import { fetchClaimHistory, fetchCurrentClaim, submitClaim, uploadImage } from '../api/user'
import EmptyState from '../components/EmptyState.vue'
import { notify } from '../store/notify'

const current = ref(null)
const history = ref([])
const showSubmit = ref(false)
const images = ref([])
const remark = ref('')
const uploading = ref(false)
const router = useRouter()

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
  if (status === 1) return '待提交'
  if (status === 2) return '待审核'
  if (status === 3) return '已完成'
  if (status === 4) return '驳回'
  return '-'
}

const statusBadgeClass = (status) => {
  if (status === 1) return 'available'
  if (status === 2) return 'pending'
  if (status === 3) return 'claimed'
  if (status === 4) return 'claimed'
  return 'available'
}

const goTasks = () => {
  router.push('/tasks')
}

onMounted(load)
</script>
