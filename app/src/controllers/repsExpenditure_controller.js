import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  list
  products

  static get targets () {
    return [
      'amount', 'salesRep', 'reason'
    ]
  }

  create (e) {
    e.preventDefault()
    if (this.amountTarget.value === '') {
      window.alert('Amount is required')
      return
    }
    if (this.salesRepTarget.value === '') {
      window.alert('Amount is required')
      return
    }
    if (this.reasonTarget.value === '') {
      window.alert('Amount is required')
      return
    }
    const req = {
      amount: parseFloat(this.amountTarget.value),
      sales_rep_phone_number: this.salesRepTarget.value,
      reason: this.reasonTarget.value
    }
    axios.post('/api/v1/accounting/reps-expenditures', req).then(resp => {
      window.location.reload()
    }).catch(err => {
      let error = err.response.data.details
      window.alert(error)
    })
  }

  remove (e) {
    let confirmed = window.confirm('This record will be permanently delete. Continue?')
    if (!confirmed) {
      return
    }
    const id = e.currentTarget.getAttribute('data-id')
    axios.delete('/api/v1/accounting/reps-expenditures/' + id).then(resp => {
      window.location.reload()
    }).catch(err => {
      let error = err.response.data.details
      console.log(err)
      window.alert(error)
    })
  }
}
