<script lang="ts" setup>
import styles from './styles.module.css'
import { BaseButton, BaseInput, InputPassVisibility } from '@/shared/ui'
import { signup } from '@/shared/api'
import { useValidateEmailAndPass } from '../../model/validations'
import { ref } from 'vue'

defineEmits<{
  changeToSignin: [val: 'signin']
}>()

// const { passIsVisible, passInputType, repeatPassIsVisible, repeatPassInputType } =
//   useInputPassVisibility()
const passInputType = ref('password')
const repeatPassInputType = ref('password')

const { validate, errEmail, errPass, errRepPass, resetErrVal } = useValidateEmailAndPass()

async function signupWithValidation(e: Event) {
  const auth = validate(e)
  if (auth) {
    try {
      await signup(auth)
    } catch (err) {
      const error = err as Error
      errRepPass.value = error.message
    }
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
    <BaseInput
      label-text="Пароль"
      :input-type="passInputType"
      :with-error="true"
      :error-message="errPass"
      auto-complete="new-password"
      input-name="password"
      :tab-index="3"
      @input="resetErrVal"
    >
      <template #icon>
        <InputPassVisibility
          :input-type="passInputType"
          @new-input-type="(val) => (passInputType = val)"
        />
      </template>
    </BaseInput>
    <BaseInput
      label-text="Подтвердите пароль"
      :input-type="repeatPassInputType"
      :with-error="true"
      :error-message="errRepPass"
      auto-complete="new-password"
      input-name="repeatPassword"
      :tab-index="4"
      @input="resetErrVal"
    >
      <template #icon>
        <InputPassVisibility
          :input-type="repeatPassInputType"
          @new-input-type="(val) => (repeatPassInputType = val)"
        />
      </template>
    </BaseInput>
    <BaseButton :class="[styles.button]">Регистрация</BaseButton>
  </form>
</template>
