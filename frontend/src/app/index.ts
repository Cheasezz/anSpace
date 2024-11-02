import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { router } from './providers'
import Particles from '@tsparticles/vue3'
import { loadSlim } from '@tsparticles/slim'

export const app = createApp(App)
  .use(Particles, {
    init: async (engine) => {
      await loadSlim(engine)
    },
  })
  .use(createPinia())
  .use(router)
