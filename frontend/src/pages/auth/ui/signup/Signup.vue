<script lang="ts" setup>
import styles from './styles.module.css'
import { BaseButton, BaseInput, PassInput } from '@/shared/ui'
import { signup } from '@/shared/api'
import { useValidateEmailAndPass } from '../../model/validations'
import { useUserStore } from '@/entities/user'

const emit = defineEmits<{
  changeToSignin: [val: 'signin']
  asyncReqInProccess: [val: boolean]
}>()

const userStore = useUserStore()

const { validate, errEmail, errPass, errRepPass, resetErrVal } = useValidateEmailAndPass()

async function signupWithValidation(e: Event) {
  const tId = setTimeout(() => emit('asyncReqInProccess', true), 2000)
  const auth = validate(e)
  if (auth) {
    try {
      await signup(auth)
      await userStore.whoAmI()
    } catch (err) {
      const error = err as Error
      console.log(err)
      errRepPass.value = error.message
    }
  } else {
    clearTimeout(tId)
  }
}
</script>

<template>
  <form
    :class="[styles.form]"
    @submit.prevent="signupWithValidation"
  >
    <h1 :class="[styles.h1]">Регистрация</h1>
    <div :class="[styles.flexRow]">
      <p :class="[styles.p]">Уже есть аккаунт?</p>
      <span
        :class="[styles.span, 'animated-underline']"
        tabindex="1"
        @click="$emit('changeToSignin', 'signin')"
        @keyup.enter="$emit('changeToSignin', 'signin')"
      >
        Войти
      </span>
    </div>
    <BaseInput
      label-text="Почта"
      :with-error="true"
      :error-message="errEmail"
      auto-complete="email"
      input-name="email"
      :tab-index="2"
      @input="resetErrVal"
    />
    <PassInput
      label-text="Пароль"
      :with-error="true"
      :error-message="errPass"
      auto-complete="new-password"
      input-name="password"
      :tab-index="3"
      @input="resetErrVal"
    />
    <PassInput
      label-text="Подтвердите пароль"
      :with-error="true"
      :error-message="errRepPass"
      auto-complete="new-password"
      input-name="repeatPassword"
      :tab-index="4"
      @input="resetErrVal"
    />
    <BaseButton :class="[styles.button]">Регистрация</BaseButton>
  </form>
</template>
