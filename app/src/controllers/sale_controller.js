import { Controller } from 'stimulus'
import { hide, show } from '../utils'
import _ from 'lodash-es'
import axios from 'axios'

export default class extends Controller {
  list
  products

  static get targets () {
    return [
      'barcodeInput', 'productSelect', 'quantityInput', 'addToListBtn', 'cartItemDiv', 'listTbl', 'itemTemplate',
      'cartTotal', 'customerName', 'phoneNumber', 'amountTender', 'paymentMethod', 'accountNumber',
      'accountNumberDiv', 'amountTenderDiv'
    ]
  }

  connect () {
    this.list = []
    this.products = []
    const that = this
    Array.prototype.forEach.call(this.productSelectTarget.options, function (opt) {
      if (opt.value === '') return
      that.products.push({
        barcode: opt.value,
        name: opt.getAttribute('data-name'),
        id: opt.getAttribute('data-id'),
        price: parseFloat(opt.getAttribute('data-price'))
      })
    })
  }

  changeProduct (evt) {
    const product = this.findProduct(evt.currentTarget.value)
    this.barcodeInputTarget.value = product.barcode
  }

  findProduct (barcode) {
    let product = null
    this.products.forEach(p => {
      if (p.barcode === barcode) {
        product = p
      }
    })

    return product
  }

  addToList () {
    let barcode = this.barcodeInputTarget.value
    if (barcode === '') {
      barcode = this.productSelectTarget.value
    } else {
      this.productSelectTarget.value = barcode
    }

    if (barcode === '') {
      window.alert('Please scan in or select a product')
      this.barcodeInputTarget.focus()
      return
    }

    let qnt = parseInt(this.quantityInputTarget.value)
    if (qnt < 1) {
      qnt = 1
    }

    for (let i = 0; i < this.list.length; i++) {
      if (this.list[i].barcode === barcode) {
        this.list[i].quantity += qnt
        this.displayList()
        return
      }
    }

    let p = this.findProduct(barcode)
    this.list.push({
      id: p.id,
      barcode: p.barcode,
      unitPrice: p.price,
      name: p.name,
      quantity: qnt
    })

    this.displayList()
  }

  barcodeEntered (evt) {
    if (evt.keyCode !== 13) {
      return
    }
    this.addToList()
  }

  removeFromList (evt) {
    let barcode = evt.currentTarget.getAttribute('data-barcode')
    _.remove(this.list, function (item) {
      return item.barcode === barcode
    })
    this.displayList()
  }

  displayList () {
    const _this = this
    this.listTblTarget.innerHTML = ''
    let cartTotal = 0

    this.list.forEach((item, i) => {
      const exRow = document.importNode(_this.itemTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerText = i + 1
      fields[1].innerText = item.name
      fields[2].innerText = item.barcode
      fields[3].innerHTML = item.quantity
      fields[4].innerHTML = item.unitPrice
      fields[5].innerHTML = (item.quantity * item.unitPrice)
      fields[6].innerHTML = `<button data-action="click->sale#removeFromList" data-barcode="${item.barcode}">Remove</button>`

      _this.listTblTarget.appendChild(exRow)
      cartTotal += item.quantity * item.unitPrice
    })

    this.cartTotalTarget.textContent = cartTotal
    this.amountTenderTarget.value = cartTotal
    if (this.list.length > 0) {
      show(this.cartItemDivTarget)
    } else {
      hide(this.cartItemDivTarget)
    }
    this.barcodeInputTarget.focus()
  }

  paymentMethodChanged (evt) {
    if (this.paymentMethodTarget.value === 'wallet') {
      show(this.accountNumberDivTarget)
      hide(this.amountTenderDivTarget)
    } else {
      hide(this.accountNumberDivTarget)
      show(this.amountTenderDivTarget)
    }
  }

  sell () {
    const amountTender = parseFloat(this.amountTenderTarget.value)
    if (amountTender < this.cartTotal) {
      window.alert('The amount tender cannot be less than the cart total')
      return
    }
    let req = {
      payment_method: this.paymentMethodTarget.value,
      account_number: this.accountNumberTarget.value,
      amount_tender: amountTender,
      customer_name: this.customerNameTarget.value,
      phone_number: this.phoneNumberTarget.value,
      items: []
    }
    this.list.forEach(item => {
      req.items.push({
        product_id: item.id,
        quantity: item.quantity
      })
    })

    const that = this

    axios.post('/api/v1/sales/sell', req).then(resp => {
      console.log(resp)
      // $('#receiptModal').modal()
      // todo: open receipt page in a new tap
      // todo: show notification
      window.location.href = `/sales/${resp.data.id}`
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
