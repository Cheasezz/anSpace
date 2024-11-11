import { me, type TUserResponce } from '@/shared/api'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const user = ref<TUserResponce>()
  const whoAmIError = ref<Error>()

  async function whoAmI() {
    try {
      user.value = await me()
    } catch (err) {
      user.value = undefined
      whoAmIError.value = err as Error
    }
  }

  return { user, whoAmIError, whoAmI }
})
