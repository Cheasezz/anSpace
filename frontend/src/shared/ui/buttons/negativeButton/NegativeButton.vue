<script lang="ts" setup>
import type { ButtonHTMLAttributes, PropType } from 'vue'
import { BaseButton } from '../..'
import styles from './styles.module.css'

defineProps({
  buttonType: {
    type: String as PropType<ButtonHTMLAttributes['type']>,
    default: 'submit',
  },
})
</script>

<template>
  <BaseButton
    :class="[styles.negativeButton]"
    :button-type="buttonType"
  >
    <template #icon>
      <slot name="icon"></slot>
    </template>
    <slot></slot>
  </BaseButton>
</template>
