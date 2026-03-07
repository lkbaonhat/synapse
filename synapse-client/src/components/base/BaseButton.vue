<template>
  <!-- §9 — No logic in template; variant/size/loading resolved in script -->
  <button
    :class="['btn', variantClass, sizeClass, { 'btn--loading': loading }]"
    :disabled="disabled || loading"
    :type="type"
    v-bind="$attrs"
  >
    <span v-if="loading" class="btn-spinner" aria-hidden="true"></span>
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// §3 — Base prefix; §4 — Generic props syntax
const props = defineProps<{
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
  disabled?: boolean
  type?: 'button' | 'submit' | 'reset'
}>()

// §9 — Complex resolved in computed, not inline in template
const variantClass = computed(() => {
  const map: Record<string, string> = {
    primary: 'btn-primary',
    secondary: 'btn-secondary',
    danger: 'btn-danger',
    ghost: 'btn-ghost'
  }
  return map[props.variant ?? 'primary']
})

const sizeClass = computed(() => {
  const map: Record<string, string> = { sm: 'btn-sm', md: '', lg: 'btn-lg' }
  return map[props.size ?? 'md']
})
</script>

<style scoped>
/* §10 — All colours via CSS tokens */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.45rem;
  font-weight: 500;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
  font-size: 0.875rem;
  padding: 0.55rem 1.25rem;
}

.btn:disabled { opacity: 0.5; cursor: not-allowed; }

/* Variants */
.btn-primary  { background: var(--color-primary); color: white; }
.btn-primary:hover:not(:disabled) { background: var(--color-primary-hover); }

.btn-secondary { background: var(--color-surface-hover); color: var(--color-text); border: 1px solid var(--color-border); }
.btn-secondary:hover:not(:disabled) { border-color: var(--color-primary); color: var(--color-primary); }

.btn-danger { background: rgba(239,68,68,0.1); color: #ef4444; border: 1px solid rgba(239,68,68,0.2); }
.btn-danger:hover:not(:disabled) { background: rgba(239,68,68,0.2); }

.btn-ghost { background: none; color: var(--color-text-muted); }
.btn-ghost:hover:not(:disabled) { background: var(--color-surface-hover); color: var(--color-text); }

/* Sizes */
.btn-sm { padding: 0.35rem 0.85rem; font-size: 0.8rem; }
.btn-lg { padding: 0.75rem 1.75rem; font-size: 1rem; }

/* Loading spinner */
.btn--loading { pointer-events: none; }
.btn-spinner {
  width: 13px;
  height: 13px;
  border: 2px solid rgba(255,255,255,0.35);
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;
}

@keyframes spin { to { transform: rotate(360deg); } }
</style>
