// commitlint.config.ts
// Enforces Conventional Commits format on every commit message.
// Runs automatically via Husky's commit-msg hook (see .husky/commit-msg).
import type { UserConfig } from '@commitlint/types'

const config: UserConfig = {
  extends: ['@commitlint/config-conventional'],

  // ── Custom rules on top of conventional ──────────────────────────────
  rules: {
    // Type must be one of the values below — no others allowed
    'type-enum': [
      2, // error level (2 = error, 1 = warn, 0 = off)
      'always',
      [
        'feat',     // ✨ New feature
        'fix',      // 🐛 Bug fix
        'docs',     // 📝 Documentation only
        'style',    // 💅 Formatting/whitespace — no logic change
        'refactor', // ♻️  Code restructure — no feature/fix
        'perf',     // ⚡ Performance improvement
        'test',     // ✅ Add/fix tests
        'build',    // 🏗  Build system / dependency changes
        'ci',       // 🤖 CI/CD config
        'chore',    // 🔧 Tooling, config, scripts
        'revert'    // ⏪ Revert a commit
      ]
    ],

    // Scope must be lowercase
    'scope-case': [2, 'always', 'lower-case'],

    // Subject (description) must not start with a capital letter
    'subject-case': [2, 'never', ['sentence-case', 'start-case', 'pascal-case', 'upper-case']],

    // Subject must not end with a period
    'subject-full-stop': [2, 'never', '.'],

    // Subject must be at least 10 characters
    'subject-min-length': [2, 'always', 10],

    // Header (type + scope + subject) max 100 chars
    'header-max-length': [2, 'always', 100],

    // Body lines max 120 chars
    'body-max-line-length': [2, 'always', 120],

    // Footer lines max 120 chars
    'footer-max-line-length': [2, 'always', 120]
  },

  // ── Prompt config (used by commitizen / pnpm commit) ─────────────────
  prompt: {
    questions: {
      type: {
        description: 'Select the type of change you are committing'
      },
      scope: {
        description: 'What is the scope of this change? (e.g. auth, deck, card, ui)'
      },
      subject: {
        description: 'Write a SHORT, IMPERATIVE description (no capital, no period at end)'
      },
      body: {
        description: 'Provide a LONGER description of the change (optional, press enter to skip)'
      },
      isBreaking: {
        description: 'Are there any BREAKING CHANGES?'
      },
      breakingBody: {
        description: 'Describe the breaking change'
      },
      breaking: {
        description: 'Describe the breaking changes (BREAKING CHANGE: ...)'
      },
      isIssueAffected: {
        description: 'Does this change affect any open issues?'
      },
      issues: {
        description: 'Add issue references (e.g. "closes #123", "refs #456")'
      }
    }
  }
}

export default config
