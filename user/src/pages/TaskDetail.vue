<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">{{ detail.title }}</div>
        <div class="panel-sub">奖励 {{ detail.reward_points }} 积分</div>
      </div>
      <button class="primary-btn" @click="handleClaim" :disabled="detail.claimed">
        {{ detail.claimed ? '已领取' : '领取任务' }}
      </button>
    </div>

    <NoticeBar :text="statusTip" />

    <div class="detail-block">
      <div class="detail-title">提交要求</div>
      <ul class="rule-list">
        <li>提交任务成果截图（可多图）</li>
        <li>填写备注说明，描述完成过程</li>
        <li>审核通过后自动发放积分</li>
      </ul>
    </div>

    <div class="detail-block">
      <div class="detail-title">任务简介</div>
      <div class="detail-text">{{ detail.summary }}</div>
    </div>

    <div class="detail-block">
      <div class="detail-title">任务详情</div>
      <div class="detail-text">{{ detail.detail }}</div>
    </div>

    <div v-if="detail.doc_url" class="detail-block">
      <div class="detail-title">文档资料</div>
      <a :href="detail.doc_url" target="_blank">{{ detail.doc_url }}</a>
    </div>

    <div class="detail-block">
      <div class="detail-title">领取规则</div>
      <ul class="rule-list">
        <li>每个任务只能被一个用户领取</li>
        <li>每个用户同一时间最多领取 1 个任务</li>
        <li>提交成果后等待管理员审核</li>
      </ul>
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
      <span class="badge" :class="detail.claimed ? 'claimed' : 'available'">
        {{ detail.claimed ? '已领取' : '可领取' }}
      </span>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchTaskDetail, claimTask } from '../api/user'
import NoticeBar from '../components/NoticeBar.vue'

const route = useRoute()
const router = useRouter()
const detail = ref({})

const statusTip = computed(() => {
  if (detail.value.claimed) {
    return '当前任务已被领取，如已领取请前往“我的任务”提交成果。'
  }
  return '任务可领取，领取后请在“我的任务”提交截图与备注。'
})

const load = async () => {
  try {
    const res = await fetchTaskDetail(route.params.id)
    detail.value = res.data || {}
  } catch (err) {
    console.error(err)
  }
}

const handleClaim = async () => {
  try {
    await claimTask(detail.value.id)
    await load()
    router.push('/my-tasks')
  } catch (err) {
    console.error(err)
  }
}

onMounted(load)
</script>
