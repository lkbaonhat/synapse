import js from '@eslint/js'
import pluginVue from 'eslint-plugin-vue'
import tsParser from '@typescript-eslint/parser'

export default [
  js.configs.recommended,
  ...pluginVue.configs['flat/recommended'],
  {
    files: ['**/*.{ts,vue}'],
    languageOptions: {
      parser: pluginVue.processor ? undefined : tsParser,
      parserOptions: {
        parser: tsParser,
        sourceType: 'module'
      }
    },
    rules: {
      // Vue-specific
      'vue/multi-word-component-names': 'off',
      'vue/no-v-html': 'warn',
      // TypeScript-friendly
      'no-unused-vars': 'off',
      'no-console': ['warn', { allow: ['warn', 'error'] }]
    }
  },
  {
    ignores: ['dist/**', 'node_modules/**', '*.config.*']
  }
]
