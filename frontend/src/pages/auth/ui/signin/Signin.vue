<script lang="ts" setup>
import styles from './styles.module.css'
import { BaseButton, BaseInput, NegativeButton, PassInput } from '@/shared/ui'
import { useValidateEmailAndPass } from '../../model/validations'
import { signin } from '@/shared/api'
import { useUserStore } from '@/entities/user'

const userStore = useUserStore()

const emit = defineEmits<{
  changeToSignup: [val: 'signup']
  changeToResetPass: [val: 'resetPass']
  asyncReqInProccess: [val: boolean]
}>()

const { validate, errEmail, errPass, resetErrVal } = useValidateEmailAndPass()

async function signinWithValidation(e: Event) {
  const tId = setTimeout(() => emit('asyncReqInProccess', true), 50)
  const auth = validate(e)
  if (auth) {
    try {
      await signin(auth)
      await userStore.whoAmI()
    } catch (err) {
      emit('asyncReqInProccess', false)
      const error = err as Error
      errPass.value = error.message
    }
  } else {
    clearTimeout(tId)
  }
}
</script>

<template>
  <form
    :class="[styles.form]"
    @submit.prevent="signinWithValidation"
  >
    <h1 :class="[styles.h1]">Авторизация</h1>
    <div :class="[styles.flexRow]">
      <p :class="[styles.p]">Еще не с нами?</p>
      <span
        :class="[styles.span, 'animated-underline']"
        tabindex="1"
        @click="$emit('changeToSignup', 'signup')"
        @keyup.enter="$emit('changeToSignup', 'signup')"
      >
        Создать аккаунт
      </span>
    </div>
    <BaseInput
      label-text="Почта"
      :tab-index="2"
      auto-complete="email"
      input-name="email"
      :with-error="true"
      :error-message="errEmail"
      @input="resetErrVal"
    />
    <PassInput
      label-text="Пароль"
      :tab-index="3"
      auto-complete="current-password"
      input-name="password"
      :with-error="true"
      :error-message="errPass"
      @input="resetErrVal"
    >
    </PassInput>
    <div :class="[styles.buttons]">
      <NegativeButton
        button-type="button"
        @click="$emit('changeToResetPass', 'resetPass')"
      >
        Забыли пароль?
      </NegativeButton>
      <BaseButton :class="[styles.button]"> Вход</BaseButton>
    </div>
  </form>
</template>
