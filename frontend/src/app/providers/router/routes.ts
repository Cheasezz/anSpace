import { AuthPage } from '@/pages/auth'
import type { RouteRecordRaw } from 'vue-router'

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'auth',
    component: AuthPage,
    meta: {
      requiresAuth: false,
    },
  },
]
