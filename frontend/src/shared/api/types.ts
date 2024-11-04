import type { TUnserializedBody } from 'feature-fetch'

export type TAccessTokenResponce = {
  accessToken: string
}

export type TUserResponce = {
  user: {
    email: string
    username: string
    passwordHash: string
  }
}

export type TErrorResponce = {
  message: string
}

export type TUserAuth = {
  email: string
  password: string
}

export type TPostParams = {
  path: string
  body?: TUnserializedBody
  protectedPath?: boolean
}
export type TGetParams = {
  path: string
  protectedPath?: boolean
}
