<script lang="ts" setup>
import styles from './styles.module.css'
import { BaseButton, BaseInput, InputPassVisibility } from '@/shared/ui'
import { useValidateEmailAndPass } from '../../model/validations'
import { signin } from '@/shared/api'
import { ref } from 'vue'

defineEmits<{
  changeToSignup: [val: 'signup']
}>()

const passInputType = ref('password')
const { validate, errEmail, errPass, resetErrVal } = useValidateEmailAndPass()

async function signinWithValidation(e: Event) {
  const auth = validate(e)
  if (auth) {
    try {
      await signin(auth)
    } catch (err) {
      const error = err as Error
      errPass.value = error.message
    }
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
    <BaseInput
      label-text="Пароль"
      :input-type="passInputType"
      :tab-index="3"
      auto-complete="current-password"
      input-name="password"
      :with-error="true"
      :error-message="errPass"
      @input="resetErrVal"
    >
      <template #icon>
        <InputPassVisibility
          :input-type="passInputType"
          @new-input-type="(val) => (passInputType = val)"
        />
      </template>
    </BaseInput>
    <BaseButton :class="[styles.button]"> Вход</BaseButton>
  </form>
</template>
