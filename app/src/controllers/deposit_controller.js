import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  list
  products

  static get targets () {
    return [
      'bank', 'amount'
    ]
  }

  create (e) {
    e.preventDefault()
    if (this.bankTarget.value === '') {
      window.alert('Bank account is required')
      return
    }
    if (this.amountTarget.value === '') {
      window.alert('Amount is required')
      return
    }
    const that = this
    const req = { bank_id: this.bankTarget.value, amount: amountTarget.value }
    axios.post('/api/v1/accounting/deposits', req).then(resp => {
      console.log(resp)
      window.location.reload()
      that.cancel()
    }).catch(err => {
      let error = err.response.data.details
      window.alert(error)
    })
  }

  cancel () {
    this.bankTarget.value = ''
    this.amountTarget.value = ''
  }
}
