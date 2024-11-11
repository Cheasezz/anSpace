<script lang="ts" setup>
import { computed, ref, type InputHTMLAttributes, type PropType } from 'vue'
import EyeClosed from '../../assets/svg/EyeClosed.vue'
import EyeOpen from '../../assets/svg/EyeOpen.vue'
import styles from './style.module.css'
import { BaseInput } from '../..'

defineProps({
  labelText: String,
  tabIndex: Number,
  autoComplete: {
    type: String as PropType<InputHTMLAttributes['autocomplete']>,
    default: 'off',
  },
  inputName: {
    type: String as PropType<InputHTMLAttributes['name']>,
  },
  withError: Boolean,
  errorMessage: String,
})

const passInputType = ref('password')

const passIsVisable = computed(() => {
  return passInputType.value === 'password' ? false : true
})

function changePassVisibility() {
  passIsVisable.value ? (passInputType.value = 'password') : (passInputType.value = 'text')
}
</script>

<template>
  <BaseInput
    :label-text="labelText"
    :input-type="passInputType"
    :with-error="true"
    :error-message="errorMessage"
    :auto-complete="autoComplete"
    :input-name="inputName"
    :tab-index="tabIndex"
  >
    <template #icon>
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
  </BaseInput>
</template>
