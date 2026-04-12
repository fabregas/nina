import { defineConfig } from 'vitepress'

export default defineConfig({
  base: "/nina/",
  title: "Nina UI",
  description: "WebAssembly UI Framework for Go",
  ignoreDeadLinks: true,
  themeConfig: {
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Documentation', link: '/guide/getting-started' }
    ],

    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'Getting started', link: '/guide/getting-started' },
          { text: 'Architecture', link: '/guide/architecture' }
        ]
      },
      {
        text: 'UI elements',
        items: [
          { text: 'Button', link: '/ui_element/button' },
          { text: 'Alert', link: '/ui_element/alert' },
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/fabregas/nina' }
    ]
  }
})
