<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">个人中心</div>
        <div class="panel-sub">账户与积分信息</div>
      </div>
    </div>

    <div class="profile-grid">
      <div class="profile-item">
        <div class="panel-sub">当前积分</div>
        <div class="profile-value">{{ profile.available_points || 0 }}</div>
      </div>
      <div class="profile-item">
        <div class="panel-sub">可提现金额</div>
        <div class="profile-value">{{ withdrawable }} 元</div>
      </div>
      <div class="profile-item">
        <div class="panel-sub">已完成任务</div>
        <div class="profile-value">{{ completedCount }}</div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-title">收款二维码</div>
      <div class="form-grid">
        <label class="full">
          微信收款二维码
          <input type="file" @change="handleUpload" :disabled="!!qrUrl" />
          <div class="upload-tip" v-if="uploading">上传中...</div>
          <div v-else-if="qrUrl" class="upload-preview">
            <img :src="qrUrl" alt="收款二维码预览" />
            <div class="upload-tip">已上传：{{ qrUrl }}</div>
            <button class="ghost-btn" type="button" @click="resetQr">重新上传</button>
          </div>
        </label>
      </div>
    </div>

    <div class="panel">
      <div class="panel-title">修改密码</div>
      <div class="form-grid">
        <label>
          原密码
          <input type="password" v-model="oldPassword" />
        </label>
        <label>
          新密码
          <input type="password" v-model="newPassword" />
        </label>
      </div>
      <div class="modal-actions">
        <button class="primary-btn" @click="save">保存修改</button>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { changePassword, fetchClaimHistory, uploadImage, updateWechatQr } from '../api/user'
import { profileState, loadProfile } from '../store/profile'
import { notify } from '../store/notify'

const profile = profileState
const oldPassword = ref('')
const newPassword = ref('')
const completedCount = ref(0)
const qrUrl = ref('')
const uploading = ref(false)

const load = async () => {
  try {
    await loadProfile()

    const resHistory = await fetchClaimHistory()
    const list = resHistory.data?.list || []
    completedCount.value = list.filter((item) => item.status === 3).length
    qrUrl.value = profile.value.wechat_qr_url || ''
  } catch (err) {
    console.error(err)
  }
}

const withdrawable = computed(() => {
  const points = profile.value.available_points || 0
  return (points / 10).toFixed(2)
})

const save = async () => {
  try {
    await changePassword({ old_password: oldPassword.value, new_password: newPassword.value })
    oldPassword.value = ''
    newPassword.value = ''
    await loadProfile()
    notify('密码已更新', 'success')
  } catch (err) {
    notify(err.message || '修改失败', 'error')
  }
}

const handleUpload = async (event) => {
  if (qrUrl.value) {
    notify('收款二维码已绑定，如需更换请点击“重新上传”', 'error')
    return
  }
  const file = event.target.files?.[0]
  if (!file) {
    notify('请选择图片文件', 'error')
    return
  }
  const formData = new FormData()
  formData.append('file', file)
  try {
    uploading.value = true
    const res = await uploadImage(formData)
    const url = res.data?.url
    if (url) {
      qrUrl.value = url
      await updateWechatQr({ wechat_qr_url: url })
      await loadProfile()
      notify('上传成功', 'success')
    }
  } catch (err) {
    notify(err.message || '上传失败', 'error')
  } finally {
    uploading.value = false
  }
}

const resetQr = () => {
  qrUrl.value = ''
}

onMounted(load)
</script>
