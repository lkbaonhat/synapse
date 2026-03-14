// §2 + §9 — Pure util function: no Vue imports, no side effects, fully unit-testable
// Used via computed() in views/components — never called inline in templates
import { marked } from 'marked'
import hljs from 'highlight.js'
import DOMPurify from 'dompurify'

marked.setOptions({
  breaks: true,
  gfm: true
})

// Custom renderer to add syntax highlighting to code blocks
const renderer = new marked.Renderer()
renderer.code = function(codeOrToken: string | { text: string; lang?: string }, language?: string | undefined): string {
  const text = typeof codeOrToken === 'string' ? codeOrToken : codeOrToken.text
  const tokenLang = typeof codeOrToken === 'string' ? language : codeOrToken.lang
  const lang = tokenLang && hljs.getLanguage(tokenLang) ? tokenLang : 'plaintext'
  const highlighted = hljs.highlight(text, { language: lang }).value
  return `<pre class="hljs"><code class="language-${lang}">${highlighted}</code></pre>`
}

marked.use({ renderer })

/**
 * Converts a markdown string to sanitised HTML.
 * @param text - Raw markdown text
 * @returns XSS-safe HTML string, safe to bind with v-html
 */
export function renderMarkdown(text: string): string {
  const rawHtml = marked.parse(text) as string
  return DOMPurify.sanitize(rawHtml, {
    ALLOWED_TAGS: [
      'p', 'br', 'strong', 'em', 'u', 's', 'code', 'pre',
      'blockquote', 'ul', 'ol', 'li', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
      'a', 'img', 'table', 'thead', 'tbody', 'tr', 'th', 'td', 'span'
    ],
    ALLOWED_ATTR: ['href', 'src', 'alt', 'class', 'title', 'target']
  })
}

/**
 * Renders a cloze-style sentence by replacing {{word}} placeholders
 * with a visible blank span — used for the study card front face.
 */
export function renderCloze(text: string, reveal = false): string {
  const replaced = text.replace(/\{\{(.+?)\}\}/g, (_match, word: string) => {
    if (reveal) {
      return `<span class="cloze-answer">${word}</span>`
    }
    return `<span class="cloze-blank" aria-label="hidden word">_____</span>`
  })
  return DOMPurify.sanitize(replaced, {
    ALLOWED_TAGS: ['span'],
    ALLOWED_ATTR: ['class', 'aria-label']
  })
}
