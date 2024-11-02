<script lang="ts" setup>
import type { ButtonHTMLAttributes, PropType } from 'vue'
import styles from './styles.module.css'

defineProps({
  buttonType: {
    type: String as PropType<ButtonHTMLAttributes['type']>,
    default: 'submit',
  },
})
</script>

<template>
  <button
    :class="[styles.button]"
    :type="buttonType"
  >
    <slot></slot>
  </button>
</template>
