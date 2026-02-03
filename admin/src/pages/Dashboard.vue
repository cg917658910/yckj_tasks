<template>
  <section>
    <section class="hero">
      <div>
        <div class="hero-title">今日概览</div>
        <div class="hero-sub">任务发布、审核、积分与提现的全局掌控</div>
      </div>
      <div class="hero-meta">
        <div class="meta-item">
          <div class="meta-label">待审核任务</div>
          <div class="meta-value">{{ stats.pendingClaims }}</div>
        </div>
        <div class="meta-item">
          <div class="meta-label">待处理提现</div>
          <div class="meta-value">{{ stats.pendingWithdrawals }}</div>
        </div>
        <div class="meta-item">
          <div class="meta-label">在线任务</div>
          <div class="meta-value">{{ stats.onlineTasks }}</div>
        </div>
      </div>
    </section>

    <section class="cards">
      <div class="card">
        <div class="card-title">本周发布任务</div>
        <div class="card-value">{{ stats.weekTasks }}</div>
        <div class="card-foot">环比 +12%</div>
      </div>
      <div class="card">
        <div class="card-title">审核通过率</div>
        <div class="card-value">82%</div>
        <div class="card-foot">驳回率 18%</div>
      </div>
      <div class="card">
        <div class="card-title">积分发放</div>
        <div class="card-value">18,240</div>
        <div class="card-foot">本周累计</div>
      </div>
      <div class="card">
        <div class="card-title">提现打款</div>
        <div class="card-value">3,540 元</div>
        <div class="card-foot">本周累计</div>
      </div>
    </section>

    <section class="grid">
      <div class="panel">
        <div class="panel-head">
          <div>
            <div class="panel-title">任务列表</div>
            <div class="panel-sub">可领取任务与状态总览</div>
          </div>
          <div class="panel-actions">
            <select>
              <option>所有状态</option>
              <option>上架</option>
              <option>下架</option>
            </select>
            <button class="ghost-btn">批量下架</button>
          </div>
        </div>
        <div class="table">
          <div class="table-head">
            <div>任务名称</div>
            <div>奖励积分</div>
            <div>领取状态</div>
            <div>发布时间</div>
            <div>操作</div>
          </div>
          <div class="table-row" v-for="item in tasks" :key="item.id">
            <div>
              <div class="task-title">{{ item.title }}</div>
              <div class="task-desc">{{ item.summary }}</div>
            </div>
            <div class="points">{{ item.points }} 积分</div>
            <div>
              <span class="chip" :class="item.statusClass">{{ item.status }}</span>
            </div>
            <div class="time">{{ item.time }}</div>
            <div>
              <button class="link">编辑</button>
              <button class="link danger">下架</button>
            </div>
          </div>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head">
          <div>
            <div class="panel-title">提现审核</div>
            <div class="panel-sub">最新申请与风控提示</div>
          </div>
          <button class="ghost-btn">查看全部</button>
        </div>
        <div class="list">
          <div class="list-item" v-for="item in withdrawals" :key="item.id">
            <div>
              <div class="list-title">{{ item.user }}</div>
              <div class="list-sub">申请 {{ item.amount }} 元 · {{ item.points }} 积分</div>
            </div>
            <div class="list-right">
              <span class="chip pending">待审核</span>
              <button class="link">处理</button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="panel approve">
      <div class="panel-head">
        <div>
          <div class="panel-title">任务审核队列</div>
          <div class="panel-sub">待审核成果与驳回记录</div>
        </div>
        <button class="ghost-btn">进入审核中心</button>
      </div>
      <div class="approve-grid">
        <div class="approve-card" v-for="item in approvals" :key="item.id">
          <div class="approve-title">{{ item.title }}</div>
          <div class="approve-user">提交人：{{ item.user }}</div>
          <div class="approve-footer">
            <span class="chip pending">待审核</span>
            <div class="approve-actions">
              <button class="link">查看</button>
              <button class="link success">通过</button>
              <button class="link danger">驳回</button>
            </div>
          </div>
        </div>
      </div>
    </section>
  </section>
</template>

<script setup>
const stats = {
  pendingClaims: 12,
  pendingWithdrawals: 7,
  onlineTasks: 34,
  weekTasks: 48,
}

const tasks = [
  {
    id: 1,
    title: '市场调研任务',
    summary: '收集用户反馈并输出调研报告',
    points: 200,
    status: '可领取',
    statusClass: 'online',
    time: '2026-02-03 14:20',
  },
  {
    id: 2,
    title: '文档翻译任务',
    summary: '将英文说明翻译为中文',
    points: 100,
    status: '已领取',
    statusClass: 'claimed',
    time: '2026-02-03 11:10',
  },
  {
    id: 3,
    title: '数据校验任务',
    summary: '核验提交数据并输出清单',
    points: 150,
    status: '待审核',
    statusClass: 'pending',
    time: '2026-02-02 19:55',
  },
]

const withdrawals = [
  { id: 1, user: '小北', amount: 120, points: 1200 },
  { id: 2, user: '阿澈', amount: 80, points: 800 },
  { id: 3, user: '林夕', amount: 50, points: 500 },
]

const approvals = [
  { id: 1, title: '门店探访任务', user: '小北' },
  { id: 2, title: '问卷整理任务', user: '阿澈' },
  { id: 3, title: '竞品截图任务', user: '林夕' },
]
</script>
