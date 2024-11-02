<script lang="ts" setup>
import styles from './styles.module.css'

const { inputType = 'text' } = defineProps<{
  labelText?: string
  inputType?: string
  tabIndex?: number
  autoComplete?: string
  inputName?: string
  withError?: boolean
  errorMessage?: string
}>()

const model = defineModel<string>()

// Disable text auto select when input focusing with tab or in firefox brows
function disableAutoSelectedText(input: HTMLInputElement) {
  input.selectionStart = input.selectionEnd
}
</script>

<template>
  <div :class="[styles.inputWrapper]">
    <label
      :class="[styles.label]"
      tabindex="-1"
    >
      {{ labelText }}
      <input
        v-model.trim="model"
        :tabindex="tabIndex"
        :class="[styles.input]"
        :type="inputType"
        :autocomplete="autoComplete"
        :name="inputName"
        @focusin="disableAutoSelectedText($event.target as HTMLInputElement)"
      />
      <slot name="icon"> </slot>
    </label>
    <p
      v-if="withError"
      :class="[styles.inputError]"
    >
      {{ errorMessage }}
    </p>
  </div>
</template>
