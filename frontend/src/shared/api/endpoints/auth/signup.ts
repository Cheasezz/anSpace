import { accessTokenName, POST } from '../../client'
import type { TAccessTokenResponce, TUserAuth } from '../../types'

export async function signup(userAuth: TUserAuth) {
  const res = await POST<TAccessTokenResponce>({ path: '/auth/signup', body: userAuth })
  localStorage.setItem(accessTokenName, res.accessToken)
}
