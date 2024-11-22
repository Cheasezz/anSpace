<script lang="ts" setup>
import { BaseButton, BaseInput, NegativeButton } from '@/shared/ui'
import styles from './styles.module.css'

const emit = defineEmits<{
  changeToSignin: [val: 'signin']
  changeToSignup: [val: 'signup']
  asyncReqInProccess: [val: boolean]
}>()
</script>

<template>
<div>
  <h1 :class="[styles.h1]">Восстановить пароль</h1>
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
  />
  <div :class="[styles.buttons]">
    <NegativeButton
      button-type="button"
      @click="$emit('changeToSignin', 'signin')"
    >
      Войти в аккаунт
    </NegativeButton>
    <BaseButton>Сохранить</BaseButton>
  </div>
</div>
</template>
