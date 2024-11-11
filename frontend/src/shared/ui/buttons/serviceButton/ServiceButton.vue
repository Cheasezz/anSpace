<script lang="ts" setup>
import { type ButtonHTMLAttributes, type PropType } from 'vue'
import styles from './styles.module.css'
import { BaseButton } from '../..'

defineProps({
  buttonType: {
    type: String as PropType<ButtonHTMLAttributes['type']>,
    default: 'submit',
  },

  serviceConnected: Boolean,
})
</script>

<template>
  <BaseButton
    :class="[styles.serviceButton, { [styles.connected]: serviceConnected }]"
    :button-type="buttonType"
  >
    <template #icon>
      <slot name="icon"> </slot>
    </template>
    <slot> </slot>
  </BaseButton>
</template>
