import { createPinia } from 'pinia'

export const useAuthStore = createPinia({
  id: 'auth',
  state: () => ({
    token: localStorage.getItem('token') || '',
    user_id: localStorage.getItem('user_id') || ''
  }),
  actions: {
    setToken(token, user_id) {
      this.token = token
      this.user_id = user_id
      localStorage.setItem('token', token)
      localStorage.setItem('user_id', user_id)
    },
    logout() {
      this.token = ''
      this.user_id = ''
      localStorage.removeItem('token')
      localStorage.removeItem('user_id')
    }
  }
})
