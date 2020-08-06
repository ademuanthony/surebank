import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  list
  products

  static get targets () {
    return [
      'bank', 'amount', 'accountName'
    ]
  }

  accountNumberChanged (e) {
    const number = e.currentTarget.value
    if (number.length !== 7) return
    axios.get('/api/v1/customers/account-name?account_number=' + number).then(resp => {
      this.accountNameTarget.textContent = resp.data.name
    }).catch(err => {
      this.accountNameTarget.textContent = ''
      let error = err.response.data.details
      window.alert(error)
    })
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
    const req = { bank_id: this.bankTarget.value, amount: parseFloat(this.amountTarget.value) }
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
