<script lang="ts" setup>
import { LitresLogo, ServiceButton } from '@/shared/ui'
</script>

<template>
  <ServiceButton
    button-type="button"
    disabled
  >
    <template #icon>
      <LitresLogo />
    </template>
    Litres
  </ServiceButton>
</template>
