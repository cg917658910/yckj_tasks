<template>
  <section class="panel">
    <div class="panel-head">
      <div>
        <div class="panel-title">积分规则</div>
        <div class="panel-sub">配置兑换比例与提现门槛</div>
      </div>
    </div>

    <div class="form-grid">
      <label>
        兑换比例（积分/元）
        <input v-model.number="form.exchange_rate" type="number" />
      </label>
      <label>
        最低提现金额
        <input v-model.number="form.min_withdraw_amount" type="number" />
      </label>
      <label>
        注册赠送积分
        <input v-model.number="form.register_bonus_points" type="number" />
      </label>
    </div>

    <div class="modal-actions">
      <button class="primary-btn" @click="save">保存规则</button>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { fetchPointsRule, updatePointsRule } from '../api/admin'
import { notify } from '../store/notify'

const form = ref({
  exchange_rate: 10,
  min_withdraw_amount: 10,
  register_bonus_points: 10,
})

const load = async () => {
  try {
    const res = await fetchPointsRule()
    form.value = { ...form.value, ...(res.data || {}) }
  } catch (err) {
    console.error(err)
  }
}

const save = async () => {
  try {
    await updatePointsRule(form.value)
    notify('规则已保存', 'success')
  } catch (err) {
    notify(err.message || '保存失败', 'error')
  }
}

onMounted(load)
</script>
