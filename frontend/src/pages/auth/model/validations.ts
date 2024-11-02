import { ref } from 'vue'
import isEmail from 'validator/es/lib/isEmail'
import type { TUserAuth } from '@/shared/api'
import { equals } from 'validator'

const errEmptyFields = 'Заполните все поля',
  errShortPass = 'Парольдолжен быть больше 11 символов',
  errIncorrectEmail = 'Некорректная почта',
  errPassMustMatch = 'Пароли должны совпадать'

export function useValidateEmailAndPass() {
  const errEmail = ref(),
    errPass = ref(),
    errRepPass = ref()

  // Accept event from submited form for authorization and validate named input
  // Relevant for signin form and signup form with named input "repeatPassword"
  function validate(e: Event): TUserAuth | false {
    const form = new FormData(e.target as HTMLFormElement),
      isSignup = form.has('repeatPassword'),
      email = form.get('email')?.toString(),
      password = form.get('password')?.toString(),
      repeatPassword = form.get('repeatPassword')?.toString()

    if (!password || !email) {
      if (!password) errPass.value = errEmptyFields
      if (!email) errEmail.value = errEmptyFields
      return false
    }
    if (password.length < 12) {
      errPass.value = errShortPass
      return false
    }
    if (!isEmail(email, { domain_specific_validation: true })) {
      errEmail.value = errIncorrectEmail
      return false
    }
    if (isSignup && typeof repeatPassword !== 'undefined') {
      if (!equals(password, repeatPassword)) {
        errRepPass.value = errPassMustMatch
        return false
      }
    }

    const user: TUserAuth = {
      email,
      password,
    }
    return user
  }

  function resetErrVal() {
    errEmail.value = ''
    errPass.value = ''
    errRepPass.value = ''
  }

  return { errEmail, errPass, errRepPass, validate, resetErrVal }
}
