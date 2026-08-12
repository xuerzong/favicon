import { debounce } from './debounce'
import './style.css'

const API_BASE_URL = 'http://localhost:8080/'

const input = document.getElementById('url-input') as HTMLInputElement
const wrapper = document.getElementById('input-wrapper') as HTMLLabelElement
const prefix = document.querySelector('.input-prefix') as HTMLSpanElement
const preview = document.getElementById('favicon-preview') as HTMLImageElement
const loading = document.getElementById('loading') as HTMLDivElement

prefix.textContent = API_BASE_URL

let debouncedUrl = 'example.com'

const handleUrlChange = debounce((nextUrl: string) => {
  debouncedUrl = nextUrl
  preview.src = `${API_BASE_URL}${debouncedUrl || 'example.com'}`
  loading.hidden = true
}, 1000)

input.addEventListener('focus', () => {
  wrapper.dataset.focused = 'true'
})

input.addEventListener('blur', () => {
  delete wrapper.dataset.focused
})

input.addEventListener('input', () => {
  loading.hidden = false
  handleUrlChange(input.value)
})
