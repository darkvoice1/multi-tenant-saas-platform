<template>
  <section class="card">
    <h2>控制台</h2>
    <p>角色：<strong>{{ role || '未登录' }}</strong></p>
    <p>租户：<strong>{{ tenantId || '-' }}</strong></p>
    <div class="actions">
      <button v-if="isAdmin" class="primary">管理员动作</button>
      <button class="ghost">成员动作</button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { getRole, getTenantId } from '../auth'

const role = computed(() => getRole())
const tenantId = computed(() => getTenantId())
const isAdmin = computed(() => role.value === 'admin')
</script>

<style scoped>
.card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 1.5rem;
  background: white;
}
.actions {
  display: flex;
  gap: 0.8rem;
  margin-top: 1rem;
}
.primary {
  padding: 0.5rem 0.9rem;
  border: none;
  background: #0ea5e9;
  color: white;
  border-radius: 8px;
}
.ghost {
  padding: 0.5rem 0.9rem;
  border: 1px solid #cbd5f5;
  background: transparent;
  border-radius: 8px;
}
</style>
