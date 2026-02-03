<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">提现中心</div>
        <div class="panel-sub">上传微信收款码并申请提现</div>
      </div>
    </div>

    <NoticeBar
      v-if="!qrUrl"
      text="请先在个人中心上传微信收款二维码，上传后即可申请提现。"
    />
    <button v-if="!qrUrl" class="ghost-btn" type="button" @click="goProfile">去个人中心上传</button>
    <div v-else class="upload-preview">
      <img :src="qrUrl" alt="收款二维码预览" />
      <div class="upload-tip">已绑定收款二维码</div>
    </div>

    <div class="form-grid">
      <label>
        提现金额
        <input type="number" v-model.number="amount" placeholder="最低提现 10 元" />
      </label>
    </div>

    <div class="modal-actions">
      <button class="primary-btn" @click="apply">申请提现</button>
    </div>

    <div class="panel">
      <div class="panel-title">提现记录</div>
      <div class="history-item" v-for="item in records" :key="item.id">
        <div>
          <div class="task-title">{{ item.amount }} 元</div>
          <div class="task-desc">状态：{{ statusText(item.status) }}</div>
        </div>
        <div class="points">{{ item.points_cost }} 积分</div>
      </div>
      <EmptyState v-if="records.length === 0" title="暂无提现记录" desc="申请提现后会显示处理进度" />
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { applyWithdrawal, fetchWithdrawals, fetchProfile } from '../api/user'
import EmptyState from '../components/EmptyState.vue'
import NoticeBar from '../components/NoticeBar.vue'
import { notify } from '../store/notify'

const router = useRouter()
const qrUrl = ref('')
const amount = ref(0)
const records = ref([])

const load = async () => {
  try {
    const profile = await fetchProfile()
    qrUrl.value = profile.data?.wechat_qr_url || ''
    const res = await fetchWithdrawals()
    records.value = res.data?.list || []
  } catch (err) {
    console.error(err)
  }
}

const goProfile = () => {
  router.push('/profile')
}

const apply = async () => {
  if (!qrUrl.value) {
    notify('请先在个人中心上传收款二维码', 'error')
    return
  }
  if (!amount.value || amount.value <= 0) {
    notify('请输入正确的提现金额', 'error')
    return
  }
  try {
    await applyWithdrawal({ amount: amount.value })
    amount.value = 0
    await load()
    notify('提现申请已提交', 'success')
  } catch (err) {
    notify(err.message || '提现失败', 'error')
  }
}

const statusText = (status) => {
  if (status === 2) return '已打款'
  if (status === 3) return '驳回'
  return '待审核'
}

onMounted(load)
</script>
