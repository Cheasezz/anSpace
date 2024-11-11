import { createRouter, createWebHistory } from 'vue-router'
import { routes } from './routes'
import { accessTokenName } from '@/shared/api/client'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem(accessTokenName)
  if (to.meta.requiresAuth && !token) next('/')
  else next()
})

export { router }
