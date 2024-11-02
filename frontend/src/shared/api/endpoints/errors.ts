/* eslint-disable no-console */
import {
  FetchError,
  isStatusCode,
  NetworkError,
  RequestError,
  type TFetchResponse,
} from 'feature-fetch'
import type { TErrorResponce } from '../types'

export function errorCheck(res: TFetchResponse<unknown, TErrorResponce, 'json'>) {
  if (res.isErr()) {
    const error = res.error

    if (isStatusCode(res, 404)) {
      console.error('Not found:', error.message)
    }

    if (error instanceof NetworkError) {
      console.error('Network error:', error.message)
      throw new Error(`Сетевая ошибка: ${error.message}`)
    } else if (error instanceof RequestError) {
      console.error('Request error:', error.message, 'Status:', error.status)

      throw new Error(`Ошибка запроса: ${(error.data as TErrorResponce).message}`)
    } else if (error instanceof FetchError) {
      console.error('Service error:', error.message)
      throw new Error(`Ошибка сервера: ${error.message}`)
    } else {
      console.error('Unexpected error:', error)
      throw new Error(`Неизвестная ошибка: ${error}`)
    }
  }
}
