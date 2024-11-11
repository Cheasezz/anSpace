<script lang="ts" setup>
import styles from './styles.module.css'
import { ShikiAuthButton } from '@/entities/shikimori'
import { KPAuthButton } from '@/entities/kinopoisk'
import { LitresAuthButton } from '@/entities/litres'
</script>

<template>
  <div :class="[styles.servicesList]">
    <h3 :class="[styles.h3]">Подключить список</h3>
    <div :class="[styles.list]">
      <ShikiAuthButton />
      <KPAuthButton />
      <LitresAuthButton />
    </div>
    <p :class="[styles.hint]">
      На данный момент к подключению доступен только Shikimori. Больше сервисов появится в будущем.
    </p>
  </div>
</template>
