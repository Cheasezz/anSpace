<script lang="ts" setup>
import { BaseButton, BaseInput, NegativeButton } from '@/shared/ui'
import styles from './styles.module.css'
import { ServicesList } from '@/widgets/servicesList'
import { accessTokenName } from '@/shared/api/client'
import { useUserStore } from '@/entities/user'

const userStore = useUserStore()

function logout() {
  localStorage.removeItem(accessTokenName)
  userStore.whoAmI()
}
</script>

<template>
  <div :class="[styles.almostThere]">
    <h1>Почти на месте</h1>
    <div :class="[styles.inputWrapper]">
      <BaseInput
        label-text="Никнейм"
        auto-complete="nickname"
        input-name="username"
        :tab-index="1"
        :with-error="true"
      />
      <p :class="[styles.hint]">Отображаемое имя на сайте</p>
    </div>

    <ServicesList />

    <div :class="[styles.buttons]">
      <BaseButton>Сохранить</BaseButton>
      <NegativeButton
        @click="logout"
        button-type="button"
      >
        Выход
      </NegativeButton>
    </div>
  </div>
</template>
