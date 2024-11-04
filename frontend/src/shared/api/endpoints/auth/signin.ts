import { accessTokenName, POST } from '../../client'
import type { TAccessTokenResponce, TUserAuth } from '../../types'

export async function signin(userAuth: TUserAuth) {
  const res = await POST<TAccessTokenResponce>({ path: '/auth/signin', body: userAuth })
  localStorage.setItem(accessTokenName, res.accessToken)
}
