<script lang="ts" setup>
import styles from './styles.module.css'
import { ref } from 'vue'
import Signin from '../signin/Signin.vue'
import Signup from '../signup/Signup.vue'
import { MainWithSpaceBg } from '@/shared/ui'

const inOrUp = ref<string>('signin')
</script>

<template>
  <MainWithSpaceBg>
    <div :class="[styles.container]">
      <Transition mode="out-in">
        <Signin
          v-if="inOrUp == 'signin'"
          @change-to-signup="(val) => (inOrUp = val)"
        />

        <Signup
          v-else-if="inOrUp == 'signup'"
          @change-to-signin="(val) => (inOrUp = val)"
        />
      </Transition>
    </div>
  </MainWithSpaceBg>
</template>

<style>
.v-move,
.v-enter-active,
.v-leave-active {
  transition: all 0.3s linear;
}

.v-enter-from {
  opacity: 0;
  transform: translateX(-40px);
}
.v-enter-to {
  opacity: 1;
}

.v-leave-to {
  opacity: 0;
  transform: translateX(40px);
}
</style>
