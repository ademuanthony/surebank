import { Controller } from 'stimulus'
import { hide, show } from '../utils'
import _ from 'lodash-es'
import axios from 'axios'

export default class extends Controller {
  list
  products

  static get targets () {
    return [
      'nameInput'
    ]
  }

  sell () {
    if (this.nameInputTarget.value === '') {
      window.alert('Branch name is required')
      return
    }
    const that = this
    const req = {name: this.nameInputTarget.value}
    axios.post('/api/v1/branches', req).then(resp => {
      console.log(resp)
      window.location.reload()
      that.cancel()
    }).catch(err => {
      let error = err.response.data.details
      window.alert(error)
    })
  }

  cancel () {
    this.list = []
    this.barcodeInputTarget.value = ''
    this.productSelectTarget.value = ''
    this.quantityInputTarget.value = 1
    this.paymentMethodTarget.value = 'cash'
    this.accountNumberTarget.value = ''
    hide(this.accountNumberDivTarget)
    show(this.amountTenderDivTarget)
    this.displayList()
  }
}
