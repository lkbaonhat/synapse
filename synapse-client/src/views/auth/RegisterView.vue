<template>
  <div class="auth-form">
    <div class="form-header">
      <h2 class="form-title">Create your account</h2>
      <p class="form-subtitle">Start your learning journey today — it's free.</p>
    </div>

    <form @submit.prevent="onSubmit" class="form" novalidate>
      <!-- Global error -->
      <Transition name="fade">
        <div v-if="authStore.error" class="alert alert-error">
          <span>⚠</span> {{ authStore.error }}
        </div>
      </Transition>

      <!-- Name -->
      <div class="form-group">
        <label for="reg-name" class="form-label">Full Name</label>
        <input
          id="reg-name"
          v-model="name"
          type="text"
          class="input"
          :class="{ 'input--error': errors.name }"
          placeholder="Jane Doe"
          autocomplete="name"
        />
        <span v-if="errors.name" class="field-error">{{ errors.name }}</span>
      </div>

      <!-- Email -->
      <div class="form-group">
        <label for="reg-email" class="form-label">Email</label>
        <input
          id="reg-email"
          v-model="email"
          type="email"
          class="input"
          :class="{ 'input--error': errors.email }"
          placeholder="you@example.com"
          autocomplete="email"
        />
        <span v-if="errors.email" class="field-error">{{ errors.email }}</span>
      </div>

      <!-- Password -->
      <div class="form-group">
        <label for="reg-password" class="form-label">Password</label>
        <div class="input-wrapper">
          <input
            id="reg-password"
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            class="input"
            :class="{ 'input--error': errors.password }"
            placeholder="Min. 6 characters"
            autocomplete="new-password"
          />
          <button type="button" class="toggle-pwd" @click="showPassword = !showPassword" tabindex="-1">
            {{ showPassword ? '🙈' : '👁' }}
          </button>
        </div>
        <span v-if="errors.password" class="field-error">{{ errors.password }}</span>
      </div>

      <!-- Confirm Password -->
      <div class="form-group">
        <label for="reg-confirm" class="form-label">Confirm Password</label>
        <input
          id="reg-confirm"
          v-model="confirmPassword"
          :type="showPassword ? 'text' : 'password'"
          class="input"
          :class="{ 'input--error': errors.confirmPassword }"
          placeholder="Repeat your password"
          autocomplete="new-password"
        />
        <span v-if="errors.confirmPassword" class="field-error">{{ errors.confirmPassword }}</span>
      </div>

      <button type="submit" class="btn btn-primary btn-full" :disabled="authStore.isLoading">
        <span v-if="authStore.isLoading" class="spinner"></span>
        <span>{{ authStore.isLoading ? 'Creating account…' : 'Create Account' }}</span>
      </button>
    </form>

    <p class="form-footer">
      Already have an account?
      <RouterLink to="/login" class="link">Sign in</RouterLink>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useForm, useField } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'

import { useAuthStore } from '@/store/useAuthStore'
import { registerSchema } from '@/types'

const router = useRouter()
const authStore = useAuthStore()

const { handleSubmit, errors } = useForm({
  validationSchema: toTypedSchema(registerSchema)
})

const { value: name } = useField<string>('name')
const { value: email } = useField<string>('email')
const { value: password } = useField<string>('password')
const { value: confirmPassword } = useField<string>('confirmPassword')
const showPassword = ref(false)

const onSubmit = handleSubmit(async (values) => {
  try {
    await authStore.register(values)
    router.push('/app/library')
  } catch {
    // error is already set in the store
  }
})
</script>

<style scoped>
.auth-form { display: flex; flex-direction: column; gap: 1.25rem; width: 100%; }
.form-title { font-size: 1.5rem; font-weight: 700; color: var(--color-text); margin: 0 0 0.4rem; }
.form-subtitle { font-size: 0.875rem; color: var(--color-text-muted); margin: 0; }
.form { display: flex; flex-direction: column; gap: 0.85rem; }
.form-group { display: flex; flex-direction: column; gap: 0.35rem; }
.form-label { font-size: 0.85rem; font-weight: 500; color: var(--color-text-muted); }
.input-wrapper { position: relative; }
.input-wrapper .input { padding-right: 2.5rem; }
.toggle-pwd { position: absolute; right: 0.75rem; top: 50%; transform: translateY(-50%); background: none; border: none; cursor: pointer; font-size: 1rem; padding: 0; line-height: 1; }
.input--error { border-color: #ef4444 !important; }
.field-error { font-size: 0.78rem; color: #ef4444; }
.alert { padding: 0.75rem 1rem; border-radius: var(--radius-sm); font-size: 0.875rem; display: flex; align-items: center; gap: 0.5rem; }
.alert-error { background: rgba(239,68,68,0.1); color: #ef4444; border: 1px solid rgba(239,68,68,0.3); }
.btn-full { width: 100%; margin-top: 0.25rem; }
.spinner { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 0.7s linear infinite; flex-shrink: 0; }
@keyframes spin { to { transform: rotate(360deg); } }
.form-footer { text-align: center; font-size: 0.875rem; color: var(--color-text-muted); margin: 0; }
.link { color: var(--color-primary); text-decoration: none; font-weight: 500; }
.link:hover { text-decoration: underline; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s, transform 0.2s; }
.fade-enter-from, .fade-leave-to { opacity: 0; transform: translateY(-4px); }
</style>
