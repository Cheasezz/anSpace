import { createApiFetchClient, type TUnserializedBody } from 'feature-fetch'
import type { TErrorResponce } from './types'
import { errorCheck } from './endpoints/errors'

export const accessTokenName = 'accessToken'

const fetchClient = createApiFetchClient({
  prefixUrl: `${import.meta.env.VITE_BACKEND_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
  fetchProps: { credentials: 'include' },
})

export async function POST<T>(
  path: string,
  body: TUnserializedBody,
  protectedPath: boolean = false,
): Promise<T> {
  if (protectedPath) {
    const res = await fetchClient.post<unknown, TErrorResponce, TUnserializedBody>(path, body, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem(accessTokenName)}`,
      },
    })

    errorCheck(res)

    return res.unwrap().data as T
  } else {
    const res = await fetchClient.post<unknown, TErrorResponce, TUnserializedBody>(path, body)
    errorCheck(res)

    return res.unwrap().data as T
  }
}
