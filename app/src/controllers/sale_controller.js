import { Controller } from 'stimulus'
import { show } from '../utils'

export default class extends Controller {
  list
  products

  static get targets () {
    return [
      'barcodeInput', 'productSelect', 'quantityInput', 'addToListBtn', 'cartItemDiv', 'listTbl', 'itemTemplate', 'cartTotal'
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

    let p = this.findProduct(barcode)
    this.list.push({
      id: p.id,
      barcode: p.barcode,
      unitPrice: p.price,
      name: p.name,
      quantity: qnt
    })

    this.displayList()
    show(this.cartItemDivTarget)
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

      _this.listTblTarget.appendChild(exRow)
      cartTotal += item.quantity * item.unitPrice
    })

    this.cartTotalTarget.textContent = cartTotal
  }
}
