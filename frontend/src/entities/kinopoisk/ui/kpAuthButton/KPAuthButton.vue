<script lang="ts" setup>
import { KPIcon, ServiceButton } from '@/shared/ui'
</script>

<template>
  <ServiceButton
    button-type="button"
    disabled
  >
    <template #icon>
      <KPIcon />
    </template>
    Kinopoisk
  </ServiceButton>
</template>
