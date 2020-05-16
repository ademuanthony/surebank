import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  list
  products

  static get targets () {
    return [
      'bank', 'name', 'number'
    ]
  }

  create (e) {
    e.preventDefault()
    if (this.bankTarget.value === '') {
      window.alert('Bank name is required')
      return
    }
    if (this.nameTarget.value === '') {
      window.alert('Account name is required')
      return
    }
    if (this.numberTarget.value === '') {
      window.alert('Account number is required')
      return
    }
    const that = this
    const req = { name: this.nameTarget.value, number: this.numberTarget.value, bank: this.bankTarget.value }
    axios.post('/api/v1/accounting/banks', req).then(resp => {
      console.log(resp)
      window.location.reload()
      that.cancel()
    }).catch(err => {
      let error = err.response.data.details
      window.alert(error)
    })
  }

  cancel () {
    this.nameTarget.value = ''
    this.bankTarget.value = ''
    this.numberTarget.value = ''
  }
}
