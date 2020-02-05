import dompurify from 'dompurify'

const Dygraph = require('../../dist/js/dygraphs.min.js')

export const appName = 'dcrextdata'

export const hide = (el) => {
  el.classList.add('d-none')
  el.classList.add('d-hide')
}

export const hideAll = (els) => {
  els.forEach(el => {
    el.classList.add('d-none')
    el.classList.add('d-hide')
  })
}

export const show = (el) => {
  el.classList.remove('d-none')
  el.classList.remove('d-hide')
}

export const showAll = (els) => {
  els.forEach(el => {
    el.classList.remove('d-none')
    el.classList.remove('d-hide')
  })
}

export const setAllValues = (targets, value) => {
  targets.forEach(el => {
    el.innerHTML = value
  })
}

export const showLoading = (loadingTarget, elementsToHide) => {
  show(loadingTarget)
  if (!elementsToHide || !elementsToHide.length) return
  elementsToHide.forEach(element => {
    hide(element)
  })
}

export const hideLoading = (loadingTarget, elementsToShow) => {
  hide(loadingTarget)
  if (!elementsToShow || !elementsToShow.length) return
  elementsToShow.forEach(element => {
    show(element)
  })
}

export const isHidden = (el) => {
  return el.classList.contains('d-none')
}

function darkenColor (colorStr) {
  // Defined in dygraph-utils.js
  var color = Dygraph.toRGB_(colorStr)
  color.r = Math.floor((255 + color.r) / 2)
  color.g = Math.floor((255 + color.g) / 2)
  color.b = Math.floor((255 + color.b) / 2)
  return 'rgb(' + color.r + ',' + color.g + ',' + color.b + ')'
}

export var options = {
  axes: { y: { axisLabelWidth: 100 } },
  axisLabelFontSize: 12,
  retainDateWindow: false,
  showRangeSelector: true,
  rangeSelectorHeight: 40,
  drawPoints: true,
  pointSize: 0.25,
  legend: 'always',
  labelsSeparateLines: true,
  highlightCircleSize: 4,
  yLabelWidth: 20,
  drawAxesAtZero: true
}

export function getRandomColor () {
  const letters = '0123456789ABCDEF'
  let color = '#'
  for (let i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * 16)]
  }
  return color
}

export function setActiveOptionBtn (opt, optTargets) {
  optTargets.forEach(li => {
    if (li.dataset.option === opt) {
      li.classList.add('active')
    } else {
      li.classList.remove('active')
    }
  })
}

export function setActiveRecordSetBtn (opt, optTargets) {
  optTargets.forEach(li => {
    if (li.dataset.option === opt) {
      li.classList.add('active')
    } else {
      li.classList.remove('active')
    }
  })
}

export function displayPillBtnOption (opt, optTargets) {
  optTargets.forEach(li => {
    if (opt === 'chart' && li.dataset.option === 'both') {
      li.classList.add('d-hide')
    } else {
      li.classList.remove('d-hide')
    }
  })
}

export function selectedOption (optTargets) {
  var key = false
  optTargets.forEach((el) => {
    if (el.classList.contains('active')) key = el.dataset.option
  })
  return key
}

export function insertQueryParam (name, value) {
  const urlParams = new URLSearchParams(window.location.search)
  const oldValue = urlParams.get(name)
  if (oldValue !== null) {
    return false
  }
  urlParams.append(name, value)
  const baseUrl = window.location.href.replace(window.location.search, '')
  window.history.pushState(window.history.state, appName, `${baseUrl}?${urlParams.toString()}`)
  return true
}

export function updateQueryParam (name, value) {
  let urlParams = new URLSearchParams(window.location.search)
  if (!urlParams.has(name)) {
    return false
  }
  urlParams.set(name, value)
  const baseUrl = window.location.href.replace(window.location.search, '')
  window.history.pushState(window.history.state, appName, `${baseUrl}?${urlParams.toString()}`)
  return true
}

export function insertOrUpdateQueryParam (name, value) {
  const urlParams = new URLSearchParams(window.location.search)
  return !urlParams.has(name) ? insertQueryParam(name, value) : updateQueryParam(name, value)
}

export function getParameterByName (name, url) {
  const urlParams = new URLSearchParams(window.location.search)
  return urlParams.get(name)
}

export function updateZoomSelector (targets, minDate, maxDate) {
  const duration = maxDate - minDate
  const days = duration / (1000 * 60 * 60 * 24)
  targets.forEach(el => {
    let showElement = false
    switch (el.dataset.option) {
      case 'day':
      case 'all':
        showElement = days >= 1
        break
      case 'week':
        showElement = days >= 7
        break
      case 'month':
        showElement = days >= 30
        break
      case 'year':
        showElement = days >= 365
        break
    }

    if (showElement) {
      show(el)
    } else {
      hide(el)
    }
  })
}

export function formatDate (date, format) {
  if (!format || format === '') {
    format = 'yyyy-MM-dd hh:mm'
  }

  let dd = date.getDate()
  let mm = date.getMonth() + 1
  let yyyy = date.getFullYear()
  let milliseconds = date.getMilliseconds()
  let seconds = date.getSeconds()
  let minutes = date.getMinutes()
  let hour = date.getHours()

  if (dd < 10) {
    dd = '0' + dd
  }

  if (mm < 10) {
    mm = '0' + mm
  }

  if (hour < 10) {
    hour = '0' + hour
  }

  if (minutes < 10) {
    minutes = '0' + minutes
  }

  if (seconds < 10) {
    seconds = '0' + seconds
  }

  let dateFormatted = format.replace('yyyy', yyyy).replace('MM', mm).replace('dd', dd)
  dateFormatted = dateFormatted.replace('hh', hour).replace('mm', minutes)
  dateFormatted = dateFormatted.replace('ss', seconds).replace('sss', milliseconds)
  return dateFormatted
}

export function getNumberOfPages (recordsCount, pageSize) {
  const rem = recordsCount % pageSize
  let pageCount = (recordsCount - rem) / pageSize
  if (rem > 0) {
    pageCount += 1
  }
  return pageCount
}
