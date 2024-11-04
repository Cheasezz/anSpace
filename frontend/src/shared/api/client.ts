import { createApiFetchClient, type TUnserializedBody } from 'feature-fetch'
import type { TErrorResponce, TGetParams, TPostParams } from './types'
import { errorCheck } from './errors'
import { contains } from 'validator'
import { refresh } from './endpoints/auth/refresh'

export const accessTokenName = 'accessToken'
const emptyAccessToken = 'empty accessToken'

const fetchClient = createApiFetchClient({
  prefixUrl: `${import.meta.env.VITE_BACKEND_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
  fetchProps: { credentials: 'include' },
})

export async function POST<T>(params: TPostParams): Promise<T> {
  try {
    if (params.protectedPath) {
      const accToken = localStorage.getItem(accessTokenName)
      if (accToken) {
        const res = await fetchClient.post<unknown, TErrorResponce, TUnserializedBody>(
          params.path,
          params.body,
          {
            headers: {
              Authorization: `Bearer ${accToken}`,
            },
          },
        )

        errorCheck(res)

        return res.unwrap().data as T
      } else {
        throw new Error(emptyAccessToken)
      }
    } else {
      const res = await fetchClient.post<unknown, TErrorResponce, TUnserializedBody>(
        params.path,
        params.body,
      )
      errorCheck(res)

      return res.unwrap().data as T
    }
  } catch (err) {
    const error = err as Error
    if (contains(error.message, 'Token is expired')) {
      await refresh()
      return await POST<T>(params)
    } else {
      throw err
    }
  }
}

export async function GET<T>(params: TGetParams): Promise<T> {
  try {
    if (params.protectedPath) {
      const accToken = localStorage.getItem(accessTokenName)
      if (accToken) {
        const res = await fetchClient.get<unknown, TErrorResponce>(params.path, {
          headers: {
            Authorization: `Bearer ${accToken}`,
          },
        })

        errorCheck(res)

        return res.unwrap().data as T
      } else {
        throw new Error(emptyAccessToken)
      }
    } else {
      const res = await fetchClient.get<unknown, TErrorResponce>(params.path)
      errorCheck(res)

      return res.unwrap().data as T
    }
  } catch (err) {
    const error = err as Error
    if (contains(error.message, 'Token is expired')) {
      await refresh()
      return await GET<T>(params)
    } else {
      throw err
    }
  }
}
