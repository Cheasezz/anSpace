<script lang="ts" setup>
import type { InputHTMLAttributes, InputTypeHTMLAttribute, PropType } from 'vue'
import styles from './styles.module.css'

defineProps({
  labelText: String,
  inputType: {
    type: String as PropType<InputTypeHTMLAttribute>,
    default: 'text',
  },
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
        :class="[
          styles.input,
          {
            [styles.err]: errorMessage,
          },
        ]"
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
