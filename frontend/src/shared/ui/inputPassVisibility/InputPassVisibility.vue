<script lang="ts" setup>
import { computed } from 'vue'
import EyeClosed from '../assets/svg/EyeClosed.vue'
import EyeOpen from '../assets/svg/EyeOpen.vue'
import styles from './style.module.css'

const prop = defineProps<{
  inputType?: string
}>()

const emit = defineEmits<{
  newInputType: [val: 'password' | 'text']
}>()

const passIsVisable = computed(() => {
  return prop.inputType === 'password' ? false : true
})

function changePassVisibility() {
  passIsVisable.value ? emit('newInputType', 'password') : emit('newInputType', 'text')
}
</script>

<template>
  <div
    :class="[styles.div]"
    @click="changePassVisibility"
  >
    <EyeClosed
      v-show="passIsVisable"
      :class="[styles.svg, styles.close]"
    />
    <EyeOpen
      v-show="!passIsVisable"
      :class="[styles.svg, styles.open]"
    />
  </div>
</template>
