import { ref } from 'vue'

export const notifyState = ref({
  show: false,
  message: '',
  type: 'error',
})

export const notify = (message, type = 'error') => {
  notifyState.value = { show: true, message, type }
  setTimeout(() => {
    notifyState.value.show = false
  }, 2500)
}
