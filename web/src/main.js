// Copyright © 2023 Ronnie Zhang (大脸怪). MIT License.

import { createApp } from 'vue'
import App from './App.vue'
import { setupDirectives } from './directives'

import { setupRouter } from './router'
import { setupStore } from './store'
import { setupNaiveDiscreteApi } from './utils'
import '@/styles/reset.css'
import '@/styles/global.css'
import 'uno.css'

async function bootstrap() {
  const app = createApp(App)
  app.config.warnHandler = (msg, instance, trace) => {
    if (msg.includes('Extraneous non-emits event listeners')) {
      console.group('⚠️ Vue Warning 捕获')
      console.log('组件名:', instance?.type?.name)
      console.log('props:', instance?.props)
      console.log('trace:', trace)
      console.groupEnd()
    }
  }
  setupStore(app)
  setupDirectives(app)
  await setupRouter(app)
  app.mount('#app')
  setupNaiveDiscreteApi()
}

bootstrap()
