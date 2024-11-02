import { accessTokenName, POST } from '../client'
import { type TAccessTokenResponce, type TUserAuth } from '../types'

export async function signup(userAuth: TUserAuth) {
  const res = await POST<TAccessTokenResponce>('/auth/signup', userAuth)

  localStorage.setItem(accessTokenName, res.accessToken)
}
