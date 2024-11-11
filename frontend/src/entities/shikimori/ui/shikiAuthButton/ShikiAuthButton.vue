<script lang="ts" setup>
import { ServiceButton, ShikiIcon } from '@/shared/ui'
import { ref } from 'vue'

const isAuth = ref(false)
function shikiAuth() {
  isAuth.value = !isAuth.value
}
</script>

<template>
  <ServiceButton
    @click="shikiAuth"
    :service-connected="isAuth"
    button-type="button"
  >
    <template #icon>
      <ShikiIcon />
    </template>
    Shikimori
  </ServiceButton>
</template>
