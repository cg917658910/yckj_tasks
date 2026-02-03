<template>
  <aside class="sidebar">
    <div class="brand">
      <div class="logo">T</div>
      <div>
        <div class="brand-title">任务管理系统</div>
        <div class="brand-sub">任务发布平台</div>
      </div>
    </div>

    <div class="profile">
      <div class="avatar">管</div>
      <div>
        <div class="name">admin</div>
        <div class="role">管理员</div>
      </div>
    </div>

    <nav class="menu">
      <RouterLink class="menu-item" to="/dashboard">仪表盘</RouterLink>
      <RouterLink class="menu-item" to="/tasks">任务管理</RouterLink>
      <RouterLink class="menu-item" to="/claims">审核中心</RouterLink>
      <RouterLink class="menu-item" to="/points">积分管理</RouterLink>
      <RouterLink class="menu-item" to="/points/history">积分历史</RouterLink>
      <RouterLink class="menu-item" to="/withdrawals">提现管理</RouterLink>
      <RouterLink class="menu-item" to="/users">用户管理</RouterLink>
      <RouterLink class="menu-item" to="/settings">系统设置</RouterLink>
    </nav>

    <div class="sidebar-footer">
      <div class="rule-card">
        <div class="rule-title">积分规则</div>
        <div class="rule-row">10 积分 = 1 元</div>
        <div class="rule-row">最低提现 10 元</div>
        <div class="rule-row">注册赠送 10 积分</div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { RouterLink } from 'vue-router'
</script>
