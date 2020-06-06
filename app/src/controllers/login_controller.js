import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  list
  products

  static get targets () {
    return [
      'rememberMe', 'email', 'password'
    ]
  }

  submit (e) {
    e.preventDefault()
    window.alert('hi')
    if (this.emailTarget.value === '') {
      window.alert('Email/Username is required')
      return
    }
    if (this.passwordTarget.value === '') {
      window.alert('Password is required')
      return
    }
    const that = this
    const req = {
      email: this.emailTarget.value,
      password: this.passwordTarget.value,
      rememberMe: this.rememberMeTarget.checked
    }
    axios.post('', req).then(resp => {
      console.log(resp)
      window.location.href = resp.data.redirect
      that.cancel()
    }).catch(err => {
      let error = err.response.data.details
      window.alert(error)
    })
  }

  cancel () {
    this.emailTarget.value = ''
    this.passwordTarget.value = ''
    this.rememberMeTarget.checked = false
  }
}
