import { accessTokenName, POST } from '../../client'
import type { TAccessTokenResponce } from '../../types'

export async function refresh() {
  const res = await POST<TAccessTokenResponce>({ path: '/auth/refresh', protectedPath: true })
  localStorage.setItem(accessTokenName, res.accessToken)
}
