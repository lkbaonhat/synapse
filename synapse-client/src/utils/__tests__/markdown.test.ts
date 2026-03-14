/**
 * TDD: Tests for markdown utilities
 *
 * Run: pnpm test src/utils/__tests__/markdown.test.ts
 */
import { describe, it, expect } from 'vitest'
import { renderMarkdown, renderCloze } from '@/utils/markdown'

describe('renderMarkdown', () => {
  it('converts bold markdown to <strong>', () => {
    const html = renderMarkdown('**bold text**')
    expect(html).toContain('<strong>bold text</strong>')
  })

  it('converts italic markdown to <em>', () => {
    const html = renderMarkdown('*italic*')
    expect(html).toContain('<em>italic</em>')
  })

  it('converts a heading to <h1>', () => {
    const html = renderMarkdown('# Hello')
    expect(html).toContain('<h1>Hello</h1>')
  })

  it('wraps code blocks in <pre><code>', () => {
    const html = renderMarkdown('```js\nconsole.log("hi")\n```')
    expect(html).toContain('<pre')
    expect(html).toContain('<code')
  })

  it('sanitizes dangerous <script> tags (XSS prevention)', () => {
    const html = renderMarkdown('<script>alert("xss")</script>')
    expect(html).not.toContain('<script>')
    expect(html).not.toContain('alert')
  })

  it('sanitizes onclick attributes', () => {
    const html = renderMarkdown('<img onclick="evil()" src="x.png">')
    expect(html).not.toContain('onclick')
  })

  it('returns a non-empty string for plain text', () => {
    const html = renderMarkdown('Hello world')
    expect(html.length).toBeGreaterThan(0)
    expect(html).toContain('Hello world')
  })
})

describe('renderCloze', () => {
  it('replaces {{word}} with a blank span when reveal=false', () => {
    const html = renderCloze('A {{closure}} accesses outer scope.', false)
    expect(html).toContain('cloze-blank')
    expect(html).not.toContain('closure')
  })

  it('replaces {{word}} with a revealed span when reveal=true', () => {
    const html = renderCloze('A {{closure}} accesses outer scope.', true)
    expect(html).toContain('cloze-answer')
    expect(html).toContain('closure')
  })

  it('handles multiple cloze blanks', () => {
    const html = renderCloze('{{A}} is {{B}}.', false)
    const blankCount = (html.match(/cloze-blank/g) ?? []).length
    expect(blankCount).toBe(2)
  })

  it('preserves text outside brackets unchanged', () => {
    const html = renderCloze('Hello {{world}}!', false)
    expect(html).toContain('Hello')
    expect(html).toContain('!')
  })

  it('sanitizes any injected HTML inside the brackets', () => {
    const html = renderCloze('{{<script>evil()</script>}}', false)
    expect(html).not.toContain('<script>')
  })
})
