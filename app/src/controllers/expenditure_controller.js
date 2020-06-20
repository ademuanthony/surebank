import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  list
  products

  static get targets () {
    return [
      'amount', 'memo'
    ]
  }

  create (e) {
    e.preventDefault()
    if (this.amountTarget.value === '') {
      window.alert('Amount is required')
      return
    }
    const that = this
    const req = { amount: parseFloat(this.amountTarget.value), memo: this.memoTarget.value }
    axios.post('/api/v1/accounting/expenditures', req).then(resp => {
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
