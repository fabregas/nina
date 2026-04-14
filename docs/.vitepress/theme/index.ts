import DefaultTheme from 'vitepress/theme'
import './style.css'
import Preview from './Preview.vue'

export default {
  extends: DefaultTheme,

  enhanceApp({ app }) {
    app.component('Preview', Preview)
  }
}
