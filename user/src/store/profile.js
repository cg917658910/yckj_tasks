import { ref } from 'vue'
import { fetchProfile } from '../api/user'

export const profileState = ref({})

export const loadProfile = async () => {
  try {
    const res = await fetchProfile()
    profileState.value = res.data || {}
  } catch (err) {
    // ignore
  }
}
