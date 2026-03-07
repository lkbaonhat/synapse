<template>
  <!-- §9 — Clean template: no logic, label/error handled in script -->
  <div class="base-input-wrapper">
    <label v-if="label" :for="inputId" class="base-input-label">
      {{ label }}
      <span v-if="required" class="required-star" aria-hidden="true">*</span>
    </label>
    <input
      :id="inputId"
      v-bind="$attrs"
      :value="modelValue"
      :class="['input', { 'input--error': !!error }]"
      :aria-describedby="error ? `${inputId}-error` : undefined"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <Transition name="fade">
      <span v-if="error" :id="`${inputId}-error`" class="base-input-error" role="alert">
        {{ error }}
      </span>
    </Transition>
    <p v-if="hint && !error" class="base-input-hint">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { nanoid } from 'nanoid'

// §3 — Base prefix; §4 — Generic props/emits syntax
const props = defineProps<{
  modelValue?: string
  label?: string
  error?: string
  hint?: string
  required?: boolean
  id?: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

// Auto-generate a unique ID if not provided — §9 (no random in template)
const inputId = computed(() => props.id ?? `input-${nanoid(6)}`)
</script>

<style scoped>
/* §10 — CSS tokens only */
.base-input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.base-input-label {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--color-text-muted);
}

.required-star { color: #ef4444; margin-left: 0.2rem; }

.base-input-error {
  font-size: 0.78rem;
  color: #ef4444;
}

.base-input-hint {
  font-size: 0.78rem;
  color: var(--color-text-muted);
  margin: 0;
}

/* Transition */
.fade-enter-active, .fade-leave-active { transition: opacity 0.15s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
