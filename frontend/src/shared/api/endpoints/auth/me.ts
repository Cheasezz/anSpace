import { GET } from '../../client'
import type { TUserResponce } from '../../types'

export async function me(): Promise<TUserResponce> {
  try {
    const res = await GET<TUserResponce>({ path: '/auth/me', protectedPath: true })
    return res
  } catch (err) {
    throw err
  }
}
