<script lang="ts" setup>
import styles from './styles.module.css'
import { ref, watch } from 'vue'
import Signin from '../signin/Signin.vue'
import Signup from '../signup/Signup.vue'
import { MainWithSpaceBg } from '@/shared/ui'
import AlmostThere from '../almostThere/AlmostThere.vue'
import { useUserStore } from '@/entities/user'
import ResetPassword from '../resetPassword/ResetPassword.vue'

const isLoading = ref<boolean>(false)
const userStore = useUserStore()
const authProcess = ref<string>('signin')
await userStore.whoAmI()

watch(
  () => userStore.user,
  () => {
    if (userStore.user) authProcess.value = 'almost'
    else authProcess.value = 'signin'
    isLoading.value = false
  },
  { immediate: true },
)

function changeLoadVal(v: boolean) {
  v ? (isLoading.value = true) : (isLoading.value = false)
}
</script>

<template>
  <MainWithSpaceBg>
    <div :class="[styles.container, { [styles.active]: isLoading }]">
      <Transition mode="out-in">
        <AlmostThere v-if="authProcess == 'almost'" />
        <ResetPassword
          v-else-if="authProcess == 'resetPass'"
          @change-to-signin="(val) => (authProcess = val)"
          @change-to-signup="(val) => (authProcess = val)"
        />
        <Signin
          v-else-if="authProcess == 'signin'"
          @change-to-signup="(val) => (authProcess = val)"
          @change-to-reset-pass="(val) => (authProcess = val)"
          @async-req-in-proccess="changeLoadVal"
        />

        <Signup
          v-else-if="authProcess == 'signup'"
          @change-to-signin="(val) => (authProcess = val)"
          @async-req-in-proccess="changeLoadVal"
        />
      </Transition>
    </div>
  </MainWithSpaceBg>
</template>

<style>
.v-move,
.v-enter-active,
.v-leave-active {
  transition: all 0.2s linear;
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
